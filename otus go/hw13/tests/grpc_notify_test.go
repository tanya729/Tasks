package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cucumber/godog"
	"github.com/golang/protobuf/ptypes"
	pkg "github.com/lenniDespero/otus-golang/hw13/internal/pkg"
	amqpClient "github.com/lenniDespero/otus-golang/hw13/internal/pkg/ampq"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type grpcNotifyEvent struct {
	err    []error
	events []*models.Event
	ampq   *amqpClient.Ampq
}

func (g *grpcNotifyEvent) iSendRequestToGRPCAPIWithEvent() error {
	now := time.Now()
	dateStart := now.Add(time.Minute * 16).Add(time.Hour)
	dateComplete := now.Add(time.Minute * 17).Add(time.Hour)
	dateStartProto, err := ptypes.TimestampProto(dateStart)
	if err != nil {
		return err
	}
	dateCompleteProto, err := ptypes.TimestampProto(dateComplete)
	if err != nil {
		return err
	}
	_, err = GrpcClient.Add(context.Background(), &pkg.EventAddRequest{Title: "test", DateStarted: dateStartProto, DateComplete: dateCompleteProto, Notice: "test notice"})
	if err != nil {
		return err
	}
	return nil
}

func (g *grpcNotifyEvent) gRPCIGetEventsFromQueue() error {
	go func(g *grpcNotifyEvent) {
		err := g.ampq.Subscribe("notifier_test", func(delivery amqp.Delivery) {
			event := &models.Event{}
			if err := json.Unmarshal(delivery.Body, event); err != nil {
				g.err = append(g.err, err)
			}
			g.events = append(g.events, event)
		})
		if err != nil {
			g.err = append(g.err, err)
		}
	}(g)
	time.Sleep(time.Duration(4) * time.Second)

	return nil
}

func (g *grpcNotifyEvent) gRPCIHaveEventInEvents() error {
	if len(g.events) == 0 {
		return errors.New("no events")
	}
	return nil
}

func (g *grpcNotifyEvent) gRPCIHasNoErrors() error {
	if len(g.err) > 0 {
		return errors.New("unexpected errors")
	}
	return nil
}

func (g *grpcNotifyEvent) startNotifier() {
	configPath := "../config/application.yml"
	conf := config.GetConfigFromFile(configPath)
	noticeConf := config.Ampq{
		Host:     conf.Ampq.Host,
		Port:     conf.Ampq.Port,
		User:     conf.Ampq.User,
		Password: conf.Ampq.Password,
		Queue:    "notice",
	}
	noticeBus, err := amqpClient.NewAmpq(&noticeConf)
	if err != nil {
		log.Fatal(err)
	}
	g.ampq = noticeBus
}

func FeatureContextGrpcNotify(s *godog.Suite) {
	g := grpcNotifyEvent{}
	s.BeforeSuite(connectionToCalendarAPI)
	s.BeforeSuite(g.startNotifier)
	s.Step(`^I send request to GRPC API with event$`, g.iSendRequestToGRPCAPIWithEvent)
	s.Step(`^gRPC I get events from queue$`, g.gRPCIGetEventsFromQueue)
	s.Step(`^gRPC I have event in Events$`, g.gRPCIHaveEventInEvents)
	s.Step(`^gRPC I has no errors$`, g.gRPCIHasNoErrors)
	s.AfterSuite(stopGrpcTest)
}
