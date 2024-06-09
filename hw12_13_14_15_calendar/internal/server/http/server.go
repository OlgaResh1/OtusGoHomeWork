package internalhttp

import (
	"context"
	"net/http"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
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

type Application interface { // TODO
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello\n"))
}

func NewServer(cfg config.Config, logger Logger, app Application) *Server {
	srv := &http.Server{
		Addr:              cfg.Http.Address,
		ReadHeaderTimeout: cfg.Http.RequestTimeout,
	}
	return &Server{logger: logger, app: app, server: srv}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	mux.HandleFunc("/", s.helloHandler)

	s.server.Handler = s.loggingMiddleware(mux)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
