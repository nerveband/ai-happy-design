function requireFigJam(feature: string): any {
  var figmaAny = figma as any;
  if (figmaAny.editorType !== 'figjam') {
    throw new Error(feature + ' is unavailable because current editorType is ' + figmaAny.editorType);
  }
  return figmaAny;
}

export async function handleFigJam(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create_sticky': return createSticky(params);
    case 'create_shape': return createShape(params);
    case 'create_connector': return createConnector(params);
    case 'get_board': return getBoard();
    default: throw new Error('Unknown figjam action: ' + action);
  }
}

async function createSticky(params: any) {
  var figmaAny = requireFigJam('FigJam sticky');
  if (typeof figmaAny.createSticky !== 'function') throw new Error('createSticky is unavailable in this Figma runtime');
  var sticky = figmaAny.createSticky();
  sticky.text = params.text || '';
  sticky.x = params.x || 0;
  sticky.y = params.y || 0;
  if (params.color) sticky.fills = [{ type: 'SOLID', color: params.color }];
  figma.currentPage.appendChild(sticky);
  return { id: sticky.id, type: sticky.type, text: sticky.text };
}

async function createShape(params: any) {
  var figmaAny = requireFigJam('FigJam shape');
  if (typeof figmaAny.createShapeWithText !== 'function') throw new Error('createShapeWithText is unavailable in this Figma runtime');
  var shape = figmaAny.createShapeWithText();
  shape.shapeType = params.shape || 'ROUNDED_RECTANGLE';
  shape.text = params.text || '';
  shape.x = params.x || 0;
  shape.y = params.y || 0;
  figma.currentPage.appendChild(shape);
  return { id: shape.id, type: shape.type, text: shape.text };
}

async function createConnector(params: any) {
  var figmaAny = requireFigJam('FigJam connector');
  if (typeof figmaAny.createConnector !== 'function') throw new Error('createConnector is unavailable in this Figma runtime');
  var connector = figmaAny.createConnector();
  connector.text = params.text || '';
  figma.currentPage.appendChild(connector);
  return { id: connector.id, type: connector.type, text: connector.text };
}

async function getBoard() {
  var figmaAny = requireFigJam('FigJam board');
  return {
    editorType: figmaAny.editorType,
    currentPage: { id: figma.currentPage.id, name: figma.currentPage.name, childCount: figma.currentPage.children.length },
  };
}
