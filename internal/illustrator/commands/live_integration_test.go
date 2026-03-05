//go:build darwin && integration

package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/illustrator/commands"
	illustratorhost "github.com/nerveband/ai-happy-design/internal/illustrator/host"
	_ "github.com/nerveband/ai-happy-design/internal/illustrator/schema"
)

func TestScratchDocumentLiveFlow(t *testing.T) {
	adapter := illustratorhost.NewAdapter()
	status := adapter.Status()
	if !status.IllustratorAppFound {
		t.Skip("Illustrator not installed")
	}
	if !status.IllustratorRunning {
		t.Skip("Illustrator not running")
	}

	executor := commands.NewExecutor()
	defer func() {
		_, _, _ = executor.Execute(commands.Request{
			Command: commonschema.Lookup("document.close"),
			Params:  map[string]any{"save": false},
		})
	}()

	mustExecute(t, executor, "app.version", map[string]any{})
	mustExecute(t, executor, "app.info", map[string]any{})
	mustExecute(t, executor, "app.user_interaction_level", map[string]any{"mode": "DONTDISPLAYALERTS"})
	mustExecute(t, executor, "document.new", map[string]any{"width": 800, "height": 600})
	mustExecute(t, executor, "document.info", map[string]any{})
	mustExecute(t, executor, "document.list", map[string]any{})
	mustExecute(t, executor, "artboard.list", map[string]any{})
	mustExecute(t, executor, "artboard.create", map[string]any{
		"name":   "Secondary Board",
		"left":   820,
		"top":    600,
		"right":  1620,
		"bottom": 0,
	})
	mustExecute(t, executor, "artboard.set_active", map[string]any{"artboardId": "Secondary Board"})
	mustExecute(t, executor, "layer.create", map[string]any{"name": "Live Layer"})
	mustExecute(t, executor, "layer.rename", map[string]any{"layerId": "Live Layer", "name": "Validation Layer"})
	mustExecute(t, executor, "layer.visibility", map[string]any{"layerId": "Validation Layer", "visible": false})
	mustExecute(t, executor, "layer.visibility", map[string]any{"layerId": "Validation Layer", "visible": true})
	mustExecute(t, executor, "layer.lock", map[string]any{"layerId": "Validation Layer", "locked": true})
	mustExecute(t, executor, "layer.lock", map[string]any{"layerId": "Validation Layer", "locked": false})
	mustExecute(t, executor, "path.create_rect", map[string]any{
		"layerId": "Validation Layer",
		"name":    "Validation Rect",
		"left":    60,
		"top":     520,
		"width":   220,
		"height":  140,
	})
	mustExecute(t, executor, "appearance.set_fill", map[string]any{"itemId": "Validation Rect", "color": "#FF5500"})
	mustExecute(t, executor, "appearance.set_stroke", map[string]any{"itemId": "Validation Rect", "color": "#112233", "width": 2})
	mustExecute(t, executor, "appearance.set_gradient", map[string]any{
		"itemId": "Validation Rect",
		"stops": []any{
			map[string]any{"offset": 0, "color": "#FF5500"},
			map[string]any{"offset": 100, "color": "#112233"},
		},
		"type": "linear",
	})
	mustExecute(t, executor, "path.transform", map[string]any{"itemId": "Validation Rect", "translateX": 10, "translateY": -10, "rotate": 5})
	mustExecute(t, executor, "path.duplicate", map[string]any{"itemId": "Validation Rect", "destinationLayerId": "Validation Layer"})
	mustExecute(t, executor, "selection.select_by_name", map[string]any{"name": "Validation Rect"})
	mustExecute(t, executor, "selection.get", map[string]any{})
	mustExecute(t, executor, "inspect.bounds", map[string]any{})
	mustExecute(t, executor, "selection.clear", map[string]any{})
	mustExecute(t, executor, "text.create", map[string]any{
		"layerId":  "Validation Layer",
		"name":     "Validation Text",
		"contents": "Live Validation",
		"left":     80,
		"top":      420,
	})
	mustExecute(t, executor, "text.set_contents", map[string]any{"itemId": "Validation Text", "contents": "Live Bridge Validation"})
	mustExecute(t, executor, "text.set_style", map[string]any{"itemId": "Validation Text", "fontSize": 24, "tracking": 20, "fillColor": "#223344"})
	mustExecute(t, executor, "inspect.fonts", map[string]any{})
	mustExecute(t, executor, "selection.select_by_name", map[string]any{"name": "Validation Text"})
	mustExecute(t, executor, "artboard.fit_to_artwork", map[string]any{})
	styles := mustExecute(t, executor, "inspect.styles", map[string]any{})
	if styleName := firstGraphicStyle(styles); styleName != "" {
		mustExecute(t, executor, "appearance.apply_graphic_style", map[string]any{"itemId": "Validation Rect", "styleName": styleName})
	}
	mustExecute(t, executor, "inspect.tree", map[string]any{})
	mustExecute(t, executor, "inspect.summary", map[string]any{})
	mustExecute(t, executor, "text.outline", map[string]any{"itemId": "Validation Text"})

	outDir := t.TempDir()
	svgPath := filepath.Join(outDir, "live-validation.svg")
	pngPath := filepath.Join(outDir, "live-validation.png")
	mustExecute(t, executor, "export.svg", map[string]any{"outputPath": svgPath})
	mustExecute(t, executor, "export.png", map[string]any{"outputPath": pngPath, "scale": 1})

	if _, err := os.Stat(svgPath); err != nil {
		t.Fatalf("expected svg export at %s: %v", svgPath, err)
	}
	if _, err := os.Stat(pngPath); err != nil {
		t.Fatalf("expected png export at %s: %v", pngPath, err)
	}

	mustExecute(t, executor, "document.close", map[string]any{"save": false})
}

func mustExecute(t *testing.T, executor *commands.Executor, name string, params map[string]any) any {
	t.Helper()

	result, _, err := executor.Execute(commands.Request{
		Command: commonschema.Lookup(name),
		Params:  params,
	})
	if err != nil {
		t.Fatalf("%s failed: %+v", name, err)
	}
	return result
}

func firstGraphicStyle(result any) string {
	root, ok := result.(map[string]any)
	if !ok {
		return ""
	}
	values, ok := root["graphicStyles"].([]any)
	if !ok || len(values) == 0 {
		return ""
	}
	name, _ := values[0].(string)
	return name
}
