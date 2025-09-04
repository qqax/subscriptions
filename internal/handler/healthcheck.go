package handler

import (
	"encoding/json"
	"net/http"
	"subscription/internal/repository/postgres"
	"time"

	"subscription/internal/logger"
)

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "subscription-api",
	})

	logger.Debug().Msg("Health check passed")
}

// ReadyCheckHandler проверяет готовность всех зависимостей
func ReadyCheckHandler(dbClient *postgres.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now().Format(time.RFC3339),
			"services":  map[string]interface{}{},
		}

		// Проверяем соединение с БД
		if err := dbClient.HealthCheck(); err != nil {
			logger.Error().Err(err).Msg("Database health check failed")
			response["status"] = "degraded"
			response["services"].(map[string]interface{})["database"] = map[string]interface{}{
				"status":  "disconnected",
				"error":   err.Error(),
				"details": "Failed to ping database",
			}
		} else {
			response["services"].(map[string]interface{})["database"] = map[string]interface{}{
				"status": "connected",
				"details": map[string]interface{}{
					"type":    "postgresql",
					"version": "13", // Можно получить реальную версию
				},
			}
		}

		// Устанавливаем соответствующий HTTP статус
		if response["status"] == "degraded" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// LiveCheckHandler проверяет что приложение работает (без зависимостей)
func LiveCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "live",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// DBStatsHandler возвращает статистику БД (для админов)
func DBStatsHandler(dbClient *postgres.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		stats := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"database":  map[string]interface{}{},
		}

		// Получаем статистику подключений
		sqlDB, err := dbClient.DB.DB()
		if err == nil {
			stats["database"].(map[string]interface{})["connections"] = map[string]interface{}{
				"open":           sqlDB.Stats().OpenConnections,
				"in_use":         sqlDB.Stats().InUse,
				"idle":           sqlDB.Stats().Idle,
				"wait_count":     sqlDB.Stats().WaitCount,
				"wait_duration":  sqlDB.Stats().WaitDuration.String(),
				"max_open_conns": sqlDB.Stats().MaxOpenConnections,
			}
		}

		// Проверяем доступность БД
		if err = dbClient.HealthCheck(); err != nil {
			stats["database"].(map[string]interface{})["status"] = "down"
			stats["database"].(map[string]interface{})["error"] = err.Error()
		} else {
			stats["database"].(map[string]interface{})["status"] = "up"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}
