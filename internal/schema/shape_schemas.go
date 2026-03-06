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
}
