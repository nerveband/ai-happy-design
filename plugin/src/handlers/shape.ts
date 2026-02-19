import { resolveStableId } from '../utils/stableId';
import { getParentById } from '../utils/getNode';
import { parentHasAutoLayout } from '../utils/layoutGuard';

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

function applyFill(node: GeometryMixin, color: any) {
  if (!color) return;
  const c = parseHexColor(color);
  node.fills = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
}

function applyLayoutProps(node: SceneNode, params: any) {
  if (params.layoutSizingHorizontal !== undefined) (node as any).layoutSizingHorizontal = params.layoutSizingHorizontal;
  if (params.layoutSizingVertical !== undefined) (node as any).layoutSizingVertical = params.layoutSizingVertical;
  if (params.layoutAlign !== undefined) (node as any).layoutAlign = params.layoutAlign;
  if (params.layoutGrow !== undefined) (node as any).layoutGrow = params.layoutGrow;
  if (params.layoutPositioning !== undefined) {
    // layoutPositioning: "ABSOLUTE" is only valid inside auto-layout frames.
    // Silently skip if parent has no auto-layout to prevent Figma errors.
    if (params.layoutPositioning !== 'ABSOLUTE' || parentHasAutoLayout(node.parent)) {
      (node as any).layoutPositioning = params.layoutPositioning;
    }
  }
}

function applyGeometryStyle(node: GeometryMixin & SceneNode, params: any) {
  const strokeColor = params.stroke ?? params.strokeColor;
  if (strokeColor) {
    const c = parseHexColor(strokeColor);
    node.strokes = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
  }

  const strokeWidth = params.strokeWidth ?? params.strokeWeight;
  if (strokeWidth !== undefined && 'strokeWeight' in node) {
    (node as any).strokeWeight = strokeWidth;
  }

  if (params.noFill === true) {
    node.fills = [];
  }

  if (params.opacity !== undefined && 'opacity' in node) {
    (node as any).opacity = Math.max(0, Math.min(1, params.opacity));
  }
}

export async function handleShape(action: string, params: any): Promise<any> {
  switch (action) {
    case 'rectangle':
    case 'rect':
    case 'add_rectangle':
    case 'create_rect':
    case 'create_rectangle': return createRectangle(params);
    case 'ellipse':
    case 'circle':
    case 'add_ellipse':
    case 'create_circle':
    case 'create_ellipse': return createEllipse(params);
    case 'polygon':
    case 'add_polygon':
    case 'create_polygon': return createPolygon(params);
    case 'star':
    case 'add_star':
    case 'create_star': return createStar(params);
    case 'line':
    case 'add_line':
    case 'create_line': return createLine(params);
    case 'svg':
    case 'from_svg':
    case 'import_svg':
    case 'create_from_svg': return createFromSvg(params);
    case 'image':
    case 'add_image':
    case 'create_image': return createImage(params);
    default: throw new Error('Unknown shape action: ' + action + '. Available: create_rectangle, create_ellipse, create_polygon, create_star, create_line, create_from_svg, create_image');
  }
}

async function getParentNode(parentId?: string): Promise<BaseNode & ChildrenMixin> {
  if (parentId) {
    const parent = await getParentById(parentId);
    if (parent) {
      return parent;
    }
  }
  return figma.currentPage;
}

async function createRectangle(params: any) {
  const rect = figma.createRectangle();
  rect.x = params.x ?? 0;
  rect.y = params.y ?? 0;
  rect.resize(params.width ?? 100, params.height ?? 100);
  if (params.name) rect.name = params.name;

  if (params.cornerRadius !== undefined) {
    if (typeof params.cornerRadius === 'object') {
      rect.topLeftRadius = params.cornerRadius.topLeft ?? 0;
      rect.topRightRadius = params.cornerRadius.topRight ?? 0;
      rect.bottomRightRadius = params.cornerRadius.bottomRight ?? 0;
      rect.bottomLeftRadius = params.cornerRadius.bottomLeft ?? 0;
    } else {
      rect.cornerRadius = params.cornerRadius;
    }
  }
  applyFill(rect, params.color ?? params.fillColor);
  applyGeometryStyle(rect, params);

  const parent = await getParentNode(params.parentId);
  parent.appendChild(rect);
  applyLayoutProps(rect, params);
  var stableId = await resolveStableId(rect, parent);
  return { id: stableId, name: rect.name, type: rect.type, parentId: params.parentId || figma.currentPage.id };
}

async function createEllipse(params: any) {
  const ellipse = figma.createEllipse();
  ellipse.x = params.x ?? 0;
  ellipse.y = params.y ?? 0;
  ellipse.resize(params.width ?? 100, params.height ?? 100);
  if (params.name) ellipse.name = params.name;
  applyFill(ellipse, params.color ?? params.fillColor);
  applyGeometryStyle(ellipse, params);

  const arcStartAngle = params.arcStartAngle;
  const arcEndAngle = params.arcEndAngle;
  if (arcStartAngle !== undefined || arcEndAngle !== undefined || params.arcSweepAngle !== undefined) {
    ellipse.arcData = {
      startingAngle: arcStartAngle ?? 0,
      endingAngle: arcEndAngle ?? (2 * Math.PI),
      innerRadius: 0,
    };
  }

  const parent = await getParentNode(params.parentId);
  parent.appendChild(ellipse);
  applyLayoutProps(ellipse, params);
  var stableId = await resolveStableId(ellipse, parent);
  return { id: stableId, name: ellipse.name, type: ellipse.type };
}

async function createPolygon(params: any) {
  const polygon = figma.createPolygon();
  polygon.x = params.x ?? 0;
  polygon.y = params.y ?? 0;
  polygon.resize(params.width ?? 100, params.height ?? 100);
  polygon.pointCount = params.pointCount ?? params.sides ?? 6;
  if (params.name) polygon.name = params.name;
  applyFill(polygon, params.color ?? params.fillColor);
  applyGeometryStyle(polygon, params);

  const parent = await getParentNode(params.parentId);
  parent.appendChild(polygon);
  applyLayoutProps(polygon, params);
  var stableId = await resolveStableId(polygon, parent);
  return { id: stableId, name: polygon.name, type: polygon.type };
}

async function createStar(params: any) {
  const star = figma.createStar();
  star.x = params.x ?? 0;
  star.y = params.y ?? 0;
  star.resize(params.width ?? 100, params.height ?? 100);
  star.pointCount = params.pointCount ?? params.points ?? 5;
  star.innerRadius = params.innerRadius ?? 0.4;
  if (params.name) star.name = params.name;
  applyFill(star, params.color ?? params.fillColor);
  applyGeometryStyle(star, params);

  const parent = await getParentNode(params.parentId);
  parent.appendChild(star);
  applyLayoutProps(star, params);
  var stableId = await resolveStableId(star, parent);
  return { id: stableId, name: star.name, type: star.type };
}

async function createLine(params: any) {
  const line = figma.createLine();

  if (params.startX !== undefined && params.startY !== undefined && params.endX !== undefined && params.endY !== undefined) {
    const dx = params.endX - params.startX;
    const dy = params.endY - params.startY;
    const length = Math.sqrt((dx * dx) + (dy * dy));
    line.x = params.startX;
    line.y = params.startY;
    line.resize(length, 0);
    line.rotation = (Math.atan2(dy, dx) * 180) / Math.PI;
  } else {
    line.x = params.x ?? 0;
    line.y = params.y ?? 0;
    line.resize(params.length ?? params.width ?? 100, 0);
    line.rotation = params.rotation ?? 0;
  }

  if (params.name) line.name = params.name;
  const strokeColor = params.stroke ?? params.color ?? params.strokeColor;
  if (strokeColor) {
    const c = parseHexColor(strokeColor);
    line.strokes = [{ type: 'SOLID', color: { r: c.r, g: c.g, b: c.b }, opacity: c.a }];
  }
  line.strokeWeight = params.strokeWeight ?? params.strokeWidth ?? 1;
  if (params.opacity !== undefined) {
    line.opacity = Math.max(0, Math.min(1, params.opacity));
  }

  const parent = await getParentNode(params.parentId);
  parent.appendChild(line);
  var stableId = await resolveStableId(line, parent);
  return { id: stableId, name: line.name, type: line.type };
}

async function createFromSvg(params: any) {
  const svg = params.svg ?? params.svgPath ?? params.svgString;
  if (!svg) throw new Error('svg string is required');

  const node = figma.createNodeFromSvg(svg);
  node.x = params.x ?? 0;
  node.y = params.y ?? 0;
  if (params.name) node.name = params.name;

  const parent = await getParentNode(params.parentId);
  parent.appendChild(node);
  var stableId = await resolveStableId(node, parent);
  return { id: stableId, name: node.name, type: node.type, width: node.width, height: node.height };
}

async function createImage(params: any) {
  const imageData = params.imageData;
  if (!imageData) throw new Error('imageData (base64) is required');

  const rect = figma.createRectangle();
  rect.x = params.x ?? 0;
  rect.y = params.y ?? 0;
  rect.resize(params.width ?? 400, params.height ?? 400);
  rect.name = params.name ?? 'Image';

  if (params.cornerRadius !== undefined) {
    rect.cornerRadius = params.cornerRadius;
  }

  // Decode base64 and create image
  const base64 = imageData.includes(',') ? imageData.split(',').slice(-1)[0] : imageData;
  const bytes = figma.base64Decode(base64);
  const image = figma.createImage(bytes);
  const scaleMode = params.scaleMode ?? 'FILL';

  rect.fills = [{
    type: 'IMAGE',
    imageHash: image.hash,
    scaleMode: scaleMode as 'FILL' | 'FIT' | 'CROP' | 'TILE',
  }];

  const parent = await getParentNode(params.parentId);
  parent.appendChild(rect);
  applyGeometryStyle(rect, params);
  applyLayoutProps(rect, params);
  var stableId = await resolveStableId(rect, parent);
  return {
    id: stableId,
    name: rect.name,
    type: rect.type,
    width: rect.width,
    height: rect.height,
    imageHash: image.hash,
    scaleMode: scaleMode,
  };
}
