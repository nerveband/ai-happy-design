package commands

import (
	"encoding/json"
	"fmt"
	"time"

	illustratorinspect "github.com/nerveband/ai-happy-design/internal/illustrator/inspect"
)

const scriptPrelude = `
function ahdDoc() {
  if (app.documents.length === 0) {
    throw new Error("No active document");
  }
  return app.activeDocument;
}

function ahdLayer(doc, id) {
  if (!id) return doc.activeLayer || doc.layers[0];
  for (var i = 0; i < doc.layers.length; i++) {
    var layer = doc.layers[i];
    if (layer.name === id || layer.note === id) return layer;
  }
  throw new Error("Layer not found: " + id);
}

function ahdPageItem(doc, id) {
  for (var i = 0; i < doc.pageItems.length; i++) {
    var item = doc.pageItems[i];
    if (item.name === id || item.note === id) return item;
  }
  throw new Error("Page item not found: " + id);
}

function ahdArtboard(doc, id) {
  if (!id) return doc.artboards[doc.artboards.getActiveArtboardIndex()];
  for (var i = 0; i < doc.artboards.length; i++) {
    var artboard = doc.artboards[i];
    if (artboard.name === id) return artboard;
  }
  throw new Error("Artboard not found: " + id);
}

function ahdRGB(hex, opacity) {
  var value = String(hex || "").replace(/^#/, "");
  if (value.length === 8) {
    value = value.substring(0, 6);
  }
  var color = new RGBColor();
  color.red = parseInt(value.substring(0, 2), 16);
  color.green = parseInt(value.substring(2, 4), 16);
  color.blue = parseInt(value.substring(4, 6), 16);
  return { color: color, opacity: opacity == null ? 100 : opacity };
}
`

func buildPlan(request Request) (executionPlan, error) {
	switch request.Command.Name {
	case "app.info":
		return scriptPlan(request.Params, `return {
  name: app.name,
  version: app.version,
  documents: app.documents.length,
  activeDocument: app.documents.length ? app.activeDocument.name : null
};`, 10*time.Second), nil
	case "app.version":
		return scriptPlan(request.Params, `return { version: app.version };`, 10*time.Second), nil
	case "app.select_tool":
		return scriptPlan(request.Params, `app.selectTool(params.tool); return { tool: params.tool };`, 10*time.Second), nil
	case "app.execute_menu":
		return scriptPlan(request.Params, `app.executeMenuCommand(params.menuCommand); return { menuCommand: params.menuCommand };`, 10*time.Second), nil
	case "app.user_interaction_level":
		return scriptPlan(request.Params, `if (params.mode) { app.userInteractionLevel = UserInteractionLevel[params.mode]; } return { mode: String(app.userInteractionLevel) };`, 10*time.Second), nil

	case "document.new":
		return scriptPlan(request.Params, `var width = params.width || 1920; var height = params.height || 1080; var doc = app.documents.add(DocumentColorSpace.RGB, width, height); return { name: doc.name, width: width, height: height };`, 20*time.Second), nil
	case "document.open":
		return scriptPlan(request.Params, `var file = new File(params.filePath); var doc = app.open(file); return { name: doc.name, path: file.fsName };`, 20*time.Second), nil
	case "document.save":
		return scriptPlan(request.Params, `var doc = ahdDoc(); doc.save(); return { name: doc.name, saved: true };`, 20*time.Second), nil
	case "document.save_as":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.saveAs(file); return { name: doc.name, path: file.fsName };`, 20*time.Second), nil
	case "document.close":
		return scriptPlan(request.Params, `var doc = ahdDoc(); doc.close(params.save ? SaveOptions.SAVECHANGES : SaveOptions.DONOTSAVECHANGES); return { closed: true };`, 20*time.Second), nil
	case "document.list":
		return scriptPlan(request.Params, `var docs = []; for (var i = 0; i < app.documents.length; i++) { docs.push({ index: i, name: app.documents[i].name }); } return docs;`, 10*time.Second), nil
	case "document.info":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { name: doc.name, artboards: doc.artboards.length, layers: doc.layers.length, pageItems: doc.pageItems.length };`, 10*time.Second), nil

	case "artboard.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.artboards.length; i++) { var artboard = doc.artboards[i]; items.push({ index: i, name: artboard.name, rect: artboard.artboardRect }); } return items;`, 10*time.Second), nil
	case "artboard.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var rect = [params.left, params.top, params.right, params.bottom]; var artboard = doc.artboards.add(rect); artboard.name = params.name; return { name: artboard.name, rect: rect };`, 20*time.Second), nil
	case "artboard.resize":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var artboard = ahdArtboard(doc, params.artboardId); var rect = artboard.artboardRect; artboard.artboardRect = [rect[0], rect[1], rect[0] + params.width, rect[1] - params.height]; return { name: artboard.name, rect: artboard.artboardRect };`, 20*time.Second), nil
	case "artboard.set_active":
		return scriptPlan(request.Params, `var doc = ahdDoc(); for (var i = 0; i < doc.artboards.length; i++) { if (doc.artboards[i].name === params.artboardId) { doc.artboards.setActiveArtboardIndex(i); return { activeArtboard: params.artboardId, index: i }; } } throw new Error("Artboard not found: " + params.artboardId);`, 10*time.Second), nil
	case "artboard.fit_to_artwork":
		return scriptPlan(request.Params, `var doc = ahdDoc(); doc.fitArtboardToSelectedArt(doc.artboards.getActiveArtboardIndex()); return { fitted: true, activeArtboardIndex: doc.artboards.getActiveArtboardIndex() };`, 20*time.Second), nil

	case "layer.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layers = []; for (var i = 0; i < doc.layers.length; i++) { var layer = doc.layers[i]; layers.push({ index: i, name: layer.name, visible: layer.visible, locked: layer.locked }); } return layers;`, 10*time.Second), nil
	case "layer.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var parent = params.parentLayerId ? ahdLayer(doc, params.parentLayerId) : doc; var layer = parent.layers.add(); layer.name = params.name; return { name: layer.name };`, 20*time.Second), nil
	case "layer.rename":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); layer.name = params.name; return { layerId: params.layerId, name: layer.name };`, 20*time.Second), nil
	case "layer.visibility":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); layer.visible = params.visible; return { layerId: params.layerId, visible: layer.visible };`, 20*time.Second), nil
	case "layer.lock":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); layer.locked = params.locked; return { layerId: params.layerId, locked: layer.locked };`, 20*time.Second), nil
	case "layer.reorder":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var relative = ahdLayer(doc, params.relativeTo); if (params.position === "before") { layer.move(relative, ElementPlacement.PLACEBEFORE); } else if (params.position === "after") { layer.move(relative, ElementPlacement.PLACEAFTER); } else { layer.move(relative, ElementPlacement.INSIDE); } return { layerId: params.layerId, relativeTo: params.relativeTo, position: params.position };`, 20*time.Second), nil

	case "selection.get":
		return scriptPlan(request.Params, `var items = []; for (var i = 0; i < app.selection.length; i++) { items.push({ index: i, name: app.selection[i].name || "", typename: app.selection[i].typename }); } return items;`, 10*time.Second), nil
	case "selection.clear":
		return scriptPlan(request.Params, `app.selection = null; return { cleared: true };`, 10*time.Second), nil
	case "selection.set_by_ids":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var selection = []; for (var i = 0; i < params.ids.length; i++) { selection.push(ahdPageItem(doc, params.ids[i])); } app.selection = selection; return { selected: selection.length };`, 20*time.Second), nil
	case "selection.select_by_name":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var selection = []; for (var i = 0; i < doc.pageItems.length; i++) { var item = doc.pageItems[i]; if (params.partial ? item.name.indexOf(params.name) !== -1 : item.name === params.name) { selection.push(item); } } app.selection = selection; return { selected: selection.length, name: params.name };`, 20*time.Second), nil

	case "path.create_rect":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var item = layer.pathItems.rectangle(params.top, params.left, params.width, params.height); item.name = params.name; if (params.cornerRadius) { item.rounded = true; } return { name: item.name, width: item.width, height: item.height };`, 20*time.Second), nil
	case "path.create_ellipse":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var item = layer.pathItems.ellipse(params.top, params.left, params.width, params.height); item.name = params.name; return { name: item.name, width: item.width, height: item.height };`, 20*time.Second), nil
	case "path.create_path":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var item = layer.pathItems.add(); item.name = params.name; item.setEntirePath(params.points); item.closed = !!params.closed; return { name: item.name, points: params.points.length };`, 20*time.Second), nil
	case "path.transform":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); if (params.translateX || params.translateY) { item.translate(params.translateX || 0, params.translateY || 0); } if (params.scaleX || params.scaleY) { item.resize(params.scaleX || 100, params.scaleY || 100); } if (params.rotate) { item.rotate(params.rotate); } return { itemId: params.itemId, transformed: true };`, 20*time.Second), nil
	case "path.duplicate":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var duplicate = item.duplicate(); if (params.destinationLayerId) { duplicate.move(ahdLayer(doc, params.destinationLayerId), ElementPlacement.INSIDE); } return { source: params.itemId, name: duplicate.name || params.itemId };`, 20*time.Second), nil

	case "text.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var text = layer.textFrames.add(); text.name = params.name; text.contents = params.contents; text.position = [params.left, params.top]; return { name: text.name, contents: text.contents };`, 20*time.Second), nil
	case "text.set_contents":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); text.contents = params.contents; return { itemId: params.itemId, contents: text.contents };`, 20*time.Second), nil
	case "text.set_style":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); var attrs = text.textRange.characterAttributes; if (params.fontFamily) { attrs.textFont = textFonts.getByName(params.fontFamily); } if (params.fontSize) { attrs.size = params.fontSize; } if (params.tracking != null) { attrs.tracking = params.tracking; } if (params.leading != null) { attrs.leading = params.leading; } if (params.fillColor) { attrs.fillColor = ahdRGB(params.fillColor).color; } return { itemId: params.itemId, styled: true };`, 20*time.Second), nil
	case "text.outline":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); var outlined = text.createOutline(); return { itemId: params.itemId, outlinedName: outlined.name || params.itemId };`, 20*time.Second), nil

	case "appearance.set_fill":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var fill = ahdRGB(params.color, params.opacity); item.filled = true; item.fillColor = fill.color; item.opacity = fill.opacity; return { itemId: params.itemId, fill: params.color };`, 20*time.Second), nil
	case "appearance.set_stroke":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var stroke = ahdRGB(params.color); item.stroked = true; item.strokeColor = stroke.color; item.strokeWidth = params.width; return { itemId: params.itemId, stroke: params.color, width: params.width };`, 20*time.Second), nil
	case "appearance.set_gradient", "appearance.apply_graphic_style":
		return selectorPlan("ahd.exec", 20*time.Second), nil

	case "action.load":
		return scriptPlan(request.Params, `app.loadAction(new File(params.filePath)); return { filePath: params.filePath, loaded: true };`, 20*time.Second), nil
	case "action.run":
		return scriptPlan(request.Params, `app.doScript(params.actionName, params.setName); return { actionName: params.actionName, setName: params.setName };`, 20*time.Second), nil
	case "action.unload":
		return scriptPlan(request.Params, `app.unloadAction(params.setName, ""); return { setName: params.setName, unloaded: true };`, 20*time.Second), nil

	case "export.png":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsPNG24(); if (params.scale) { options.horizontalScale = params.scale * 100; options.verticalScale = params.scale * 100; } doc.exportFile(file, ExportType.PNG24, options); return { outputPath: file.fsName, format: "png" };`, 30*time.Second), nil
	case "export.jpg":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsJPEG(); if (params.quality) { options.qualitySetting = params.quality; } doc.exportFile(file, ExportType.JPEG, options); return { outputPath: file.fsName, format: "jpg" };`, 30*time.Second), nil
	case "export.svg":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsSVG(); doc.exportFile(file, ExportType.SVG, options); return { outputPath: file.fsName, format: "svg" };`, 30*time.Second), nil
	case "export.pdf":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new PDFSaveOptions(); if (params.preset) { options.pDFPreset = params.preset; } doc.saveAs(file, options); return { outputPath: file.fsName, format: "pdf" };`, 30*time.Second), nil
	case "export.ai":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new IllustratorSaveOptions(); doc.saveAs(file, options); return { outputPath: file.fsName, format: "ai" };`, 30*time.Second), nil

	case "inspect.tree", "inspect.styles", "inspect.bounds", "inspect.fonts", "inspect.summary":
		return selectorPlan(illustratorinspect.SelectorFor(request.Command.Name), 20*time.Second), nil
	default:
		return executionPlan{}, fmt.Errorf("%s execution is not wired yet", request.Command.Name)
	}
}

func scriptPlan(params map[string]any, body string, timeout time.Duration) executionPlan {
	return executionPlan{
		mode:    "script",
		timeout: timeout,
		script:  wrapScript(params, body),
	}
}

func selectorPlan(selector string, timeout time.Duration) executionPlan {
	return executionPlan{
		mode:     "selector",
		selector: selector,
		timeout:  timeout,
	}
}

func wrapScript(params map[string]any, body string) string {
	payload, err := json.Marshal(params)
	if err != nil {
		payload = []byte("{}")
	}
	return fmt.Sprintf(`(function () {
var params = %s;
%s
%s
}())`, string(payload), scriptPrelude, body)
}
