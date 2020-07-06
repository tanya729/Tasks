package models

import (
	"time"
)

//Event structure
type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title" validate:"required"`
	DateEdited   time.Time `json:"date_edited"`
	EditorID     int64     `json:"editor_id"`
	DateCreated  time.Time `json:"date_created"`
	CreatorID    int64     `json:"creator_id"`
	DateStarted  time.Time `json:"date_started"`
	DateComplete time.Time `json:"date_comlete"`
	Notice       string    `json:"notice"`
	Deleted      bool      `json:"deleted"`
}
