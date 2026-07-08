package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Pattern: `^[0-9]+:[0-9]+$`}
	for _, s := range []Schema{
		{Command: "tokens.preset_tailwind", Description: "Return or write Tailwind-compatible design token preset", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "outputPath", Type: "string"}}},
		{Command: "tokens.preset_shadcn", Description: "Return or write shadcn-compatible design token preset", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "outputPath", Type: "string"}}},
		{Command: "tokens.preset_material", Description: "Return or write Material-compatible design token preset", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "outputPath", Type: "string"}}},
		{Command: "tokens.setup_system", Description: "Generate a batch payload for setting up token variables/styles", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "preset", Type: "string"}, {Name: "outputPath", Type: "string"}}},
		{Command: "design_system.health", Description: "Score design-system health from a live or supplied design-system snapshot", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "spec", Type: "object"}, {Name: "specPath", Type: "string"}}},
		{Command: "component.analyze_set", Description: "Analyze component set variants and likely UI states", Params: []Param{nodeID, {Name: "componentSet", Type: "object"}}, Safety: "local", Idempotency: "idempotent"},
		{Command: "component.arrange_set", Description: "Generate arrangement plan for component variants", Params: []Param{nodeID, {Name: "componentSet", Type: "object"}}, Safety: "local", Idempotency: "idempotent"},
		{Command: "parity.audit_component", Description: "Audit a Figma component snapshot against local component code/spec", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "figmaNode", Type: "object"}, {Name: "figmaNodePath", Type: "string"}, {Name: "codeSpec", Type: "object"}, {Name: "codeSpecPath", Type: "string"}}},
		{Command: "document.get_editor_context", Description: "Get current Figma editor type and runtime feature flags", Params: nil, Safety: "read"},
		{Command: "slides.get_current", Description: "Get current Figma Slides context where supported", Params: nil, Safety: "read"},
		{Command: "slides.create_slide", Description: "Create a slide where supported by Figma Slides", Params: []Param{{Name: "row", Type: "number"}, {Name: "col", Type: "number"}, {Name: "name", Type: "string"}}},
		{Command: "slides.set_background", Description: "Set slide background where supported", Params: []Param{nodeID, {Name: "color", Type: "string"}}},
		{Command: "slides.add_text", Description: "Add text to a slide where supported", Params: []Param{nodeID, {Name: "text", Type: "string", Required: true}, {Name: "x", Type: "number"}, {Name: "y", Type: "number"}}},
		{Command: "slides.reorder", Description: "Reorder slides where supported", Params: []Param{nodeID, {Name: "row", Type: "number"}, {Name: "col", Type: "number"}}},
		{Command: "figjam.create_sticky", Description: "Create a FigJam sticky where supported", Params: []Param{{Name: "text", Type: "string"}, {Name: "x", Type: "number"}, {Name: "y", Type: "number"}, {Name: "color", Type: "string"}}},
		{Command: "figjam.create_shape", Description: "Create a FigJam shape where supported", Params: []Param{{Name: "shape", Type: "string"}, {Name: "text", Type: "string"}, {Name: "x", Type: "number"}, {Name: "y", Type: "number"}}},
		{Command: "figjam.create_connector", Description: "Create a FigJam connector where supported", Params: []Param{{Name: "startNodeId", Type: "string"}, {Name: "endNodeId", Type: "string"}, {Name: "text", Type: "string"}}},
		{Command: "figjam.get_board", Description: "Get FigJam board summary where supported", Params: nil, Safety: "read"},
		{Command: "packaging.generate_skills", Description: "Generate portable AI tool skill/rules files from the internal catalog", Safety: "local", Idempotency: "idempotent", Params: []Param{{Name: "outputDir", Type: "string"}, {Name: "targets", Type: "array"}}},
	} {
		Register(s)
	}
}
