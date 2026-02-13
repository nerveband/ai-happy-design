package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design-v2/internal/figma"
)

// BulkOperation represents a single operation in a bulk execute request.
type BulkOperation struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

// RegisterBulkTool registers the "bulk" tool for executing multiple operations atomically.
func RegisterBulkTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("bulk",
		mcp.WithDescription("Execute multiple Figma operations in sequence. Provide a JSON array of operations."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("execute")),
		mcp.WithString("operations", mcp.Required(), mcp.Description("JSON array of operations: [{\"command\": \"...\", \"params\": {...}}, ...]")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "execute":
			opsStr, errResult := requireStringArg(args, "operations")
			if errResult != nil {
				return errResult, nil
			}

			var ops []BulkOperation
			if err := json.Unmarshal([]byte(opsStr), &ops); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("invalid operations JSON: %v", err)), nil
			}

			if len(ops) == 0 {
				return mcp.NewToolResultError("operations array is empty"), nil
			}

			results := make([]interface{}, 0, len(ops))
			for i, op := range ops {
				result, err := commander.SendCommand(op.Command, op.Params)
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("operation %d (%s) failed: %v", i, op.Command, err)), nil
				}
				results = append(results, map[string]interface{}{
					"index":   i,
					"command": op.Command,
					"result":  result,
				})
			}

			text, err := formatResult(results)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to format results: %v", err)), nil
			}
			return mcp.NewToolResultText(text), nil

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown bulk action: %s", action)), nil
		}
	})
}
