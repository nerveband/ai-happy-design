package schema

func init() {
	Register(Schema{
		Command:     "boolean.union",
		Description: "Boolean union of multiple nodes (merge shapes together)",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Desc: "Comma-separated node IDs (minimum 2)"},
			{Name: "name", Type: "string", Desc: "Name for the resulting boolean node"},
		},
	})

	Register(Schema{
		Command:     "boolean.subtract",
		Description: "Boolean subtract (cut second+ shapes from first shape)",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Desc: "Comma-separated node IDs (minimum 2, first is base)"},
			{Name: "name", Type: "string", Desc: "Name for the resulting boolean node"},
		},
	})

	Register(Schema{
		Command:     "boolean.intersect",
		Description: "Boolean intersect (keep only overlapping areas)",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Desc: "Comma-separated node IDs (minimum 2)"},
			{Name: "name", Type: "string", Desc: "Name for the resulting boolean node"},
		},
	})

	Register(Schema{
		Command:     "boolean.exclude",
		Description: "Boolean exclude (keep non-overlapping areas, XOR)",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Desc: "Comma-separated node IDs (minimum 2)"},
			{Name: "name", Type: "string", Desc: "Name for the resulting boolean node"},
		},
	})

	Register(Schema{
		Command:     "boolean.flatten",
		Description: "Flatten node(s) into a single vector",
		Params: []Param{
			{Name: "nodeId", Type: "string", Aliases: []string{"nodeIds"}, Required: true, Desc: "Node ID or comma-separated node IDs to flatten"},
		},
	})
}
