import { serializeNode, serializeNodeSummary } from '../utils/serialize';

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
    case 'children':
    case 'get_children':
    case 'list_children':
    case 'get_tree_nodes':
    case 'tree':
    case 'get_tree': return getTree(params);
    default: throw new Error('Unknown node action: ' + action + '. Available: get_info, create_frame, move, resize, rotate, set_opacity, set_blend_mode, set_visibility, set_locked, rename, delete, clone, set_corner_radius, get_tree');
  }
}

async function getSceneNode(nodeId: string): Promise<SceneNode> {
  const node = await figma.getNodeByIdAsync(nodeId) as SceneNode | null;
  if (!node) throw new Error(`Node not found: ${nodeId}`);
  return node;
}

async function getInfo(params: any) {
  const { nodeId, depth } = params;
  const node = await getSceneNode(nodeId);
  return serializeNode(node, 0, depth ?? 3);
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
  if (clipsContent !== undefined) frame.clipsContent = clipsContent;

  if (parentId) {
    var parent = await figma.getNodeByIdAsync(parentId);
    if (parent && 'appendChild' in parent) (parent as FrameNode).appendChild(frame);
  }

  return { id: frame.id, name: frame.name, type: frame.type, x: frame.x, y: frame.y };
}

async function moveNode(params: any) {
  const { nodeId, x, y, relative } = params;
  const node = await getSceneNode(nodeId);

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
  const node = await getSceneNode(nodeId);
  if (!('resize' in node)) throw new Error(`Node ${nodeId} cannot be resized`);
  (node as any).resize(width ?? node.width, height ?? (node as any).height);
  return { id: node.id, name: node.name, width: (node as any).width, height: (node as any).height };
}

async function rotateNode(params: any) {
  const { nodeId, rotation, relative } = params;
  const node = await getSceneNode(nodeId);

  if (relative) {
    node.rotation += rotation;
  } else {
    node.rotation = rotation;
  }

  return { id: node.id, name: node.name, rotation: node.rotation };
}

async function setOpacity(params: any) {
  const { nodeId, opacity } = params;
  const node = await getSceneNode(nodeId);
  if (!('opacity' in node)) throw new Error(`Node ${nodeId} does not support opacity`);
  node.opacity = Math.max(0, Math.min(1, opacity));
  return { id: node.id, name: node.name, opacity: node.opacity };
}

async function setBlendMode(params: any) {
  const { nodeId, blendMode } = params;
  const node = await getSceneNode(nodeId);
  if (!('blendMode' in node)) throw new Error(`Node ${nodeId} does not support blend mode`);
  (node as any).blendMode = blendMode;
  return { id: node.id, name: node.name, blendMode: (node as any).blendMode };
}

async function setVisibility(params: any) {
  const { nodeId, visible } = params;
  const node = await getSceneNode(nodeId);
  node.visible = visible;
  return { id: node.id, name: node.name, visible: node.visible };
}

async function setLocked(params: any) {
  const { nodeId, locked } = params;
  const node = await getSceneNode(nodeId);
  node.locked = locked;
  return { id: node.id, name: node.name, locked: node.locked };
}

async function renameNode(params: any) {
  const { nodeId, name } = params;
  const node = await getSceneNode(nodeId);
  node.name = name;
  return { id: node.id, name: node.name };
}

async function deleteNode(params: any) {
  const { nodeId } = params;
  const node = await getSceneNode(nodeId);
  const info = serializeNodeSummary(node);
  node.remove();
  return { deleted: info };
}

async function cloneNode(params: any) {
  const { nodeId, x, y, parentId } = params;
  const node = await getSceneNode(nodeId);
  const clone = node.clone();

  if (x !== undefined) clone.x = x;
  if (y !== undefined) clone.y = y;

  if (parentId) {
    const parent = await figma.getNodeByIdAsync(parentId);
    if (parent && 'appendChild' in parent) (parent as FrameNode).appendChild(clone);
  }

  return { id: clone.id, name: clone.name, type: clone.type, sourceId: node.id };
}

async function getTree(params: any) {
  const { nodeId, depth } = params;
  const node = await figma.getNodeByIdAsync(nodeId);
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  // If it's a page node, load it first so children are accessible
  if (node.type === 'PAGE') {
    await (node as PageNode).loadAsync();
  }

  return serializeNode(node as SceneNode, 0, depth ?? 3);
}

async function setCornerRadius(params: any) {
  const { nodeId, radius, topLeft, topRight, bottomRight, bottomLeft } = params;
  const node = await getSceneNode(nodeId);
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
