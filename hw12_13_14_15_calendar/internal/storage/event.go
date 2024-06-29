package storage

import (
	"errors"
	"time"
)

type EventID int

const NotValidID EventID = -1

type EventOwnerID int

const NotValidOwnerID EventOwnerID = -1

type Event struct {
	ID            EventID       `json:"id,omitempty"`
	OwnerID       EventOwnerID  `json:"userid,omitempty"`
	Title         string        `json:"title,omitempty"`
	Description   string        `json:"description,omitempty"`
	StartDateTime time.Time     `json:"time,omitempty"`
	Duration      time.Duration `json:"duration,omitempty"`
	TimeToNotify  time.Duration `json:"timenotify,omitempty"`
}

var (
	ErrNotExistsEvent = errors.New("event not found")
	ErrNotValidEvent  = errors.New("event not valid")
	ErrDateBusy       = errors.New("date already busy")
	ErrUserNotValid   = errors.New("user is not valid")
)
