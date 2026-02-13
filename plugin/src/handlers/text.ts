import { loadFont, loadNodeFonts, resolveFontFamily } from '../utils/fonts';

export async function handleText(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create': return createText(params);
    case 'set_content': return setContent(params);
    case 'set_font': return setFont(params);
    case 'set_size': return setSize(params);
    case 'set_weight': return setWeight(params);
    case 'set_color': return setColor(params);
    case 'set_align': return setAlign(params);
    case 'set_line_height': return setLineHeight(params);
    case 'set_letter_spacing': return setLetterSpacing(params);
    case 'set_decoration': return setDecoration(params);
    case 'set_case': return setCase(params);
    case 'set_paragraph_spacing': return setParagraphSpacing(params);
    case 'get_content': return getContent(params);
    default: throw new Error(`Unknown text action: ${action}`);
  }
}

function getTextNode(nodeId: string): TextNode {
  const node = figma.getNodeById(nodeId);
  if (!node || node.type !== 'TEXT') throw new Error(`Node ${nodeId} is not a text node`);
  return node as TextNode;
}

async function createText(params: any) {
  const { x = 0, y = 0, content = '', fontFamily = 'Inter', fontStyle = 'Regular', fontSize = 16, color, name, parentId, width } = params;

  const family = resolveFontFamily(fontFamily);
  await loadFont(family, fontStyle);

  const text = figma.createText();
  text.fontName = { family, style: fontStyle };
  text.fontSize = fontSize;
  text.x = x;
  text.y = y;
  if (name) text.name = name;
  if (content) text.characters = content;
  if (color) {
    text.fills = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }
  if (width !== undefined) {
    text.resize(width, text.height);
    text.textAutoResize = 'HEIGHT';
  }

  if (parentId) {
    const parent = figma.getNodeById(parentId);
    if (parent && 'appendChild' in parent) (parent as FrameNode).appendChild(text);
  }

  return { id: text.id, name: text.name, type: text.type, width: text.width, height: text.height };
}

async function setContent(params: any) {
  const { nodeId, content } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);
  node.characters = content;
  return { id: node.id, name: node.name, characters: node.characters };
}

async function setFont(params: any) {
  const { nodeId, fontFamily, fontStyle = 'Regular', rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  const family = resolveFontFamily(fontFamily);
  await loadFont(family, fontStyle);

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeFontName(rangeStart, rangeEnd, { family, style: fontStyle });
  } else {
    await loadNodeFonts(node);
    node.fontName = { family, style: fontStyle };
  }
  return { id: node.id, name: node.name };
}

async function setSize(params: any) {
  const { nodeId, size, rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeFontSize(rangeStart, rangeEnd, size);
  } else {
    node.fontSize = size;
  }
  return { id: node.id, name: node.name };
}

async function setWeight(params: any) {
  const { nodeId, weight, rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);

  // Map numeric weights to style names
  const weightMap: Record<number, string> = {
    100: 'Thin', 200: 'Extra Light', 300: 'Light', 400: 'Regular',
    500: 'Medium', 600: 'Semi Bold', 700: 'Bold', 800: 'Extra Bold', 900: 'Black',
  };
  const styleName = typeof weight === 'string' ? weight : (weightMap[weight] || 'Regular');

  const currentFont = node.fontName;
  const family = currentFont === figma.mixed ? 'Inter' : currentFont.family;
  await loadFont(family, styleName);

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeFontName(rangeStart, rangeEnd, { family, style: styleName });
  } else {
    node.fontName = { family, style: styleName };
  }
  return { id: node.id, name: node.name };
}

async function setColor(params: any) {
  const { nodeId, color, rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  const fill: SolidPaint = {
    type: 'SOLID',
    color: { r: color.r, g: color.g, b: color.b },
    opacity: color.a ?? 1,
  };

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeFills(rangeStart, rangeEnd, [fill]);
  } else {
    node.fills = [fill];
  }
  return { id: node.id, name: node.name };
}

async function setAlign(params: any) {
  const { nodeId, horizontal, vertical } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  if (horizontal) node.textAlignHorizontal = horizontal;
  if (vertical) node.textAlignVertical = vertical;
  return { id: node.id, name: node.name };
}

async function setLineHeight(params: any) {
  const { nodeId, value, unit = 'PIXELS', rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  const lineHeight: LineHeight = unit === 'AUTO'
    ? { unit: 'AUTO' }
    : { value, unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeLineHeight(rangeStart, rangeEnd, lineHeight);
  } else {
    node.lineHeight = lineHeight;
  }
  return { id: node.id, name: node.name };
}

async function setLetterSpacing(params: any) {
  const { nodeId, value, unit = 'PIXELS', rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  const letterSpacing: LetterSpacing = { value, unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeLetterSpacing(rangeStart, rangeEnd, letterSpacing);
  } else {
    node.letterSpacing = letterSpacing;
  }
  return { id: node.id, name: node.name };
}

async function setDecoration(params: any) {
  const { nodeId, decoration, rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeTextDecoration(rangeStart, rangeEnd, decoration);
  } else {
    node.textDecoration = decoration;
  }
  return { id: node.id, name: node.name };
}

async function setCase(params: any) {
  const { nodeId, textCase, rangeStart, rangeEnd } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);

  if (rangeStart !== undefined && rangeEnd !== undefined) {
    node.setRangeTextCase(rangeStart, rangeEnd, textCase);
  } else {
    node.textCase = textCase;
  }
  return { id: node.id, name: node.name };
}

async function setParagraphSpacing(params: any) {
  const { nodeId, spacing } = params;
  const node = getTextNode(nodeId);
  await loadNodeFonts(node);
  node.paragraphSpacing = spacing;
  return { id: node.id, name: node.name };
}

async function getContent(params: any) {
  const { nodeId } = params;
  const node = getTextNode(nodeId);
  return {
    id: node.id,
    name: node.name,
    characters: node.characters,
    fontSize: node.fontSize,
    fontName: node.fontName,
    textAlignHorizontal: node.textAlignHorizontal,
    textAlignVertical: node.textAlignVertical,
    lineHeight: node.lineHeight,
    letterSpacing: node.letterSpacing,
    textDecoration: node.textDecoration,
    textCase: node.textCase,
  };
}
