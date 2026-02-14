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
    default: throw new Error('Unknown layout action: ' + action + '. Available: set_auto_layout, set_padding, set_spacing, set_alignment, set_sizing, set_constraints, set_layout_wrap, set_wrap, remove_auto_layout');
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
