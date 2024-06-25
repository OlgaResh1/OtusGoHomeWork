package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/gorilla/mux"                                                      //nolint:depguard
)

type Server struct { // TODO
	server *http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) (storage.EventID, error)
	UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error
	RemoveEvent(ctx context.Context, id storage.EventID) error
	GetEventsAll(ctx context.Context, ownerID storage.EventOwnerID) ([]storage.Event, error)
	GetEventsForDay(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForWeek(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForMonth(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
}

const ErrorParceDate = "error parse date "

func NewServer(cfg config.Config, logger Logger, app Application) *Server {
	srv := &http.Server{
		Addr:              cfg.HTTP.Address,
		ReadHeaderTimeout: cfg.HTTP.RequestTimeout,
	}
	return &Server{logger: logger, app: app, server: srv}
}

func (s *Server) Router() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/{user_id:[0-9]+}/events/create", s.createEventHandler).Methods("POST")
	mux.HandleFunc("/{user_id:[0-9]+}/events/{id:[0-9]+}", s.updateEventHandler).Methods("PUT")
	mux.HandleFunc("/{user_id:[0-9]+}/events/{id:[0-9]+}", s.removeEventHandler).Methods("DELETE")
	mux.HandleFunc("/{user_id:[0-9]+}/eventsForDay/{year}/{month}/{day}", s.getEventForDayHandler).Methods("GET")
	mux.HandleFunc("/{user_id:[0-9]+}/eventsForWeek/{year}/{month}/{day}", s.getEventForWeekHandler).Methods("GET")
	mux.HandleFunc("/{user_id:[0-9]+}/eventsForMonth/{year}/{month}", s.getEventForMonthHandler).Methods("GET")
	mux.HandleFunc("/{user_id:[0-9]+}/events", s.getEventHandler).Methods("GET")
	return mux
}

func (s *Server) Start(_ context.Context) error {
	s.server.Handler = s.loggingMiddleware(s.Router())
	s.logger.Info(fmt.Sprintf("HTTP Server started %s", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) createEventHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	userID, err := strconv.Atoi(v["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse userid " + v["user_id"])
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error read request", "error", err)
		return
	}
	event := storage.Event{}

	err = json.Unmarshal(body, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse json", "error", err)
		return
	}
	event.OwnerID = storage.EventOwnerID(userID)
	eventID, err := s.app.CreateEvent(req.Context(), event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("error create event", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"id": %d }`, int(eventID))))
}

func (s *Server) updateEventHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	eventID, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse id " + v["id"])
		return
	}
	userID, err := strconv.Atoi(v["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse userid " + v["user_id"])
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error read request", "error", err)
		return
	}
	event := storage.Event{}

	err = json.Unmarshal(body, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse json", "error", err)
		return
	}
	event.OwnerID = storage.EventOwnerID(userID)
	err = s.app.UpdateEvent(req.Context(), storage.EventID(eventID), event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("error update event", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) removeEventHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)

	eventID, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse id " + v["id"])
		return
	}
	err = s.app.RemoveEvent(req.Context(), storage.EventID(eventID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getEventHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)

	userID, err := strconv.Atoi(v["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse user_id " + v["user_id"])
		return
	}

	res, err := s.app.GetEventsAll(req.Context(), storage.EventOwnerID(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jres, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jres)
}

func (s *Server) getEventIntervalHandler(w http.ResponseWriter, req *http.Request,
	funcInterval func(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error),
	date time.Time,
) {
	v := mux.Vars(req)
	userID, err := strconv.Atoi(v["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("error parse user_id " + v["user_id"])
		return
	}

	res, err := funcInterval(req.Context(), storage.EventOwnerID(userID), date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jres, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jres)
}

func (s *Server) getEventForDayHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	datestring := fmt.Sprintf("%s-%s-%s", v["year"], v["month"], v["day"])
	date, err := time.Parse("2006-1-2", datestring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error(ErrorParceDate+datestring, "error", err)
		return
	}
	s.getEventIntervalHandler(w, req, s.app.GetEventsForDay, date)
}

func (s *Server) getEventForWeekHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	datestring := fmt.Sprintf("%s-%s-%s", v["year"], v["month"], v["day"])
	date, err := time.Parse("2006-1-2", datestring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error(ErrorParceDate+datestring, "error", err)
		return
	}
	s.getEventIntervalHandler(w, req, s.app.GetEventsForWeek, date)
}

func (s *Server) getEventForMonthHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	datestring := fmt.Sprintf("%s-%s-%s", v["year"], v["month"], "1")
	date, err := time.Parse("2006-1-2", datestring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error(ErrorParceDate+datestring, "error", err)
		return
	}
	s.getEventIntervalHandler(w, req, s.app.GetEventsForMonth, date)
}
