package schema

func init() {
	Register(Schema{
		Command:     "design.compute_tokens",
		Description: "Compute design tokens (font sizes, spacing, padding, layout, card widths) for any canvas dimensions. Call FIRST before creating any design. Returns concrete pixel values.",
		Params: []Param{
			{Name: "width", Type: "number", Required: true, Desc: "Canvas width in pixels", Min: Ptr(100), Max: Ptr(10000)},
			{Name: "height", Type: "number", Required: true, Desc: "Canvas height in pixels", Min: Ptr(100), Max: Ptr(10000)},
			{Name: "dpi", Type: "number", Desc: "DPI (72 for screen, 300 for print)", Default: 72.0},
		},
	})

	Register(Schema{
		Command:     "design_system.analyze",
		Description: "Analyze the current Figma file's styles, variables, and components. Returns categorized rules for colors, typography, spacing, effects, and components with reuse guidance.",
		Params:      []Param{},
	})
}
