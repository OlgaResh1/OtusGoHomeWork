package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage"                      //nolint:depguard
	memorystorage "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	"github.com/stretchr/testify/require"                                                              //nolint:depguard
)

func TestHttpServer(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{}
	cfg.HTTP.Address = "localhost:8080"
	logger := &logger.Logger{}
	memst, err := memorystorage.New(cfg)
	require.NoError(t, err)

	app := app.New(logger, memst)
	server := NewServer(cfg, logger, app)

	serv := httptest.NewServer(server.Router())
	defer serv.Close()

	client := &http.Client{}
	var eventID1, eventID2 storage.EventID
	t.Run("test createEvent1", func(t *testing.T) {
		jevent := `{"title":"event1", "description": "desc1", "time":"2024-05-13T00:00:00Z" }`
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, serv.URL+"/5/events/create",
			bytes.NewReader([]byte(jevent)))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		result, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code for is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
		eventResult := &storage.Event{}
		err = json.Unmarshal(result, eventResult)
		require.Nil(t, err)
		eventID1 = eventResult.ID
	})
	t.Run("test createEvent2", func(t *testing.T) {
		jevent := `{"title":"event2", "description": "desc2", "time":"2024-05-23T10:00:00Z" }`
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, serv.URL+"/5/events/create",
			bytes.NewReader([]byte(jevent)))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		result, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
		eventResult := &storage.Event{}
		err = json.Unmarshal(result, eventResult)
		require.Nil(t, err)
		eventID2 = eventResult.ID
	})

	t.Run("test getEvents1", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serv.URL+"/5/eventsForDay/2024/05/13",
			bytes.NewReader([]byte{}))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		result, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code for is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
		var eventResult []*storage.Event
		err = json.Unmarshal(result, &eventResult)
		require.Nil(t, err)
		require.Len(t, eventResult, 1)
		require.Equal(t, eventResult[0].ID, eventID1)
	})
	t.Run("test getEvents2", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serv.URL+"/5/eventsForWeek/2024/05/20",
			bytes.NewReader([]byte{}))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		result, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code for is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
		var eventResult []*storage.Event
		err = json.Unmarshal(result, &eventResult)
		require.Nil(t, err)
		require.Len(t, eventResult, 1)
		require.Equal(t, eventResult[0].ID, eventID2)
	})
	t.Run("test getEvents3", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serv.URL+"/5/eventsForMonth/2024/05",
			bytes.NewReader([]byte{}))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		result, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code for is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
		var eventResult []*storage.Event
		err = json.Unmarshal(result, &eventResult)
		require.Nil(t, err)
		require.Len(t, eventResult, 2)
	})
	t.Run("test updateEvent1", func(t *testing.T) {
		jevent := `{"title":"event1-upd", "description": "desc1", "time":"2024-05-13T10:00:00Z" }`
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, serv.URL+"/5/events/"+strconv.Itoa(int(eventID1)),
			bytes.NewReader([]byte(jevent)))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code for is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("test deleteEvent1", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serv.URL+"/5/events/"+strconv.Itoa(int(eventID2)),
			bytes.NewReader([]byte{}))
		require.Nil(t, err)
		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code for is wrong. Have: %d, want: %d.", resp.StatusCode, http.StatusOK)
		}
	})
}
