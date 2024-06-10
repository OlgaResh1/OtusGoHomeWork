package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string, format string, isAddSource bool) *Logger {
	logConfig := &slog.HandlerOptions{
		AddSource:   isAddSource,
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
