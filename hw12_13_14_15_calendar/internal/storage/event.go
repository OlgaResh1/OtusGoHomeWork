package storage

import (
	"errors"
	"time"
)

type EventID int

const NotValidID EventID = -1

type EventOwnerID int

type Event struct {
	ID            EventID
	OwnerID       EventOwnerID
	Title         string
	Description   string
	StartDateTime time.Time
	Duration      time.Duration
	TimeToNotify  time.Duration
}

var (
	ErrNotExistsEvent = errors.New("event not found")
	ErrNotValidEvent  = errors.New("event not valid")
	ErrDateBusy       = errors.New("date already busy")
	ErrUserNotValid   = errors.New("user is not valid")
)
