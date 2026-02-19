import { serializeNode, serializeNodeSummary } from '../utils/serialize';
import { serializeNodeCompact } from '../utils/serializeCompact';
import { resolveStableId } from '../utils/stableId';
import { getSceneNodeById, getNodeById, getParentById } from '../utils/getNode';

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
  const hex = raw.length === 3 ? raw.split('').map(ch => ch + ch).join('') : raw;
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

export async function handleNode(action: string, params: any): Promise<any> {
  switch (action) {
    case 'info':
    case 'get_node':
    case 'get':
    case 'get_info': return getInfo(params);
    case 'create':
    case 'add_frame':
    case 'new_frame':
    case 'create_frame': return createFrame(params);
    case 'move': return moveNode(params);
    case 'resize': return resizeNode(params);
    case 'rotate': return rotateNode(params);
    case 'set_opacity': return setOpacity(params);
    case 'set_blend_mode': return setBlendMode(params);
    case 'set_visibility': return setVisibility(params);
    case 'set_locked': return setLocked(params);
    case 'rename': return renameNode(params);
    case 'delete': return deleteNode(params);
    case 'clone': return cloneNode(params);
    case 'set_corner_radius': return setCornerRadius(params);
    case 'modify': return modifyNode(params);
    case 'set_mask': return setMask(params);
    case 'children':
    case 'get_children':
    case 'list_children':
    case 'get_tree_nodes':
    case 'tree':
    case 'get_tree': return getTree(params);
    default: throw new Error('Unknown node action: ' + action + '. Available: get_info, create_frame, move, resize, rotate, set_opacity, set_blend_mode, set_visibility, set_locked, rename, delete, clone, set_corner_radius, get_tree, modify, set_mask');
  }
}

function getParentChain(node: BaseNode): Array<{id: string, name: string, type: string}> {
  var chain: Array<{id: string, name: string, type: string}> = [];
  var current = node.parent;
  while (current && current.type !== 'DOCUMENT') {
    chain.push({ id: current.id, name: current.name, type: current.type });
    current = current.parent;
  }
  return chain;
}

async function getInfo(params: any) {
  const { nodeId, depth } = params;
  const node = await getSceneNodeById(nodeId);
  var result = serializeNode(node, 0, depth ?? 3);
  result.parentChain = getParentChain(node);
  return result;
}

async function createFrame(params: any) {
  var _x = params.x;
  var _y = params.y;
  var width = params.width || 100;
  var height = params.height || 100;
  var name = params.name;
  var parentId = params.parentId;
  var color = params.color;
  var clipsContent = params.clipsContent;
  var autoPosition = params.autoPosition;

  // Auto-position root frames to avoid overlapping existing ones
  if (!parentId && (_x === undefined || _x === null || autoPosition === true)) {
    var page = figma.currentPage;
    var maxRight = 0;
    var hasFrames = false;
    for (var i = 0; i < page.children.length; i++) {
      var child = page.children[i];
      var right = child.x + child.width;
      if (right > maxRight) maxRight = right;
      hasFrames = true;
    }
    if (hasFrames && (_x === undefined || _x === null || autoPosition === true)) {
      _x = Math.ceil((maxRight + 80) / 8) * 8; // 80px gap, snapped to 8px grid
    } else if (_x === undefined || _x === null) {
      _x = 0;
    }
  }
  if (_x === undefined || _x === null) _x = 0;
  if (_y === undefined || _y === null) _y = 0;

  var frame = figma.createFrame();
  frame.x = _x;
  frame.y = _y;
  frame.resize(width, height);
  if (name) frame.name = name;
  if (color) {
    var c = parseHexColor(color);
    frame.fills = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
  }
  if (params.noFill === true) {
    frame.fills = [];
  }
  var strokeColor = params.stroke ?? params.strokeColor;
  if (strokeColor) {
    var stroke = parseHexColor(strokeColor);
    frame.strokes = [{ type: 'SOLID', color: { r: stroke.r, g: stroke.g, b: stroke.b }, opacity: stroke.a }];
  }
  var strokeWidth = params.strokeWidth ?? params.strokeWeight;
  if (strokeWidth !== undefined) {
    frame.strokeWeight = strokeWidth;
  }
  if (params.cornerRadius !== undefined) {
    frame.cornerRadius = params.cornerRadius;
  }
  if (params.opacity !== undefined) {
    frame.opacity = Math.max(0, Math.min(1, params.opacity));
  }
  if (clipsContent !== undefined) frame.clipsContent = clipsContent;

  // Auto-layout
  var layoutMode = params.layoutMode ?? params.direction;
  if (layoutMode) {
    frame.layoutMode = layoutMode === 'HORIZONTAL' ? 'HORIZONTAL' : 'VERTICAL';

    var itemSpacing = params.itemSpacing ?? params.spacing;
    if (itemSpacing !== undefined) frame.itemSpacing = itemSpacing;

    var padding = params.padding;
    if (padding !== undefined) {
      if (typeof padding === 'number') {
        frame.paddingTop = padding;
        frame.paddingRight = padding;
        frame.paddingBottom = padding;
        frame.paddingLeft = padding;
      }
    }
    if (params.paddingTop !== undefined) frame.paddingTop = params.paddingTop;
    if (params.paddingRight !== undefined) frame.paddingRight = params.paddingRight;
    if (params.paddingBottom !== undefined) frame.paddingBottom = params.paddingBottom;
    if (params.paddingLeft !== undefined) frame.paddingLeft = params.paddingLeft;

    var primaryAlign = params.primaryAxisAlignItems ?? params.primaryAxisAlign;
    if (primaryAlign) frame.primaryAxisAlignItems = primaryAlign;

    var counterAlign = params.counterAxisAlignItems ?? params.counterAxisAlign;
    if (counterAlign) frame.counterAxisAlignItems = counterAlign;

    var primarySizing = params.primaryAxisSizingMode ?? params.primaryAxisSizing;
    if (primarySizing) frame.primaryAxisSizingMode = primarySizing;

    var counterSizing = params.counterAxisSizingMode ?? params.counterAxisSizing;
    if (counterSizing) frame.counterAxisSizingMode = counterSizing;

    var wrap = params.layoutWrap ?? params.wrap;
    if (wrap !== undefined) {
      frame.layoutWrap = wrap === true || wrap === 'WRAP' ? 'WRAP' : 'NO_WRAP';
    }
  }

  var container: BaseNode & ChildrenMixin;
  if (parentId) {
    const parent = await getParentById(parentId);
    if (parent) {
      parent.appendChild(frame);
      container = parent;
    } else {
      container = figma.currentPage;
    }
  } else {
    container = figma.currentPage;
  }

  var stableId = await resolveStableId(frame, container);
  return { id: stableId, name: frame.name, type: frame.type, x: frame.x, y: frame.y, width: frame.width, height: frame.height, layoutMode: frame.layoutMode, parentId: parentId || figma.currentPage.id };
}

async function moveNode(params: any) {
  const { nodeId, x, y, relative } = params;
  const node = await getSceneNodeById(nodeId);

  if (relative) {
    node.x += x ?? 0;
    node.y += y ?? 0;
  } else {
    if (x !== undefined) node.x = x;
    if (y !== undefined) node.y = y;
  }

  return { id: node.id, name: node.name, x: node.x, y: node.y };
}

async function resizeNode(params: any) {
  const { nodeId, width, height } = params;
  const node = await getSceneNodeById(nodeId);
  if (!('resize' in node)) throw new Error(`Node ${nodeId} cannot be resized`);
  (node as any).resize(width ?? node.width, height ?? (node as any).height);
  return { id: node.id, name: node.name, width: (node as any).width, height: (node as any).height };
}

async function rotateNode(params: any) {
  const { nodeId, rotation, relative } = params;
  const node = await getSceneNodeById(nodeId);
  if (!('rotation' in node)) throw new Error(`Node ${nodeId} does not support rotation`);
  const rotatable = node as SceneNode & { rotation: number };

  if (relative) {
    rotatable.rotation += rotation;
  } else {
    rotatable.rotation = rotation;
  }

  return { id: node.id, name: node.name, rotation: rotatable.rotation };
}

async function setOpacity(params: any) {
  const { nodeId, opacity } = params;
  const node = await getSceneNodeById(nodeId);
  if (!('opacity' in node)) throw new Error(`Node ${nodeId} does not support opacity`);
  node.opacity = Math.max(0, Math.min(1, opacity));
  return { id: node.id, name: node.name, opacity: node.opacity };
}

async function setBlendMode(params: any) {
  const { nodeId, blendMode } = params;
  const node = await getSceneNodeById(nodeId);
  if (!('blendMode' in node)) throw new Error(`Node ${nodeId} does not support blend mode`);
  (node as any).blendMode = blendMode;
  return { id: node.id, name: node.name, blendMode: (node as any).blendMode };
}

async function setVisibility(params: any) {
  const { nodeId, visible } = params;
  const node = await getSceneNodeById(nodeId);
  node.visible = visible;
  return { id: node.id, name: node.name, visible: node.visible };
}

async function setLocked(params: any) {
  const { nodeId, locked } = params;
  const node = await getSceneNodeById(nodeId);
  node.locked = locked;
  return { id: node.id, name: node.name, locked: node.locked };
}

async function renameNode(params: any) {
  const { nodeId, name } = params;
  const node = await getSceneNodeById(nodeId);
  node.name = name;
  return { id: node.id, name: node.name };
}

async function deleteNode(params: any) {
  const { nodeId } = params;
  const node = await getSceneNodeById(nodeId);
  const info = serializeNodeSummary(node);
  node.remove();
  return { deleted: info };
}

async function cloneNode(params: any) {
  const { nodeId, x, y, parentId } = params;
  const node = await getSceneNodeById(nodeId);
  const clone = node.clone();

  if (x !== undefined) clone.x = x;
  if (y !== undefined) clone.y = y;

  var container: BaseNode & ChildrenMixin;
  if (parentId) {
    const parent = await getParentById(parentId);
    if (parent) {
      parent.appendChild(clone);
      container = parent;
    } else {
      container = figma.currentPage;
    }
  } else {
    container = clone.parent as BaseNode & ChildrenMixin || figma.currentPage;
  }

  var stableId = await resolveStableId(clone, container);
  return { id: stableId, name: clone.name, type: clone.type, sourceId: node.id };
}

async function getTree(params: any) {
  const { nodeId, depth, compact } = params;
  const node = await getNodeById(nodeId);

  // If it's a page node, load it first so children are accessible
  if (node.type === 'PAGE') {
    await (node as PageNode).loadAsync();
  }

  if (compact) {
    return serializeNodeCompact(node as SceneNode, depth ?? 3);
  }
  return serializeNode(node as SceneNode, 0, depth ?? 3);
}

async function modifyNode(params: any) {
  var nodeId = params.nodeId;
  if (!nodeId) throw new Error('nodeId is required for modify');
  var node = await getSceneNodeById(nodeId);
  var changed: string[] = [];

  // Position
  if (params.x !== undefined) { node.x = params.x; changed.push('x'); }
  if (params.y !== undefined) { node.y = params.y; changed.push('y'); }

  // Size
  if (params.width !== undefined || params.height !== undefined) {
    if ('resize' in node) {
      (node as any).resize(
        params.width !== undefined ? params.width : (node as any).width,
        params.height !== undefined ? params.height : (node as any).height
      );
      changed.push('size');
    }
  }

  // Fill color
  var color = params.color || params.fillColor;
  if (color) {
    if ('fills' in node) {
      var c = parseHexColor(color);
      (node as any).fills = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
      changed.push('color');
    }
  }

  // Opacity
  if (params.opacity !== undefined && 'opacity' in node) {
    node.opacity = Math.max(0, Math.min(1, params.opacity));
    changed.push('opacity');
  }

  // Corner radius
  if (params.cornerRadius !== undefined && 'cornerRadius' in node) {
    (node as any).cornerRadius = params.cornerRadius;
    changed.push('cornerRadius');
  }

  // Visibility
  if (params.visible !== undefined) {
    node.visible = params.visible;
    changed.push('visible');
  }

  // Name
  if (params.name !== undefined) {
    node.name = params.name;
    changed.push('name');
  }

  // Rotation
  if (params.rotation !== undefined && 'rotation' in node) {
    (node as any).rotation = params.rotation;
    changed.push('rotation');
  }

  // Text content (for text nodes)
  var textContent = params.characters || params.text;
  if (textContent !== undefined && node.type === 'TEXT') {
    var textNode = node as TextNode;
    await figma.loadFontAsync(textNode.fontName as FontName);
    textNode.characters = textContent;
    changed.push('characters');
  }

  // Font size
  if (params.fontSize !== undefined && node.type === 'TEXT') {
    var textNode2 = node as TextNode;
    await figma.loadFontAsync(textNode2.fontName as FontName);
    textNode2.fontSize = params.fontSize;
    changed.push('fontSize');
  }

  // Font family
  if (params.fontFamily !== undefined && node.type === 'TEXT') {
    var textNode3 = node as TextNode;
    var style = params.fontStyle || 'Regular';
    await figma.loadFontAsync({ family: params.fontFamily, style: style });
    textNode3.fontName = { family: params.fontFamily, style: style };
    changed.push('fontFamily');
  }

  // Text alignment
  if (params.textAlignHorizontal !== undefined && node.type === 'TEXT') {
    (node as TextNode).textAlignHorizontal = params.textAlignHorizontal;
    changed.push('textAlignHorizontal');
  }

  // isMask
  if (params.isMask !== undefined) {
    (node as any).isMask = params.isMask;
    changed.push('isMask');
  }

  // Layout sizing
  if (params.layoutSizingHorizontal !== undefined && 'layoutSizingHorizontal' in node) {
    (node as any).layoutSizingHorizontal = params.layoutSizingHorizontal;
    changed.push('layoutSizingHorizontal');
  }
  if (params.layoutSizingVertical !== undefined && 'layoutSizingVertical' in node) {
    (node as any).layoutSizingVertical = params.layoutSizingVertical;
    changed.push('layoutSizingVertical');
  }

  // Blend mode
  if (params.blendMode !== undefined && 'blendMode' in node) {
    (node as any).blendMode = params.blendMode;
    changed.push('blendMode');
  }

  return { id: node.id, name: node.name, type: node.type, modified: changed };
}

async function setMask(params: any) {
  var maskNodeId = params.nodeId;
  var targetIds = params.targetIds;
  if (!maskNodeId) throw new Error('nodeId (mask shape) is required');
  if (!targetIds || !Array.isArray(targetIds) || targetIds.length === 0) {
    throw new Error('targetIds (array of node IDs to mask) is required');
  }

  var maskNode = await getSceneNodeById(maskNodeId);

  // Collect target nodes
  var targets: SceneNode[] = [];
  for (var i = 0; i < targetIds.length; i++) {
    targets.push(await getSceneNodeById(targetIds[i]));
  }

  // All nodes must share the same parent for grouping
  var parent = maskNode.parent;
  if (!parent || !('appendChild' in parent)) {
    throw new Error('Mask node must have a valid parent that supports children');
  }

  // Create a group containing the mask shape + targets
  var allNodes: SceneNode[] = [maskNode];
  for (var j = 0; j < targets.length; j++) {
    allNodes.push(targets[j]);
  }

  var group = figma.group(allNodes, parent as BaseNode & ChildrenMixin);
  group.name = params.name || 'Masked Group';

  // Set isMask on the mask shape (must be first child in group)
  // figma.group preserves order from the array, so maskNode is already first
  (maskNode as any).isMask = true;

  return {
    id: group.id,
    name: group.name,
    type: group.type,
    maskNodeId: maskNode.id,
    maskedNodeIds: targets.map(function(t) { return t.id; }),
  };
}

async function setCornerRadius(params: any) {
  const { nodeId, radius, topLeft, topRight, bottomRight, bottomLeft } = params;
  const node = await getSceneNodeById(nodeId);
  if (!('cornerRadius' in node)) throw new Error(`Node ${nodeId} does not support corner radius`);

  const rectNode = node as RectangleNode;
  if (radius !== undefined) {
    rectNode.cornerRadius = radius;
  } else {
    if (topLeft !== undefined) rectNode.topLeftRadius = topLeft;
    if (topRight !== undefined) rectNode.topRightRadius = topRight;
    if (bottomRight !== undefined) rectNode.bottomRightRadius = bottomRight;
    if (bottomLeft !== undefined) rectNode.bottomLeftRadius = bottomLeft;
  }

  return { id: node.id, name: node.name };
}
