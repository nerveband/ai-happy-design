package schema

import "testing"

func TestRemainingRoadmapSchemasAreRegistered(t *testing.T) {
	for _, command := range []string{
		"tokens.preset_tailwind",
		"tokens.preset_shadcn",
		"tokens.preset_material",
		"tokens.setup_system",
		"design_system.health",
		"component.analyze_set",
		"component.arrange_set",
		"parity.audit_component",
		"document.get_editor_context",
		"slides.get_current",
		"slides.create_slide",
		"slides.set_background",
		"slides.add_text",
		"slides.reorder",
		"figjam.create_sticky",
		"figjam.create_shape",
		"figjam.create_connector",
		"figjam.get_board",
	} {
		if Lookup(command) == nil {
			t.Errorf("missing schema %s", command)
		}
	}
}
