import { getNodeById, getSceneNodeById } from '../utils/getNode';

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
      var colRow: RowsColsLayoutGrid = {
        pattern: pattern === 'ROWS' ? 'ROWS' : 'COLUMNS',
        sectionSize: g.sectionSize || g.size || 0,
        visible: g.visible !== false,
        color: g.color ? parseHexColor(g.color) : { r: 0.06, g: 0.45, b: 1, a: 0.1 },
        alignment: g.alignment || 'STRETCH',
        gutterSize: g.gutterSize || g.gutter || 20,
        count: g.count || 12,
        offset: g.offset || 0,
      };
      layoutGrids.push(colRow);
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
