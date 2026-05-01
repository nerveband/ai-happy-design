package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunLocalCommandTokensExport(t *testing.T) {
	result, err := runLocalCommand("tokens.export", map[string]interface{}{"preset": "minimal"})
	if err != nil {
		t.Fatalf("tokens.export failed: %v", err)
	}
	out := result.(map[string]interface{})
	if out["preset"] != "minimal" {
		t.Fatalf("unexpected preset: %#v", out["preset"])
	}
}

func TestRunLocalCommandTokensExportVariablesConfig(t *testing.T) {
	dir := t.TempDir()
	varsPath := filepath.Join(dir, "vars.json")
	cssPath := filepath.Join(dir, "tokens.css")
	swiftPath := filepath.Join(dir, "FigmaTokens.swift")
	configPath := filepath.Join(dir, "tokens.config.json")
	vars := []byte(`{"variables":[{"name":"Color/Primary","resolvedType":"COLOR","valuesByMode":{"m":{"r":0.1,"g":0.2,"b":0.3,"a":1}}},{"name":"Space/Md","resolvedType":"FLOAT","valuesByMode":{"m":16}}]}`)
	if err := os.WriteFile(varsPath, vars, 0644); err != nil {
		t.Fatal(err)
	}
	config := []byte(`{"variablesFile":"` + varsPath + `","outputs":{"css":"` + cssPath + `","swift":"` + swiftPath + `"}}`)
	if err := os.WriteFile(configPath, config, 0644); err != nil {
		t.Fatal(err)
	}
	result, err := runLocalCommand("tokens.export", map[string]interface{}{"configPath": configPath})
	if err != nil {
		t.Fatalf("tokens.export failed: %v", err)
	}
	out := result.(map[string]interface{})
	if out["preset"] != "figma_variables" {
		t.Fatalf("expected figma_variables preset, got %#v", out["preset"])
	}
	if _, err := os.Stat(cssPath); err != nil {
		t.Fatalf("css output missing: %v", err)
	}
	if _, err := os.Stat(swiftPath); err != nil {
		t.Fatalf("swift output missing: %v", err)
	}
}

func TestRunLocalCommandAccessibilityAudit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ops.json")
	data := []byte(`[{"command":"text.create","params":{"name":"Tiny","fontSize":8}}]`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	result, err := runLocalCommand("document.accessibility_audit", map[string]interface{}{"file": path})
	if err != nil {
		t.Fatalf("audit failed: %v", err)
	}
	out := result.(map[string]interface{})
	if out["ok"].(bool) {
		t.Fatalf("expected audit finding")
	}
}

func TestRunLocalCommandAccessibilityAuditWCAGCoverage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ops.json")
	data := []byte(`[{"command":"text.create","params":{"name":"Error Message","fontSize":10,"color":"#777777","backgroundColor":"#777777","lineHeight":100}},{"command":"shape.create_rectangle","params":{"name":"Ghost","opacity":0.1}}]`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	result, err := runLocalCommand("document.accessibility_audit", map[string]interface{}{"file": path})
	if err != nil {
		t.Fatalf("audit failed: %v", err)
	}
	out := result.(map[string]interface{})
	codes := map[string]bool{}
	for _, raw := range out["findings"].([]map[string]interface{}) {
		codes[raw["code"].(string)] = true
	}
	for _, code := range []string{"WCAG_TEXT_CONTRAST", "WCAG_LINE_HEIGHT", "WCAG_NON_TEXT_CONTRAST", "WCAG_COLOR_ONLY"} {
		if !codes[code] {
			t.Fatalf("expected %s in findings %#v", code, out["findings"])
		}
	}
}

func TestRunLocalCommandParityCompareNodeProperties(t *testing.T) {
	spec := map[string]interface{}{
		"colors": map[string]interface{}{"primary": "#000"},
		"nodes":  []interface{}{map[string]interface{}{"name": "Button", "type": "FRAME", "properties": map[string]interface{}{"colors": map[string]interface{}{"primary": "#000"}}}},
	}
	result, err := runLocalCommand("parity.compare_code", map[string]interface{}{"codeSpec": spec, "threshold": 0.0})
	if err != nil {
		t.Fatal(err)
	}
	out := result.(map[string]interface{})
	if out["checked"].(int) != 1 {
		t.Fatalf("unexpected checked: %#v", out["checked"])
	}
	if len(out["findings"].([]map[string]interface{})) == 0 {
		t.Fatalf("expected deeper parity findings")
	}
}

func TestApplySimpleJQ(t *testing.T) {
	got, err := applySimpleJQ(map[string]interface{}{"summary": map[string]interface{}{"ok": true}}, ".summary.ok")
	if err != nil {
		t.Fatal(err)
	}
	if got != true {
		t.Fatalf("got %#v", got)
	}
}
