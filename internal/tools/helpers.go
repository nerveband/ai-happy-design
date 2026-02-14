package tools

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// sendCommand is a convenience that sends a Figma command and returns an MCP
// text result or an MCP error result.
func sendCommand(commander *figma.Commander, command string, params map[string]interface{}) (*mcp.CallToolResult, error) {
	result, err := commander.SendCommand(command, params)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	text, err := formatResult(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format result: %v", err)), nil
	}
	return mcp.NewToolResultText(text), nil
}

// formatResult converts an arbitrary result to a human-readable JSON string.
func formatResult(result interface{}) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// getStringArg extracts a string argument with a default.
func getStringArg(args map[string]interface{}, key, defaultVal string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

// getFloat64Arg extracts a float64 argument with a default.
func getFloat64Arg(args map[string]interface{}, key string, defaultVal float64) float64 {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		}
	}
	return defaultVal
}

// getBoolArg extracts a boolean argument with a default.
func getBoolArg(args map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}

// requireStringArg extracts a required string argument or returns an error result.
func requireStringArg(args map[string]interface{}, key string) (string, *mcp.CallToolResult) {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s, nil
		}
	}
	return "", mcp.NewToolResultError(fmt.Sprintf("required argument %q is missing or empty", key))
}
