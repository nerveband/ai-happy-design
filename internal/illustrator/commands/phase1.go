package commands

import (
	"fmt"
	"time"
)

const phase1ScriptPrelude = `
function ahdView(doc) {
  if (!doc.views || doc.views.length === 0) {
    throw new Error("No active view");
  }
  return doc.views[0];
}

function ahdViewSummary(doc) {
  var view = ahdView(doc);
  return {
    bounds: view.bounds ? [view.bounds[0], view.bounds[1], view.bounds[2], view.bounds[3]] : null,
    centerPoint: view.centerPoint ? [view.centerPoint[0], view.centerPoint[1]] : null,
    rotateAngle: view.rotateAngle,
    screenMode: String(view.screenMode),
    zoom: view.zoom,
    documentScreenMode: doc.getScreenMode ? String(doc.getScreenMode()) : null,
    documentViewMode: doc.getViewMode ? String(doc.getViewMode()) : null,
    rulerVisible: doc.isRulerVisible ? doc.isRulerVisible() : null,
    transparencyGridVisible: doc.isTransparencyGridVisible ? doc.isTransparencyGridVisible() : null
  };
}

function ahdMatrixValue(spec) {
  var matrix = app.getIdentityMatrix();
  if (!spec) return matrix;
  if (spec.mValueA != null) matrix.mValueA = spec.mValueA;
  if (spec.mValueB != null) matrix.mValueB = spec.mValueB;
  if (spec.mValueC != null) matrix.mValueC = spec.mValueC;
  if (spec.mValueD != null) matrix.mValueD = spec.mValueD;
  if (spec.mValueTX != null) matrix.mValueTX = spec.mValueTX;
  if (spec.mValueTY != null) matrix.mValueTY = spec.mValueTY;
  return matrix;
}

function ahdMatrixSummary(matrix) {
  return {
    mValueA: matrix.mValueA,
    mValueB: matrix.mValueB,
    mValueC: matrix.mValueC,
    mValueD: matrix.mValueD,
    mValueTX: matrix.mValueTX,
    mValueTY: matrix.mValueTY
  };
}

function ahdRectValue(spec) {
  if (!spec) return null;
  return [spec.left, spec.top, spec.right, spec.bottom];
}

function ahdRGBColor(hex) {
  return ahdRGB(hex).color;
}

function ahdCMYKColor(spec) {
  var color = new CMYKColor();
  color.cyan = spec.cyan || 0;
  color.magenta = spec.magenta || 0;
  color.yellow = spec.yellow || 0;
  color.black = spec.black || 0;
  return color;
}

function ahdColorFromMode(spec) {
  if (!spec) {
    throw new Error("Color specification is required");
  }
  if (spec.colorMode === "CMYK") {
    return ahdCMYKColor(spec);
  }
  if (!spec.hex) {
    throw new Error("RGB colors require a hex field");
  }
  return ahdRGBColor(spec.hex);
}

function ahdSwatch(doc, id) {
  for (var i = 0; i < doc.swatches.length; i++) {
    var swatch = doc.swatches[i];
    if (swatch.name === id) return swatch;
  }
  throw new Error("Swatch not found: " + id);
}

function ahdSpot(doc, id) {
  for (var i = 0; i < doc.spots.length; i++) {
    var spot = doc.spots[i];
    if (spot.name === id) return spot;
  }
  ahdTryRedraw();
  for (var j = 0; j < doc.spots.length; j++) {
    var retrySpot = doc.spots[j];
    if (retrySpot.name === id) return retrySpot;
  }
  throw new Error("Spot not found: " + id);
}

function ahdCharacterStyle(doc, id) {
  for (var i = 0; i < doc.characterStyles.length; i++) {
    var style = doc.characterStyles[i];
    if (style.name === id) return style;
  }
  throw new Error("Character style not found: " + id);
}

function ahdParagraphStyle(doc, id) {
  for (var i = 0; i < doc.paragraphStyles.length; i++) {
    var style = doc.paragraphStyles[i];
    if (style.name === id) return style;
  }
  throw new Error("Paragraph style not found: " + id);
}

function ahdSymbol(doc, id) {
  for (var i = 0; i < doc.symbols.length; i++) {
    var symbol = doc.symbols[i];
    if (symbol.name === id) return symbol;
  }
  throw new Error("Symbol not found: " + id);
}

function ahdSymbolItem(doc, id) {
  for (var i = 0; i < doc.symbolItems.length; i++) {
    var item = doc.symbolItems[i];
    if (item.name === id || item.note === id) return item;
  }
  ahdTryRedraw();
  for (var j = 0; j < doc.symbolItems.length; j++) {
    var retryItem = doc.symbolItems[j];
    if (retryItem.name === id || retryItem.note === id) return retryItem;
  }
  throw new Error("Symbol item not found: " + id);
}

function ahdPlacedItem(doc, id) {
  for (var i = 0; i < doc.placedItems.length; i++) {
    var item = doc.placedItems[i];
    if (item.name === id || item.note === id) return item;
  }
  ahdTryRedraw();
  for (var j = 0; j < doc.placedItems.length; j++) {
    var retryItem = doc.placedItems[j];
    if (retryItem.name === id || retryItem.note === id) return retryItem;
  }
  throw new Error("Placed item not found: " + id);
}

function ahdRasterItem(doc, id) {
  for (var i = 0; i < doc.rasterItems.length; i++) {
    var item = doc.rasterItems[i];
    if (item.name === id || item.note === id) return item;
  }
  ahdTryRedraw();
  for (var j = 0; j < doc.rasterItems.length; j++) {
    var retryItem = doc.rasterItems[j];
    if (retryItem.name === id || retryItem.note === id) return retryItem;
  }
  throw new Error("Raster item not found: " + id);
}

function ahdDataSet(doc, id) {
  for (var i = 0; i < doc.dataSets.length; i++) {
    var item = doc.dataSets[i];
    if (item.name === id) return item;
  }
  throw new Error("Dataset not found: " + id);
}

function ahdVariable(doc, id) {
  for (var i = 0; i < doc.variables.length; i++) {
    var item = doc.variables[i];
    if (item.name === id) return item;
  }
  throw new Error("Variable not found: " + id);
}

function ahdDataSetSummary(item, index) {
  return {
    index: index,
    name: item.name || ""
  };
}

function ahdVariableSummary(item, index) {
  return {
    index: index,
    name: item.name || "",
    kind: item.kind ? String(item.kind) : null,
    pageItems: item.pageItems ? item.pageItems.length : 0
  };
}

function ahdRepeatItem(collection, id, label) {
  for (var i = 0; i < collection.length; i++) {
    var item = collection[i];
    if (item.name === id) return item;
  }
  throw new Error(label + " not found: " + id);
}

function ahdRepeatSummary(item, kind) {
  return {
    name: item.name || "",
    kind: kind,
    typename: item.typename
  };
}

function ahdSwatchSummary(swatch) {
  return {
    name: swatch.name,
    color: ahdColorValue(swatch.color)
  };
}

function ahdSpotSummary(spot) {
  return {
    name: spot.name,
    color: ahdColorValue(spot.color),
    colorType: String(spot.colorType),
    spotKind: String(spot.spotKind)
  };
}

function ahdSymbolSummary(symbol) {
  return {
    name: symbol.name || "",
    typename: symbol.typename
  };
}

function ahdSymbolItemSummary(item) {
  return {
    name: item.name || "",
    symbol: item.symbol ? item.symbol.name : null,
    width: item.width,
    height: item.height,
    position: item.position ? [item.position[0], item.position[1]] : null
  };
}

function ahdPlacedSummary(item) {
  var filePath = null;
  try {
    if (item.file) filePath = item.file.fsName;
  } catch (err) {}
  return {
    name: item.name || "",
    filePath: filePath,
    bounds: ahdBounds(item),
    position: item.position ? [item.position[0], item.position[1]] : null,
    typename: item.typename
  };
}

function ahdRasterSummary(item) {
  var filePath = null;
  try {
    if (item.file) filePath = item.file.fsName;
  } catch (err) {}
  return {
    name: item.name || "",
    filePath: filePath,
    embedded: item.embedded === true,
    imageColorSpace: item.imageColorSpace ? String(item.imageColorSpace) : null,
    status: item.status ? String(item.status) : null,
    bounds: ahdBounds(item)
  };
}

function ahdTracingSummary(pluginItem) {
  ahdTryRedraw();
  if (!pluginItem.tracing) {
    return {
      pluginItemName: pluginItem.name || "",
      tracing: false
    };
  }
  return {
    pluginItemName: pluginItem.name || "",
    tracing: true,
    typename: pluginItem.typename || null,
    tracingMode: pluginItem.tracing.tracingOptions && pluginItem.tracing.tracingOptions.tracingMode ? String(pluginItem.tracing.tracingOptions.tracingMode) : null
  };
}

function ahdTracingPlugin(doc, id) {
  var item = ahdPageItem(doc, id);
  if (!item.tracing) {
    throw new Error("Tracing plugin item not found: " + id);
  }
  return item;
}

function ahdScreenModeValue(name) {
  if (!ScreenMode || ScreenMode[name] == null) {
    throw new Error("Unknown screen mode: " + name);
  }
  return ScreenMode[name];
}

function ahdPerspectivePlaneValue(name) {
  switch (name) {
    case "GRIDLEFTPLANETYPE":
      if (typeof GRIDLEFTPLANETYPE !== "undefined") return GRIDLEFTPLANETYPE;
      break;
    case "GRIDRIGHTPLANETYPE":
      if (typeof GRIDRIGHTPLANETYPE !== "undefined") return GRIDRIGHTPLANETYPE;
      break;
    case "GRIDFLOORPLANETYPE":
      if (typeof GRIDFLOORPLANETYPE !== "undefined") return GRIDFLOORPLANETYPE;
      break;
    case "INVALIDGRIDPLANETYPE":
      if (typeof INVALIDGRIDPLANETYPE !== "undefined") return INVALIDGRIDPLANETYPE;
      break;
  }
  if (PerspectiveGridPlaneType && PerspectiveGridPlaneType[name] != null) {
    return PerspectiveGridPlaneType[name];
  }
  throw new Error("Unknown perspective plane: " + name);
}

function ahdPerspectiveGridTypeValue(name) {
  switch (name) {
    case "OnePointPerspectiveGridType": return OnePointPerspectiveGridType;
    case "TwoPointPerspectiveGridType": return TwoPointPerspectiveGridType;
    case "ThreePointPerspectiveGridType": return ThreePointPerspectiveGridType;
    case "InvalidPerspectiveGridType": return InvalidPerspectiveGridType;
  }
  throw new Error("Unknown perspective grid type: " + name);
}

function ahdSymbolRegistrationPointValue(name) {
  if (SymbolRegistrationPoint && SymbolRegistrationPoint[name] != null) {
    return SymbolRegistrationPoint[name];
  }
  throw new Error("Unknown symbol registration point: " + name);
}

function ahdImageColorSpaceValue(name) {
  if (ImageColorSpace && ImageColorSpace[name] != null) {
    return ImageColorSpace[name];
  }
  throw new Error("Unknown image color space: " + name);
}

function ahdColorConvertPurposeValue(name) {
  if (ColorConvertPurpose && ColorConvertPurpose[name] != null) {
    return ColorConvertPurpose[name];
  }
  throw new Error("Unknown color conversion purpose: " + name);
}

function ahdVariableKindValue(name) {
  if (VariableKind && VariableKind[name] != null) {
    return VariableKind[name];
  }
  throw new Error("Unknown variable kind: " + name);
}

function ahdDocumentPresetTypeValue(name) {
  if (DocumentPresetType && DocumentPresetType[name] != null) {
    return DocumentPresetType[name];
  }
  throw new Error("Unknown document preset type: " + name);
}

function ahdLibraryTypeValue(name) {
  var libraryTypeName = name || "IllustratorArtwork";
  if (typeof this[libraryTypeName] !== "undefined") {
    return this[libraryTypeName];
  }
  if (LibraryType && LibraryType[libraryTypeName] != null) {
    return LibraryType[libraryTypeName];
  }
  throw new Error("Unknown library type: " + libraryTypeName);
}

function ahdDocumentLayoutStyleValue(name) {
  if (!name) return DocumentLayoutStyle.CASCADE;
  if (DocumentLayoutStyle && DocumentLayoutStyle[name] != null) {
    return DocumentLayoutStyle[name];
  }
  throw new Error("Unknown document layout style: " + name);
}

function ahdTransformationValue(name) {
  if (!name) return undefined;
  if (Transformation && Transformation[name] != null) {
    return Transformation[name];
  }
  throw new Error("Unknown transformation anchor: " + name);
}

function ahdZOrderMethodValue(name) {
  if (ZOrderMethod && ZOrderMethod[name] != null) {
    return ZOrderMethod[name];
  }
  throw new Error("Unknown z-order method: " + name);
}

function ahdPlacementValue(name) {
  switch (name || "inside") {
    case "inside":
      return ElementPlacement.INSIDE;
    case "before":
      return ElementPlacement.PLACEBEFORE;
    case "after":
      return ElementPlacement.PLACEAFTER;
    case "at_beginning":
      return ElementPlacement.PLACEATBEGINNING;
    case "at_end":
      return ElementPlacement.PLACEATEND;
  }
  throw new Error("Unknown placement mode: " + name);
}

function ahdPlacementTarget(doc, id) {
  if (!id) return doc;
  try {
    return ahdLayer(doc, id);
  } catch (layerErr) {}
  return ahdPageItem(doc, id);
}

function ahdGraphicStyle(doc, id) {
  for (var i = 0; i < doc.graphicStyles.length; i++) {
    var style = doc.graphicStyles[i];
    if (style.name === id) return style;
  }
  throw new Error("Graphic style not found: " + id);
}

function ahdExportTypeValue(name) {
  switch (name) {
    case "gif": return ExportType.GIF;
    case "jpeg": return ExportType.JPEG;
    case "photoshop": return ExportType.PHOTOSHOP != null ? ExportType.PHOTOSHOP : ExportType.Photoshop;
    case "png8": return ExportType.PNG8;
    case "png24": return ExportType.PNG24;
    case "svg": return ExportType.SVG;
    case "tiff": return ExportType.TIFF;
    case "webp": return ExportType.WEBP;
    case "autocad": return ExportType.AutoCAD != null ? ExportType.AutoCAD : ExportType.AUTOCAD;
  }
  throw new Error("Unknown export type: " + name);
}

function ahdApplyArtboardExport(doc, artboardId, options) {
  if (!artboardId) return null;
  var index = ahdArtboardIndex(doc, artboardId);
  doc.artboards.setActiveArtboardIndex(index);
  if (options.artBoardClipping != null) {
    options.artBoardClipping = true;
  }
  if (options.artboardRange != null) {
    options.saveMultipleArtboards = true;
    options.artboardRange = String(index + 1);
  }
  return index;
}

function ahdImageCaptureOptions(spec) {
  var options = new ImageCaptureOptions();
  if (!spec) return options;
  if (spec.antiAliasing != null) options.antiAliasing = spec.antiAliasing;
  if (spec.matte != null) options.matte = spec.matte;
  if (spec.matteColor) options.matteColor = ahdRGBColor(spec.matteColor);
  if (spec.resolution != null) options.resolution = spec.resolution;
  if (spec.transparency != null) options.transparency = spec.transparency;
  return options;
}

function ahdSetPrintJobOptions(options, spec) {
  if (!spec) return;
  var job = new PrintJobOptions();
  if (spec.name) job.name = spec.name;
  if (spec.copies != null) job.copies = spec.copies;
  if (spec.collate != null) job.collate = spec.collate;
  if (spec.printAllArtboards != null) job.printAllArtboards = spec.printAllArtboards;
  if (spec.artboardRange) job.artboardRange = spec.artboardRange;
  if (spec.printAsBitmap != null) job.printAsBitmap = spec.printAsBitmap;
  if (spec.reversePages != null) job.reversePages = spec.reversePages;
  options.jobOptions = job;
}

function ahdRepeatConfigValue(kind, spec) {
  var config;
  switch (kind) {
    case "grid":
      config = new GridRepeatConfig();
      break;
    case "radial":
      config = new RadialRepeatConfig();
      break;
    case "symmetry":
      config = new SymmetryRepeatConfig();
      break;
    default:
      config = {};
      break;
  }
  if (!spec) return config;
  for (var key in spec) {
    if (spec.hasOwnProperty(key) && spec[key] != null) {
      config[key] = spec[key];
    }
  }
  return config;
}
`

func buildPhase1Plan(request Request) (executionPlan, error) {
	switch request.Command.Name {
	case "app.beep":
		return scriptPlan(request.Params, `app.beep(); return { beeped: true };`, 10*time.Second), nil
	case "app.copy":
		return scriptPlanActivated(request.Params, `app.copy(); return { copied: true, selectionCount: ahdSelectionItems().length };`, 10*time.Second), nil
	case "app.cut":
		return scriptPlanActivated(request.Params, `app.cut(); return { cut: true };`, 10*time.Second), nil
	case "app.paste":
		return scriptPlanActivated(request.Params, `app.paste(); return { pasted: true, selectionCount: ahdSelectionItems().length };`, 10*time.Second), nil
	case "app.undo":
		return scriptPlan(request.Params, `app.undo(); return { undone: true };`, 10*time.Second), nil
	case "app.redo":
		return scriptPlan(request.Params, `app.redo(); return { redone: true };`, 10*time.Second), nil
	case "app.redraw":
		return scriptPlan(request.Params, `var documentCount = app.documents.length; if (documentCount > 0) { app.redraw(); } return { redrawn: documentCount > 0, documents: documentCount };`, 10*time.Second), nil
	case "app.convert_sample_color":
		return scriptPlan(request.Params, `var converted = app.convertSampleColor(ahdImageColorSpaceValue(params.sourceColorSpace), params.sourceColor, ahdImageColorSpaceValue(params.destColorSpace), ahdColorConvertPurposeValue(params.purpose || "defaultpurpose"), !!params.sourceHasAlpha, !!params.destHasAlpha); return { sourceColorSpace: params.sourceColorSpace, destColorSpace: params.destColorSpace, sourceColor: params.sourceColor, converted: converted };`, 10*time.Second), nil
	case "app.translate_placeholder_text":
		return scriptPlan(request.Params, `return { source: params.text, translated: app.translatePlaceholderText(params.text) };`, 10*time.Second), nil
	case "app.preset_lists":
		return scriptPlan(request.Params, `var printers = []; var ppdFiles = []; var printerEnumeration = "ok"; var ppdEnumeration = "ok"; try { var printerList = app.printerList || []; for (var i = 0; i < printerList.length; i++) { printers.push(printerList[i].name || String(printerList[i])); } printers.sort(); } catch (printerErr) { printerEnumeration = "unavailable"; } try { var ppdList = app.PPDFileList || []; for (var j = 0; j < ppdList.length; j++) { ppdFiles.push(ppdList[j].name || String(ppdList[j])); } ppdFiles.sort(); } catch (ppdErr) { ppdEnumeration = "unavailable"; } return { pdfPresets: ahdArrayLikeStrings(app.PDFPresetsList), printPresets: ahdArrayLikeStrings(app.printPresetsList), startupPresets: ahdArrayLikeStrings(app.startupPresetsList), tracingPresets: ahdArrayLikeStrings(app.tracingPresetList), colorSettings: ahdArrayLikeStrings(app.colorSettingsList), printers: printers, printerEnumeration: printerEnumeration, ppdFiles: ppdFiles, ppdEnumeration: ppdEnumeration };`, 15*time.Second), nil
	case "app.get_preset_file":
		return scriptPlan(request.Params, `var file = app.getPresetFileOfType(ahdDocumentPresetTypeValue(params.presetType)); return { presetType: params.presetType, filePath: file ? file.fsName : null };`, 10*time.Second), nil
	case "app.get_preset_settings":
		return scriptPlan(request.Params, `var settings = app.getPresetSettings(params.preset); return { preset: params.preset, settings: settings };`, 10*time.Second), nil
	case "app.load_color_settings":
		return scriptPlan(request.Params, `if (params.disable === true) { app.loadColorSettings(""); return { disabled: true, filePath: "" }; } if (!params.filePath) { throw new Error("filePath is required unless disable=true"); } var file = new File(params.filePath); app.loadColorSettings(file); return { disabled: false, filePath: file.fsName, loaded: true };`, 10*time.Second), nil
	case "app.show_presets":
		return scriptPlan(request.Params, `var file = new File(params.filePath); app.showPresets(file); return { filePath: file.fsName, shown: true };`, 10*time.Second), nil

	case "workspace.save":
		return scriptPlan(request.Params, `app.saveWorkspace(params.name); return { name: params.name, saved: true };`, 10*time.Second), nil
	case "workspace.switch":
		return scriptPlan(request.Params, `app.switchWorkspace(params.name); return { name: params.name, switched: true };`, 10*time.Second), nil
	case "workspace.reset":
		return scriptPlan(request.Params, `app.resetWorkspace(); return { reset: true };`, 10*time.Second), nil
	case "workspace.delete":
		return scriptPlan(request.Params, `app.deleteWorkspace(params.name); return { name: params.name, deleted: true };`, 10*time.Second), nil

	case "document.activate":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); doc.activate(); return { documentId: doc.name, activated: true };`, 10*time.Second), nil
	case "document.arrange":
		return scriptPlanActivated(request.Params, `var layoutStyle = params.layoutStyle || "CASCADE"; var menuCommand = layoutStyle === "CASCADE" ? "cascade" : "tile"; app.executeMenuCommand(menuCommand); return { layoutStyle: layoutStyle, menuCommand: menuCommand, arranged: true, documents: app.documents.length };`, 15*time.Second), nil
	case "document.write_as_library":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var file = new File(params.filePath); var libraryType = params.libraryType || "IllustratorArtwork"; var ok = doc.writeAsLibrary(file, ahdLibraryTypeValue(libraryType)); return { documentId: doc.name, filePath: file.fsName, libraryType: libraryType, written: ok === true };`, 20*time.Second), nil
	case "document.export_pdf_preset":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var file = new File(params.filePath); var snapshot = ahdFolderSnapshot(file.parent, [".joboptions"]); doc.exportPDFPreset(file); return { documentId: doc.name, filePath: ahdResolveCreatedSibling(file, [".joboptions"], snapshot), exported: true, presetType: "pdf" };`, 20*time.Second), nil
	case "document.import_pdf_preset":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var file = new File(params.filePath); if (params.replacingPreset === true) { doc.importPDFPreset(file, true); } else { doc.importPDFPreset(file); } return { documentId: doc.name, filePath: file.fsName, imported: true, replacingPreset: params.replacingPreset === true, presetType: "pdf" };`, 20*time.Second), nil
	case "document.export_print_preset":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var file = new File(params.filePath); var snapshot = ahdFolderSnapshot(file.parent, [".prst"]); doc.exportPrintPreset(file); return { documentId: doc.name, filePath: ahdResolveCreatedSibling(file, [".prst"], snapshot), exported: true, presetType: "print" };`, 20*time.Second), nil
	case "document.import_print_preset":
		return scriptPlan(request.Params, `var doc = ahdDocument(params.documentId); var file = new File(params.filePath); doc.importPrintPreset(params.printPreset, file); return { documentId: doc.name, filePath: file.fsName, imported: true, printPreset: params.printPreset, presetType: "print" };`, 20*time.Second), nil

	case "preference.get":
		return scriptPlan(request.Params, `var value = null; switch (params.valueType) { case "boolean": value = app.preferences.getBooleanPreference(params.key); break; case "integer": value = app.preferences.getIntegerPreference(params.key); break; case "real": value = app.preferences.getRealPreference(params.key); break; default: value = app.preferences.getStringPreference(params.key); break; } return { key: params.key, valueType: params.valueType, value: value };`, 10*time.Second), nil
	case "preference.set":
		return scriptPlan(request.Params, `var outValue = null; switch (params.valueType) { case "boolean": outValue = !!params.value.boolean; app.preferences.setBooleanPreference(params.key, outValue); break; case "integer": outValue = Math.round(params.value.number); app.preferences.setIntegerPreference(params.key, outValue); break; case "real": outValue = Number(params.value.number); app.preferences.setRealPreference(params.key, outValue); break; default: outValue = params.value.string || ""; app.preferences.setStringPreference(params.key, outValue); break; } return { key: params.key, valueType: params.valueType, value: outValue, set: true };`, 10*time.Second), nil
	case "preference.delete":
		return scriptPlan(request.Params, `app.preferences.removePreference(params.key); return { key: params.key, deleted: true };`, 10*time.Second), nil

	case "view.info":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return ahdViewSummary(doc);`, 10*time.Second), nil
	case "view.set_screen_mode":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var view = ahdView(doc); view.screenMode = ahdScreenModeValue(params.mode); return ahdViewSummary(doc);`, 10*time.Second), nil
	case "view.set_zoom":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var view = ahdView(doc); view.zoom = params.zoom; return ahdViewSummary(doc);`, 10*time.Second), nil
	case "view.set_ruler_visibility":
		return scriptPlanActivated(request.Params, `var doc = ahdDoc(); var visible = doc.isRulerVisible ? doc.isRulerVisible() : null; if (visible == null) { throw new Error("Ruler visibility is unavailable in this host"); } if (visible !== params.visible) { app.executeMenuCommand("ruler"); } return ahdViewSummary(doc);`, 15*time.Second), nil
	case "view.set_transparency_grid_visibility":
		return scriptPlanActivated(request.Params, `var doc = ahdDoc(); var visible = doc.isTransparencyGridVisible ? doc.isTransparencyGridVisible() : null; if (visible == null) { throw new Error("Transparency grid visibility is unavailable in this host"); } if (visible !== params.visible) { app.executeMenuCommand("TransparencyGrid Menu Item"); } return ahdViewSummary(doc);`, 15*time.Second), nil
	case "view.set_center":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var view = ahdView(doc); view.centerPoint = [params.x, params.y]; return ahdViewSummary(doc);`, 10*time.Second), nil
	case "view.rotate":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var view = ahdView(doc); view.rotateAngle = params.angle; return ahdViewSummary(doc);`, 10*time.Second), nil

	case "matrix.identity":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.getIdentityMatrix());`, 10*time.Second), nil
	case "matrix.rotation":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.getRotationMatrix(params.angle));`, 10*time.Second), nil
	case "matrix.scale":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.getScaleMatrix(params.scaleX, params.scaleY != null ? params.scaleY : params.scaleX));`, 10*time.Second), nil
	case "matrix.translation":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.getTranslationMatrix(params.deltaX, params.deltaY != null ? params.deltaY : 0));`, 10*time.Second), nil
	case "matrix.concatenate":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.concatenateMatrix(ahdMatrixValue(params.matrix), ahdMatrixValue(params.secondMatrix)));`, 10*time.Second), nil
	case "matrix.concatenate_rotation":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.concatenateRotationMatrix(ahdMatrixValue(params.matrix), params.angle));`, 10*time.Second), nil
	case "matrix.concatenate_scale":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.concatenateScaleMatrix(ahdMatrixValue(params.matrix), params.scaleX, params.scaleY != null ? params.scaleY : params.scaleX));`, 10*time.Second), nil
	case "matrix.concatenate_translation":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.concatenateTranslationMatrix(ahdMatrixValue(params.matrix), params.deltaX, params.deltaY != null ? params.deltaY : 0));`, 10*time.Second), nil
	case "matrix.invert":
		return scriptPlan(request.Params, `return ahdMatrixSummary(app.invertMatrix(ahdMatrixValue(params.matrix)));`, 10*time.Second), nil
	case "matrix.compare":
		return scriptPlan(request.Params, `return { equal: app.isEqualMatrix(ahdMatrixValue(params.matrix), ahdMatrixValue(params.secondMatrix)) === true };`, 10*time.Second), nil
	case "matrix.singular":
		return scriptPlan(request.Params, `return { singular: app.isSingularMatrix(ahdMatrixValue(params.matrix)) === true };`, 10*time.Second), nil

	case "perspective.show":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { shown: doc.showPerspectiveGrid() === true };`, 10*time.Second), nil
	case "perspective.hide":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { hidden: doc.hidePerspectiveGrid() === true };`, 10*time.Second), nil
	case "perspective.get_active_plane":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { plane: String(doc.getPerspectiveActivePlane()) };`, 10*time.Second), nil
	case "perspective.set_active_plane":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { plane: params.plane, active: doc.setActivePlane(ahdPerspectivePlaneValue(params.plane)) === true };`, 10*time.Second), nil
	case "perspective.select_preset":
		return scriptPlan(request.Params, `var doc = ahdDoc(); return { gridType: params.gridType, presetName: params.presetName, selected: doc.selectPerspectivePreset(ahdPerspectiveGridTypeValue(params.gridType), params.presetName) === true };`, 10*time.Second), nil
	case "perspective.import_preset":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); if (params.presetName) { doc.importPerspectiveGridPreset(file, params.presetName); } else { doc.importPerspectiveGridPreset(file); } return { filePath: file.fsName, presetName: params.presetName || null, imported: true };`, 20*time.Second), nil
	case "perspective.export_preset":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.exportPerspectiveGridPreset(file); return { filePath: file.fsName, exported: true };`, 20*time.Second), nil

	case "artboard.delete":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var index = ahdArtboardIndex(doc, params.artboardId); var artboard = doc.artboards[index]; var name = artboard.name || params.artboardId; artboard.remove(); return { artboardId: params.artboardId, name: name, deleted: true, remaining: doc.artboards.length };`, 20*time.Second), nil
	case "artboard.rearrange":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var fn = doc.rearrangeArtboards || doc.rearrangeArboards; if (!fn) throw new Error("Artboard rearrangement is unavailable"); var layoutName = params.artboardLayout || "GridByRow"; var rowsOrCols = params.artboardRowsOrCols != null ? params.artboardRowsOrCols : 1; var spacing = params.artboardSpacing != null ? params.artboardSpacing : 20; var moveArtwork = params.moveArtwork != null ? params.moveArtwork : true; var ok = fn.call(doc, DocumentArtboardLayout[layoutName], rowsOrCols, spacing, moveArtwork); return { rearranged: ok === true, artboardLayout: layoutName, artboardRowsOrCols: rowsOrCols, artboardSpacing: spacing, moveArtwork: moveArtwork };`, 20*time.Second), nil

	case "selection.select_active_artboard_objects":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var ok = doc.selectObjectsOnActiveArtboard(); return { selected: ok === true, selectionCount: ahdSelectionItems().length };`, 10*time.Second), nil

	case "page_item.remove":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var name = item.name || params.itemId; item.remove(); return { itemId: params.itemId, name: name, deleted: true };`, 20*time.Second), nil
	case "page_item.duplicate":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var target = params.targetId ? ahdPlacementTarget(doc, params.targetId) : undefined; var duplicate = target ? item.duplicate(target, ahdPlacementValue(params.placement || "inside")) : item.duplicate(); var duplicateId = params.name || ((item.name || params.itemId) + " Copy"); ahdStampIdentifier(duplicate, duplicateId); ahdTryRedraw(); return ahdPageItemSummary(duplicate);`, 20*time.Second), nil
	case "page_item.move":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var target = ahdPlacementTarget(doc, params.targetId); item.move(target, ahdPlacementValue(params.placement || "inside")); ahdTryRedraw(); return { itemId: params.itemId, targetId: params.targetId, placement: params.placement || "inside", moved: true, bounds: ahdBounds(item) };`, 20*time.Second), nil
	case "page_item.resize":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); item.resize(params.scaleX, params.scaleY, params.changePositions, params.changeFillPatterns, params.changeFillGradients, params.changeStrokePatterns, params.changeLineWidths, ahdTransformationValue(params.anchor)); ahdTryRedraw(); return { itemId: params.itemId, scaleX: params.scaleX, scaleY: params.scaleY, bounds: ahdBounds(item) };`, 20*time.Second), nil
	case "page_item.rotate":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); item.rotate(params.angle, params.changePositions, params.changeFillPatterns, params.changeFillGradients, params.changeStrokePatterns, ahdTransformationValue(params.anchor)); ahdTryRedraw(); return { itemId: params.itemId, angle: params.angle, bounds: ahdBounds(item) };`, 20*time.Second), nil
	case "page_item.transform":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); item.transform(ahdMatrixValue(params.matrix), params.changePositions, params.changeFillPatterns, params.changeFillGradients, params.changeStrokePatterns, params.changeLineWidths, ahdTransformationValue(params.anchor)); ahdTryRedraw(); return { itemId: params.itemId, matrix: ahdMatrixSummary(ahdMatrixValue(params.matrix)), bounds: ahdBounds(item) };`, 20*time.Second), nil
	case "page_item.translate":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); item.translate(params.deltaX || 0, params.deltaY || 0, params.transformObjects, params.transformFillPatterns, params.transformFillGradients, params.transformStrokePatterns); ahdTryRedraw(); return { itemId: params.itemId, deltaX: params.deltaX || 0, deltaY: params.deltaY || 0, bounds: ahdBounds(item) };`, 20*time.Second), nil
	case "page_item.z_order":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); item.zOrder(ahdZOrderMethodValue(params.method)); ahdTryRedraw(); return { itemId: params.itemId, method: params.method, reordered: true };`, 20*time.Second), nil
	case "page_item.bring_in_perspective":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); item.bringInPerspective(params.posX, params.posY, ahdPerspectivePlaneValue(params.plane)); ahdTryRedraw(); return { itemId: params.itemId, posX: params.posX, posY: params.posY, plane: params.plane, bounds: ahdBounds(item) };`, 20*time.Second), nil

	case "path.create_polygon":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var item = layer.pathItems.polygon(params.top, params.left, params.radius, params.sides); item.name = params.name; return { name: item.name, radius: params.radius, sides: params.sides, width: item.width, height: item.height };`, 20*time.Second), nil
	case "path.create_rounded_rect":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var item = layer.pathItems.roundedRectangle(params.top, params.left, params.width, params.height, params.horizontalRadius != null ? params.horizontalRadius : 15, params.verticalRadius != null ? params.verticalRadius : 20, params.reversed === true); item.name = params.name; return { name: item.name, width: item.width, height: item.height, horizontalRadius: params.horizontalRadius != null ? params.horizontalRadius : 15, verticalRadius: params.verticalRadius != null ? params.verticalRadius : 20, reversed: params.reversed === true };`, 20*time.Second), nil
	case "path.create_star":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var item = layer.pathItems.star(params.top, params.left, params.outerRadius, params.innerRadius, params.points); item.name = params.name; return { name: item.name, outerRadius: params.outerRadius, innerRadius: params.innerRadius, points: params.points, width: item.width, height: item.height };`, 20*time.Second), nil
	case "path.set_entire_path":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); if (!item.setEntirePath) throw new Error("Item does not support setEntirePath: " + params.itemId); item.setEntirePath(params.points); if (params.closed != null) item.closed = params.closed; return { itemId: params.itemId, pointCount: params.points.length, closed: item.closed === true, bounds: ahdBounds(item) };`, 20*time.Second), nil

	case "text.create_area":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var layer = ahdLayer(doc, params.layerId); var path = layer.pathItems.rectangle(params.top, params.left, params.width, params.height); path.name = params.name + " Path"; var text = doc.textFrames.areaText(path); text.name = params.name; text.contents = params.contents; return { name: text.name, contents: text.contents, width: path.width, height: path.height };`, 30*time.Second), nil
	case "text.create_on_path":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var path = ahdPageItem(doc, params.pathItemId); var text = doc.textFrames.pathText(path); text.name = params.name; text.contents = params.contents; return { name: text.name, pathItemId: params.pathItemId, contents: text.contents };`, 30*time.Second), nil
	case "text.change_case":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); if (!text.textRange) throw new Error("Item does not expose a textRange: " + params.itemId); text.textRange.changeCaseTo(CaseChangeType[params.caseType]); return { itemId: params.itemId, caseType: params.caseType, contents: text.contents };`, 20*time.Second), nil
	case "text.thread":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); var next = ahdPageItem(doc, params.nextItemId); text.nextFrame = next; return { itemId: params.itemId, nextItemId: params.nextItemId, threaded: true };`, 30*time.Second), nil
	case "text.convert_to_area":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); var converted = text.convertPointObjectToAreaObject(); var result = converted || text; return { itemId: params.itemId, convertedName: result.name || params.itemId, kind: result.kind ? String(result.kind) : null };`, 30*time.Second), nil
	case "text.convert_to_point":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var text = ahdPageItem(doc, params.itemId); var converted = text.convertAreaObjectToPointObject(); var result = converted || text; return { itemId: params.itemId, convertedName: result.name || params.itemId, kind: result.kind ? String(result.kind) : null };`, 30*time.Second), nil

	case "style.character.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var styles = []; for (var i = 0; i < doc.characterStyles.length; i++) { styles.push({ index: i, name: doc.characterStyles[i].name || "", typename: doc.characterStyles[i].typename }); } return styles;`, 10*time.Second), nil
	case "style.character.apply":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var style = ahdCharacterStyle(doc, params.styleName); var text = ahdPageItem(doc, params.itemId); style.applyTo(text.textRange); return { styleName: style.name, itemId: params.itemId, applied: true };`, 20*time.Second), nil
	case "style.character.import":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.importCharacterStyles(file); return { filePath: file.fsName, imported: true, count: doc.characterStyles.length };`, 20*time.Second), nil
	case "style.paragraph.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var styles = []; for (var i = 0; i < doc.paragraphStyles.length; i++) { styles.push({ index: i, name: doc.paragraphStyles[i].name || "", typename: doc.paragraphStyles[i].typename }); } return styles;`, 10*time.Second), nil
	case "style.paragraph.apply":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var style = ahdParagraphStyle(doc, params.styleName); var text = ahdPageItem(doc, params.itemId); for (var i = 0; i < text.paragraphs.length; i++) { style.applyTo(text.paragraphs[i], params.clearOverrides === true); } return { styleName: style.name, itemId: params.itemId, paragraphs: text.paragraphs.length, applied: true };`, 20*time.Second), nil
	case "style.paragraph.import":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.importParagraphStyles(file); return { filePath: file.fsName, imported: true, count: doc.paragraphStyles.length };`, 20*time.Second), nil
	case "style.graphic.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var styles = []; for (var i = 0; i < doc.graphicStyles.length; i++) { styles.push({ index: i, name: doc.graphicStyles[i].name || "", typename: doc.graphicStyles[i].typename }); } return styles;`, 10*time.Second), nil
	case "style.graphic.apply":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var style = ahdGraphicStyle(doc, params.styleName); var item = ahdPageItem(doc, params.itemId); style.applyTo(item); ahdTryRedraw(); return { styleName: style.name, itemId: params.itemId, applied: true, item: ahdPageItemSummary(item) };`, 20*time.Second), nil
	case "style.graphic.merge":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var style = ahdGraphicStyle(doc, params.styleName); var item = ahdPageItem(doc, params.itemId); var merged = style.mergeTo(item); ahdTryRedraw(); return { styleName: style.name, itemId: params.itemId, merged: merged === true || merged == null, item: ahdPageItemSummary(item) };`, 20*time.Second), nil
	case "style.graphic.remove":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var style = ahdGraphicStyle(doc, params.styleName); var name = style.name || params.styleName; style.remove(); return { styleName: params.styleName, name: name, removed: true, remaining: doc.graphicStyles.length };`, 20*time.Second), nil

	case "swatch.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.swatches.length; i++) { items.push(ahdSwatchSummary(doc.swatches[i])); } return items;`, 10*time.Second), nil
	case "swatch.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var swatch = doc.swatches.add(); swatch.name = params.name; swatch.color = ahdColorFromMode(params); return ahdSwatchSummary(swatch);`, 20*time.Second), nil
	case "swatch.delete":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var swatch = ahdSwatch(doc, params.swatchId); var name = swatch.name; swatch.remove(); return { swatchId: params.swatchId, name: name, deleted: true };`, 20*time.Second), nil

	case "spot.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.spots.length; i++) { items.push(ahdSpotSummary(doc.spots[i])); } return items;`, 10*time.Second), nil
	case "spot.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var spot = doc.spots.add(); spot.name = params.name; spot.colorType = ColorModel[params.colorType || "SPOT"]; spot.color = ahdColorFromMode(params); return ahdSpotSummary(spot);`, 20*time.Second), nil
	case "spot.delete":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var spot = ahdSpot(doc, params.spotId); var name = spot.name; spot.remove(); return { spotId: params.spotId, name: name, deleted: true };`, 20*time.Second), nil

	case "symbol.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var symbols = []; var symbolItems = []; for (var i = 0; i < doc.symbols.length; i++) { symbols.push(ahdSymbolSummary(doc.symbols[i])); } for (var j = 0; j < doc.symbolItems.length; j++) { symbolItems.push(ahdSymbolItemSummary(doc.symbolItems[j])); } return { symbols: symbols, symbolItems: symbolItems };`, 10*time.Second), nil
	case "symbol.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var source = ahdPageItem(doc, params.itemId); var symbol; if (params.registrationPoint) { symbol = doc.symbols.add(source, ahdSymbolRegistrationPointValue(params.registrationPoint)); } else { symbol = doc.symbols.add(source); } if (params.name) { symbol.name = params.name; } return ahdSymbolSummary(symbol);`, 20*time.Second), nil
	case "symbol.place":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var symbol = ahdSymbol(doc, params.symbolId); var item = doc.symbolItems.add(symbol); ahdStampIdentifier(item, params.name || item.name || params.symbolId); if (params.left != null || params.top != null) { item.position = [params.left != null ? params.left : item.position[0], params.top != null ? params.top : item.position[1]]; } ahdTryRedraw(); return ahdSymbolItemSummary(item);`, 20*time.Second), nil
	case "symbol.break_link":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdSymbolItem(doc, params.itemId); var expanded = item.breakLink(); ahdTryRedraw(); return { itemId: params.itemId, broken: true, expanded: expanded ? ahdPageItemSummary(expanded) : null };`, 20*time.Second), nil

	case "placed.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.placedItems.length; i++) { items.push(ahdPlacedSummary(doc.placedItems[i])); } return items;`, 10*time.Second), nil
	case "placed.place":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = doc.placedItems.add(); item.file = new File(params.filePath); ahdStampIdentifier(item, params.name || item.name || params.filePath); if (params.left != null || params.top != null) { item.position = [params.left != null ? params.left : 0, params.top != null ? params.top : doc.height]; } ahdTryRedraw(); if (params.embed) { item.embed(); ahdTryRedraw(); return { filePath: params.filePath, embedded: true }; } return ahdPlacedSummary(item);`, 30*time.Second), nil
	case "placed.embed":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPlacedItem(doc, params.itemId); var name = item.name || params.itemId; item.embed(); ahdTryRedraw(); return { itemId: params.itemId, name: name, embedded: true };`, 30*time.Second), nil
	case "placed.relink":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPlacedItem(doc, params.itemId); var file = new File(params.filePath); item.relink(file); ahdTryRedraw(); return ahdPlacedSummary(item);`, 30*time.Second), nil
	case "placed.trace":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPlacedItem(doc, params.itemId); ahdTryRedraw(); var plugin = item.trace(); ahdStampIdentifier(plugin, params.name || plugin.name || params.itemId + " Trace"); ahdTryRedraw(); if (params.presetName) plugin.tracing.tracingOptions.loadFromPreset(params.presetName); ahdTryRedraw(); var summary = ahdTracingSummary(plugin); if (params.expand) { var group = plugin.tracing.expandTracing(params.viewed === true); summary.expanded = true; summary.groupName = group.name || ""; summary.groupType = group.typename; } return summary;`, 90*time.Second), nil

	case "raster.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.rasterItems.length; i++) { items.push(ahdRasterSummary(doc.rasterItems[i])); } return items;`, 10*time.Second), nil
	case "raster.rasterize":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var source = ahdPageItem(doc, params.itemId); var options = new RasterizeOptions(); if (params.resolution != null) options.resolution = params.resolution; var raster; if (params.clipBounds) { raster = doc.rasterize(source, ahdRectValue(params.clipBounds), options); } else { raster = doc.rasterize(source, undefined, options); } ahdStampIdentifier(raster, params.name || raster.name || params.itemId + " Raster"); ahdTryRedraw(); return ahdRasterSummary(raster);`, 45*time.Second), nil
	case "raster.trace":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdRasterItem(doc, params.itemId); ahdTryRedraw(); var plugin = item.trace(); ahdStampIdentifier(plugin, params.name || plugin.name || params.itemId + " Trace"); ahdTryRedraw(); if (params.presetName) plugin.tracing.tracingOptions.loadFromPreset(params.presetName); ahdTryRedraw(); var summary = ahdTracingSummary(plugin); if (params.expand) { var group = plugin.tracing.expandTracing(params.viewed === true); summary.expanded = true; summary.groupName = group.name || ""; summary.groupType = group.typename; } return summary;`, 90*time.Second), nil
	case "raster.colorize":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdRasterItem(doc, params.itemId); item.colorize(ahdColorFromMode(params)); ahdTryRedraw(); return ahdRasterSummary(item);`, 20*time.Second), nil
	case "raster.release_tracing":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var plugin = ahdTracingPlugin(doc, params.itemId); var released = plugin.tracing.releaseTracing(); return { itemId: params.itemId, releasedType: released.typename, releasedName: released.name || "" };`, 45*time.Second), nil

	case "repeat.grid.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.gridRepeatItems.length; i++) { items.push(ahdRepeatSummary(doc.gridRepeatItems[i], "grid")); } return items;`, 10*time.Second), nil
	case "repeat.grid.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var source = ahdPageItem(doc, params.itemId); var item = doc.gridRepeatItems.add(source, ahdRepeatConfigValue("grid", params.config)); if (params.name) item.name = params.name; return ahdRepeatSummary(item, "grid");`, 30*time.Second), nil
	case "repeat.grid.update":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdRepeatItem(doc.gridRepeatItems, params.repeatId, "Grid repeat"); var result = item.setGridConfiguration(ahdRepeatConfigValue("grid", params.config), GridRepeatUpdate[params.state || "GRIDALL"]); return { repeatId: params.repeatId, kind: "grid", state: params.state || "GRIDALL", result: result };`, 30*time.Second), nil

	case "repeat.radial.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.radialRepeatItems.length; i++) { items.push(ahdRepeatSummary(doc.radialRepeatItems[i], "radial")); } return items;`, 10*time.Second), nil
	case "repeat.radial.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var source = ahdPageItem(doc, params.itemId); var item = doc.radialRepeatItems.add(source, ahdRepeatConfigValue("radial", params.config)); if (params.name) item.name = params.name; return ahdRepeatSummary(item, "radial");`, 30*time.Second), nil
	case "repeat.radial.update":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdRepeatItem(doc.radialRepeatItems, params.repeatId, "Radial repeat"); item.setRadialConfiguration(ahdRepeatConfigValue("radial", params.config), RadialRepeatUpdate[params.state || "RADIALALL"]); return { repeatId: params.repeatId, kind: "radial", state: params.state || "RADIALALL", updated: true };`, 30*time.Second), nil

	case "repeat.symmetry.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.symmetryRepeatItems.length; i++) { items.push(ahdRepeatSummary(doc.symmetryRepeatItems[i], "symmetry")); } return items;`, 10*time.Second), nil
	case "repeat.symmetry.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var source = ahdPageItem(doc, params.itemId); var item = doc.symmetryRepeatItems.add(source, ahdRepeatConfigValue("symmetry", params.config)); if (params.name) item.name = params.name; return ahdRepeatSummary(item, "symmetry");`, 30*time.Second), nil
	case "repeat.symmetry.update":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdRepeatItem(doc.symmetryRepeatItems, params.repeatId, "Symmetry repeat"); item.setSymmetryConfiguration(ahdRepeatConfigValue("symmetry", params.config), SymmetryRepeatUpdate[params.state || "SYMMETRYALL"]); return { repeatId: params.repeatId, kind: "symmetry", state: params.state || "SYMMETRYALL", updated: true };`, 30*time.Second), nil

	case "dataset.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.dataSets.length; i++) { items.push(ahdDataSetSummary(doc.dataSets[i], i)); } return items;`, 10*time.Second), nil
	case "dataset.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = doc.dataSets.add(); if (params.name) item.name = params.name; return ahdDataSetSummary(item, doc.dataSets.length - 1);`, 20*time.Second), nil
	case "dataset.delete":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdDataSet(doc, params.datasetId); var name = item.name || params.datasetId; item.remove(); return { datasetId: params.datasetId, name: name, deleted: true };`, 20*time.Second), nil
	case "dataset.apply":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var dataSet = ahdDataSet(doc, params.datasetId); dataSet.display(); app.redraw(); return { datasetId: params.datasetId, displayed: true };`, 20*time.Second), nil
	case "dataset.update":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var dataSet = ahdDataSet(doc, params.datasetId); dataSet.update(); return { datasetId: params.datasetId, updated: true };`, 20*time.Second), nil
	case "dataset.import":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.importVariables(file); return { filePath: file.fsName, imported: true, variables: doc.variables.length, dataSets: doc.dataSets.length };`, 30*time.Second), nil
	case "dataset.export":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.exportVariables(file); return { filePath: file.fsName, exported: true, variables: doc.variables.length, dataSets: doc.dataSets.length };`, 30*time.Second), nil

	case "variable.list":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var items = []; for (var i = 0; i < doc.variables.length; i++) { items.push(ahdVariableSummary(doc.variables[i], i)); } return items;`, 10*time.Second), nil
	case "variable.create":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = doc.variables.add(); item.kind = ahdVariableKindValue(params.kind); if (params.name) item.name = params.name; return ahdVariableSummary(item, doc.variables.length - 1);`, 20*time.Second), nil
	case "variable.delete":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdVariable(doc, params.variableId); var name = item.name || params.variableId; item.remove(); return { variableId: params.variableId, name: name, deleted: true };`, 20*time.Second), nil
	case "variable.bind_visibility":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); var variable = ahdVariable(doc, params.variableId); item.visibilityVariable = variable; return { itemId: params.itemId, variableId: params.variableId, binding: "visibility", kind: variable.kind ? String(variable.kind) : null };`, 20*time.Second), nil
	case "variable.bind_text":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); if (item.typename !== "TextFrame") { throw new Error("Text variable bindings require a TextFrame item"); } var variable = ahdVariable(doc, params.variableId); item.contentVariable = variable; return { itemId: params.itemId, variableId: params.variableId, binding: "text", kind: variable.kind ? String(variable.kind) : null };`, 20*time.Second), nil
	case "variable.bind_content":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var item = ahdPageItem(doc, params.itemId); if (item.contentVariable === undefined) { throw new Error("Item does not support content-variable bindings: " + params.itemId); } var variable = ahdVariable(doc, params.variableId); item.contentVariable = variable; return { itemId: params.itemId, variableId: params.variableId, binding: "content", kind: variable.kind ? String(variable.kind) : null };`, 20*time.Second), nil
	case "variable.import":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.importVariables(file); return { filePath: file.fsName, imported: true, variables: doc.variables.length, dataSets: doc.dataSets.length };`, 30*time.Second), nil
	case "variable.export":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.filePath); doc.exportVariables(file); return { filePath: file.fsName, exported: true, variables: doc.variables.length, dataSets: doc.dataSets.length };`, 30*time.Second), nil
	case "trace.preset.list":
		return scriptPlan(request.Params, `return { presets: ahdArrayLikeStrings(app.tracingPresetList) };`, 10*time.Second), nil
	case "trace.preset.store":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var plugin = ahdTracingPlugin(doc, params.itemId); var stored = plugin.tracing.tracingOptions.storeToPreset(params.presetName) === true; return { itemId: params.itemId, presetName: params.presetName, stored: stored };`, 20*time.Second), nil

	case "capture.image":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); if (params.clipBounds && params.options) { doc.imageCapture(file, ahdRectValue(params.clipBounds), ahdImageCaptureOptions(params.options)); } else if (params.clipBounds) { doc.imageCapture(file, ahdRectValue(params.clipBounds)); } else if (params.options) { doc.imageCapture(file, undefined, ahdImageCaptureOptions(params.options)); } else { doc.imageCapture(file); } return { outputPath: ahdResolvedOutputPath(file, null), captured: true };`, 30*time.Second), nil
	case "capture.window":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); doc.windowCapture(file, [params.width, params.height]); return { outputPath: ahdResolvedOutputPath(file, [".tiff", ".tif"]), captured: true, width: params.width, height: params.height };`, 30*time.Second), nil

	case "print.presets":
		return scriptPlan(request.Params, `return { presets: ahdArrayLikeStrings(app.printPresetsList) };`, 20*time.Second), nil
	case "print.devices":
		return scriptPlan(request.Params, `return { printers: [], printerEnumeration: "skipped", ppdFiles: [], ppdEnumeration: "skipped", presets: ahdArrayLikeStrings(app.printPresetsList) };`, 10*time.Second), nil
	case "print.run":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var options = new PrintOptions(); if (params.printerName) options.printerName = params.printerName; if (params.PPDName) options.PPDName = params.PPDName; if (params.printPreset) options.printPreset = params.printPreset; ahdSetPrintJobOptions(options, params.jobOptions); doc.print(options); return { printed: true, printerName: params.printerName || null, PPDName: params.PPDName || null, printPreset: params.printPreset || null };`, 45*time.Second), nil

	case "export.gif":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsGIF(); if (params.scale) { options.horizontalScale = params.scale * 100; options.verticalScale = params.scale * 100; } ahdApplyArtboardExport(doc, params.artboardId, options); doc.exportFile(file, ahdExportTypeValue("gif"), options); return { outputPath: ahdResolvedOutputPath(file, [".gif"]), format: "gif", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.png8":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsPNG8(); if (params.scale) { options.horizontalScale = params.scale * 100; options.verticalScale = params.scale * 100; } if (params.colorCount != null) options.colorCount = params.colorCount; ahdApplyArtboardExport(doc, params.artboardId, options); doc.exportFile(file, ahdExportTypeValue("png8"), options); return { outputPath: ahdResolvedOutputPath(file, [".png"]), format: "png8", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.tiff":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsTIFF(); if (params.resolution != null) options.resolution = params.resolution; ahdApplyArtboardExport(doc, params.artboardId, options); doc.exportFile(file, ahdExportTypeValue("tiff"), options); return { outputPath: ahdResolvedOutputPath(file, [".tif", ".tiff"]), format: "tiff", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.webp":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsWebP(); if (params.losslessCompression != null) options.losslessCompression = params.losslessCompression; if (params.imageQuality != null) options.imageQuality = params.imageQuality; if (params.isTransparent != null) options.isTransparent = params.isTransparent; if (params.ppi != null) options.pPI = params.ppi; doc.exportFile(file, ahdExportTypeValue("webp"), options); return { outputPath: ahdResolvedOutputPath(file, [".webp"]), format: "webp" };`, 30*time.Second), nil
	case "export.photoshop":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsPhotoshop(); if (params.resolution != null) options.resolution = params.resolution; if (params.editableText != null) options.editableText = params.editableText; if (params.maximumEditability != null) options.maximumEditability = params.maximumEditability; ahdApplyArtboardExport(doc, params.artboardId, options); doc.exportFile(file, ahdExportTypeValue("photoshop"), options); return { outputPath: ahdResolvedOutputPath(file, [".psd"]), format: "photoshop", artboardId: params.artboardId || null };`, 30*time.Second), nil
	case "export.autocad":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new ExportOptionsAutoCAD(); doc.exportFile(file, ahdExportTypeValue("autocad"), options); return { outputPath: ahdResolvedOutputPath(file, [".dwg", ".dxf"]), format: "autocad" };`, 30*time.Second), nil
	case "export.eps":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new EPSSaveOptions(); doc.saveAs(file, options); return { outputPath: ahdResolvedOutputPath(file, [".eps"]), format: "eps" };`, 30*time.Second), nil
	case "export.fxg":
		return scriptPlan(request.Params, `var doc = ahdDoc(); var file = new File(params.outputPath); var options = new FXGSaveOptions(); doc.saveAs(file, options); return { outputPath: ahdResolvedOutputPath(file, [".fxg"]), format: "fxg" };`, 30*time.Second), nil
	default:
		return executionPlan{}, fmt.Errorf("%s execution is not wired yet", request.Command.Name)
	}
}
