package tools

import (
	"context"
	"encoding/json"
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
			format := getStringArg(args, "format", "PNG")
			return sendExportCommand(commander, nodeId, format, getFloat64Arg(args, "scale", 1))

		case "svg":
			return sendExportCommand(commander, nodeId, "SVG", getFloat64Arg(args, "scale", 1))

		case "pdf":
			return sendExportCommand(commander, nodeId, "PDF", getFloat64Arg(args, "scale", 1))

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown export action: %s", action)), nil
		}
	})
}

// sendExportCommand sends an export command and returns an MCP image content block
// so the LLM can visually inspect the result.
func sendExportCommand(commander *figma.Commander, nodeId, format string, scale float64) (*mcp.CallToolResult, error) {
	result, err := commander.SendCommand("export_node_as_image", map[string]interface{}{
		"nodeId": nodeId,
		"format": format,
		"scale":  scale,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Extract fields from the plugin response: {id, name, format, scale, size, data}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		// Fallback to text-only response
		text, _ := formatResult(result)
		return mcp.NewToolResultText(text), nil
	}

	base64Data, _ := resultMap["data"].(string)
	nodeName, _ := resultMap["name"].(string)
	nodeID, _ := resultMap["id"].(string)
	sizeBytes, _ := resultMap["size"].(float64)

	if base64Data == "" {
		return mcp.NewToolResultError("export returned empty data"), nil
	}

	// For SVG, return as text (it's not an image content block)
	if format == "SVG" {
		text, _ := formatResult(result)
		return mcp.NewToolResultText(text), nil
	}

	// Build metadata text
	meta := map[string]interface{}{
		"id":     nodeID,
		"name":   nodeName,
		"format": format,
		"scale":  scale,
		"size":   sizeBytes,
	}
	metaJSON, _ := json.MarshalIndent(meta, "", "  ")

	// Determine MIME type
	mimeType := "image/png"
	switch format {
	case "JPG":
		mimeType = "image/jpeg"
	case "PDF":
		mimeType = "application/pdf"
	}

	// Return both text metadata and the image content block
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: mcp.ContentTypeText,
				Text: string(metaJSON),
			},
			mcp.ImageContent{
				Type:     mcp.ContentTypeImage,
				Data:     base64Data,
				MIMEType: mimeType,
			},
		},
	}, nil
}
