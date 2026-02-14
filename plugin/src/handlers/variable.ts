export async function handleVariable(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create': return createVariable(params);
    case 'list':
    case 'get':
    case 'list_all':
    case 'list_variables':
    case 'get_all': return getAllVariables(params);
    case 'set':
    case 'update':
    case 'update_value':
    case 'set_value': return setVariableValue(params);
    case 'bind': return bindVariable(params);
    case 'unbind': return unbindVariable(params);
    case 'new_collection':
    case 'add_collection':
    case 'create_collection': return createCollection(params);
    case 'list_collections':
    case 'get_all_collections':
    case 'collections':
    case 'get_collections': return getCollections(params);
    case 'delete': return deleteVariable(params);
    default: throw new Error('Unknown variable action: ' + action + '. Available: create, get_all, set_value, bind, unbind, create_collection, get_collections, delete');
  }
}

async function getOrCreateCollection(collectionId?: string, collectionName?: string) {
  if (collectionId) {
    var byId = await figma.variables.getVariableCollectionByIdAsync(collectionId);
    if (byId) return byId;
  }

  var collections = await figma.variables.getLocalVariableCollectionsAsync();

  if (collectionName) {
    var existing = collections.find(function(c) { return c.name === collectionName; });
    if (existing) return existing;
    return figma.variables.createVariableCollection(collectionName);
  }

  if (collections.length > 0) return collections[0];
  return figma.variables.createVariableCollection('Default Collection');
}

async function createVariable(params: any) {
  const name = params.name;
  const resolvedType = params.resolvedType ?? params.type ?? 'COLOR';
  if (!name) throw new Error('name is required');

  const collection = await getOrCreateCollection(params.collectionId, params.collectionName);
  const variable = figma.variables.createVariable(name, collection, resolvedType);

  if (params.value !== undefined) {
    const defaultMode = collection.defaultModeId || collection.modes[0]?.modeId;
    if (defaultMode) {
      let value = params.value;
      if (typeof value === 'string') {
        try {
          value = JSON.parse(value);
        } catch {
          // keep as string
        }
      }
      variable.setValueForMode(defaultMode, value);
    }
  }

  return { id: variable.id, name: variable.name, resolvedType: variable.resolvedType };
}

async function getAllVariables(_params: any) {
  var variables = await figma.variables.getLocalVariablesAsync();
  return {
    variables: variables.map(function(v) {
      return {
        id: v.id,
        name: v.name,
        resolvedType: v.resolvedType,
        valuesByMode: v.valuesByMode,
      };
    }),
    count: variables.length,
  };
}

async function setVariableValue(params: any) {
  const { variableId, modeId } = params;
  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error(`Variable not found: ${variableId}`);

  const targetMode = modeId || Object.keys(variable.valuesByMode)[0];
  let value = params.value;
  if (typeof value === 'string') {
    try {
      value = JSON.parse(value);
    } catch {
      // keep as plain string
    }
  }

  variable.setValueForMode(targetMode, value);
  return { id: variable.id, name: variable.name, modeId: targetMode };
}

async function bindVariable(params: any) {
  const { nodeId, variableId } = params;
  const field = params.field ?? params.property ?? 'fills';

  const node = await figma.getNodeByIdAsync(nodeId) as SceneNode | null;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error(`Variable not found: ${variableId}`);

  if (field === 'fills' && 'fills' in node) {
    const fills = JSON.parse(JSON.stringify(node.fills));
    if (fills.length > 0) {
      fills[0] = figma.variables.setBoundVariableForPaint(fills[0], 'color', variable);
      (node as MinimalFillsMixin).fills = fills;
    }
  } else if ('setBoundVariable' in node) {
    (node as SceneNode).setBoundVariable(field as VariableBindableNodeField, variable);
  } else {
    throw new Error(`Cannot bind variable to field "${field}" on node ${nodeId}`);
  }

  return { id: node.id, name: node.name, field, variableId };
}

async function unbindVariable(params: any) {
  const { nodeId } = params;
  const field = params.field ?? params.property ?? 'fills';
  const node = await figma.getNodeByIdAsync(nodeId) as SceneNode | null;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  if (field === 'fills' && 'fills' in node) {
    const fills = JSON.parse(JSON.stringify(node.fills));
    if (fills.length > 0) {
      try {
        const variablesAny = figma.variables as any;
        if (typeof variablesAny.removeBoundVariableForPaint === 'function') {
          fills[0] = variablesAny.removeBoundVariableForPaint(fills[0], 'color');
        } else {
          throw new Error('removeBoundVariableForPaint is unavailable');
        }
      } catch {
        // Older API fallback: overwrite paint without variable binding metadata.
        const base = fills[0];
        if (base?.type === 'SOLID') {
          fills[0] = {
            type: 'SOLID',
            color: base.color,
            opacity: base.opacity,
          };
        }
      }
      (node as MinimalFillsMixin).fills = fills;
    }
  } else if ('setBoundVariable' in node) {
    (node as SceneNode).setBoundVariable(field as VariableBindableNodeField, null);
  }

  return { id: node.id, name: node.name, field };
}

async function createCollection(params: any) {
  const { name, modes } = params;
  const collection = figma.variables.createVariableCollection(name);

  let modeNames: string[] = [];
  if (Array.isArray(modes)) {
    modeNames = modes;
  } else if (typeof modes === 'string') {
    modeNames = modes.split(',').map((m: string) => m.trim()).filter(Boolean);
  }

  if (modeNames.length > 0) {
    collection.renameMode(collection.modes[0].modeId, modeNames[0]);
    for (let i = 1; i < modeNames.length; i++) {
      collection.addMode(modeNames[i]);
    }
  }

  return {
    id: collection.id,
    name: collection.name,
    modes: collection.modes.map(m => ({ modeId: m.modeId, name: m.name })),
  };
}

async function getCollections(_params: any) {
  var collections = await figma.variables.getLocalVariableCollectionsAsync();
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
  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error(`Variable not found: ${variableId}`);

  const info = { id: variable.id, name: variable.name };
  variable.remove();
  return { deleted: info };
}
