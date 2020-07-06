package calendar

import (
	"time"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
)

// CalendarInterface interface
type CalendarInterface interface {

	// Add event to calendar.
	Add(title string, dateStarted time.Time, dateComplete time.Time, notice string, userId int64) (string, error)

	//Edit event data in calendar
	Edit(id string, event models.Event, userId int64) error

	//GetEvents return all events
	GetEvents() ([]models.Event, error)

	//GetEventByID return event with ID
	GetEventByID(id string) ([]models.Event, error)

	//Delete will mark event as deleted
	Delete(id string) error

	//GetEventsByStartPeriod return events where date start between NOW+timeBefore and NOW+timeBefore+timeLength
	GetEventsByStartPeriod(timeBefore string, timeLength string) ([]models.Event, error)
}
