package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterPaintTool registers the "paint" tool for fill, stroke, and gradient operations.
func RegisterPaintTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("paint",
		mcp.WithDescription("Paint operations: fills, strokes, gradients, images. Use action parameter to select operation."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("set_solid", "set_gradient", "set_image", "set_image_fill", "set_image_url", "set_image_fill_from_url", "add_fill", "remove_fill", "get_fills", "set_stroke")),
		mcp.WithString("nodeId", mcp.Description("Target node ID")),
		mcp.WithString("color", mcp.Description("Hex color string (e.g. #FF0000)")),
		mcp.WithNumber("opacity", mcp.Description("Opacity from 0 to 1")),
		mcp.WithString("gradientType", mcp.Description("Gradient type: LINEAR, RADIAL, ANGULAR, DIAMOND"),
			mcp.Enum("LINEAR", "RADIAL", "ANGULAR", "DIAMOND")),
		mcp.WithString("stops", mcp.Description("JSON array of gradient stops [{position, color}]")),
		mcp.WithString("url", mcp.Description("URL of the image to use as fill")),
		mcp.WithString("imageUrl", mcp.Description("URL of the image to use as fill")),
		mcp.WithString("imageData", mcp.Description("Base64-encoded image data (raw base64 or data URL prefix).")),
		mcp.WithNumber("timeoutMs", mcp.Description("Timeout in milliseconds for URL image fetch/createImageAsync path.")),
		mcp.WithString("scaleMode", mcp.Description("Image scale mode: FILL, FIT, TILE, CROP"),
			mcp.Enum("FILL", "FIT", "TILE", "CROP")),
		mcp.WithNumber("strokeWeight", mcp.Description("Stroke weight in pixels")),
		mcp.WithString("strokeAlign", mcp.Description("Stroke alignment: INSIDE, OUTSIDE, CENTER"),
			mcp.Enum("INSIDE", "OUTSIDE", "CENTER")),
		mcp.WithNumber("fillIndex", mcp.Description("Index of fill to remove (for remove_fill)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "set_solid":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			color, errResult := requireStringArg(args, "color")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_fill_color", map[string]interface{}{
				"nodeId":  nodeId,
				"color":   color,
				"opacity": getFloat64Arg(args, "opacity", 1.0),
			})

		case "set_gradient":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_gradient", map[string]interface{}{
				"nodeId":       nodeId,
				"gradientType": getStringArg(args, "gradientType", "LINEAR"),
				"stops":        getStringArg(args, "stops", "[]"),
			})

		case "set_image", "set_image_fill":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_image_fill", map[string]interface{}{
				"nodeId":    nodeId,
				"imageData": getStringArg(args, "imageData", ""),
				"scaleMode": getStringArg(args, "scaleMode", "FILL"),
			})

		case "set_image_url", "set_image_fill_from_url":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			url := getStringArg(args, "url", "")
			if url == "" {
				url = getStringArg(args, "imageUrl", "")
			}
			return sendCommand(commander, "set_image_fill_url", map[string]interface{}{
				"nodeId":    nodeId,
				"url":       url,
				"imageUrl":  url,
				"scaleMode": getStringArg(args, "scaleMode", "FILL"),
				"timeoutMs": getFloat64Arg(args, "timeoutMs", 8000),
			})

		case "add_fill":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "add_fill", map[string]interface{}{
				"nodeId":  nodeId,
				"color":   getStringArg(args, "color", "#000000"),
				"opacity": getFloat64Arg(args, "opacity", 1.0),
			})

		case "remove_fill":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "remove_fill", map[string]interface{}{
				"nodeId": nodeId,
				"index":  getFloat64Arg(args, "fillIndex", 0),
			})

		case "get_fills":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "get_fills", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "set_stroke":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_stroke_color", map[string]interface{}{
				"nodeId":       nodeId,
				"color":        getStringArg(args, "color", "#000000"),
				"opacity":      getFloat64Arg(args, "opacity", 1.0),
				"strokeWeight": getFloat64Arg(args, "strokeWeight", 1),
				"strokeAlign":  getStringArg(args, "strokeAlign", "INSIDE"),
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown paint action: %s", action)), nil
		}
	})
}
