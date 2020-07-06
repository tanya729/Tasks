package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	amqpClient "github.com/lenniDespero/otus-golang/hw13/internal/pkg/ampq"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	"github.com/streadway/amqp"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TestNotify struct {
	err      error
	httpResp *http.Response
	ampq     *amqpClient.Ampq
	events   []*models.Event
}

func (t *TestNotify) iSendRequestToAPIWithEvent() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/add", apiName, apiPort)
	form := url.Values{}
	form.Add("title", "test title")
	form.Add("notice", "test title notice")
	now := time.Now()
	dateStart := now.Add(time.Minute * 16).Add(time.Hour)
	dateComplete := now.Add(time.Minute * 17).Add(time.Hour)
	form.Add("dateStarted", dateStart.Format("2006-01-02 15:04:05"))
	form.Add("dateComplete", dateComplete.Format("2006-01-02 15:04:05"))
	resp, err := http.PostForm(apiUrl, form)
	if err != nil {
		t.err = err
		return err
	}
	t.httpResp = resp
	return nil
}

func (t *TestNotify) iGetEventsFromQueue() error {
	go func(t *TestNotify) {
		err := t.ampq.Subscribe("notifier_test", func(delivery amqp.Delivery) {
			event := &models.Event{}
			if err := json.Unmarshal(delivery.Body, event); err != nil {
				t.err = err
			}
			t.events = append(t.events, event)
		})
		if err != nil {
			t.err = err
		}
	}(t)
	time.Sleep(time.Duration(4) * time.Second)

	return nil
}

func (t *TestNotify) iHaveEventInEvents() error {
	if len(t.events) == 0 {
		return errors.New("no events")
	}
	return nil
}

func (t *TestNotify) iHasNoErrors() error {
	if t.err != nil {
		return t.err
	}
	return nil
}

func (t *TestNotify) start() {
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
		t.err = err
	}
	t.ampq = noticeBus
}

func (t TestNotify) stop() {
	ClearDatabase()
}

func FeatureContextNotice(s *godog.Suite) {
	t := &TestNotify{}
	s.BeforeSuite(t.start)
	s.Step(`^I send request to API with event$`, t.iSendRequestToAPIWithEvent)
	s.Step(`^I get events from queue$`, t.iGetEventsFromQueue)
	s.Step(`^I have event in Events$`, t.iHaveEventInEvents)
	s.Step(`^I has no errors$`, t.iHasNoErrors)
	s.AfterFeature(func(*messages.GherkinDocument) {
		t.stop()
	})
}
