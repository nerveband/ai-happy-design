package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nerveband/ai-happy-design/internal/config"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{Use: "profile", Short: "Manage local CLI profiles"}
var profileListCmd = &cobra.Command{Use: "list", Short: "List profiles", RunE: func(cmd *cobra.Command, args []string) error {
	profiles, active, err := listProfiles()
	if err != nil {
		return err
	}
	return printJSON(map[string]interface{}{"profiles": profiles, "active": active})
}}
var profileUseCmd = &cobra.Command{Use: "use <name>", Short: "Set active profile", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := os.MkdirAll(config.Dir(), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(config.Dir(), "active-profile"), []byte(args[0]), 0644); err != nil {
		return err
	}
	return printJSON(map[string]interface{}{"active": args[0]})
}}
var profileInspectCmd = &cobra.Command{Use: "inspect [name]", Short: "Inspect profile", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		_, active, _ := listProfiles()
		name = active
	}
	if name == "" {
		return fmt.Errorf("no profile selected")
	}
	path := filepath.Join(config.Dir(), "profiles", name+".json")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return printJSON(map[string]interface{}{"name": name, "path": path, "exists": false})
	}
	if err != nil {
		return err
	}
	var payload map[string]interface{}
	_ = json.Unmarshal(raw, &payload)
	redactSecrets(payload)
	return printJSON(map[string]interface{}{"name": name, "path": path, "exists": true, "profile": payload})
}}

func listProfiles() ([]string, string, error) {
	dir := filepath.Join(config.Dir(), "profiles")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []string{}, activeProfile(), nil
	}
	if err != nil {
		return nil, "", err
	}
	out := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			out = append(out, entry.Name()[:len(entry.Name())-5])
		}
	}
	return out, activeProfile(), nil
}

func activeProfile() string {
	raw, err := os.ReadFile(filepath.Join(config.Dir(), "active-profile"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(raw))
}

func redactSecrets(m map[string]interface{}) {
	for key, value := range m {
		switch key {
		case "token", "apiKey", "password", "secret":
			m[key] = "[redacted]"
		default:
			if child, ok := value.(map[string]interface{}); ok {
				redactSecrets(child)
			}
		}
	}
}

func init() {
	profileCmd.AddCommand(profileListCmd, profileUseCmd, profileInspectCmd)
	rootCmd.AddCommand(profileCmd)
}
