/**
 * Serializes a Figma node into a plain JSON-safe object.
 * Handles circular references and special Figma types.
 */

export interface SerializedNode {
  id: string;
  name: string;
  type: string;
  visible: boolean;
  locked: boolean;
  x?: number;
  y?: number;
  width?: number;
  height?: number;
  rotation?: number;
  opacity?: number;
  blendMode?: string;
  fills?: any[];
  strokes?: any[];
  strokeWeight?: number;
  strokeAlign?: string;
  cornerRadius?: number | typeof figma.mixed;
  effects?: any[];
  constraints?: any;
  layoutMode?: string;
  paddingLeft?: number;
  paddingRight?: number;
  paddingTop?: number;
  paddingBottom?: number;
  itemSpacing?: number;
  children?: SerializedNode[];
  characters?: string;
  fontSize?: number | typeof figma.mixed;
  fontName?: any;
  textAlignHorizontal?: string;
  textAlignVertical?: string;
  componentId?: string;
  mainComponentId?: string;
  componentSlotId?: string;
  isSlot?: boolean;
  parentChain?: Array<{id: string, name: string, type: string}>;
}

export function serializeNode(node: SceneNode, depth: number = 0, maxDepth: number = 10): SerializedNode {
  const result: SerializedNode = {
    id: node.id,
    name: node.name,
    type: node.type,
    visible: node.visible,
    locked: node.locked,
  };

  // Position and size
  if ('x' in node) result.x = node.x;
  if ('y' in node) result.y = node.y;
  if ('width' in node) result.width = node.width;
  if ('height' in node) result.height = node.height;
  if ('rotation' in node) result.rotation = node.rotation;

  // Appearance
  if ('opacity' in node) result.opacity = node.opacity;
  if ('blendMode' in node) result.blendMode = node.blendMode;
  if ('fills' in node) result.fills = serializePaints(node.fills);
  if ('strokes' in node) result.strokes = serializePaints(node.strokes);
  if ('strokeWeight' in node) result.strokeWeight = node.strokeWeight as number;
  if ('strokeAlign' in node) result.strokeAlign = node.strokeAlign;
  if ('cornerRadius' in node) result.cornerRadius = node.cornerRadius;
  if ('effects' in node) result.effects = serializeEffects(node.effects);

  // Layout
  if ('constraints' in node) result.constraints = node.constraints;
  if ('layoutMode' in node) {
    const frame = node as FrameNode;
    result.layoutMode = frame.layoutMode;
    result.paddingLeft = frame.paddingLeft;
    result.paddingRight = frame.paddingRight;
    result.paddingTop = frame.paddingTop;
    result.paddingBottom = frame.paddingBottom;
    result.itemSpacing = frame.itemSpacing;
  }

  // Text
  if (node.type === 'TEXT') {
    const textNode = node as TextNode;
    result.characters = textNode.characters;
    result.fontSize = textNode.fontSize;
    result.fontName = textNode.fontName;
    result.textAlignHorizontal = textNode.textAlignHorizontal;
    result.textAlignVertical = textNode.textAlignVertical;
  }

  // Components
  if (node.type === 'COMPONENT') {
    result.componentId = (node as ComponentNode).id;
  }
  if (node.type === 'INSTANCE') {
    const inst = node as InstanceNode;
    result.mainComponentId = inst.mainComponent?.id;
  }
  var anyNode = node as any;
  if (anyNode.componentSlotId != null) result.componentSlotId = anyNode.componentSlotId;
  if (node.type === 'SLOT' || anyNode.isSlot === true) result.isSlot = true;

  // Children (with depth limit)
  if ('children' in node && depth < maxDepth) {
    result.children = (node as ChildrenMixin & SceneNode).children.map(
      (child: SceneNode) => serializeNode(child, depth + 1, maxDepth)
    );
  }

  return result;
}

function serializePaints(paints: readonly Paint[] | typeof figma.mixed): any[] {
  if (paints === figma.mixed) return [];
  return paints.map(p => JSON.parse(JSON.stringify(p)));
}

function serializeEffects(effects: readonly Effect[]): any[] {
  return effects.map(e => JSON.parse(JSON.stringify(e)));
}

/**
 * Lightweight node summary (id + name + type only) for list operations.
 */
export function serializeNodeSummary(node: SceneNode): { id: string; name: string; type: string } {
  return { id: node.id, name: node.name, type: node.type };
}
