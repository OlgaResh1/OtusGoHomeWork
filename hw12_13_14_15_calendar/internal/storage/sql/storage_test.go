package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/require"
)

func testEvent(ownerID storage.EventOwnerID, eventTimeString string) storage.Event {
	eventTime, err := time.Parse("02.01.2006 15:04:05", eventTimeString)
	if err != nil {
		return storage.Event{}
	}
	return storage.Event{
		OwnerID:       ownerID,
		Title:         "test sql event",
		Description:   "test sql event description",
		StartDateTime: eventTime,
		Duration:      0,
		TimeToNotify:  0,
	}
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	cfg.Storage.Type = "sql"
	cfg.SQL.Dsn = "host=localhost port=5432 user=user1 password=pass1 dbname=calendardb sslmode=disable"

	s, err := New(ctx, cfg)
	require.NoError(t, err)
	defer s.Close(ctx)

	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	err = goose.Up(s.db, "../../../migrations")
	require.NoError(t, err)

	defer goose.Down(s.db, "../../../migrations")

	userID := storage.EventOwnerID(time.Now().Second())

	eventID1, err := s.CreateEvent(ctx, testEvent(userID, "13.05.2024 12:00:00"))
	require.NoError(t, err)
	require.NotNil(t, eventID1)

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
