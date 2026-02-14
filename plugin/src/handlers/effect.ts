function parseHexColor(color: any, fallback = { r: 0, g: 0, b: 0, a: 0.25 }) {
  if (color && typeof color === 'object' && typeof color.r === 'number') {
    return {
      r: color.r,
      g: color.g,
      b: color.b,
      a: typeof color.a === 'number' ? color.a : 0.25,
    };
  }

  if (typeof color !== 'string') return fallback;
  const raw = color.trim().replace(/^#/, '');
  if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) return fallback;
  const hex = raw.length === 3 ? raw.split('').map(ch => ch + ch).join('') : raw;
  const hasAlpha = hex.length === 8;
  const n = parseInt(hex, 16);
  if (Number.isNaN(n)) return fallback;

  return {
    r: ((n >> (hasAlpha ? 24 : 16)) & 0xff) / 255,
    g: ((n >> (hasAlpha ? 16 : 8)) & 0xff) / 255,
    b: ((n >> (hasAlpha ? 8 : 0)) & 0xff) / 255,
    a: hasAlpha ? (n & 0xff) / 255 : 1,
  };
}

export async function handleEffect(action: string, params: any): Promise<any> {
  switch (action) {
    case 'set':
    case 'set_effect':
    case 'apply_effects':
    case 'add_effect':
    case 'set_effects': return setEffects(params);
    case 'shadow':
    case 'drop_shadow':
    case 'add_shadow': return addShadow(params);
    case 'blur':
    case 'add_gaussian_blur':
    case 'add_blur': return addBlur(params);
    case 'apply_style': return applyStyle(params);
    case 'delete':
    case 'remove_all':
    case 'clear_effects':
    case 'remove':
    case 'remove_effect': return removeEffect(params);
    case 'list_effects':
    case 'get':
    case 'list':
    case 'get_effects': return getEffects(params);
    default: throw new Error('Unknown effect action: ' + action + '. Available: set_effects, add_shadow, add_blur, apply_style, remove, remove_effect, get_effects');
  }
}

async function getEffectNode(nodeId: string): Promise<SceneNode & BlendMixin> {
  const node = await figma.getNodeByIdAsync(nodeId) as SceneNode | null;
  if (!node || !('effects' in node)) throw new Error(`Node ${nodeId} does not support effects`);
  return node as SceneNode & BlendMixin;
}

async function setEffects(params: any) {
  const { nodeId } = params;
  const node = await getEffectNode(nodeId);
  let effects = params.effects;
  if (typeof effects === 'string') {
    try {
      effects = JSON.parse(effects);
    } catch {
      effects = [];
    }
  }
  if (!Array.isArray(effects)) effects = [];

  node.effects = effects.map((e: any) => buildEffect(e));
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function addShadow(params: any) {
  const { nodeId, radius, spread, visible = true } = params;
  const shadowType = params.type ?? params.shadowType ?? 'DROP_SHADOW';
  const node = await getEffectNode(nodeId);
  const color = parseHexColor(params.color);

  const shadow: DropShadowEffect | InnerShadowEffect = {
    type: shadowType === 'INNER_SHADOW' ? 'INNER_SHADOW' : 'DROP_SHADOW',
    color,
    offset: {
      x: params.offset?.x ?? params.offsetX ?? 0,
      y: params.offset?.y ?? params.offsetY ?? 4,
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
  const { nodeId, radius = 10, visible = true } = params;
  const blurType = params.type ?? params.blurType ?? 'LAYER_BLUR';
  const node = await getEffectNode(nodeId);

  const blur: Effect = blurType === 'BACKGROUND_BLUR'
    ? { type: 'BACKGROUND_BLUR', radius, visible }
    : { type: 'LAYER_BLUR', radius, visible };

  node.effects = [...node.effects, blur];
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function applyStyle(params: any) {
  const { nodeId, styleId } = params;
  const node = await figma.getNodeByIdAsync(nodeId) as SceneNode | null;
  if (!node || !('effectStyleId' in node)) throw new Error(`Node ${nodeId} does not support effect styles`);
  await (node as any).setEffectStyleIdAsync(styleId ?? '');
  return { id: node.id, name: node.name, effectStyleId: (node as GeometryMixin).effectStyleId };
}

async function removeEffect(params: any) {
  const { nodeId, index } = params;
  const node = await getEffectNode(nodeId);

  if (index === undefined) {
    node.effects = [];
    return { id: node.id, name: node.name, effectCount: 0 };
  }

  const effects = [...node.effects];
  if (index < 0 || index >= effects.length) throw new Error(`Effect index ${index} out of range`);
  effects.splice(index, 1);
  node.effects = effects;
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function getEffects(params: any) {
  const { nodeId } = params;
  const node = await getEffectNode(nodeId);
  return {
    id: node.id,
    name: node.name,
    effects: JSON.parse(JSON.stringify(node.effects)),
  };
}

function buildEffect(e: any): Effect {
  switch (e.type) {
    case 'DROP_SHADOW':
    case 'INNER_SHADOW': {
      const color = (typeof e.color === 'string')
        ? parseHexColor(e.color)
        : (e.color || { r: 0, g: 0, b: 0, a: 0.25 });
      const offset = e.offset || { x: e.offsetX ?? 0, y: e.offsetY ?? 4 };
      return {
        type: e.type,
        color: { r: color.r, g: color.g, b: color.b, a: typeof color.a === 'number' ? color.a : 0.25 },
        offset: { x: offset.x ?? 0, y: offset.y ?? 4 },
        radius: e.radius ?? 4,
        spread: e.spread ?? 0,
        visible: e.visible ?? true,
        blendMode: e.blendMode || 'NORMAL',
      };
    }
    case 'LAYER_BLUR':
      return { type: 'LAYER_BLUR', radius: e.radius ?? 10, visible: e.visible ?? true };
    case 'BACKGROUND_BLUR':
      return { type: 'BACKGROUND_BLUR', radius: e.radius ?? 10, visible: e.visible ?? true };
    default:
      throw new Error(`Unknown effect type: ${e.type}`);
  }
}
