package calendar

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/storage"
	stor "github.com/lenniDespero/otus-golang/hw13/internal/pkg/types"
)

type inputData struct {
	Title        string
	DateStarted  time.Time
	DateComplete time.Time
	Notice       string
}

func prepareCalendar() *Calendar {
	calendar := New(storage.New())

	events := []inputData{
		{
			Title:        "first title",
			DateStarted:  time.Date(2020, time.January, 1, 13, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 1, 14, 20, 12, 0, time.Local),
			Notice:       "test notice"},
		{

			Title:        "first title",
			DateStarted:  time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 2, 20, 20, 12, 0, time.Local),
			Notice:       "test notice 2"},
		{

			Title:        "first title",
			DateStarted:  time.Date(2020, time.January, 3, 9, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 3, 23, 20, 12, 0, time.Local)},
	}
	for _, event := range events {
		_, err := calendar.Add(event.Title, event.DateStarted, event.DateComplete, event.Notice, 13)
		if err != nil {
			log.Fatalf("unexpected error: %s", err.Error())
		}
	}
	return calendar
}

func getLastId(calendar *Calendar) string {
	events, err := calendar.GetEvents()
	if err != nil {
		log.Fatalf("unexpected error: %s", err.Error())
	}
	return events[len(events)-1].ID
}

func TestNew(t *testing.T) {
	calendar := New(storage.New())
	calendar2 := New(storage.New())
	if !reflect.DeepEqual(calendar, calendar2) {
		t.Errorf("Not equal data in storage: %v, %v", calendar, calendar2)
	}
}

func TestCalendar_Add(t *testing.T) {
	calendar := prepareCalendar()
	_, err := calendar.Add("some event", time.Date(2021, time.January, 3, 9, 20, 12, 0, time.Local), time.Date(2021, time.January, 3, 10, 20, 12, 0, time.Local), "something", 13)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestCalendar_GetEventByID(t *testing.T) {
	calendar := prepareCalendar()
	testData := inputData{
		Title:        "some event",
		DateStarted:  time.Date(2021, time.January, 3, 9, 20, 12, 0, time.Local),
		DateComplete: time.Date(2021, time.January, 3, 10, 20, 12, 0, time.Local),
		Notice:       "something",
	}
	id, err := calendar.Add(testData.Title, testData.DateStarted, testData.DateComplete, testData.Notice, 13)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	event, err := calendar.GetEventByID(id)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if event[0].Title != testData.Title ||
		event[0].DateStarted != testData.DateStarted ||
		event[0].DateComplete != testData.DateComplete {
		t.Errorf("event has wrong data")
	}
}

func TestCalendar_GetEvents(t *testing.T) {
	calendar := prepareCalendar()
	_, err := calendar.GetEvents()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestCalendar_Delete(t *testing.T) {
	calendar := prepareCalendar()
	id := getLastId(calendar)
	err := calendar.Delete(id)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	_, err = calendar.GetEventByID(id)
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrEventDeleted)
	} else if err != stor.ErrEventDeleted {
		t.Errorf("expected error: %s, get: %s", stor.ErrEventDeleted, err.Error())
	}
}

func TestCalendar_Edit(t *testing.T) {
	calendar := prepareCalendar()
	id := getLastId(calendar)
	events, err := calendar.GetEventByID(id)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	events[0].Title = "new Title"
	err = calendar.Edit(id, events[0], 14)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
