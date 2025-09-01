package main

import (
	"gorm.io/gorm"
	"subscription/config"
	"subscription/db"
	"subscription/logger"
	"subscription/models"

	"gorm.io/gen"
)

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal().Err(err).Msg("loading configuration")
	}

	if err := logger.Init(config.LogLevel); err != nil {
		logger.Fatal().Err(err).Msg("initializing logger")
	}

	client, err := db.NewClient(db.Params{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		Name:     config.DBName,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("connecting to database")
	}

	if err = generateModels(client); err != nil {
		logger.Fatal().Err(err).Msg("generating database models")
	}

	if err = migrateDatabase(client.DB); err != nil {
		logger.Fatal().Err(err).Msg("migrating database")
	}

	logger.Info().Msg("Model generation and migration completed successfully")
}

func generateModels(client *db.Client) error {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./pkg/db/tables",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(client.DB)

	g.ApplyBasic(models.Service{}, models.Subscription{})

	g.Execute()
	return nil
}

func migrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Service{}, &models.Subscription{},
	); err != nil {
		return err
	}

	return nil
}
