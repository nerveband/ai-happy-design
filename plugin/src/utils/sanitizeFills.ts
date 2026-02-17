// Known-safe property lists per fill type.
// Figma's internal read-back may include boundVariables, imageRef, and other
// properties that cause errors when re-assigned via node.fills = [...].

const SOLID_KEYS = ['type', 'color', 'opacity', 'visible', 'blendMode'] as const;
const GRADIENT_KEYS = ['type', 'gradientStops', 'gradientTransform', 'opacity', 'visible', 'blendMode'] as const;
const IMAGE_KEYS = ['type', 'imageHash', 'scaleMode', 'opacity', 'visible', 'blendMode',
  'imageTransform', 'scalingFactor', 'rotation', 'filters'] as const;

function pick(obj: any, keys: readonly string[]): any {
  const out: any = {};
  for (const k of keys) {
    if (obj[k] !== undefined) out[k] = obj[k];
  }
  return out;
}

export function sanitizeFill(f: any): any {
  switch (f.type) {
    case 'SOLID':
      return pick(f, SOLID_KEYS);
    case 'GRADIENT_LINEAR':
    case 'GRADIENT_RADIAL':
    case 'GRADIENT_ANGULAR':
    case 'GRADIENT_DIAMOND':
      return pick(f, GRADIENT_KEYS);
    case 'IMAGE':
      return pick(f, IMAGE_KEYS);
    default:
      return { ...f };
  }
}

export function sanitizeFills(fills: readonly any[]): any[] {
  return fills.map(sanitizeFill);
}
