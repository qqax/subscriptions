package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	appLogger "subscription/internal/logger" // Алиас для вашего логгера
)

// Client wraps the GORM DB instance with connection management.
type Client struct {
	*gorm.DB
}

// Params holds settings for the database connection.
type Params struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// NewClient creates a new database connection with zerolog integration.
func NewClient(p Params) (*Client, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		p.Host, p.User, p.Password, p.Name, p.Port, p.SSLMode,
	)

	appLogger.Debug().
		Str("host", p.Host).
		Str("port", p.Port).
		Str("dbname", p.Name).
		Msg("Connecting to database")

	// GORM config with zerolog adapter
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   NewZerologLogger(),
	}

	// Connection with retries
	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			// Test connection
			if sqlDB, pingErr := db.DB(); pingErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					break
				}
			}
		}

		if i < maxRetries-1 {
			appLogger.Warn().
				Err(err).
				Int("attempt", i+1).
				Msg("Failed to connect to database, retrying...")
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	appLogger.Info().Msg("Database connection established successfully")

	return &Client{DB: db}, nil
}

// HealthCheck проверяет соединение с БД.
func (c *Client) HealthCheck() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("getting sql.DB: %w", err)
	}
	return sqlDB.Ping()
}

func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		appLogger.Error().Err(err).Msg("Failed to get underlying SQL DB for closing")
		return fmt.Errorf("getting sql.DB: %w", err)
	}

	appLogger.Info().Msg("Closing database connection")

	if err = sqlDB.Close(); err != nil {
		appLogger.Error().Err(err).Msg("Failed to close database connection")
		return fmt.Errorf("closing database connection: %w", err)
	}

	appLogger.Info().Msg("Database connection closed successfully")
	return nil
}

// Migrate runs schema migrations for the given models.
func (c *Client) Migrate(models ...interface{}) error {
	appLogger.Info().Msg("Starting database migration")

	if err := c.DB.AutoMigrate(models...); err != nil {
		appLogger.Error().
			Err(err).
			Msg("Schema migration failed")
		return fmt.Errorf("schema migration: %w", err)
	}

	appLogger.Info().Msg("Schema migration completed successfully")
	return nil
}

// WithTx executes a function within a transaction.
func (c *Client) WithTx(fn func(tx *gorm.DB) error) error {
	return c.DB.Transaction(fn)
}
