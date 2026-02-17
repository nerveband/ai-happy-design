// Known-safe property lists per effect type.
// When reading node.effects and writing back, Figma's internal representation
// may include extra properties (boundVariables, _internal, etc.) that cause
// "Invalid format for effects" when re-assigned. This sanitizer strips them.

const SHADOW_KEYS = ['type', 'color', 'offset', 'radius', 'spread', 'visible', 'blendMode'] as const;
const BLUR_KEYS = ['type', 'blurType', 'radius', 'visible'] as const;
const NOISE_KEYS = ['type', 'noiseType', 'color', 'noiseSize', 'density', 'visible', 'secondaryColor', 'opacity', 'blendMode'] as const;
const TEXTURE_KEYS = ['type', 'noiseSize', 'radius', 'clipToShape', 'visible'] as const;
const GLASS_KEYS = ['type', 'visible', 'lightIntensity', 'lightAngle', 'refraction', 'depth', 'dispersion', 'radius'] as const;

function pick(obj: any, keys: readonly string[]): any {
  const out: any = {};
  for (const k of keys) {
    if (obj[k] !== undefined) out[k] = obj[k];
  }
  return out;
}

export function sanitizeEffect(e: any): any {
  switch (e.type) {
    case 'DROP_SHADOW':
    case 'INNER_SHADOW':
      return pick(e, SHADOW_KEYS);
    case 'LAYER_BLUR':
    case 'BACKGROUND_BLUR':
      return pick(e, BLUR_KEYS);
    case 'NOISE':
      return pick(e, NOISE_KEYS);
    case 'TEXTURE':
      return pick(e, TEXTURE_KEYS);
    case 'GLASS':
      return pick(e, GLASS_KEYS);
    default:
      // Unknown future type — pass through as-is for forward compatibility.
      return { ...e };
  }
}

export function sanitizeEffects(effects: readonly any[]): any[] {
  return effects.map(sanitizeEffect);
}
