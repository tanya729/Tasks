package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/monitor"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/storage/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"net"
	"net/http"
	"os/user"
	"strconv"
	"time"

	pkg "github.com/lenniDespero/otus-golang/hw13/internal/pkg"

	"github.com/lenniDespero/otus-golang/hw13/internal/calendar"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/spf13/pflag"

	"github.com/golang/protobuf/ptypes"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type calendarpb struct {
	сalendar calendar.Calendar
	stats    *prometheus.SummaryVec
	codes    *prometheus.CounterVec
}

func (c calendarpb) Edit(ctx context.Context, e *pkg.EventEditRequest) (*pkg.EventEditResponse, error) {
	logger.Info(fmt.Sprintf("Got request Edit %v", e.String()))
	startFunc := time.Now()
	startDate, err := ptypes.Timestamp(e.Event.DateStarted)
	if err != nil {
		c.codes.WithLabelValues(codes.InvalidArgument.String(), "Edit").Inc()
		return nil, err
	}
	endDate, err := ptypes.Timestamp(e.Event.DateComplete)
	if err != nil {
		c.codes.WithLabelValues(codes.InvalidArgument.String(), "Edit").Inc()
		return nil, err
	}
	userId, err := getUserId()
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "Edit").Inc()
		return nil, err
	}
	err = c.сalendar.Edit(e.Id, models.Event{ID: e.Event.Id, Title: e.Event.Title, DateStarted: startDate, DateComplete: endDate, Notice: e.Event.Notice}, userId)
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "Edit").Inc()
		return nil, err
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("Edit").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "Edit").Inc()
	return &pkg.EventEditResponse{}, nil
}

func (c calendarpb) Delete(ctx context.Context, e *pkg.EventDeleteRequest) (*pkg.EventDeleteResponse, error) {
	logger.Info(fmt.Sprintf("Got request Delete %v", e.String()))
	startFunc := time.Now()
	err := c.сalendar.Delete(e.Id)
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "Delete").Inc()
		return nil, err
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("GetDelete").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "Delete").Inc()
	return &pkg.EventDeleteResponse{}, nil
}

func (c calendarpb) Get(ctx context.Context, e *pkg.EventGetByIdRequest) (*pkg.EventGetByIdResponse, error) {
	logger.Info(fmt.Sprintf("Got request Get %v", e.String()))
	startFunc := time.Now()
	ev, err := c.сalendar.GetEventByID(e.Id)
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "Get").Inc()
		return nil, err
	}
	respEvents := make([]*pkg.Event, 0, len(ev))
	for _, row := range ev {
		respEvents = append(respEvents, convertToProtoEvent(&row))
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("Get").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "Get").Inc()
	return &pkg.EventGetByIdResponse{Events: respEvents}, nil
}

func (c calendarpb) GetAll(ctx context.Context, e *pkg.EventGetAllRequest) (*pkg.EventGetAllResponse, error) {
	logger.Info(fmt.Sprintf("Got request GetAll %v", e.String()))
	startFunc := time.Now()
	ev, err := c.сalendar.GetEvents()
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "GetAll").Inc()
		return nil, err
	}
	respEvents := make([]*pkg.Event, 0, len(ev))
	for _, row := range ev {
		respEvents = append(respEvents, convertToProtoEvent(&row))
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("GetAll").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "GetAll").Inc()
	return &pkg.EventGetAllResponse{Events: respEvents}, nil
}

func (c calendarpb) Add(ctx context.Context, e *pkg.EventAddRequest) (*pkg.EventAddResponse, error) {
	startFunc := time.Now()
	startDate, err := ptypes.Timestamp(e.DateStarted)
	if err != nil {
		c.codes.WithLabelValues(codes.InvalidArgument.String(), "Add").Inc()
		return nil, err
	}
	endDate, err := ptypes.Timestamp(e.DateComplete)
	if err != nil {
		c.codes.WithLabelValues(codes.InvalidArgument.String(), "Add").Inc()
		return nil, err
	}
	userId, err := getUserId()
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "Add").Inc()
		return nil, err
	}
	logger.Info(fmt.Sprintf("Got request Add %v", models.Event{Title: e.Title, DateStarted: startDate, DateComplete: endDate}))
	id, err := c.сalendar.Add(e.Title, startDate.Local(), endDate.Local(), e.Notice, userId)
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("Add").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "Add").Inc()
	return &pkg.EventAddResponse{Id: id}, err
}

func (c calendarpb) GetDay(ctx context.Context, e *pkg.EventsGetDayRequest) (*pkg.EventsGetDayResponse, error) {
	logger.Info(fmt.Sprintf("Got request GetDay events %v", e.String()))
	startFunc := time.Now()
	now := time.Now()
	endDay := time.Until(now.Add(time.Duration(24) * time.Hour))
	timeLength := strconv.FormatInt(int64(endDay.Round(time.Minute).Minutes()), 10)
	ev, err := c.сalendar.GetEventsByStartPeriod("0", timeLength)
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "GetDay").Inc()
		return nil, err
	}
	respEvents := make([]*pkg.Event, 0, len(ev))
	for _, row := range ev {
		respEvents = append(respEvents, convertToProtoEvent(&row))
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("GetDay").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "GetDay").Inc()
	return &pkg.EventsGetDayResponse{Events: respEvents}, nil
}

func (c calendarpb) GetWeek(ctx context.Context, e *pkg.EventsGetWeekRequest) (*pkg.EventsGetWeekResponse, error) {
	logger.Info(fmt.Sprintf("Got request GetWeek events %v", e.String()))
	startFunc := time.Now()
	now := time.Now()
	endDay := time.Until(now.Add(time.Duration(24) * time.Hour * 7))
	timeLength := strconv.FormatInt(int64(endDay.Round(time.Minute).Minutes()), 10)
	ev, err := c.сalendar.GetEventsByStartPeriod("0", timeLength)
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "GetWeek").Inc()
		return nil, err
	}
	respEvents := make([]*pkg.Event, 0, len(ev))
	for _, row := range ev {
		respEvents = append(respEvents, convertToProtoEvent(&row))
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("GetWeek").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "GetWeek").Inc()
	return &pkg.EventsGetWeekResponse{Events: respEvents}, nil
}

func (c calendarpb) GetMonth(ctx context.Context, e *pkg.EventsGetMonthRequest) (*pkg.EventsGetMonthResponse, error) {
	logger.Info(fmt.Sprintf("Got request GetMonth events %v", e.String()))
	startFunc := time.Now()
	now := time.Now()
	endDay := time.Until(now.Add(time.Duration(24) * time.Hour * 30))
	timeLength := strconv.FormatInt(int64(endDay.Round(time.Minute).Minutes()), 10)
	ev, err := c.сalendar.GetEventsByStartPeriod("0", timeLength)
	if err != nil {
		c.codes.WithLabelValues(codes.FailedPrecondition.String(), "GetMonth").Inc()
		return nil, err
	}
	respEvents := make([]*pkg.Event, 0, len(ev))
	for _, row := range ev {
		respEvents = append(respEvents, convertToProtoEvent(&row))
	}
	dur := time.Since(startFunc)
	c.stats.WithLabelValues("GetMonth").Observe(dur.Seconds())
	c.codes.WithLabelValues(codes.OK.String(), "GetMonth").Inc()
	return &pkg.EventsGetMonthResponse{Events: respEvents}, nil
}

func convertToProtoEvent(event *models.Event) *pkg.Event {
	dateStart, err := ptypes.TimestampProto(event.DateStarted)
	if err != nil {
		logger.Fatal("Cant't convert %v to timestamp proto", event.DateStarted)
	}
	dateComplete, err := ptypes.TimestampProto(event.DateComplete)
	if err != nil {
		logger.Fatal("Cant't convert %v to timestamp proto", event.DateStarted)
	}
	return &pkg.Event{
		Id:           event.ID,
		Title:        event.Title,
		DateStarted:  dateStart,
		DateComplete: dateComplete,
		Notice:       event.Notice,
	}
}

func StartGrpcServer(calendar calendar.Calendar, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		logger.Fatal("failed to listen %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pkg.RegisterEventServiceServer(grpcServer,
		&calendarpb{сalendar: calendar,
			stats: monitor.NewSummaryVec("calendar_grpc", "stats", "Get events statistics"),
			codes: monitor.NewCodesVec("calendar_grpc"),
		})
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 2112),
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatal("Unable to start a http server", "error", err)
		}
	}()
	logger.Info("Prometheus HTTP server started")
	grpcServer.Serve(lis)
}

func getUserId() (int64, error) {
	user, err := user.Current()
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(user.Uid, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func main() {
	var configPath = flag.String("config", "../../config/application.yml", "path to configuration flag")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	conf := config.GetConfigFromFile(*configPath)
	logger.Init(conf.Log.LogLevel, conf.Log.LogFile)
	storage, err := sql.New(&conf.DBConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	calendar := calendar.New(storage)
	logger.Info("GRPC server start")
	StartGrpcServer(*calendar, conf.GrpcServer.Port)
}
