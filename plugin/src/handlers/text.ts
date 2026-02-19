import { loadFont, loadNodeFonts, resolveFontFamily } from '../utils/fonts';
import { getTextNodeById, getParentById } from '../utils/getNode';
import { resolveStableId } from '../utils/stableId';

function parseHexColor(color: any, fallback = { r: 0, g: 0, b: 0, a: 1 }) {
  if (color && typeof color === 'object' && typeof color.r === 'number') {
    return {
      r: color.r,
      g: color.g,
      b: color.b,
      a: typeof color.a === 'number' ? color.a : 1,
    };
  }

  if (typeof color !== 'string') return fallback;
  const raw = color.trim().replace(/^#/, '');
  if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) return fallback;

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

export async function handleText(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create_text':
    case 'add_text':
    case 'add':
    case 'new':
    case 'create': return createText(params);
    case 'set_text':
    case 'update_text':
    case 'edit':
    case 'set_content': return setContent(params);
    case 'font':
    case 'set_font_family':
    case 'set_font': return setFont(params);
    case 'size':
    case 'set_font_size':
    case 'font_size':
    case 'set_size': return setSize(params);
    case 'weight':
    case 'set_font_weight':
    case 'font_weight':
    case 'bold':
    case 'set_weight': return setWeight(params);
    case 'color':
    case 'set_text_color':
    case 'text_color':
    case 'set_color': return setColor(params);
    case 'align':
    case 'set_text_align':
    case 'alignment':
    case 'set_align': return setAlign(params);
    case 'set_spacing': return setSpacing(params);
    case 'set_line_height': return setLineHeight(params);
    case 'set_letter_spacing': return setLetterSpacing(params);
    case 'set_decoration': return setDecoration(params);
    case 'set_case': return setCase(params);
    case 'set_paragraph_spacing': return setParagraphSpacing(params);
    case 'content':
    case 'get_text':
    case 'read':
    case 'get_content': return getContent(params);
    case 'get_segments': return getSegments(params);
    case 'load_font': return loadFontAction(params);
    case 'set_style_id': return setStyleId(params);
    case 'list_fonts':
    case 'available_fonts':
    case 'list_available_fonts': return listFonts(params);
    case 'set_range_style':
    case 'range_style':
    case 'style_ranges': return setRangeStyle(params);
    default: throw new Error('Unknown text action: ' + action + '. Available: create, set_content, set_font, set_size, set_weight, set_color, set_align, set_spacing, set_line_height, set_letter_spacing, set_decoration, set_case, set_paragraph_spacing, get_content, get_segments, load_font, set_style_id, list_fonts, set_range_style');
  }
}

async function createText(params: any) {
  const content = params.content ?? params.text ?? '';
  const fontFamily = params.fontFamily ?? params.family ?? 'Inter';
  const fontStyle = params.fontStyle ?? params.style ?? 'Regular';
  const fontSize = params.fontSize ?? params.size ?? 16;
  const width = params.width;

  const family = resolveFontFamily(fontFamily);
  await loadFont(family, fontStyle);

  const text = figma.createText();
  text.fontName = { family, style: fontStyle };
  text.fontSize = fontSize;
  text.x = params.x ?? 0;
  text.y = params.y ?? 0;
  if (params.name) text.name = params.name;
  text.characters = String(content);

  if (params.color) {
    const c = parseHexColor(params.color);
    text.fills = [{
      type: 'SOLID',
      color: { r: c.r, g: c.g, b: c.b },
      opacity: c.a,
    }];
  }

  if (width !== undefined) {
    text.resize(width, text.height);
    text.textAutoResize = 'HEIGHT';
  } else {
    text.textAutoResize = 'WIDTH_AND_HEIGHT';
  }

  var container: BaseNode & ChildrenMixin;
  if (params.parentId) {
    const parent = await getParentById(params.parentId);
    if (parent) {
      parent.appendChild(text);
      container = parent;
    } else {
      container = figma.currentPage;
    }
  } else {
    container = figma.currentPage;
  }

  var stableId = await resolveStableId(text, container);

  // Auto-layout child properties
  var parentHasAutoLayout = container && 'layoutMode' in container && (container as any).layoutMode !== 'NONE';

  // Auto-fix: text in auto-layout with HEIGHT resize needs layoutSizingVertical = HUG
  if (parentHasAutoLayout && text.textAutoResize === 'HEIGHT' && params.layoutSizingVertical === undefined) {
    (text as any).layoutSizingVertical = 'HUG';
  }

  if (params.layoutSizingHorizontal !== undefined) (text as any).layoutSizingHorizontal = params.layoutSizingHorizontal;
  if (params.layoutSizingVertical !== undefined) (text as any).layoutSizingVertical = params.layoutSizingVertical;
  if (params.layoutAlign !== undefined) (text as any).layoutAlign = params.layoutAlign;
  if (params.layoutGrow !== undefined) (text as any).layoutGrow = params.layoutGrow;

  // Forward textAlignHorizontal on create
  var textAlign = params.textAlignHorizontal ?? params.textAlign;
  if (textAlign) text.textAlignHorizontal = textAlign;

  // Forward lineHeight on create
  if (params.lineHeight !== undefined) {
    var lhUnit = params.lineHeightUnit ?? 'PERCENT';
    text.lineHeight = lhUnit === 'AUTO'
      ? { unit: 'AUTO' }
      : { value: params.lineHeight, unit: lhUnit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };
  }

  if (params.letterSpacing !== undefined) {
    var lsUnit = params.letterSpacingUnit ?? 'PIXELS';
    text.letterSpacing = { value: params.letterSpacing, unit: lsUnit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };
  }

  if (params.textCase !== undefined) {
    text.textCase = params.textCase;
  }

  if (params.opacity !== undefined) {
    text.opacity = Math.max(0, Math.min(1, params.opacity));
  }

  return { id: stableId, name: text.name, type: text.type, width: text.width, height: text.height, parentId: params.parentId || figma.currentPage.id };
}

async function setContent(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);
  node.characters = String(params.content ?? params.text ?? '');
  return { id: node.id, name: node.name, characters: node.characters };
}

async function setFont(params: any) {
  const node = await getTextNodeById(params.nodeId);
  const fontFamily = params.fontFamily ?? params.family ?? 'Inter';
  const fontStyle = params.fontStyle ?? params.style ?? 'Regular';
  const family = resolveFontFamily(fontFamily);
  await loadFont(family, fontStyle);

  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeFontName(params.rangeStart, params.rangeEnd, { family, style: fontStyle });
  } else {
    await loadNodeFonts(node);
    node.fontName = { family, style: fontStyle };
  }
  return { id: node.id, name: node.name };
}

async function setSize(params: any) {
  const node = await getTextNodeById(params.nodeId);
  const size = params.size ?? params.fontSize;
  await loadNodeFonts(node);

  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeFontSize(params.rangeStart, params.rangeEnd, size);
  } else {
    node.fontSize = size;
  }
  return { id: node.id, name: node.name };
}

async function setWeight(params: any) {
  const node = await getTextNodeById(params.nodeId);
  const weight = params.weight ?? params.fontWeight;

  const weightMap: Record<number, string> = {
    100: 'Thin',
    200: 'Extra Light',
    300: 'Light',
    400: 'Regular',
    500: 'Medium',
    600: 'Semi Bold',
    700: 'Bold',
    800: 'Extra Bold',
    900: 'Black',
  };
  const numericWeight = typeof weight === 'string' ? parseInt(weight, 10) : weight;
  const styleName = typeof weight === 'string' && Number.isNaN(numericWeight)
    ? weight
    : (weightMap[numericWeight] || 'Regular');

  const currentFont = node.fontName;
  const family = currentFont === figma.mixed ? 'Inter' : currentFont.family;
  await loadFont(family, styleName);

  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeFontName(params.rangeStart, params.rangeEnd, { family, style: styleName });
  } else {
    node.fontName = { family, style: styleName };
  }
  return { id: node.id, name: node.name };
}

async function setColor(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const c = parseHexColor(params.color);
  const fill: SolidPaint = {
    type: 'SOLID',
    color: { r: c.r, g: c.g, b: c.b },
    opacity: c.a,
  };

  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeFills(params.rangeStart, params.rangeEnd, [fill]);
  } else {
    node.fills = [fill];
  }
  return { id: node.id, name: node.name };
}

async function setAlign(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const horizontal = params.horizontal ?? params.textAlign ?? params.textAlignHorizontal;
  const vertical = params.vertical ?? params.textAlignVertical;
  if (horizontal) node.textAlignHorizontal = horizontal;
  if (vertical) node.textAlignVertical = vertical;
  return { id: node.id, name: node.name };
}

async function setSpacing(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  if (params.letterSpacing !== undefined) {
    const unit = params.letterSpacingUnit ?? 'PIXELS';
    node.letterSpacing = { value: params.letterSpacing, unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };
  }
  if (params.lineHeight !== undefined) {
    const unit = params.lineHeightUnit ?? 'PIXELS';
    node.lineHeight = unit === 'AUTO'
      ? { unit: 'AUTO' }
      : { value: params.lineHeight, unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };
  }
  if (params.paragraphSpacing !== undefined) {
    node.paragraphSpacing = params.paragraphSpacing;
  }

  return { id: node.id, name: node.name };
}

async function setLineHeight(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const value = params.value ?? params.lineHeight;
  const unit = params.unit ?? params.lineHeightUnit ?? 'PIXELS';
  const lineHeight: LineHeight = unit === 'AUTO'
    ? { unit: 'AUTO' }
    : { value, unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };

  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeLineHeight(params.rangeStart, params.rangeEnd, lineHeight);
  } else {
    node.lineHeight = lineHeight;
  }
  return { id: node.id, name: node.name };
}

async function setLetterSpacing(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const value = params.value ?? params.letterSpacing;
  const unit = params.unit ?? params.letterSpacingUnit ?? 'PIXELS';
  const letterSpacing: LetterSpacing = {
    value,
    unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS',
  };

  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeLetterSpacing(params.rangeStart, params.rangeEnd, letterSpacing);
  } else {
    node.letterSpacing = letterSpacing;
  }
  return { id: node.id, name: node.name };
}

async function setDecoration(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const decoration = params.decoration ?? params.textDecoration;
  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeTextDecoration(params.rangeStart, params.rangeEnd, decoration);
  } else {
    node.textDecoration = decoration;
  }
  return { id: node.id, name: node.name };
}

async function setCase(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const textCase = params.textCase ?? params.case;
  if (params.rangeStart !== undefined && params.rangeEnd !== undefined) {
    node.setRangeTextCase(params.rangeStart, params.rangeEnd, textCase);
  } else {
    node.textCase = textCase;
  }
  return { id: node.id, name: node.name };
}

async function setParagraphSpacing(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);
  node.paragraphSpacing = params.spacing ?? params.paragraphSpacing ?? 0;
  return { id: node.id, name: node.name };
}

async function getContent(params: any) {
  const node = await getTextNodeById(params.nodeId);
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

async function getSegments(params: any) {
  const node = await getTextNodeById(params.nodeId);
  const property = params.property;
  const properties = property
    ? [property]
    : ['fontName', 'fontSize', 'fills', 'textDecoration', 'textCase', 'lineHeight', 'letterSpacing'];
  const segments = node.getStyledTextSegments(properties as any);
  return { id: node.id, name: node.name, segments };
}

async function loadFontAction(params: any) {
  const fontFamily = params.fontFamily ?? params.family ?? 'Inter';
  const fontStyle = params.fontStyle ?? params.style ?? 'Regular';
  const family = resolveFontFamily(fontFamily);
  await loadFont(family, fontStyle);
  return { loaded: true, family, style: fontStyle };
}

async function setStyleId(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await node.setTextStyleIdAsync(params.styleId ?? '');
  return { id: node.id, name: node.name, textStyleId: node.textStyleId };
}

async function listFonts(params: any) {
  const fonts = await figma.listAvailableFontsAsync();
  const familyFilter = params.family ?? params.fontFamily;

  // Group by family
  const familyMap = new Map<string, string[]>();
  for (const font of fonts) {
    if (familyFilter && !font.fontName.family.toLowerCase().includes(familyFilter.toLowerCase())) {
      continue;
    }
    const existing = familyMap.get(font.fontName.family);
    if (existing) {
      existing.push(font.fontName.style);
    } else {
      familyMap.set(font.fontName.family, [font.fontName.style]);
    }
  }

  const result: { family: string; styles: string[] }[] = [];
  for (const [family, styles] of familyMap) {
    result.push({ family, styles });
  }

  return { fonts: result, count: result.length };
}

/**
 * Resolve a range specification to {start, end} indices.
 * Exported for testing.
 */
export function resolveRange(
  range: { match?: string; start?: number; end?: number },
  text: string,
): { start: number; end: number } | null {
  let start: number;
  let end: number;

  if (typeof range.match === 'string') {
    const idx = text.indexOf(range.match);
    if (idx === -1) return null;
    start = idx;
    end = idx + range.match.length;
  } else if (typeof range.start === 'number' && typeof range.end === 'number') {
    start = range.start;
    end = range.end;
  } else {
    return null;
  }

  // Clamp to text bounds
  const len = text.length;
  start = Math.max(0, Math.min(start, len));
  end = Math.max(start, Math.min(end, len));
  if (start === end) return null;

  return { start, end };
}

async function setRangeStyle(params: any) {
  const node = await getTextNodeById(params.nodeId);
  await loadNodeFonts(node);

  const ranges: any[] = params.ranges;
  if (!Array.isArray(ranges) || ranges.length === 0) {
    throw new Error('set_range_style requires a non-empty "ranges" array');
  }

  const text = node.characters;
  let rangesApplied = 0;

  for (const range of ranges) {
    const resolved = resolveRange(range, text);
    if (!resolved) continue;

    const { start, end } = resolved;

    // --- Font name (bold / italic / fontFamily / fontStyle) ---
    if (range.bold !== undefined || range.italic !== undefined || range.fontFamily || range.fontStyle) {
      // Get current font at range start to use as base
      let baseFontName: FontName;
      const rangeFont = node.getRangeFontName(start, start + 1);
      if (rangeFont === figma.mixed) {
        const overall = node.fontName;
        baseFontName = overall === figma.mixed
          ? { family: 'Inter', style: 'Regular' }
          : overall;
      } else {
        baseFontName = rangeFont;
      }

      let family = range.fontFamily
        ? resolveFontFamily(range.fontFamily)
        : baseFontName.family;

      let style: string;
      if (range.fontStyle) {
        // Explicit fontStyle takes priority
        style = range.fontStyle;
      } else {
        // Build style from bold/italic flags
        const isBold = range.bold !== undefined ? range.bold : /Bold/i.test(baseFontName.style);
        const isItalic = range.italic !== undefined ? range.italic : /Italic/i.test(baseFontName.style);

        if (isBold && isItalic) style = 'Bold Italic';
        else if (isBold) style = 'Bold';
        else if (isItalic) style = 'Italic';
        else style = 'Regular';
      }

      await loadFont(family, style);
      node.setRangeFontName(start, end, { family, style });
    }

    // --- Font size ---
    if (range.fontSize !== undefined) {
      node.setRangeFontSize(start, end, range.fontSize);
    }

    // --- Color ---
    if (range.color !== undefined) {
      const c = parseHexColor(range.color);
      const fill: SolidPaint = {
        type: 'SOLID',
        color: { r: c.r, g: c.g, b: c.b },
        opacity: c.a,
      };
      node.setRangeFills(start, end, [fill]);
    }

    // --- Letter spacing ---
    if (range.letterSpacing !== undefined) {
      const unit = range.letterSpacingUnit ?? 'PIXELS';
      node.setRangeLetterSpacing(start, end, {
        value: range.letterSpacing,
        unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS',
      });
    }

    // --- Line height ---
    if (range.lineHeight !== undefined) {
      const unit = range.lineHeightUnit ?? 'PIXELS';
      const lh: LineHeight = unit === 'AUTO'
        ? { unit: 'AUTO' }
        : { value: range.lineHeight, unit: unit === 'PERCENT' ? 'PERCENT' : 'PIXELS' };
      node.setRangeLineHeight(start, end, lh);
    }

    // --- Text decoration ---
    if (range.textDecoration !== undefined) {
      node.setRangeTextDecoration(start, end, range.textDecoration);
    }

    // --- Text case ---
    if (range.textCase !== undefined) {
      node.setRangeTextCase(start, end, range.textCase);
    }

    rangesApplied++;
  }

  return { id: node.id, name: node.name, rangesApplied };
}
