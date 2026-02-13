package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design-v2/internal/config"
	"github.com/nerveband/ai-happy-design-v2/internal/mcp"
	"github.com/nerveband/ai-happy-design-v2/internal/ws"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "ai-happy-design",
	Short: "AI Happy Design - Figma MCP server and CLI",
	Long: `A single binary that works as both an MCP server for AI editors
and a CLI for direct Figma manipulation via a WebSocket relay.`,
	Version: version,
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server (stdio transport)",
	Long: `Starts the MCP server on stdio for use with AI editors (Claude Code,
Cursor, Windsurf, etc). Also starts a WebSocket relay server in the
background for communicating with the Figma plugin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		wsServer := ws.NewServer(cfg.Port)
		go wsServer.Start()
		return mcp.StartServer(wsServer)
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect [channel-key]",
	Short: "Connect to a Figma plugin channel",
	Long: `Connect to the WebSocket relay and join a specific channel.
The channel key should match what the Figma plugin is using
(e.g. happy-unicorn-42).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		url := fmt.Sprintf("ws://%s:%d", cfg.ServerHost, cfg.Port)
		client := ws.NewClient(url)
		return client.JoinChannel(args[0])
	},
}

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "Start WebSocket relay server only",
	Long:  `Starts just the WebSocket relay server without the MCP server. Useful for debugging or running the relay separately.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		wsServer := ws.NewServer(cfg.Port)
		fmt.Printf("Starting WebSocket relay on port %d...\n", cfg.Port)
		return wsServer.Start()
	},
}

func main() {
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(wsCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
