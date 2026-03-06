//go:build darwin && integration

package commands_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/illustrator/commands"
	illustratorhost "github.com/nerveband/ai-happy-design/internal/illustrator/host"
	_ "github.com/nerveband/ai-happy-design/internal/illustrator/schema"
)

func TestPhase1LiveFlow(t *testing.T) {
	adapter := illustratorhost.NewAdapter()
	status := adapter.Status()
	if !status.IllustratorAppFound {
		t.Skip("Illustrator not installed")
	}

	t.Run("utility", func(t *testing.T) {
		withFreshPhase1Session(t, adapter, func(executor *commands.Executor) {
			mustExecute(t, executor, "matrix.identity", map[string]any{})
			mustExecute(t, executor, "matrix.compare", map[string]any{
				"matrix":       identityMatrix(),
				"secondMatrix": identityMatrix(),
			})
			mustExecute(t, executor, "print.presets", map[string]any{})
			printDevices := mustExecute(t, executor, "print.devices", map[string]any{})
			if got := stringField(printDevices, "printerEnumeration"); got != "skipped" {
				t.Fatalf("print.devices printerEnumeration = %q, want skipped", got)
			}
			if got := stringField(printDevices, "ppdEnumeration"); got != "skipped" {
				t.Fatalf("print.devices ppdEnumeration = %q, want skipped", got)
			}
			mustExecute(t, executor, "app.translate_placeholder_text", map[string]any{
				"text": "Lorem ipsum",
			})
			presetLists := mustExecute(t, executor, "app.preset_lists", map[string]any{})
			if got := stringField(presetLists, "printerEnumeration"); got != "unavailable" {
				t.Fatalf("app.preset_lists printerEnumeration = %q, want unavailable", got)
			}
			if got := stringField(presetLists, "ppdEnumeration"); got != "unavailable" {
				t.Fatalf("app.preset_lists ppdEnumeration = %q, want unavailable", got)
			}
			mustExecute(t, executor, "app.get_preset_file", map[string]any{
				"presetType": "Print",
			})
			if presetName := firstNamedPresetFromField(presetLists, "startupPresets"); presetName != "" {
				optionalExecute(t, executor, "app.get_preset_settings", map[string]any{
					"preset": presetName,
				})
			}

			prefKey := "AHD/Phase1/LivePreference"
			mustExecute(t, executor, "preference.set", map[string]any{
				"key":       prefKey,
				"valueType": "string",
				"value": map[string]any{
					"string": "phase1",
				},
			})
			mustExecute(t, executor, "preference.get", map[string]any{
				"key":       prefKey,
				"valueType": "string",
			})
			mustExecute(t, executor, "preference.delete", map[string]any{
				"key": prefKey,
			})
		})
	})

	t.Run("geometry_and_view", func(t *testing.T) {
		withFreshPhase1Session(t, adapter, func(executor *commands.Executor) {
			mustExecute(t, executor, "document.new", map[string]any{
				"width":              900,
				"height":             700,
				"artboards":          3,
				"artboardLayout":     "GridByRow",
				"artboardSpacing":    24,
				"artboardRowsOrCols": 1,
				"colorSpace":         "RGB",
			})
			mustExecute(t, executor, "app.redraw", map[string]any{})
			viewInfo := mustExecute(t, executor, "view.info", map[string]any{})
			rulerVisible := boolField(viewInfo, "rulerVisible")
			gridVisible := boolField(viewInfo, "transparencyGridVisible")
			mustExecute(t, executor, "view.set_zoom", map[string]any{"zoom": 1.25})
			mustExecute(t, executor, "view.set_ruler_visibility", map[string]any{"visible": !rulerVisible})
			mustExecute(t, executor, "view.set_ruler_visibility", map[string]any{"visible": rulerVisible})
			mustExecute(t, executor, "view.set_transparency_grid_visibility", map[string]any{"visible": !gridVisible})
			mustExecute(t, executor, "view.set_transparency_grid_visibility", map[string]any{"visible": gridVisible})
			mustExecute(t, executor, "view.set_center", map[string]any{"x": 300, "y": 300})
			mustExecute(t, executor, "view.rotate", map[string]any{"angle": 3})
			mustExecute(t, executor, "perspective.get_active_plane", map[string]any{})
			mustExecute(t, executor, "perspective.show", map[string]any{})
			mustExecute(t, executor, "perspective.hide", map[string]any{})
			mustExecute(t, executor, "artboard.rearrange", map[string]any{
				"artboardLayout":     "GridByRow",
				"artboardRowsOrCols": 1,
				"artboardSpacing":    20,
				"moveArtwork":        true,
			})
			mustExecute(t, executor, "artboard.delete", map[string]any{
				"artboardId": "Artboard 3",
			})
			mustExecute(t, executor, "layer.create", map[string]any{"name": "Phase1 Geometry"})
			mustExecute(t, executor, "layer.create", map[string]any{"name": "Phase1 Ops"})
			mustExecute(t, executor, "path.create_polygon", map[string]any{
				"layerId": "Phase1 Geometry",
				"name":    "Polygon Source",
				"left":    150,
				"top":     560,
				"radius":  70,
				"sides":   6,
			})
			mustExecute(t, executor, "path.create_star", map[string]any{
				"layerId":     "Phase1 Geometry",
				"name":        "Star Source",
				"left":        360,
				"top":         560,
				"outerRadius": 70,
				"innerRadius": 24,
				"points":      5,
			})
			mustExecute(t, executor, "path.create_rounded_rect", map[string]any{
				"layerId":          "Phase1 Geometry",
				"name":             "Rounded Rect",
				"left":             520,
				"top":              620,
				"width":            160,
				"height":           120,
				"horizontalRadius": 18,
				"verticalRadius":   12,
			})
			mustExecute(t, executor, "path.create_path", map[string]any{
				"layerId": "Phase1 Geometry",
				"name":    "Editable Path",
				"points": []any{
					[]any{120.0, 420.0},
					[]any{380.0, 390.0},
				},
				"closed": false,
			})
			mustExecute(t, executor, "path.set_entire_path", map[string]any{
				"itemId": "Editable Path",
				"points": []any{
					[]any{120.0, 420.0},
					[]any{250.0, 470.0},
					[]any{380.0, 390.0},
				},
				"closed": false,
			})
			duplicate := mustExecute(t, executor, "page_item.duplicate", map[string]any{
				"itemId":    "Polygon Source",
				"targetId":  "Phase1 Ops",
				"placement": "inside",
			})
			duplicateName := stringField(duplicate, "name")
			if duplicateName == "" {
				t.Fatalf("page_item.duplicate returned no name")
			}
			mustExecute(t, executor, "page_item.move", map[string]any{
				"itemId":    duplicateName,
				"targetId":  "Phase1 Geometry",
				"placement": "inside",
			})
			mustExecute(t, executor, "page_item.translate", map[string]any{
				"itemId":                 duplicateName,
				"deltaX":                 80,
				"deltaY":                 -40,
				"transformObjects":       true,
				"transformFillGradients": true,
			})
			mustExecute(t, executor, "page_item.resize", map[string]any{
				"itemId":           duplicateName,
				"scaleX":           120,
				"scaleY":           85,
				"changePositions":  true,
				"changeLineWidths": 120,
				"anchor":           "CENTER",
			})
			mustExecute(t, executor, "page_item.rotate", map[string]any{
				"itemId":          duplicateName,
				"angle":           10,
				"changePositions": true,
				"anchor":          "CENTER",
			})
			mustExecute(t, executor, "page_item.transform", map[string]any{
				"itemId": duplicateName,
				"matrix": map[string]any{
					"mValueA":  1.0,
					"mValueB":  0.0,
					"mValueC":  0.0,
					"mValueD":  1.0,
					"mValueTX": 25.0,
					"mValueTY": -15.0,
				},
				"changePositions": true,
				"anchor":          "CENTER",
			})
			mustExecute(t, executor, "page_item.z_order", map[string]any{
				"itemId": duplicateName,
				"method": "BRINGTOFRONT",
			})
			optionalExecute(t, executor, "page_item.bring_in_perspective", map[string]any{
				"itemId": duplicateName,
				"posX":   240,
				"posY":   300,
				"plane":  "GRIDLEFTPLANETYPE",
			})
			mustExecute(t, executor, "page_item.remove", map[string]any{
				"itemId": duplicateName,
			})
			mustExecute(t, executor, "selection.select_active_artboard_objects", map[string]any{})
			mustExecute(t, executor, "selection.clear", map[string]any{})
		})
	})

	t.Run("text_and_styles", func(t *testing.T) {
		withFreshPhase1Session(t, adapter, func(executor *commands.Executor) {
			mustExecute(t, executor, "document.new", map[string]any{
				"width":      900,
				"height":     700,
				"colorSpace": "RGB",
			})
			mustExecute(t, executor, "layer.create", map[string]any{"name": "Phase1 Text"})
			mustExecute(t, executor, "path.create_path", map[string]any{
				"layerId": "Phase1 Text",
				"name":    "Path Line",
				"points": []any{
					[]any{120.0, 420.0},
					[]any{380.0, 390.0},
				},
				"closed": false,
			})
			mustExecute(t, executor, "selection.clear", map[string]any{})
			mustExecute(t, executor, "text.create", map[string]any{
				"layerId":  "Phase1 Text",
				"name":     "Point Text",
				"contents": "Phase 1",
				"left":     90,
				"top":      300,
			})
			mustExecute(t, executor, "text.create_area", map[string]any{
				"layerId":  "Phase1 Text",
				"name":     "Area Text One",
				"contents": "Area text one for live validation.",
				"left":     520,
				"top":      480,
				"width":    180,
				"height":   90,
			})
			mustExecute(t, executor, "text.create_area", map[string]any{
				"layerId":  "Phase1 Text",
				"name":     "Area Text Two",
				"contents": "Area text two for threading.",
				"left":     720,
				"top":      480,
				"width":    140,
				"height":   90,
			})
			mustExecute(t, executor, "text.create_area", map[string]any{
				"layerId":  "Phase1 Text",
				"name":     "Convertible Area Text",
				"contents": "Convertible area text.",
				"left":     520,
				"top":      340,
				"width":    180,
				"height":   80,
			})
			mustExecute(t, executor, "path.create_rounded_rect", map[string]any{
				"layerId":          "Phase1 Text",
				"name":             "Graphic Style Target",
				"left":             130,
				"top":              540,
				"width":            120,
				"height":           80,
				"horizontalRadius": 12,
				"verticalRadius":   12,
			})
			mustExecute(t, executor, "text.create_on_path", map[string]any{
				"name":       "Path Text",
				"contents":   "Path text validation",
				"pathItemId": "Path Line",
			})
			mustExecute(t, executor, "text.change_case", map[string]any{
				"itemId":   "Point Text",
				"caseType": "UPPERCASE",
			})
			charStyles := mustExecute(t, executor, "style.character.list", map[string]any{})
			if styleName := firstNamedItem(charStyles); styleName != "" {
				mustExecute(t, executor, "style.character.apply", map[string]any{
					"styleName": styleName,
					"itemId":    "Path Text",
				})
			}
			paraStyles := mustExecute(t, executor, "style.paragraph.list", map[string]any{})
			if styleName := firstNamedItem(paraStyles); styleName != "" {
				mustExecute(t, executor, "style.paragraph.apply", map[string]any{
					"styleName":      styleName,
					"itemId":         "Area Text Two",
					"clearOverrides": true,
				})
			}
			graphicStyles := mustExecute(t, executor, "style.graphic.list", map[string]any{})
			if styleName := firstNamedItem(graphicStyles); styleName != "" {
				mustExecute(t, executor, "style.graphic.apply", map[string]any{
					"styleName": styleName,
					"itemId":    "Graphic Style Target",
				})
				optionalExecute(t, executor, "style.graphic.merge", map[string]any{
					"styleName": styleName,
					"itemId":    "Graphic Style Target",
				})
			}
			mustExecute(t, executor, "text.thread", map[string]any{
				"itemId":     "Area Text One",
				"nextItemId": "Area Text Two",
			})
			mustExecute(t, executor, "text.convert_to_area", map[string]any{
				"itemId": "Point Text",
			})
			mustExecute(t, executor, "text.convert_to_point", map[string]any{
				"itemId": "Convertible Area Text",
			})
		})
	})

	t.Run("colors_symbols_and_data", func(t *testing.T) {
		withFreshPhase1Session(t, adapter, func(executor *commands.Executor) {
			mustExecute(t, executor, "document.new", map[string]any{
				"width":      900,
				"height":     700,
				"colorSpace": "RGB",
			})
			mustExecute(t, executor, "layer.create", map[string]any{"name": "Phase1 Data"})
			mustExecute(t, executor, "path.create_polygon", map[string]any{
				"layerId": "Phase1 Data",
				"name":    "Visibility Source",
				"left":    150,
				"top":     560,
				"radius":  70,
				"sides":   6,
			})
			mustExecute(t, executor, "text.create", map[string]any{
				"layerId":  "Phase1 Data",
				"name":     "Variable Text",
				"contents": "Dataset One",
				"left":     90,
				"top":      300,
			})
			mustExecute(t, executor, "swatch.create", map[string]any{
				"name":      "Phase1 Swatch",
				"colorMode": "RGB",
				"hex":       "#FF5500",
			})
			mustExecute(t, executor, "swatch.list", map[string]any{})
			mustExecute(t, executor, "swatch.delete", map[string]any{"swatchId": "Phase1 Swatch"})
			mustExecute(t, executor, "spot.create", map[string]any{
				"name":      "Phase1 Spot",
				"colorMode": "RGB",
				"hex":       "#2255FF",
				"colorType": "SPOT",
			})
			mustExecute(t, executor, "spot.list", map[string]any{})
			mustExecute(t, executor, "spot.delete", map[string]any{"spotId": "Phase1 Spot"})
			optionalExecute(t, executor, "document.write_as_library", map[string]any{
				"filePath": filepath.Join(t.TempDir(), "phase1-symbols.ai"),
			})
			mustExecute(t, executor, "symbol.create", map[string]any{
				"itemId":            "Visibility Source",
				"name":              "Phase1 Symbol",
				"registrationPoint": "SYMBOLCENTERPOINT",
			})
			mustExecute(t, executor, "symbol.place", map[string]any{
				"symbolId": "Phase1 Symbol",
				"name":     "Phase1 Symbol Item",
				"left":     360,
				"top":      560,
			})
			mustExecute(t, executor, "symbol.break_link", map[string]any{
				"itemId": "Phase1 Symbol Item",
			})
			mustExecute(t, executor, "symbol.list", map[string]any{})

			mustExecute(t, executor, "variable.create", map[string]any{
				"name": "Phase1TextVariable",
				"kind": "TEXTUAL",
			})
			mustExecute(t, executor, "variable.bind_text", map[string]any{
				"variableId": "Phase1TextVariable",
				"itemId":     "Variable Text",
			})
			mustExecute(t, executor, "variable.create", map[string]any{
				"name": "Phase1VisibilityVariable",
				"kind": "VISIBILITY",
			})
			mustExecute(t, executor, "variable.bind_visibility", map[string]any{
				"variableId": "Phase1VisibilityVariable",
				"itemId":     "Visibility Source",
			})
			mustExecute(t, executor, "variable.list", map[string]any{})
			mustExecute(t, executor, "dataset.create", map[string]any{
				"name": "Phase1Dataset",
			})
			mustExecute(t, executor, "dataset.list", map[string]any{})
			mustExecute(t, executor, "dataset.apply", map[string]any{
				"datasetId": "Phase1Dataset",
			})
			mustExecute(t, executor, "dataset.update", map[string]any{
				"datasetId": "Phase1Dataset",
			})
			mustExecute(t, executor, "dataset.delete", map[string]any{
				"datasetId": "Phase1Dataset",
			})
			mustExecute(t, executor, "variable.delete", map[string]any{
				"variableId": "Phase1TextVariable",
			})
			mustExecute(t, executor, "variable.delete", map[string]any{
				"variableId": "Phase1VisibilityVariable",
			})
		})
	})

	t.Run("repeat_capture_export_and_trace", func(t *testing.T) {
		withFreshPhase1Session(t, adapter, func(executor *commands.Executor) {
			mustExecute(t, executor, "document.new", map[string]any{
				"width":      900,
				"height":     700,
				"colorSpace": "RGB",
			})
			mustExecute(t, executor, "document.new", map[string]any{
				"width":      400,
				"height":     300,
				"colorSpace": "RGB",
			})
			docs := mustExecute(t, executor, "document.list", map[string]any{})
			if docID := firstNamedItem(docs); docID != "" {
				mustExecute(t, executor, "document.activate", map[string]any{
					"documentId": docID,
				})
			}
			mustExecute(t, executor, "document.arrange", map[string]any{
				"layoutStyle": "CASCADE",
			})
			closePhase1Doc(t, executor)
			waitForIllustratorReady(t, executor)
			mustExecute(t, executor, "app.redraw", map[string]any{})
			mustExecute(t, executor, "layer.create", map[string]any{"name": "Phase1 Export"})
			mustExecute(t, executor, "path.create_polygon", map[string]any{
				"layerId": "Phase1 Export",
				"name":    "Polygon Source",
				"left":    150,
				"top":     560,
				"radius":  70,
				"sides":   6,
			})
			mustExecute(t, executor, "path.create_star", map[string]any{
				"layerId":     "Phase1 Export",
				"name":        "Star Source",
				"left":        360,
				"top":         560,
				"outerRadius": 70,
				"innerRadius": 24,
				"points":      5,
			})
			mustExecute(t, executor, "path.create_star", map[string]any{
				"layerId":     "Phase1 Export",
				"name":        "Raster Source",
				"left":        520,
				"top":         560,
				"outerRadius": 60,
				"innerRadius": 20,
				"points":      5,
			})
			mustExecute(t, executor, "symbol.create", map[string]any{
				"itemId":            "Polygon Source",
				"name":              "Phase1 Symbol",
				"registrationPoint": "SYMBOLCENTERPOINT",
			})
			mustExecute(t, executor, "symbol.place", map[string]any{
				"symbolId": "Phase1 Symbol",
				"name":     "Phase1 Symbol Item",
				"left":     660,
				"top":      620,
			})
			mustExecute(t, executor, "repeat.grid.create", map[string]any{
				"itemId": "Polygon Source",
				"name":   "Grid Repeat",
				"config": map[string]any{"horizontalSpacing": 18, "verticalSpacing": 12},
			})
			mustExecute(t, executor, "repeat.grid.update", map[string]any{
				"repeatId": "Grid Repeat",
				"config":   map[string]any{"horizontalSpacing": 22},
				"state":    "HORIZONTALSPACING",
			})
			mustExecute(t, executor, "repeat.grid.list", map[string]any{})
			mustExecute(t, executor, "repeat.radial.create", map[string]any{
				"itemId": "Star Source",
				"name":   "Radial Repeat",
				"config": map[string]any{"numberOfInstances": 4, "radius": 100},
			})
			mustExecute(t, executor, "repeat.radial.update", map[string]any{
				"repeatId": "Radial Repeat",
				"config":   map[string]any{"numberOfInstances": 5},
				"state":    "NUMBEROFINSTANCES",
			})
			mustExecute(t, executor, "repeat.radial.list", map[string]any{})
			mustExecute(t, executor, "repeat.symmetry.create", map[string]any{
				"itemId": "Phase1 Symbol Item",
				"name":   "Symmetry Repeat",
				"config": map[string]any{"axisRotationAngleInRadians": 0.5},
			})
			mustExecute(t, executor, "repeat.symmetry.update", map[string]any{
				"repeatId": "Symmetry Repeat",
				"config":   map[string]any{"axisRotationAngleInRadians": 1.0},
				"state":    "AXISROTATION",
			})
			mustExecute(t, executor, "repeat.symmetry.list", map[string]any{})

			outDir := t.TempDir()
			capturePath := filepath.Join(outDir, "capture.png")
			windowPath := filepath.Join(outDir, "window.tif")
			webpPath := filepath.Join(outDir, "phase1.webp")
			psdPath := filepath.Join(outDir, "phase1.psd")
			tiffPath := filepath.Join(outDir, "phase1.tif")
			png8Path := filepath.Join(outDir, "phase1-indexed")
			varsPath := filepath.Join(outDir, "variables.xml")
			pdfPresetPath := filepath.Join(outDir, "phase1.joboptions")
			printPresetPath := filepath.Join(outDir, "phase1.prst")

			imageCapture := mustExecute(t, executor, "capture.image", map[string]any{
				"outputPath": capturePath,
				"options": map[string]any{
					"antiAliasing": true,
					"resolution":   144,
					"transparency": true,
				},
			})
			windowCapture := mustExecute(t, executor, "capture.window", map[string]any{
				"outputPath": windowPath,
				"width":      800,
				"height":     600,
			})
			pdfPresetExport := mustExecute(t, executor, "document.export_pdf_preset", map[string]any{
				"filePath": pdfPresetPath,
			})
			pdfPresetArtifact := stringField(pdfPresetExport, "filePath")
			if pdfPresetArtifact == "" {
				t.Fatalf("document.export_pdf_preset returned no filePath")
			}
			printPresetExport := mustExecute(t, executor, "document.export_print_preset", map[string]any{
				"filePath": printPresetPath,
			})
			printPresetArtifact := stringField(printPresetExport, "filePath")
			if printPresetArtifact == "" {
				t.Fatalf("document.export_print_preset returned no filePath")
			}
			png8Export := mustExecute(t, executor, "export.png8", map[string]any{
				"outputPath": png8Path,
				"scale":      1,
				"colorCount": 64,
			})
			png8File := stringField(png8Export, "outputPath")
			if png8File == "" {
				t.Fatalf("export.png8 returned no outputPath")
			}
			tiffExport := mustExecute(t, executor, "export.tiff", map[string]any{
				"outputPath": tiffPath,
				"resolution": 144,
				"artboardId": "Artboard 1",
			})
			tiffArtifact := stringField(tiffExport, "outputPath")
			if tiffArtifact == "" {
				t.Fatalf("export.tiff returned no outputPath")
			}
			webpExport := mustExecute(t, executor, "export.webp", map[string]any{
				"outputPath":          webpPath,
				"losslessCompression": false,
				"imageQuality":        80,
				"isTransparent":       true,
				"ppi":                 144,
			})
			webpArtifact := stringField(webpExport, "outputPath")
			if webpArtifact == "" {
				t.Fatalf("export.webp returned no outputPath")
			}
			psdExport := mustExecute(t, executor, "export.photoshop", map[string]any{
				"outputPath":         psdPath,
				"resolution":         144,
				"editableText":       true,
				"maximumEditability": true,
			})
			psdArtifact := stringField(psdExport, "outputPath")
			if psdArtifact == "" {
				t.Fatalf("export.photoshop returned no outputPath")
			}
			mustExecute(t, executor, "placed.place", map[string]any{
				"filePath": png8File,
				"name":     "Placed Capture",
				"left":     120,
				"top":      180,
			})
			mustExecute(t, executor, "placed.list", map[string]any{})
			placedTrace := mustExecute(t, executor, "placed.trace", map[string]any{
				"itemId":     "Placed Capture",
				"name":       "Placed Trace",
				"presetName": "Default",
			})
			mustExecute(t, executor, "trace.preset.list", map[string]any{})
			mustExecute(t, executor, "placed.place", map[string]any{
				"filePath": png8File,
				"name":     "Placed Embed",
				"left":     260,
				"top":      180,
			})
			mustExecute(t, executor, "placed.embed", map[string]any{
				"itemId": "Placed Embed",
			})
			mustExecute(t, executor, "raster.rasterize", map[string]any{
				"itemId":     "Raster Source",
				"name":       "Rasterized Star",
				"resolution": 144,
			})
			mustExecute(t, executor, "raster.colorize", map[string]any{
				"itemId":    "Rasterized Star",
				"colorMode": "RGB",
				"hex":       "#2288FF",
			})
			mustExecute(t, executor, "raster.list", map[string]any{})
			rasterTrace := mustExecute(t, executor, "raster.trace", map[string]any{
				"itemId":     "Rasterized Star",
				"name":       "Raster Trace",
				"presetName": "Default",
			})
			if pluginItem := stringField(rasterTrace, "pluginItemName"); pluginItem != "" {
				optionalExecute(t, executor, "trace.preset.store", map[string]any{
					"itemId":     pluginItem,
					"presetName": "AHD Phase1 Temp Preset",
				})
				mustExecute(t, executor, "raster.release_tracing", map[string]any{
					"itemId": pluginItem,
				})
			}
			if pluginItem := stringField(placedTrace, "pluginItemName"); pluginItem != "" {
				_ = pluginItem
			}

			mustExecute(t, executor, "dataset.export", map[string]any{"filePath": varsPath})
			mustExecute(t, executor, "variable.export", map[string]any{"filePath": varsPath})
			mustExecute(t, executor, "dataset.import", map[string]any{"filePath": varsPath})
			mustExecute(t, executor, "variable.import", map[string]any{"filePath": varsPath})
			optionalExecute(t, executor, "document.import_pdf_preset", map[string]any{
				"filePath":        pdfPresetArtifact,
				"replacingPreset": false,
			})
			printPresets := optionalExecute(t, executor, "print.presets", map[string]any{})
			if presetName := firstNamedPreset(printPresets); presetName != "" {
				optionalExecute(t, executor, "document.import_print_preset", map[string]any{
					"filePath":    printPresetArtifact,
					"printPreset": presetName,
				})
			}

			imageArtifact := stringField(imageCapture, "outputPath")
			if imageArtifact == "" {
				imageArtifact = capturePath
			}
			windowArtifact := stringField(windowCapture, "outputPath")
			if windowArtifact == "" {
				windowArtifact = windowPath
			}
			for _, path := range []string{imageArtifact, windowArtifact, png8File, tiffArtifact, webpArtifact, psdArtifact, varsPath, pdfPresetArtifact, printPresetArtifact} {
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expected artifact at %s: %v", path, err)
				}
			}
		})
	})
}

func identityMatrix() map[string]any {
	return map[string]any{
		"mValueA":  1.0,
		"mValueB":  0.0,
		"mValueC":  0.0,
		"mValueD":  1.0,
		"mValueTX": 0.0,
		"mValueTY": 0.0,
	}
}

func firstNamedItem(result any) string {
	items, ok := result.([]any)
	if !ok || len(items) == 0 {
		return ""
	}
	for _, item := range items {
		root, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := root["name"].(string)
		if name != "" {
			return name
		}
	}
	return ""
}

func firstNamedPreset(result any) string {
	root, ok := result.(map[string]any)
	if !ok {
		return ""
	}
	items, ok := root["presets"].([]any)
	if !ok || len(items) == 0 {
		return ""
	}
	name, _ := items[0].(string)
	return name
}

func firstNamedPresetFromField(result any, field string) string {
	root, ok := result.(map[string]any)
	if !ok {
		return ""
	}
	items, ok := root[field].([]any)
	if !ok || len(items) == 0 {
		return ""
	}
	name, _ := items[0].(string)
	return name
}

func boolField(result any, field string) bool {
	root, ok := result.(map[string]any)
	if !ok {
		return false
	}
	value, _ := root[field].(bool)
	return value
}

func withFreshPhase1Session(t *testing.T, adapter *illustratorhost.Adapter, fn func(executor *commands.Executor)) {
	t.Helper()

	executor := commands.NewExecutor()
	ensureFreshPhase1Session(t, adapter, executor)
	defer closeAllPhase1Docs(t, adapter, false)
	fn(executor)
}

func closePhase1Doc(t *testing.T, executor *commands.Executor) {
	t.Helper()

	_, _, _ = executor.Execute(commands.Request{
		Command: commonschema.Lookup("document.close"),
		Params:  map[string]any{"save": false},
	})
}

func closeAllPhase1Docs(t *testing.T, adapter *illustratorhost.Adapter, fatal bool) {
	t.Helper()

	_, err := adapter.ExecuteJavaScript(`var previousLevel = app.userInteractionLevel; try { app.userInteractionLevel = UserInteractionLevel.DONTDISPLAYALERTS; while (app.documents.length > 0) { app.documents[app.documents.length - 1].close(SaveOptions.DONOTSAVECHANGES); } "ok"; } finally { try { app.userInteractionLevel = previousLevel; } catch (restoreErr) {} }`, 45*time.Second)
	if err != nil {
		if fatal {
			t.Fatalf("close Illustrator docs: %v", err)
		}
	}
}

func ensureFreshPhase1Session(t *testing.T, adapter *illustratorhost.Adapter, executor *commands.Executor) {
	t.Helper()

	if adapter.Status().IllustratorRunning {
		_ = adapter.Activate()
		waitForIllustratorReady(t, executor)
		if closeAllPhase1DocsReady(adapter) == nil {
			return
		}
		t.Log("close-doc cleanup failed on the running Illustrator process; falling back to restart")
	}

	restartIllustrator(t, adapter)
	waitForIllustratorReady(t, executor)
	closeAllPhase1Docs(t, adapter, true)
}

func closeAllPhase1DocsReady(adapter *illustratorhost.Adapter) error {
	_, err := adapter.ExecuteJavaScript(`var previousLevel = app.userInteractionLevel; try { app.userInteractionLevel = UserInteractionLevel.DONTDISPLAYALERTS; while (app.documents.length > 0) { app.documents[app.documents.length - 1].close(SaveOptions.DONOTSAVECHANGES); } "ok"; } finally { try { app.userInteractionLevel = previousLevel; } catch (restoreErr) {} }`, 45*time.Second)
	return err
}

func waitForIllustratorReady(t *testing.T, executor *commands.Executor) {
	t.Helper()

	deadline := time.Now().Add(45 * time.Second)
	lastErr := ""
	for time.Now().Before(deadline) {
		_, _, err := executor.Execute(commands.Request{
			Command: commonschema.Lookup("app.version"),
			Params:  map[string]any{},
		})
		if err == nil {
			return
		}
		lastErr = err.Message
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("Illustrator did not become ready in time: %s", lastErr)
}

func optionalExecute(t *testing.T, executor *commands.Executor, name string, params map[string]any) any {
	t.Helper()

	result, _, err := executor.Execute(commands.Request{
		Command: commonschema.Lookup(name),
		Params:  params,
	})
	if err != nil {
		t.Logf("%s optional validation failed: %+v", name, err)
		return nil
	}
	return result
}

func stringField(result any, field string) string {
	root, ok := result.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := root[field].(string)
	return value
}

func restartIllustrator(t *testing.T, adapter *illustratorhost.Adapter) {
	t.Helper()

	_ = adapter.Quit()
	quitDeadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(quitDeadline) {
		if !adapter.Status().IllustratorRunning {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if adapter.Status().IllustratorRunning {
		t.Log("Illustrator did not quit cleanly; continuing with the existing process")
	} else {
		if err := adapter.Open(); err != nil {
			t.Fatalf("open Illustrator: %v", err)
		}
	}

	deadline := time.Now().Add(90 * time.Second)
	nextOpenAttempt := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		if err := adapter.Activate(); err != nil {
			t.Logf("activate Illustrator: %v", err)
		}
		status := adapter.Status()
		if status.IllustratorRunning {
			if _, err := adapter.ExecuteJavaScript("app.version;", 5*time.Second); err == nil {
				return
			}
		} else if time.Now().After(nextOpenAttempt) {
			if err := adapter.Open(); err != nil {
				t.Logf("re-open Illustrator: %v", err)
			}
			nextOpenAttempt = time.Now().Add(12 * time.Second)
		}
		time.Sleep(3 * time.Second)
	}
	t.Fatalf("Illustrator did not become ready after restart")
}
