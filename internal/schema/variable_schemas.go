package schema

func init() {
	Register(Schema{
		Command:     "variable.create",
		Description: "Create a variable in a collection",
		Params: []Param{
			{Name: "name", Type: "string", Required: true, Desc: "Variable name"},
			{Name: "collectionId", Type: "string", Desc: "Collection ID (uses first collection or creates default if omitted)"},
			{Name: "collectionName", Type: "string", Desc: "Collection name (finds existing or creates new)"},
			{Name: "resolvedType", Type: "string", Aliases: []string{"type"}, Desc: "Variable type", Enum: []string{"COLOR", "FLOAT", "STRING", "BOOLEAN"}, Default: "COLOR"},
			{Name: "value", Type: "string", Desc: "Initial value (JSON for COLOR objects, plain string/number/boolean otherwise)"},
		},
	})

	Register(Schema{
		Command:     "variable.get_all",
		Aliases:     []string{"variable.list"},
		Description: "Get all local variables with their values by mode",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "variable.set_value",
		Description: "Set a variable's value for a specific mode",
		Params: []Param{
			{Name: "variableId", Type: "string", Required: true, Desc: "Variable ID"},
			{Name: "value", Type: "string", Required: true, Desc: "Value to set (JSON for COLOR objects)"},
			{Name: "modeId", Type: "string", Desc: "Mode ID (uses first available mode if omitted)"},
		},
	})

	Register(Schema{
		Command:     "variable.bind",
		Description: "Bind a variable to a node property",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "variableId", Type: "string", Required: true, Desc: "Variable ID to bind"},
			{Name: "field", Type: "string", Aliases: []string{"property"}, Desc: "Property to bind to", Default: "fills"},
		},
	})

	Register(Schema{
		Command:     "variable.unbind",
		Description: "Unbind a variable from a node property",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "field", Type: "string", Aliases: []string{"property"}, Desc: "Property to unbind", Default: "fills"},
		},
	})

	Register(Schema{
		Command:     "variable.create_collection",
		Description: "Create a variable collection with optional modes",
		Params: []Param{
			{Name: "name", Type: "string", Required: true, Desc: "Collection name"},
			{Name: "modes", Type: "string", Desc: "Comma-separated mode names (e.g. 'Light,Dark')"},
		},
	})

	Register(Schema{
		Command:     "variable.get_collections",
		Aliases:     []string{"variable.list_collections"},
		Description: "Get all variable collections with their modes",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "variable.delete",
		Description: "Delete a variable",
		Params: []Param{
			{Name: "variableId", Type: "string", Required: true, Desc: "Variable ID to delete"},
		},
	})

	Register(Schema{
		Command:     "variable.resolve_for_consumer",
		Aliases:     []string{"variable.resolve"},
		Description: "Resolve a variable's value, following aliases recursively",
		Params: []Param{
			{Name: "variableId", Type: "string", Required: true, Desc: "Variable ID to resolve"},
			{Name: "nodeId", Type: "string", Desc: "Optional node ID for context-aware resolution", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "modeId", Type: "string", Desc: "Specific mode ID to resolve for"},
		},
	})

	Register(Schema{
		Command:     "variable.add_mode",
		Aliases:     []string{"variable.create_mode"},
		Description: "Add a new mode to a variable collection",
		Params: []Param{
			{Name: "collectionId", Type: "string", Required: true, Desc: "Collection ID"},
			{Name: "name", Type: "string", Required: true, Aliases: []string{"modeName"}, Desc: "Name for the new mode"},
		},
	})

	Register(Schema{
		Command:     "variable.rename_mode",
		Description: "Rename an existing mode in a variable collection",
		Params: []Param{
			{Name: "collectionId", Type: "string", Required: true, Desc: "Collection ID"},
			{Name: "modeId", Type: "string", Required: true, Desc: "Mode ID to rename"},
			{Name: "name", Type: "string", Required: true, Aliases: []string{"newName"}, Desc: "New mode name"},
		},
	})

	Register(Schema{
		Command:     "variable.delete_mode",
		Aliases:     []string{"variable.remove_mode"},
		Description: "Delete a mode from a variable collection",
		Params: []Param{
			{Name: "collectionId", Type: "string", Required: true, Desc: "Collection ID"},
			{Name: "modeId", Type: "string", Required: true, Desc: "Mode ID to delete"},
		},
	})
}
