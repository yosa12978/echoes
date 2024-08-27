package logging

import (
	"io"
	"log/slog"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type slogLogger struct {
	logger *slog.Logger
}

func NewJsonLogger(w io.Writer) Logger {
	return &slogLogger{
		logger: slog.New(slog.NewJSONHandler(w, nil)),
	}
}

func NewTextLogger(w io.Writer) Logger {
	return &slogLogger{
		logger: slog.New(slog.NewTextHandler(w, nil)),
	}
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}
