package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

// RegisterConnectTool registers the "connect" tool for managing plugin connections.
func RegisterConnectTool(s *server.MCPServer, commander *figma.Commander, wsServer *ws.Server) {
	tool := mcp.NewTool("connect",
		mcp.WithDescription("Connection management: join a channel, check status, disconnect."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("join", "status", "disconnect")),
		mcp.WithString("channelKey", mcp.Description("Channel key to join (e.g. happy-unicorn-42)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "join":
			channelKey := getStringArg(args, "channelKey", "")
			if channelKey == "" {
				channelKey = ws.GenerateChannelKey()
			}
			status := "waiting"
			if wsServer.HasChannel(channelKey) {
				status = "connected"
			}
			return mcp.NewToolResultText(fmt.Sprintf(
				"Channel: %s\nStatus: %s\nUse this channel key in the Figma plugin to connect.",
				channelKey, status,
			)), nil

		case "status":
			activeChannel := wsServer.GetChannelKey()
			if activeChannel == "" {
				return mcp.NewToolResultText("Status: disconnected\nNo Figma plugin is currently connected.\nStart the Figma plugin and use the channel key to connect."), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf(
				"Status: connected\nActive channel: %s\nConnected clients: %d",
				activeChannel, wsServer.ChannelCount(),
			)), nil

		case "disconnect":
			return mcp.NewToolResultText("Disconnect is handled by the Figma plugin. Close the plugin or leave the channel there."), nil

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown connect action: %s", action)), nil
		}
	})
}
