package main

import (
	"subscription/config"
	"subscription/logger"
)

func main() {
	err := config.Load()
	if err != nil {
		logger.Fatal().Msgf("failed to load env file: %v", err)
	}

	err = logger.Init(config.LogLevel)
	if err != nil {
		logger.Fatal().Msgf("failed to init logger: %v", err)
	}

}
