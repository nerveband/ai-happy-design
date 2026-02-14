package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterShapeTool registers the "shape" tool for creating geometric shapes.
func RegisterShapeTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("shape",
		mcp.WithDescription("Create geometric shapes: rectangles, ellipses, polygons, stars, lines, and SVG paths."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create_rectangle", "create_ellipse", "create_polygon", "create_star", "create_line", "create_from_svg", "create_image")),
		mcp.WithString("name", mcp.Description("Name for the created shape")),
		mcp.WithNumber("x", mcp.Description("X position")),
		mcp.WithNumber("y", mcp.Description("Y position")),
		mcp.WithNumber("width", mcp.Description("Width of the shape")),
		mcp.WithNumber("height", mcp.Description("Height of the shape")),
		mcp.WithString("parentId", mcp.Description("Parent node ID to insert into")),
		mcp.WithString("color", mcp.Description("Fill color as hex string")),
		mcp.WithNumber("cornerRadius", mcp.Description("Corner radius for rectangles")),
		mcp.WithNumber("pointCount", mcp.Description("Number of points for polygon/star")),
		mcp.WithNumber("innerRadius", mcp.Description("Inner radius ratio for stars (0-1)")),
		mcp.WithString("svgPath", mcp.Description("SVG path data string")),
		mcp.WithString("imageData", mcp.Description("Base64-encoded image data (PNG/JPG) for create_image")),
		mcp.WithString("scaleMode", mcp.Description("Image scale mode: FILL, FIT, CROP, TILE (default FILL)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create_rectangle":
			return sendCommand(commander, "create_rectangle", map[string]interface{}{
				"name":         getStringArg(args, "name", "Rectangle"),
				"x":            getFloat64Arg(args, "x", 0),
				"y":            getFloat64Arg(args, "y", 0),
				"width":        getFloat64Arg(args, "width", 100),
				"height":       getFloat64Arg(args, "height", 100),
				"parentId":     getStringArg(args, "parentId", ""),
				"color":        getStringArg(args, "color", ""),
				"cornerRadius": getFloat64Arg(args, "cornerRadius", 0),
			})

		case "create_ellipse":
			return sendCommand(commander, "create_ellipse", map[string]interface{}{
				"name":     getStringArg(args, "name", "Ellipse"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"width":    getFloat64Arg(args, "width", 100),
				"height":   getFloat64Arg(args, "height", 100),
				"parentId": getStringArg(args, "parentId", ""),
				"color":    getStringArg(args, "color", ""),
			})

		case "create_polygon":
			return sendCommand(commander, "create_polygon", map[string]interface{}{
				"name":       getStringArg(args, "name", "Polygon"),
				"x":          getFloat64Arg(args, "x", 0),
				"y":          getFloat64Arg(args, "y", 0),
				"width":      getFloat64Arg(args, "width", 100),
				"height":     getFloat64Arg(args, "height", 100),
				"parentId":   getStringArg(args, "parentId", ""),
				"color":      getStringArg(args, "color", ""),
				"pointCount": getFloat64Arg(args, "pointCount", 6),
			})

		case "create_star":
			return sendCommand(commander, "create_star", map[string]interface{}{
				"name":        getStringArg(args, "name", "Star"),
				"x":           getFloat64Arg(args, "x", 0),
				"y":           getFloat64Arg(args, "y", 0),
				"width":       getFloat64Arg(args, "width", 100),
				"height":      getFloat64Arg(args, "height", 100),
				"parentId":    getStringArg(args, "parentId", ""),
				"color":       getStringArg(args, "color", ""),
				"pointCount":  getFloat64Arg(args, "pointCount", 5),
				"innerRadius": getFloat64Arg(args, "innerRadius", 0.5),
			})

		case "create_line":
			return sendCommand(commander, "create_line", map[string]interface{}{
				"name":     getStringArg(args, "name", "Line"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"width":    getFloat64Arg(args, "width", 100),
				"parentId": getStringArg(args, "parentId", ""),
				"color":    getStringArg(args, "color", ""),
			})

		case "create_from_svg":
			svgPath, errResult := requireStringArg(args, "svgPath")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_from_svg", map[string]interface{}{
				"name":     getStringArg(args, "name", "SVG"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"svgPath":  svgPath,
				"parentId": getStringArg(args, "parentId", ""),
			})

		case "create_image":
			imageData, errResult := requireStringArg(args, "imageData")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_image", map[string]interface{}{
				"name":         getStringArg(args, "name", "Image"),
				"x":            getFloat64Arg(args, "x", 0),
				"y":            getFloat64Arg(args, "y", 0),
				"width":        getFloat64Arg(args, "width", 400),
				"height":       getFloat64Arg(args, "height", 400),
				"parentId":     getStringArg(args, "parentId", ""),
				"imageData":    imageData,
				"scaleMode":    getStringArg(args, "scaleMode", "FILL"),
				"cornerRadius": getFloat64Arg(args, "cornerRadius", 0),
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown shape action: %s", action)), nil
		}
	})
}
