import { serializeNodeSummary } from '../utils/serialize';

export async function handleLayer(action: string, params: any): Promise<any> {
  switch (action) {
    case 'set_order': return setOrder(params);
    case 'bring_forward': return bringForward(params);
    case 'send_backward': return sendBackward(params);
    case 'bring_to_front': return bringToFront(params);
    case 'send_to_back': return sendToBack(params);
    case 'group': return groupNodes(params);
    case 'ungroup': return ungroupNodes(params);
    case 'move_to_parent': return moveToParent(params);
    default: throw new Error(`Unknown layer action: ${action}`);
  }
}

function getSceneNode(nodeId: string): SceneNode {
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);
  return node;
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
  const node = getSceneNode(nodeId);
  const parent = node.parent;
  if (!parent || !('insertChild' in parent)) throw new Error('Cannot reorder in this parent');

  const maxIndex = (parent as FrameNode).children.length - 1;
  const targetIndex = Math.max(0, Math.min(index, maxIndex));
  (parent as FrameNode).insertChild(targetIndex, node);

  return { id: node.id, name: node.name, index: targetIndex };
}

async function bringForward(params: any) {
  const { nodeId } = params;
  const node = getSceneNode(nodeId);
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
  const node = getSceneNode(nodeId);
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
  const node = getSceneNode(nodeId);
  const parent = node.parent;
  if (!parent || !('children' in parent)) throw new Error('Cannot reorder');

  const maxIdx = (parent as FrameNode).children.length - 1;
  (parent as FrameNode).insertChild(maxIdx, node);

  return { id: node.id, name: node.name, index: getNodeIndex(node) };
}

async function sendToBack(params: any) {
  const { nodeId } = params;
  const node = getSceneNode(nodeId);
  const parent = node.parent;
  if (!parent || !('insertChild' in parent)) throw new Error('Cannot reorder');

  (parent as FrameNode).insertChild(0, node);

  return { id: node.id, name: node.name, index: 0 };
}

async function groupNodes(params: any) {
  const { nodeIds, name } = params;
  if (!nodeIds || nodeIds.length === 0) throw new Error('No nodes specified for grouping');

  const nodes = nodeIds.map((id: string) => {
    const n = figma.getNodeById(id) as SceneNode;
    if (!n) throw new Error(`Node not found: ${id}`);
    return n;
  });

  // All nodes must share the same parent
  const parent = nodes[0].parent;
  if (!parent) throw new Error('Nodes must have a parent');

  const group = figma.group(nodes, parent as FrameNode | PageNode);
  if (name) group.name = name;

  return { id: group.id, name: group.name, childCount: group.children.length };
}

async function ungroupNodes(params: any) {
  const { nodeId } = params;
  const node = figma.getNodeById(nodeId);
  if (!node || node.type !== 'GROUP') throw new Error(`Node ${nodeId} is not a group`);

  const group = node as GroupNode;
  const parent = group.parent;
  if (!parent || !('children' in parent)) throw new Error('Group has no valid parent');

  const children = [...group.children];
  const ungrouped = children.map(child => serializeNodeSummary(child));

  figma.ungroup(group);

  return { ungrouped, count: ungrouped.length };
}

async function moveToParent(params: any) {
  const { nodeId, parentId, index } = params;
  const node = getSceneNode(nodeId);
  const parent = figma.getNodeById(parentId);
  if (!parent || !('appendChild' in parent)) throw new Error(`Invalid parent: ${parentId}`);

  if (index !== undefined) {
    (parent as FrameNode).insertChild(index, node);
  } else {
    (parent as FrameNode).appendChild(node);
  }

  return { id: node.id, name: node.name, parentId: parent.id };
}
