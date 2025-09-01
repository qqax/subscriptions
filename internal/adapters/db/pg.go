package db

import (
	"fmt"
	"subscription/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Client wraps the GORM DB instance.
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
}

// NewClient creates a new database connection using the provided parameters.
func NewClient(p Params) (*Client, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		p.Host, p.User, p.Password, p.Name, p.Port,
	)
	logger.Debug().
		Str("host", p.Host).
		Str("port", p.Port).
		Str("dbname", p.Name).
		Msg("Connecting to database")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	logger.Info().Msg("Database connection established")

	return &Client{DB: db}, nil
}

// Migrate runs schema migrations for the given models.
func (c *Client) Migrate(models ...interface{}) error {
	if err := c.AutoMigrate(models...); err != nil {
		logger.Error().
			Err(err).
			Msg("schema migration failed")
		return fmt.Errorf("schema migration: %w", err)
	}

	logger.Info().Msg("schema migration completed successfully")
	return nil
}
