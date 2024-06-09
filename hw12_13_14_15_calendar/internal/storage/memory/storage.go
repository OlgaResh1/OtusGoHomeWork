package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events     map[storage.EventId]storage.Event
	counterIds storage.EventId
	mu         sync.RWMutex
}

func New(cfg config.Config) (*Storage, error) {
	s := &Storage{}
	s.events = make(map[storage.EventId]storage.Event)
	return s, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.EventId, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(event.Title) == 0 || event.StartDateTime.IsZero() {
		return 0, storage.ErrNotValidEvent
	}

	for _, ev := range s.events {
		if event.OwnerId == ev.OwnerId && event.StartDateTime == ev.StartDateTime {
			return 0, storage.ErrDateBusy
		}
	}
	event.Id = s.counterIds
	s.counterIds++

	s.events[event.Id] = event
	return event.Id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id storage.EventId, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	exist, ok := s.events[id]
	if !ok {
		return storage.ErrNotExistsEvent
	}
	if len(event.Title) == 0 || event.StartDateTime.IsZero() {
		return storage.ErrNotValidEvent
	}
	for _, ev := range s.events {
		if ev.Id != id && event.OwnerId == ev.OwnerId &&
			event.StartDateTime == ev.StartDateTime {
			return storage.ErrDateBusy
		}
	}

	if event.OwnerId != exist.OwnerId {
		return storage.ErrUserNotValid
	}

	exist.Title = event.Title
	exist.Description = event.Description
	exist.StartDateTime = event.StartDateTime
	exist.Duration = event.Duration
	exist.TimeToNotify = event.TimeToNotify

	s.events[id] = event
	return nil
}

func (s *Storage) RemoveEvent(ctx context.Context, id storage.EventId) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEventsAll(ctx context.Context, ownerId storage.EventOwnerId) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.OwnerId == ownerId {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForDay(ctx context.Context, ownerId storage.EventOwnerId, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.OwnerId == ownerId &&
			event.StartDateTime.Day() == date.Day() &&
			event.StartDateTime.Month() == date.Month() &&
			event.StartDateTime.Year() == date.Year() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForWeek(ctx context.Context, ownerId storage.EventOwnerId, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)

	dateYear, dateWeek := date.ISOWeek()

	for _, event := range s.events {
		if event.OwnerId == ownerId {
			year, week := event.StartDateTime.ISOWeek()
			if year == dateYear && week == dateWeek {
				result = append(result, event)
			}
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForMonth(ctx context.Context, ownerId storage.EventOwnerId, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.OwnerId == ownerId && event.StartDateTime.Year() == date.Year() {
			result = append(result, event)
		}
	}
	return result, nil
}
