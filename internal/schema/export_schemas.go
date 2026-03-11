package schema

func init() {
	Register(Schema{
		Command:     "export.image",
		Description: "Export a node as PNG or JPG image. Auto-saves to OS temp dir and returns file path.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Node ID to export", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "format", Type: "string", Desc: "Image format", Enum: []string{"PNG", "JPG"}, Default: "PNG"},
			{Name: "scale", Type: "number", Desc: "Export scale factor", Min: Ptr(0.01), Max: Ptr(4), Default: 2.0},
		},
	})

	Register(Schema{
		Command:     "export.svg",
		Description: "Export a node as SVG",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Node ID to export", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "svgIdAttribute", Type: "boolean", Desc: "Include id attributes in SVG elements"},
			{Name: "svgOutlineText", Type: "boolean", Desc: "Outline text as paths"},
			{Name: "svgSimplifyStroke", Type: "boolean", Desc: "Simplify strokes"},
		},
	})

	Register(Schema{
		Command:     "export.pdf",
		Description: "Export a node as PDF",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Node ID to export", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "export.batch",
		Aliases:     []string{"export.batch_export"},
		Description: "Export multiple frames at once. If no nodeIds, exports all top-level frames.",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Desc: "Comma-separated node IDs to export (all top-level frames if omitted)"},
			{Name: "format", Type: "string", Desc: "Export format", Enum: []string{"PNG", "JPG", "SVG"}, Default: "PNG"},
			{Name: "scale", Type: "number", Desc: "Export scale factor", Min: Ptr(0.01), Max: Ptr(4), Default: 2.0},
		},
	})
}
