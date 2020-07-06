package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	"net/http"
	"net/url"
	"os"
	"time"
)

type listsEventsTest struct {
	Resp   []*http.Response
	Errs   []error
	Events []models.Event
}

func (t *listsEventsTest) iSendRequestToAPIForCycleWithEventsForDayWeekAndMonth() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/add", apiName, apiPort)

	days := []int64{0, 6, 20}
	for _, day := range days {
		form := url.Values{}
		form.Add("title", "test title")
		form.Add("notice", "test title notice")
		now := time.Now()
		dateStart := now.Add(time.Minute*16 + time.Duration(day)*time.Hour*24)
		dateComplete := now.Add(time.Minute*17 + time.Duration(day)*time.Hour*24)
		form.Add("dateStarted", dateStart.Format("2006-01-02 15:04:05"))
		form.Add("dateComplete", dateComplete.Format("2006-01-02 15:04:05"))
		resp, err := http.PostForm(apiUrl, form)
		if err != nil {
			t.Errs = append(t.Errs, err)
			return err
		}
		t.Resp = append(t.Resp, resp)
	}
	return nil
}

func (t *listsEventsTest) hasNoErrorsInTheseCases() error {
	if len(t.Errs) > 0 {
		return t.Errs[0]
	}
	return nil
}

func (t *listsEventsTest) iSendRequestWithTypeDay() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/events/day", apiName, apiPort)

	resp, err := http.Get(apiUrl)
	t.Errs = []error{}
	if err != nil {
		t.Errs = append(t.Errs, err)
		return err
	}
	t.Resp = append([]*http.Response{}, resp)

	return nil
}

func (t *listsEventsTest) iGetEventsWithEvents(arg1 int) error {
	var events []models.Event
	err := json.NewDecoder(t.Resp[0].Body).Decode(&events)
	if err != nil {
		return err
	}
	if len(events) != arg1 {
		return errors.New(fmt.Sprintf("incorrect events count, expected %d, got %d", arg1, len(events)))
	}
	return nil
}

func (t *listsEventsTest) iHaveNotErrors() error {
	if len(t.Errs) > 0 {
		return t.Errs[0]
	}
	return nil
}

func (t *listsEventsTest) iSendRequestWithTypeWeek() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/events/week", apiName, apiPort)

	resp, err := http.Get(apiUrl)
	t.Errs = []error{}
	if err != nil {
		t.Errs = append(t.Errs, err)
		return err
	}
	t.Resp = append([]*http.Response{}, resp)
	return nil
}

func (t *listsEventsTest) iSendRequestWithTypeMonth() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/events/month", apiName, apiPort)

	resp, err := http.Get(apiUrl)
	t.Errs = []error{}
	if err != nil {
		t.Errs = append(t.Errs, err)
		return err
	}
	t.Resp = append([]*http.Response{}, resp)

	return nil
}

func (t *listsEventsTest) clear() {
	ClearDatabase()
}

func FeatureContextListEvents(s *godog.Suite) {
	list := listsEventsTest{}
	s.Step(`^I send request to API for cycle with (\d+) events for day, week, and month$`, list.iSendRequestToAPIForCycleWithEventsForDayWeekAndMonth)
	s.Step(`^has no errors in these cases$`, list.hasNoErrorsInTheseCases)
	s.Step(`^I send request with type day$`, list.iSendRequestWithTypeDay)
	s.Step(`^I get Events with (\d+) events$`, list.iGetEventsWithEvents)
	s.Step(`^I have not errors$`, list.iHaveNotErrors)
	s.Step(`^I send request with type week$`, list.iSendRequestWithTypeWeek)
	s.Step(`^I have not errors$`, list.iHaveNotErrors)
	s.Step(`^I get Events with (\d+) events$`, list.iGetEventsWithEvents)
	s.Step(`^I send  request with type month$`, list.iSendRequestWithTypeMonth)
	s.Step(`^I have not errors$`, list.iHaveNotErrors)
	s.Step(`^I get Events with (\d+) events$`, list.iGetEventsWithEvents)
	s.AfterFeature(func(*messages.GherkinDocument) {
		list.clear()
	})
}
