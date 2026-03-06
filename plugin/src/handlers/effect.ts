import { getSceneNodeById } from '../utils/getNode';
import { sanitizeEffects } from '../utils/sanitizeEffects';

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
  // Generic "add" dispatches based on params.type
  if (action === 'add') {
    var effectType = (params.type || 'DROP_SHADOW').toUpperCase();
    if (effectType === 'DROP_SHADOW' || effectType === 'SHADOW') return addShadow(params);
    if (effectType === 'INNER_SHADOW') return addShadow(Object.assign({}, params, { inner: true }));
    if (effectType === 'LAYER_BLUR' || effectType === 'BACKGROUND_BLUR' || effectType === 'BLUR') return addBlur(params);
    if (effectType === 'NOISE') return addNoise(params);
    if (effectType === 'TEXTURE') return addTexture(params);
    if (effectType === 'GLASS') return applyGlass(params);
    throw new Error('Unknown effect type: ' + params.type + '. Available: DROP_SHADOW, INNER_SHADOW, LAYER_BLUR, BACKGROUND_BLUR, NOISE, TEXTURE, GLASS');
  }
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
    case 'add_noise': return addNoise(params);
    case 'add_texture': return addTexture(params);
    case 'add_glass': return addNativeGlass(params);
    case 'apply_glass':
    case 'glass': return applyGlass(params);
    default: throw new Error('Unknown effect action: ' + action + '. Available: add, set_effects, add_shadow, add_blur, apply_style, remove, remove_effect, get_effects, add_noise, add_texture, apply_glass, add_glass');
  }
}

async function getEffectNode(nodeId: string): Promise<SceneNode & BlendMixin> {
  const node = await getSceneNodeById(nodeId);
  if (!('effects' in node)) throw new Error(`Node ${nodeId} does not support effects`);
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
  var shadowType = params.type ?? params.shadowType ?? 'DROP_SHADOW';
  if (params.inner === true) shadowType = 'INNER_SHADOW';
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

  const nextEffects = sanitizeEffects(node.effects);
  nextEffects.push(shadow);
  node.effects = nextEffects;
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function addBlur(params: any) {
  const { nodeId, radius = 10, visible = true } = params;
  const blurType = params.type ?? params.blurType ?? 'LAYER_BLUR';
  const node = await getEffectNode(nodeId);

  const blur: BlurEffectNormal = blurType === 'BACKGROUND_BLUR'
    ? { type: 'BACKGROUND_BLUR', blurType: 'NORMAL', radius, visible }
    : { type: 'LAYER_BLUR', blurType: 'NORMAL', radius, visible };

  const nextEffects = sanitizeEffects(node.effects);
  nextEffects.push(blur);
  node.effects = nextEffects;
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function applyStyle(params: any) {
  const { nodeId, styleId } = params;
  const node = await getSceneNodeById(nodeId);
  if (!('effectStyleId' in node) || !('setEffectStyleIdAsync' in node)) {
    throw new Error(`Node ${nodeId} does not support effect styles`);
  }
  const styledNode = node as SceneNode & BlendMixin;
  await styledNode.setEffectStyleIdAsync(styleId ?? '');
  return { id: node.id, name: node.name, effectStyleId: styledNode.effectStyleId };
}

async function removeEffect(params: any) {
  const { nodeId, index } = params;
  const node = await getEffectNode(nodeId);

  if (index === undefined) {
    node.effects = [];
    return { id: node.id, name: node.name, effectCount: 0 };
  }

  const effects = sanitizeEffects(node.effects);
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

async function addNoise(params: any) {
  var nodeId = params.nodeId;
  var node = await getEffectNode(nodeId);
  var noiseType = params.noiseType || 'monotone';
  var color = parseHexColor(params.color, { r: 1, g: 1, b: 1, a: 0.25 });
  var noiseSize = params.noiseSize || 100;
  var density = params.density !== undefined ? params.density : 0.3;
  var visible = params.visible !== undefined ? params.visible : true;

  // Build the noise effect object. Figma Beta API uses these as regular effects.
  // Since noise effects may not be in all plugin typings, we construct manually.
  var upperType = noiseType.toUpperCase();
  var noiseEffect: any = {
    type: 'NOISE',
    noiseType: upperType,
    color: color,
    noiseSize: noiseSize,
    density: density,
    visible: visible,
  };

  // DUOTONE has secondaryColor; MULTITONE has opacity
  if (upperType === 'DUOTONE' && params.secondaryColor) {
    noiseEffect.secondaryColor = parseHexColor(params.secondaryColor);
  }
  if (upperType === 'MULTITONE' && params.opacity !== undefined) {
    noiseEffect.opacity = params.opacity;
  }

  // Try with blendMode first (newer Figma versions), fall back without it
  var blendMode = params.blendMode || 'SOFT_LIGHT';
  var nextEffects = sanitizeEffects(node.effects);
  var effectWithBlend: any = {};
  for (var k in noiseEffect) {
    effectWithBlend[k] = noiseEffect[k];
  }
  effectWithBlend.blendMode = blendMode;
  nextEffects.push(effectWithBlend);
  try {
    node.effects = nextEffects;
  } catch (e: any) {
    // blendMode may not be supported in this Figma version — retry without it
    var fallbackEffects = sanitizeEffects(node.effects);
    fallbackEffects.push(noiseEffect);
    try {
      node.effects = fallbackEffects;
    } catch (e2: any) {
      throw new Error('Noise effect failed: ' + (e2.message || String(e2)));
    }
  }
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function addTexture(params: any) {
  var nodeId = params.nodeId;
  var node = await getEffectNode(nodeId);
  var noiseSize = params.noiseSize || 100;
  var radius = params.radius || 0;
  var visible = params.visible !== undefined ? params.visible : true;

  var textureEffect: any = {
    type: 'TEXTURE',
    noiseSize: noiseSize,
    radius: radius,
    clipToShape: params.clipToShape !== undefined ? params.clipToShape : true,
    visible: visible,
  };

  var nextEffects = sanitizeEffects(node.effects);
  nextEffects.push(textureEffect);
  try {
    node.effects = nextEffects;
  } catch (e: any) {
    throw new Error('Texture effect failed: ' + (e.message || String(e)));
  }
  return { id: node.id, name: node.name, effectCount: node.effects.length };
}

async function addNativeGlass(params: any) {
  var nodeId = params.nodeId;
  var node = await getEffectNode(nodeId);

  var glassEffect: any = {
    type: 'GLASS',
    visible: params.visible !== undefined ? params.visible : true,
    lightIntensity: params.lightIntensity !== undefined ? params.lightIntensity : 0.5,
    lightAngle: params.lightAngle !== undefined ? params.lightAngle : 45,
    refraction: params.refraction !== undefined ? params.refraction : 0.5,
    depth: params.depth !== undefined ? params.depth : 1,
    dispersion: params.dispersion !== undefined ? params.dispersion : 0,
    radius: params.radius !== undefined ? params.radius : 0,
  };

  var nextEffects = sanitizeEffects(node.effects);
  nextEffects.push(glassEffect);
  try {
    node.effects = nextEffects;
  } catch (e: any) {
    throw new Error('Native glass effect failed: ' + (e.message || String(e)));
  }
  return {
    id: node.id,
    name: node.name,
    glass: glassEffect,
    effectCount: node.effects.length,
  };
}

async function applyGlass(params: any) {
  var nodeId = params.nodeId;
  var node = await getEffectNode(nodeId);
  var intensity = params.intensity || 'medium';

  // Map intensity presets to native glass params
  var nativeParams: any = { lightIntensity: 0.5, lightAngle: 45, refraction: 0.5, depth: 1, dispersion: 0, radius: 0 };
  switch (intensity) {
    case 'light':
      nativeParams = { lightIntensity: 0.3, lightAngle: 45, refraction: 0.3, depth: 1, dispersion: 0, radius: 8 };
      break;
    case 'heavy':
      nativeParams = { lightIntensity: 0.8, lightAngle: 45, refraction: 0.7, depth: 2, dispersion: 0.2, radius: 16 };
      break;
    default: // medium
      nativeParams = { lightIntensity: 0.5, lightAngle: 45, refraction: 0.5, depth: 1.5, dispersion: 0.1, radius: 12 };
  }

  // Try native GLASS effect first
  var glassEffect: any = {
    type: 'GLASS',
    visible: true,
    lightIntensity: nativeParams.lightIntensity,
    lightAngle: nativeParams.lightAngle,
    refraction: nativeParams.refraction,
    depth: nativeParams.depth,
    dispersion: nativeParams.dispersion,
    radius: nativeParams.radius,
  };

  var nextEffects = sanitizeEffects(node.effects);
  nextEffects.push(glassEffect);
  try {
    node.effects = nextEffects;
    return {
      id: node.id,
      name: node.name,
      glass: { mode: 'native', intensity: intensity, effect: glassEffect },
    };
  } catch (e: any) {
    // Native GLASS not supported — fall back to simulated glass
  }

  // Simulated glass fallback
  var tint = params.tint || '#FFFFFF';
  var tintColor = parseHexColor(tint, { r: 1, g: 1, b: 1, a: 1 });
  var fillOpacity: number, blurRadius: number, strokeOpacity: number;
  switch (intensity) {
    case 'light':
      fillOpacity = 0.08; blurRadius = 20; strokeOpacity = 0.10;
      break;
    case 'heavy':
      fillOpacity = 0.15; blurRadius = 40; strokeOpacity = 0.15;
      break;
    default:
      fillOpacity = 0.10; blurRadius = 30; strokeOpacity = 0.12;
  }

  if ('fills' in node) {
    (node as any).fills = [{
      type: 'SOLID',
      color: { r: tintColor.r, g: tintColor.g, b: tintColor.b },
      opacity: fillOpacity,
    }];
  }

  var simEffects: any[] = sanitizeEffects(node.effects);
  simEffects.push({
    type: 'BACKGROUND_BLUR',
    blurType: 'NORMAL',
    radius: blurRadius,
    visible: true,
  });
  node.effects = simEffects;

  if ('strokes' in node) {
    (node as any).strokes = [{
      type: 'SOLID',
      color: { r: tintColor.r, g: tintColor.g, b: tintColor.b },
      opacity: strokeOpacity,
    }];
    (node as any).strokeWeight = 1;
    if ('strokeAlign' in node) {
      (node as any).strokeAlign = 'INSIDE';
    }
  }

  return {
    id: node.id,
    name: node.name,
    glass: { mode: 'simulated', intensity: intensity, fillOpacity: fillOpacity, blurRadius: blurRadius, strokeOpacity: strokeOpacity },
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
      return { type: 'LAYER_BLUR', blurType: 'NORMAL', radius: e.radius ?? 10, visible: e.visible ?? true };
    case 'BACKGROUND_BLUR':
      return { type: 'BACKGROUND_BLUR', blurType: 'NORMAL', radius: e.radius ?? 10, visible: e.visible ?? true };
    default:
      // Unknown/beta effect type — pass through as-is for forward compatibility
      return { ...e };
  }
}
