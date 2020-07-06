package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	amqpClient "github.com/lenniDespero/otus-golang/hw13/internal/pkg/ampq"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/storage/sql"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/types"
	"github.com/spf13/pflag"
	"log"
	"strconv"
	"time"
)

type scheduler struct {
	storage types.LimitedStorageInterface
	ampq    *amqpClient.Ampq
	ctx     context.Context
}

func newScheduler(storage types.LimitedStorageInterface, ampq *amqpClient.Ampq, ctx context.Context) *scheduler {
	return &scheduler{storage, ampq, ctx}
}

func (s *scheduler) Start(conf config.Scheduler) {
	period, err := strconv.Atoi(conf.Period)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug(fmt.Sprintf("Set ticker on seconds : %d", period))
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	logger.Debug(fmt.Sprintf("Started at : %v", time.Now()))
	defer ticker.Stop()
	localCtx, cancel := context.WithCancel(s.ctx)
	for ; true; <-ticker.C {
		select {
		case <-localCtx.Done():
			return
		default:
			logger.Debug("Get current events")

			events, err := s.storage.GetEventsByStartPeriod(conf.BeforeTime, conf.EventTime, s.ctx)
			if err != nil {
				logger.Fatal(err.Error())
				cancel()
			}
			for _, event := range events {
				msg, err := json.Marshal(event)
				if err != nil {
					logger.Fatal(fmt.Sprintf("Failed to encode event: %s", err.Error()))
					cancel()
					continue
				}
				if err := s.ampq.Publish(msg); err != nil {
					logger.Fatal(fmt.Sprintf("Failed to publish event: %s", err.Error()))
					cancel()
					continue
				}
				logger.Debug(fmt.Sprintf("message %s published", msg))
			}
			continue
		}
	}
}

func main() {
	var configPath = flag.String("config", "../config/application.yml", "path to configuration flag")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	conf := config.GetConfigFromFile(*configPath)
	logger.Init(conf.Log.LogLevel, conf.Log.LogFile)
	storage, err := sql.New(&conf.DBConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	amqpBus, err := amqpClient.NewAmpq(&conf.Ampq)
	if err != nil {
		logger.Fatal(err.Error())
	}
	scheduler := newScheduler(storage, amqpBus, context.Background())
	if err != nil {
		log.Fatalf("Scheduler init error: %s", err)
	}
	scheduler.Start(conf.Scheduler)
}
