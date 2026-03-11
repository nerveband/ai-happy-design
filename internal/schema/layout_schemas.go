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
			{Name: "horizontal", Type: "string", Aliases: []string{"layoutSizingHorizontal"}, Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "vertical", Type: "string", Aliases: []string{"layoutSizingVertical"}, Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "primaryAxis", Type: "string", Aliases: []string{"primaryAxisSizing"}, Desc: "Primary axis sizing mode (frame-level)", Enum: []string{"FIXED", "AUTO"}},
			{Name: "counterAxis", Type: "string", Aliases: []string{"counterAxisSizing"}, Desc: "Counter axis sizing mode (frame-level)", Enum: []string{"FIXED", "AUTO"}},
			{Name: "width", Type: "number", Desc: "Explicit width (for FIXED mode)", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Desc: "Explicit height (for FIXED mode)", Min: Ptr(1), Max: Ptr(10000)},
		},
	})

	Register(Schema{
		Command:     "layout.set_padding",
		Description: "Set padding on a frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "padding", Type: "number", Aliases: []string{"all"}, Desc: "Uniform padding for all sides", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingTop", Type: "number", Aliases: []string{"top"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingRight", Type: "number", Aliases: []string{"right"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingBottom", Type: "number", Aliases: []string{"bottom"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingLeft", Type: "number", Aliases: []string{"left"}, Min: Ptr(0), Max: Ptr(500)},
		},
	})

	Register(Schema{
		Command:     "layout.set_spacing",
		Description: "Set item spacing in an auto-layout frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "itemSpacing", Type: "number", Aliases: []string{"spacing"}, Desc: "Space between children", Min: Ptr(0), Max: Ptr(500)},
		},
	})

	Register(Schema{
		Command:     "layout.set_alignment",
		Description: "Set axis alignment in an auto-layout frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "primaryAxisAlignItems", Type: "string", Aliases: []string{"primaryAxisAlign", "primary"}, Enum: []string{"MIN", "CENTER", "MAX", "SPACE_BETWEEN"}},
			{Name: "counterAxisAlignItems", Type: "string", Aliases: []string{"counterAxisAlign", "counter"}, Enum: []string{"MIN", "CENTER", "MAX", "BASELINE"}},
		},
	})

	Register(Schema{
		Command:     "layout.set_wrap",
		Description: "Set layout wrap mode on an auto-layout frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "layoutWrap", Type: "string", Aliases: []string{"wrap"}, Enum: []string{"WRAP", "NO_WRAP"}},
		},
	})

	Register(Schema{
		Command:     "layout.set_constraints",
		Description: "Set layout constraints on a node (for non-auto-layout parents)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "horizontal", Type: "string", Aliases: []string{"constraintHorizontal"}, Desc: "Horizontal constraint", Enum: []string{"MIN", "CENTER", "MAX", "STRETCH", "SCALE"}},
			{Name: "vertical", Type: "string", Aliases: []string{"constraintVertical"}, Desc: "Vertical constraint", Enum: []string{"MIN", "CENTER", "MAX", "STRETCH", "SCALE"}},
		},
	})

	Register(Schema{
		Command:     "layout.check_overlaps",
		Description: "Check for overlapping children in a frame. Returns overlap details.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Frame node ID to check", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layout.set_grid",
		Aliases:     []string{"layout.add_grid"},
		Description: "Set or add layout grids on a frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "grids", Type: "array", Desc: "Array of grid objects [{pattern, count, gutterSize, alignment, sectionSize, offset, color, visible}]"},
			{Name: "pattern", Type: "string", Aliases: []string{"type"}, Desc: "Grid pattern (for single grid shorthand)", Enum: []string{"COLUMNS", "ROWS", "GRID"}},
			{Name: "count", Type: "number", Desc: "Number of columns/rows", Default: 12.0},
			{Name: "gutterSize", Type: "number", Aliases: []string{"gutter"}, Desc: "Gutter size between columns/rows", Default: 20.0},
			{Name: "alignment", Type: "string", Desc: "Grid alignment", Enum: []string{"MIN", "CENTER", "MAX", "STRETCH"}, Default: "STRETCH"},
			{Name: "sectionSize", Type: "number", Aliases: []string{"size"}, Desc: "Section size (column/row width/height)"},
			{Name: "offset", Type: "number", Desc: "Grid offset"},
			{Name: "append", Type: "boolean", Desc: "Append to existing grids instead of replacing"},
		},
	})

	Register(Schema{
		Command:     "layout.get_grids",
		Description: "Get layout grids on a frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layout.remove_grids",
		Aliases:     []string{"layout.clear_grids"},
		Description: "Remove layout grids from a frame. If index specified, removes only that grid.",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "index", Type: "number", Desc: "Index of specific grid to remove (omit to remove all)"},
		},
	})
}
