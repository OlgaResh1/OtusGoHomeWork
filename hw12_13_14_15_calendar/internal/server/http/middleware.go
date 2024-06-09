package internalhttp

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *Server) logHTTPRequest(r *http.Request, d time.Duration, statusCode int) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("error split host and port: %v", err)
		return
	}
	t := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	msg := fmt.Sprintf("%s [%s] %s %s %s %d %d %s", ip, t, r.Method, r.URL.String(), r.Proto, statusCode,
		d.Microseconds(), r.UserAgent())

	s.logger.Info(msg, "source", "http")
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		s.logHTTPRequest(r, duration, lrw.statusCode)
	})
}
