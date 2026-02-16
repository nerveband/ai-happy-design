package config

import (
	"os"
	"strconv"
	"time"
)

const DefaultPort = 3055

// Config holds all configuration values for the application.
type Config struct {
	Port        int
	ServerHost  string
	IdleTimeout time.Duration // 0 = never auto-shutdown
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	port := DefaultPort
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	host := "localhost"
	if h := os.Getenv("SERVER_HOST"); h != "" {
		host = h
	}

	idleTimeout := 15 * time.Minute
	if t := os.Getenv("AHD_IDLE_TIMEOUT"); t != "" {
		if t == "0" || t == "off" || t == "disabled" {
			idleTimeout = 0
		} else if d, err := time.ParseDuration(t); err == nil {
			idleTimeout = d
		} else if mins, err := strconv.Atoi(t); err == nil {
			idleTimeout = time.Duration(mins) * time.Minute
		}
	}

	return &Config{
		Port:        port,
		ServerHost:  host,
		IdleTimeout: idleTimeout,
	}
}
