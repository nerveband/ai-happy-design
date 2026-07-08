import { getNodeById, getSceneNodeById } from '../utils/getNode';
import { loadFont, resolveFontFamily } from '../utils/fonts';

export async function handleLayout(action: string, params: any): Promise<any> {
  switch (action) {
    case 'auto_layout':
    case 'set_layout':
    case 'enable_auto_layout':
    case 'get_layout':
    case 'set_auto_layout': return setAutoLayout(params);
    case 'padding':
    case 'set_padding': return setPadding(params);
    case 'spacing':
    case 'set_spacing': return setSpacing(params);
    case 'alignment':
    case 'align':
    case 'set_alignment': return setAlignment(params);
    case 'sizing':
    case 'set_size':
    case 'set_sizing': return setSizing(params);
    case 'constraints':
    case 'set_constraints': return setConstraints(params);
    case 'wrap':
    case 'set_layout_wrap': return setLayoutWrap(params);
    case 'set_wrap': return setLayoutWrap(params);
    case 'remove_layout':
    case 'remove_auto_layout': return removeAutoLayout(params);
    case 'check_overlaps':
    case 'detect_overlaps':
    case 'overlaps': return checkOverlaps(params);
    case 'grid':
    case 'set_grid':
    case 'set_layout_grid':
    case 'add_grid': return setGrid(params);
    case 'get_grid':
    case 'get_grids':
    case 'get_layout_grids': return getGrids(params);
    case 'remove_grid':
    case 'remove_grids':
    case 'clear_grids':
    case 'remove_layout_grids': return removeGrids(params);
    case 'set_grid_container':
    case 'grid_container': return setGridContainer(params);
    case 'set_grid_tracks':
    case 'grid_tracks': return setGridTracks(params);
    case 'set_grid_child_position':
    case 'grid_child_position': return setGridChildPosition(params);
    case 'get_grid_layout':
    case 'grid_layout': return getGridLayout(params);
    case 'reorder_grid_rows': return reorderGridTracks(params, 'row');
    case 'reorder_grid_columns': return reorderGridTracks(params, 'column');
    case 'pricing_grid': return createPricingGrid(params);
    default: throw new Error('Unknown layout action: ' + action + '. Available: set_auto_layout, set_padding, set_spacing, set_alignment, set_sizing, set_constraints, set_layout_wrap, set_wrap, remove_auto_layout, check_overlaps, set_grid, get_grids, remove_grids');
  }
}

async function getFrameNode(nodeId: string): Promise<FrameNode> {
  const node = await getNodeById(nodeId);
  if (node.type !== 'FRAME' && node.type !== 'COMPONENT' && node.type !== 'COMPONENT_SET') {
    throw new Error(`Node ${nodeId} is not a frame-like node`);
  }
  return node as FrameNode;
}

async function setAutoLayout(params: any) {
  const node = await getFrameNode(params.nodeId);
  const direction = params.direction ?? params.layoutMode;
  if (direction) node.layoutMode = direction === 'NONE' ? 'NONE' : (direction === 'HORIZONTAL' ? 'HORIZONTAL' : 'VERTICAL');

  const spacing = params.spacing ?? params.itemSpacing;
  if (spacing !== undefined) node.itemSpacing = spacing;

  const padding = params.padding;
  if (padding !== undefined) {
    if (typeof padding === 'number') {
      node.paddingTop = padding;
      node.paddingRight = padding;
      node.paddingBottom = padding;
      node.paddingLeft = padding;
    } else {
      node.paddingTop = padding.top ?? node.paddingTop;
      node.paddingRight = padding.right ?? node.paddingRight;
      node.paddingBottom = padding.bottom ?? node.paddingBottom;
      node.paddingLeft = padding.left ?? node.paddingLeft;
    }
  } else {
    if (params.paddingTop !== undefined) node.paddingTop = params.paddingTop;
    if (params.paddingRight !== undefined) node.paddingRight = params.paddingRight;
    if (params.paddingBottom !== undefined) node.paddingBottom = params.paddingBottom;
    if (params.paddingLeft !== undefined) node.paddingLeft = params.paddingLeft;
  }

  const primary = params.primaryAxisAlign ?? params.primaryAxisAlignItems ?? params.alignment?.primary;
  const counter = params.counterAxisAlign ?? params.counterAxisAlignItems ?? params.alignment?.counter;
  if (primary) node.primaryAxisAlignItems = primary;
  if (counter) node.counterAxisAlignItems = counter;

  const wrap = params.wrap ?? params.layoutWrap;
  if (wrap !== undefined) {
    node.layoutWrap = wrap === true || wrap === 'WRAP' ? 'WRAP' : 'NO_WRAP';
  }

  return { id: node.id, name: node.name, layoutMode: node.layoutMode };
}

async function setPadding(params: any) {
  const node = await getFrameNode(params.nodeId);
  const all = params.all ?? params.padding;
  if (all !== undefined && typeof all === 'number') {
    node.paddingTop = all;
    node.paddingRight = all;
    node.paddingBottom = all;
    node.paddingLeft = all;
  } else {
    if (params.top !== undefined) node.paddingTop = params.top;
    if (params.right !== undefined) node.paddingRight = params.right;
    if (params.bottom !== undefined) node.paddingBottom = params.bottom;
    if (params.left !== undefined) node.paddingLeft = params.left;
    if (params.paddingTop !== undefined) node.paddingTop = params.paddingTop;
    if (params.paddingRight !== undefined) node.paddingRight = params.paddingRight;
    if (params.paddingBottom !== undefined) node.paddingBottom = params.paddingBottom;
    if (params.paddingLeft !== undefined) node.paddingLeft = params.paddingLeft;
  }

  return {
    id: node.id,
    name: node.name,
    padding: {
      top: node.paddingTop,
      right: node.paddingRight,
      bottom: node.paddingBottom,
      left: node.paddingLeft,
    },
  };
}

async function setSpacing(params: any) {
  const node = await getFrameNode(params.nodeId);
  node.itemSpacing = params.spacing ?? params.itemSpacing ?? 0;
  return { id: node.id, name: node.name, itemSpacing: node.itemSpacing };
}

async function setAlignment(params: any) {
  const node = await getFrameNode(params.nodeId);
  const primary = params.primary ?? params.primaryAxisAlign ?? params.primaryAxisAlignItems;
  const counter = params.counter ?? params.counterAxisAlign ?? params.counterAxisAlignItems;
  if (primary) node.primaryAxisAlignItems = primary;
  if (counter) node.counterAxisAlignItems = counter;
  return {
    id: node.id,
    name: node.name,
    primaryAxisAlignItems: node.primaryAxisAlignItems,
    counterAxisAlignItems: node.counterAxisAlignItems,
  };
}

async function setSizing(params: any) {
  const node = await getSceneNodeById(params.nodeId);

  // Support layoutSizingHorizontal/layoutSizingVertical on any node (text, frame, etc.)
  const horizontal = params.horizontal ?? params.layoutSizingHorizontal;
  const vertical = params.vertical ?? params.layoutSizingVertical;
  if (horizontal && 'layoutSizingHorizontal' in node) {
    (node as any).layoutSizingHorizontal = horizontal;
  }
  if (vertical && 'layoutSizingVertical' in node) {
    (node as any).layoutSizingVertical = vertical;
  }

  // Legacy frame-level axis sizing (only for frame-like nodes)
  const primaryAxis = params.primaryAxis ?? params.primaryAxisSizing;
  const counterAxis = params.counterAxis ?? params.counterAxisSizing;
  if (primaryAxis && 'primaryAxisSizingMode' in node) {
    (node as any).primaryAxisSizingMode = primaryAxis;
  }
  if (counterAxis && 'counterAxisSizingMode' in node) {
    (node as any).counterAxisSizingMode = counterAxis;
  }

  if (params.width !== undefined && 'resize' in node) {
    (node as any).resize(params.width, (node as any).height);
  }
  if (params.height !== undefined && 'resize' in node) {
    (node as any).resize((node as any).width, params.height);
  }

  return { id: node.id, name: node.name, width: (node as any).width, height: (node as any).height };
}

async function setConstraints(params: any) {
  const node = await getSceneNodeById(params.nodeId);
  if (!('constraints' in node)) throw new Error(`Node ${params.nodeId} does not support constraints`);

  const constraints: Constraints = {
    horizontal: params.horizontal || node.constraints.horizontal,
    vertical: params.vertical || node.constraints.vertical,
  };
  node.constraints = constraints;

  return { id: node.id, name: node.name, constraints: node.constraints };
}

async function setLayoutWrap(params: any) {
  const node = await getFrameNode(params.nodeId);
  const wrap = params.wrap ?? params.layoutWrap;
  node.layoutWrap = wrap === true || wrap === 'WRAP' ? 'WRAP' : 'NO_WRAP';
  return { id: node.id, name: node.name, layoutWrap: node.layoutWrap };
}

async function removeAutoLayout(params: any) {
  const node = await getFrameNode(params.nodeId);
  node.layoutMode = 'NONE';
  return { id: node.id, name: node.name, layoutMode: node.layoutMode };
}

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
  var raw = color.trim().replace(/^#/, '');
  if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) return fallback;
  var hex = raw.length === 3 ? raw.split('').map(function(ch: string) { return ch + ch; }).join('') : raw;
  var hasAlpha = hex.length === 8;
  var n = parseInt(hex, 16);
  if (Number.isNaN(n)) return fallback;
  return {
    r: ((n >> (hasAlpha ? 24 : 16)) & 0xff) / 255,
    g: ((n >> (hasAlpha ? 16 : 8)) & 0xff) / 255,
    b: ((n >> (hasAlpha ? 8 : 0)) & 0xff) / 255,
    a: hasAlpha ? (n & 0xff) / 255 : 1,
  };
}

async function setGrid(params: any) {
  var node = await getFrameNode(params.nodeId);
  var grids = params.grids;

  if (!Array.isArray(grids)) {
    // Single grid shorthand
    grids = [params];
  }

  var layoutGrids: LayoutGrid[] = [];
  for (var i = 0; i < grids.length; i++) {
    var g = grids[i];
    var pattern = (g.pattern || g.type || 'COLUMNS').toUpperCase();

    if (pattern === 'GRID') {
      var gridItem: GridLayoutGrid = {
        pattern: 'GRID',
        sectionSize: g.sectionSize || g.size || 10,
        visible: g.visible !== false,
        color: g.color ? parseHexColor(g.color) : { r: 0.06, g: 0.45, b: 1, a: 0.1 },
      };
      layoutGrids.push(gridItem);
    } else {
      var align = g.alignment || 'STRETCH';
      var colRow: any = {
        pattern: pattern === 'ROWS' ? 'ROWS' : 'COLUMNS',
        visible: g.visible !== false,
        color: g.color ? parseHexColor(g.color) : { r: 0.06, g: 0.45, b: 1, a: 0.1 },
        alignment: align,
        gutterSize: g.gutterSize || g.gutter || 20,
        count: g.count || 12,
      };
      // sectionSize only valid when alignment is not STRETCH
      if (g.sectionSize != null || g.size != null) {
        colRow.sectionSize = g.sectionSize || g.size;
      }
      if (g.offset != null) {
        colRow.offset = g.offset;
      }
      layoutGrids.push(colRow as LayoutGrid);
    }
  }

  if (params.append) {
    var existing = JSON.parse(JSON.stringify(node.layoutGrids));
    for (var j = 0; j < layoutGrids.length; j++) {
      existing.push(layoutGrids[j]);
    }
    node.layoutGrids = existing;
  } else {
    node.layoutGrids = layoutGrids;
  }

  return { id: node.id, name: node.name, gridCount: node.layoutGrids.length };
}

async function getGrids(params: any) {
  var node = await getFrameNode(params.nodeId);
  return {
    id: node.id,
    name: node.name,
    layoutGrids: JSON.parse(JSON.stringify(node.layoutGrids)),
  };
}

async function removeGrids(params: any) {
  var node = await getFrameNode(params.nodeId);
  if (params.index !== undefined) {
    var grids = JSON.parse(JSON.stringify(node.layoutGrids));
    if (params.index < 0 || params.index >= grids.length) {
      throw new Error('Grid index ' + params.index + ' out of range');
    }
    grids.splice(params.index, 1);
    node.layoutGrids = grids;
  } else {
    node.layoutGrids = [];
  }
  return { id: node.id, name: node.name, gridCount: node.layoutGrids.length };
}

function requireGridContainer(node: any, nodeId: string) {
  var required = ['layoutMode', 'gridRowCount', 'gridColumnCount', 'gridRowGap', 'gridColumnGap'];
  for (var i = 0; i < required.length; i++) {
    if (!(required[i] in node)) {
      throw new Error('Node ' + nodeId + ' does not support Figma grid layout containers');
    }
  }
}

function applyTrackSizes(tracks: any, sizes: any) {
  if (sizes == null) return;
  var parsed = typeof sizes === 'string' ? JSON.parse(sizes) : sizes;
  if (!Array.isArray(parsed)) throw new Error('Grid track sizes must be an array');
  for (var i = 0; i < parsed.length && i < tracks.length; i++) {
    var item = parsed[i];
    if (!item || typeof item !== 'object') continue;
    if (item.type) tracks[i].type = String(item.type).toUpperCase();
    if (item.value != null) tracks[i].value = item.value;
  }
}

async function setGridContainer(params: any) {
  var node = await getFrameNode(params.nodeId);
  var gridNode = node as any;
  requireGridContainer(gridNode, params.nodeId);

  gridNode.layoutMode = 'GRID';
  if (params.gridRowCount != null) gridNode.gridRowCount = params.gridRowCount;
  if (params.gridColumnCount != null) gridNode.gridColumnCount = params.gridColumnCount;
  if (params.gridRowGap != null) gridNode.gridRowGap = params.gridRowGap;
  if (params.gridColumnGap != null) gridNode.gridColumnGap = params.gridColumnGap;
  applyTrackSizes(gridNode.gridRowSizes, params.gridRowsSizing || params.gridRowSizes);
  applyTrackSizes(gridNode.gridColumnSizes, params.gridColumnsSizing || params.gridColumnSizes);

  return getGridLayout(params);
}

async function setGridTracks(params: any) {
  var node = await getFrameNode(params.nodeId);
  var gridNode = node as any;
  requireGridContainer(gridNode, params.nodeId);
  if (gridNode.layoutMode !== 'GRID') gridNode.layoutMode = 'GRID';

  applyTrackSizes(gridNode.gridRowSizes, params.gridRowsSizing || params.gridRowSizes);
  applyTrackSizes(gridNode.gridColumnSizes, params.gridColumnsSizing || params.gridColumnSizes);
  return getGridLayout(params);
}

async function setGridChildPosition(params: any) {
  var node = await getSceneNodeById(params.nodeId);
  var gridChild = node as any;
  if (typeof gridChild.setGridChildPosition !== 'function') {
    throw new Error('Node ' + params.nodeId + ' does not support Figma grid child positioning');
  }

  var row = params.gridRowAnchorIndex ?? params.rowIndex ?? params.row ?? 0;
  var column = params.gridColumnAnchorIndex ?? params.columnIndex ?? params.column ?? 0;
  gridChild.setGridChildPosition(row, column);
  if (params.gridRowSpan != null) gridChild.gridRowSpan = params.gridRowSpan;
  if (params.gridColumnSpan != null) gridChild.gridColumnSpan = params.gridColumnSpan;
  if (params.gridChildHorizontalAlign != null && 'gridChildHorizontalAlign' in gridChild) {
    gridChild.gridChildHorizontalAlign = params.gridChildHorizontalAlign;
  }
  if (params.gridChildVerticalAlign != null && 'gridChildVerticalAlign' in gridChild) {
    gridChild.gridChildVerticalAlign = params.gridChildVerticalAlign;
  }

  return {
    id: node.id,
    name: node.name,
    gridRowAnchorIndex: gridChild.gridRowAnchorIndex,
    gridColumnAnchorIndex: gridChild.gridColumnAnchorIndex,
    gridRowSpan: gridChild.gridRowSpan,
    gridColumnSpan: gridChild.gridColumnSpan,
    gridChildHorizontalAlign: gridChild.gridChildHorizontalAlign,
    gridChildVerticalAlign: gridChild.gridChildVerticalAlign,
  };
}

async function getGridLayout(params: any) {
  var node = await getFrameNode(params.nodeId);
  var gridNode = node as any;
  requireGridContainer(gridNode, params.nodeId);
  return {
    id: node.id,
    name: node.name,
    layoutMode: gridNode.layoutMode,
    gridRowCount: gridNode.gridRowCount,
    gridColumnCount: gridNode.gridColumnCount,
    gridRowGap: gridNode.gridRowGap,
    gridColumnGap: gridNode.gridColumnGap,
    gridRowSizes: JSON.parse(JSON.stringify(gridNode.gridRowSizes || [])),
    gridColumnSizes: JSON.parse(JSON.stringify(gridNode.gridColumnSizes || [])),
    gridAutoFlow: gridNode.gridAutoFlow,
    gridAutoRows: JSON.parse(JSON.stringify(gridNode.gridAutoRows || null)),
    gridAutoColumns: JSON.parse(JSON.stringify(gridNode.gridAutoColumns || null)),
    gridColumnCountSizingMode: gridNode.gridColumnCountSizingMode,
    gridRowCountSizingMode: gridNode.gridRowCountSizingMode,
  };
}

async function reorderGridTracks(params: any, axis: string) {
  var node = await getFrameNode(params.nodeId);
  var gridNode = node as any;
  requireGridContainer(gridNode, params.nodeId);
  var field = axis === 'row' ? 'gridRowSizes' : 'gridColumnSizes';
  var tracks = gridNode[field];
  if (!Array.isArray(tracks)) {
    throw new Error('Grid ' + axis + ' track reordering is unavailable in this Figma runtime');
  }
  var fromIndex = params.fromIndex;
  var toIndex = params.toIndex;
  if (fromIndex < 0 || fromIndex >= tracks.length || toIndex < 0 || toIndex >= tracks.length) {
    throw new Error('Grid ' + axis + ' reorder index out of range');
  }
  var next = tracks.slice();
  var moved = next.splice(fromIndex, 1)[0];
  next.splice(toIndex, 0, moved);
  gridNode[field] = next;
  return getGridLayout(params);
}

function fontStyleForWeight(weight: any, fallback: string) {
  if (typeof weight === 'string' && weight.trim() && Number.isNaN(parseInt(weight, 10))) return weight;
  var n = typeof weight === 'string' ? parseInt(weight, 10) : weight;
  var weightMap: Record<number, string> = {
    100: 'Thin', 200: 'Extra Light', 300: 'Light', 400: 'Regular', 500: 'Medium',
    600: 'Semi Bold', 700: 'Bold', 800: 'Extra Bold', 900: 'Black',
  };
  return weightMap[n] || fallback;
}

async function addGridText(parent: BaseNode & ChildrenMixin, name: string, content: string, style: any, width: number) {
  var family = resolveFontFamily(style.fontFamily || 'Inter');
  var fontStyle = style.fontStyle || fontStyleForWeight(style.fontWeight, 'Regular');
  await loadFont(family, fontStyle);
  var text = figma.createText();
  text.name = name;
  text.fontName = { family, style: fontStyle };
  text.fontSize = style.fontSize || 16;
  text.characters = String(content || '');
  text.resize(width, text.height);
  text.textAutoResize = 'HEIGHT';
  if (style.lineHeight !== undefined) {
    var lh = style.lineHeight <= 4 ? style.lineHeight * 100 : style.lineHeight;
    text.lineHeight = { value: lh, unit: 'PERCENT' };
  }
  if (style.letterSpacing !== undefined) text.letterSpacing = { value: style.letterSpacing, unit: 'PIXELS' };
  if (style.textTransform === 'uppercase') text.textCase = 'UPPER';
  if (style.color) {
    var c = parseHexColor(style.color);
    text.fills = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
  }
  parent.appendChild(text);
  return text;
}

function mergeStyle(base: any, override: any) {
  var out: any = {};
  base = base || {};
  override = override || {};
  for (var key in base) out[key] = base[key];
  for (var key2 in override) out[key2] = override[key2];
  return out;
}

function setFrameFill(frame: FrameNode, color: any) {
  if (color === 'transparent' || color === false) {
    frame.fills = [];
    return;
  }
  if (!color) return;
  var c = parseHexColor(color);
  frame.fills = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
}

async function createPricingGrid(params: any) {
  var parent = params.parentId ? await getFrameNode(params.parentId) : figma.currentPage;
  var cards = params.cards || [];
  if (!Array.isArray(cards) || cards.length === 0) throw new Error('cards array is required');
  var x = params.x || 0;
  var y = params.y || 0;
  var width = params.width || 1000;
  var columns = params.columns || 4;
  var gap = params.gap == null ? 28 : params.gap;
  var rowGap = params.rowGap == null ? 42 : params.rowGap;
  var cardWidth = (width - gap * (columns - 1)) / columns;
  var fontFamily = params.fontFamily || 'Inter';
  var gold = params.gold || '#E5AD43';
  var cream = params.cream || '#F8F4EC';
  var rule = params.rule || '#D89E36';
  var titleStyle = mergeStyle({ fontFamily: fontFamily, fontSize: 36, fontWeight: 700, lineHeight: 0.95, color: gold, letterSpacing: 1.5, textTransform: 'uppercase' }, params.titleStyle);
  var priceStyle = mergeStyle({ fontFamily: fontFamily, fontSize: 24, fontWeight: 700, lineHeight: 1.05, color: cream }, params.priceStyle);
  var bodyStyle = mergeStyle({ fontFamily: fontFamily, fontSize: 19, fontWeight: 400, lineHeight: 1.22, color: cream }, params.bodyStyle);
  var noteStyle = mergeStyle({ fontFamily: fontFamily, fontSize: 17, fontWeight: 500, lineHeight: 1.2, color: cream }, params.noteStyle);
  var created: any[] = [];
  var rowHeights: number[] = [];

  for (var i = 0; i < cards.length; i++) {
    var row = Math.floor(i / columns);
    var col = i % columns;
    var card = cards[i];
    var cx = x + col * (cardWidth + gap);
    var cy = y;
    for (var r = 0; r < row; r++) cy += (rowHeights[r] || 286) + rowGap;
    var frame = figma.createFrame();
    frame.name = card.name || String(card.tier || 'Pricing card');
    frame.x = cx;
    frame.y = cy;
    frame.resize(cardWidth, card.minHeight || params.cardMinHeight || 286);
    frame.layoutMode = 'VERTICAL';
    frame.primaryAxisSizingMode = 'AUTO';
    frame.counterAxisSizingMode = 'FIXED';
    frame.itemSpacing = card.gap == null ? 7 : card.gap;
    frame.paddingLeft = card.paddingLeft == null ? 0 : card.paddingLeft;
    frame.paddingRight = card.paddingRight == null ? 18 : card.paddingRight;
    frame.paddingTop = card.paddingTop == null ? 10 : card.paddingTop;
    frame.paddingBottom = card.paddingBottom == null ? 10 : card.paddingBottom;
    setFrameFill(frame, card.background || 'transparent');
    var stroke = parseHexColor(card.rule || rule);
    frame.strokes = [{ type: 'SOLID', color: { r: stroke.r, g: stroke.g, b: stroke.b }, opacity: stroke.a }];
    frame.strokeLeftWeight = card.hideLeftRule ? 0 : (card.ruleWidth || 2);
    frame.strokeRightWeight = 0;
    frame.strokeTopWeight = 0;
    frame.strokeBottomWeight = 0;
    parent.appendChild(frame);
    var innerWidth = cardWidth - frame.paddingLeft - frame.paddingRight;
    var heading = card.heading || card.title || card.tier;
    if (heading) await addGridText(frame, frame.name + ' title', String(heading), titleStyle, innerWidth);
    if (card.price) await addGridText(frame, frame.name + ' price', String(card.price), priceStyle, innerWidth);
    var benefits = card.benefits || card.bullets || [];
    if (Array.isArray(benefits) && benefits.length > 0) {
      await addGridText(frame, frame.name + ' benefits', benefits.map(function(b: any) { return '• ' + String(b); }).join('\n'), bodyStyle, innerWidth);
    }
    if (card.eligibility) {
      var eligibility = String(card.eligibility);
      if (eligibility.toLowerCase().indexOf('eligibility:') !== 0) eligibility = 'Eligibility: ' + eligibility;
      await addGridText(frame, frame.name + ' eligibility', eligibility, noteStyle, innerWidth);
    }
    rowHeights[row] = Math.max(rowHeights[row] || 0, frame.height);
    created.push({ id: frame.id, name: frame.name, x: frame.x, y: frame.y, width: frame.width, height: frame.height });
  }
  return { cards: created, count: created.length, rows: rowHeights.length, rowHeights: rowHeights };
}

async function checkOverlaps(params: any) {
  const node = await getFrameNode(params.nodeId);
  const children = node.children;
  const overlaps: any[] = [];

  for (let i = 0; i < children.length; i++) {
    for (let j = i + 1; j < children.length; j++) {
      const a = children[i];
      const b = children[j];

      // Skip absolute-positioned children (they're intentionally overlapping)
      if ('layoutPositioning' in a && (a as any).layoutPositioning === 'ABSOLUTE') continue;
      if ('layoutPositioning' in b && (b as any).layoutPositioning === 'ABSOLUTE') continue;

      const ax1 = a.x, ay1 = a.y, ax2 = a.x + a.width, ay2 = a.y + a.height;
      const bx1 = b.x, by1 = b.y, bx2 = b.x + b.width, by2 = b.y + b.height;

      const overlapX = Math.max(0, Math.min(ax2, bx2) - Math.max(ax1, bx1));
      const overlapY = Math.max(0, Math.min(ay2, by2) - Math.max(ay1, by1));
      const overlapArea = overlapX * overlapY;

      if (overlapArea > 0) {
        overlaps.push({
          nodeA: { id: a.id, name: a.name, bounds: { x: ax1, y: ay1, width: a.width, height: a.height } },
          nodeB: { id: b.id, name: b.name, bounds: { x: bx1, y: by1, width: b.width, height: b.height } },
          overlapArea,
          overlapRect: {
            x: Math.max(ax1, bx1),
            y: Math.max(ay1, by1),
            width: overlapX,
            height: overlapY,
          },
        });
      }
    }
  }

  return { overlaps, count: overlaps.length, hasOverlaps: overlaps.length > 0 };
}
