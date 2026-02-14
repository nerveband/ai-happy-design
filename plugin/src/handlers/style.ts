import { getSceneNodeById } from '../utils/getNode';

function parseHexColor(color: any, fallback = { r: 0, g: 0, b: 0, a: 1 }) {
  if (color && typeof color === 'object' && typeof color.r === 'number') {
    return {
      r: color.r,
      g: color.g,
      b: color.b,
      a: typeof color.a === 'number' ? color.a : 1,
    };
  }

  if (typeof color !== 'string') {
    return fallback;
  }

  const raw = color.trim().replace(/^#/, '');
  if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) {
    return fallback;
  }

  const hex = raw.length === 3
    ? raw.split('').map(ch => ch + ch).join('')
    : raw;

  const hasAlpha = hex.length === 8;
  const n = parseInt(hex, 16);
  if (Number.isNaN(n)) return fallback;

  return {
    r: ((n >> (hasAlpha ? 24 : 16)) & 0xff) / 255,
    g: ((n >> (hasAlpha ? 16 : 8)) & 0xff) / 255,
    b: ((n >> (hasAlpha ? 8 : 0)) & 0xff) / 255,
    a: hasAlpha ? (n & 0xff) / 255 : 1,
  };
}

export async function handleStyle(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create_paint_style':
    case 'new_paint':
    case 'create_paint': return createPaintStyle(params);
    case 'create_text_style':
    case 'new_text':
    case 'create_text': return createTextStyle(params);
    case 'create_effect_style':
    case 'new_effect':
    case 'create_effect': return createEffectStyle(params);
    case 'apply_style':
    case 'set_style':
    case 'apply': return applyStyle(params);
    case 'list':
    case 'list_styles':
    case 'get_styles':
    case 'list_all':
    case 'get_all': return getAllStyles(params);
    case 'delete':
    case 'delete_style':
    case 'remove': return removeStyle(params);
    default: throw new Error('Unknown style action: ' + action + '. Available: create_paint, create_text, create_effect, apply, get_all, remove');
  }
}

async function createPaintStyle(params: any) {
  const { name, color, description } = params;
  const style = figma.createPaintStyle();
  style.name = name || 'Paint Style';
  if (description) style.description = description;

  const c = parseHexColor(color);
  style.paints = [{
    type: 'SOLID',
    color: { r: c.r, g: c.g, b: c.b },
    opacity: c.a,
  }];

  return { id: style.id, name: style.name, type: 'PAINT' };
}

async function createTextStyle(params: any) {
  const { name, description } = params;
  const style = figma.createTextStyle();
  style.name = name || 'Text Style';
  if (description) style.description = description;
  return { id: style.id, name: style.name, type: 'TEXT' };
}

async function createEffectStyle(params: any) {
  const { name, description } = params;
  const style = figma.createEffectStyle();
  style.name = name || 'Effect Style';
  if (description) style.description = description;
  return { id: style.id, name: style.name, type: 'EFFECT' };
}

async function applyStyle(params: any) {
  const { nodeId, styleId, styleType, target } = params;
  const node = await getSceneNodeById(nodeId);

  const resolvedType = (styleType || target || 'FILL').toUpperCase();
  if (resolvedType === 'TEXT') {
    if (!('textStyleId' in node)) throw new Error(`Node ${nodeId} does not support text styles`);
    await (node as TextNode).setTextStyleIdAsync(styleId || '');
  } else if (resolvedType === 'EFFECT') {
    if (!('effectStyleId' in node)) throw new Error(`Node ${nodeId} does not support effect styles`);
    await (node as any).setEffectStyleIdAsync(styleId || '');
  } else if (resolvedType === 'STROKE') {
    if (!('strokeStyleId' in node)) throw new Error(`Node ${nodeId} does not support stroke styles`);
    (node as GeometryMixin).strokeStyleId = styleId || '';
  } else {
    if (!('fillStyleId' in node)) throw new Error(`Node ${nodeId} does not support fill styles`);
    await (node as any).setFillStyleIdAsync(styleId || '');
  }

  return { id: node.id, name: node.name, styleId, styleType: resolvedType };
}

async function getAllStyles(_params: any) {
  var rawPaint = await figma.getLocalPaintStylesAsync();
  var paintStyles = rawPaint.map(function(s) {
    return { id: s.id, name: s.name, type: 'PAINT', description: s.description };
  });
  var rawText = await figma.getLocalTextStylesAsync();
  var textStyles = rawText.map(function(s) {
    return { id: s.id, name: s.name, type: 'TEXT', description: s.description };
  });
  var rawEffect = await figma.getLocalEffectStylesAsync();
  var effectStyles = rawEffect.map(function(s) {
    return { id: s.id, name: s.name, type: 'EFFECT', description: s.description };
  });
  var rawGrid = await figma.getLocalGridStylesAsync();
  var gridStyles = rawGrid.map(function(s) {
    return { id: s.id, name: s.name, type: 'GRID', description: s.description };
  });

  return {
    paintStyles: paintStyles,
    textStyles: textStyles,
    effectStyles: effectStyles,
    gridStyles: gridStyles,
    total: paintStyles.length + textStyles.length + effectStyles.length + gridStyles.length,
  };
}

async function removeStyle(params: any) {
  const { nodeId, styleType, target } = params;
  const node = await getSceneNodeById(nodeId);

  const resolvedType = (styleType || target || 'FILL').toUpperCase();
  if (resolvedType === 'TEXT' && 'textStyleId' in node) {
    await (node as TextNode).setTextStyleIdAsync('');
  } else if (resolvedType === 'EFFECT' && 'effectStyleId' in node) {
    await (node as any).setEffectStyleIdAsync('');
  } else if (resolvedType === 'STROKE' && 'strokeStyleId' in node) {
    (node as GeometryMixin).strokeStyleId = '';
  } else if ('fillStyleId' in node) {
    await (node as any).setFillStyleIdAsync('');
  } else {
    throw new Error(`Node ${nodeId} does not support ${resolvedType} styles`);
  }

  return { id: node.id, name: node.name, removed: resolvedType };
}
