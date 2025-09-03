package config

import (
	"github.com/joho/godotenv"
	"subscription/internal/logger"
)

var (
	LogLevel = "info"

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
)

// Load initializes the application's configuration by loading environment variables.
// It attempts to load variables from a .env file and then overrides them with
// environment variables if available. Missing required variables will cause the
// application to panic.
func Load() error {
	// Attempt to load variables from a .env file
	if err := godotenv.Load(); err != nil {
		logger.Error().Msg("No .env file found, using environment variables")
	}

	// Load and set environment variables with fallback values
	LogLevel = optionalEnvStr("LOG_LEVEL", LogLevel)
	DBHost = mustEnvStr("DB_HOST")
	DBPort = optionalEnvStr("DB_PORT", DBPort)
	DBUser = optionalEnvStr("DB_USER", DBUser)
	DBPassword = optionalEnvStr("DB_PASSWORD", DBPassword)
	DBName = optionalEnvStr("DB_NAME", DBName)

	return nil
}
