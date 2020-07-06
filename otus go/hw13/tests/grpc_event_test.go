package main

import (
	"context"
	"errors"
	"github.com/cucumber/godog"
	"github.com/golang/protobuf/ptypes"
	pkg "github.com/lenniDespero/otus-golang/hw13/internal/pkg"
)

type grpcAddEvent struct {
	addResponse *pkg.EventAddResponse
	getResponse *pkg.EventGetByIdResponse
	err         error
}

func (g *grpcAddEvent) gRPCISendCreateRequestToAPI() error {
	var err error
	g.addResponse, err = GrpcClient.Add(context.Background(), &pkg.EventAddRequest{Title: "test", DateStarted: ptypes.TimestampNow(), DateComplete: ptypes.TimestampNow(), Notice: "test notice"})
	if err != nil {
		return err
	}
	return nil
}

func (g *grpcAddEvent) gRPCaddedEventWillBeReturnedWithIdOfTheEvent() error {
	id := g.addResponse.GetId()
	if id == "" {
		return errors.New("no id")
	}
	return nil
}

func (g *grpcAddEvent) gRPCISendGetRequestWithEventIdToAPI() error {
	var err error
	g.getResponse, err = GrpcClient.Get(context.Background(), &pkg.EventGetByIdRequest{Id: g.addResponse.GetId()})
	if err != nil {
		return err
	}
	return nil
}

func (g *grpcAddEvent) gRPCIGetResponseWithEvent() error {
	events := g.getResponse.Events
	if len(events) != 1 {
		return errors.New("expected 1 event")
	}
	return nil
}

func (g *grpcAddEvent) gRPCISendGetRequestWithNonExistingEventIdToAPI() error {
	g.getResponse, g.err = GrpcClient.Get(context.Background(), &pkg.EventGetByIdRequest{Id: g.addResponse.GetId() + "-wrongId"})
	if g.getResponse != nil {
		return errors.New("expected error")
	}
	return nil
}

func (g *grpcAddEvent) gRPCIGetResponseWithErrorCodeEventNotFound() error {
	if g.err.Error() != "rpc error: code = Unknown desc = event not found" {
		return errors.New("unexpected error")
	}
	return nil
}

func FeatureContextGrpcEvent(s *godog.Suite) {
	g := grpcAddEvent{}
	s.BeforeSuite(connectionToCalendarAPI)
	s.Step(`^gRPC I send Create request to API$`, g.gRPCISendCreateRequestToAPI)
	s.Step(`^gRPC added event will be returned with id of the event$`, g.gRPCaddedEventWillBeReturnedWithIdOfTheEvent)
	s.Step(`^gRPC I send Get request with event id to API$`, g.gRPCISendGetRequestWithEventIdToAPI)
	s.Step(`^gRPC I get response with event$`, g.gRPCIGetResponseWithEvent)
	s.Step(`^gRPC I send Get request with non existing event id to API$`, g.gRPCISendGetRequestWithNonExistingEventIdToAPI)
	s.Step(`^gRPC I get response with error code \'Event not found$`, g.gRPCIGetResponseWithErrorCodeEventNotFound)
	s.AfterSuite(stopGrpcTest)
}

func stopGrpcTest() {
	ClearDatabase()
	closeClient()
}
