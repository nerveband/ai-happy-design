import { getSceneNodeById } from '../utils/getNode';

export async function handleBoolean(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create_union':
    case 'merge':
    case 'combine':
    case 'union': return booleanOp(params, 'UNION');
    case 'create_subtract':
    case 'minus':
    case 'difference':
    case 'subtract': return booleanOp(params, 'SUBTRACT');
    case 'create_intersect':
    case 'intersect_shapes':
    case 'intersect': return booleanOp(params, 'INTERSECT');
    case 'create_exclude':
    case 'xor':
    case 'exclude': return booleanOp(params, 'EXCLUDE');
    case 'flatten_selection':
    case 'flatten': return flattenNode(params);
    default: throw new Error('Unknown boolean action: ' + action + '. Available: union, subtract, intersect, exclude, flatten');
  }
}

async function booleanOp(params: any, operation: 'UNION' | 'SUBTRACT' | 'INTERSECT' | 'EXCLUDE') {
  const rawNodeIds = params.nodeIds;
  const nodeIds = Array.isArray(rawNodeIds)
    ? rawNodeIds
    : String(rawNodeIds || '').split(',').map((id: string) => id.trim()).filter(Boolean);
  const { name } = params;
  if (!nodeIds || nodeIds.length < 2) throw new Error('Boolean operations require at least 2 nodes');

  const nodes = await Promise.all(nodeIds.map((id: string) => getSceneNodeById(id)));

  // All nodes must share the same parent
  const parent = nodes[0].parent;
  if (!parent) throw new Error('Nodes must have a parent');

  let result: BooleanOperationNode;
  switch (operation) {
    case 'SUBTRACT':
      result = figma.subtract(nodes, parent as FrameNode | PageNode);
      break;
    case 'INTERSECT':
      result = figma.intersect(nodes, parent as FrameNode | PageNode);
      break;
    case 'EXCLUDE':
      result = figma.exclude(nodes, parent as FrameNode | PageNode);
      break;
    case 'UNION':
    default:
      result = figma.union(nodes, parent as FrameNode | PageNode);
      break;
  }
  if (name) result.name = name;
  return { id: result.id, name: result.name, type: result.type, operation };
}

async function flattenNode(params: any) {
  const rawNodeIds = params.nodeIds ?? params.nodeId;
  const nodeIds = Array.isArray(rawNodeIds)
    ? rawNodeIds
    : String(rawNodeIds || '').split(',').map((id: string) => id.trim()).filter(Boolean);
  if (!nodeIds || nodeIds.length === 0) throw new Error('No nodes specified for flatten');

  const nodes = await Promise.all(nodeIds.map((id: string) => getSceneNodeById(id)));

  const flattened = figma.flatten(nodes);
  return { id: flattened.id, name: flattened.name, type: flattened.type };
}
