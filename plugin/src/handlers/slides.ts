import { loadFont, resolveFontFamily } from '../utils/fonts';

function requireSlides(feature: string): any {
  var figmaAny = figma as any;
  if (figmaAny.editorType !== 'slides') {
    throw new Error(feature + ' is unavailable because current editorType is ' + figmaAny.editorType);
  }
  return figmaAny;
}

function parseHexColor(color: any, fallback = { r: 1, g: 1, b: 1, a: 1 }) {
  if (typeof color !== 'string') return fallback;
  var raw = color.trim().replace(/^#/, '');
  if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) return fallback;
  var hex = raw.length === 3 ? raw.split('').map(function(ch: string) { return ch + ch; }).join('') : raw;
  var n = parseInt(hex, 16);
  var hasAlpha = hex.length === 8;
  if (Number.isNaN(n)) return fallback;
  return {
    r: ((n >> (hasAlpha ? 24 : 16)) & 0xff) / 255,
    g: ((n >> (hasAlpha ? 16 : 8)) & 0xff) / 255,
    b: ((n >> (hasAlpha ? 8 : 0)) & 0xff) / 255,
    a: hasAlpha ? (n & 0xff) / 255 : 1,
  };
}

export async function handleSlides(action: string, params: any): Promise<any> {
  switch (action) {
    case 'get_current': return getCurrent();
    case 'create_slide': return createSlide(params);
    case 'set_background': return setBackground(params);
    case 'add_text': return addText(params);
    case 'reorder': return reorder(params);
    default: throw new Error('Unknown slides action: ' + action);
  }
}

async function getCurrent() {
  var figmaAny = requireSlides('Slides context');
  return { editorType: figmaAny.editorType, currentPage: { id: figma.currentPage.id, name: figma.currentPage.name } };
}

async function createSlide(params: any) {
  var figmaAny = requireSlides('Slide creation');
  if (typeof figmaAny.createSlide !== 'function') throw new Error('createSlide is unavailable in this Figma runtime');
  var slide = figmaAny.createSlide(params.row, params.col);
  if (params.name) slide.name = params.name;
  return { id: slide.id, name: slide.name, type: slide.type };
}

async function setBackground(params: any) {
  var figmaAny = requireSlides('Slide background');
  var node = await figma.getNodeByIdAsync(params.nodeId);
  if (!node || !('fills' in node)) throw new Error('Slide node does not support fills');
  var c = parseHexColor(params.color || '#FFFFFF');
  (node as any).fills = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
  return { id: (node as any).id, editorType: figmaAny.editorType };
}

async function addText(params: any) {
  requireSlides('Slide text');
  var parent = await figma.getNodeByIdAsync(params.nodeId);
  if (!parent || !('appendChild' in parent)) throw new Error('Slide parent not found');
  var family = resolveFontFamily(params.fontFamily || 'Inter');
  var style = params.fontStyle || 'Regular';
  await loadFont(family, style);
  var text = figma.createText();
  text.characters = String(params.text || '');
  text.fontName = { family: family, style: style };
  text.fontSize = params.fontSize || 32;
  text.x = params.x || 0;
  text.y = params.y || 0;
  (parent as any).appendChild(text);
  return { id: text.id, name: text.name, type: text.type };
}

async function reorder(params: any) {
  var figmaAny = requireSlides('Slide reorder');
  var node = await figma.getNodeByIdAsync(params.nodeId);
  if (!node) throw new Error('Slide not found');
  if (typeof (node as any).setSlideGridPosition === 'function') {
    (node as any).setSlideGridPosition(params.row || 0, params.col || 0);
    return { id: (node as any).id, row: params.row || 0, col: params.col || 0 };
  }
  return { id: (node as any).id, editorType: figmaAny.editorType, reordered: false, reason: 'setSlideGridPosition unavailable' };
}
