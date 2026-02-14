package ws

import (
	"encoding/json"
	"testing"
)

// --- Legacy route tests ---

func TestResolveCommandRoute_Legacy(t *testing.T) {
	domain, action, err := resolveCommandRoute("set_fill_color", map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if domain != "paint" || action != "set_solid" {
		t.Fatalf("unexpected route: %s.%s", domain, action)
	}
}

func TestResolveCommandRoute_AllLegacyRoutes(t *testing.T) {
	for cmd, expected := range legacyCommandRoutes {
		domain, action, err := resolveCommandRoute(cmd, map[string]interface{}{})
		if err != nil {
			t.Errorf("legacy route %q: unexpected error: %v", cmd, err)
			continue
		}
		if domain != expected.Domain || action != expected.Action {
			t.Errorf("legacy route %q: expected %s.%s, got %s.%s",
				cmd, expected.Domain, expected.Action, domain, action)
		}
	}
}

// --- Dot-notation tests ---

func TestResolveCommandRoute_DomainActionPassthrough(t *testing.T) {
	domain, action, err := resolveCommandRoute("paint.set_gradient", map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if domain != "paint" || action != "set_gradient" {
		t.Fatalf("unexpected route: %s.%s", domain, action)
	}
}

func TestResolveCommandRoute_DotNotationVariants(t *testing.T) {
	tests := []struct {
		command string
		domain  string
		action  string
	}{
		{"node.create_frame", "node", "create_frame"},
		{"text.create", "text", "create"},
		{"document.find_free_space", "document", "find_free_space"},
		{"shape.create_rectangle", "shape", "create_rectangle"},
		{"export.image", "export", "image"},
		{"layout.set_auto_layout", "layout", "set_auto_layout"},
		{"page.set_current", "page", "set_current"},
		{"effect.add_shadow", "effect", "add_shadow"},
		{"boolean.union", "boolean", "union"},
		{"variable.create", "variable", "create"},
		{"component.create", "component", "create"},
		{"style.apply", "style", "apply"},
		{"layer.group", "layer", "group"},
		{"connect.something", "connect", "something"},
	}

	for _, tt := range tests {
		domain, action, err := resolveCommandRoute(tt.command, map[string]interface{}{})
		if err != nil {
			t.Errorf("dot-notation %q: unexpected error: %v", tt.command, err)
			continue
		}
		if domain != tt.domain || action != tt.action {
			t.Errorf("dot-notation %q: expected %s.%s, got %s.%s",
				tt.command, tt.domain, tt.action, domain, action)
		}
	}
}

// --- Export format routing ---

func TestResolveCommandRoute_ExportFormat(t *testing.T) {
	tests := []struct {
		format string
		action string
	}{
		{format: "SVG", action: "svg"},
		{format: "PDF", action: "pdf"},
		{format: "PNG", action: "image"},
		{format: "", action: "image"},
	}

	for _, tt := range tests {
		_, action, err := resolveCommandRoute("export_node_as_image", map[string]interface{}{"format": tt.format})
		if err != nil {
			t.Fatalf("unexpected error for format %q: %v", tt.format, err)
		}
		if action != tt.action {
			t.Fatalf("expected action %q for format %q, got %q", tt.action, tt.format, action)
		}
	}
}

// --- Error cases ---

func TestResolveCommandRoute_UnknownCommand(t *testing.T) {
	_, _, err := resolveCommandRoute("nonexistent_command_xyz", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if err.Error() != "unknown command: nonexistent_command_xyz" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestResolveCommandRoute_EmptyString(t *testing.T) {
	_, _, err := resolveCommandRoute("", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestResolveCommandRoute_DotOnly(t *testing.T) {
	// A single dot should not split into valid domain.action
	_, _, err := resolveCommandRoute(".", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for dot-only command")
	}
}

func TestResolveCommandRoute_LeadingDot(t *testing.T) {
	// ".action" — dot at index 0, should not match dot-notation path
	_, _, err := resolveCommandRoute(".action", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for leading-dot command")
	}
}

func TestResolveCommandRoute_TrailingDot(t *testing.T) {
	// "domain." — dot at last position, should not match dot-notation path
	_, _, err := resolveCommandRoute("domain.", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for trailing-dot command")
	}
}

func TestResolveCommandRoute_NilParams(t *testing.T) {
	// Legacy route with nil params should still work
	domain, action, err := resolveCommandRoute("create_frame", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if domain != "node" || action != "create_frame" {
		t.Fatalf("unexpected route: %s.%s", domain, action)
	}
}

// --- Wrapped response extraction tests ---

func TestExtractWrappedResponse_ResultPresenceFalseValue(t *testing.T) {
	raw := map[string]interface{}{
		"id":     "abc-123",
		"result": false,
	}
	data, _ := json.Marshal(raw)
	msg := &Message{Type: "message", Message: data}

	out, ok := extractWrappedResponse(msg)
	if !ok {
		t.Fatal("expected wrapped response to be recognized")
	}
	if out.ID != "abc-123" || out.Type != "response" {
		t.Fatalf("unexpected extracted response: %#v", out)
	}
	if string(out.Result) != "false" {
		t.Fatalf("expected result to preserve false value, got %s", string(out.Result))
	}
}

func TestExtractWrappedResponse_Error(t *testing.T) {
	raw := map[string]interface{}{
		"id":    "abc-456",
		"error": "boom",
	}
	data, _ := json.Marshal(raw)
	msg := &Message{Type: "message", Message: data}

	out, ok := extractWrappedResponse(msg)
	if !ok {
		t.Fatal("expected wrapped error to be recognized")
	}
	if out.ID != "abc-456" || out.Type != "error" || out.Error != "boom" {
		t.Fatalf("unexpected extracted error: %#v", out)
	}
}

func TestExtractWrappedResponse_NoIdOrResult(t *testing.T) {
	raw := map[string]interface{}{
		"something": "else",
	}
	data, _ := json.Marshal(raw)
	msg := &Message{Type: "message", Message: data}

	_, ok := extractWrappedResponse(msg)
	if ok {
		t.Fatal("expected non-response message to not be recognized")
	}
}

func TestExtractWrappedResponse_NullResult(t *testing.T) {
	raw := map[string]interface{}{
		"id":     "abc-789",
		"result": nil,
	}
	data, _ := json.Marshal(raw)
	msg := &Message{Type: "message", Message: data}

	out, ok := extractWrappedResponse(msg)
	if !ok {
		t.Fatal("expected wrapped response with null result to be recognized")
	}
	if out.ID != "abc-789" || out.Type != "response" {
		t.Fatalf("unexpected extracted response: %#v", out)
	}
}
