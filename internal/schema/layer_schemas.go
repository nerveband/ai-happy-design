package schema

func init() {
	Register(Schema{
		Command:     "layer.set_order",
		Description: "Set the layer order (z-index) of a node within its parent",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "index", Type: "number", Required: true, Desc: "Target z-index position", Min: Ptr(0)},
		},
	})

	Register(Schema{
		Command:     "layer.bring_forward",
		Description: "Move a node one layer forward in its parent",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layer.send_backward",
		Description: "Move a node one layer backward in its parent",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layer.bring_to_front",
		Description: "Move a node to the front (topmost layer) in its parent",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layer.send_to_back",
		Description: "Move a node to the back (bottommost layer) in its parent",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layer.group",
		Description: "Group multiple nodes together. All nodes must share the same parent.",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Desc: "Comma-separated node IDs to group"},
			{Name: "name", Type: "string", Desc: "Name for the new group"},
		},
	})

	Register(Schema{
		Command:     "layer.ungroup",
		Description: "Ungroup a group node, releasing its children back to the parent",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Group node ID to ungroup", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "layer.insert_child",
		Description: "Insert a child node into a parent at a specific index",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Desc: "Parent node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "childId", Type: "string", Required: true, Desc: "Child node ID to insert", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "index", Type: "number", Desc: "Position to insert at (default 0)", Min: Ptr(0), Default: 0.0},
		},
	})
}
