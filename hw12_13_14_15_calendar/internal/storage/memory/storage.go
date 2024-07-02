package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events     map[storage.EventID]storage.Event
	counterIDs storage.EventID
	mu         sync.RWMutex
}

func New(_ config.Config) (*Storage, error) {
	s := &Storage{}
	s.events = make(map[storage.EventID]storage.Event)
	return s, nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) (storage.EventID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(event.Title) == 0 || event.StartDateTime.IsZero() {
		return 0, storage.ErrNotValidEvent
	}

	for _, ev := range s.events {
		if event.OwnerID == ev.OwnerID && event.StartDateTime == ev.StartDateTime {
			return 0, storage.ErrDateBusy
		}
	}
	event.ID = s.counterIDs
	s.counterIDs++

	s.events[event.ID] = event
	return event.ID, nil
}

func (s *Storage) UpdateEvent(_ context.Context, id storage.EventID, event storage.Event) error {
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
		if ev.ID != id && event.OwnerID == ev.OwnerID &&
			event.StartDateTime == ev.StartDateTime {
			return storage.ErrDateBusy
		}
	}

	if event.OwnerID != exist.OwnerID {
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

func (s *Storage) RemoveEvent(_ context.Context, id storage.EventID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

func (s *Storage) RemoveOldEvents(_ context.Context, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range s.events {
		if event.StartDateTime.Before(date) {
			delete(s.events, event.ID)
		}
	}
	return nil
}

func (s *Storage) GetEventsAll(_ context.Context, ownerID storage.EventOwnerID) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.OwnerID == ownerID {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForDay(_ context.Context, ownerID storage.EventOwnerID, date time.Time,
) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.OwnerID == ownerID &&
			event.StartDateTime.Day() == date.Day() &&
			event.StartDateTime.Month() == date.Month() &&
			event.StartDateTime.Year() == date.Year() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForWeek(_ context.Context, ownerID storage.EventOwnerID, date time.Time,
) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)

	dateYear, dateWeek := date.ISOWeek()

	for _, event := range s.events {
		if event.OwnerID == ownerID {
			year, week := event.StartDateTime.ISOWeek()
			if year == dateYear && week == dateWeek {
				result = append(result, event)
			}
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForMonth(_ context.Context, ownerID storage.EventOwnerID, date time.Time,
) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.OwnerID == ownerID && event.StartDateTime.Year() == date.Year() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForNotification(_ context.Context, startDate time.Time, endDate time.Time,
) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.TimeToNotify > 0 {
			notify := event.StartDateTime.Add(-event.TimeToNotify)
			if notify.After(startDate) && notify.Before(endDate) {
				result = append(result, event)
			}
		}
	}
	return result, nil
}
