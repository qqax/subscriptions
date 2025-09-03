// cmd/server/main.go
package main

import (
	"net/http"
	"os"
	"subscription/internal/repository/postgres/models"

	"subscription/core/service"
	ogenServer "subscription/internal/api/generated"
	ogenAdapter "subscription/internal/handler/ogen"
	"subscription/internal/logger"
	"subscription/internal/repository/postgres"

	"github.com/rs/zerolog/log"
)

func main() {
	// Инициализация логгера
	if err := logger.InitSimple("info"); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	// Подключение к БД
	dbConfig := postgres.Params{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "server"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dbClient, err := postgres.NewClient(dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbClient.Close()

	// Миграции
	if err = dbClient.Migrate(&models.Subscription{}); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Репозиторий
	repo := postgres.NewSubscriptionRepository(dbClient.DB)

	// Сервис (ядро)
	subscriptionService := service.NewSubscriptionService(repo)

	// Ogen adapter
	adapter := ogenAdapter.NewOgenAdapter(subscriptionService)

	// Create ogen server
	server, err := ogenServer.NewServer(adapter)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create ogen server")
	}

	// Add middleware
	httpHandler := addMiddleware(server)

	// Запуск сервера
	logger.Info().Msg("Starting server on :8000")
	if err = http.ListenAndServe(":8000", httpHandler); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

func addMiddleware(handler http.Handler) http.Handler {
	// Add your middleware here
	return handler
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
