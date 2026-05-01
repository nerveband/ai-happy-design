import { getSceneNodeById } from '../utils/getNode';

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
    case 'resolve':
    case 'resolve_value':
    case 'resolve_for_consumer': return resolveForConsumer(params);
    case 'add_mode':
    case 'create_mode': return addMode(params);
    case 'rename_mode': return renameMode(params);
    case 'delete_mode':
    case 'remove_mode': return deleteMode(params);
    case 'extend':
    case 'extend_collection': return extendCollection(params);
    case 'extend_library':
    case 'extend_library_collection': return extendLibraryCollection(params);
    case 'get_values_for_collection':
    case 'values_for_collection': return getValuesForCollection(params);
    case 'remove_mode_override': return removeModeOverride(params);
    case 'remove_collection_overrides':
    case 'remove_overrides_for_variable': return removeCollectionOverrides(params);
    case 'overrides':
    case 'get_overrides': return getOverrides(params);
    default: throw new Error('Unknown variable action: ' + action + '. Available: create, get_all, set_value, bind, unbind, create_collection, get_collections, delete, resolve_for_consumer, add_mode, rename_mode, delete_mode, extend_collection, extend_library_collection, get_values_for_collection, remove_mode_override, remove_collection_overrides, get_overrides');
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

function variablePlanHint(err: any): Error {
  var message = err && err.message ? String(err.message) : String(err);
  if (message.indexOf('Limited to') >= 0 || message.indexOf('outside of enterprise plan') >= 0 || message.indexOf('pricing tier') >= 0) {
    return new Error(message + ' Hint: this Figma plan or file mode does not allow the requested variable collection or mode operation.');
  }
  return new Error(message);
}

async function getCollectionById(collectionId: string) {
  if (!collectionId) throw new Error('collectionId is required');
  var collection = await figma.variables.getVariableCollectionByIdAsync(collectionId);
  if (!collection) throw new Error('Collection not found: ' + collectionId);
  return collection;
}

function serializeCollection(c: any) {
  return {
    id: c.id,
    key: c.key,
    name: c.name,
    remote: c.remote,
    isExtension: c.isExtension,
    parentVariableCollectionId: c.parentVariableCollectionId,
    rootVariableCollectionId: c.rootVariableCollectionId,
    modes: c.modes.map(function(m: any) {
      return { modeId: m.modeId, name: m.name, parentModeId: m.parentModeId };
    }),
    variableIds: c.variableIds,
    variableOverrides: c.variableOverrides,
  };
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

  const node = await getSceneNodeById(nodeId);

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
  const node = await getSceneNodeById(nodeId);

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

async function resolveForConsumer(params: any) {
  var variableId = params.variableId;
  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error('Variable not found: ' + variableId);

  var nodeId = params.nodeId;
  var node = nodeId ? await getSceneNodeById(nodeId) : null;

  // Resolve using the node's bound mode or specified mode
  var modeId = params.modeId;
  if (!modeId && node) {
    // Try to get the mode from the node's explicit variable mode settings
    var collectionId = variable.variableCollectionId;
    var collection = await figma.variables.getVariableCollectionByIdAsync(collectionId);
    if (collection) {
      modeId = collection.defaultModeId || collection.modes[0]?.modeId;
    }
  }

  if (!modeId) {
    // Fall back to first available mode
    var keys = Object.keys(variable.valuesByMode);
    modeId = keys.length > 0 ? keys[0] : undefined;
  }

  if (!modeId) throw new Error('No mode available to resolve variable');

  var rawValue = variable.valuesByMode[modeId];

  // If the raw value is a variable alias, resolve it recursively
  var resolvedValue = rawValue;
  var depth = 0;
  while (resolvedValue && typeof resolvedValue === 'object' && 'type' in resolvedValue && (resolvedValue as any).type === 'VARIABLE_ALIAS' && depth < 10) {
    var aliasId = (resolvedValue as any).id;
    var aliasVar = await figma.variables.getVariableByIdAsync(aliasId);
    if (!aliasVar) break;
    var aliasKeys = Object.keys(aliasVar.valuesByMode);
    if (aliasKeys.length === 0) break;
    resolvedValue = aliasVar.valuesByMode[modeId] || aliasVar.valuesByMode[aliasKeys[0]];
    depth++;
  }

  return {
    id: variable.id,
    name: variable.name,
    resolvedType: variable.resolvedType,
    modeId: modeId,
    rawValue: rawValue,
    resolvedValue: resolvedValue,
  };
}

async function addMode(params: any) {
  var collectionId = params.collectionId;
  var modeName = params.name || params.modeName;
  if (!collectionId) throw new Error('collectionId is required');
  if (!modeName) throw new Error('name is required');

  var collection = await figma.variables.getVariableCollectionByIdAsync(collectionId);
  if (!collection) throw new Error('Collection not found: ' + collectionId);

  collection.addMode(modeName);
  return {
    id: collection.id,
    name: collection.name,
    modes: collection.modes.map(function(m) { return { modeId: m.modeId, name: m.name }; }),
  };
}

async function renameMode(params: any) {
  var collectionId = params.collectionId;
  var modeId = params.modeId;
  var newName = params.name || params.newName;
  if (!collectionId) throw new Error('collectionId is required');
  if (!modeId) throw new Error('modeId is required');
  if (!newName) throw new Error('name is required');

  var collection = await figma.variables.getVariableCollectionByIdAsync(collectionId);
  if (!collection) throw new Error('Collection not found: ' + collectionId);

  collection.renameMode(modeId, newName);
  return {
    id: collection.id,
    name: collection.name,
    modes: collection.modes.map(function(m) { return { modeId: m.modeId, name: m.name }; }),
  };
}

async function deleteMode(params: any) {
  var collectionId = params.collectionId;
  var modeId = params.modeId;
  if (!collectionId) throw new Error('collectionId is required');
  if (!modeId) throw new Error('modeId is required');

  var collection = await figma.variables.getVariableCollectionByIdAsync(collectionId);
  if (!collection) throw new Error('Collection not found: ' + collectionId);

  collection.removeMode(modeId);
  return {
    id: collection.id,
    name: collection.name,
    modes: collection.modes.map(function(m) { return { modeId: m.modeId, name: m.name }; }),
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

async function extendCollection(params: any) {
  var collection = await getCollectionById(params.collectionId);
  var collectionAny = collection as any;
  if (typeof collectionAny.extend !== 'function') {
    throw new Error('Variable collection extend is unavailable in this Figma runtime');
  }
  try {
    var extended = collectionAny.extend(params.name || params.collectionName || (collection.name + ' Extension'));
    return serializeCollection(extended);
  } catch (err: any) {
    throw variablePlanHint(err);
  }
}

async function extendLibraryCollection(params: any) {
  var variablesAny = figma.variables as any;
  if (typeof variablesAny.extendLibraryCollectionByKeyAsync !== 'function') {
    throw new Error('Figma API extendLibraryCollectionByKeyAsync is unavailable in this runtime');
  }
  var key = params.collectionKey || params.key;
  if (!key) throw new Error('collectionKey is required');
  try {
    var extended = await variablesAny.extendLibraryCollectionByKeyAsync(key, params.name || params.collectionName || 'Library Extension');
    return serializeCollection(extended);
  } catch (err: any) {
    throw variablePlanHint(err);
  }
}

async function getValuesForCollection(params: any) {
  var variableId = params.variableId;
  if (!variableId) throw new Error('variableId is required');
  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error('Variable not found: ' + variableId);
  var collection = await getCollectionById(params.collectionId);
  var variableAny = variable as any;
  if (typeof variableAny.valuesByModeForCollectionAsync !== 'function') {
    throw new Error('Variable valuesByModeForCollectionAsync is unavailable in this Figma runtime');
  }
  var values = await variableAny.valuesByModeForCollectionAsync(collection);
  return { id: variable.id, name: variable.name, collectionId: collection.id, valuesByMode: values };
}

async function removeModeOverride(params: any) {
  var variableId = params.variableId;
  var modeId = params.extendedModeId || params.modeId;
  if (!variableId) throw new Error('variableId is required');
  if (!modeId) throw new Error('modeId or extendedModeId is required');
  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error('Variable not found: ' + variableId);
  var variableAny = variable as any;
  if (typeof variableAny.removeOverrideForMode !== 'function') {
    throw new Error('Variable removeOverrideForMode is unavailable in this Figma runtime');
  }
  variableAny.removeOverrideForMode(modeId);
  return { id: variable.id, name: variable.name, removedModeOverride: modeId };
}

async function removeCollectionOverrides(params: any) {
  var collection = await getCollectionById(params.collectionId);
  var collectionAny = collection as any;
  if (typeof collectionAny.removeOverridesForVariable !== 'function') {
    throw new Error('Extended collection removeOverridesForVariable is unavailable in this Figma runtime');
  }
  var variableId = params.variableId;
  if (!variableId) throw new Error('variableId is required');
  var variable = await figma.variables.getVariableByIdAsync(variableId);
  if (!variable) throw new Error('Variable not found: ' + variableId);
  collectionAny.removeOverridesForVariable(variable);
  return { collectionId: collection.id, variableId: variable.id, removed: true };
}

async function getOverrides(params: any) {
  var collection = await getCollectionById(params.collectionId);
  var collectionAny = collection as any;
  return {
    id: collection.id,
    name: collection.name,
    isExtension: collectionAny.isExtension,
    variableOverrides: collectionAny.variableOverrides || {},
  };
}
