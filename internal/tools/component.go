package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterComponentTool registers the "component" tool for component and variant operations.
func RegisterComponentTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("component",
		mcp.WithDescription("Component operations: create, instantiate, create variant sets, get local/remote components, get/set overrides."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create", "create_instance", "create_set", "get_local", "get_remote", "get_overrides", "set_overrides")),
		mcp.WithString("nodeId", mcp.Description("Target node ID")),
		mcp.WithString("componentKey", mcp.Description("Component key for creating instances")),
		mcp.WithString("name", mcp.Description("Name for the component")),
		mcp.WithString("nodeIds", mcp.Description("Comma-separated node IDs (for create_set)")),
		mcp.WithString("overrides", mcp.Description("JSON object of property overrides")),
		mcp.WithNumber("x", mcp.Description("X position")),
		mcp.WithNumber("y", mcp.Description("Y position")),
		mcp.WithString("parentId", mcp.Description("Parent node ID")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_component_from_node", map[string]interface{}{
				"nodeId": nodeId,
				"name":   getStringArg(args, "name", ""),
			})

		case "create_instance":
			componentKey, errResult := requireStringArg(args, "componentKey")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_component_instance", map[string]interface{}{
				"componentKey": componentKey,
				"x":            getFloat64Arg(args, "x", 0),
				"y":            getFloat64Arg(args, "y", 0),
				"parentId":     getStringArg(args, "parentId", ""),
			})

		case "create_set":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_component_set", map[string]interface{}{
				"nodeIds": nodeIds,
				"name":    getStringArg(args, "name", ""),
			})

		case "get_local":
			return sendCommand(commander, "get_local_components", map[string]interface{}{})

		case "get_remote":
			return sendCommand(commander, "get_remote_components", map[string]interface{}{})

		case "get_overrides":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "get_overrides", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "set_overrides":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_overrides", map[string]interface{}{
				"nodeId":    nodeId,
				"overrides": getStringArg(args, "overrides", "{}"),
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown component action: %s", action)), nil
		}
	})
}
