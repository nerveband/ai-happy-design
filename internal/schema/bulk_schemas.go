package schema

func init() {
	Register(Schema{
		Command:     "bulk.execute",
		Description: "Execute multiple operations with retry and optional interpolation. Use for batch workflows.",
		Params: []Param{
			{Name: "operations", Type: "array", Required: true, Desc: "JSON array of operations [{name?, command, params}]"},
			{Name: "failFast", Type: "boolean", Desc: "Stop on first failure (default false)", Default: false},
			{Name: "retries", Type: "number", Desc: "Number of retries per failed operation", Min: Ptr(0), Max: Ptr(5), Default: 0.0},
			{Name: "retryDelayMs", Type: "number", Desc: "Delay between retries in milliseconds", Min: Ptr(0), Max: Ptr(10000), Default: 500.0},
			{Name: "interpolate", Type: "boolean", Desc: "Enable step result interpolation via ${{steps.NAME.result.field}}", Default: false},
		},
	})

	Register(Schema{
		Command:     "bulk.slide",
		Description: "Composite: creates a full social media slide. Auto-expands into multiple primitive ops.",
		Params: []Param{
			{Name: "canvas", Type: "string", Desc: "Canvas dimensions as WxH string (e.g. '1080x1920')", Default: "1080x1920"},
			{Name: "bg", Type: "string", Desc: "Background color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "gradient", Type: "object", Desc: "Gradient specification {stops, angle}"},
			{Name: "elements", Type: "array", Required: true, Desc: "Array of element definitions: eyebrow, headline, body, bar, counter, cta, url, stats, progress, arabic"},
		},
	})

	Register(Schema{
		Command:     "bulk.banner",
		Description: "Composite: creates an email banner. Auto-expands into multiple primitive ops.",
		Params: []Param{
			{Name: "canvas", Type: "string", Desc: "Canvas dimensions as WxH string (e.g. '600x200')", Default: "600x200"},
			{Name: "bg", Type: "string", Desc: "Background color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "gradient", Type: "object", Desc: "Gradient specification {stops, angle}"},
			{Name: "dividerX", Type: "number", Desc: "X position for divider line"},
			{Name: "elements", Type: "array", Required: true, Desc: "Array of element definitions: headline, subtitle"},
		},
	})
}
