package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// toolDescriptions provides self-documentation for all tools.
var toolDescriptions = map[string]map[string]string{
	"paint": {
		"set_solid":               "Set a solid fill color on a node. Params: nodeId, color (hex), opacity",
		"set_gradient":            "Set a gradient fill. Params: nodeId, gradientType, stops (JSON)",
		"set_image":               "Set an image fill from base64 data. Params: nodeId, imageData (base64 or data URL), scaleMode",
		"set_image_fill":          "Alias of set_image for MCP/CLI parity.",
		"set_image_url":           "Set an image fill from URL. Params: nodeId, url, scaleMode. Plugin tries createImageAsync first, then fetch fallback.",
		"set_image_fill_from_url": "Alias of set_image_url for MCP/CLI parity.",
		"add_fill":                "Add a fill to existing fills. Params: nodeId, color, opacity",
		"remove_fill":             "Remove a fill by index. Params: nodeId, fillIndex",
		"get_fills":               "Get all fills on a node. Params: nodeId",
		"set_stroke":              "Set stroke color and weight. Params: nodeId, color, opacity, strokeWeight, strokeAlign",
	},
	"shape": {
		"create_rectangle": "Create a rectangle. Params: name, x, y, width, height, parentId, color, cornerRadius",
		"create_ellipse":   "Create an ellipse. Params: name, x, y, width, height, parentId, color",
		"create_polygon":   "Create a polygon. Params: name, x, y, width, height, parentId, color, pointCount",
		"create_star":      "Create a star. Params: name, x, y, width, height, parentId, color, pointCount, innerRadius",
		"create_line":      "Create a line. Params: name, x, y, width, parentId, color",
		"create_from_svg":  "Create shape from SVG path. Params: name, x, y, svgPath, parentId",
	},
	"text": {
		"create":         "Create a text node. Params: text, x, y, fontFamily, fontStyle, fontSize, parentId, name",
		"set_content":    "Set text content. Params: nodeId, text",
		"set_font":       "Set font family and style. Params: nodeId, fontFamily, fontStyle",
		"set_size":       "Set font size. Params: nodeId, fontSize",
		"set_weight":     "Set font weight. Params: nodeId, fontWeight",
		"set_color":      "Set text color. Params: nodeId, color",
		"set_align":      "Set text alignment. Params: nodeId, textAlign",
		"set_spacing":    "Set letter/line/paragraph spacing. Params: nodeId, letterSpacing, lineHeight, paragraphSpacing",
		"set_case":       "Set text case. Params: nodeId, textCase",
		"set_decoration": "Set text decoration. Params: nodeId, textDecoration",
		"get_segments":   "Get styled text segments. Params: nodeId",
		"load_font":      "Load a font for use. Params: fontFamily, fontStyle",
		"set_style_id":   "Apply a text style. Params: nodeId, styleId",
	},
	"layout": {
		"set_auto_layout": "Enable auto-layout on a frame. Params: nodeId, direction, itemSpacing, padding*, primaryAxisAlign, counterAxisAlign",
		"set_padding":     "Set padding. Params: nodeId, padding (uniform) or paddingTop/Right/Bottom/Left",
		"set_spacing":     "Set item spacing. Params: nodeId, itemSpacing",
		"set_alignment":   "Set axis alignment. Params: nodeId, primaryAxisAlign, counterAxisAlign",
		"set_sizing":      "Set axis sizing mode. Params: nodeId, primaryAxisSizing, counterAxisSizing",
		"set_wrap":        "Set layout wrap. Params: nodeId, layoutWrap",
		"set_constraints": "Set layout constraints. Params: nodeId, constraintHorizontal, constraintVertical",
	},
	"node": {
		"get_info":       "Get node information. Params: nodeId",
		"get_tree":       "Get node tree. Params: nodeId, depth",
		"create_frame":   "Create a frame. Params: name, x, y, width, height, parentId, color",
		"move":           "Move a node. Params: nodeId, x, y",
		"resize":         "Resize a node. Params: nodeId, width, height",
		"rotate":         "Rotate a node. Params: nodeId, rotation (degrees)",
		"set_opacity":    "Set opacity. Params: nodeId, opacity",
		"set_blend_mode": "Set blend mode. Params: nodeId, blendMode",
		"set_visibility": "Set visibility. Params: nodeId, visible",
		"set_locked":     "Set lock state. Params: nodeId, locked",
		"rename":         "Rename a node. Params: nodeId, name",
		"delete":         "Delete a node. Params: nodeId",
		"clone":          "Clone a node. Params: nodeId",
	},
	"layer": {
		"set_order":      "Set layer order. Params: nodeId, index",
		"bring_forward":  "Bring one layer forward. Params: nodeId",
		"send_backward":  "Send one layer backward. Params: nodeId",
		"bring_to_front": "Bring to front. Params: nodeId",
		"send_to_back":   "Send to back. Params: nodeId",
		"group":          "Group nodes. Params: nodeIds (comma-separated)",
		"ungroup":        "Ungroup a group. Params: nodeId",
		"insert_child":   "Insert a child into a parent. Params: parentId, childId, index",
	},
	"component": {
		"create":          "Create component from node. Params: nodeId, name",
		"create_instance": "Create component instance. Params: componentKey, x, y, parentId",
		"create_set":      "Create component set (variants). Params: nodeIds, name",
		"get_local":       "Get local components. No params",
		"get_remote":      "Get remote components. No params",
		"get_overrides":   "Get instance overrides. Params: nodeId",
		"set_overrides":   "Set instance overrides. Params: nodeId, overrides (JSON)",
	},
	"style": {
		"create_paint":  "Create paint style. Params: name, color, description",
		"create_text":   "Create text style. Params: name, description",
		"create_effect": "Create effect style. Params: name, description",
		"apply":         "Apply style to node. Params: nodeId, styleId, styleType",
		"get_all":       "Get all styles. No params",
		"remove":        "Remove style from node. Params: nodeId, styleType",
	},
	"variable": {
		"create":            "Create a variable. Params: name, collectionId, resolvedType",
		"get_all":           "Get all variables. No params",
		"set_value":         "Set variable value. Params: variableId, value, modeId",
		"bind":              "Bind variable to node. Params: nodeId, variableId, field",
		"unbind":            "Unbind variable from node. Params: nodeId, field",
		"create_collection": "Create variable collection. Params: name",
	},
	"effect": {
		"set_effects": "Set all effects on a node. Params: nodeId, effects (JSON array)",
		"add_shadow":  "Add shadow effect. Params: nodeId, shadowType, color, offsetX, offsetY, radius, spread",
		"add_blur":    "Add blur effect. Params: nodeId, blurType, radius",
		"apply_style": "Apply effect style. Params: nodeId, styleId",
		"remove":      "Remove all effects. Params: nodeId",
	},
	"boolean": {
		"union":     "Boolean union. Params: nodeIds (comma-separated)",
		"subtract":  "Boolean subtract. Params: nodeIds (comma-separated)",
		"intersect": "Boolean intersect. Params: nodeIds (comma-separated)",
		"exclude":   "Boolean exclude. Params: nodeIds (comma-separated)",
		"flatten":   "Flatten a node. Params: nodeId",
	},
	"page": {
		"create":      "Create a page. Params: name",
		"delete":      "Delete a page. Params: pageId",
		"rename":      "Rename a page. Params: pageId, name",
		"duplicate":   "Duplicate a page. Params: pageId",
		"set_current": "Set current page. Params: pageId",
		"get_all":     "Get all pages. No params",
	},
	"document": {
		"get_info":      "Get document info. No params",
		"get_selection": "Get current selection. No params",
		"set_selection": "Set selection. Params: nodeIds (comma-separated)",
		"scan_text":     "Scan all text nodes. No params",
		"scan_by_type":  "Scan nodes by type. Params: nodeType",
		"get_styles":    "Get all document styles. No params",
		"focus":         "Focus viewport on node. Params: nodeId",
	},
	"export": {
		"image": "Export as image (PNG/JPG). Params: nodeId, format, scale",
		"svg":   "Export as SVG. Params: nodeId, scale",
		"pdf":   "Export as PDF. Params: nodeId, scale",
	},
	"bulk": {
		"execute": "Execute multiple operations with retry and optional interpolation. Params: operations (JSON array of {name?, command, params}), continueOnError, retries, retryDelayMs, interpolate",
	},
	"connect": {
		"join":       "Join a channel. Params: channelKey",
		"status":     "Get connection status. No params",
		"disconnect": "Disconnect from channel. No params",
	},
	"design": {
		"compute_tokens": "Compute design tokens (font sizes, spacing, padding, layout, card widths) for any canvas dimensions. Call FIRST before creating any design. Returns concrete pixel values — no math needed. Params: *width, *height",
	},
}

// ToolCatalog returns a copy of the server tool/action descriptions for
// machine-readable discovery in CLI or external integrations.
func ToolCatalog() map[string]map[string]string {
	out := make(map[string]map[string]string, len(toolDescriptions))
	for tool, actions := range toolDescriptions {
		actionCopy := make(map[string]string, len(actions))
		for action, desc := range actions {
			actionCopy[action] = desc
		}
		out[tool] = actionCopy
	}
	return out
}

// RegisterDescribeTool registers the "describe" tool for self-documentation.
func RegisterDescribeTool(s *server.MCPServer, _ *figma.Commander) {
	tool := mcp.NewTool("describe",
		mcp.WithDescription("Get help and documentation for tools and actions. Use this to discover available operations."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("tool", "action", "all", "catalog", "design_guide", "setup")),
		mcp.WithString("toolName", mcp.Description("Name of the tool to describe")),
		mcp.WithString("actionName", mcp.Description("Name of the action within a tool")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "tool":
			toolName, errResult := requireStringArg(args, "toolName")
			if errResult != nil {
				return errResult, nil
			}
			actions, ok := toolDescriptions[toolName]
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("unknown tool: %s", toolName)), nil
			}
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Tool: %s\nActions:\n", toolName))
			for name, desc := range actions {
				sb.WriteString(fmt.Sprintf("  %s - %s\n", name, desc))
			}
			return mcp.NewToolResultText(sb.String()), nil

		case "action":
			toolName := getStringArg(args, "toolName", "")
			actionName := getStringArg(args, "actionName", "")
			if toolName == "" || actionName == "" {
				return mcp.NewToolResultError("both toolName and actionName are required"), nil
			}
			actions, ok := toolDescriptions[toolName]
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("unknown tool: %s", toolName)), nil
			}
			desc, ok := actions[actionName]
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("unknown action %s in tool %s", actionName, toolName)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("%s.%s: %s", toolName, actionName, desc)), nil

		case "all":
			var sb strings.Builder
			sb.WriteString("AI Happy Design v2 - Available Tools and Actions\n")
			sb.WriteString("=================================================\n\n")
			sb.WriteString("Tip: for LLM-friendly examples and playbook, call describe with action='catalog'.\n\n")
			for toolName, actions := range toolDescriptions {
				sb.WriteString(fmt.Sprintf("%s:\n", toolName))
				for name, desc := range actions {
					sb.WriteString(fmt.Sprintf("  %s - %s\n", name, desc))
				}
				sb.WriteString("\n")
			}
			return mcp.NewToolResultText(sb.String()), nil

		case "catalog":
			catalog := LLMCatalog()
			data, err := json.MarshalIndent(catalog, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to build catalog: %v", err)), nil
			}
			return mcp.NewToolResultText(string(data)), nil

		case "design_guide":
			guide := DesignGuide()
			data, err := json.MarshalIndent(guide, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to build design guide: %v", err)), nil
			}
			return mcp.NewToolResultText(string(data)), nil

		case "setup":
			return mcp.NewToolResultText(SetupInstructions()), nil

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown describe action: %s", action)), nil
		}
	})
}
