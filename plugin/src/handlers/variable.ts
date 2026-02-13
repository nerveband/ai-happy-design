export async function handleVariable(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create': return createVariable(params);
    case 'get_all': return getAllVariables(params);
    case 'set_value': return setVariableValue(params);
    case 'bind': return bindVariable(params);
    case 'create_collection': return createCollection(params);
    case 'get_collections': return getCollections(params);
    case 'delete': return deleteVariable(params);
    default: throw new Error(`Unknown variable action: ${action}`);
  }
}

async function createVariable(params: any) {
  const { name, collectionId, resolvedType = 'COLOR' } = params;

  const collection = figma.variables.getVariableCollectionById(collectionId);
  if (!collection) throw new Error(`Collection not found: ${collectionId}`);

  const variable = figma.variables.createVariable(name, collection, resolvedType);

  return { id: variable.id, name: variable.name, resolvedType: variable.resolvedType };
}

async function getAllVariables(_params: any) {
  const variables = figma.variables.getLocalVariables();
  return {
    variables: variables.map(v => ({
      id: v.id,
      name: v.name,
      resolvedType: v.resolvedType,
      valuesByMode: v.valuesByMode,
    })),
    count: variables.length,
  };
}

async function setVariableValue(params: any) {
  const { variableId, modeId, value } = params;
  const variable = figma.variables.getVariableById(variableId);
  if (!variable) throw new Error(`Variable not found: ${variableId}`);

  variable.setValueForMode(modeId, value);
  return { id: variable.id, name: variable.name, modeId };
}

async function bindVariable(params: any) {
  const { nodeId, variableId, field } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  const variable = figma.variables.getVariableById(variableId);
  if (!variable) throw new Error(`Variable not found: ${variableId}`);

  if (field === 'fills' && 'fills' in node) {
    const fills = JSON.parse(JSON.stringify(node.fills));
    if (fills.length > 0) {
      fills[0] = figma.variables.setBoundVariableForPaint(fills[0], 'color', variable);
      (node as MinimalFillsMixin).fills = fills;
    }
  } else if ('boundVariables' in node) {
    (node as SceneNode).setBoundVariable(field as VariableBindableNodeField, variable);
  } else {
    throw new Error(`Cannot bind variable to field "${field}" on node ${nodeId}`);
  }

  return { id: node.id, name: node.name, field, variableId };
}

async function createCollection(params: any) {
  const { name, modes } = params;
  const collection = figma.variables.createVariableCollection(name);

  // Rename default mode if specified
  if (modes && modes.length > 0) {
    collection.renameMode(collection.modes[0].modeId, modes[0]);
    // Add additional modes
    for (let i = 1; i < modes.length; i++) {
      collection.addMode(modes[i]);
    }
  }

  return {
    id: collection.id,
    name: collection.name,
    modes: collection.modes.map(m => ({ modeId: m.modeId, name: m.name })),
  };
}

async function getCollections(_params: any) {
  const collections = figma.variables.getLocalVariableCollections();
  return {
    collections: collections.map(c => ({
      id: c.id,
      name: c.name,
      modes: c.modes.map(m => ({ modeId: m.modeId, name: m.name })),
      variableIds: c.variableIds,
    })),
    count: collections.length,
  };
}

async function deleteVariable(params: any) {
  const { variableId } = params;
  const variable = figma.variables.getVariableById(variableId);
  if (!variable) throw new Error(`Variable not found: ${variableId}`);

  const info = { id: variable.id, name: variable.name };
  variable.remove();
  return { deleted: info };
}
