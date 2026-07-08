package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRemainingRoadmapLocalCommands(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("AHD_CONFIG_DIR", tmp)

	for _, tc := range []struct {
		command string
		params  map[string]interface{}
		check   string
	}{
		{"tokens.preset_tailwind", map[string]interface{}{"outputPath": filepath.Join(tmp, "tw.json")}, "tokens"},
		{"tokens.preset_shadcn", map[string]interface{}{}, "tokens"},
		{"tokens.preset_material", map[string]interface{}{}, "tokens"},
		{"tokens.setup_system", map[string]interface{}{"preset": "tailwind"}, "batch"},
		{"design_system.health", map[string]interface{}{"spec": map[string]interface{}{"nodes": []interface{}{map[string]interface{}{"name": "Button", "type": "COMPONENT"}}}}, "score"},
		{"component.analyze_set", map[string]interface{}{"componentSet": map[string]interface{}{"name": "Button", "variants": []interface{}{"hover", "disabled"}}}, "states"},
		{"component.arrange_set", map[string]interface{}{"componentSet": map[string]interface{}{"variants": []interface{}{"default", "hover"}}}, "layout"},
		{"parity.audit_component", map[string]interface{}{"figmaNode": map[string]interface{}{"name": "Button"}, "codeSpec": map[string]interface{}{"name": "Button"}}, "score"},
		{"verify.visual", map[string]interface{}{"artifactPath": writeTempArtifact(t, tmp)}, "ok"},
	} {
		t.Run(tc.command, func(t *testing.T) {
			handled, result, err := handleLocalCommand(tc.command, tc.params)
			if err != nil {
				t.Fatalf("%s: %v", tc.command, err)
			}
			if !handled {
				t.Fatalf("%s was not handled locally", tc.command)
			}
			raw, _ := json.Marshal(result)
			var m map[string]interface{}
			if err := json.Unmarshal(raw, &m); err != nil {
				t.Fatalf("result is not object json: %v", err)
			}
			if _, ok := m[tc.check]; !ok {
				t.Fatalf("missing %q in %#v", tc.check, m)
			}
		})
	}
}

func TestJobLedgerRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("AHD_CONFIG_DIR", tmp)
	created, err := createJobRecord("test.kind", map[string]interface{}{"hello": "world"})
	if err != nil {
		t.Fatal(err)
	}
	listed, err := listJobRecords()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID {
		t.Fatalf("unexpected jobs: %#v", listed)
	}
	got, err := getJobRecord(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "test.kind" {
		t.Fatalf("unexpected job: %#v", got)
	}
	if err := cancelJobRecord(created.ID); err != nil {
		t.Fatal(err)
	}
}

func writeTempArtifact(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "artifact.png")
	if err := os.WriteFile(path, []byte("png"), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}
