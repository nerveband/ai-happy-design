import { serializeNodeSummary } from '../utils/serialize';
import { getSceneNodeById, getNodeById, getParentById } from '../utils/getNode';

export async function handleLayer(action: string, params: any): Promise<any> {
  switch (action) {
    case 'reorder':
    case 'order':
    case 'set_z_order':
    case 'set_order': return setOrder(params);
    case 'forward':
    case 'bring_forward': return bringForward(params);
    case 'backward':
    case 'send_backward': return sendBackward(params);
    case 'to_front':
    case 'front':
    case 'bring_to_front': return bringToFront(params);
    case 'to_back':
    case 'back':
    case 'send_to_back': return sendToBack(params);
    case 'group': return groupNodes(params);
    case 'ungroup': return ungroupNodes(params);
    case 'reparent':
    case 'move_to':
    case 'move_to_parent': return moveToParent(params);
    case 'insert_child': return insertChild(params);
    default: throw new Error('Unknown layer action: ' + action + '. Available: set_order, bring_forward, send_backward, bring_to_front, send_to_back, group, ungroup, move_to_parent, insert_child');
  }
}

function getParentChildren(node: SceneNode): readonly SceneNode[] {
  const parent = node.parent;
  if (!parent || !('children' in parent)) throw new Error('Node has no valid parent');
  return parent.children;
}

function getNodeIndex(node: SceneNode): number {
  const children = getParentChildren(node);
  return children.indexOf(node);
}

async function setOrder(params: any) {
  const { nodeId, index } = params;
  const node = await getSceneNodeById(nodeId);
  const parent = node.parent;
  if (!parent || !('insertChild' in parent)) throw new Error('Cannot reorder in this parent');

  const maxIndex = (parent as FrameNode).children.length - 1;
  const targetIndex = Math.max(0, Math.min(index, maxIndex));
  (parent as FrameNode).insertChild(targetIndex, node);

  return { id: node.id, name: node.name, index: targetIndex };
}

async function bringForward(params: any) {
  const { nodeId } = params;
  const node = await getSceneNodeById(nodeId);
  const parent = node.parent;
  if (!parent || !('children' in parent)) throw new Error('Cannot reorder');

  const idx = getNodeIndex(node);
  const maxIdx = (parent as FrameNode).children.length - 1;
  if (idx < maxIdx) {
    (parent as FrameNode).insertChild(idx + 1, node);
  }

  return { id: node.id, name: node.name, index: getNodeIndex(node) };
}

async function sendBackward(params: any) {
  const { nodeId } = params;
  const node = await getSceneNodeById(nodeId);
  const parent = node.parent;
  if (!parent || !('children' in parent)) throw new Error('Cannot reorder');

  const idx = getNodeIndex(node);
  if (idx > 0) {
    (parent as FrameNode).insertChild(idx - 1, node);
  }

  return { id: node.id, name: node.name, index: getNodeIndex(node) };
}

async function bringToFront(params: any) {
  const { nodeId } = params;
  const node = await getSceneNodeById(nodeId);
  const parent = node.parent;
  if (!parent || !('children' in parent)) throw new Error('Cannot reorder');

  const maxIdx = (parent as FrameNode).children.length - 1;
  (parent as FrameNode).insertChild(maxIdx, node);

  return { id: node.id, name: node.name, index: getNodeIndex(node) };
}

async function sendToBack(params: any) {
  const { nodeId } = params;
  const node = await getSceneNodeById(nodeId);
  const parent = node.parent;
  if (!parent || !('insertChild' in parent)) throw new Error('Cannot reorder');

  (parent as FrameNode).insertChild(0, node);

  return { id: node.id, name: node.name, index: 0 };
}

async function groupNodes(params: any) {
  const rawNodeIds = params.nodeIds;
  const nodeIds = Array.isArray(rawNodeIds)
    ? rawNodeIds
    : String(rawNodeIds || '').split(',').map((id: string) => id.trim()).filter(Boolean);
  const { name } = params;
  if (!nodeIds || nodeIds.length === 0) throw new Error('No nodes specified for grouping');

  const nodes = await Promise.all(nodeIds.map(async (id: string) => {
    const n = await getSceneNodeById(id);
    return n;
  }));

  // All nodes must share the same parent
  const parent = nodes[0].parent;
  if (!parent) throw new Error('Nodes must have a parent');

  const group = figma.group(nodes, parent as FrameNode | PageNode);
  if (name) group.name = name;

  return { id: group.id, name: group.name, childCount: group.children.length };
}

async function ungroupNodes(params: any) {
  const { nodeId } = params;
  const node = await getNodeById(nodeId);
  if (node.type !== 'GROUP') throw new Error(`Node ${nodeId} is not a group`);

  const group = node as GroupNode;
  const parent = group.parent;
  if (!parent || !('children' in parent)) throw new Error('Group has no valid parent');

  const children = group.children.slice();
  const ungrouped = children.map(child => serializeNodeSummary(child));

  figma.ungroup(group);

  return { ungrouped, count: ungrouped.length };
}

async function insertChild(params: any) {
  const { parentId, childId, index } = params;
  const parent = await getParentById(parentId);
  if (!parent) throw new Error(`Invalid parent: ${parentId}`);
  const child = await getSceneNodeById(childId);

  (parent as FrameNode).insertChild(index ?? 0, child);

  return { parentId: parent.id, childId: child.id, index: index ?? 0 };
}

async function moveToParent(params: any) {
  const { nodeId, parentId, index } = params;
  const node = await getSceneNodeById(nodeId);
  const parent = await getParentById(parentId);
  if (!parent) throw new Error(`Invalid parent: ${parentId}`);

  if (index !== undefined) {
    (parent as FrameNode).insertChild(index, node);
  } else {
    (parent as FrameNode).appendChild(node);
  }

  return { id: node.id, name: node.name, parentId: parent.id };
}
