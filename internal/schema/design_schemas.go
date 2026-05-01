package schema

func init() {
	Register(Schema{
		Command:       "design.compute_tokens",
		Description:   "Compute responsive design tokens for a target canvas",
		Safety:        "local",
		Idempotency:   "idempotent",
		RequiresFigma: false,
		Params: []Param{
			{Name: "width", Type: "number", Required: true, Desc: "Canvas width in pixels", Min: Ptr(1)},
			{Name: "height", Type: "number", Required: true, Desc: "Canvas height in pixels", Min: Ptr(1)},
			{Name: "dpi", Type: "number", Desc: "Optional DPI. Defaults to 72 for screen work."},
		},
	})
	Register(Schema{
		Command:       "tokens.export",
		Aliases:       []string{"tokens.export_json", "tokens.export_css", "tokens.export_tailwind", "tokens.export_swift", "tokens.export_android"},
		Description:   "Export a deterministic token preset or computed token snapshot",
		Safety:        "local",
		Idempotency:   "idempotent",
		RequiresFigma: false,
		Params: []Param{
			{Name: "preset", Type: "string", Desc: "Token preset name", Enum: []string{"minimal", "shadcn", "tailwind", "from_canvas"}},
			{Name: "format", Type: "string", Desc: "Output format", Enum: []string{"json", "css", "tailwind", "swift", "android"}},
			{Name: "outputPath", Type: "string", Desc: "Optional file path to write"},
			{Name: "configPath", Type: "string", Desc: "Optional JSON config with outputs map"},
			{Name: "variablesFile", Type: "string", Desc: "JSON file from variable.get_all"},
			{Name: "variables", Type: "array", Desc: "Inline variable.get_all variables array"},
			{Name: "outputs", Type: "object", Desc: "Format to output path map"},
		},
	})
	Register(Schema{
		Command:       "parity.compare_code",
		Description:   "Compare a code-spec JSON contract for semantic completeness before design-code parity checks",
		Safety:        "local",
		Idempotency:   "idempotent",
		RequiresFigma: false,
		Params: []Param{
			{Name: "specPath", Type: "string", Desc: "Path to a code-spec JSON file"},
			{Name: "codeSpecPath", Type: "string", Desc: "Alias for specPath"},
			{Name: "codeSpec", Type: "object", Desc: "Inline code-spec object"},
			{Name: "threshold", Type: "number", Desc: "Minimum passing parity score"},
		},
	})
}
