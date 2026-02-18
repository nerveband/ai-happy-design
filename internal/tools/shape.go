package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		mcp.WithString("stroke", mcp.Description("Stroke color as hex string")),
		mcp.WithNumber("strokeWidth", mcp.Description("Stroke width in px")),
		mcp.WithBoolean("noFill", mcp.Description("Remove fill on create")),
		mcp.WithNumber("opacity", mcp.Description("Opacity from 0 to 1")),
		mcp.WithNumber("cornerRadius", mcp.Description("Corner radius for rectangles")),
		mcp.WithNumber("pointCount", mcp.Description("Number of points for polygon/star")),
		mcp.WithNumber("innerRadius", mcp.Description("Inner radius ratio for stars (0-1)")),
		mcp.WithString("svgPath", mcp.Description("SVG path data string")),
		mcp.WithString("imageData", mcp.Description("Base64-encoded image data (PNG/JPG) for create_image")),
		mcp.WithString("scaleMode", mcp.Description("Image scale mode: FILL, FIT, CROP, TILE (default FILL)")),
		mcp.WithString("layoutSizingHorizontal", mcp.Description("Horizontal sizing in auto-layout: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL")),
		mcp.WithString("layoutSizingVertical", mcp.Description("Vertical sizing in auto-layout: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL")),
		mcp.WithString("layoutAlign", mcp.Description("Layout alignment: STRETCH or INHERIT"),
			mcp.Enum("STRETCH", "INHERIT")),
		mcp.WithNumber("layoutGrow", mcp.Description("Layout grow: 0 or 1")),
		mcp.WithString("layoutPositioning", mcp.Description("Layout positioning: ABSOLUTE for decorative overlays"),
			mcp.Enum("AUTO", "ABSOLUTE")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create_rectangle":
			params := map[string]interface{}{
				"name":         getStringArg(args, "name", "Rectangle"),
				"x":            getFloat64Arg(args, "x", 0),
				"y":            getFloat64Arg(args, "y", 0),
				"width":        getFloat64Arg(args, "width", 100),
				"height":       getFloat64Arg(args, "height", 100),
				"parentId":     getStringArg(args, "parentId", ""),
				"color":        getStringArg(args, "color", ""),
				"cornerRadius": getFloat64Arg(args, "cornerRadius", 0),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "noFill") {
				params["noFill"] = getBoolArg(args, "noFill", false)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "layoutSizingHorizontal") {
				params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "")
			}
			if hasArg(args, "layoutSizingVertical") {
				params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "")
			}
			if hasArg(args, "layoutAlign") {
				params["layoutAlign"] = getStringArg(args, "layoutAlign", "")
			}
			if hasArg(args, "layoutGrow") {
				params["layoutGrow"] = getFloat64Arg(args, "layoutGrow", 0)
			}
			if hasArg(args, "layoutPositioning") {
				params["layoutPositioning"] = getStringArg(args, "layoutPositioning", "")
			}
			return sendCommand(commander, "create_rectangle", params)

		case "create_ellipse":
			params := map[string]interface{}{
				"name":     getStringArg(args, "name", "Ellipse"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"width":    getFloat64Arg(args, "width", 100),
				"height":   getFloat64Arg(args, "height", 100),
				"parentId": getStringArg(args, "parentId", ""),
				"color":    getStringArg(args, "color", ""),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "noFill") {
				params["noFill"] = getBoolArg(args, "noFill", false)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "layoutSizingHorizontal") {
				params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "")
			}
			if hasArg(args, "layoutSizingVertical") {
				params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "")
			}
			if hasArg(args, "layoutAlign") {
				params["layoutAlign"] = getStringArg(args, "layoutAlign", "")
			}
			if hasArg(args, "layoutGrow") {
				params["layoutGrow"] = getFloat64Arg(args, "layoutGrow", 0)
			}
			if hasArg(args, "layoutPositioning") {
				params["layoutPositioning"] = getStringArg(args, "layoutPositioning", "")
			}
			return sendCommand(commander, "create_ellipse", params)

		case "create_polygon":
			params := map[string]interface{}{
				"name":       getStringArg(args, "name", "Polygon"),
				"x":          getFloat64Arg(args, "x", 0),
				"y":          getFloat64Arg(args, "y", 0),
				"width":      getFloat64Arg(args, "width", 100),
				"height":     getFloat64Arg(args, "height", 100),
				"parentId":   getStringArg(args, "parentId", ""),
				"color":      getStringArg(args, "color", ""),
				"pointCount": getFloat64Arg(args, "pointCount", 6),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "noFill") {
				params["noFill"] = getBoolArg(args, "noFill", false)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "layoutSizingHorizontal") {
				params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "")
			}
			if hasArg(args, "layoutSizingVertical") {
				params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "")
			}
			if hasArg(args, "layoutAlign") {
				params["layoutAlign"] = getStringArg(args, "layoutAlign", "")
			}
			if hasArg(args, "layoutGrow") {
				params["layoutGrow"] = getFloat64Arg(args, "layoutGrow", 0)
			}
			if hasArg(args, "layoutPositioning") {
				params["layoutPositioning"] = getStringArg(args, "layoutPositioning", "")
			}
			return sendCommand(commander, "create_polygon", params)

		case "create_star":
			params := map[string]interface{}{
				"name":        getStringArg(args, "name", "Star"),
				"x":           getFloat64Arg(args, "x", 0),
				"y":           getFloat64Arg(args, "y", 0),
				"width":       getFloat64Arg(args, "width", 100),
				"height":      getFloat64Arg(args, "height", 100),
				"parentId":    getStringArg(args, "parentId", ""),
				"color":       getStringArg(args, "color", ""),
				"pointCount":  getFloat64Arg(args, "pointCount", 5),
				"innerRadius": getFloat64Arg(args, "innerRadius", 0.5),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "noFill") {
				params["noFill"] = getBoolArg(args, "noFill", false)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "layoutSizingHorizontal") {
				params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "")
			}
			if hasArg(args, "layoutSizingVertical") {
				params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "")
			}
			if hasArg(args, "layoutAlign") {
				params["layoutAlign"] = getStringArg(args, "layoutAlign", "")
			}
			if hasArg(args, "layoutGrow") {
				params["layoutGrow"] = getFloat64Arg(args, "layoutGrow", 0)
			}
			if hasArg(args, "layoutPositioning") {
				params["layoutPositioning"] = getStringArg(args, "layoutPositioning", "")
			}
			return sendCommand(commander, "create_star", params)

		case "create_line":
			params := map[string]interface{}{
				"name":     getStringArg(args, "name", "Line"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"width":    getFloat64Arg(args, "width", 100),
				"parentId": getStringArg(args, "parentId", ""),
				"color":    getStringArg(args, "color", ""),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			return sendCommand(commander, "create_line", params)

		case "create_from_svg":
			svgPath, errResult := requireStringArg(args, "svgPath")
			if errResult != nil {
				return errResult, nil
			}
			// Auto-resolve file paths: /path/to/icon.svg or ~/icon.svg → read SVG content
			svgContent := svgPath
			if strings.HasPrefix(svgPath, "file://") || strings.HasPrefix(svgPath, "/") || strings.HasPrefix(svgPath, "~/") {
				rawPath := strings.TrimPrefix(svgPath, "file://")
				if strings.HasPrefix(rawPath, "~/") {
					home, _ := os.UserHomeDir()
					rawPath = filepath.Join(home, rawPath[2:])
				}
				data, readErr := os.ReadFile(rawPath)
				if readErr != nil {
					return mcp.NewToolResultError(fmt.Sprintf("cannot read SVG file %s: %v", rawPath, readErr)), nil
				}
				svgContent = string(data)
			}
			return sendCommand(commander, "create_from_svg", map[string]interface{}{
				"name":     getStringArg(args, "name", "SVG"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"svgPath":  svgContent,
				"parentId": getStringArg(args, "parentId", ""),
			})

		case "create_image":
			imageData, errResult := requireStringArg(args, "imageData")
			if errResult != nil {
				return errResult, nil
			}
			params := map[string]interface{}{
				"name":         getStringArg(args, "name", "Image"),
				"x":            getFloat64Arg(args, "x", 0),
				"y":            getFloat64Arg(args, "y", 0),
				"width":        getFloat64Arg(args, "width", 400),
				"height":       getFloat64Arg(args, "height", 400),
				"parentId":     getStringArg(args, "parentId", ""),
				"imageData":    imageData,
				"scaleMode":    getStringArg(args, "scaleMode", "FILL"),
				"cornerRadius": getFloat64Arg(args, "cornerRadius", 0),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "noFill") {
				params["noFill"] = getBoolArg(args, "noFill", false)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "layoutSizingHorizontal") {
				params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "")
			}
			if hasArg(args, "layoutSizingVertical") {
				params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "")
			}
			if hasArg(args, "layoutAlign") {
				params["layoutAlign"] = getStringArg(args, "layoutAlign", "")
			}
			if hasArg(args, "layoutGrow") {
				params["layoutGrow"] = getFloat64Arg(args, "layoutGrow", 0)
			}
			if hasArg(args, "layoutPositioning") {
				params["layoutPositioning"] = getStringArg(args, "layoutPositioning", "")
			}
			return sendCommand(commander, "create_image", params)

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown shape action: %s", action)), nil
		}
	})
}
