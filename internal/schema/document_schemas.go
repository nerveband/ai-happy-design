package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Pattern: `^[0-9]+:[0-9]+$`}
	Register(Schema{Command: "document.get_info", Description: "Get document information", Params: nil})
	Register(Schema{Command: "document.get_selection", Description: "Get current selection", Params: nil})
	Register(Schema{Command: "document.get_focused_node", Aliases: []string{"devmode.get_focused_node", "devmode.get_context", "devmode.get_selection_context"}, Description: "Get the current Dev Mode focused node", Params: nil})
	Register(Schema{Command: "document.set_selection", Description: "Set current selection", Params: []Param{{Name: "nodeIds", Type: "array", Required: true}}})
	Register(Schema{Command: "document.scan_text", Description: "Scan text nodes", Params: []Param{{Name: "query", Type: "string"}, {Name: "pageId", Type: "string"}}})
	Register(Schema{Command: "document.scan_by_type", Description: "Scan nodes by type", Params: []Param{{Name: "nodeType", Type: "string", Required: true}, {Name: "pageId", Type: "string"}}})
	Register(Schema{Command: "document.get_styles", Description: "Get document styles", Params: nil})
	Register(Schema{Command: "document.find_by_name", Description: "Find nodes by name", Params: []Param{{Name: "query", Type: "string", Required: true}, {Name: "pageId", Type: "string"}}})
	Register(Schema{Command: "document.find_by_type", Description: "Find nodes by type", Params: []Param{{Name: "nodeType", Type: "string", Required: true}, {Name: "pageId", Type: "string"}}})
	Register(Schema{Command: "document.zoom_to", Description: "Zoom viewport to a node", Params: []Param{nodeID}})
	Register(Schema{Command: "document.find_free_space", Description: "Find free canvas space for a new frame", Params: []Param{{Name: "width", Type: "number"}, {Name: "height", Type: "number"}}})
	Register(Schema{Command: "document.find_nodes", Description: "Find nodes by query, type, or text content", Params: []Param{{Name: "query", Type: "string"}, {Name: "nodeType", Type: "string"}, {Name: "textContent", Type: "string"}}})
	Register(Schema{Command: "document.lint", Description: "Lint a node for design quality", Params: []Param{nodeID}})
	Register(Schema{
		Command:       "document.accessibility_audit",
		Description:   "Run a deterministic accessibility audit over a batch payload or live document",
		Safety:        "read",
		Idempotency:   "idempotent",
		RequiresFigma: false,
		Params:        []Param{{Name: "file", Type: "string", Desc: "Batch JSON file for local audit"}, {Name: "batchFile", Type: "string", Desc: "Alias for file"}},
	})
}
