package schema

func init() {
	Register(Schema{
		Command:     "layout.set_auto_layout",
		Description: "Set auto-layout properties on a frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "direction", Type: "string", Required: true, Enum: []string{"HORIZONTAL", "VERTICAL", "NONE"}, Desc: "HORIZONTAL=row, VERTICAL=column, NONE=remove auto-layout"},
			{Name: "itemSpacing", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "padding", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingTop", Type: "number", Min: Ptr(0)}, {Name: "paddingRight", Type: "number", Min: Ptr(0)},
			{Name: "paddingBottom", Type: "number", Min: Ptr(0)}, {Name: "paddingLeft", Type: "number", Min: Ptr(0)},
			{Name: "primaryAxisAlignItems", Type: "string", Aliases: []string{"primaryAxisAlign"}, Enum: []string{"MIN", "CENTER", "MAX", "SPACE_BETWEEN"}},
			{Name: "counterAxisAlignItems", Type: "string", Aliases: []string{"counterAxisAlign"}, Enum: []string{"MIN", "CENTER", "MAX", "BASELINE"}},
			{Name: "layoutWrap", Type: "string", Enum: []string{"WRAP", "NO_WRAP"}},
		},
	})

	Register(Schema{
		Command:     "layout.set_sizing",
		Description: "Set sizing behavior of a node within auto-layout",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "horizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "vertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
		},
	})

	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	for _, s := range []Schema{
		{Command: "layout.set_grid_container", Description: "Set Figma grid auto-layout container properties", Params: []Param{nodeID, {Name: "gridRowCount", Type: "number", Min: Ptr(1)}, {Name: "gridColumnCount", Type: "number", Min: Ptr(1)}, {Name: "gridRowGap", Type: "number", Min: Ptr(0)}, {Name: "gridColumnGap", Type: "number", Min: Ptr(0)}, {Name: "gridRowsSizing", Type: "array"}, {Name: "gridColumnsSizing", Type: "array"}}},
		{Command: "layout.set_grid_tracks", Description: "Set Figma grid auto-layout row and column track sizing", Params: []Param{nodeID, {Name: "gridRowsSizing", Type: "array"}, {Name: "gridColumnsSizing", Type: "array"}}},
		{Command: "layout.set_grid_child_position", Description: "Set a direct child position and span inside a Figma grid auto-layout container", Params: []Param{nodeID, {Name: "gridRowAnchorIndex", Type: "number", Min: Ptr(0)}, {Name: "gridColumnAnchorIndex", Type: "number", Min: Ptr(0)}, {Name: "gridRowSpan", Type: "number", Min: Ptr(1)}, {Name: "gridColumnSpan", Type: "number", Min: Ptr(1)}, {Name: "gridChildHorizontalAlign", Type: "string"}, {Name: "gridChildVerticalAlign", Type: "string"}}},
		{Command: "layout.get_grid_layout", Description: "Get Figma grid auto-layout container properties", Params: []Param{nodeID}},
	} {
		Register(s)
	}
}
