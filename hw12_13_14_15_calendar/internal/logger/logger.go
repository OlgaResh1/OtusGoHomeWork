package logger

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string, format string, isAddSource bool) *Logger {
	fmt.Println(level, format)

	logConfig := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}
	logHandler := slog.NewTextHandler(os.Stderr, logConfig)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return &Logger{logger: logger}
}

func (l Logger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (l Logger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// 66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
func (l Logger) LogHTTPRequest(r *http.Request, d time.Duration, statusCode int) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("error split host and port: %v", err)
		return
	}

	t := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	log.Printf("%s [%s] %s %s %s %d %d %s", ip, t, r.Method, r.URL.String(), r.Proto, statusCode, d.Microseconds(), r.UserAgent())
}
