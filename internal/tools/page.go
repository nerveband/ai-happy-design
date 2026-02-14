package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterPageTool registers the "page" tool for page management.
func RegisterPageTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("page",
		mcp.WithDescription("Page operations: create, delete, rename, duplicate, set current, get all pages."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create", "delete", "rename", "duplicate", "set_current", "get_all")),
		mcp.WithString("pageId", mcp.Description("Target page ID")),
		mcp.WithString("name", mcp.Description("Page name")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create":
			return sendCommand(commander, "create_page", map[string]interface{}{
				"name": getStringArg(args, "name", "New Page"),
			})

		case "delete":
			pageId, errResult := requireStringArg(args, "pageId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "delete_page", map[string]interface{}{
				"pageId": pageId,
			})

		case "rename":
			pageId, errResult := requireStringArg(args, "pageId")
			if errResult != nil {
				return errResult, nil
			}
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "rename_page", map[string]interface{}{
				"pageId": pageId,
				"name":   name,
			})

		case "duplicate":
			pageId, errResult := requireStringArg(args, "pageId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "duplicate_page", map[string]interface{}{
				"pageId": pageId,
			})

		case "set_current":
			pageId, errResult := requireStringArg(args, "pageId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_current_page", map[string]interface{}{
				"pageId": pageId,
			})

		case "get_all":
			return sendCommand(commander, "get_pages", map[string]interface{}{})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown page action: %s", action)), nil
		}
	})
}
