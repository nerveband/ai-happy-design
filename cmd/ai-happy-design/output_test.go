package main

import "testing"

func TestDiagnosticsEnabledOnlyForHumanOutput(t *testing.T) {
	original := globalOutputFormat
	defer func() { globalOutputFormat = original }()

	for _, format := range []string{"", "json", "jsonl"} {
		globalOutputFormat = format
		if diagnosticsEnabled() {
			t.Fatalf("diagnostics should be disabled for machine output format %q", format)
		}
	}

	globalOutputFormat = "text"
	if !diagnosticsEnabled() {
		t.Fatal("diagnostics should be enabled for text output")
	}
}
