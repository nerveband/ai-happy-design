package schema

func init() {
	Register(Schema{
		Command:     "component.create",
		Description: "Create a component from an existing node, or create an empty component",
		Params: []Param{
			{Name: "nodeId", Type: "string", Desc: "Node ID to convert to component (optional — omit for empty component)", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string", Desc: "Component name"},
			{Name: "x", Type: "number", Desc: "X position (for empty component)"},
			{Name: "y", Type: "number", Desc: "Y position (for empty component)"},
			{Name: "width", Type: "number", Desc: "Width (for empty component)", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Desc: "Height (for empty component)", Min: Ptr(1), Max: Ptr(10000)},
		},
	})

	Register(Schema{
		Command:     "component.create_instance",
		Description: "Create an instance of a component",
		Params: []Param{
			{Name: "componentId", Type: "string", Desc: "Local component node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "componentKey", Type: "string", Desc: "Component key (for remote/published components)"},
			{Name: "x", Type: "number", Desc: "X position"},
			{Name: "y", Type: "number", Desc: "Y position"},
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Desc: "Parent node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string", Desc: "Instance name override"},
		},
	})

	Register(Schema{
		Command:     "component.create_set",
		Description: "Create a component set (variants) from multiple components",
		Params: []Param{
			{Name: "nodeIds", Type: "string", Required: true, Aliases: []string{"componentIds"}, Desc: "Comma-separated component node IDs"},
			{Name: "name", Type: "string", Desc: "Component set name"},
		},
	})

	Register(Schema{
		Command:     "component.get_local",
		Description: "List all local components on the current page",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "component.get_remote",
		Description: "List remote (published) components used in the current page",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "component.get_overrides",
		Description: "Get the overrides of a component instance",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Instance node ID", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "component.set_overrides",
		Description: "Set overrides on a component instance",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Instance node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "overrides", Type: "object", Required: true, Desc: "JSON object of property overrides to apply"},
		},
	})

	Register(Schema{
		Command:     "component.detach_instance",
		Description: "Detach a component instance, converting it to a regular frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Instance node ID to detach", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "component.reset_instance",
		Description: "Reset all overrides on a component instance back to main component defaults",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Instance node ID to reset", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "component.swap_instance",
		Description: "Swap a component instance to use a different main component",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Instance node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "newComponentId", Type: "string", Required: true, Desc: "New component node ID to swap to", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "component.get_property_definitions",
		Aliases:     []string{"component.get_properties"},
		Description: "Get component property definitions for a component or component set",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Component or component set node ID", Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "component.add_property_definition",
		Aliases:     []string{"component.add_property"},
		Description: "Add a property definition to a component or component set",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Component or component set node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "propertyName", Type: "string", Required: true, Aliases: []string{"name"}, Desc: "Property name"},
			{Name: "type", Type: "string", Desc: "Property type", Enum: []string{"TEXT", "BOOLEAN", "VARIANT", "INSTANCE_SWAP"}, Default: "TEXT"},
			{Name: "defaultValue", Type: "string", Desc: "Default value for the property"},
			{Name: "variantOptions", Type: "array", Desc: "Variant options array (for VARIANT type)"},
			{Name: "preferredValues", Type: "array", Desc: "Preferred values array (for INSTANCE_SWAP type)"},
		},
	})

	Register(Schema{
		Command:     "component.delete_property_definition",
		Aliases:     []string{"component.delete_property", "component.remove_property"},
		Description: "Delete a property definition from a component or component set",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Component or component set node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "propertyName", Type: "string", Required: true, Aliases: []string{"name"}, Desc: "Property name to delete"},
		},
	})
}
