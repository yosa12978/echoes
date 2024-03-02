package logging

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Fields(fields map[string]interface{})
	Error(err error)
	Debug(msg string)
	Info(msg string)
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

type cmdLogger struct {
	logger zerolog.Logger
}

func New(service string) Logger {
	l := new(cmdLogger)
	l.logger = zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Str("service", service).
		Logger()
	return l
}

func (l *cmdLogger) Fields(fields map[string]interface{}) {
	l.logger.Info().Fields(fields).Send()
}

func (l *cmdLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *cmdLogger) Error(err error) {
	l.logger.Err(err).Send()
}

func (l *cmdLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *cmdLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *cmdLogger) Printf(format string, v ...interface{}) {
	l.logger.Info().Msg(fmt.Sprintf(format, v...))
}

func (l *cmdLogger) Trace(msg string) {
	l.logger.Trace().Msg(msg)
}
