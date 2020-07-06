package main

import (
	"context"
	"flag"
	"time"

	pkg "github.com/lenniDespero/otus-golang/hw13/internal/pkg"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/spf13/pflag"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

func main() {
	var configPath = flag.String("config", "../../config/application.yml", "path to configuration flag")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	conf := config.GetConfigFromFile(*configPath)
	logger.Init(conf.Log.LogLevel, conf.Log.LogFile)

	cc, err := grpc.Dial(conf.GrpcServer.Host+":"+conf.GrpcServer.Port, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("could not connect: %v", err)
	}
	defer cc.Close()

	c := pkg.NewEventServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	logger.Info("Send request Add")
	answer1, err := c.Add(ctx, &pkg.EventAddRequest{Title: "test", DateStarted: ptypes.TimestampNow(), DateComplete: ptypes.TimestampNow(), Notice: "test notice"})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("ok")

	logger.Info("Send request Get")
	answer2, err := c.Get(ctx, &pkg.EventGetByIdRequest{Id: answer1.Id})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug(answer2.String())

	event := answer2.Events[0]
	event.Notice = "new notice"
	logger.Info("Send request Edit")
	_, err = c.Edit(ctx, &pkg.EventEditRequest{Id: event.Id, Event: event})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("ok")

	logger.Info("Send request GetAll")
	answer4, err := c.GetAll(ctx, &pkg.EventGetAllRequest{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(answer4.String())

	logger.Info("Send request Delete")
	_, err = c.Delete(ctx, &pkg.EventDeleteRequest{Id: answer1.Id})
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("ok")
}
