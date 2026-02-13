export async function handleLayout(action: string, params: any): Promise<any> {
  switch (action) {
    case 'set_auto_layout': return setAutoLayout(params);
    case 'set_padding': return setPadding(params);
    case 'set_spacing': return setSpacing(params);
    case 'set_alignment': return setAlignment(params);
    case 'set_sizing': return setSizing(params);
    case 'set_constraints': return setConstraints(params);
    case 'set_layout_wrap': return setLayoutWrap(params);
    case 'remove_auto_layout': return removeAutoLayout(params);
    default: throw new Error(`Unknown layout action: ${action}`);
  }
}

function getFrameNode(nodeId: string): FrameNode {
  const node = figma.getNodeById(nodeId);
  if (!node || (node.type !== 'FRAME' && node.type !== 'COMPONENT' && node.type !== 'COMPONENT_SET')) {
    throw new Error(`Node ${nodeId} is not a frame-like node`);
  }
  return node as FrameNode;
}

async function setAutoLayout(params: any) {
  const { nodeId, direction, spacing, padding, alignment, wrap } = params;
  const node = getFrameNode(nodeId);

  node.layoutMode = direction === 'HORIZONTAL' ? 'HORIZONTAL' : 'VERTICAL';

  if (spacing !== undefined) node.itemSpacing = spacing;

  if (padding !== undefined) {
    if (typeof padding === 'number') {
      node.paddingTop = padding;
      node.paddingRight = padding;
      node.paddingBottom = padding;
      node.paddingLeft = padding;
    } else {
      node.paddingTop = padding.top ?? 0;
      node.paddingRight = padding.right ?? 0;
      node.paddingBottom = padding.bottom ?? 0;
      node.paddingLeft = padding.left ?? 0;
    }
  }

  if (alignment) {
    if (alignment.primary) node.primaryAxisAlignItems = alignment.primary;
    if (alignment.counter) node.counterAxisAlignItems = alignment.counter;
  }

  if (wrap !== undefined) {
    node.layoutWrap = wrap ? 'WRAP' : 'NO_WRAP';
  }

  return { id: node.id, name: node.name, layoutMode: node.layoutMode };
}

async function setPadding(params: any) {
  const { nodeId, top, right, bottom, left, all } = params;
  const node = getFrameNode(nodeId);

  if (all !== undefined) {
    node.paddingTop = all;
    node.paddingRight = all;
    node.paddingBottom = all;
    node.paddingLeft = all;
  } else {
    if (top !== undefined) node.paddingTop = top;
    if (right !== undefined) node.paddingRight = right;
    if (bottom !== undefined) node.paddingBottom = bottom;
    if (left !== undefined) node.paddingLeft = left;
  }

  return {
    id: node.id, name: node.name,
    padding: { top: node.paddingTop, right: node.paddingRight, bottom: node.paddingBottom, left: node.paddingLeft },
  };
}

async function setSpacing(params: any) {
  const { nodeId, spacing } = params;
  const node = getFrameNode(nodeId);
  node.itemSpacing = spacing;
  return { id: node.id, name: node.name, itemSpacing: node.itemSpacing };
}

async function setAlignment(params: any) {
  const { nodeId, primary, counter } = params;
  const node = getFrameNode(nodeId);
  if (primary) node.primaryAxisAlignItems = primary;
  if (counter) node.counterAxisAlignItems = counter;
  return { id: node.id, name: node.name, primaryAxisAlignItems: node.primaryAxisAlignItems, counterAxisAlignItems: node.counterAxisAlignItems };
}

async function setSizing(params: any) {
  const { nodeId, primaryAxis, counterAxis, width, height } = params;
  const node = getFrameNode(nodeId);

  if (primaryAxis) node.primaryAxisSizingMode = primaryAxis;
  if (counterAxis) node.counterAxisSizingMode = counterAxis;
  if (width !== undefined) node.resize(width, node.height);
  if (height !== undefined) node.resize(node.width, height);

  return { id: node.id, name: node.name, width: node.width, height: node.height };
}

async function setConstraints(params: any) {
  const { nodeId, horizontal, vertical } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node || !('constraints' in node)) throw new Error(`Node ${nodeId} does not support constraints`);

  const constraints: Constraints = {
    horizontal: horizontal || node.constraints.horizontal,
    vertical: vertical || node.constraints.vertical,
  };
  node.constraints = constraints;

  return { id: node.id, name: node.name, constraints: node.constraints };
}

async function setLayoutWrap(params: any) {
  const { nodeId, wrap } = params;
  const node = getFrameNode(nodeId);
  node.layoutWrap = wrap ? 'WRAP' : 'NO_WRAP';
  return { id: node.id, name: node.name, layoutWrap: node.layoutWrap };
}

async function removeAutoLayout(params: any) {
  const { nodeId } = params;
  const node = getFrameNode(nodeId);
  node.layoutMode = 'NONE';
  return { id: node.id, name: node.name, layoutMode: node.layoutMode };
}
