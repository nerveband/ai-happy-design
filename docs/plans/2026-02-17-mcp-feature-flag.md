# MCP Feature Flag Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Gate the MCP server behind a TOML config file flag (`mcp.enabled = false` by default), so users are directed to the CLI. The `register` command auto-enables MCP.

**Architecture:** Add a TOML config file at `~/.config/ai-happy-design/config.toml`. Extend `internal/config/` to read it, with env vars overriding file values. Add a `config` CLI subcommand for get/set/init/path. Gate `mcp` command on `mcp.enabled`. Make `register` auto-enable MCP before registering.

**Tech Stack:** Go, `github.com/BurntSushi/toml`, cobra CLI

---

### Task 1: Add TOML dependency

**Files:**
- Modify: `go.mod`

**Step 1: Add the dependency**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go get github.com/BurntSushi/toml@latest
```

Expected: `go.mod` updated with `github.com/BurntSushi/toml`

**Step 2: Verify**

Run: `grep BurntSushi go.mod`
Expected: Line containing `github.com/BurntSushi/toml`

**Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "deps: add BurntSushi/toml for config file support"
```

---

### Task 2: Extend config package with TOML file support

**Files:**
- Modify: `internal/config/config.go`
- Create: `internal/config/config_test.go`

**Step 1: Write the failing tests**

Create `internal/config/config_test.go`:

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Ensure no env vars interfere
	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("AHD_IDLE_TIMEOUT")
	os.Unsetenv("AHD_CONFIG_DIR")

	// Point config dir to a temp dir with no config file
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

	// Set a value
	if err := Set("mcp.enabled", "true"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get it back
	val, err := Get("mcp.enabled")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "true" {
		t.Errorf("expected 'true', got %q", val)
	}

	// Verify it persisted to file
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

	// Create first
	Init(false)

	// Try again without force
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
```

**Step 2: Run tests to verify they fail**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/config/ -v -count=1
```

Expected: Compilation errors (MCPEnabled, Dir, Path, Set, Get, Init don't exist yet)

**Step 3: Implement the config package**

Rewrite `internal/config/config.go` to:

```go
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

// Dir returns the config directory path. Respects AHD_CONFIG_DIR env var.
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

// Path returns the full path to the config file.
func Path() string {
	return filepath.Join(Dir(), "config.toml")
}

// Load reads configuration from the TOML file, then applies env var overrides.
func Load() *Config {
	cfg := &Config{
		Port:        DefaultPort,
		ServerHost:  "localhost",
		IdleTimeout: 15 * time.Minute,
		MCPEnabled:  false,
	}

	// Read TOML file
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

	// Env vars override file values
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

// loadFile reads the TOML config file, returning defaults if missing.
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

// saveFile writes the fileConfig to disk as TOML.
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

// Known config keys and their section.field mapping.
var knownKeys = map[string][2]string{
	"mcp.enabled":      {"mcp", "enabled"},
	"server.port":      {"server", "port"},
	"server.host":      {"server", "host"},
	"server.idle_timeout": {"server", "idle_timeout"},
}

// Get retrieves a config value by dotted key.
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

// Set updates a config value by dotted key and persists to disk.
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

// Init creates a default config file. Returns error if file exists and force is false.
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
```

**Step 4: Run tests to verify they pass**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/config/ -v -count=1
```

Expected: All tests PASS

**Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add TOML config file support with Get/Set/Init"
```

---

### Task 3: Add `config` CLI subcommand

**Files:**
- Create: `cmd/ai-happy-design/config_cmd.go`
- Modify: `cmd/ai-happy-design/main.go` (add command registration)

**Step 1: Create the config command file**

Create `cmd/ai-happy-design/config_cmd.go`:

```go
package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify ai-happy-design configuration stored at ~/.config/ai-happy-design/config.toml`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Long: `Get a config value by dotted key.

Valid keys: mcp.enabled, server.port, server.host, server.idle_timeout`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := config.Get(args[0])
		if err != nil {
			return err
		}
		fmt.Println(val)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long: `Set a config value by dotted key and persist to disk.

Valid keys: mcp.enabled, server.port, server.host, server.idle_timeout

Examples:
  ai-happy-design config set mcp.enabled true
  ai-happy-design config set server.port 4000`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Set(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("%s = %s\n", args[0], args[1])
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Path())
	},
}

var configInitForce bool

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default config file",
	Long:  `Creates the config file with default values. Use --force to overwrite an existing file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Init(configInitForce); err != nil {
			return err
		}
		fmt.Printf("Config created at %s\n", config.Path())
		return nil
	},
}
```

**Step 2: Register the config command in main.go**

In `cmd/ai-happy-design/main.go`, add to the `main()` function, after the existing `rootCmd.AddCommand(...)` calls and before `rootCmd.Execute()`:

```go
	configInitCmd.Flags().BoolVar(&configInitForce, "force", false, "Overwrite existing config file")
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)
```

**Step 3: Verify it compiles**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go build ./cmd/ai-happy-design/
```

Expected: No errors

**Step 4: Commit**

```bash
git add cmd/ai-happy-design/config_cmd.go cmd/ai-happy-design/main.go
git commit -m "feat: add config CLI subcommand (get/set/path/init)"
```

---

### Task 4: Gate the `mcp` command on `mcp.enabled`

**Files:**
- Modify: `cmd/ai-happy-design/main.go` (the `mcpCmd` RunE)

**Step 1: Add the gate check**

In the `mcpCmd` `RunE` function (around line 49), add the MCP enabled check at the very top, before `loadConfig()`:

Replace the existing `mcpCmd` RunE body:
```go
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadConfig()
```

With:
```go
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadConfig()

		if !cfg.MCPEnabled {
			fmt.Fprintln(os.Stderr, "MCP server is disabled by default. To enable it:")
			fmt.Fprintln(os.Stderr, "  ai-happy-design config set mcp.enabled true")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Or register with an editor (auto-enables MCP):")
			fmt.Fprintln(os.Stderr, "  ai-happy-design register")
			os.Exit(1)
		}
```

**Step 2: Verify it compiles**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go build ./cmd/ai-happy-design/
```

Expected: No errors

**Step 3: Commit**

```bash
git add cmd/ai-happy-design/main.go
git commit -m "feat: gate mcp command on mcp.enabled config flag"
```

---

### Task 5: Make `register` auto-enable MCP

**Files:**
- Modify: `cmd/ai-happy-design/register.go`

**Step 1: Add auto-enable at the start of runRegister**

At the top of `runRegister()` (after getting `binaryPath` and `home`, around line 52 before `editors := detectEditors(home)`), add:

```go
	// Auto-enable MCP when registering with editors
	cfg := config.Load()
	if !cfg.MCPEnabled {
		if err := config.Set("mcp.enabled", "true"); err != nil {
			fmt.Fprintf(os.Stderr, "  [!] Warning: could not enable MCP in config: %v\n", err)
		} else {
			fmt.Println("  [+] MCP enabled in config")
		}
	}
```

Also add the config import. Add to the imports:
```go
	"github.com/nerveband/ai-happy-design/internal/config"
```

**Step 2: Verify it compiles**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go build ./cmd/ai-happy-design/
```

Expected: No errors

**Step 3: Commit**

```bash
git add cmd/ai-happy-design/register.go
git commit -m "feat: register command auto-enables MCP in config"
```

---

### Task 6: Build, sign, and verify end-to-end

**Files:** None (verification only)

**Step 1: Run all Go tests**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./... -count=1
```

Expected: All tests PASS

**Step 2: Build and sign**

Run:
```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go build -o ai-happy-design ./cmd/ai-happy-design/ && cp ai-happy-design /tmp/ai-happy-design && codesign -f -s - /tmp/ai-happy-design
```

Expected: Build succeeds, signing succeeds

**Step 3: Verify MCP is gated**

Run:
```bash
/tmp/ai-happy-design mcp 2>&1; echo "exit: $?"
```

Expected: Error message about MCP being disabled, exit code 1

**Step 4: Verify config commands work**

Run:
```bash
/tmp/ai-happy-design config path
/tmp/ai-happy-design config init --force
/tmp/ai-happy-design config get mcp.enabled
/tmp/ai-happy-design config set mcp.enabled true
/tmp/ai-happy-design config get mcp.enabled
```

Expected: Path printed, config created, `false`, set confirmation, `true`

**Step 5: Reset for clean state**

Run:
```bash
/tmp/ai-happy-design config set mcp.enabled false
```
