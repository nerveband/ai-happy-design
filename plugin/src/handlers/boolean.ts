export async function handleBoolean(action: string, params: any): Promise<any> {
  switch (action) {
    case 'union': return booleanOp(params, 'UNION');
    case 'subtract': return booleanOp(params, 'SUBTRACT');
    case 'intersect': return booleanOp(params, 'INTERSECT');
    case 'exclude': return booleanOp(params, 'EXCLUDE');
    case 'flatten': return flattenNode(params);
    default: throw new Error(`Unknown boolean action: ${action}`);
  }
}

async function booleanOp(params: any, operation: 'UNION' | 'SUBTRACT' | 'INTERSECT' | 'EXCLUDE') {
  const { nodeIds, name } = params;
  if (!nodeIds || nodeIds.length < 2) throw new Error('Boolean operations require at least 2 nodes');

  const nodes = nodeIds.map((id: string) => {
    const node = figma.getNodeById(id) as SceneNode;
    if (!node) throw new Error(`Node not found: ${id}`);
    return node;
  });

  // All nodes must share the same parent
  const parent = nodes[0].parent;
  if (!parent) throw new Error('Nodes must have a parent');

  // Figma boolean operations: group nodes then set boolean operation
  const booleanGroup = figma.union(nodes, parent as FrameNode | PageNode);

  // Change the boolean operation type
  if (operation !== 'UNION') {
    switch (operation) {
      case 'SUBTRACT':
        const subtracted = figma.subtract(nodes, parent as FrameNode | PageNode);
        if (name) subtracted.name = name;
        return { id: subtracted.id, name: subtracted.name, type: subtracted.type, operation };
      case 'INTERSECT':
        const intersected = figma.intersect(nodes, parent as FrameNode | PageNode);
        if (name) intersected.name = name;
        return { id: intersected.id, name: intersected.name, type: intersected.type, operation };
      case 'EXCLUDE':
        const excluded = figma.exclude(nodes, parent as FrameNode | PageNode);
        if (name) excluded.name = name;
        return { id: excluded.id, name: excluded.name, type: excluded.type, operation };
    }
  }

  if (name) booleanGroup.name = name;
  return { id: booleanGroup.id, name: booleanGroup.name, type: booleanGroup.type, operation };
}

async function flattenNode(params: any) {
  const { nodeIds } = params;
  if (!nodeIds || nodeIds.length === 0) throw new Error('No nodes specified for flatten');

  const nodes = nodeIds.map((id: string) => {
    const node = figma.getNodeById(id) as SceneNode;
    if (!node) throw new Error(`Node not found: ${id}`);
    return node;
  });

  const flattened = figma.flatten(nodes);
  return { id: flattened.id, name: flattened.name, type: flattened.type };
}
