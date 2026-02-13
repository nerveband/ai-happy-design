export async function handlePaint(action: string, params: any): Promise<any> {
  switch (action) {
    case 'set_solid': return setSolid(params);
    case 'set_gradient': return setGradient(params);
    case 'set_image': return setImage(params);
    case 'set_image_url': return setImageUrl(params);
    case 'add_fill': return addFill(params);
    case 'remove_fill': return removeFill(params);
    case 'get_fills': return getFills(params);
    case 'set_stroke': return setStroke(params);
    default: throw new Error(`Unknown paint action: ${action}`);
  }
}

function getNode(nodeId: string): SceneNode & MinimalFillsMixin {
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node || !('fills' in node)) throw new Error(`Invalid node or node does not support fills: ${nodeId}`);
  return node as SceneNode & MinimalFillsMixin;
}

async function setSolid(params: any) {
  const { nodeId, color } = params;
  const node = getNode(nodeId);
  node.fills = [{
    type: 'SOLID',
    color: { r: color.r, g: color.g, b: color.b },
    opacity: color.a ?? 1,
  }];
  return { id: node.id, name: node.name };
}

async function setGradient(params: any) {
  const { nodeId, type, stops } = params;
  const node = getNode(nodeId);

  const typeMap: Record<string, string> = {
    LINEAR: 'GRADIENT_LINEAR',
    RADIAL: 'GRADIENT_RADIAL',
    ANGULAR: 'GRADIENT_ANGULAR',
    DIAMOND: 'GRADIENT_DIAMOND',
  };

  const paint: GradientPaint = {
    type: (typeMap[type] || 'GRADIENT_LINEAR') as any,
    gradientStops: stops.map((s: any) => ({
      position: s.position,
      color: { r: s.color.r, g: s.color.g, b: s.color.b, a: s.color.a ?? 1 },
    })),
    gradientTransform: [[1, 0, 0], [0, 1, 0]],
  };

  node.fills = [paint];
  return { id: node.id, name: node.name, type };
}

async function setImage(params: any) {
  const { nodeId, imageData, scaleMode } = params;
  const node = getNode(nodeId);

  const bytes = figma.base64Decode(imageData);
  const image = figma.createImage(bytes);
  node.fills = [{
    type: 'IMAGE',
    imageHash: image.hash,
    scaleMode: scaleMode || 'FILL',
  }];
  return { id: node.id, name: node.name };
}

async function setImageUrl(params: any) {
  const { nodeId, url, scaleMode } = params;
  const node = getNode(nodeId);

  const image = await figma.createImageAsync(url);
  node.fills = [{
    type: 'IMAGE',
    imageHash: image.hash,
    scaleMode: scaleMode || 'FILL',
  }];
  return { id: node.id, name: node.name };
}

async function addFill(params: any) {
  const { nodeId, type, color, stops, imageData, index } = params;
  const node = getNode(nodeId);

  let newFill: Paint;
  if (type === 'SOLID') {
    newFill = {
      type: 'SOLID',
      color: { r: color.r, g: color.g, b: color.b },
      opacity: color.a ?? 1,
    };
  } else if (type === 'GRADIENT_LINEAR' || type === 'GRADIENT_RADIAL' || type === 'GRADIENT_ANGULAR' || type === 'GRADIENT_DIAMOND') {
    newFill = {
      type: type as any,
      gradientStops: stops.map((s: any) => ({
        position: s.position,
        color: { r: s.color.r, g: s.color.g, b: s.color.b, a: s.color.a ?? 1 },
      })),
      gradientTransform: [[1, 0, 0], [0, 1, 0]],
    };
  } else if (type === 'IMAGE' && imageData) {
    const bytes = figma.base64Decode(imageData);
    const img = figma.createImage(bytes);
    newFill = { type: 'IMAGE', imageHash: img.hash, scaleMode: 'FILL' };
  } else {
    throw new Error(`Unsupported fill type: ${type}`);
  }

  const fills = [...(node.fills as Paint[])];
  if (index !== undefined) fills.splice(index, 0, newFill);
  else fills.push(newFill);
  node.fills = fills;
  return { id: node.id, name: node.name, fillCount: fills.length };
}

async function removeFill(params: any) {
  const { nodeId, index } = params;
  const node = getNode(nodeId);
  const fills = [...(node.fills as Paint[])];
  if (index < 0 || index >= fills.length) throw new Error(`Fill index ${index} out of range`);
  fills.splice(index, 1);
  node.fills = fills;
  return { id: node.id, name: node.name, fillCount: fills.length };
}

async function getFills(params: any) {
  const { nodeId } = params;
  const node = getNode(nodeId);
  return { id: node.id, name: node.name, fills: JSON.parse(JSON.stringify(node.fills)) };
}

async function setStroke(params: any) {
  const { nodeId, color, weight, alignment } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node || !('strokes' in node)) throw new Error(`Invalid node for strokes: ${nodeId}`);

  (node as any).strokes = [{
    type: 'SOLID',
    color: { r: color.r, g: color.g, b: color.b },
    opacity: color.a ?? 1,
  }];
  if ('strokeWeight' in node) (node as any).strokeWeight = weight ?? 1;
  if ('strokeAlign' in node && alignment) (node as any).strokeAlign = alignment;
  return { id: node.id, name: node.name };
}
