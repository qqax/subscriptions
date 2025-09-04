package config

import (
	"github.com/joho/godotenv"
	"subscription/internal/logger"
)

const (
	DefaultLogLevel = "info"
	DefaultSSLMode  = "disable"
)

var (
	LogLevel string

	ServerPort string
	ServerHost string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
)

// Load initializes the application's configuration by loading environment variables.
// It attempts to load variables from a .env.local file and then overrides them with
// environment variables if available. Missing required variables will cause the
// application to panic.
func Load() error {
	// Attempt to load variables from a .env.local file
	if err := godotenv.Load(); err != nil {
		logger.Error().Msg("No .env.local file found, using environment variables")
	}

	// Load and set environment variables with fallback values
	LogLevel = optionalEnvStr("LOG_LEVEL", DefaultLogLevel)

	DBHost = mustEnvStr("DB_HOST")
	DBPort = mustEnvStr("DB_PORT")
	DBUser = mustEnvStr("DB_USER")
	DBPassword = mustEnvStr("DB_PASSWORD")
	DBName = mustEnvStr("DB_NAME")
	SSLMode = optionalEnvStr("DB_SSLMODE", DefaultSSLMode)

	ServerHost = mustEnvStr("SERVER_HOST")
	ServerPort = mustEnvStr("SERVER_PORT")

	return nil
}
