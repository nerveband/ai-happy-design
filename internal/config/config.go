package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const DefaultPort = 3055

// Config holds all configuration values for the application.
type Config struct {
	Port        int
	ServerHost  string
	IdleTimeout time.Duration // 0 = never auto-shutdown
	MCPEnabled  bool
}

// fileConfig mirrors the TOML file structure.
type fileConfig struct {
	MCP    mcpConfig    `toml:"mcp"`
	Server serverConfig `toml:"server"`
}

type mcpConfig struct {
	Enabled bool `toml:"enabled"`
}

type serverConfig struct {
	Port        int    `toml:"port"`
	Host        string `toml:"host"`
	IdleTimeout string `toml:"idle_timeout"`
}

// Dir returns the configuration directory path.
// Respects AHD_CONFIG_DIR env var, otherwise defaults to ~/.config/ai-happy-design.
func Dir() string {
	if d := os.Getenv("AHD_CONFIG_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".config", "ai-happy-design")
	}
	return filepath.Join(home, ".config", "ai-happy-design")
}

// Path returns the full path to the config.toml file.
func Path() string {
	return filepath.Join(Dir(), "config.toml")
}

// Load reads configuration from the TOML file (if present) and environment
// variables. Environment variables override TOML values for backwards
// compatibility.
func Load() *Config {
	cfg := &Config{
		Port:        DefaultPort,
		ServerHost:  "localhost",
		IdleTimeout: 15 * time.Minute,
		MCPEnabled:  false,
	}

	// Load from TOML file (silently ignore if missing or unparseable).
	var fc fileConfig
	if _, err := toml.DecodeFile(Path(), &fc); err == nil {
		cfg.MCPEnabled = fc.MCP.Enabled
		if fc.Server.Port > 0 {
			cfg.Port = fc.Server.Port
		}
		if fc.Server.Host != "" {
			cfg.ServerHost = fc.Server.Host
		}
		if fc.Server.IdleTimeout != "" {
			cfg.IdleTimeout = parseIdleTimeout(fc.Server.IdleTimeout)
		}
	}

	// Environment variables override TOML values.
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			cfg.Port = parsed
		}
	}
	if h := os.Getenv("SERVER_HOST"); h != "" {
		cfg.ServerHost = h
	}
	if t := os.Getenv("AHD_IDLE_TIMEOUT"); t != "" {
		cfg.IdleTimeout = parseIdleTimeout(t)
	}

	return cfg
}

func parseIdleTimeout(s string) time.Duration {
	if s == "0" || s == "off" || s == "disabled" {
		return 0
	}
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	if mins, err := strconv.Atoi(s); err == nil {
		return time.Duration(mins) * time.Minute
	}
	return 15 * time.Minute
}

// loadFile reads and parses the TOML config file, returning defaults if the
// file does not exist.
func loadFile() (*fileConfig, error) {
	fc := &fileConfig{
		Server: serverConfig{
			Port:        DefaultPort,
			Host:        "localhost",
			IdleTimeout: "15m",
		},
	}
	_, err := toml.DecodeFile(Path(), fc)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return fc, nil
}

// saveFile writes the fileConfig to disk as TOML, creating the directory if
// needed.
func saveFile(fc *fileConfig) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	f, err := os.Create(Path())
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()
	enc := toml.NewEncoder(f)
	return enc.Encode(fc)
}

// knownKeys maps dotted key names to their TOML section and field.
var knownKeys = map[string][2]string{
	"mcp.enabled":         {"mcp", "enabled"},
	"server.port":         {"server", "port"},
	"server.host":         {"server", "host"},
	"server.idle_timeout": {"server", "idle_timeout"},
}

// Get retrieves a configuration value by its dotted key name from the TOML
// file.
func Get(key string) (string, error) {
	if _, ok := knownKeys[key]; !ok {
		return "", fmt.Errorf("unknown config key %q. Valid keys: %s", key, validKeysString())
	}
	fc, err := loadFile()
	if err != nil {
		return "", err
	}
	switch key {
	case "mcp.enabled":
		return fmt.Sprintf("%v", fc.MCP.Enabled), nil
	case "server.port":
		return strconv.Itoa(fc.Server.Port), nil
	case "server.host":
		return fc.Server.Host, nil
	case "server.idle_timeout":
		return fc.Server.IdleTimeout, nil
	}
	return "", fmt.Errorf("unknown key %q", key)
}

// Set updates a configuration value by its dotted key name and writes the
// change to the TOML file.
func Set(key, value string) error {
	if _, ok := knownKeys[key]; !ok {
		return fmt.Errorf("unknown config key %q. Valid keys: %s", key, validKeysString())
	}
	fc, err := loadFile()
	if err != nil {
		return err
	}
	switch key {
	case "mcp.enabled":
		v := strings.ToLower(value)
		fc.MCP.Enabled = v == "true" || v == "1" || v == "yes"
	case "server.port":
		p, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid port value %q: %w", value, err)
		}
		fc.Server.Port = p
	case "server.host":
		fc.Server.Host = value
	case "server.idle_timeout":
		fc.Server.IdleTimeout = value
	}
	return saveFile(fc)
}

// Init creates a default config.toml file. If force is false and the file
// already exists, it returns an error.
func Init(force bool) error {
	path := Path()
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config file already exists at %s (use --force to overwrite)", path)
		}
	}
	fc := &fileConfig{
		MCP: mcpConfig{Enabled: false},
		Server: serverConfig{
			Port:        DefaultPort,
			Host:        "localhost",
			IdleTimeout: "15m",
		},
	}
	return saveFile(fc)
}

func validKeysString() string {
	keys := make([]string, 0, len(knownKeys))
	for k := range knownKeys {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
