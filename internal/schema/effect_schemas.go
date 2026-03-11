package schema

func init() {
	Register(Schema{
		Command:     "effect.add",
		Description: "Add an effect by type. Agent-friendly wrapper around the specific effect.add_* commands.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "type", Type: "string", Enum: []string{"DROP_SHADOW", "SHADOW", "INNER_SHADOW", "LAYER_BLUR", "BACKGROUND_BLUR", "BLUR", "NOISE", "TEXTURE", "GLASS"}, Default: "DROP_SHADOW"},
			{Name: "color", Type: "string", Desc: "Effect color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`, Default: "#00000040"},
			{Name: "offsetX", Type: "number", Desc: "Horizontal offset", Default: 0.0},
			{Name: "offsetY", Type: "number", Desc: "Vertical offset", Default: 4.0},
			{Name: "radius", Type: "number", Min: Ptr(0), Max: Ptr(200), Default: 4.0},
			{Name: "spread", Type: "number", Min: Ptr(-100), Max: Ptr(100), Default: 0.0},
			{Name: "noiseType", Type: "string", Enum: []string{"MONOTONE", "DUOTONE", "MULTITONE"}},
			{Name: "secondaryColor", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "noiseSize", Type: "number", Min: Ptr(0), Max: Ptr(1000)},
			{Name: "density", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "blendMode", Type: "string"},
			{Name: "visible", Type: "boolean"},
			{Name: "intensity", Type: "string", Enum: []string{"light", "medium", "heavy"}},
		},
	})

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

	Register(Schema{
		Command:     "effect.set_effects",
		Description: "Set all effects on a node (replaces existing effects)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "effects", Type: "string", Required: true, Desc: "JSON array of effect objects (or actual array). Each: {type, color, offset, radius, spread, visible, blendMode}"},
		},
	})

	Register(Schema{
		Command:     "effect.apply_style",
		Description: "Apply an effect style to a node by style ID",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "styleId", Type: "string", Required: true, Desc: "Effect style ID to apply"},
		},
	})

	Register(Schema{
		Command:     "effect.remove",
		Aliases:     []string{"effect.clear_effects", "effect.remove_all"},
		Description: "Remove effects from a node. If index specified, removes only that effect; otherwise removes all.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "index", Type: "number", Desc: "Index of specific effect to remove (omit to remove all)"},
		},
	})

	Register(Schema{
		Command:     "effect.get_effects",
		Aliases:     []string{"effect.list"},
		Description: "Get all effects on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "effect.add_noise",
		Aliases:     []string{"noise"},
		Description: "Add a noise overlay effect (requires Figma Beta API)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "noiseType", Type: "string", Desc: "Noise type", Enum: []string{"MONOTONE", "DUOTONE", "MULTITONE"}, Default: "MONOTONE"},
			{Name: "color", Type: "string", Desc: "Noise color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "secondaryColor", Type: "string", Desc: "Secondary color for DUOTONE", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "noiseSize", Type: "number", Desc: "Noise grain size", Min: Ptr(0), Max: Ptr(1000), Default: 100.0},
			{Name: "density", Type: "number", Desc: "Noise density", Min: Ptr(0), Max: Ptr(1), Default: 0.3},
			{Name: "opacity", Type: "number", Desc: "Noise opacity (MULTITONE only)", Min: Ptr(0), Max: Ptr(1)},
			{Name: "blendMode", Type: "string", Desc: "Blend mode for noise", Default: "SOFT_LIGHT"},
			{Name: "visible", Type: "boolean", Default: true},
		},
	})

	Register(Schema{
		Command:     "effect.add_texture",
		Aliases:     []string{"texture"},
		Description: "Add a texture effect (requires Figma Beta API)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "noiseSize", Type: "number", Desc: "Texture size", Default: 100.0},
			{Name: "radius", Type: "number", Desc: "Texture radius", Default: 0.0},
			{Name: "clipToShape", Type: "boolean", Desc: "Clip texture to shape bounds", Default: true},
			{Name: "visible", Type: "boolean", Default: true},
		},
	})
}
