package config

import (
	"github.com/rs/zerolog/log"
	"os"
)

// optionalEnvStr retrieves the value of the environment variable named by key.
// If the variable is present in the environment, its value is returned.
// Otherwise, the fallback value is returned.
func optionalEnvStr(key, fallback string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return fallback
}

// mustEnvStr retrieves the value of the environment variable named by key.
// If the variable is present in the environment, its value is returned.
// If the variable is not present, the function logs a fatal error and terminates the program.
func mustEnvStr(key string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	} else {
		log.Panic().Msgf("Required environment variable %s not found", key)
		return "" // This line is never reached.
	}
}
