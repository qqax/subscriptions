package postgres

import (
	"context"
	"errors"
	"time"

	gormLogger "gorm.io/gorm/logger"
	appLogger "subscription/internal/logger" // Алиас для вашего логгера
)

// ZerologLogger implements gorm.Logger.Interface for zerolog
type ZerologLogger struct {
	LogLevel gormLogger.LogLevel
}

// NewZerologLogger creates a new GORM logger adapter for zerolog
func NewZerologLogger() gormLogger.Interface {
	return &ZerologLogger{
		LogLevel: gormLogger.Warn, // Default level
	}
}

// LogMode set log mode
func (l *ZerologLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *ZerologLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		appLogger.Info().Msgf(msg, data...)
	}
}

// Warn logs warn messages
func (l *ZerologLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		appLogger.Warn().Msgf(msg, data...)
	}
}

// Error logs error messages
func (l *ZerologLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		appLogger.Error().Msgf(msg, data...)
	}
}

// Trace logs SQL queries, timing, and errors
func (l *ZerologLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound):
		appLogger.Error().
			Err(err).
			Str("sql", sql).
			Int64("rows", rows).
			Dur("elapsed", elapsed).
			Msg("SQL Error")
	case elapsed > time.Second: // Slow query
		appLogger.Warn().
			Str("sql", sql).
			Int64("rows", rows).
			Dur("elapsed", elapsed).
			Msg("Slow SQL Query")
	case l.LogLevel == gormLogger.Info:
		appLogger.Debug().
			Str("sql", sql).
			Int64("rows", rows).
			Dur("elapsed", elapsed).
			Msg("SQL Query")
	}
}
