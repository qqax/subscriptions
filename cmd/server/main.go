// cmd/server/main.go
package main

import (
	"net/http"
	"subscription/internal/config"
	"subscription/internal/handler"
	"subscription/internal/repository/postgres/models"

	"subscription/core/service"
	ogenServer "subscription/internal/api/generated"
	ogenAdapter "subscription/internal/handler/ogen"
	"subscription/internal/logger"
	"subscription/internal/repository/postgres"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Инициализация логгера
	if err := logger.InitSimple(config.LogLevel); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	// Подключение к БД
	dbConfig := postgres.Params{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		Name:     config.DBName,
		SSLMode:  config.SSLMode,
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

	// Create ogen server.
	server, err := ogenServer.NewServer(adapter)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create ogen server")
	}

	// Add middleware
	httpHandler := handler.AddMiddleware(server)

	// Запуск сервера
	logger.Info().Msgf("Starting server on %s:%s", config.ServerHost, config.ServerPort)
	if err = http.ListenAndServe(config.ServerHost+":"+config.ServerPort, httpHandler); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}
