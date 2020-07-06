package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/monitor"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/storage/sql"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lenniDespero/otus-golang/hw13/internal/calendar"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"

	"github.com/gorilla/mux"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/spf13/pflag"
)

type MyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	calendar calendar.CalendarInterface
	ctx      context.Context
	stats    *prometheus.SummaryVec
	codes    *prometheus.CounterVec
}

type eventRequest struct {
	ID           string
	Title        string
	Notice       string
	DateStarted  time.Time
	DateComplete time.Time
}

func (err *MyError) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Message)
}

func sendResponse(msg []byte, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", msg)
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
	cal := calendar.New(storage)
	logger.Info("Calendar was created")
	InitServer(conf.HttpListen.Ip, conf.HttpListen.Port, cal)
}

//Init http server
func InitServer(listenIp string, listenPort string, calendar calendar.CalendarInterface) {
	server := &Server{
		calendar: calendar,
		ctx:      context.Background(),
		stats:    monitor.NewSummaryVec("calendar_http", "stats", "Get events statistics"),
		codes:    monitor.NewCodesVec("calendar_http"),
	}
	router := mux.NewRouter()
	router.HandleFunc("/hello", server.hello).Methods("GET")
	router.HandleFunc("/add", server.add).Methods("POST")
	router.HandleFunc("/edit/{id}", server.edit).Methods("POST")
	router.HandleFunc("/get", server.get).Methods("GET")
	router.HandleFunc("/get/{id}", server.getById).Methods("GET")
	router.HandleFunc("/delete/{id}", server.delete).Methods("POST")
	router.HandleFunc("/events", server.events).Queries("time_before", "{time_before}", "time_length", "{time_length}").Methods("GET")
	router.HandleFunc("/events/{type}", server.eventsPlan).Methods("GET")

	srv := &http.Server{
		Addr:         listenIp + ":" + listenPort,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 2113),
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatal("Unable to start a http server", "error", err)
		}
	}()
	logger.Info("Prometheus HTTP server started")
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error while starting HTTP server", "error", err)
		}
	}()
	logger.Info("HTTP server started on host: " + listenIp + ", port: " + listenPort)

	<-done
	logger.Info("HTTP server stopped")

	ctx, cancel := context.WithTimeout(server.ctx, 5*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed", "error", err)
	}
	logger.Info("HTTP server exited properly")
}

func (s Server) hello(w http.ResponseWriter, r *http.Request) {
	logger.Info("Incoming message",
		"host", r.Host,
		"url", r.URL.Path)
	startFunc := time.Now()
	message := []byte(`{"Message":"Hello world"}`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", message)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.Path).Inc()
}

func (s Server) add(w http.ResponseWriter, r *http.Request) {
	logger.Info("Ger request add")
	startFunc := time.Now()
	data := &eventRequest{}
	var err error
	data.Title = r.FormValue("title")
	data.Notice = r.FormValue("notice")
	data.DateStarted, err = time.Parse("2006-01-02 15:04:05", r.FormValue("dateStarted"))
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
		return
	}
	data.DateComplete, err = time.Parse("2006-01-02 15:04:05", r.FormValue("dateComplete"))
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
		return
	}
	if data.DateStarted.After(data.DateComplete) {
		msg, _ := json.Marshal(MyError{http.StatusBadRequest, "Check dates date_complete before date_started"})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
		return
	}
	ip := r.RemoteAddr
	eventId, err := s.calendar.Add(data.Title, data.DateStarted, data.DateComplete, data.Notice, iptoInt(ip))
	if err != nil {
		if err == types.ErrDateBusy {
			msg, _ := json.Marshal(MyError{http.StatusBadRequest, err.Error()})
			sendResponse(msg, http.StatusBadRequest, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.URL.Path).Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
			return
		}
	}
	message := fmt.Sprintf(`{"id":"%s"}`, eventId)
	msg := []byte(message)
	sendResponse(msg, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.Path).Inc()
}

func (s Server) edit(w http.ResponseWriter, r *http.Request) {
	logger.Info("Ger request edit")
	startFunc := time.Now()
	data := &eventRequest{}
	vars := mux.Vars(r)
	var err error
	data.ID = r.FormValue("id")
	data.Title = r.FormValue("title")
	data.Notice = r.FormValue("notice")
	data.DateStarted, err = time.Parse("2006-01-02 15:04:05", r.FormValue("dateStarted"))
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/edit/id").Inc()
		return
	}
	data.DateComplete, err = time.Parse("2006-01-02 15:04:05", r.FormValue("dateComplete"))
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/edit/id").Inc()
		return
	}
	if data.DateStarted.After(data.DateComplete) {
		msg, _ := json.Marshal(MyError{http.StatusBadRequest, "Check dates date_complete before date_started"})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/edit/id").Inc()
		return
	}
	ip := r.RemoteAddr
	err = s.calendar.Edit(vars["id"], models.Event{ID: data.ID, Title: data.Title, DateStarted: data.DateStarted, DateComplete: data.DateComplete, Notice: data.Notice}, iptoInt(ip))
	if err != nil {
		if err == types.ErrDateBusy {
			msg, _ := json.Marshal(MyError{http.StatusBadRequest, err.Error()})
			sendResponse(msg, http.StatusBadRequest, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusBadRequest), "/edit/id").Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/edit/id").Inc()
			return
		}
	}
	message := fmt.Sprintf(`{"Message":"Event with id %s was changed"}`, vars["id"])
	msg := []byte(message)
	sendResponse(msg, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), "/edit/id").Inc()
}

func (s Server) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger.Info("Incoming message get",
		"host", r.Host,
		"url", r.URL.Path)
	startFunc := time.Now()
	events, err := s.calendar.GetEvents()
	if err != nil {
		if err == types.ErrNotFound {
			msg, _ := json.Marshal(MyError{http.StatusNotFound, err.Error()})
			sendResponse(msg, http.StatusNotFound, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.URL.Path).Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
			return
		}
	}
	ev, err := json.Marshal(events)
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
		return
	}
	sendResponse(ev, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.Path).Inc()
}

func (s Server) getById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	logger.Info("Incoming message get by id",
		"host", r.Host,
		"url", r.URL.Path)
	startFunc := time.Now()
	event, err := s.calendar.GetEventByID(vars["id"])
	if err != nil {
		if err == types.ErrNotFound {
			msg, _ := json.Marshal(MyError{http.StatusNotFound, err.Error()})
			sendResponse(msg, http.StatusNotFound, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusNotFound), "/get/id").Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/get/id").Inc()
			return
		}
	}
	ev, err := json.Marshal(event)
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/get/id").Inc()
		return
	}
	sendResponse(ev, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), "/get/id").Inc()
}

func (s Server) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger.Info("Incoming message delete id",
		"host", r.Host,
		"url", r.URL.Path)
	startFunc := time.Now()
	err := s.calendar.Delete(vars["id"])
	if err != nil {
		if err == types.ErrNotFound {
			msg, _ := json.Marshal(MyError{http.StatusNotFound, err.Error()})
			sendResponse(msg, http.StatusNotFound, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusNotFound), "/delete/id").Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), "/delete/id").Inc()
			return
		}
	}
	message := fmt.Sprintf(`{"Message":"event with id %s was deleted"}`, vars["id"])
	msg := []byte(message)
	sendResponse(msg, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), "/delete/id").Inc()
}

func (s Server) events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger.Info("Incoming message events",
		"host", r.Host,
		"url", r.URL.Path)
	startFunc := time.Now()
	v := r.URL.Query()
	timeBefore := v.Get("time_before")
	timeLength := v.Get("time_length")
	if timeBefore == "" || timeLength == "" {
		msg, _ := json.Marshal(MyError{http.StatusBadRequest, "Query parameters time_before and time_length required"})
		sendResponse(msg, http.StatusBadRequest, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.URL.Path).Inc()
		return
	}
	events, err := s.calendar.GetEventsByStartPeriod(timeBefore, timeLength)
	if err != nil {
		if err == types.ErrNotFound {
			msg, _ := json.Marshal(MyError{http.StatusNotFound, err.Error()})
			sendResponse(msg, http.StatusNotFound, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.URL.Path).Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
			return
		}
	}
	ev, err := json.Marshal(events)
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
		return
	}
	sendResponse(ev, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.Path).Inc()
}

func (s Server) eventsPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger.Info("Incoming message events",
		"host", r.Host,
		"url", r.URL.Path)
	startFunc := time.Now()
	timeBefore := "0"
	timeLength := ""
	vars := mux.Vars(r)
	now := time.Now()
	if vars["type"] == "day" {
		endDay := time.Until(now.Add(time.Duration(24) * time.Hour))
		timeLength = strconv.FormatInt(int64(endDay.Round(time.Minute).Minutes()), 10)
	} else if vars["type"] == "week" {
		endDay := time.Until(now.Add(time.Duration(24) * time.Hour * 7))
		timeLength = strconv.FormatInt(int64(endDay.Round(time.Minute).Minutes()), 10)
	} else if vars["type"] == "month" {
		endDay := time.Until(now.Add(time.Duration(24) * time.Hour * 30))
		timeLength = strconv.FormatInt(int64(endDay.Round(time.Minute).Minutes()), 10)
	} else {
		msg, _ := json.Marshal(MyError{http.StatusBadRequest, "Type must be day | week| month"})
		sendResponse(msg, http.StatusBadRequest, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.URL.Path).Inc()
		return
	}
	events, err := s.calendar.GetEventsByStartPeriod(timeBefore, timeLength)
	if err != nil {
		if err == types.ErrNotFound {
			msg, _ := json.Marshal(MyError{http.StatusNotFound, err.Error()})
			sendResponse(msg, http.StatusNotFound, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.URL.Path).Inc()
			return
		} else {
			msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
			sendResponse(msg, http.StatusInternalServerError, w)
			s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
			return
		}
	}
	ev, err := json.Marshal(events)
	if err != nil {
		msg, _ := json.Marshal(MyError{http.StatusInternalServerError, err.Error()})
		sendResponse(msg, http.StatusInternalServerError, w)
		s.codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.URL.Path).Inc()
		return
	}
	sendResponse(ev, http.StatusOK, w)
	dur := time.Since(startFunc)
	s.stats.WithLabelValues(r.URL.Path).Observe(dur.Seconds())
	s.codes.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.Path).Inc()
}

func iptoInt(ip string) int64 {
	bits := strings.Split(ip, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum int64
	// left shifting 24,16,8,0 and bitwise OR
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}
