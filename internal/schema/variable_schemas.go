package schema

func init() {
	for _, s := range []Schema{
		{Command: "variable.extend_collection", Description: "Enterprise/Beta Figma API: extend a local variable collection", Params: []Param{{Name: "collectionId", Type: "string", Required: true}, {Name: "name", Type: "string"}}},
		{Command: "variable.extend_library_collection", Description: "Enterprise/Beta Figma API: extend a published library variable collection by key", Params: []Param{{Name: "collectionKey", Type: "string", Required: true}, {Name: "name", Type: "string"}}},
		{Command: "variable.get_values_for_collection", Description: "Beta Figma API: get inherited and overridden values for a variable in a collection", Params: []Param{{Name: "variableId", Type: "string", Required: true}, {Name: "collectionId", Type: "string", Required: true}}},
		{Command: "variable.remove_mode_override", Description: "Beta Figma API: remove one extended mode override from a variable", Params: []Param{{Name: "variableId", Type: "string", Required: true}, {Name: "modeId", Type: "string", Required: true}}},
		{Command: "variable.remove_collection_overrides", Description: "Beta Figma API: remove all overrides for one variable in an extended collection", Params: []Param{{Name: "collectionId", Type: "string", Required: true}, {Name: "variableId", Type: "string", Required: true}}},
		{Command: "variable.get_overrides", Description: "Beta Figma API: get override map for an extended variable collection", Params: []Param{{Name: "collectionId", Type: "string", Required: true}}},
	} {
		Register(s)
	}
}
