package storage

import (
	"errors"
	"time"
)

type EventId int

const NotValidId EventId = -1

type EventOwnerId int

type Event struct {
	Id            EventId
	OwnerId       EventOwnerId
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
