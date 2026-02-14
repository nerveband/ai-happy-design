package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterExportTool registers the "export" tool for exporting nodes as images.
func RegisterExportTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("export",
		mcp.WithDescription("Export operations: export nodes as PNG/JPG/SVG/PDF images."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("image", "svg", "pdf")),
		mcp.WithString("nodeId", mcp.Required(), mcp.Description("Node ID to export")),
		mcp.WithNumber("scale", mcp.Description("Export scale (1x, 2x, etc)")),
		mcp.WithString("format", mcp.Description("Image format: PNG, JPG, SVG, PDF"),
			mcp.Enum("PNG", "JPG", "SVG", "PDF")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")
		nodeId := getStringArg(args, "nodeId", "")

		switch action {
		case "image":
			return sendCommand(commander, "export_node_as_image", map[string]interface{}{
				"nodeId": nodeId,
				"format": getStringArg(args, "format", "PNG"),
				"scale":  getFloat64Arg(args, "scale", 1),
			})

		case "svg":
			return sendCommand(commander, "export_node_as_image", map[string]interface{}{
				"nodeId": nodeId,
				"format": "SVG",
				"scale":  getFloat64Arg(args, "scale", 1),
			})

		case "pdf":
			return sendCommand(commander, "export_node_as_image", map[string]interface{}{
				"nodeId": nodeId,
				"format": "PDF",
				"scale":  getFloat64Arg(args, "scale", 1),
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown export action: %s", action)), nil
		}
	})
}
