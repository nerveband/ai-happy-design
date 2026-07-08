package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	pageID := Param{Name: "pageId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	nodeIDs := Param{Name: "nodeIds", Type: "array", Required: true}
	color := Param{Name: "color", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`}

	for _, s := range []Schema{
		{Command: "page.create", Description: "Create a page", Params: []Param{{Name: "name", Type: "string"}}},
		{Command: "page.delete", Description: "Delete a page", Params: []Param{pageID}},
		{Command: "page.rename", Description: "Rename a page", Params: []Param{pageID, {Name: "name", Type: "string", Required: true}}},
		{Command: "page.set_current", Description: "Set current page", Params: []Param{pageID}},
		{Command: "page.get_all", Description: "Get all pages", Params: nil},
		{Command: "page.get_current", Description: "Get current page", Params: nil},
		{Command: "page.duplicate", Description: "Duplicate a page", Params: []Param{pageID, {Name: "name", Type: "string"}}},

		{Command: "export.image", Description: "Export node as PNG/JPG", Params: []Param{nodeID, {Name: "format", Type: "string", Enum: []string{"PNG", "JPG"}}, {Name: "scale", Type: "number", Min: Ptr(0.1), Max: Ptr(8)}}},
		{Command: "export.svg", Description: "Export node as SVG", Params: []Param{nodeID}},
		{Command: "export.pdf", Description: "Export node as PDF", Params: []Param{nodeID}},
		{Command: "export.json", Description: "Export node as structured JSON", Params: []Param{nodeID, {Name: "depth", Type: "number", Min: Ptr(1), Max: Ptr(50)}}},
		{Command: "export.batch_export", Description: "Export multiple nodes", Params: []Param{nodeIDs, {Name: "format", Type: "string", Enum: []string{"PNG", "JPG", "SVG", "PDF", "JSON"}}, {Name: "scale", Type: "number"}}},

		{Command: "boolean.union", Description: "Boolean union", Params: []Param{nodeIDs, {Name: "name", Type: "string"}}},
		{Command: "boolean.subtract", Description: "Boolean subtract", Params: []Param{nodeIDs, {Name: "name", Type: "string"}}},
		{Command: "boolean.intersect", Description: "Boolean intersect", Params: []Param{nodeIDs, {Name: "name", Type: "string"}}},
		{Command: "boolean.exclude", Description: "Boolean exclude", Params: []Param{nodeIDs, {Name: "name", Type: "string"}}},
		{Command: "boolean.flatten", Description: "Flatten nodes", Params: []Param{nodeIDs}},

		{Command: "style.create_paint", Description: "Create a paint style", Params: []Param{{Name: "name", Type: "string", Required: true}, color, {Name: "description", Type: "string"}}},
		{Command: "style.create_text", Description: "Create a text style", Params: []Param{{Name: "name", Type: "string", Required: true}, {Name: "description", Type: "string"}}},
		{Command: "style.create_effect", Description: "Create an effect style", Params: []Param{{Name: "name", Type: "string", Required: true}, {Name: "description", Type: "string"}}},
		{Command: "style.apply", Description: "Apply a style", Params: []Param{nodeID, {Name: "styleId", Type: "string", Required: true}, {Name: "styleType", Type: "string"}}},
		{Command: "style.get_all", Description: "Get all local styles", Params: nil},
		{Command: "style.remove", Description: "Remove a style from a node", Params: []Param{nodeID, {Name: "styleType", Type: "string", Required: true}}},

		{Command: "variable.create", Description: "Create a variable", Params: []Param{{Name: "name", Type: "string", Required: true}, {Name: "collectionId", Type: "string"}, {Name: "collectionName", Type: "string"}, {Name: "resolvedType", Type: "string"}}},
		{Command: "variable.get_all", Description: "Get all variables", Params: nil},
		{Command: "variable.set_value", Description: "Set variable value", Params: []Param{{Name: "variableId", Type: "string", Required: true}, {Name: "value", Type: "object"}, {Name: "modeId", Type: "string"}}},
		{Command: "variable.bind", Description: "Bind variable to node", Params: []Param{nodeID, {Name: "variableId", Type: "string", Required: true}, {Name: "field", Type: "string"}}},
		{Command: "variable.unbind", Description: "Unbind variable from node", Params: []Param{nodeID, {Name: "field", Type: "string"}}},
		{Command: "variable.create_collection", Description: "Create a variable collection", Params: []Param{{Name: "name", Type: "string", Required: true}}},
		{Command: "variable.get_collections", Description: "Get variable collections", Params: nil},
		{Command: "variable.delete", Description: "Delete a variable", Params: []Param{{Name: "variableId", Type: "string", Required: true}}},
		{Command: "variable.resolve_for_consumer", Description: "Resolve a variable for a consuming node", Params: []Param{{Name: "variableId", Type: "string", Required: true}, nodeID}},
		{Command: "variable.add_mode", Description: "Add a variable collection mode", Params: []Param{{Name: "collectionId", Type: "string", Required: true}, {Name: "name", Type: "string", Required: true}}},
		{Command: "variable.rename_mode", Description: "Rename a variable collection mode", Params: []Param{{Name: "collectionId", Type: "string", Required: true}, {Name: "modeId", Type: "string", Required: true}, {Name: "name", Type: "string", Required: true}}},
		{Command: "variable.delete_mode", Description: "Delete a variable collection mode", Params: []Param{{Name: "collectionId", Type: "string", Required: true}, {Name: "modeId", Type: "string", Required: true}}},

		{Command: "component.create", Description: "Create a component", Params: []Param{nodeID, {Name: "name", Type: "string"}}},
		{Command: "component.create_instance", Description: "Create a component instance", Params: []Param{{Name: "componentId", Type: "string"}, {Name: "componentKey", Type: "string"}, {Name: "x", Type: "number"}, {Name: "y", Type: "number"}}},
		{Command: "component.create_set", Description: "Create a component set", Params: []Param{nodeIDs, {Name: "name", Type: "string"}}},
		{Command: "component.get_local", Description: "Get local components", Params: nil},
		{Command: "component.get_remote", Description: "Get remote components", Params: nil},
		{Command: "component.get_overrides", Description: "Get instance overrides", Params: []Param{nodeID}},
		{Command: "component.set_overrides", Description: "Set instance overrides", Params: []Param{nodeID, {Name: "overrides", Type: "object", Required: true}}},
		{Command: "component.detach_instance", Description: "Detach an instance", Params: []Param{nodeID}},
		{Command: "component.reset_instance", Description: "Reset an instance", Params: []Param{nodeID}},
		{Command: "component.swap_instance", Description: "Swap an instance component", Params: []Param{nodeID, {Name: "newComponentId", Type: "string", Required: true}}},
		{Command: "component.get_property_definitions", Description: "Get component property definitions", Params: []Param{nodeID}},
		{Command: "component.add_property_definition", Description: "Add component property definition", Params: []Param{nodeID, {Name: "name", Type: "string", Required: true}, {Name: "type", Type: "string", Required: true}}},
		{Command: "component.delete_property_definition", Description: "Delete component property definition", Params: []Param{nodeID, {Name: "propertyName", Type: "string", Required: true}}},
		{Command: "component.create_slot", Description: "Create a component slot where supported by Figma", Params: []Param{nodeID, {Name: "name", Type: "string", Required: true}, {Name: "x", Type: "number"}, {Name: "y", Type: "number"}, {Name: "width", Type: "number"}, {Name: "height", Type: "number"}}},
		{Command: "component.reset_slot", Description: "Reset a component slot where supported by Figma", Params: []Param{nodeID}},
		{Command: "component.get_slots", Description: "Get component slots where supported by Figma", Params: []Param{nodeID}, Safety: "read"},
	} {
		Register(s)
	}
}
