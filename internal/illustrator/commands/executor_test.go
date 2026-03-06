package commands

import "testing"

func TestCommandWarningsForExperimentalPhase1Commands(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"document.write_as_library",
		"page_item.bring_in_perspective",
		"trace.preset.store",
	} {
		warnings := commandWarnings(name)
		if len(warnings) != 1 {
			t.Fatalf("%s warnings = %d, want 1", name, len(warnings))
		}
		if warnings[0].Code != "EXPERIMENTAL_COMMAND" {
			t.Fatalf("%s warning code = %q, want EXPERIMENTAL_COMMAND", name, warnings[0].Code)
		}
	}
}

func TestCommandWarningsForStableCommands(t *testing.T) {
	t.Parallel()

	if warnings := commandWarnings("document.new"); len(warnings) != 0 {
		t.Fatalf("document.new warnings = %d, want 0", len(warnings))
	}
}
