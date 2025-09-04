// cmd/server/main.go
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"subscription/internal/config"
	"subscription/internal/handler"
	"subscription/internal/repository/postgres/models"
	"syscall"
	"time"

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

	// Migrations
	if err = dbClient.Migrate(&models.Subscription{}); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Repository
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

	// Add middlewares
	httpHandler := handler.AddMiddleware(server)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/live", handler.LiveCheckHandler)
	mux.HandleFunc("/ready", handler.ReadyCheckHandler(dbClient))
	mux.HandleFunc("/admin/db-stats", handler.DBStatsHandler(dbClient))
	mux.Handle("/", httpHandler)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         config.ServerHost + ":" + config.ServerPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel for graceful shutdown
	shutdownChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		logger.Info().Str("address", srv.Addr).Msgf("Starting server on %s:%s", config.ServerHost, config.ServerPort)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			shutdownChan <- err
		}
	}()

	// Канал для системных сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Ожидаем сигнал завершения или ошибку сервера
	select {
	case sig := <-sigChan:
		logger.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err = <-shutdownChan:
		logger.Error().Err(err).Msg("Server error occurred")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Info().Msg("Shutting down server gracefully...")

	// Останавливаем сервер
	if err = srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to shutdown server gracefully")
	} else {
		logger.Info().Msg("Server stopped gracefully")
	}

	// Проверяем есть ли еще ошибки
	select {
	case err = <-shutdownChan:
		if err != nil {
			logger.Error().Err(err).Msg("Additional error occurred during shutdown process")
		} else {
			logger.Info().Msg("Shutdown process completed")
		}
	default:
		logger.Info().Msg("Shutdown completed cleanly without any additional errors")
	}
}
