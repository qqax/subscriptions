package config

import (
	"github.com/rs/zerolog/log"
	"os"
)

func optionalEnvStr(key, fallback string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}

	return fallback
}

func mustEnvStr(key string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	} else {
		log.Panic().Msgf("Required environment variable %s not found", key)
		return ""
	}
}
