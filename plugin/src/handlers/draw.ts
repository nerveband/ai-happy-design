import { getNodeById, getParentById, getSceneNodeById } from '../utils/getNode';
import { resolveStableId } from '../utils/stableId';

export async function handleDraw(action: string, params: any): Promise<any> {
  switch (action) {
    case 'text_path':
    case 'create_text_path': return createTextPath(params);
    case 'transform_group':
    case 'create_transform_group': return createTransformGroup(params);
    case 'load_brush':
    case 'load_brushes': return loadBrushes(params);
    case 'variable_width_stroke':
    case 'set_variable_width_stroke': return setVariableWidthStroke(params);
    case 'brush_stroke':
    case 'set_brush_stroke': return setBrushStroke(params);
    case 'dynamic_stroke':
    case 'set_dynamic_stroke': return setDynamicStroke(params);
    case 'pattern_fill':
    case 'set_pattern_fill': return setPatternFill(params);
    case 'pattern_stroke':
    case 'set_pattern_stroke': return setPatternStroke(params);
    default: throw new Error('Unknown draw action: ' + action + '. Available: create_text_path, create_transform_group, load_brushes, set_variable_width_stroke, set_brush_stroke, set_dynamic_stroke, set_pattern_fill, set_pattern_stroke');
  }
}

function requireFigmaMethod(name: string): Function {
  var fn = (figma as any)[name];
  if (typeof fn !== 'function') {
    throw new Error('Figma API ' + name + ' is unavailable in this runtime');
  }
  return fn;
}

function parseList(value: any, fallback: string[] = []): string[] {
  if (Array.isArray(value)) return value;
  if (typeof value === 'string') {
    return value.split(',').map(function(part: string) { return part.trim(); }).filter(Boolean);
  }
  return fallback;
}

function parseJsonValue(value: any, fallback: any) {
  if (value == null) return fallback;
  if (typeof value !== 'string') return value;
  try {
    return JSON.parse(value);
  } catch {
    return fallback;
  }
}

function createPatternPaint(params: any): any {
  if (!params.sourceNodeId) throw new Error('sourceNodeId is required');
  return {
    type: 'PATTERN',
    sourceNodeId: params.sourceNodeId,
    tileType: params.tileType || 'RECTANGULAR',
    scalingFactor: params.scalingFactor == null ? 1 : params.scalingFactor,
    spacing: params.spacing || { x: params.spacingX || 0, y: params.spacingY || 0 },
    horizontalAlignment: params.horizontalAlignment || 'START',
    visible: params.visible !== false,
    opacity: params.opacity == null ? 1 : params.opacity,
    blendMode: params.blendMode || 'NORMAL',
  };
}

async function createTextPath(params: any) {
  var createTextPathFn = requireFigmaMethod('createTextPath');
  var node = await getSceneNodeById(params.nodeId);
  var textPath = createTextPathFn.call(figma, node, params.startSegment || 0, params.startPosition == null ? 0 : params.startPosition);

  if (params.characters != null || params.text != null) {
    var fontName = textPath.fontName;
    if (fontName && fontName !== figma.mixed) {
      await figma.loadFontAsync(fontName);
    }
    textPath.characters = String(params.characters != null ? params.characters : params.text);
  }
  if (params.name) textPath.name = params.name;
  if (params.fontSize != null) textPath.fontSize = params.fontSize;

  return { id: textPath.id, name: textPath.name, type: textPath.type };
}

async function createTransformGroup(params: any) {
  var transformGroupFn = requireFigmaMethod('transformGroup');
  var nodeIds = parseList(params.nodeIds);
  if (nodeIds.length === 0) throw new Error('nodeIds is required');

  var nodes: SceneNode[] = [];
  for (var i = 0; i < nodeIds.length; i++) {
    nodes.push(await getSceneNodeById(nodeIds[i]));
  }

  var parent = params.parentId ? await getParentById(params.parentId) : figma.currentPage;
  if (!parent) throw new Error('Parent not found: ' + params.parentId);
  var index = params.index == null ? parent.children.length : params.index;
  var modifiers = parseJsonValue(params.modifiers, []);
  if (!Array.isArray(modifiers)) throw new Error('modifiers must be an array');

  var group = transformGroupFn.call(figma, nodes, parent, index, modifiers);
  if (params.name) group.name = params.name;
  var stableId = await resolveStableId(group, parent);
  return { id: stableId, name: group.name, type: group.type, childCount: group.children.length };
}

async function loadBrushes(params: any) {
  var loadBrushesFn = requireFigmaMethod('loadBrushesAsync');
  var brushTypes = parseList(params.brushTypes || params.types || params.brushType || params.type, ['STRETCH', 'SCATTER']);
  for (var i = 0; i < brushTypes.length; i++) {
    await loadBrushesFn.call(figma, String(brushTypes[i]).toUpperCase());
  }
  return { loaded: brushTypes.map(function(t: string) { return String(t).toUpperCase(); }) };
}

async function setVariableWidthStroke(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  if (!('variableWidthStrokeProperties' in node)) {
    throw new Error('Node ' + params.nodeId + ' does not support variable width strokes');
  }

  var widthProfile = String(params.widthProfile || params.profile || 'UNIFORM').toUpperCase();
  if (widthProfile === 'CUSTOM') {
    var points = parseJsonValue(params.variableWidthPoints || params.points, []);
    if (!Array.isArray(points)) throw new Error('variableWidthPoints must be an array');
    (node as any).variableWidthStrokeProperties = { widthProfile: 'CUSTOM', variableWidthPoints: points };
  } else {
    (node as any).variableWidthStrokeProperties = { widthProfile: widthProfile };
  }
  return { id: node.id, name: node.name, variableWidthStrokeProperties: (node as any).variableWidthStrokeProperties };
}

async function setBrushStroke(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  if (!('complexStrokeProperties' in node)) {
    throw new Error('Node ' + params.nodeId + ' does not support brush strokes');
  }

  var brushType = String(params.brushType || 'STRETCH').toUpperCase();
  if (typeof (figma as any).loadBrushesAsync === 'function') {
    await (figma as any).loadBrushesAsync(brushType);
  }

  if (brushType === 'SCATTER') {
    (node as any).complexStrokeProperties = {
      type: 'BRUSH',
      brushType: 'SCATTER',
      brushName: params.brushName || 'BUBBLEGUM',
      gap: params.gap == null ? 1 : params.gap,
      wiggle: params.wiggle == null ? 0 : params.wiggle,
      sizeJitter: params.sizeJitter == null ? 0 : params.sizeJitter,
      angularJitter: params.angularJitter == null ? 0 : params.angularJitter,
      rotation: params.rotation == null ? 0 : params.rotation,
    };
  } else {
    (node as any).complexStrokeProperties = {
      type: 'BRUSH',
      brushType: 'STRETCH',
      brushName: params.brushName || 'HEIST',
      direction: params.direction || 'FORWARD',
    };
  }
  return { id: node.id, name: node.name, complexStrokeProperties: (node as any).complexStrokeProperties };
}

async function setDynamicStroke(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  if (!('complexStrokeProperties' in node)) {
    throw new Error('Node ' + params.nodeId + ' does not support dynamic strokes');
  }
  var smoothen = params.smoothen == null ? 0 : params.smoothen;
  if (params.smoothen == null) smoothen = smoothen + 0.5;
  (node as any).complexStrokeProperties = {
    type: 'DYNAMIC',
    frequency: params.frequency == null ? 1 : params.frequency,
    wiggle: params.wiggle == null ? 1 : params.wiggle,
    smoothen: smoothen,
  };
  return { id: node.id, name: node.name, complexStrokeProperties: (node as any).complexStrokeProperties };
}

async function setPatternFill(params: any) {
  var node = await getNodeById(params.nodeId);
  if (!('setFillsAsync' in node)) throw new Error('Node ' + params.nodeId + ' does not support async pattern fills');
  var paint = createPatternPaint(params);
  var append = params.append === true;
  var fills = append && Array.isArray((node as any).fills) ? JSON.parse(JSON.stringify((node as any).fills)) : [];
  fills.push(paint);
  await (node as any).setFillsAsync(fills);
  return { id: node.id, name: node.name, fillCount: fills.length };
}

async function setPatternStroke(params: any) {
  var node = await getNodeById(params.nodeId);
  if (!('setStrokesAsync' in node)) throw new Error('Node ' + params.nodeId + ' does not support async pattern strokes');
  var paint = createPatternPaint(params);
  var append = params.append === true;
  var strokes = append && Array.isArray((node as any).strokes) ? JSON.parse(JSON.stringify((node as any).strokes)) : [];
  strokes.push(paint);
  await (node as any).setStrokesAsync(strokes);
  return { id: node.id, name: node.name, strokeCount: strokes.length };
}
