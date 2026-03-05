package commonvalidate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/commonvalidate"
	_ "github.com/nerveband/ai-happy-design/internal/illustrator/schema"
)

func TestValidateCommandNormalizesAliases(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("text.create")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"name": "Headline",
		"text": "Hello",
		"left": 12,
		"top":  24,
	}, t.TempDir())

	if len(result.Errors) != 0 {
		t.Fatalf("expected no validation errors, got %+v", result.Errors)
	}
	if _, ok := result.Params["contents"]; !ok {
		t.Fatalf("expected canonical contents param in normalized params: %+v", result.Params)
	}
	if result.Params["contents"] != "Hello" {
		t.Fatalf("unexpected contents value: %#v", result.Params["contents"])
	}
}

func TestValidateCommandRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("export.png")
	cwd := t.TempDir()
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"path": "../escape.png",
	}, cwd)

	if len(result.Errors) == 0 {
		t.Fatal("expected traversal attempt to be rejected")
	}
}

func TestValidateCommandFuzzyEnum(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("app.user_interaction_level")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"mode": "displayalerts",
	}, t.TempDir())

	if len(result.Errors) != 0 {
		t.Fatalf("expected fuzzy enum normalization, got errors %+v", result.Errors)
	}
	if result.Params["mode"] != "DISPLAYALERTS" {
		t.Fatalf("expected DISPLAYALERTS, got %#v", result.Params["mode"])
	}
}

func TestSafePathReturnsJoinedPath(t *testing.T) {
	t.Parallel()

	cwd := t.TempDir()
	got, err := commonvalidate.SafePath(cwd, filepath.Join("exports", "poster.png"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(cwd, "exports", "poster.png")
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestSafePathRejectsAbsolutePath(t *testing.T) {
	t.Parallel()

	cwd := t.TempDir()
	absolute := filepath.Join(cwd, "poster.png")
	if _, err := commonvalidate.SafePath(cwd, absolute); err == nil {
		t.Fatal("expected absolute path to be rejected")
	}
	if _, err := os.Stat(cwd); err != nil {
		t.Fatalf("temp dir should still exist: %v", err)
	}
}
