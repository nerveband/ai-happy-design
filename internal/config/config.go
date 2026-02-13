package config

import (
	"os"
	"strconv"
)

// Config holds all configuration values for the application.
type Config struct {
	Port       int
	ServerHost string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	port := 3055
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	host := "localhost"
	if h := os.Getenv("SERVER_HOST"); h != "" {
		host = h
	}

	return &Config{
		Port:       port,
		ServerHost: host,
	}
}
