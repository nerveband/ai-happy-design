import { serializeNode, serializeNodeSummary } from '../utils/serialize';

export async function handleNode(action: string, params: any): Promise<any> {
  switch (action) {
    case 'get_info': return getInfo(params);
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
    default: throw new Error(`Unknown node action: ${action}`);
  }
}

function getSceneNode(nodeId: string): SceneNode {
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);
  return node;
}

async function getInfo(params: any) {
  const { nodeId, depth } = params;
  const node = getSceneNode(nodeId);
  return serializeNode(node, 0, depth ?? 3);
}

async function createFrame(params: any) {
  const { x = 0, y = 0, width = 100, height = 100, name, parentId, color, clipsContent } = params;
  const frame = figma.createFrame();
  frame.x = x;
  frame.y = y;
  frame.resize(width, height);
  if (name) frame.name = name;
  if (color) {
    frame.fills = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }
  if (clipsContent !== undefined) frame.clipsContent = clipsContent;

  if (parentId) {
    const parent = figma.getNodeById(parentId);
    if (parent && 'appendChild' in parent) (parent as FrameNode).appendChild(frame);
  }

  return { id: frame.id, name: frame.name, type: frame.type };
}

async function moveNode(params: any) {
  const { nodeId, x, y, relative } = params;
  const node = getSceneNode(nodeId);

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
  const node = getSceneNode(nodeId);
  if (!('resize' in node)) throw new Error(`Node ${nodeId} cannot be resized`);
  (node as any).resize(width ?? node.width, height ?? (node as any).height);
  return { id: node.id, name: node.name, width: (node as any).width, height: (node as any).height };
}

async function rotateNode(params: any) {
  const { nodeId, rotation, relative } = params;
  const node = getSceneNode(nodeId);

  if (relative) {
    node.rotation += rotation;
  } else {
    node.rotation = rotation;
  }

  return { id: node.id, name: node.name, rotation: node.rotation };
}

async function setOpacity(params: any) {
  const { nodeId, opacity } = params;
  const node = getSceneNode(nodeId);
  if (!('opacity' in node)) throw new Error(`Node ${nodeId} does not support opacity`);
  node.opacity = Math.max(0, Math.min(1, opacity));
  return { id: node.id, name: node.name, opacity: node.opacity };
}

async function setBlendMode(params: any) {
  const { nodeId, blendMode } = params;
  const node = getSceneNode(nodeId);
  if (!('blendMode' in node)) throw new Error(`Node ${nodeId} does not support blend mode`);
  (node as any).blendMode = blendMode;
  return { id: node.id, name: node.name, blendMode: (node as any).blendMode };
}

async function setVisibility(params: any) {
  const { nodeId, visible } = params;
  const node = getSceneNode(nodeId);
  node.visible = visible;
  return { id: node.id, name: node.name, visible: node.visible };
}

async function setLocked(params: any) {
  const { nodeId, locked } = params;
  const node = getSceneNode(nodeId);
  node.locked = locked;
  return { id: node.id, name: node.name, locked: node.locked };
}

async function renameNode(params: any) {
  const { nodeId, name } = params;
  const node = getSceneNode(nodeId);
  node.name = name;
  return { id: node.id, name: node.name };
}

async function deleteNode(params: any) {
  const { nodeId } = params;
  const node = getSceneNode(nodeId);
  const info = serializeNodeSummary(node);
  node.remove();
  return { deleted: info };
}

async function cloneNode(params: any) {
  const { nodeId, x, y, parentId } = params;
  const node = getSceneNode(nodeId);
  const clone = node.clone();

  if (x !== undefined) clone.x = x;
  if (y !== undefined) clone.y = y;

  if (parentId) {
    const parent = figma.getNodeById(parentId);
    if (parent && 'appendChild' in parent) (parent as FrameNode).appendChild(clone);
  }

  return { id: clone.id, name: clone.name, type: clone.type, sourceId: node.id };
}

async function setCornerRadius(params: any) {
  const { nodeId, radius, topLeft, topRight, bottomRight, bottomLeft } = params;
  const node = getSceneNode(nodeId);
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
