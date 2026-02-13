import { serializeNodeSummary } from '../utils/serialize';

export async function handleComponent(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create': return createComponent(params);
    case 'create_instance': return createInstance(params);
    case 'create_set': return createComponentSet(params);
    case 'get_local': return getLocalComponents(params);
    case 'get_remote': return getRemoteComponents(params);
    case 'detach_instance': return detachInstance(params);
    case 'reset_instance': return resetInstance(params);
    case 'swap_instance': return swapInstance(params);
    default: throw new Error(`Unknown component action: ${action}`);
  }
}

async function createComponent(params: any) {
  const { nodeId, name } = params;

  let component: ComponentNode;

  if (nodeId) {
    // Convert existing node to component
    const node = figma.getNodeById(nodeId) as SceneNode;
    if (!node) throw new Error(`Node not found: ${nodeId}`);
    component = figma.createComponentFromNode(node);
  } else {
    // Create empty component
    component = figma.createComponent();
    component.resize(params.width ?? 100, params.height ?? 100);
    if (params.x !== undefined) component.x = params.x;
    if (params.y !== undefined) component.y = params.y;
  }

  if (name) component.name = name;

  return { id: component.id, name: component.name, type: component.type };
}

async function createInstance(params: any) {
  const { componentId, x, y, parentId, name } = params;
  const component = figma.getNodeById(componentId);
  if (!component || component.type !== 'COMPONENT') throw new Error(`Not a component: ${componentId}`);

  const instance = (component as ComponentNode).createInstance();
  if (x !== undefined) instance.x = x;
  if (y !== undefined) instance.y = y;
  if (name) instance.name = name;

  if (parentId) {
    const parent = figma.getNodeById(parentId);
    if (parent && 'appendChild' in parent) (parent as FrameNode).appendChild(instance);
  }

  return { id: instance.id, name: instance.name, type: instance.type, mainComponentId: componentId };
}

async function createComponentSet(params: any) {
  const { componentIds, name } = params;
  if (!componentIds || componentIds.length === 0) throw new Error('No components specified');

  const components = componentIds.map((id: string) => {
    const node = figma.getNodeById(id);
    if (!node || node.type !== 'COMPONENT') throw new Error(`Not a component: ${id}`);
    return node as ComponentNode;
  });

  const set = figma.combineAsVariants(components, figma.currentPage);
  if (name) set.name = name;

  return { id: set.id, name: set.name, type: set.type, variantCount: set.children.length };
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
  const node = figma.getNodeById(nodeId);
  if (!node || node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const detached = (node as InstanceNode).detachInstance();
  return { id: detached.id, name: detached.name, type: detached.type };
}

async function resetInstance(params: any) {
  const { nodeId } = params;
  const node = figma.getNodeById(nodeId);
  if (!node || node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  (node as InstanceNode).resetOverrides();
  return { id: node.id, name: node.name };
}

async function swapInstance(params: any) {
  const { nodeId, newComponentId } = params;
  const node = figma.getNodeById(nodeId);
  if (!node || node.type !== 'INSTANCE') throw new Error(`Not an instance: ${nodeId}`);

  const newComponent = figma.getNodeById(newComponentId);
  if (!newComponent || newComponent.type !== 'COMPONENT') throw new Error(`Not a component: ${newComponentId}`);

  (node as InstanceNode).swapComponent(newComponent as ComponentNode);
  return { id: node.id, name: node.name, newComponentId };
}
