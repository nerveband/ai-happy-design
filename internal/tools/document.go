package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterDocumentTool registers the "document" tool for document-level operations.
func RegisterDocumentTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("document",
		mcp.WithDescription("Document operations: get info, get/set selection, scan text nodes, scan by type, get styles, focus viewport."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("get_info", "get_selection", "set_selection", "scan_text", "scan_by_type", "get_styles", "focus")),
		mcp.WithString("nodeId", mcp.Description("Node ID to focus on")),
		mcp.WithString("nodeIds", mcp.Description("Comma-separated node IDs (for set_selection)")),
		mcp.WithString("nodeType", mcp.Description("Node type to scan for (e.g. TEXT, FRAME, RECTANGLE)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "get_info":
			return sendCommand(commander, "get_document_info", map[string]interface{}{})

		case "get_selection":
			return sendCommand(commander, "get_selection", map[string]interface{}{})

		case "set_selection":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_selection", map[string]interface{}{
				"nodeIds": nodeIds,
			})

		case "scan_text":
			return sendCommand(commander, "scan_text_nodes", map[string]interface{}{})

		case "scan_by_type":
			nodeType, errResult := requireStringArg(args, "nodeType")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "scan_by_type", map[string]interface{}{
				"nodeType": nodeType,
			})

		case "get_styles":
			return sendCommand(commander, "get_styles", map[string]interface{}{})

		case "focus":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "focus_node", map[string]interface{}{
				"nodeId": nodeId,
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown document action: %s", action)), nil
		}
	})
}
