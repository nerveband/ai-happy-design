package schema

func init() {
	Register(Schema{
		Command:     "shape.create",
		Description: "Create a shape by type. Agent-friendly wrapper around the specific shape.create_* commands.",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "type", Type: "string", Enum: []string{"RECTANGLE", "RECT", "ELLIPSE", "CIRCLE", "POLYGON", "STAR", "LINE", "VECTOR", "SVG", "IMAGE"}, Default: "RECTANGLE"},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Aliases: []string{"bg"}, Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "fillColor", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "stroke", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWidth", Type: "number", Min: Ptr(0), Max: Ptr(1000)},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "pointCount", Type: "number", Min: Ptr(3), Max: Ptr(100)},
			{Name: "sides", Type: "number", Min: Ptr(3), Max: Ptr(100)},
			{Name: "points", Type: "number", Min: Ptr(3), Max: Ptr(100)},
			{Name: "innerRadius", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "startX", Type: "number"}, {Name: "startY", Type: "number"},
			{Name: "endX", Type: "number"}, {Name: "endY", Type: "number"},
			{Name: "length", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "rotation", Type: "number"},
			{Name: "svg", Type: "string"},
			{Name: "svgPath", Type: "string"},
			{Name: "svgString", Type: "string"},
			{Name: "imageData", Type: "string"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})

	Register(Schema{
		Command:     "shape.create_rectangle",
		Aliases:     []string{"rect"},
		Description: "Create a rectangle",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Aliases: []string{"bg"}, Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "constrainProportions", Type: "boolean", Desc: "Lock aspect ratio"},
		},
	})

	Register(Schema{
		Command:     "shape.create_ellipse",
		Description: "Create an ellipse or circle",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})

	Register(Schema{
		Command:     "shape.create_image",
		Description: "Create a rectangle with an image fill (one-step convenience)",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "imageData", Type: "string", Required: true, Desc: "Base64 image data, data URL, file path, or HTTP(S) URL"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})

	Register(Schema{
		Command:     "shape.create_polygon",
		Description: "Create a polygon shape",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Aliases: []string{"fillColor"}, Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "pointCount", Type: "number", Aliases: []string{"sides"}, Desc: "Number of sides", Min: Ptr(3), Max: Ptr(100), Default: 6.0},
		},
	})

	Register(Schema{
		Command:     "shape.create_star",
		Description: "Create a star shape",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Aliases: []string{"fillColor"}, Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "pointCount", Type: "number", Aliases: []string{"points"}, Desc: "Number of points", Min: Ptr(3), Max: Ptr(100), Default: 5.0},
			{Name: "innerRadius", Type: "number", Desc: "Inner radius ratio (0-1)", Min: Ptr(0), Max: Ptr(1), Default: 0.4},
		},
	})

	Register(Schema{
		Command:     "shape.create_line",
		Description: "Create a line",
		Params: []Param{
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number", Desc: "X position (or use startX/startY/endX/endY)"},
			{Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"length"}, Desc: "Line length", Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Aliases: []string{"stroke", "strokeColor"}, Desc: "Line color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWeight", Type: "number", Aliases: []string{"strokeWidth"}, Desc: "Line thickness", Min: Ptr(0), Max: Ptr(100), Default: 1.0},
			{Name: "rotation", Type: "number", Desc: "Rotation angle in degrees"},
			{Name: "startX", Type: "number", Desc: "Start X (alternative to x/y/width)"},
			{Name: "startY", Type: "number", Desc: "Start Y"},
			{Name: "endX", Type: "number", Desc: "End X"},
			{Name: "endY", Type: "number", Desc: "End Y"},
		},
	})

	Register(Schema{
		Command:     "shape.create_from_svg",
		Description: "Create a shape from an SVG string",
		Params: []Param{
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "svg", Type: "string", Required: true, Aliases: []string{"svgPath", "svgString"}, Desc: "SVG string to import"},
		},
	})

	Register(Schema{
		Command:     "shape.create_vector",
		Description: "Create a vector shape with path data",
		Params: []Param{
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "vectorPaths", Type: "array", Aliases: []string{"paths"}, Desc: "Array of vector path objects [{windingRule, data}]"},
			{Name: "vectorNetwork", Type: "object", Desc: "Vector network {vertices, segments, regions}"},
			{Name: "color", Type: "string", Aliases: []string{"fillColor"}, Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeColor", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWeight", Type: "number", Min: Ptr(0), Max: Ptr(100)},
		},
	})
}
