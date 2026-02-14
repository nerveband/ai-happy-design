package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterVariableTool registers the "variable" tool for Figma variables and collections.
func RegisterVariableTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("variable",
		mcp.WithDescription("Variable operations: create variables, get all, set values, bind/unbind to nodes, create collections."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create", "get_all", "set_value", "bind", "unbind", "create_collection")),
		mcp.WithString("nodeId", mcp.Description("Target node ID (for bind/unbind)")),
		mcp.WithString("name", mcp.Description("Variable or collection name")),
		mcp.WithString("variableId", mcp.Description("Variable ID")),
		mcp.WithString("collectionId", mcp.Description("Collection ID")),
		mcp.WithString("resolvedType", mcp.Description("Variable type: COLOR, FLOAT, STRING, BOOLEAN"),
			mcp.Enum("COLOR", "FLOAT", "STRING", "BOOLEAN")),
		mcp.WithString("value", mcp.Description("Value to set (JSON encoded for complex types)")),
		mcp.WithString("field", mcp.Description("Node field to bind to (e.g. fills, opacity)")),
		mcp.WithString("modeId", mcp.Description("Mode ID for setting values")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create":
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			collectionId, errResult := requireStringArg(args, "collectionId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_variable", map[string]interface{}{
				"name":         name,
				"collectionId": collectionId,
				"resolvedType": getStringArg(args, "resolvedType", "COLOR"),
			})

		case "get_all":
			return sendCommand(commander, "get_variables", map[string]interface{}{})

		case "set_value":
			variableId, errResult := requireStringArg(args, "variableId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_variable_value", map[string]interface{}{
				"variableId": variableId,
				"value":      getStringArg(args, "value", ""),
				"modeId":     getStringArg(args, "modeId", ""),
			})

		case "bind":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			variableId, errResult := requireStringArg(args, "variableId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "bind_variable", map[string]interface{}{
				"nodeId":     nodeId,
				"variableId": variableId,
				"field":      getStringArg(args, "field", "fills"),
			})

		case "unbind":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "unbind_variable", map[string]interface{}{
				"nodeId": nodeId,
				"field":  getStringArg(args, "field", "fills"),
			})

		case "create_collection":
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_variable_collection", map[string]interface{}{
				"name": name,
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown variable action: %s", action)), nil
		}
	})
}
