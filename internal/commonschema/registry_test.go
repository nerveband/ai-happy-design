package commonschema_test

import (
	"testing"

	"github.com/nerveband/ai-happy-design/internal/commonschema"
	_ "github.com/nerveband/ai-happy-design/internal/illustrator/schema"
)

func TestLookupSupportsCommandAliases(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("app.get_info")
	if command == nil {
		t.Fatal("expected alias lookup to resolve app.get_info")
	}
	if command.Name != "app.info" {
		t.Fatalf("expected canonical command app.info, got %q", command.Name)
	}
}

func TestLookupParamSupportsAliases(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("text.create")
	if command == nil {
		t.Fatal("expected text.create schema to be registered")
	}
	param := commonschema.LookupParam(command, "text")
	if param == nil {
		t.Fatal("expected alias lookup for text param")
	}
	if param.Name != "contents" {
		t.Fatalf("expected canonical param contents, got %q", param.Name)
	}
}
