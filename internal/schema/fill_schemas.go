package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	Register(Schema{Command: "fill.apply_shader", Description: "Apply a shader paint fill to a node where supported by Figma", Params: []Param{nodeID, {Name: "shaderId", Type: "string", Required: true}, {Name: "uniforms", Type: "object"}}})
}
