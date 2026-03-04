package schema

func init() {
	Register(Schema{
		Command:     "effect.add_shadow",
		Aliases:     []string{"shadow"},
		Description: "Add a drop shadow to a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Desc: "Shadow color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`, Default: "#00000040"},
			{Name: "offsetX", Type: "number", Desc: "Horizontal offset", Default: 0.0},
			{Name: "offsetY", Type: "number", Desc: "Vertical offset", Default: 4.0},
			{Name: "radius", Type: "number", Desc: "Blur radius", Min: Ptr(0), Max: Ptr(200), Default: 4.0},
			{Name: "spread", Type: "number", Desc: "Spread distance", Min: Ptr(-100), Max: Ptr(100), Default: 0.0},
			{Name: "type", Type: "string", Enum: []string{"DROP_SHADOW", "INNER_SHADOW"}, Default: "DROP_SHADOW"},
		},
	})

	Register(Schema{
		Command:     "effect.add_blur",
		Aliases:     []string{"blur"},
		Description: "Add blur to a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "radius", Type: "number", Min: Ptr(0), Max: Ptr(200), Default: 10.0},
			{Name: "type", Type: "string", Enum: []string{"LAYER_BLUR", "BACKGROUND_BLUR"}, Default: "LAYER_BLUR"},
		},
	})

	Register(Schema{
		Command:     "effect.apply_glass",
		Aliases:     []string{"glass"},
		Description: "Apply glass morphism effect (background blur + semi-transparent fill + stroke)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "intensity", Type: "string", Enum: []string{"light", "medium", "heavy"}, Default: "medium"},
		},
	})
}
