package schema

func init() {
	Register(Schema{
		Command:     "page.create",
		Description: "Create a new page in the document",
		Params: []Param{
			{Name: "name", Type: "string", Desc: "Page name"},
		},
	})

	Register(Schema{
		Command:     "page.delete",
		Description: "Delete a page (cannot delete the only page)",
		Params: []Param{
			{Name: "pageId", Type: "string", Required: true, Desc: "Page node ID"},
		},
	})

	Register(Schema{
		Command:     "page.rename",
		Description: "Rename a page",
		Params: []Param{
			{Name: "pageId", Type: "string", Required: true, Desc: "Page node ID"},
			{Name: "name", Type: "string", Required: true, Desc: "New page name"},
		},
	})

	Register(Schema{
		Command:     "page.duplicate",
		Description: "Duplicate a page",
		Params: []Param{
			{Name: "pageId", Type: "string", Required: true, Desc: "Page node ID to duplicate"},
			{Name: "name", Type: "string", Desc: "Name for the duplicated page"},
		},
	})

	Register(Schema{
		Command:     "page.set_current",
		Aliases:     []string{"page.switch", "page.navigate", "page.go_to"},
		Description: "Switch to a specific page (set as current)",
		Params: []Param{
			{Name: "pageId", Type: "string", Required: true, Desc: "Page node ID to switch to"},
		},
	})

	Register(Schema{
		Command:     "page.get_all",
		Aliases:     []string{"page.list"},
		Description: "Get all pages in the document",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "page.get_current",
		Description: "Get the currently active page",
		Params:      []Param{},
	})
}
