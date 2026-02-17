package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("AHD_IDLE_TIMEOUT")
	os.Unsetenv("AHD_CONFIG_DIR")

	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	cfg := Load()
	if cfg.Port != 3055 {
		t.Errorf("expected port 3055, got %d", cfg.Port)
	}
	if cfg.ServerHost != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.ServerHost)
	}
	if cfg.MCPEnabled {
		t.Error("expected MCP disabled by default")
	}
}

func TestLoadFromTOML(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("AHD_IDLE_TIMEOUT")

	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	toml := `[mcp]
enabled = true

[server]
port = 4000
host = "example.com"
idle_timeout = "5m"
`
	os.WriteFile(filepath.Join(tmp, "config.toml"), []byte(toml), 0644)

	cfg := Load()
	if !cfg.MCPEnabled {
		t.Error("expected MCP enabled from TOML")
	}
	if cfg.Port != 4000 {
		t.Errorf("expected port 4000, got %d", cfg.Port)
	}
	if cfg.ServerHost != "example.com" {
		t.Errorf("expected host example.com, got %s", cfg.ServerHost)
	}
}

func TestEnvOverridesToml(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	toml := `[server]
port = 4000
`
	os.WriteFile(filepath.Join(tmp, "config.toml"), []byte(toml), 0644)

	os.Setenv("PORT", "5000")
	defer os.Unsetenv("PORT")

	cfg := Load()
	if cfg.Port != 5000 {
		t.Errorf("expected env PORT=5000 to override TOML port=4000, got %d", cfg.Port)
	}
}

func TestConfigDir(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	dir := Dir()
	if dir != tmp {
		t.Errorf("expected dir %s, got %s", tmp, dir)
	}
}

func TestConfigPath(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	p := Path()
	expected := filepath.Join(tmp, "config.toml")
	if p != expected {
		t.Errorf("expected path %s, got %s", expected, p)
	}
}

func TestSetAndGet(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	if err := Set("mcp.enabled", "true"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := Get("mcp.enabled")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "true" {
		t.Errorf("expected 'true', got %q", val)
	}

	cfg := Load()
	if !cfg.MCPEnabled {
		t.Error("expected MCP enabled after Set")
	}
}

func TestSetServerPort(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")
	os.Unsetenv("PORT")

	if err := Set("server.port", "9999"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	cfg := Load()
	if cfg.Port != 9999 {
		t.Errorf("expected port 9999, got %d", cfg.Port)
	}
}

func TestGetUnknownKey(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	_, err := Get("nonexistent.key")
	if err == nil {
		t.Error("expected error for unknown key")
	}
}

func TestSetUnknownKey(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	err := Set("nonexistent.key", "value")
	if err == nil {
		t.Error("expected error for unknown key")
	}
}

func TestInitCreatesDefaultFile(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	if err := Init(false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	path := filepath.Join(tmp, "config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("config file is empty")
	}
}

func TestInitRefusesOverwrite(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	Init(false)

	err := Init(false)
	if err == nil {
		t.Error("expected error when config already exists")
	}
}

func TestInitForceOverwrites(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("AHD_CONFIG_DIR", tmp)
	defer os.Unsetenv("AHD_CONFIG_DIR")

	Init(false)
	err := Init(true)
	if err != nil {
		t.Errorf("expected force init to succeed, got: %v", err)
	}
}
