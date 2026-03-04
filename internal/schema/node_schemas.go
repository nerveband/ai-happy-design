package schema

func init() {
	Register(Schema{
		Command:     "node.create_frame",
		Aliases:     []string{"frame"},
		Description: "Create a frame (Figma's equivalent of a div)",
		Params: []Param{
			{Name: "name", Type: "string", Desc: "Semantic name for the frame (never leave as default)"},
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Desc: "Parent node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "x", Type: "number", Desc: "X position"},
			{Name: "y", Type: "number", Desc: "Y position"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Desc: "Frame width", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Desc: "Frame height", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "color", Type: "string", Aliases: []string{"bg"}, Desc: "Fill color hex (NOT fillColor)", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Desc: "Corner radius", Min: Ptr(0), Max: Ptr(500)},
			{Name: "opacity", Type: "number", Desc: "Opacity", Min: Ptr(0), Max: Ptr(1)},
			{Name: "layoutMode", Type: "string", Desc: "Auto-layout direction", Enum: []string{"HORIZONTAL", "VERTICAL"}},
			{Name: "itemSpacing", Type: "number", Desc: "Spacing between auto-layout children", Min: Ptr(0), Max: Ptr(500)},
			{Name: "padding", Type: "number", Desc: "Uniform padding for all sides", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingTop", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingRight", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingBottom", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingLeft", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "primaryAxisAlignItems", Type: "string", Enum: []string{"MIN", "CENTER", "MAX", "SPACE_BETWEEN"}},
			{Name: "counterAxisAlignItems", Type: "string", Enum: []string{"MIN", "CENTER", "MAX", "BASELINE"}},
			{Name: "layoutSizingHorizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutSizingVertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutWrap", Type: "string", Enum: []string{"WRAP", "NO_WRAP"}},
			{Name: "clipsContent", Type: "boolean", Desc: "Whether children clip to frame bounds"},
			{Name: "structural", Type: "boolean", Desc: "Structural frame: removes default fill, enables clipping. Use for wrappers and auto-layout containers."},
			{Name: "stroke", Type: "string", Desc: "Stroke color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWeight", Type: "number", Aliases: []string{"sw"}, Min: Ptr(0), Max: Ptr(100)},
			{Name: "minWidth", Type: "number", Desc: "Minimum width (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxWidth", Type: "number", Desc: "Maximum width (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "minHeight", Type: "number", Desc: "Minimum height (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxHeight", Type: "number", Desc: "Maximum height (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
		},
	})

	Register(Schema{
		Command:     "node.modify",
		Aliases:     []string{"modify"},
		Description: "Modify any property of an existing node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "color", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "visible", Type: "boolean"},
			{Name: "name", Type: "string"},
			{Name: "rotation", Type: "number", Min: Ptr(0), Max: Ptr(360)},
			{Name: "text", Type: "string"},
			{Name: "fontSize", Type: "number", Aliases: []string{"sz"}, Min: Ptr(4), Max: Ptr(500), SemanticTokens: true},
			{Name: "fontFamily", Type: "string", Aliases: []string{"ff"}},
			{Name: "isMask", Type: "boolean"},
			{Name: "layoutSizingHorizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutSizingVertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "blendMode", Type: "string", Enum: []string{"NORMAL", "MULTIPLY", "SCREEN", "OVERLAY", "DARKEN", "LIGHTEN", "COLOR_DODGE", "COLOR_BURN", "HARD_LIGHT", "SOFT_LIGHT", "DIFFERENCE", "EXCLUSION", "HUE", "SATURATION", "COLOR", "LUMINOSITY"}},
			{Name: "minWidth", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxWidth", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "minHeight", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxHeight", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})

	Register(Schema{
		Command:     "node.move",
		Description: "Move a node to new coordinates",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "x", Type: "number", Required: true}, {Name: "y", Type: "number", Required: true},
		},
	})

	Register(Schema{
		Command:     "node.resize",
		Description: "Resize a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "width", Type: "number", Required: true, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Required: true, Min: Ptr(1), Max: Ptr(10000)},
		},
	})

	Register(Schema{
		Command:     "node.delete",
		Description: "Delete a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "node.get_info",
		Description: "Get information about a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "node.get_tree",
		Description: "Get the node tree hierarchy",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "depth", Type: "number", Desc: "Max depth to traverse", Min: Ptr(1), Max: Ptr(20)},
			{Name: "compact", Type: "boolean", Desc: "Return flat array instead of nested tree (3-5x fewer tokens)"},
		},
	})
}
