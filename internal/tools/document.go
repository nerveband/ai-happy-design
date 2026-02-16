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
		mcp.WithDescription("Document operations: get info, get/set selection, scan text nodes, scan by type, get styles, focus viewport, find free canvas space, find_nodes (unified search by name/type/text content)."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("get_info", "get_selection", "set_selection", "scan_text", "scan_by_type", "get_styles", "focus", "find_free_space", "find_nodes")),
		mcp.WithString("nodeId", mcp.Description("Node ID to focus on")),
		mcp.WithString("nodeIds", mcp.Description("Comma-separated node IDs (for set_selection)")),
		mcp.WithString("nodeType", mcp.Description("Node type to scan for (e.g. TEXT, FRAME, RECTANGLE)")),
		mcp.WithString("query", mcp.Description("Search query for find_nodes (matches node names)")),
		mcp.WithString("textContent", mcp.Description("Search text content for find_nodes (matches TEXT node characters)")),
		mcp.WithNumber("width", mcp.Description("Desired frame width in pixels (for find_free_space, default 1080)")),
		mcp.WithNumber("height", mcp.Description("Desired frame height in pixels (for find_free_space, default 1080)")),
		mcp.WithNumber("gap", mcp.Description("Gap between frames in pixels (for find_free_space, default 100)")),
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

		case "find_free_space":
			params := map[string]interface{}{}
			if w := getFloat64Arg(args, "width", 0); w > 0 {
				params["width"] = w
			}
			if h := getFloat64Arg(args, "height", 0); h > 0 {
				params["height"] = h
			}
			if g := getFloat64Arg(args, "gap", 0); g > 0 {
				params["gap"] = g
			}
			return sendCommand(commander, "find_free_space", params)

		case "find_nodes":
			params := map[string]interface{}{}
			if hasArg(args, "query") {
				params["query"] = getStringArg(args, "query", "")
			}
			if hasArg(args, "nodeType") {
				params["type"] = getStringArg(args, "nodeType", "")
			}
			if hasArg(args, "textContent") {
				params["textContent"] = getStringArg(args, "textContent", "")
			}
			return sendCommand(commander, "document.find_nodes", params)

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown document action: %s", action)), nil
		}
	})
}
