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
	"reflect"
	"strconv"
	"time"
)

type testStruct struct {
	resp     Response
	err      error
	httpResp *http.Response
}

type Response struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Id      string         `json:"id"`
	Events  []models.Event `json:"events"`
}

func (t *testStruct) iSendCreateRequestToAPI() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/add", apiName, apiPort)
	form := url.Values{}
	form.Add("title", "test title")
	form.Add("notice", "test title notice")
	now := time.Now()
	dateStart := now.Add(time.Minute * 16)
	dateComplete := now.Add(time.Minute * 17)
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

func (t *testStruct) addedEventWillBeReturnedWithIdOfTheEvent() error {
	var p Response
	err := json.NewDecoder(t.httpResp.Body).Decode(&p)
	if err != nil {
		return err
	}
	t.resp = p
	if t.resp.Id == "" {
		err = errors.New("empty Id error")
		t.err = err
		return err
	}
	return nil
}

func (t *testStruct) getErrorHasNoErrorsInBothCases() error {
	var err error
	if t.httpResp.StatusCode != 200 {
		err = errors.New(t.httpResp.Status)
		return err
	}
	if t.resp.Code != 0 {
		err = errors.New(strconv.Itoa(t.resp.Code) + ": " + t.resp.Message)
		return err
	}
	if t.resp.Message != "" {
		err = errors.New(strconv.Itoa(t.resp.Code) + ": " + t.resp.Message)
		return err
	}
	return nil
}

func (t *testStruct) iSendGetRequestWithEventIdToAPI() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/get/", apiName, apiPort)
	apiUrl = apiUrl + t.resp.Id

	resp, err := http.Get(apiUrl)
	if err != nil {
		t.err = err
		return err
	}
	t.httpResp = resp

	return nil
}

func (t *testStruct) iGetResponseWithEvent() error {
	var events []models.Event
	err := json.NewDecoder(t.httpResp.Body).Decode(&events)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(events, []models.Event{}) {
		return errors.New("empty events")
	}
	return nil
}

func (t *testStruct) getErrorHasNoErrors() error {
	if t.httpResp.StatusCode != 200 {
		return errors.New(t.httpResp.Status)
	}
	return nil
}

func (t *testStruct) iSendGetRequestWithNonExistingEventIdToAPI() error {
	apiName := os.Getenv("API_NAME")
	apiPort := os.Getenv("HTTP_PORT")
	apiUrl := fmt.Sprintf("http://%s:%s/get/", apiName, apiPort)
	apiUrl = apiUrl + "strange_id"

	resp, err := http.Get(apiUrl)
	if err != nil {
		t.err = err
		return err
	}
	t.httpResp = resp

	return nil
}

func (t *testStruct) iGetResponseWithErrorCodeEventNotFound() error {
	if t.httpResp.StatusCode != 404 {
		return errors.New(fmt.Sprintf("expected code 404 get %d", t.httpResp.StatusCode))
	}
	return nil
}

func FeatureContextHttpEvent(s *godog.Suite) {
	tStruct := testStruct{}
	s.Step(`^I send Create request to API$`, tStruct.iSendCreateRequestToAPI)
	s.Step(`^added event will be returned with id of the event$`, tStruct.addedEventWillBeReturnedWithIdOfTheEvent)
	s.Step(`^GetError has no errors in both cases$`, tStruct.getErrorHasNoErrorsInBothCases)
	s.Step(`^I send Get request with event id to API$`, tStruct.iSendGetRequestWithEventIdToAPI)
	s.Step(`^I get response with event$`, tStruct.iGetResponseWithEvent)
	s.Step(`^GetError has no errors$`, tStruct.getErrorHasNoErrors)
	s.Step(`^I send Get request with non existing event id to API$`, tStruct.iSendGetRequestWithNonExistingEventIdToAPI)
	s.Step(`^I get response with error code \'Event not found\'$`, tStruct.iGetResponseWithErrorCodeEventNotFound)

	s.AfterFeature(func(*messages.GherkinDocument) {
		tStruct.stop()
	})
}

func (t testStruct) stop() {
	ClearDatabase()
}
