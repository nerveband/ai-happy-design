package schema

func init() {
	Register(Schema{
		Command:     "shape.create_rectangle",
		Aliases:     []string{"rect"},
		Description: "Create a rectangle",
		Params: []Param{
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
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
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
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
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
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
