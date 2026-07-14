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

	Register(Schema{
		Command:     "layout.audit",
		Description: "Read-only layout audit for a node subtree. Measures bounds, overflow, clipping, overlap, gaps, and text sizing; review fixes, apply one batch, then re-audit before taking a screenshot.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "depth", Type: "number", Desc: "Maximum subtree depth to inspect", Min: Ptr(0), Max: Ptr(50)},
			{Name: "maxNodes", Type: "number", Desc: "Maximum nodes to inspect", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "minGap", Type: "number", Desc: "Minimum recommended gap between aligned siblings", Min: Ptr(0), Max: Ptr(500), Default: 4.0},
			{Name: "compact", Type: "boolean", Desc: "Return actionable findings without repeated geometry evidence"},
			{Name: "detailed", Type: "boolean", Desc: "Include full geometry and measurement evidence"},
		},
	})

	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	for _, s := range []Schema{
		{Command: "layout.set_grid_container", Description: "Set Figma grid auto-layout container properties", Params: []Param{nodeID, {Name: "gridRowCount", Type: "number", Min: Ptr(1)}, {Name: "gridColumnCount", Type: "number", Min: Ptr(1)}, {Name: "gridRowGap", Type: "number", Min: Ptr(0)}, {Name: "gridColumnGap", Type: "number", Min: Ptr(0)}, {Name: "gridRowsSizing", Type: "array"}, {Name: "gridColumnsSizing", Type: "array"}}},
		{Command: "layout.set_grid_tracks", Description: "Set Figma grid auto-layout row and column track sizing", Params: []Param{nodeID, {Name: "gridRowsSizing", Type: "array"}, {Name: "gridColumnsSizing", Type: "array"}}},
		{Command: "layout.set_grid_child_position", Description: "Set a direct child position and span inside a Figma grid auto-layout container", Params: []Param{nodeID, {Name: "gridRowAnchorIndex", Type: "number", Min: Ptr(0)}, {Name: "gridColumnAnchorIndex", Type: "number", Min: Ptr(0)}, {Name: "gridRowSpan", Type: "number", Min: Ptr(1)}, {Name: "gridColumnSpan", Type: "number", Min: Ptr(1)}, {Name: "gridChildHorizontalAlign", Type: "string"}, {Name: "gridChildVerticalAlign", Type: "string"}}},
		{Command: "layout.get_grid_layout", Description: "Get Figma grid auto-layout container properties", Params: []Param{nodeID}},
		{Command: "layout.reorder_grid_rows", Description: "Reorder Figma grid auto-layout row tracks where supported", Params: []Param{nodeID, {Name: "fromIndex", Type: "number", Required: true, Min: Ptr(0)}, {Name: "toIndex", Type: "number", Required: true, Min: Ptr(0)}}},
		{Command: "layout.reorder_grid_columns", Description: "Reorder Figma grid auto-layout column tracks where supported", Params: []Param{nodeID, {Name: "fromIndex", Type: "number", Required: true, Min: Ptr(0)}, {Name: "toIndex", Type: "number", Required: true, Min: Ptr(0)}}},
	} {
		Register(s)
	}
}
