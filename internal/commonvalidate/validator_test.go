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

func TestValidateCommandRejectsQueryLikeIdentifier(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("path.transform")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"itemId": "rect?nodeId=123",
	}, t.TempDir())

	if len(result.Errors) == 0 {
		t.Fatal("expected query-like identifier to be rejected")
	}
}

func TestValidateCommandRejectsFragmentLikeIdentifier(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("appearance.apply_graphic_style")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"itemId":    "Validation Rect",
		"styleName": "Brand Style#1",
	}, t.TempDir())

	if len(result.Errors) == 0 {
		t.Fatal("expected fragment-like identifier to be rejected")
	}
}

func TestValidateCommandValidatesGradientStopsSchema(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("appearance.set_gradient")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"itemId": "Hero Rect",
		"stops": []any{
			map[string]any{"offset": 0.0, "color": "#FF5500"},
			map[string]any{"offset": 100.0, "color": "#112233"},
		},
		"type": "linear",
	}, t.TempDir())

	if len(result.Errors) != 0 {
		t.Fatalf("expected valid gradient stops, got %+v", result.Errors)
	}
}

func TestValidateCommandRejectsInvalidGradientStops(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("appearance.set_gradient")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"itemId": "Hero Rect",
		"stops": []any{
			map[string]any{"offset": 0.0, "color": "#FF5500"},
			map[string]any{"offset": 120.0},
		},
	}, t.TempDir())

	if len(result.Errors) == 0 {
		t.Fatal("expected invalid gradient stop payload to fail")
	}
}

func TestValidateCommandRejectsMalformedPathPoints(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("path.create_path")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"name":   "Line",
		"points": []any{[]any{0.0, 0.0}, []any{10.0}},
	}, t.TempDir())

	if len(result.Errors) == 0 {
		t.Fatal("expected malformed point array to fail")
	}
}

func TestValidateCommandAcceptsDocumentNewArtboardLayout(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("document.new")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"width":              1440.0,
		"height":             900.0,
		"artboards":          4.0,
		"artboardLayout":     "GridByRow",
		"artboardSpacing":    24.0,
		"artboardRowsOrCols": 1.0,
		"colorSpace":         "RGB",
		"preset":             "Web",
	}, t.TempDir())

	if len(result.Errors) != 0 {
		t.Fatalf("expected document.new layout payload to validate, got %+v", result.Errors)
	}
}

func TestValidateCommandRejectsImpossibleDocumentNewGridLayout(t *testing.T) {
	t.Parallel()

	command := commonschema.Lookup("document.new")
	result := commonvalidate.ValidateCommand(command, map[string]any{
		"artboards":          2.0,
		"artboardLayout":     "GridByRow",
		"artboardRowsOrCols": 2.0,
	}, t.TempDir())

	if len(result.Errors) == 0 {
		t.Fatal("expected impossible grid layout payload to be rejected")
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

func TestSafePathRejectsURLLikeSyntax(t *testing.T) {
	t.Parallel()

	cwd := t.TempDir()
	if _, err := commonvalidate.SafePath(cwd, "file:///tmp/poster.ai"); err == nil {
		t.Fatal("expected URL-like path to be rejected")
	}
	if _, err := commonvalidate.SafePath(cwd, "%2e%2e/escape.ai"); err == nil {
		t.Fatal("expected encoded traversal to be rejected")
	}
}
