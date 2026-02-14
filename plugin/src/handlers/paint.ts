import { getSceneNodeById } from '../utils/getNode';

function parseHexColor(color: any, fallback = { r: 0, g: 0, b: 0, a: 1 }) {
  if (color && typeof color === 'object' && typeof color.r === 'number') {
    return {
      r: color.r,
      g: color.g,
      b: color.b,
      a: typeof color.a === 'number' ? color.a : 1,
    };
  }

  if (typeof color !== 'string') return fallback;
  const raw = color.trim().replace(/^#/, '');
  if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) return fallback;

  const hex = raw.length === 3
    ? raw.split('').map(ch => ch + ch).join('')
    : raw;
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

function parseStops(rawStops: any): ColorStop[] {
  if (Array.isArray(rawStops)) return rawStops as ColorStop[];
  if (typeof rawStops === 'string') {
    try {
      const parsed = JSON.parse(rawStops);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  }
  return [];
}

export async function handlePaint(action: string, params: any): Promise<any> {
  switch (action) {
    case 'set_color':
    case 'set_fill':
    case 'fill':
    case 'color':
    case 'set_solid': return setSolid(params);
    case 'gradient':
    case 'set_gradient_fill':
    case 'set_gradient': return setGradient(params);
    case 'image':
    case 'image_fill':
    case 'set_image_fill':
    case 'set_image': return setImage(params);
    case 'set_image_url':
    case 'image_url':
    case 'image_from_url':
    case 'set_image_fill_from_url': return setImageUrl(params);
    case 'add_fill': return addFill(params);
    case 'remove_fill': return removeFill(params);
    case 'fills':
    case 'list_fills':
    case 'get_fills': return getFills(params);
    case 'stroke':
    case 'set_stroke_color':
    case 'set_stroke': return setStroke(params);
    default: throw new Error('Unknown paint action: ' + action + '. Available: set_solid, set_gradient, set_image_fill, set_image, set_image_fill_from_url, set_image_url, add_fill, remove_fill, get_fills, set_stroke');
  }
}

async function getFillNodeAsync(nodeId: string): Promise<SceneNode & MinimalFillsMixin> {
  const node = await getSceneNodeById(nodeId);
  if (!('fills' in node)) throw new Error(`Invalid node or node does not support fills: ${nodeId}`);
  return node as SceneNode & MinimalFillsMixin;
}

function decodeBase64ToBytes(imageData: string): Uint8Array {
  const base64 = imageData.includes(',')
    ? imageData.split(',').slice(-1)[0]
    : imageData;
  return figma.base64Decode(base64);
}

async function withTimeout<T>(promise: Promise<T>, timeoutMs: number, label: string): Promise<T> {
  let timer: ReturnType<typeof setTimeout> | undefined;
  try {
    return await Promise.race([
      promise,
      new Promise<T>((_, reject) => {
        timer = setTimeout(() => reject(new Error(`${label} timed out after ${timeoutMs}ms`)), timeoutMs);
      }),
    ]);
  } finally {
    if (timer) clearTimeout(timer);
  }
}

async function createImageFromUrlWithFallback(url: string, timeoutMs: number): Promise<Image> {
  try {
    return await withTimeout(figma.createImageAsync(url), timeoutMs, 'figma.createImageAsync');
  } catch {
    const response = await withTimeout(fetch(url), timeoutMs, 'fetch(url)');
    if (!response.ok) {
      throw new Error(`Failed to fetch image URL: ${response.status} ${response.statusText}`);
    }
    const bytes = new Uint8Array(await response.arrayBuffer());
    return figma.createImage(bytes);
  }
}

async function setSolid(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  const c = parseHexColor(params.color, { r: 0, g: 0, b: 0, a: params.opacity ?? 1 });
  node.fills = [{
    type: 'SOLID',
    color: { r: c.r, g: c.g, b: c.b },
    opacity: typeof params.opacity === 'number' ? params.opacity : c.a,
  }];
  return { id: node.id, name: node.name };
}

async function setGradient(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  const gradientType = params.type ?? params.gradientType ?? 'LINEAR';
  const stops = parseStops(params.stops);

  const typeMap: Record<string, GradientPaint['type']> = {
    LINEAR: 'GRADIENT_LINEAR',
    RADIAL: 'GRADIENT_RADIAL',
    ANGULAR: 'GRADIENT_ANGULAR',
    DIAMOND: 'GRADIENT_DIAMOND',
  };

  const paint: GradientPaint = {
    type: typeMap[gradientType] || 'GRADIENT_LINEAR',
    gradientStops: stops.map((s: any) => {
      const c = parseHexColor(s?.color, { r: 0, g: 0, b: 0, a: 1 });
      return {
        position: s?.position ?? 0,
        color: { r: c.r, g: c.g, b: c.b, a: c.a },
      };
    }),
    gradientTransform: params.handlePositions
      ? [[1, 0, 0], [0, 1, 0]]
      : [[1, 0, 0], [0, 1, 0]],
  };

  node.fills = [paint];
  return { id: node.id, name: node.name, type: gradientType };
}

async function setImage(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  const imageData = params.imageData;
  if (!imageData) throw new Error('imageData is required');

  const bytes = decodeBase64ToBytes(imageData);
  const image = figma.createImage(bytes);
  node.fills = [{
    type: 'IMAGE',
    imageHash: image.hash,
    scaleMode: params.scaleMode || 'FILL',
  }];
  return { id: node.id, name: node.name };
}

async function setImageUrl(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  const url = params.url ?? params.imageUrl;
  if (!url) throw new Error('url is required');

  const timeoutMs = params.timeoutMs ?? 8000;
  const image = await createImageFromUrlWithFallback(url, timeoutMs);
  node.fills = [{
    type: 'IMAGE',
    imageHash: image.hash,
    scaleMode: params.scaleMode || 'FILL',
  }];
  return { id: node.id, name: node.name };
}

async function addFill(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  const fillType = params.type ?? 'SOLID';
  let newFill: Paint;

  if (fillType === 'SOLID') {
    const c = parseHexColor(params.color, { r: 0, g: 0, b: 0, a: params.opacity ?? 1 });
    newFill = {
      type: 'SOLID',
      color: { r: c.r, g: c.g, b: c.b },
      opacity: typeof params.opacity === 'number' ? params.opacity : c.a,
    };
  } else if (fillType === 'GRADIENT_LINEAR' || fillType === 'GRADIENT_RADIAL' || fillType === 'GRADIENT_ANGULAR' || fillType === 'GRADIENT_DIAMOND') {
    const stops = parseStops(params.stops);
    newFill = {
      type: fillType as any,
      gradientStops: stops.map((s: any) => {
        const c = parseHexColor(s?.color, { r: 0, g: 0, b: 0, a: 1 });
        return {
          position: s?.position ?? 0,
          color: { r: c.r, g: c.g, b: c.b, a: c.a },
        };
      }),
      gradientTransform: [[1, 0, 0], [0, 1, 0]],
    };
  } else if (fillType === 'IMAGE' && params.imageData) {
    const bytes = figma.base64Decode(params.imageData);
    const img = figma.createImage(bytes);
    newFill = { type: 'IMAGE', imageHash: img.hash, scaleMode: params.scaleMode || 'FILL' };
  } else {
    throw new Error(`Unsupported fill type: ${fillType}`);
  }

  const fills = [...(node.fills as Paint[])];
  if (params.index !== undefined) fills.splice(params.index, 0, newFill);
  else fills.push(newFill);
  node.fills = fills;
  return { id: node.id, name: node.name, fillCount: fills.length };
}

async function removeFill(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  const index = params.index ?? params.fillIndex ?? 0;
  const fills = [...(node.fills as Paint[])];
  if (index < 0 || index >= fills.length) throw new Error(`Fill index ${index} out of range`);
  fills.splice(index, 1);
  node.fills = fills;
  return { id: node.id, name: node.name, fillCount: fills.length };
}

async function getFills(params: any) {
  const node = await getFillNodeAsync(params.nodeId);
  return { id: node.id, name: node.name, fills: JSON.parse(JSON.stringify(node.fills)) };
}

async function setStroke(params: any) {
  const node = await getSceneNodeById(params.nodeId);
  if (!('strokes' in node)) throw new Error(`Invalid node for strokes: ${params.nodeId}`);

  const c = parseHexColor(params.color, { r: 0, g: 0, b: 0, a: params.opacity ?? 1 });
  (node as GeometryMixin).strokes = [{
    type: 'SOLID',
    color: { r: c.r, g: c.g, b: c.b },
    opacity: c.a,
  }];
  if ('strokeWeight' in node) (node as any).strokeWeight = params.weight ?? params.strokeWeight ?? 1;
  if ('strokeAlign' in node) {
    const align = params.alignment ?? params.strokeAlign;
    if (align) (node as any).strokeAlign = align;
  }
  return { id: node.id, name: node.name };
}
