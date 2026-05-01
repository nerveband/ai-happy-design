package schema

import "testing"

func TestHighLevelCommandSchemasAreRegistered(t *testing.T) {
	tests := []struct {
		command        string
		requiredParam  string
		optionalParams []string
	}{
		{
			command:        "text.measure",
			requiredParam:  "text",
			optionalParams: []string{"width", "fontFamily", "fontStyle", "fontSize"},
		},
		{
			command:        "text.fit_box",
			requiredParam:  "text",
			optionalParams: []string{"width", "height", "minFontSize", "maxFontSize"},
		},
		{
			command:        "text.create_rich_block",
			requiredParam:  "width",
			optionalParams: []string{"parentId", "heading", "price", "bullets", "benefits", "eligibility"},
		},
		{
			command:        "layout.pricing_grid",
			requiredParam:  "width",
			optionalParams: []string{"parentId", "columns", "gap", "rowGap", "cards"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			s := Lookup(tt.command)
			if s == nil {
				t.Fatalf("expected schema for %s to be registered", tt.command)
			}
			if s.Command != tt.command {
				t.Fatalf("expected canonical command %q, got %q", tt.command, s.Command)
			}

			required := LookupParam(s, tt.requiredParam)
			if required == nil {
				t.Fatalf("expected param %q to be registered", tt.requiredParam)
			}
			if !required.Required {
				t.Fatalf("expected param %q to be required", tt.requiredParam)
			}

			for _, name := range tt.optionalParams {
				if LookupParam(s, name) == nil {
					t.Fatalf("expected param %q to be registered", name)
				}
			}
		})
	}
}

func TestHighLevelCommandSchemaAliases(t *testing.T) {
	s := Lookup("text.rich")
	if s == nil {
		t.Fatal("expected text.rich alias to resolve")
	}
	if s.Command != "text.create_rich_block" {
		t.Fatalf("expected text.rich to resolve to text.create_rich_block, got %q", s.Command)
	}

	parent := LookupParam(s, "pid")
	if parent == nil {
		t.Fatal("expected pid alias to resolve for parentId")
	}
	if parent.Name != "parentId" {
		t.Fatalf("expected pid to resolve to parentId, got %q", parent.Name)
	}
}
