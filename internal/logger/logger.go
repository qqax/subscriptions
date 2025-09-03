package logger

import (
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	once     sync.Once
	instance *zerolog.Logger
)

// Config holds logger configuration
type Config struct {
	Level      string
	Output     string // "console" or "json"
	TimeFormat string
	Caller     bool
	Color      bool
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Output:     "console",
		TimeFormat: time.RFC3339,
		Caller:     true,
		Color:      true,
	}
}

// Init initializes the global logger with configuration
func Init(cfg Config) error {
	var initErr error
	once.Do(func() {
		level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
		if err != nil {
			initErr = err
			return
		}
		zerolog.SetGlobalLevel(level)

		// Set time field format
		zerolog.TimeFieldFormat = cfg.TimeFormat

		if cfg.Output == "json" {
			// JSON output
			log.Logger = zerolog.New(os.Stderr).
				With().
				Timestamp().
				Logger()
		} else {
			// Console output with customization
			writer := zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: cfg.TimeFormat,
				NoColor:    !cfg.Color,
			}
			log.Logger = zerolog.New(writer).
				With().
				Timestamp().
				Logger()
		}

		// Add caller information if enabled
		if cfg.Caller {
			log.Logger = log.Logger.With().Caller().Logger()
		}

		// Add runtime information
		log.Logger = log.Logger.With().
			Str("go_version", runtime.Version()).
			Logger()

		instance = &log.Logger

		Info().
			Str("level", cfg.Level).
			Str("output", cfg.Output).
			Msg("Logger initialized successfully")
	})
	return initErr
}

// InitSimple initializes logger with simple string level
func InitSimple(logLevel string) error {
	return Init(Config{
		Level:  logLevel,
		Output: "console",
		Caller: true,
		Color:  true,
	})
}

// Get returns the logger instance
func Get() *zerolog.Logger {
	if instance == nil {
		// Initialize with defaults if not initialized
		Init(DefaultConfig())
	}
	return instance
}

// With creates a new logger with additional context
func With() zerolog.Context {
	return Get().With()
}

// Trace returns a logger event with trace level.
func Trace() *zerolog.Event {
	return Get().Trace()
}

// Debug returns a logger event with debug level.
func Debug() *zerolog.Event {
	return Get().Debug()
}

// Info returns a logger event with info level.
func Info() *zerolog.Event {
	return Get().Info()
}

// Warn returns a logger event with warn level.
func Warn() *zerolog.Event {
	return Get().Warn()
}

// Error returns a logger event with error level.
func Error() *zerolog.Event {
	return Get().Error()
}

// Fatal returns a logger event with fatal level.
func Fatal() *zerolog.Event {
	return Get().Fatal()
}

// Panic returns a logger event with panic level.
func Panic() *zerolog.Event {
	return Get().Panic()
}

// Log returns a logger event with specified level.
func Log() *zerolog.Event {
	return Get().Log()
}

// Sync flushes any buffered log entries
func Sync() error {
	// zerolog doesn't require explicit flushing for most writers
	return nil
}

// Helper functions for common logging patterns

// WithError creates a logger with error context
func WithError(err error) *zerolog.Logger {
	logger := Get().With().Err(err).Logger()
	return &logger
}

// WithField creates a logger with a single field
func WithField(key, value string) *zerolog.Logger {
	logger := Get().With().Str(key, value).Logger()
	return &logger
}

// WithFields creates a logger with multiple fields
func WithFields(fields map[string]interface{}) *zerolog.Logger {
	context := Get().With()
	for key, value := range fields {
		context = context.Interface(key, value)
	}
	logger := context.Logger()
	return &logger
}

// Business domain specific helpers

// WithRequestID creates a logger with request ID
func WithRequestID(requestID string) *zerolog.Logger {
	return WithField("request_id", requestID)
}

// WithUserID creates a logger with user ID
func WithUserID(userID string) *zerolog.Logger {
	return WithField("user_id", userID)
}

// WithService creates a logger with service name
func WithService(service string) *zerolog.Logger {
	return WithField("service", service)
}

// Alternative approach: Simple helper methods that return the event

// WithDuration adds a duration field to the event
func WithDuration(e *zerolog.Event, d time.Duration) *zerolog.Event {
	return e.Dur("duration", d)
}

// WithTimestamp adds a timestamp field to the event
func WithTimestamp(e *zerolog.Event, t time.Time) *zerolog.Event {
	return e.Time("timestamp", t)
}

// WithBool adds a boolean field to the event
func WithBool(e *zerolog.Event, key string, value bool) *zerolog.Event {
	return e.Bool(key, value)
}

// WithInt adds an integer field to the event
func WithInt(e *zerolog.Event, key string, value int) *zerolog.Event {
	return e.Int(key, value)
}

// WithString adds a string field to the event
func WithString(e *zerolog.Event, key, value string) *zerolog.Event {
	return e.Str(key, value)
}

// WithFloat adds a float field to the event
func WithFloat(e *zerolog.Event, key string, value float64) *zerolog.Event {
	return e.Float64(key, value)
}

// WithErr adds an error field to the event
func WithErr(e *zerolog.Event, err error) *zerolog.Event {
	return e.Err(err)
}

// LogDuration helper for common duration logging pattern
func LogDuration(d time.Duration) {
	Get().Info().Dur("duration", d).Msg("Operation duration")
}

// LogOperationDuration times and logs an operation
func LogOperationDuration(operationName string, start time.Time) {
	duration := time.Since(start)
	Get().Info().
		Str("operation", operationName).
		Dur("duration_ms", duration).
		Msg("Operation completed")
}
