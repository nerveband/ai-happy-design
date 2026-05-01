package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	nodeIDs := Param{Name: "nodeIds", Type: "array", Required: true}
	parentID := Param{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`}

	for _, s := range []Schema{
		{Command: "draw.create_text_path", Description: "Beta Figma API: create text on a vector path with runtime guards", Params: []Param{nodeID, {Name: "characters", Type: "string"}, {Name: "startSegment", Type: "number"}, {Name: "startPosition", Type: "number", Min: Ptr(0), Max: Ptr(1)}, {Name: "fontSize", Type: "number"}, {Name: "name", Type: "string"}}},
		{Command: "draw.create_transform_group", Description: "Beta Figma API: create a transform group with runtime guards", Params: []Param{nodeIDs, parentID, {Name: "index", Type: "number", Min: Ptr(0)}, {Name: "modifiers", Type: "array"}, {Name: "name", Type: "string"}}},
		{Command: "draw.load_brushes", Description: "Beta Figma API: load built-in Draw brushes for this plugin run", Params: []Param{{Name: "brushType", Type: "string", Enum: []string{"STRETCH", "SCATTER"}}, {Name: "brushTypes", Type: "array"}}},
		{Command: "draw.set_variable_width_stroke", Description: "Beta Figma API: set variable-width stroke properties with runtime guards", Params: []Param{nodeID, {Name: "widthProfile", Type: "string", Enum: []string{"UNIFORM", "WEDGE", "TAPER", "QUARTER_TAPER", "EYE", "MIRRORED_TAPER", "CUSTOM"}}, {Name: "variableWidthPoints", Type: "array"}}},
		{Command: "draw.set_brush_stroke", Description: "Beta Figma API: set stretch or scatter brush stroke with runtime guards", Params: []Param{nodeID, {Name: "brushType", Type: "string", Enum: []string{"STRETCH", "SCATTER"}}, {Name: "brushName", Type: "string"}, {Name: "gap", Type: "number"}, {Name: "wiggle", Type: "number"}, {Name: "direction", Type: "string"}}},
		{Command: "draw.set_dynamic_stroke", Description: "Beta Figma API: set dynamic stroke properties with runtime guards", Params: []Param{nodeID, {Name: "frequency", Type: "number"}, {Name: "wiggle", Type: "number"}, {Name: "smoothen", Type: "number"}}},
		{Command: "draw.set_pattern_fill", Description: "Beta Figma API: set pattern fill using async paint APIs", Params: []Param{nodeID, {Name: "sourceNodeId", Type: "string", Required: true}, {Name: "tileType", Type: "string"}, {Name: "scalingFactor", Type: "number"}, {Name: "append", Type: "boolean"}}},
		{Command: "draw.set_pattern_stroke", Description: "Beta Figma API: set pattern stroke using async paint APIs", Params: []Param{nodeID, {Name: "sourceNodeId", Type: "string", Required: true}, {Name: "tileType", Type: "string"}, {Name: "scalingFactor", Type: "number"}, {Name: "append", Type: "boolean"}}},
	} {
		Register(s)
	}
}
