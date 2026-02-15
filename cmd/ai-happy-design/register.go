package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type editorConfig struct {
	Name       string
	ConfigPath string // relative to home dir
	Detected   bool
}

var registerEditor string
var registerForce bool

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register MCP server with AI editors",
	Long: `Auto-detects installed AI editors and registers ai-happy-design as an MCP server.
Supports: Claude Code, Claude Desktop, Cursor, Windsurf, VS Code (Copilot), Zed.
Use --editor to register with a specific editor only.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRegister()
	},
}

func runRegister() error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find binary path: %w", err)
	}
	// Resolve symlinks
	binaryPath, err = filepath.EvalSymlinks(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to resolve binary path: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to find home directory: %w", err)
	}

	editors := detectEditors(home)

	if registerEditor != "" {
		// Filter to requested editor
		found := false
		for i := range editors {
			if strings.EqualFold(editors[i].Name, registerEditor) {
				editors = []editorConfig{editors[i]}
				editors[0].Detected = true // force register
				found = true
				break
			}
		}
		if !found {
			names := make([]string, 0)
			for _, e := range editors {
				names = append(names, e.Name)
			}
			return fmt.Errorf("unknown editor %q. Supported: %s", registerEditor, strings.Join(names, ", "))
		}
	}

	registered := 0
	skipped := 0

	for _, editor := range editors {
		if !editor.Detected && !registerForce {
			continue
		}

		var regErr error
		switch editor.Name {
		case "Claude Code":
			regErr = registerClaudeCode(binaryPath)
		case "Claude Desktop":
			regErr = registerJSONConfig(home, editor, binaryPath)
		case "Cursor":
			regErr = registerJSONConfig(home, editor, binaryPath)
		case "Windsurf":
			regErr = registerJSONConfig(home, editor, binaryPath)
		case "VS Code":
			regErr = registerJSONConfig(home, editor, binaryPath)
		case "Zed":
			regErr = registerZed(home, editor, binaryPath)
		default:
			continue
		}

		if regErr != nil {
			fmt.Fprintf(os.Stderr, "  [!] %s: %v\n", editor.Name, regErr)
			skipped++
		} else {
			fmt.Printf("  [+] %s: registered\n", editor.Name)
			registered++
		}
	}

	if registered == 0 && skipped == 0 {
		fmt.Println("No supported AI editors detected.")
		fmt.Println("Use --editor <name> to register manually, or --force to register all.")
		fmt.Println("Supported: Claude Code, Claude Desktop, Cursor, Windsurf, VS Code, Zed")
		return nil
	}

	fmt.Printf("\nRegistered: %d, Skipped: %d\n", registered, skipped)
	return nil
}

func detectEditors(home string) []editorConfig {
	editors := []editorConfig{
		{Name: "Claude Code", ConfigPath: "", Detected: isCommandAvailable("claude")},
		{Name: "Claude Desktop", ConfigPath: claudeDesktopConfigPath(), Detected: false},
		{Name: "Cursor", ConfigPath: filepath.Join(".cursor", "mcp.json"), Detected: false},
		{Name: "Windsurf", ConfigPath: filepath.Join(".codeium", "windsurf", "mcp_config.json"), Detected: false},
		{Name: "VS Code", ConfigPath: filepath.Join(".vscode", "mcp.json"), Detected: false},
		{Name: "Zed", ConfigPath: filepath.Join(".config", "zed", "settings.json"), Detected: false},
	}

	for i := range editors {
		if editors[i].Name == "Claude Code" {
			continue // already checked via command
		}
		if editors[i].ConfigPath != "" {
			fullPath := filepath.Join(home, editors[i].ConfigPath)
			// Check if the config file OR its parent directory exists
			if _, err := os.Stat(fullPath); err == nil {
				editors[i].Detected = true
			} else if _, err := os.Stat(filepath.Dir(fullPath)); err == nil {
				editors[i].Detected = true
			}
		}
	}

	return editors
}

func claudeDesktopConfigPath() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join("Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "windows":
		return filepath.Join("AppData", "Roaming", "Claude", "claude_desktop_config.json")
	default:
		return filepath.Join(".config", "claude", "claude_desktop_config.json")
	}
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// registerClaudeCode uses `claude mcp add-json` for native registration.
func registerClaudeCode(binaryPath string) error {
	mcpJSON := map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"mcp"},
		"type":    "stdio",
	}
	data, err := json.Marshal(mcpJSON)
	if err != nil {
		return err
	}

	cmd := exec.Command("claude", "mcp", "add-json", "ai-happy-design", "--scope", "user")
	cmd.Stdin = strings.NewReader(string(data))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// registerJSONConfig handles editors that use a JSON config file with mcpServers key.
func registerJSONConfig(home string, editor editorConfig, binaryPath string) error {
	configPath := filepath.Join(home, editor.ConfigPath)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load existing config or create new one
	var config map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse existing config %s: %w", configPath, err)
		}
	} else {
		config = make(map[string]interface{})
	}

	// Get or create mcpServers section
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	// Check if already registered
	if _, exists := mcpServers["ai-happy-design"]; exists && !registerForce {
		return fmt.Errorf("already registered (use --force to overwrite)")
	}

	// Add our MCP server entry
	mcpServers["ai-happy-design"] = map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"mcp"},
	}
	config["mcpServers"] = mcpServers

	// Write back
	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write config %s: %w", configPath, err)
	}

	return nil
}

// stripJSONC removes single-line (//) and multi-line (/* */) comments from JSONC.
func stripJSONC(data []byte) []byte {
	// Remove multi-line comments
	re := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	data = re.ReplaceAll(data, nil)
	// Remove single-line comments (but not inside strings)
	re2 := regexp.MustCompile(`(?m)^\s*//.*$`)
	data = re2.ReplaceAll(data, nil)
	// Remove trailing commas before } or ]
	re3 := regexp.MustCompile(`,\s*([\]}])`)
	data = re3.ReplaceAll(data, []byte("$1"))
	return data
}

// registerZed handles Zed's settings.json which uses JSONC (JSON with comments).
func registerZed(home string, editor editorConfig, binaryPath string) error {
	configPath := filepath.Join(home, editor.ConfigPath)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var config map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err == nil {
		// Zed uses JSONC — strip comments before parsing
		cleaned := stripJSONC(data)
		if err := json.Unmarshal(cleaned, &config); err != nil {
			return fmt.Errorf("failed to parse existing config %s: %w", configPath, err)
		}
	} else {
		config = make(map[string]interface{})
	}

	// Zed uses "context_servers" key
	ctxServers, ok := config["context_servers"].(map[string]interface{})
	if !ok {
		ctxServers = make(map[string]interface{})
	}

	if _, exists := ctxServers["ai-happy-design"]; exists && !registerForce {
		return fmt.Errorf("already registered (use --force to overwrite)")
	}

	ctxServers["ai-happy-design"] = map[string]interface{}{
		"command":  binaryPath,
		"args":     []string{"mcp"},
		"settings": map[string]interface{}{},
	}
	config["context_servers"] = ctxServers

	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write config %s: %w", configPath, err)
	}

	return nil
}
