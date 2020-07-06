package storage

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	stor "github.com/lenniDespero/otus-golang/hw13/internal/pkg/types"
)

func prepareStorage() *Storage {
	storage := New()

	storage.events = map[string]models.Event{
		"1": {
			ID:           "1",
			Title:        "first title",
			DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
			DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
			EditorID:     12,
			CreatorID:    13,
			Deleted:      false,
			DateStarted:  time.Date(2020, time.January, 1, 13, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 1, 14, 20, 12, 0, time.Local)},
		"2": {
			ID:           "2",
			Title:        "first title",
			DateCreated:  time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			DateEdited:   time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			EditorID:     12,
			CreatorID:    13,
			Deleted:      false,
			DateStarted:  time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 2, 20, 20, 12, 0, time.Local)},
		"3": {
			ID:           "3",
			Title:        "first title",
			DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
			DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
			EditorID:     12,
			CreatorID:    13,
			Deleted:      true,
			DateStarted:  time.Date(2020, time.January, 3, 9, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 3, 23, 20, 12, 0, time.Local)},
	}
	return storage
}

func TestNew(t *testing.T) {
	storage1 := New()
	storage2 := New()
	if !reflect.DeepEqual(storage1, storage2) {
		t.Errorf("Not equal data in storage: %v, %v", storage1, storage2)
	}
}

func TestStorage_Add(t *testing.T) {
	storage := New()
	newEvent := models.Event{
		ID:           "",
		Title:        "first title",
		DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		EditorID:     12,
		CreatorID:    13,
		Deleted:      true,
		DateStarted:  time.Date(2020, time.January, 2, 21, 20, 12, 0, time.Local),
		DateComplete: time.Date(2020, time.January, 2, 22, 20, 12, 0, time.Local)}
	_, err := storage.Add(newEvent, context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestStorage_Add_With_Prepared(t *testing.T) {
	storage := prepareStorage()
	newEvent := models.Event{
		ID:           "1",
		Title:        "first title",
		DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		EditorID:     12,
		CreatorID:    13,
		Deleted:      true,
		DateStarted:  time.Date(2020, time.January, 2, 21, 20, 12, 0, time.Local),
		DateComplete: time.Date(2020, time.January, 2, 22, 20, 12, 0, time.Local)}
	_, err := storage.Add(newEvent, context.Background())
	if err == nil {
		t.Errorf("expected error: %s, but get nil", stor.ErrEventIdExists)
	} else if err != stor.ErrEventIdExists {
		t.Errorf("expected error: %s, get %s", stor.ErrEventIdExists, err.Error())
	}
}

func TestStorage_Add_Bad_Event(t *testing.T) {
	storage := prepareStorage()
	badEvent := models.Event{
		ID:           "",
		Title:        "first title",
		DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		EditorID:     12,
		CreatorID:    13,
		Deleted:      true,
		DateStarted:  time.Date(2020, time.January, 2, 21, 20, 12, 0, time.Local),
		DateComplete: time.Date(2020, time.January, 3, 10, 20, 12, 0, time.Local)}
	_, err := storage.Add(badEvent, context.Background())
	if err == nil {
		t.Errorf("expected error: %s, but get nil", stor.ErrDateBusy)
	} else if err != stor.ErrDateBusy {
		t.Errorf("expected error: %s, get %s", stor.ErrDateBusy, err.Error())
	}
}

func TestStorage_Delete(t *testing.T) {
	storage := prepareStorage()
	err := storage.Delete("1", context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	err = storage.Delete("1", context.Background())
	if err == nil {
		t.Errorf("expected error: %s", stor.ErrEventDeleted)
	} else if err != stor.ErrEventDeleted {
		t.Errorf("expected error: %s, get %s", stor.ErrEventDeleted, err.Error())
	}
	err = storage.Delete("4", context.Background())
	if err == nil {
		t.Errorf("expected error: %s", stor.ErrNotFound)
	} else if err != stor.ErrNotFound {
		t.Errorf("expected error: %s, get : %s", stor.ErrNotFound, err.Error())
	}
}

func TestStorage_Delete_Deleted(t *testing.T) {
	storage := prepareStorage()

	err := storage.Delete("3", context.Background())
	if err == nil {
		t.Errorf("expected error: %s", stor.ErrEventDeleted)
	} else if err != stor.ErrEventDeleted {
		t.Errorf("expected error: %s, get %s", stor.ErrEventDeleted, err.Error())
	}
}

func TestStorage_Delete_Wrong(t *testing.T) {
	storage := prepareStorage()
	err := storage.Delete("4", context.Background())
	if err == nil {
		t.Errorf("expected error: %s", stor.ErrNotFound)
	} else if err != stor.ErrNotFound {
		t.Errorf("expected error: %s, get : %s", stor.ErrNotFound, err.Error())
	}
}

func TestStorage_Edit(t *testing.T) {
	storage := prepareStorage()
	event, err := storage.GetEventByID("2", context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	event[0].Title = "new Title"
	err = storage.Edit(event[0].ID, event[0], context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	newEvent, err := storage.GetEventByID("2", context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !reflect.DeepEqual(newEvent[0], event[0]) {
		t.Errorf("not equal events: %v, %v", event[0], newEvent[0])
	}
}

func TestStorage_Edit_Change_ID(t *testing.T) {
	storage := prepareStorage()
	id := "2"
	newId := "13"
	event, err := storage.GetEventByID(id, context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	event[0].ID = newId
	err = storage.Edit(id, event[0], context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	event, err = storage.GetEventByID(id, context.Background())
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrNotFound)
	} else if err != stor.ErrNotFound {
		t.Errorf("expected error: %s, get : %s", stor.ErrNotFound, err.Error())
	}
	event, err = storage.GetEventByID(newId, context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestStorage_Edit_Not_Found(t *testing.T) {
	event := models.Event{
		ID:           "4",
		Title:        "first title",
		DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		EditorID:     12,
		CreatorID:    13,
		Deleted:      true,
		DateStarted:  time.Date(2020, time.January, 2, 21, 20, 12, 0, time.Local),
		DateComplete: time.Date(2020, time.January, 2, 22, 20, 12, 0, time.Local)}
	storage := prepareStorage()
	id := "12"
	err := storage.Edit(id, event, context.Background())
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrNotFound)
	}
	if err != stor.ErrNotFound {
		t.Errorf("expected error: %s, get : %s", stor.ErrNotFound, err.Error())
	}
}

func TestStorage_Edit_Deleted(t *testing.T) {
	event := models.Event{
		ID:           "4",
		Title:        "first title",
		DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		EditorID:     12,
		CreatorID:    13,
		Deleted:      true,
		DateStarted:  time.Date(2020, time.January, 2, 21, 20, 12, 0, time.Local),
		DateComplete: time.Date(2020, time.January, 2, 22, 20, 12, 0, time.Local)}
	storage := prepareStorage()
	id := "3"
	err := storage.Edit(id, event, context.Background())
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrEventDeleted)
	}
	if err != stor.ErrEventDeleted {
		t.Errorf("expected error: %s, get : %s", stor.ErrEventDeleted, err.Error())
	}
}

func TestStorage_GetEventByID(t *testing.T) {
	testEvent := models.Event{
		ID:           "1",
		Title:        "first title",
		DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
		EditorID:     12,
		CreatorID:    13,
		Deleted:      false,
		DateStarted:  time.Date(2020, time.January, 1, 13, 20, 12, 0, time.Local),
		DateComplete: time.Date(2020, time.January, 1, 14, 20, 12, 0, time.Local)}
	storage := prepareStorage()
	e, err := storage.GetEventByID("1", context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !reflect.DeepEqual(e[0], testEvent) {
		t.Errorf("not equal events: %v, %v", e[0], testEvent)
	}
}

func TestStorage_GetEventByID_Deleted(t *testing.T) {
	storage := prepareStorage()
	_, err := storage.GetEventByID("3", context.Background())
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrEventDeleted)
	} else if err != stor.ErrEventDeleted {
		t.Errorf("expected error: %s, get: %s", stor.ErrEventDeleted, err.Error())
	}
}

func TestStorage_GetEventByID_NotFound(t *testing.T) {
	storage := prepareStorage()
	_, err := storage.GetEventByID("13", context.Background())
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrNotFound)
	} else if err != stor.ErrNotFound {
		t.Errorf("expected error: %s, get: %s", stor.ErrNotFound, err.Error())
	}
}

func TestStorage_GetEvents(t *testing.T) {
	storage := prepareStorage()
	testEvents := []models.Event{
		{
			ID:           "1",
			Title:        "first title",
			DateCreated:  time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
			DateEdited:   time.Date(2020, time.January, 1, 12, 20, 12, 0, time.Local),
			EditorID:     12,
			CreatorID:    13,
			Deleted:      false,
			DateStarted:  time.Date(2020, time.January, 1, 13, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 1, 14, 20, 12, 0, time.Local)},
		{
			ID:           "2",
			Title:        "first title",
			DateCreated:  time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			DateEdited:   time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			EditorID:     12,
			CreatorID:    13,
			Deleted:      false,
			DateStarted:  time.Date(2020, time.January, 2, 10, 20, 12, 0, time.Local),
			DateComplete: time.Date(2020, time.January, 2, 20, 20, 12, 0, time.Local)}}
	events, err := storage.GetEvents(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(events) != 2 {
		t.Errorf("expected %d events get %d", 2, len(events))
	}
	if !reflect.DeepEqual(events, testEvents) {
		t.Errorf("events %v not equal %v", events, testEvents)
	}
}

func TestStorage_GetEvents_Empty(t *testing.T) {
	storage := New()

	_, err := storage.GetEvents(context.Background())
	if err == nil {
		t.Errorf("expected error: %s, get nil", stor.ErrNotFound)
	}
	if err != stor.ErrNotFound {
		t.Errorf("expected error: %s, get: %s", stor.ErrNotFound, err.Error())
	}
}

func Test_inTimeSpan(t *testing.T) {
	test := []struct {
		start  string
		end    string
		check  string
		isTrue bool
	}{
		{"23:00", "05:00", "04:00", true},
		{"23:00", "05:00", "23:30", true},
		{"23:00", "05:00", "20:00", false},
		{"10:00", "21:00", "11:00", true},
		{"10:00", "21:00", "22:00", false},
		{"10:00", "21:00", "03:00", false},
		{"22:00", "02:00", "00:00", true},
		{"10:00", "21:00", "10:00", true},
		{"10:00", "21:00", "21:00", true},
		{"23:00", "05:00", "06:00", false},
		{"23:00", "05:00", "23:00", true},
		{"23:00", "05:00", "05:00", true},
		{"10:00", "21:00", "10:00", true},
		{"10:00", "21:00", "21:00", true},
		{"10:00", "10:00", "09:00", false},
		{"10:00", "10:00", "11:00", false},
		{"10:00", "10:00", "10:00", true},
	}
	newLayout := "15:04"
	for _, row := range test {
		check, _ := time.Parse(newLayout, row.check)
		start, _ := time.Parse(newLayout, row.start)
		end, _ := time.Parse(newLayout, row.end)
		result := inTimeSpan(start, end, check)
		if result != row.isTrue {
			t.Errorf("get %t, expected %t on row: {%s, %s, %s, %t}", result, row.isTrue, row.start, row.end, row.check, row.isTrue)
		}
	}
}
