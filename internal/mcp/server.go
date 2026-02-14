package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

// StartServer creates the MCP server, registers all tools, and starts the
// stdio transport. It blocks until the process is terminated.
func StartServer(wsServer *ws.Server) error {
	s := server.NewMCPServer(
		"AI Happy Design",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all domain-grouped tools
	tools.RegisterAll(s, wsServer)

	// Start stdio transport (blocks until stdin closes)
	return server.ServeStdio(s)
}
