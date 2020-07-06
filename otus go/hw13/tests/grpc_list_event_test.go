package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/golang/protobuf/ptypes"
	pkg "github.com/lenniDespero/otus-golang/hw13/internal/pkg"
	"time"
)

type grpcListEvent struct {
	events []*pkg.Event
	err    []error
}

func (g *grpcListEvent) iSendRequestToGRPCForCycleWithEventsForDayWeekAndMonth(arg1 int) error {
	days := []int64{0, 6, 20}
	for _, day := range days {
		now := time.Now()
		dateStart := now.Add(time.Minute*16 + time.Duration(day)*time.Hour*24)
		dateComplete := now.Add(time.Minute*17 + time.Duration(day)*time.Hour*24)
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
			g.err = append(g.err, err)
		}
	}
	return nil
}

func (g *grpcListEvent) gRPCHasNoErrorsInTheseCases() error {
	if len(g.err) > 0 {
		return errors.New("expected 0 errors")
	}
	return nil
}

func (g *grpcListEvent) gRPCISendRequestWithTypeDay() error {
	answer, err := GrpcClient.GetDay(context.Background(), &pkg.EventsGetDayRequest{})
	if err != nil {
		g.err = append(g.err, err)
	}
	g.events = answer.Events
	return nil
}

func (g *grpcListEvent) gRPCIGetEventsWithEvents(arg1 int) error {
	if len(g.events) != arg1 {
		return errors.New(fmt.Sprintf("Expected %d events, got %d", arg1, len(g.events)))
	}
	return nil
}

func (g *grpcListEvent) gRPCIHaveNotErrors() error {
	if len(g.err) > 0 {
		return errors.New("expected 0 errors")
	}
	return nil
}

func (g *grpcListEvent) gRPCISendRequestWithTypeWeek() error {
	answer, err := GrpcClient.GetWeek(context.Background(), &pkg.EventsGetWeekRequest{})
	if err != nil {
		g.err = append(g.err, err)
	}
	g.events = answer.Events
	return nil
}

func (g *grpcListEvent) gRPCISendRequestWithTypeMonth() error {
	answer, err := GrpcClient.GetMonth(context.Background(), &pkg.EventsGetMonthRequest{})
	if err != nil {
		g.err = append(g.err, err)
	}
	g.events = answer.Events
	return nil
}

func FeatureContextGrpcList(s *godog.Suite) {
	g := grpcListEvent{}
	s.BeforeSuite(connectionToCalendarAPI)
	s.Step(`^I send request to GRPC for cycle with (\d+) events for day, week, and month$`, g.iSendRequestToGRPCForCycleWithEventsForDayWeekAndMonth)
	s.Step(`^gRPC has no errors in these cases$`, g.gRPCHasNoErrorsInTheseCases)
	s.Step(`^gRPC I send request with type day$`, g.gRPCISendRequestWithTypeDay)
	s.Step(`^gRPC I get Events with (\d+) events`, g.gRPCIGetEventsWithEvents)
	s.Step(`^gRPC I have not errors$`, g.gRPCIHaveNotErrors)
	s.Step(`^gRPC I send request with type week$`, g.gRPCISendRequestWithTypeWeek)
	s.Step(`^gRPC I get Events with (\d+) events$`, g.gRPCIGetEventsWithEvents)
	s.Step(`^gRPC I have not errors$`, g.gRPCIHaveNotErrors)
	s.Step(`^gRPC I send  request with type month$`, g.gRPCISendRequestWithTypeMonth)
	s.Step(`^gRPC I get Events with (\d+) events$`, g.gRPCIGetEventsWithEvents)
	s.Step(`^gRPC I have not errors$`, g.gRPCIHaveNotErrors)
	s.AfterSuite(stopGrpcTest)
}
