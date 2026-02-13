export async function handleEffect(action: string, params: any): Promise<any> {
  switch (action) {
    case 'set_effects': return setEffects(params);
    case 'add_shadow': return addShadow(params);
    case 'add_blur': return addBlur(params);
    case 'remove_effect': return removeEffect(params);
    case 'get_effects': return getEffects(params);
    default: throw new Error(`Unknown effect action: ${action}`);
  }
}

function getEffectNode(nodeId: string): SceneNode & BlendMixin {
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node || !('effects' in node)) throw new Error(`Node ${nodeId} does not support effects`);
  return node as SceneNode & BlendMixin;
}

async function setEffects(params: any) {
  const { nodeId, effects } = params;
  const node = getEffectNode(nodeId);

  node.effects = effects.map((e: any) => buildEffect(e));
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function addShadow(params: any) {
  const { nodeId, type = 'DROP_SHADOW', color, offset, radius, spread, visible = true } = params;
  const node = getEffectNode(nodeId);

  const shadow: DropShadowEffect | InnerShadowEffect = {
    type: type === 'INNER_SHADOW' ? 'INNER_SHADOW' : 'DROP_SHADOW',
    color: {
      r: color?.r ?? 0,
      g: color?.g ?? 0,
      b: color?.b ?? 0,
      a: color?.a ?? 0.25,
    },
    offset: {
      x: offset?.x ?? 0,
      y: offset?.y ?? 4,
    },
    radius: radius ?? 4,
    spread: spread ?? 0,
    visible,
    blendMode: 'NORMAL',
  };

  node.effects = [...node.effects, shadow];
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function addBlur(params: any) {
  const { nodeId, type = 'LAYER_BLUR', radius = 10, visible = true } = params;
  const node = getEffectNode(nodeId);

  let blur: Effect;
  if (type === 'BACKGROUND_BLUR') {
    blur = { type: 'BACKGROUND_BLUR', radius, visible };
  } else {
    blur = { type: 'LAYER_BLUR', radius, visible };
  }

  node.effects = [...node.effects, blur];
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function removeEffect(params: any) {
  const { nodeId, index } = params;
  const node = getEffectNode(nodeId);

  const effects = [...node.effects];
  if (index < 0 || index >= effects.length) throw new Error(`Effect index ${index} out of range`);
  effects.splice(index, 1);
  node.effects = effects;

  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function getEffects(params: any) {
  const { nodeId } = params;
  const node = getEffectNode(nodeId);
  return {
    id: node.id,
    name: node.name,
    effects: JSON.parse(JSON.stringify(node.effects)),
  };
}

function buildEffect(e: any): Effect {
  switch (e.type) {
    case 'DROP_SHADOW':
    case 'INNER_SHADOW':
      return {
        type: e.type,
        color: e.color || { r: 0, g: 0, b: 0, a: 0.25 },
        offset: e.offset || { x: 0, y: 4 },
        radius: e.radius ?? 4,
        spread: e.spread ?? 0,
        visible: e.visible ?? true,
        blendMode: e.blendMode || 'NORMAL',
      };
    case 'LAYER_BLUR':
      return {
        type: 'LAYER_BLUR',
        radius: e.radius ?? 10,
        visible: e.visible ?? true,
      };
    case 'BACKGROUND_BLUR':
      return {
        type: 'BACKGROUND_BLUR',
        radius: e.radius ?? 10,
        visible: e.visible ?? true,
      };
    default:
      throw new Error(`Unknown effect type: ${e.type}`);
  }
}
