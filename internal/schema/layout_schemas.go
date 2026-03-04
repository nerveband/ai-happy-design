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
}
