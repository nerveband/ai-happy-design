import { serializeNode, serializeNodeSummary } from '../utils/serialize';

export async function handleDocument(action: string, params: any): Promise<any> {
  switch (action) {
    case 'info':
    case 'get_document':
    case 'get_info': return getDocInfo(params);
    case 'selection':
    case 'get_selected':
    case 'selected':
    case 'get_selection': return getSelection(params);
    case 'select':
    case 'set_selected':
    case 'set_selection': return setSelection(params);
    case 'search_text':
    case 'find_text':
    case 'scan_text': return scanText(params);
    case 'search_by_type':
    case 'scan_type':
    case 'scan_by_type': return findByType(params);
    case 'styles':
    case 'get_all_styles':
    case 'list_styles':
    case 'get_styles': return getStyles(params);
    case 'search':
    case 'search_by_name':
    case 'find_by_name': return findByName(params);
    case 'search_type':
    case 'find_type':
    case 'find_by_type': return findByType(params);
    case 'focus': return zoomTo(params);
    case 'zoom':
    case 'zoom_to_fit':
    case 'zoom_to': return zoomTo(params);
    default: throw new Error('Unknown document action: ' + action + '. Available: get_info, get_selection, set_selection, scan_text, scan_by_type, get_styles, find_by_name, find_by_type, focus, zoom_to');
  }
}

async function getDocInfo(_params: any) {
  var doc = figma.root;
  var pages = [];
  for (var i = 0; i < doc.children.length; i++) {
    var p = doc.children[i];
    await p.loadAsync();
    pages.push({ id: p.id, name: p.name, childCount: p.children.length });
  }
  return {
    name: doc.name,
    id: doc.id,
    pageCount: doc.children.length,
    pages: pages,
    currentPage: { id: figma.currentPage.id, name: figma.currentPage.name },
  };
}

async function getSelection(_params: any) {
  const selection = figma.currentPage.selection;
  return {
    nodes: selection.map(node => serializeNode(node, 0, 2)),
    count: selection.length,
  };
}

async function setSelection(params: any) {
  const rawNodeIds = params.nodeIds;
  const nodeIds = Array.isArray(rawNodeIds)
    ? rawNodeIds
    : String(rawNodeIds || '').split(',').map((id: string) => id.trim()).filter(Boolean);

  const nodes = await Promise.all(nodeIds.map(async (id: string) => {
    const node = await figma.getNodeByIdAsync(id) as SceneNode | null;
    if (!node) throw new Error(`Node not found: ${id}`);
    return node;
  }));
  figma.currentPage.selection = nodes;
  return { selected: nodes.map(n => serializeNodeSummary(n)), count: nodes.length };
}

async function scanText(params: any) {
  var pageId = params.pageId;
  var query = params.query;

  var page = pageId ? await figma.getNodeByIdAsync(pageId) as PageNode | null : figma.currentPage;
  if (!page || page.type !== 'PAGE') throw new Error('Invalid page');
  await page.loadAsync();

  const textNodes = page.findAllWithCriteria({ types: ['TEXT'] }) as TextNode[];

  let results = textNodes.map(node => ({
    id: node.id,
    name: node.name,
    characters: node.characters,
    x: node.x,
    y: node.y,
    width: node.width,
    height: node.height,
  }));

  // Filter by query if provided
  if (query) {
    const lowerQuery = query.toLowerCase();
    results = results.filter(r => r.characters.toLowerCase().includes(lowerQuery));
  }

  return { texts: results, count: results.length };
}

async function getStyles(_params: any) {
  var rawPaint = await figma.getLocalPaintStylesAsync();
  var paintStyles = rawPaint.map(function(s) {
    return {
      id: s.id, name: s.name, type: 'PAINT',
      description: s.description,
      paints: JSON.parse(JSON.stringify(s.paints)),
    };
  });

  var rawText = await figma.getLocalTextStylesAsync();
  var textStyles = rawText.map(function(s) {
    return {
      id: s.id, name: s.name, type: 'TEXT',
      description: s.description,
      fontSize: s.fontSize,
      fontName: s.fontName,
      lineHeight: s.lineHeight,
      letterSpacing: s.letterSpacing,
    };
  });

  var rawEffect = await figma.getLocalEffectStylesAsync();
  var effectStyles = rawEffect.map(function(s) {
    return {
      id: s.id, name: s.name, type: 'EFFECT',
      description: s.description,
      effects: JSON.parse(JSON.stringify(s.effects)),
    };
  });

  var rawGrid = await figma.getLocalGridStylesAsync();
  var gridStyles = rawGrid.map(function(s) {
    return {
      id: s.id, name: s.name, type: 'GRID',
      description: s.description,
      grids: JSON.parse(JSON.stringify(s.layoutGrids)),
    };
  });

  return {
    paintStyles: paintStyles, textStyles: textStyles, effectStyles: effectStyles, gridStyles: gridStyles,
    total: paintStyles.length + textStyles.length + effectStyles.length + gridStyles.length,
  };
}

async function findByName(params: any) {
  var name = params.name;
  var type = params.type;
  var exact = params.exact === true;
  var page = figma.currentPage;
  await page.loadAsync();

  let nodes: SceneNode[];
  if (type) {
    nodes = page.findAllWithCriteria({ types: [type] }) as SceneNode[];
  } else {
    nodes = page.findAll() as SceneNode[];
  }

  const results = nodes.filter(n => {
    if (exact) return n.name === name;
    return n.name.toLowerCase().includes(name.toLowerCase());
  });

  return {
    nodes: results.slice(0, 100).map(n => serializeNodeSummary(n)),
    count: results.length,
    truncated: results.length > 100,
  };
}

async function findByType(params: any) {
  var type = params.type || params.nodeType;
  var page = figma.currentPage;
  await page.loadAsync();

  const nodes = page.findAllWithCriteria({ types: [type] }) as SceneNode[];

  return {
    nodes: nodes.slice(0, 100).map(n => serializeNodeSummary(n)),
    count: nodes.length,
    truncated: nodes.length > 100,
  };
}

async function zoomTo(params: any) {
  var rawNodeIds = params.nodeIds || params.nodeId;
  const nodeIds = Array.isArray(rawNodeIds)
    ? rawNodeIds
    : String(rawNodeIds || '').split(',').map((id: string) => id.trim()).filter(Boolean);

  if (nodeIds.length === 0) {
    throw new Error('No node IDs provided');
  }

  const nodes = await Promise.all(nodeIds.map(async (id: string) => {
    const node = await figma.getNodeByIdAsync(id) as SceneNode | null;
    if (!node) throw new Error(`Node not found: ${id}`);
    return node;
  }));

  figma.viewport.scrollAndZoomIntoView(nodes);
  return { zoomed: nodes.map(n => serializeNodeSummary(n)) };
}
