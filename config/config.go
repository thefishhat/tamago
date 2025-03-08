package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration for the editor and all its components.
type Config struct {
	// Addr is the address the server will listen on.
	Addr string `envconfig:"SERVER_URL" default:":8080"`
}

// LoadConfig loads the configuration from environment variables or the .env file.
// It uses the "TAMAGO_" prefix for all environment variables.
// The default values are set in the [Config] struct.
func LoadConfig() Config {
	_ = godotenv.Load()

	var cfg Config
	err := envconfig.Process("TAMAGO_", &cfg)
	if err != nil {
		log.Fatalf("Failed to process env vars: %v", err)
	}

	return cfg
}
