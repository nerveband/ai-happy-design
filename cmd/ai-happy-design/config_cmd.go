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

var configSourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Show config precedence and source paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printJSON(map[string]interface{}{
			"precedence": []string{"flags", "environment", "profile", "config file", "defaults"},
			"configDir":  config.Dir(),
			"configPath": config.Path(),
			"env":        []string{"PORT", "SERVER_HOST", "AHD_IDLE_TIMEOUT", "AHD_MCP_ENABLED", "AHD_CHANNEL", "AHD_CONFIG_DIR"},
		})
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
