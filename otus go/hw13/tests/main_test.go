package main

import (
	"github.com/cucumber/godog"
	pkg "github.com/lenniDespero/otus-golang/hw13/internal/pkg"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/storage/sql"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		godog.SuiteContext(s)
		FeatureContextHttpEvent(s)
		FeatureContextListEvents(s)
		FeatureContextNotice(s)
		FeatureContextGrpcEvent(s)
		FeatureContextGrpcList(s)
		FeatureContextGrpcNotify(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"feature"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func ClearDatabase() {
	configPath := "../config/application.yml"
	conf := config.GetConfigFromFile(configPath)
	storage, err := sql.New(&conf.DBConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = storage.ClearDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	storage.ConnPool.Close()
}

// Global variables for multiple tests execution
var (
	ClientConn *grpc.ClientConn
	GrpcClient pkg.EventServiceClient
)

func connectionToCalendarAPI() {
	configPath := "../config/application.yml"
	conf := config.GetConfigFromFile(configPath)
	var err error
	ClientConn, err = grpc.Dial(conf.GrpcServer.Host+":"+conf.GrpcServer.Port, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(errors.Wrap(err, "could not connect gRPC server").Error())
	}
	GrpcClient = pkg.NewEventServiceClient(ClientConn)
	if GrpcClient == nil {
		logger.Fatal(errors.New("failed creating calendar API client").Error())
	}
}

func closeClient() {
	if ClientConn != nil {
		_ = ClientConn.Close()
	}
}
