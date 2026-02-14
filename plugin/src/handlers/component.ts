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
    default: throw new Error('Unknown component action: ' + action + '. Available: create, create_instance, create_set, get_local, get_remote, get_overrides, set_overrides, detach_instance, reset_instance, swap_instance');
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

  (node as InstanceNode).resetOverrides();
  return { id: node.id, name: node.name };
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
