// cmd/generate/main.go
package main

import (
	"log"
	"os"

	"subscription/internal/logger"
	"subscription/internal/repository/postgres"
	"subscription/internal/repository/postgres/models"

	"gorm.io/gen"
	"gorm.io/gorm"
)

// main генерирует GORM модели из структуры БД
func main() {
	// Инициализация логгера
	if err := logger.InitSimple("info"); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Получение конфигурации из переменных окружения
	dbConfig := postgres.Params{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "subscriptions"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Подключение к БД
	client, err := postgres.NewClient(dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer client.Close()

	// Генерация моделей
	if err := generateModels(client.DB); err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate models")
	}

	// Миграции (опционально)
	if err := migrateDatabase(client.DB); err != nil {
		logger.Fatal().Err(err).Msg("Failed to migrate database")
	}

	logger.Info().Msg("Model generation completed successfully")
}

// generateModels генерирует код моделей GORM
func generateModels(db *gorm.DB) error {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/repository/postgres/generated",
		ModelPkgPath:  "github.com/your-org/your-app/internal/repository/postgres/generated",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)

	// Регистрируем модели для генерации
	g.ApplyBasic(
		models.Subscription{},
		// Добавьте другие модели по необходимости
	)

	g.Execute()
	return nil
}

// migrateDatabase применяет миграции
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Subscription{},
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
