package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
	stor "github.com/lenniDespero/otus-golang/hw13/internal/pkg/types"
)

// Storage struct
type Storage struct {
	events map[string]models.Event
}

//New returns new storage
func New() *Storage {
	return &Storage{events: make(map[string]models.Event)}
}

// Add models to storage.
func (storage *Storage) Add(event models.Event, _ context.Context) (string, error) {
	for _, e := range storage.events {
		if inTimeSpan(e.DateStarted, e.DateComplete, event.DateStarted) ||
			inTimeSpan(e.DateStarted, e.DateComplete, event.DateComplete) ||
			inTimeSpan(event.DateStarted, event.DateComplete, e.DateStarted) ||
			inTimeSpan(event.DateStarted, event.DateComplete, e.DateComplete) {
			return "", stor.ErrDateBusy
		}
	}
	if event.ID == "" {
		id := uuid.New()
		event.ID = id.String()
	}
	_, ok := storage.events[event.ID]
	if ok {
		return "", stor.ErrEventIdExists
	}
	storage.events[event.ID] = event
	return event.ID, nil
}

// Edit models data in data storage
func (storage *Storage) Edit(id string, event models.Event, _ context.Context) error {
	e, ok := storage.events[id]
	if !ok {
		return stor.ErrNotFound
	} else if e.Deleted == true {
		return stor.ErrEventDeleted
	}
	if event.ID != id {
		delete(storage.events, id)
	}
	storage.events[event.ID] = event
	return nil
}

// GetEvents return all events
func (storage *Storage) GetEvents(_ context.Context) ([]models.Event, error) {
	if len(storage.events) > 0 {
		events := make([]models.Event, 0, len(storage.events))
		for _, e := range storage.events {
			if !e.Deleted {
				events = append(events, e)
			}
		}
		if len(events) > 0 {
			return events, nil
		}
	}
	return []models.Event{}, stor.ErrNotFound
}

//GetEventByID return models with ID
func (storage *Storage) GetEventByID(id string, _ context.Context) ([]models.Event, error) {
	e, ok := storage.events[id]
	if !ok {
		return []models.Event{}, stor.ErrNotFound
	} else if e.Deleted == true {
		return []models.Event{}, stor.ErrEventDeleted
	}
	return []models.Event{e}, nil
}

//Delete will mark models as deleted
func (storage *Storage) Delete(id string, _ context.Context) error {
	e, ok := storage.events[id]
	if !ok {
		return stor.ErrNotFound
	} else if e.Deleted == true {
		return stor.ErrEventDeleted
	}
	e.Deleted = true
	storage.events[id] = e
	return nil
}

func (storage *Storage) GetEventsByStartPeriod(timeBefore string, timeLength string, ctx context.Context) ([]models.Event, error) {
	timeBeforeInt, err := strconv.Atoi(timeBefore)
	if err != nil {
		return nil, err
	}
	timeLengthInt, err := strconv.Atoi(timeLength)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	dateBefore := now.Local().Add(time.Minute * time.Duration(timeBeforeInt))
	dateAfter := now.Local().Add(time.Minute * time.Duration(timeBeforeInt+timeLengthInt))
	events := make([]models.Event, 0, len(storage.events))
	for _, e := range storage.events {
		if inTimeSpan(dateBefore, dateAfter, e.DateStarted) {
			events = append(events, e)
		}
	}
	return events, nil
}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
