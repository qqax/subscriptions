package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

// Init initializes the logger with the specified log level.
// If the log level is an empty string, it defaults to "info".
// The logger outputs to stderr in a human-readable format with timestamps and caller information.
func Init(logLevel string) error {
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(level)

	cw := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(cw).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Info().Str("log_level", logLevel).Msg("Logger initialized from env")
	return nil
}

// Trace returns a logger event with trace level.
func Trace() *zerolog.Event {
	return log.Trace()
}

// Debug returns a logger event with debug level.
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info returns a logger event with info level.
func Info() *zerolog.Event {
	return log.Info()
}

// Warn returns a logger event with warn level.
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error returns a logger event with error level.
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal returns a logger event with fatal level.
func Fatal() *zerolog.Event {
	return log.Fatal()
}

// Panic returns a logger event with panic level.
func Panic() *zerolog.Event {
	return log.Panic()
}
