package config

import (
	"github.com/joho/godotenv"
	"subscription/logger"
)

var (
	LogLevel = "info"

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
)

func Load() error {
	if err := godotenv.Load(); err != nil {
		logger.Error().Msg("No .env file found, using environment variables")
	}

	LogLevel = optionalEnvStr("LOG_LEVEL", LogLevel)
	DBHost = mustEnvStr("DB_HOST")
	DBPort = optionalEnvStr("DB_PORT", DBPort)
	DBUser = optionalEnvStr("DB_USER", DBUser)
	DBPassword = optionalEnvStr("DB_PASSWORD", DBPassword)
	DBName = optionalEnvStr("DB_NAME", DBName)

	return nil
}
