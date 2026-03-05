package commands

import (
	"encoding/json"
	"fmt"
	"time"
)

const scriptPrelude = `
function ahdDoc() {
  if (app.documents.length === 0) {
    throw new Error("No active document");
  }
  return app.activeDocument;
}

function ahdDocument(id) {
  if (!id) return ahdDoc();
  for (var i = 0; i < app.documents.length; i++) {
    var doc = app.documents[i];
    if (doc.name === id) return doc;
    try {
      if (doc.fullName && (doc.fullName.fsName === id || doc.fullName.fullName === id)) return doc;
    } catch (err) {}
  }
  throw new Error("Document not found: " + id);
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

function ahdArtboardIndex(doc, id) {
  if (!id) return doc.artboards.getActiveArtboardIndex();
  for (var i = 0; i < doc.artboards.length; i++) {
    if (doc.artboards[i].name === id) return i;
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

function ahdHexByte(value) {
  var rounded = Math.round(Number(value || 0));
  if (rounded < 0) rounded = 0;
  if (rounded > 255) rounded = 255;
  var hex = rounded.toString(16).toUpperCase();
  return hex.length === 1 ? "0" + hex : hex;
}

function ahdColorValue(color) {
  if (!color) return null;
  try {
    if (color.typename === "RGBColor") {
      return "#" + ahdHexByte(color.red) + ahdHexByte(color.green) + ahdHexByte(color.blue);
    }
    if (color.typename === "GrayColor") {
      return "Gray(" + color.gray + ")";
    }
    if (color.typename === "CMYKColor") {
      return "CMYK(" + [color.cyan, color.magenta, color.yellow, color.black].join(",") + ")";
    }
    if (color.typename === "SpotColor" && color.spot) {
      return "Spot(" + color.spot.name + ")";
    }
    if (color.typename === "GradientColor" && color.gradient) {
      return "Gradient(" + color.gradient.name + ")";
    }
    return color.typename || String(color);
  } catch (err) {
    return "UnknownColor";
  }
}

function ahdBounds(item) {
  return {
    geometric: item.geometricBounds ? [item.geometricBounds[0], item.geometricBounds[1], item.geometricBounds[2], item.geometricBounds[3]] : null,
    visible: item.visibleBounds ? [item.visibleBounds[0], item.visibleBounds[1], item.visibleBounds[2], item.visibleBounds[3]] : null,
    control: item.controlBounds ? [item.controlBounds[0], item.controlBounds[1], item.controlBounds[2], item.controlBounds[3]] : null
  };
}

function ahdSelectionItems() {
  var items = [];
  if (app.documents.length === 0) return items;
  try {
    if (!app.selection) return items;
    for (var i = 0; i < app.selection.length; i++) {
      items.push(app.selection[i]);
    }
  } catch (err) {
    return items;
  }
  return items;
}

function ahdDocumentPath(doc) {
  try {
    if (doc.fullName) return doc.fullName.fsName;
  } catch (err) {}
  return null;
}

function ahdDocumentColorSpace(doc) {
  try {
    return doc.documentColorSpace === DocumentColorSpace.CMYK ? "CMYK" : "RGB";
  } catch (err) {
    return null;
  }
}

function ahdPageItemSummary(item) {
  return {
    name: item.name || "",
    typename: item.typename,
    layer: item.layer ? item.layer.name : null,
    locked: item.locked === true,
    hidden: item.hidden === true || item.visible === false,
    bounds: ahdBounds(item)
  };
}

function ahdCountMapToArray(counts) {
  var items = [];
  for (var key in counts) {
    if (counts.hasOwnProperty(key)) {
      items.push({ name: key, count: counts[key] });
    }
  }
  items.sort(function (a, b) {
    if (a.name < b.name) return -1;
    if (a.name > b.name) return 1;
    return 0;
  });
  return items;
}

function ahdArrayLikeStrings(values) {
  var items = [];
  if (!values) return items;
  try {
    if (values.length != null) {
      for (var i = 0; i < values.length; i++) {
        items.push(String(values[i]));
      }
      return items;
    }
  } catch (err) {}
  try {
    for (var key in values) {
      if (values.hasOwnProperty(key)) {
        items.push(String(values[key]));
      }
    }
  } catch (err) {}
  return items;
}

function ahdCollectFonts(doc) {
  var counts = {};
  var fontName = "";
  for (var i = 0; i < doc.textFrames.length; i++) {
    fontName = "Unknown";
    try {
      if (doc.textFrames[i].textRange && doc.textFrames[i].textRange.characterAttributes && doc.textFrames[i].textRange.characterAttributes.textFont) {
        fontName = doc.textFrames[i].textRange.characterAttributes.textFont.name || "Unknown";
      }
    } catch (err) {
      fontName = "Unknown";
    }
    if (!counts[fontName]) counts[fontName] = 0;
    counts[fontName] += 1;
  }
  return ahdCountMapToArray(counts);
}

function ahdCollectStyles(doc) {
  var fills = {};
  var strokes = {};
  var graphicStyles = [];
  var item;
  var fillKey;
  var strokeKey;
  for (var i = 0; i < doc.pageItems.length; i++) {
    item = doc.pageItems[i];
    if (item.filled) {
      fillKey = ahdColorValue(item.fillColor) || "None";
      if (!fills[fillKey]) fills[fillKey] = 0;
      fills[fillKey] += 1;
    }
    if (item.stroked) {
      strokeKey = ahdColorValue(item.strokeColor) || "None";
      if (!strokes[strokeKey]) strokes[strokeKey] = 0;
      strokes[strokeKey] += 1;
    }
  }
  for (var j = 0; j < doc.graphicStyles.length; j++) {
    graphicStyles.push(doc.graphicStyles[j].name);
  }
  graphicStyles.sort();
  return {
    fills: ahdCountMapToArray(fills),
    strokes: ahdCountMapToArray(strokes),
    graphicStyles: graphicStyles
  };
}

function ahdLayerTree(layer) {
  var node = {
    name: layer.name,
    visible: layer.visible,
    locked: layer.locked,
    items: [],
    layers: []
  };
  for (var i = 0; i < layer.pageItems.length; i++) {
    node.items.push(ahdPageItemSummary(layer.pageItems[i]));
  }
  for (var j = 0; j < layer.layers.length; j++) {
    node.layers.push(ahdLayerTree(layer.layers[j]));
  }
  return node;
}

function ahdGraphicStyle(doc, name) {
  for (var i = 0; i < doc.graphicStyles.length; i++) {
    if (doc.graphicStyles[i].name === name) return doc.graphicStyles[i];
  }
  throw new Error("Graphic style not found: " + name);
}

function ahdApplyGradient(doc, item, stops, kind) {
  if (!stops || stops.length < 2) {
    throw new Error("Gradient requires at least two stops");
  }
  var gradient = doc.gradients.add();
  gradient.name = "AHD Gradient " + new Date().getTime();
  gradient.type = kind === "radial" ? GradientType.RADIAL : GradientType.LINEAR;
  while (gradient.gradientStops.length < stops.length) {
    gradient.gradientStops.add();
  }
  for (var i = 0; i < stops.length; i++) {
    var spec = stops[i] || {};
    var stop = gradient.gradientStops[i];
    stop.color = ahdRGB(spec.color).color;
    if (spec.offset != null) {
      stop.rampPoint = spec.offset;
    } else {
      stop.rampPoint = stops.length === 1 ? 0 : (100 * i) / (stops.length - 1);
    }
  }
  var gradientColor = new GradientColor();
  gradientColor.gradient = gradient;
  item.filled = true;
  item.fillColor = gradientColor;
  return gradient.name;
}

function ahdFileExtension(path) {
  var value = String(path || "");
  var index = value.lastIndexOf(".");
  if (index === -1) return "";
  return value.substring(index + 1).toLowerCase();
}

function ahdPathWithSuffix(path, suffix, defaultExtension) {
  var value = String(path || "");
  var index = value.lastIndexOf(".");
  if (index === -1) return value + suffix + "." + defaultExtension;
  return value.substring(0, index) + suffix + value.substring(index);
}
`

func buildPlan(request Request) (executionPlan, error) {
	switch request.Command.Name {
	case "app.info":
		return scriptPlan(request.Params, `return {
  name: app.name,
  version: app.version,
  scriptingVersion: app.scriptingVersion || null,
  documents: app.documents.length,
  activeDocument: app.documents.length ? app.activeDocument.name : null,
  selectionCount: ahdSelectionItems().length,
  startupPresets: ahdArrayLikeStrings(app.startupPresetsList)
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
		return scriptPlan(request.Params, `var width = params.width || 1920; var height = params.height || 1080; var artboards = params.artboards || 1; var colorSpace = params.colorSpace === "CMYK" ? DocumentColorSpace.CMYK : DocumentColorSpace.RGB; var layoutName = params.artboardLayout || "GridByRow"; var artboardLayout = DocumentArtboardLayout[layoutName]; var artboardSpacing = params.artboardSpacing != null ? params.artboardSpacing : 20; var artboardRowsOrCols = params.artboardRowsOrCols != null ? params.artboardRowsOrCols : 1; var doc; if (params.preset) { var preset = new DocumentPreset(); preset.colorMode = colorSpace; preset.width = width; preset.height = height; preset.numArtboards = artboards; preset.artboardLayout = artboardLayout; preset.artboardSpacing = artboardSpacing; preset.artboardRowsOrCols = artboardRowsOrCols; doc = app.documents.addDocument(params.preset, preset, false); } else { doc = app.documents.add(colorSpace, width, height, artboards, artboardLayout, artboardSpacing, artboardRowsOrCols); } return { name: doc.name, width: width, height: height, colorSpace: params.colorSpace || "RGB", artboards: doc.artboards.length, preset: params.preset || null, artboardLayout: layoutName, artboardSpacing: artboardSpacing, artboardRowsOrCols: artboardRowsOrCols };`, 30*time.Second), nil
	case "document.open":
		return scriptPlan(request.Params, `var file = new File(params.filePath); var doc = app.open(file); return { name: doc.name, path: file.fsName };`, 20*time.Second), nil
	case "document.save":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); doc.save(); return { name: doc.name, path: ahdDocumentPath(doc), saved: true };`, 20*time.Second), nil
	case "document.save_as":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var file = new File(params.filePath); var format = (params.format || ahdFileExtension(file.name) || "ai").toLowerCase(); if (format === "pdf") { var pdfOptions = new PDFSaveOptions(); doc.saveAs(file, pdfOptions); } else { var aiOptions = new IllustratorSaveOptions(); doc.saveAs(file, aiOptions); } return { name: doc.name, path: file.fsName, format: format };`, 20*time.Second), nil
	case "document.close":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var name = doc.name; doc.close(params.save ? SaveOptions.SAVECHANGES : SaveOptions.DONOTSAVECHANGES); return { closed: true, name: name };`, 20*time.Second), nil
	case "document.list":
		return scriptPlan(request.Params, `var docs = []; for (var i = 0; i < app.documents.length; i++) { var doc = app.documents[i]; docs.push({ index: i, name: doc.name, path: ahdDocumentPath(doc), artboards: doc.artboards.length, pageItems: doc.pageItems.length, colorSpace: ahdDocumentColorSpace(doc) }); } return docs;`, 10*time.Second), nil
	case "document.info":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); return { name: doc.name, path: ahdDocumentPath(doc), artboards: doc.artboards.length, layers: doc.layers.length, pageItems: doc.pageItems.length, colorSpace: ahdDocumentColorSpace(doc), width: doc.width, height: doc.height };`, 10*time.Second), nil

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
		return scriptPlan(request.Params, `var selection = ahdSelectionItems(); var items = []; for (var i = 0; i < selection.length; i++) { items.push({ index: i, name: selection[i].name || "", typename: selection[i].typename }); } return items;`, 10*time.Second), nil
	case "selection.clear":
		return scriptPlan(request.Params, `app.selection = null; return { cleared: true };`, 10*time.Second), nil
	case "selection.set_by_ids":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var selection = []; for (var i = 0; i < params.ids.length; i++) { selection.push(ahdPageItem(doc, params.ids[i])); } app.selection = selection; return { selected: selection.length };`, 20*time.Second), nil
	case "selection.select_by_name":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var selection = []; for (var i = 0; i < doc.pageItems.length; i++) { var item = doc.pageItems[i]; if (params.partial ? item.name.indexOf(params.name) !== -1 : item.name === params.name) { selection.push(item); } } app.selection = selection; return { selected: selection.length, name: params.name };`, 20*time.Second), nil

	case "path.create_rect":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var radius = params.cornerRadius || 0; var item = radius > 0 ? layer.pathItems.roundedRectangle(params.top, params.left, params.width, params.height, radius, radius) : layer.pathItems.rectangle(params.top, params.left, params.width, params.height); item.name = params.name; return { name: item.name, width: item.width, height: item.height, cornerRadius: radius };`, 20*time.Second), nil
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
	case "appearance.set_gradient":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var gradientName = ahdApplyGradient(doc, item, params.stops, params.type || "linear"); return { itemId: params.itemId, gradient: gradientName, stops: params.stops.length, type: params.type || "linear" };`, 20*time.Second), nil
	case "appearance.apply_graphic_style":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var style = ahdGraphicStyle(doc, params.styleName); style.applyTo(item); return { itemId: params.itemId, styleName: style.name };`, 20*time.Second), nil

	case "action.load":
		return scriptPlan(request.Params, `app.loadAction(new File(params.filePath)); return { filePath: params.filePath, loaded: true };`, 20*time.Second), nil
	case "action.run":
		return scriptPlan(request.Params, `app.doScript(params.actionName, params.setName); return { actionName: params.actionName, setName: params.setName };`, 20*time.Second), nil
	case "action.unload":
		return scriptPlan(request.Params, `app.unloadAction(params.setName, ""); return { setName: params.setName, unloaded: true };`, 20*time.Second), nil

	case "export.png":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsPNG24(); if (params.scale) { options.horizontalScale = params.scale * 100; options.verticalScale = params.scale * 100; } if (params.artboardId) { doc.artboards.setActiveArtboardIndex(ahdArtboardIndex(doc, params.artboardId)); options.artBoardClipping = true; } doc.exportFile(file, ExportType.PNG24, options); return { outputPath: file.fsName, format: "png", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.jpg":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsJPEG(); if (params.scale) { options.horizontalScale = params.scale * 100; options.verticalScale = params.scale * 100; } if (params.quality) { options.qualitySetting = params.quality; } if (params.artboardId) { doc.artboards.setActiveArtboardIndex(ahdArtboardIndex(doc, params.artboardId)); options.artBoardClipping = true; } doc.exportFile(file, ExportType.JPEG, options); return { outputPath: file.fsName, format: "jpg", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.svg":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsSVG(); var outputPath = file.fsName; if (params.artboardId) { options.saveMultipleArtboards = true; options.artboardRange = String(ahdArtboardIndex(doc, params.artboardId) + 1); outputPath = ahdPathWithSuffix(file.fsName, "_" + params.artboardId, "svg"); } doc.exportFile(file, ExportType.SVG, options); return { outputPath: outputPath, format: "svg", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.pdf":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new PDFSaveOptions(); if (params.preset) { options.pDFPreset = params.preset; } doc.saveAs(file, options); return { outputPath: file.fsName, format: "pdf" };`, 30*time.Second), nil
	case "export.ai":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new IllustratorSaveOptions(); doc.saveAs(file, options); return { outputPath: file.fsName, format: "ai" };`, 30*time.Second), nil

	case "inspect.tree":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var tree = []; for (var i = 0; i < doc.layers.length; i++) { tree.push(ahdLayerTree(doc.layers[i])); } return { document: doc.name, layers: tree };`, 20*time.Second), nil
	case "inspect.styles":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return ahdCollectStyles(doc);`, 20*time.Second), nil
	case "inspect.bounds":
		return scriptPlan(request.Params, `var selection = ahdSelectionItems(); var items = []; for (var i = 0; i < selection.length; i++) { items.push(ahdPageItemSummary(selection[i])); } return { count: items.length, items: items };`, 20*time.Second), nil
	case "inspect.fonts":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { fonts: ahdCollectFonts(doc) };`, 20*time.Second), nil
	case "inspect.summary":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var styles = ahdCollectStyles(doc); return { name: doc.name, artboards: doc.artboards.length, layers: doc.layers.length, pageItems: doc.pageItems.length, selectionCount: ahdSelectionItems().length, fonts: ahdCollectFonts(doc), fills: styles.fills, strokes: styles.strokes, graphicStyles: styles.graphicStyles };`, 20*time.Second), nil
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
