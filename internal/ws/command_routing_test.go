package ws

import (
	"encoding/json"
	"testing"
)

func TestResolveCommandRoute_Legacy(t *testing.T) {
	domain, action, err := resolveCommandRoute("set_fill_color", map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if domain != "paint" || action != "set_solid" {
		t.Fatalf("unexpected route: %s.%s", domain, action)
	}
}

func TestResolveCommandRoute_DomainActionPassthrough(t *testing.T) {
	domain, action, err := resolveCommandRoute("paint.set_gradient", map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if domain != "paint" || action != "set_gradient" {
		t.Fatalf("unexpected route: %s.%s", domain, action)
	}
}

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
