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

	Register(Schema{
		Command:     "text.set_font",
		Description: "Set font family and style on a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "fontFamily", Type: "string", Aliases: []string{"family"}, Desc: "Font family name", Default: "Inter"},
			{Name: "fontStyle", Type: "string", Aliases: []string{"style"}, Desc: "Font style (e.g. Regular, Bold, Italic)", Default: "Regular"},
			{Name: "rangeStart", Type: "number", Desc: "Character range start (for partial styling)"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end (for partial styling)"},
		},
	})

	Register(Schema{
		Command:     "text.set_size",
		Description: "Set font size on a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "fontSize", Type: "number", Required: true, Aliases: []string{"size"}, Desc: "Font size in pixels", Min: Ptr(4), Max: Ptr(500)},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})

	Register(Schema{
		Command:     "text.set_weight",
		Description: "Set font weight on a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "fontWeight", Type: "string", Required: true, Aliases: []string{"weight"}, Desc: "Font weight (100-900 numeric or style name like Bold)"},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})

	Register(Schema{
		Command:     "text.set_color",
		Description: "Set text fill color",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Required: true, Desc: "Text color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})

	Register(Schema{
		Command:     "text.set_align",
		Description: "Set text alignment",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "textAlign", Type: "string", Aliases: []string{"horizontal", "textAlignHorizontal"}, Desc: "Horizontal alignment", Enum: []string{"LEFT", "CENTER", "RIGHT", "JUSTIFIED"}},
			{Name: "textAlignVertical", Type: "string", Aliases: []string{"vertical"}, Desc: "Vertical alignment", Enum: []string{"TOP", "CENTER", "BOTTOM"}},
		},
	})

	Register(Schema{
		Command:     "text.set_spacing",
		Description: "Set letter spacing, line height, and paragraph spacing on a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "letterSpacing", Type: "number", Desc: "Letter spacing value", Min: Ptr(-5), Max: Ptr(20)},
			{Name: "letterSpacingUnit", Type: "string", Enum: []string{"PIXELS", "PERCENT"}, Default: "PIXELS"},
			{Name: "lineHeight", Type: "number", Desc: "Line height value"},
			{Name: "lineHeightUnit", Type: "string", Enum: []string{"PIXELS", "PERCENT", "AUTO"}, Default: "PIXELS"},
			{Name: "paragraphSpacing", Type: "number", Desc: "Paragraph spacing in pixels", Min: Ptr(0)},
		},
	})

	Register(Schema{
		Command:     "text.set_case",
		Description: "Set text case transformation",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "textCase", Type: "string", Required: true, Aliases: []string{"case"}, Enum: []string{"ORIGINAL", "UPPER", "LOWER", "TITLE"}},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})

	Register(Schema{
		Command:     "text.set_decoration",
		Description: "Set text decoration (underline, strikethrough)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "textDecoration", Type: "string", Required: true, Aliases: []string{"decoration"}, Enum: []string{"NONE", "UNDERLINE", "STRIKETHROUGH"}},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})

	Register(Schema{
		Command:     "text.get_segments",
		Description: "Get styled text segments (character ranges with different styles)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "property", Type: "string", Desc: "Specific property to get segments for (omit for all)", Enum: []string{"fontName", "fontSize", "fills", "textDecoration", "textCase", "lineHeight", "letterSpacing"}},
		},
	})

	Register(Schema{
		Command:     "text.load_font",
		Description: "Load a font for use in text nodes",
		Params: []Param{
			{Name: "fontFamily", Type: "string", Aliases: []string{"family"}, Desc: "Font family name", Default: "Inter"},
			{Name: "fontStyle", Type: "string", Aliases: []string{"style"}, Desc: "Font style", Default: "Regular"},
		},
	})

	Register(Schema{
		Command:     "text.set_style_id",
		Description: "Apply a text style to a text node by style ID",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "styleId", Type: "string", Required: true, Desc: "Text style ID to apply (empty string to clear)"},
		},
	})

	Register(Schema{
		Command:     "text.list_fonts",
		Aliases:     []string{"text.available_fonts"},
		Description: "List available fonts, optionally filtered by family name",
		Params: []Param{
			{Name: "fontFamily", Type: "string", Aliases: []string{"family"}, Desc: "Filter substring for font family name"},
		},
	})

	Register(Schema{
		Command:     "text.set_range_style",
		Description: "Apply multiple styles to character ranges within a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "ranges", Type: "array", Required: true, Desc: "Array of range objects: [{match|start+end, bold, italic, color, fontSize, fontFamily, fontStyle, letterSpacing, lineHeight, textDecoration, textCase, hyperlink}]"},
		},
	})

	Register(Schema{
		Command:     "text.set_opentype_features",
		Aliases:     []string{"text.set_opentype"},
		Description: "Set OpenType feature flags on a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "features", Type: "object", Required: true, Desc: "Object of OpenType features (e.g. {liga:true, smcp:true})"},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})

	Register(Schema{
		Command:     "text.get_opentype_features",
		Aliases:     []string{"text.get_opentype"},
		Description: "Get OpenType feature flags from a text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "rangeStart", Type: "number", Desc: "Character range start"},
			{Name: "rangeEnd", Type: "number", Desc: "Character range end"},
		},
	})
}
