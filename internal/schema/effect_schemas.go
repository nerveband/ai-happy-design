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
			{Name: "blurType", Type: "string", Enum: []string{"LAYER_BLUR", "BACKGROUND_BLUR"}},
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

	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	Register(Schema{Command: "effect.set_effects", Description: "Set all effects on a node", Params: []Param{nodeID, {Name: "effects", Type: "array", Required: true}}})
	Register(Schema{Command: "effect.apply_style", Description: "Apply an effect style", Params: []Param{nodeID, {Name: "styleId", Type: "string", Required: true}}})
	Register(Schema{Command: "effect.remove_effect", Description: "Remove effects from a node", Params: []Param{nodeID}})
	Register(Schema{Command: "effect.get_effects", Description: "Get effects from a node", Params: []Param{nodeID}})
	Register(Schema{Command: "effect.add_noise", Description: "Add a noise effect where supported by Figma", Params: []Param{nodeID, {Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)}, {Name: "size", Type: "number", Min: Ptr(0)}, {Name: "noiseSize", Type: "number", Min: Ptr(0)}, {Name: "noiseSizeX", Type: "number", Min: Ptr(0)}, {Name: "noiseSizeY", Type: "number", Min: Ptr(0)}, {Name: "noiseSizeVector", Type: "object"}}})
	Register(Schema{Command: "effect.add_texture", Description: "Add a texture effect where supported by Figma", Params: []Param{nodeID, {Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)}, {Name: "size", Type: "number", Min: Ptr(0)}, {Name: "noiseSize", Type: "number", Min: Ptr(0)}, {Name: "noiseSizeX", Type: "number", Min: Ptr(0)}, {Name: "noiseSizeY", Type: "number", Min: Ptr(0)}, {Name: "noiseSizeVector", Type: "object"}}})
	Register(Schema{Command: "effect.add_glass", Description: "Add native glass effect where supported by Figma", Params: []Param{nodeID, {Name: "preset", Type: "string", Enum: []string{"light", "medium", "heavy"}}, {Name: "radius", Type: "number", Min: Ptr(0)}}})
	Register(Schema{Command: "effect.list_shaders", Description: "List local shader resources where supported by Figma", Params: nil, Safety: "read"})
	Register(Schema{Command: "effect.import_shader", Description: "Import a shader resource where supported by Figma", Params: []Param{{Name: "key", Type: "string"}, {Name: "url", Type: "string"}, {Name: "name", Type: "string"}}})
	Register(Schema{Command: "effect.apply_shader_effect", Description: "Apply a shader effect to a node where supported by Figma", Params: []Param{nodeID, {Name: "shaderId", Type: "string", Required: true}, {Name: "uniforms", Type: "object"}}})
}
