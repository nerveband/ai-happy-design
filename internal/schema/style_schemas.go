package schema

func init() {
	Register(Schema{
		Command:     "style.create_paint",
		Description: "Create a paint style (reusable fill color)",
		Params: []Param{
			{Name: "name", Type: "string", Required: true, Desc: "Style name"},
			{Name: "color", Type: "string", Desc: "Fill color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "description", Type: "string", Desc: "Style description"},
		},
	})

	Register(Schema{
		Command:     "style.create_text",
		Description: "Create a text style (reusable typography settings)",
		Params: []Param{
			{Name: "name", Type: "string", Required: true, Desc: "Style name"},
			{Name: "description", Type: "string", Desc: "Style description"},
		},
	})

	Register(Schema{
		Command:     "style.create_effect",
		Description: "Create an effect style (reusable shadow/blur settings)",
		Params: []Param{
			{Name: "name", Type: "string", Required: true, Desc: "Style name"},
			{Name: "description", Type: "string", Desc: "Style description"},
		},
	})

	Register(Schema{
		Command:     "style.apply",
		Description: "Apply a style to a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "styleId", Type: "string", Required: true, Desc: "Style ID to apply"},
			{Name: "styleType", Type: "string", Aliases: []string{"target"}, Desc: "Type of style to apply", Enum: []string{"FILL", "STROKE", "TEXT", "EFFECT"}, Default: "FILL"},
		},
	})

	Register(Schema{
		Command:     "style.get_all",
		Aliases:     []string{"style.list"},
		Description: "Get all local styles (paint, text, effect, grid)",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "style.remove",
		Description: "Remove a style binding from a node (does not delete the style itself)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "styleType", Type: "string", Aliases: []string{"target"}, Desc: "Type of style to remove", Enum: []string{"FILL", "STROKE", "TEXT", "EFFECT"}, Default: "FILL"},
		},
	})
}
