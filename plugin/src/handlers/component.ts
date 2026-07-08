import { resolveStableId } from '../utils/stableId';
import { getSceneNodeById, getNodeById, getParentById } from '../utils/getNode';

export async function handleComponent(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create': return createComponent(params);
    case 'instantiate':
    case 'create_from_component':
    case 'add_instance':
    case 'create_instance': return createInstance(params);
    case 'create_component_set':
    case 'create_variants':
    case 'create_set': return createComponentSet(params);
    case 'list':
    case 'list_local':
    case 'list_components':
    case 'get_components':
    case 'get_local': return getLocalComponents(params);
    case 'list_remote':
    case 'get_remote_components':
    case 'get_remote': return getRemoteComponents(params);
    case 'overrides':
    case 'list_overrides':
    case 'get_overrides': return getOverrides(params);
    case 'set_overrides': return setOverrides(params);
    case 'detach_instance': return detachInstance(params);
    case 'reset_instance': return resetInstance(params);
    case 'swap_instance': return swapInstance(params);
    case 'get_property_definitions':
    case 'get_properties':
    case 'property_definitions': return getPropertyDefinitions(params);
    case 'set_property_definitions':
    case 'set_properties':
    case 'add_property':
    case 'add_property_definition': return addPropertyDefinition(params);
    case 'delete_property':
    case 'remove_property':
    case 'delete_property_definition': return deletePropertyDefinition(params);
    case 'create_slot': return createSlot(params);
    case 'reset_slot': return resetSlot(params);
    case 'get_slots': return getSlots(params);
    default: throw new Error('Unknown component action: ' + action + '. Available: create, create_instance, create_set, get_local, get_remote, get_overrides, set_overrides, detach_instance, reset_instance, swap_instance, get_property_definitions, add_property_definition, delete_property_definition, create_slot, reset_slot, get_slots');
  }
}

async function createComponent(params: any) {
  const { nodeId, name } = params;

  let component: ComponentNode;

  if (nodeId) {
    // Convert existing node to component
    const node = await getSceneNodeById(nodeId);
    component = figma.createComponentFromNode(node);
  } else {
    // Create empty component
    component = figma.createComponent();
    component.resize(params.width ?? 100, params.height ?? 100);
    if (params.x !== undefined) component.x = params.x;
    if (params.y !== undefined) component.y = params.y;
  }

  if (name) component.name = name;

  var stableId = await resolveStableId(component, figma.currentPage);
  return { id: stableId, name: component.name, type: component.type };
}

async function createInstance(params: any) {
  const { componentId, componentKey, x, y, parentId, name } = params;
  let componentNode: ComponentNode | null = null;

  if (componentId) {
    const found = await getNodeById(componentId);
    if (found.type === 'COMPONENT') {
      componentNode = found as ComponentNode;
    }
  }

  if (!componentNode && componentKey) {
    componentNode = await figma.importComponentByKeyAsync(componentKey);
  }

  if (!componentNode) {
    throw new Error(`Not a component: ${componentId || componentKey}`);
  }

  const instance = componentNode.createInstance();
  if (x !== undefined) instance.x = x;
  if (y !== undefined) instance.y = y;
  if (name) instance.name = name;

  var container: BaseNode & ChildrenMixin;
  if (parentId) {
    const parent = await getParentById(parentId);
    if (parent) {
      (parent as FrameNode).appendChild(instance);
      container = parent as FrameNode;
    } else {
      container = figma.currentPage;
    }
  } else {
    container = figma.currentPage;
  }

  var stableId = await resolveStableId(instance, container);
  return { id: stableId, name: instance.name, type: instance.type, mainComponentId: componentNode.id };
}

async function createComponentSet(params: any) {
  const rawComponentIds = params.componentIds ?? params.nodeIds;
  const componentIds = Array.isArray(rawComponentIds)
    ? rawComponentIds
    : String(rawComponentIds || '').split(',').map((id: string) => id.trim()).filter(Boolean);
  const { name } = params;
  if (!componentIds || componentIds.length === 0) throw new Error('No components specified');

  const components = await Promise.all(componentIds.map(async (id: string) => {
    const node = await getNodeById(id);
    if (node.type !== 'COMPONENT') throw new Error(`Not a component: ${id}`);
    return node as ComponentNode;
  }));

  const set = figma.combineAsVariants(components, figma.currentPage);
  if (name) set.name = name;

  var stableId = await resolveStableId(set as unknown as SceneNode, figma.currentPage);
  return { id: stableId, name: set.name, type: set.type, variantCount: set.children.length };
}

async function getLocalComponents(_params: any) {
  const components = figma.currentPage.findAllWithCriteria({ types: ['COMPONENT'] });
  return {
    components: components.map(c => ({
      id: c.id,
      name: c.name,
      description: (c as ComponentNode).description,
    })),
    count: components.length,
  };
}

async function getRemoteComponents(_params: any) {
  // Note: getRemoteComponentsAsync is a proposed API and may not be available
  try {
    // We can list component instances to discover remote components
    const instances = figma.currentPage.findAllWithCriteria({ types: ['INSTANCE'] });
    const remoteComponents = new Map<string, any>();

    for (const inst of instances) {
      const instance = inst as InstanceNode;
      const main = instance.mainComponent;
      if (main && main.remote) {
        remoteComponents.set(main.id, {
          id: main.id,
          name: main.name,
          description: main.description,
          remote: true,
        });
      }
    }

    return { components: Array.from(remoteComponents.values()), count: remoteComponents.size };
  } catch {
    return { components: [], count: 0, error: 'Remote component listing not available' };
  }
}

async function detachInstance(params: any) {
  const { nodeId } = params;
  const node = await getNodeById(nodeId);
  if (node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const detached = (node as InstanceNode).detachInstance();
  return { id: detached.id, name: detached.name, type: detached.type };
}

async function resetInstance(params: any) {
  const { nodeId } = params;
  const node = await getNodeById(nodeId);
  if (node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const instanceAny = node as any;
  if (typeof instanceAny.removeOverrides === 'function') {
    instanceAny.removeOverrides();
    return { id: node.id, name: node.name, resetMethod: 'removeOverrides' };
  }
  if (typeof instanceAny.resetOverrides === 'function') {
    instanceAny.resetOverrides();
    return { id: node.id, name: node.name, resetMethod: 'resetOverrides' };
  }
  throw new Error('Instance override reset is unavailable in this Figma runtime');
}

async function swapInstance(params: any) {
  const { nodeId, newComponentId } = params;
  const node = await getNodeById(nodeId);
  if (node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const newComponent = await getNodeById(newComponentId);
  if (newComponent.type !== 'COMPONENT') throw new Error(`Not a component: ${newComponentId}`);

  (node as InstanceNode).swapComponent(newComponent as ComponentNode);
  return { id: node.id, name: node.name, newComponentId };
}

async function getPropertyDefinitions(params: any) {
  var nodeId = params.nodeId;
  var node = await getNodeById(nodeId);
  if (node.type !== 'COMPONENT' && node.type !== 'COMPONENT_SET') {
    throw new Error('Node ' + nodeId + ' is not a component or component set');
  }
  var comp = node as ComponentNode | ComponentSetNode;
  return {
    id: comp.id,
    name: comp.name,
    componentPropertyDefinitions: comp.componentPropertyDefinitions,
  };
}

async function addPropertyDefinition(params: any) {
  var nodeId = params.nodeId;
  var node = await getNodeById(nodeId);
  if (node.type !== 'COMPONENT' && node.type !== 'COMPONENT_SET') {
    throw new Error('Node ' + nodeId + ' is not a component or component set');
  }
  var comp = node as ComponentNode | ComponentSetNode;

  var propertyName = params.propertyName || params.name;
  if (!propertyName) throw new Error('propertyName is required');

  var propType = (params.type || 'TEXT').toUpperCase();
  var definition: any = { type: propType };

  if (propType === 'TEXT') {
    definition.defaultValue = params.defaultValue || '';
  } else if (propType === 'BOOLEAN') {
    definition.defaultValue = params.defaultValue !== false;
  } else if (propType === 'VARIANT') {
    definition.defaultValue = params.defaultValue || '';
    if (params.variantOptions) {
      definition.variantOptions = params.variantOptions;
    }
  } else if (propType === 'INSTANCE_SWAP') {
    definition.defaultValue = params.defaultValue || '';
    if (params.preferredValues) {
      definition.preferredValues = params.preferredValues;
    }
  }

  comp.addComponentProperty(propertyName, definition.type, definition.defaultValue);
  return {
    id: comp.id,
    name: comp.name,
    propertyName: propertyName,
    componentPropertyDefinitions: comp.componentPropertyDefinitions,
  };
}

async function deletePropertyDefinition(params: any) {
  var nodeId = params.nodeId;
  var node = await getNodeById(nodeId);
  if (node.type !== 'COMPONENT' && node.type !== 'COMPONENT_SET') {
    throw new Error('Node ' + nodeId + ' is not a component or component set');
  }
  var comp = node as ComponentNode | ComponentSetNode;

  var propertyName = params.propertyName || params.name;
  if (!propertyName) throw new Error('propertyName is required');

  comp.deleteComponentProperty(propertyName);
  return {
    id: comp.id,
    name: comp.name,
    deleted: propertyName,
    componentPropertyDefinitions: comp.componentPropertyDefinitions,
  };
}

async function getOverrides(params: any) {
  const { nodeId } = params;
  const node = await getNodeById(nodeId);
  if (node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const instance = node as InstanceNode;
  return {
    id: instance.id,
    name: instance.name,
    componentProperties: instance.componentProperties,
    mainComponent: instance.mainComponent ? {
      id: instance.mainComponent.id,
      name: instance.mainComponent.name,
    } : null,
  };
}

async function setOverrides(params: any) {
  const { nodeId } = params;
  const node = await getNodeById(nodeId);
  if (node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const instance = node as InstanceNode;
  let overrides = params.overrides;
  if (typeof overrides === 'string') {
    try {
      overrides = JSON.parse(overrides);
    } catch (err: any) {
      throw new Error(`Invalid overrides JSON: ${err.message || String(err)}`);
    }
  }
  if (!overrides || typeof overrides !== 'object') {
    throw new Error('overrides must be an object or JSON object string');
  }

  instance.setProperties(overrides as Record<string, string>);
  return { id: instance.id, name: instance.name, applied: overrides };
}

async function createSlot(params: any) {
  var parent = await getNodeById(params.nodeId);
  var parentAny = parent as any;
  if (typeof parentAny.createSlot !== 'function' && typeof (figma as any).createSlot !== 'function') {
    throw new Error('Component slots are unavailable in this Figma runtime');
  }
  var slot = typeof parentAny.createSlot === 'function'
    ? parentAny.createSlot(params.name)
    : (figma as any).createSlot();
  if (params.name) slot.name = params.name;
  if (params.x != null) slot.x = params.x;
  if (params.y != null) slot.y = params.y;
  if (params.width != null && params.height != null && typeof slot.resize === 'function') {
    slot.resize(params.width, params.height);
  }
  if ('appendChild' in parentAny && slot.parent !== parentAny) {
    parentAny.appendChild(slot);
  }
  return { id: slot.id, name: slot.name, type: slot.type };
}

async function resetSlot(params: any) {
  var slot = await getNodeById(params.nodeId);
  var slotAny = slot as any;
  if (typeof slotAny.resetSlot !== 'function' && typeof slotAny.resetOverrides !== 'function') {
    throw new Error('Component slot reset is unavailable in this Figma runtime');
  }
  if (typeof slotAny.resetSlot === 'function') slotAny.resetSlot();
  else slotAny.resetOverrides();
  return { id: slot.id, name: slot.name, reset: true };
}

async function getSlots(params: any) {
  var node = await getNodeById(params.nodeId);
  var nodeAny = node as any;
  var slots: any[] = [];
  if (typeof nodeAny.getSlots === 'function') {
    slots = nodeAny.getSlots();
  } else if ('findAll' in nodeAny) {
    slots = nodeAny.findAll(function(child: any) {
      return child.type === 'SLOT' || child.isSlot === true || child.componentSlotId != null;
    });
  } else {
    throw new Error('Component slots are unavailable in this Figma runtime');
  }
  return {
    id: node.id,
    name: node.name,
    slots: slots.map(function(slot: any) {
      return { id: slot.id, name: slot.name, type: slot.type, componentSlotId: slot.componentSlotId };
    }),
    count: slots.length,
  };
}
