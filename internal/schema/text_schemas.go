package schema

func init() {
	Register(Schema{
		Command:     "text.create",
		Aliases:     []string{"text"},
		Description: "Create a text node",
		Params: []Param{
			{Name: "text", Type: "string", Required: true, Aliases: []string{"content"}, Desc: "Text content to display"},
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string", Desc: "Semantic name for the text node"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "fontSize", Type: "number", Aliases: []string{"sz"}, Desc: "Font size in pixels", Min: Ptr(4), Max: Ptr(500), Default: 16.0, SemanticTokens: true},
			{Name: "fontFamily", Type: "string", Aliases: []string{"ff"}, Desc: "Font family name", Default: "Inter"},
			{Name: "fontStyle", Type: "string", Aliases: []string{"fs"}, Desc: "Font style", Default: "Regular",
				Enum: []string{"Thin", "ExtraLight", "Light", "Regular", "Medium", "SemiBold", "Bold", "ExtraBold", "Black", "Italic", "Bold Italic"}},
			{Name: "color", Type: "string", Desc: "Text color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`, Default: "#000000"},
			{Name: "lineHeight", Type: "number", Aliases: []string{"lh"}, Desc: "Line height percentage (e.g. 150 = 150%%)", Min: Ptr(50), Max: Ptr(300), AutoFix: "lineHeightUnit:PERCENT"},
			{Name: "lineHeightUnit", Type: "string", Enum: []string{"PIXELS", "PERCENT", "AUTO"}},
			{Name: "letterSpacing", Type: "number", Aliases: []string{"ls"}, Min: Ptr(-5), Max: Ptr(20)},
			{Name: "width", Type: "number", Desc: "Text box width (enables wrapping)", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "textAlign", Type: "string", Aliases: []string{"textAlignHorizontal"}, Enum: []string{"LEFT", "CENTER", "RIGHT", "JUSTIFIED"}},
			{Name: "textCase", Type: "string", Enum: []string{"ORIGINAL", "UPPER", "LOWER", "TITLE"}},
			{Name: "textDecoration", Type: "string", Enum: []string{"NONE", "UNDERLINE", "STRIKETHROUGH"}},
			{Name: "maxLines", Type: "number", Desc: "Maximum lines before truncation with ellipsis", Min: Ptr(1), Max: Ptr(100)},
			{Name: "layoutSizingHorizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutSizingVertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
		},
	})

	Register(Schema{
		Command:     "text.set_content",
		Description: "Set text content of an existing text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "content", Type: "string", Required: true, Aliases: []string{"text"}},
		},
	})
}
