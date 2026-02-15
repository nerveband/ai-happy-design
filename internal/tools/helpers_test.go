package tools

import "testing"

func TestHasArg_Present(t *testing.T) {
	args := map[string]interface{}{
		"color": "#FF0000",
		"width": 100.0,
	}
	if !hasArg(args, "color") {
		t.Error("expected hasArg to return true for existing key 'color'")
	}
	if !hasArg(args, "width") {
		t.Error("expected hasArg to return true for existing key 'width'")
	}
}

func TestHasArg_ZeroValue(t *testing.T) {
	args := map[string]interface{}{
		"empty":  "",
		"zero":   0.0,
		"falsy":  false,
		"nilVal": nil,
	}
	for key := range args {
		if !hasArg(args, key) {
			t.Errorf("expected hasArg to return true for zero-valued key %q", key)
		}
	}
}

func TestHasArg_Missing(t *testing.T) {
	args := map[string]interface{}{
		"color": "#FF0000",
	}
	if hasArg(args, "missing") {
		t.Error("expected hasArg to return false for missing key")
	}
	if hasArg(args, "") {
		t.Error("expected hasArg to return false for empty key")
	}
}

func TestHasArg_EmptyMap(t *testing.T) {
	args := map[string]interface{}{}
	if hasArg(args, "anything") {
		t.Error("expected hasArg to return false for empty map")
	}
}
