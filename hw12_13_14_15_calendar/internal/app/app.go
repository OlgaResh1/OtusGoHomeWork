package app

import (
	"context"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct { // TODO
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Storage interface {
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, event storage.Event) (storage.EventID, error)
	UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error
	RemoveEvent(ctx context.Context, id storage.EventID) error
	GetEventsAll(ctx context.Context, ownerID storage.EventOwnerID) ([]storage.Event, error)
	GetEventsForDay(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForWeek(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForMonth(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, title string) (storage.EventID, error) {
	return a.storage.CreateEvent(ctx, storage.Event{Title: title})
}

func (a *App) UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}
