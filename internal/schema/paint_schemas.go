package schema

func init() {
	Register(Schema{
		Command:     "paint.set_solid",
		Aliases:     []string{"fill"},
		Description: "Set a solid fill color on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Required: true, Desc: "Hex color (#RRGGBB or #RGB)", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "opacity", Type: "number", Desc: "Opacity", Min: Ptr(0), Max: Ptr(1), Default: 1.0},
		},
	})

	Register(Schema{
		Command:     "paint.set_gradient",
		Aliases:     []string{"gradient"},
		Description: "Set a gradient fill on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "gradientType", Type: "string", Enum: []string{"LINEAR", "RADIAL", "ANGULAR", "DIAMOND"}, Default: "LINEAR"},
			{Name: "stops", Type: "string", Required: true, Desc: "JSON array of gradient stops: [{position:0, color:\"#FF0000\"}, {position:1, color:\"#0000FF\"}]"},
			{Name: "angle", Type: "number", Desc: "Gradient angle in degrees (0=top, 90=right)", Min: Ptr(0), Max: Ptr(360)},
		},
	})

	Register(Schema{
		Command:     "paint.set_stroke",
		Aliases:     []string{"stroke"},
		Description: "Set stroke on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Required: true, Desc: "Stroke color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWeight", Type: "number", Aliases: []string{"sw"}, Min: Ptr(0), Max: Ptr(100), Default: 1.0},
			{Name: "strokeAlign", Type: "string", Enum: []string{"INSIDE", "CENTER", "OUTSIDE"}, Default: "INSIDE"},
			{Name: "dashPattern", Type: "array", Desc: "Dash pattern [dash, gap] e.g. [10, 5]"},
			{Name: "strokeCap", Type: "string", Enum: []string{"NONE", "ROUND", "SQUARE", "ARROW_LINES", "ARROW_EQUILATERAL"}},
			{Name: "strokeJoin", Type: "string", Enum: []string{"MITER", "BEVEL", "ROUND"}},
		},
	})

	Register(Schema{
		Command:     "paint.set_image",
		Aliases:     []string{"paint.set_image_fill"},
		Description: "Set an image fill on a node from base64 data",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "imageData", Type: "string", Required: true, Desc: "Base64-encoded image data or data URL"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
		},
	})

	Register(Schema{
		Command:     "paint.set_image_url",
		Aliases:     []string{"paint.set_image_fill_from_url"},
		Description: "Set an image fill on a node from a URL. Plugin tries createImageAsync first, then fetch fallback.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "url", Type: "string", Required: true, Aliases: []string{"imageUrl"}, Desc: "URL of the image"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
			{Name: "timeoutMs", Type: "number", Desc: "Timeout in milliseconds for image fetch", Default: 8000.0},
		},
	})

	Register(Schema{
		Command:     "paint.add_fill",
		Description: "Add a fill to existing fills (does not replace)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Desc: "Fill color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "opacity", Type: "number", Desc: "Fill opacity", Min: Ptr(0), Max: Ptr(1)},
			{Name: "type", Type: "string", Desc: "Fill type", Enum: []string{"SOLID", "GRADIENT_LINEAR", "GRADIENT_RADIAL", "GRADIENT_ANGULAR", "GRADIENT_DIAMOND", "IMAGE"}, Default: "SOLID"},
			{Name: "stops", Type: "string", Desc: "Gradient stops JSON (for gradient types)"},
			{Name: "angle", Type: "number", Aliases: []string{"rotation"}, Desc: "Gradient angle (for gradient types)"},
			{Name: "imageData", Type: "string", Desc: "Base64 image data (for IMAGE type)"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
			{Name: "index", Type: "number", Desc: "Insert position in fill stack (appends if omitted)"},
		},
	})

	Register(Schema{
		Command:     "paint.remove_fill",
		Aliases:     []string{"paint.remove_fills", "paint.clear_fill", "paint.delete_fill"},
		Description: "Remove a fill by index",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "index", Type: "number", Aliases: []string{"fillIndex"}, Desc: "Index of fill to remove (default 0)", Default: 0.0},
		},
	})

	Register(Schema{
		Command:     "paint.get_fills",
		Description: "Get all fills on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})
}
