package schema

func init() {
	Register(Schema{
		Command:     "document.get_info",
		Description: "Get document information including all pages",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "document.get_selection",
		Description: "Get the currently selected nodes on the active page",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "document.set_selection",
		Description: "Set the selection to specific nodes",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Desc: "Comma-separated node IDs to select"},
		},
	})

	Register(Schema{
		Command:     "document.scan_text",
		Description: "Scan all text nodes on the current page, optionally filtering by query",
		Params: []Param{
			{Name: "pageId", Type: "string", Desc: "Page ID to scan (defaults to current page)"},
			{Name: "query", Type: "string", Desc: "Text content filter (case-insensitive substring match)"},
		},
	})

	Register(Schema{
		Command:     "document.scan_by_type",
		Description: "Find all nodes of a specific type on the current page",
		Params: []Param{
			{Name: "nodeType", Type: "string", Required: true, Aliases: []string{"type"}, Desc: "Figma node type to scan for", Enum: []string{"FRAME", "TEXT", "RECTANGLE", "ELLIPSE", "LINE", "POLYGON", "STAR", "VECTOR", "GROUP", "COMPONENT", "COMPONENT_SET", "INSTANCE", "BOOLEAN_OPERATION", "SECTION"}},
		},
	})

	Register(Schema{
		Command:     "document.get_styles",
		Description: "Get all document styles (paint, text, effect, grid)",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "document.focus",
		Aliases:     []string{"document.zoom_to"},
		Description: "Focus the viewport on specific nodes (scroll and zoom to fit)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Aliases: []string{"nodeIds"}, Desc: "Node ID or comma-separated node IDs to focus on"},
		},
	})

	Register(Schema{
		Command:     "document.find_free_space",
		Description: "Find available free space on the current page for placing new frames. Returns suggested x,y and lists existing frames.",
		Params: []Param{
			{Name: "width", Type: "number", Desc: "Requested width for the new frame", Default: 1080.0},
			{Name: "height", Type: "number", Desc: "Requested height for the new frame", Default: 1080.0},
			{Name: "gap", Type: "number", Desc: "Gap between frames", Default: 100.0},
		},
	})

	Register(Schema{
		Command:     "document.find_nodes",
		Description: "Unified search: find nodes by name, type, and/or text content",
		Params: []Param{
			{Name: "query", Type: "string", Desc: "Name filter (case-insensitive substring match)"},
			{Name: "type", Type: "string", Desc: "Node type filter", Enum: []string{"FRAME", "TEXT", "RECTANGLE", "ELLIPSE", "LINE", "POLYGON", "STAR", "VECTOR", "GROUP", "COMPONENT", "COMPONENT_SET", "INSTANCE", "BOOLEAN_OPERATION", "SECTION"}},
			{Name: "textContent", Type: "string", Desc: "Text content filter (only matches TEXT nodes)"},
			{Name: "pageId", Type: "string", Desc: "Page ID to search in (defaults to current page)"},
		},
	})

	Register(Schema{
		Command:     "document.lint",
		Aliases:     []string{"document.validate", "document.check"},
		Description: "Lint a node tree for common issues (overlaps, overflow, default names, text sizing)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Root node ID to lint", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})
}
