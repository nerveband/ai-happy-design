export async function handleShape(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create_rectangle': return createRectangle(params);
    case 'create_ellipse': return createEllipse(params);
    case 'create_polygon': return createPolygon(params);
    case 'create_star': return createStar(params);
    case 'create_line': return createLine(params);
    case 'create_from_svg': return createFromSvg(params);
    default: throw new Error(`Unknown shape action: ${action}`);
  }
}

function getParentNode(parentId?: string): FrameNode | PageNode | GroupNode {
  if (parentId) {
    const parent = figma.getNodeById(parentId);
    if (!parent || !('appendChild' in parent)) throw new Error(`Invalid parent node: ${parentId}`);
    return parent as FrameNode | PageNode | GroupNode;
  }
  return figma.currentPage;
}

async function createRectangle(params: any) {
  const { x = 0, y = 0, width = 100, height = 100, cornerRadius, name, parentId, color } = params;
  const rect = figma.createRectangle();
  rect.x = x;
  rect.y = y;
  rect.resize(width, height);
  if (name) rect.name = name;
  if (cornerRadius !== undefined) {
    if (typeof cornerRadius === 'object') {
      rect.topLeftRadius = cornerRadius.topLeft ?? 0;
      rect.topRightRadius = cornerRadius.topRight ?? 0;
      rect.bottomRightRadius = cornerRadius.bottomRight ?? 0;
      rect.bottomLeftRadius = cornerRadius.bottomLeft ?? 0;
    } else {
      rect.cornerRadius = cornerRadius;
    }
  }
  if (color) {
    rect.fills = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }

  const parent = getParentNode(parentId);
  parent.appendChild(rect);

  return { id: rect.id, name: rect.name, type: rect.type };
}

async function createEllipse(params: any) {
  const { x = 0, y = 0, width = 100, height = 100, name, parentId, color, arcStartAngle, arcEndAngle, arcSweepAngle } = params;
  const ellipse = figma.createEllipse();
  ellipse.x = x;
  ellipse.y = y;
  ellipse.resize(width, height);
  if (name) ellipse.name = name;
  if (color) {
    ellipse.fills = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }
  if (arcStartAngle !== undefined || arcEndAngle !== undefined || arcSweepAngle !== undefined) {
    ellipse.arcData = {
      startingAngle: arcStartAngle ?? 0,
      endingAngle: arcEndAngle ?? (2 * Math.PI),
      innerRadius: 0,
    };
  }

  const parent = getParentNode(parentId);
  parent.appendChild(ellipse);

  return { id: ellipse.id, name: ellipse.name, type: ellipse.type };
}

async function createPolygon(params: any) {
  const { x = 0, y = 0, width = 100, height = 100, pointCount = 3, name, parentId, color } = params;
  const polygon = figma.createPolygon();
  polygon.x = x;
  polygon.y = y;
  polygon.resize(width, height);
  polygon.pointCount = pointCount;
  if (name) polygon.name = name;
  if (color) {
    polygon.fills = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }

  const parent = getParentNode(parentId);
  parent.appendChild(polygon);

  return { id: polygon.id, name: polygon.name, type: polygon.type };
}

async function createStar(params: any) {
  const { x = 0, y = 0, width = 100, height = 100, pointCount = 5, innerRadius = 0.4, name, parentId, color } = params;
  const star = figma.createStar();
  star.x = x;
  star.y = y;
  star.resize(width, height);
  star.pointCount = pointCount;
  star.innerRadius = innerRadius;
  if (name) star.name = name;
  if (color) {
    star.fills = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }

  const parent = getParentNode(parentId);
  parent.appendChild(star);

  return { id: star.id, name: star.name, type: star.type };
}

async function createLine(params: any) {
  const { x = 0, y = 0, length = 100, rotation = 0, name, parentId, color, strokeWeight = 1 } = params;
  const line = figma.createLine();
  line.x = x;
  line.y = y;
  line.resize(length, 0);
  line.rotation = rotation;
  if (name) line.name = name;
  if (color) {
    line.strokes = [{ type: 'SOLID', color: { r: color.r, g: color.g, b: color.b }, opacity: color.a ?? 1 }];
  }
  line.strokeWeight = strokeWeight;

  const parent = getParentNode(parentId);
  parent.appendChild(line);

  return { id: line.id, name: line.name, type: line.type };
}

async function createFromSvg(params: any) {
  const { svg, x = 0, y = 0, name, parentId } = params;
  const node = figma.createNodeFromSvg(svg);
  node.x = x;
  node.y = y;
  if (name) node.name = name;

  const parent = getParentNode(parentId);
  parent.appendChild(node);

  return { id: node.id, name: node.name, type: node.type, width: node.width, height: node.height };
}
