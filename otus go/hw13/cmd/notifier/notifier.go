package main

import (
	"encoding/json"
	"flag"
	"fmt"
	amqpClient "github.com/lenniDespero/otus-golang/hw13/internal/pkg/ampq"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type notifier struct {
	ampq       *amqpClient.Ampq
	noticeAmpq *amqpClient.Ampq
	stat       *prometheus.CounterVec
}

func newNotifier(ampq *amqpClient.Ampq, notice *amqpClient.Ampq) *notifier {
	return &notifier{ampq, notice, monitor.NewCounterVec("calendar_notifier", "sender", "Sending notice")}
}

func (n *notifier) Start() error {
	err := n.ampq.Subscribe("notifier", func(delivery amqp.Delivery) {
		event := &models.Event{}
		if err := json.Unmarshal(delivery.Body, event); err != nil {
			logger.Error(fmt.Sprintf("Failed to parse message: %s", err.Error()))
			return
		}
		n.notify(event)
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *notifier) notify(event *models.Event) {
	fmt.Printf("Get event: %v", event)
	msg, err := json.Marshal(event)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to encode event: %s", err.Error()))
	}
	if err := n.noticeAmpq.Publish(msg); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to publish event: %s", err.Error()))
	}
	logger.Debug(fmt.Sprintf("message %s published to notice", msg))
	n.stat.WithLabelValues("count").Inc()
}

func main() {
	var configPath = flag.String("config", "../config/application.yml", "path to configuration flag")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	conf := config.GetConfigFromFile(*configPath)
	logger.Init(conf.Log.LogLevel, conf.Log.LogFile)
	amqpBus, err := amqpClient.NewAmpq(&conf.Ampq)
	if err != nil {
		logger.Fatal(err.Error())
	}
	noticeConf := config.Ampq{
		Host:     conf.Ampq.Host,
		Port:     conf.Ampq.Port,
		User:     conf.Ampq.User,
		Password: conf.Ampq.Password,
		Queue:    "notice",
	}
	noticeBus, err := amqpClient.NewAmpq(&noticeConf)
	if err != nil {
		logger.Fatal(err.Error())
	}

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 2114),
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatal("Unable to start a http server", "error", err)
		}
	}()
	notifierAgent := newNotifier(amqpBus, noticeBus)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Scheduler init error: %s", err.Error()))
	}
	notifierAgent.Start()
}
