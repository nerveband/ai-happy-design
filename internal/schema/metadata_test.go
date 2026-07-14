package schema

import "testing"

func TestAllSchemasExposeSafetyMetadata(t *testing.T) {
	validSafety := map[string]bool{"read": true, "write": true, "destructive": true, "local": true}
	validIdempotency := map[string]bool{"idempotent": true, "non_idempotent": true, "unknown": true}

	for _, s := range All {
		if !validSafety[s.Safety] {
			t.Errorf("%s missing valid safety metadata: %q", s.Command, s.Safety)
		}
		if !validIdempotency[s.Idempotency] {
			t.Errorf("%s missing valid idempotency metadata: %q", s.Command, s.Idempotency)
		}
		if s.Safety == "read" && s.SupportsDryRun {
			t.Errorf("%s read command should not advertise dry-run", s.Command)
		}
		if s.Safety == "write" && !s.SupportsDryRun {
			t.Errorf("%s write command should advertise dry-run", s.Command)
		}
		if s.RequiresRelay && !s.RequiresFigma {
			t.Errorf("%s requiring relay should also require Figma", s.Command)
		}
	}
}

func TestRoadmapSchemasAreRegistered(t *testing.T) {
	for _, command := range []string{
		"document.screenshot",
		"document.screenshot_selection",
		"verify.visual",
		"motion.get_styles",
		"motion.apply_style",
		"motion.remove_style",
		"motion.get_animations",
		"motion.apply_keyframes",
		"motion.remove_keyframes",
		"motion.set_timeline_duration",
		"effect.list_shaders",
		"effect.import_shader",
		"effect.apply_shader_effect",
		"fill.apply_shader",
		"component.create_slot",
		"component.reset_slot",
		"component.get_slots",
		"layout.reorder_grid_rows",
		"layout.reorder_grid_columns",
		"layout.audit",
	} {
		if Lookup(command) == nil {
			t.Errorf("missing roadmap schema %s", command)
		}
	}
}
