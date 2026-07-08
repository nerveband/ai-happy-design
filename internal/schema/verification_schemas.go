package schema

func init() {
	for _, s := range []Schema{
		{Command: "document.screenshot", Description: "Capture the current page or a target node as a screenshot artifact", Params: []Param{{Name: "nodeId", Type: "string"}, {Name: "pageId", Type: "string"}, {Name: "scale", Type: "number", Min: Ptr(0.01)}, {Name: "format", Type: "string", Enum: []string{"PNG", "JPG"}}, {Name: "outputDir", Type: "string"}}, Safety: "read"},
		{Command: "document.screenshot_selection", Description: "Capture the current selection as screenshot artifacts", Params: []Param{{Name: "scale", Type: "number", Min: Ptr(0.01)}, {Name: "format", Type: "string", Enum: []string{"PNG", "JPG"}}, {Name: "outputDir", Type: "string"}}, Safety: "read"},
		{Command: "verify.visual", Description: "Describe a visual verification artifact and expected inspection loop", Params: []Param{{Name: "artifactPath", Type: "string", Required: true}, {Name: "target", Type: "string"}, {Name: "notes", Type: "string"}}, Safety: "local"},
	} {
		Register(s)
	}
}
