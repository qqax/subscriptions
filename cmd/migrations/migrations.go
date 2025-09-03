package main

import (
	"gorm.io/gen"
	"gorm.io/gorm"
	"subscription/config"
	"subscription/internal/adapters/db"
	"subscription/internal/adapters/db/models"
	"subscription/pkg/logger"
)

// main initializes configuration, logging, and database connectivity,
// then proceeds to generate GORM models and perform schema migrations.
// The application will terminate if any of these steps fail.
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

// generateModels uses GORM's code generator to create model definitions
// from the domain models. It outputs generated code under the
// "./pkg/db/tables" directory using options suitable for tests,
// interfaces, and default queries.
func generateModels(client *db.Client) error {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./pkg/db/tables",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(client.DB)

	// Include core domain models for code generation
	g.ApplyBasic(models.Service{}, &models.Price{}, &models.ServiceStatus{}, models.Subscription{})

	g.Execute()
	return nil
}

// migrateDatabase applies schema migrations using GORM's AutoMigrate method.
// It ensures that the database schema matches the models for Service, Price,
// ServiceStatus, and Subscription. Any migration failure is returned.
func migrateDatabase(gdb *gorm.DB) error {
	if err := gdb.AutoMigrate(
		&models.Service{}, &models.Price{}, &models.ServiceStatus{}, &models.Subscription{},
	); err != nil {
		return err
	}

	return nil
}
