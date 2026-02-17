export interface CompactNode {
  id: string;
  type: string;
  name: string;
  x: number;
  y: number;
  w: number;
  h: number;
  childCount: number;
  parentId: string | null;
  depth: number;
}

export function serializeNodeCompact(
  node: SceneNode | PageNode,
  maxDepth: number = 3,
  parentId: string | null = null,
  currentDepth: number = 0
): CompactNode[] {
  const result: CompactNode[] = [];

  const entry: CompactNode = {
    id: node.id,
    type: node.type,
    name: node.name,
    x: 'x' in node ? (node as any).x : 0,
    y: 'y' in node ? (node as any).y : 0,
    w: 'width' in node ? (node as any).width : 0,
    h: 'height' in node ? (node as any).height : 0,
    childCount: 'children' in node ? (node as any).children.length : 0,
    parentId: parentId,
    depth: currentDepth,
  };
  result.push(entry);

  if ('children' in node && currentDepth < maxDepth) {
    var children = (node as any).children as SceneNode[];
    for (var i = 0; i < children.length; i++) {
      var childNodes = serializeNodeCompact(children[i], maxDepth, node.id, currentDepth + 1);
      for (var j = 0; j < childNodes.length; j++) {
        result.push(childNodes[j]);
      }
    }
  }

  return result;
}
