package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/stretchr/testify/require"                                         //nolint:depguard
)

func testEvent(ownerID storage.EventOwnerID, eventTimeString string) storage.Event {
	eventTime, err := time.Parse("02.01.2006 15:04:05", eventTimeString)
	if err != nil {
		return storage.Event{}
	}
	return storage.Event{
		OwnerID:       ownerID,
		Title:         "test event",
		Description:   "test event description",
		StartDateTime: eventTime,
		Duration:      0,
		TimeToNotify:  0,
	}
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()

	s, err := New(cfg)
	require.NoError(t, err)

	var userID storage.EventOwnerID = 2

	eventID1, err := s.CreateEvent(ctx, testEvent(userID, "13.05.2024 12:00:00"))
	require.NoError(t, err)

	eventID2, err := s.CreateEvent(ctx, testEvent(userID, "16.05.2024 20:00:00"))
	require.NoError(t, err)

	_, err = s.CreateEvent(ctx, testEvent(userID, "23.05.2024 23:00:00"))
	require.NoError(t, err)

	time1, _ := time.Parse("02.01.2006 15:04:05", "13.05.2024 10:00:00")
	events, err := s.GetEventsForDay(ctx, userID, time1)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, eventID1, events[0].ID)

	events, err = s.GetEventsForWeek(ctx, userID, time1)
	require.NoError(t, err)
	require.Len(t, events, 2)

	events, err = s.GetEventsForMonth(ctx, userID, time1)
	require.NoError(t, err)
	require.Len(t, events, 3)

	events, err = s.GetEventsAll(ctx, userID)
	require.NoError(t, err)
	require.Len(t, events, 3)

	err = s.UpdateEvent(ctx, eventID1, testEvent(userID, "26.05.2024 12:00:00"))
	require.NoError(t, err)

	events, err = s.GetEventsForWeek(ctx, userID, time1)
	require.NoError(t, err)
	require.Len(t, events, 1)

	err = s.RemoveEvent(ctx, eventID2)
	require.NoError(t, err)

	err = s.Close(ctx)
	require.NoError(t, err)
}

func TestStorageErrors(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	s, err := New(cfg)
	require.NoError(t, err)

	var userID storage.EventOwnerID = 3

	eventID1, err := s.CreateEvent(ctx, testEvent(userID, "13.05.2024 12:00:00"))
	require.NoError(t, err)

	_, err = s.CreateEvent(ctx, testEvent(userID, "13.05.2024 12:00:00"))
	require.ErrorIs(t, err, storage.ErrDateBusy)

	_, err = s.CreateEvent(ctx, storage.Event{})
	require.ErrorIs(t, err, storage.ErrNotValidEvent)

	err = s.UpdateEvent(ctx, eventID1+1, testEvent(userID, "26.05.2024 11:00:00"))
	require.ErrorIs(t, err, storage.ErrNotExistsEvent)

	err = s.UpdateEvent(ctx, eventID1, testEvent(userID+1, "26.05.2024 11:00:00"))
	require.ErrorIs(t, err, storage.ErrUserNotValid)

	err = s.UpdateEvent(ctx, eventID1, storage.Event{})
	require.ErrorIs(t, err, storage.ErrNotValidEvent)

	err = s.Close(ctx)
	require.NoError(t, err)
}
