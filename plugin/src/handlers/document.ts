import { serializeNode, serializeNodeSummary } from '../utils/serialize';

export async function handleDocument(action: string, params: any): Promise<any> {
  switch (action) {
    case 'get_info': return getDocInfo(params);
    case 'get_selection': return getSelection(params);
    case 'set_selection': return setSelection(params);
    case 'scan_text': return scanText(params);
    case 'get_styles': return getStyles(params);
    case 'find_by_name': return findByName(params);
    case 'find_by_type': return findByType(params);
    case 'zoom_to': return zoomTo(params);
    default: throw new Error(`Unknown document action: ${action}`);
  }
}

async function getDocInfo(_params: any) {
  const doc = figma.root;
  return {
    name: doc.name,
    id: doc.id,
    pageCount: doc.children.length,
    pages: doc.children.map(p => ({ id: p.id, name: p.name, childCount: p.children.length })),
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
  const { nodeIds } = params;
  const nodes = nodeIds.map((id: string) => {
    const node = figma.getNodeById(id) as SceneNode;
    if (!node) throw new Error(`Node not found: ${id}`);
    return node;
  });
  figma.currentPage.selection = nodes;
  return { selected: nodes.map(n => serializeNodeSummary(n)), count: nodes.length };
}

async function scanText(params: any) {
  const { pageId, query } = params;

  const page = pageId ? figma.getNodeById(pageId) as PageNode : figma.currentPage;
  if (!page || page.type !== 'PAGE') throw new Error('Invalid page');

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
  const paintStyles = figma.getLocalPaintStyles().map(s => ({
    id: s.id, name: s.name, type: 'PAINT',
    description: s.description,
    paints: JSON.parse(JSON.stringify(s.paints)),
  }));

  const textStyles = figma.getLocalTextStyles().map(s => ({
    id: s.id, name: s.name, type: 'TEXT',
    description: s.description,
    fontSize: s.fontSize,
    fontName: s.fontName,
    lineHeight: s.lineHeight,
    letterSpacing: s.letterSpacing,
  }));

  const effectStyles = figma.getLocalEffectStyles().map(s => ({
    id: s.id, name: s.name, type: 'EFFECT',
    description: s.description,
    effects: JSON.parse(JSON.stringify(s.effects)),
  }));

  const gridStyles = figma.getLocalGridStyles().map(s => ({
    id: s.id, name: s.name, type: 'GRID',
    description: s.description,
    grids: JSON.parse(JSON.stringify(s.layoutGrids)),
  }));

  return {
    paintStyles, textStyles, effectStyles, gridStyles,
    total: paintStyles.length + textStyles.length + effectStyles.length + gridStyles.length,
  };
}

async function findByName(params: any) {
  const { name, type, exact = false } = params;
  const page = figma.currentPage;

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
  const { type } = params;
  const page = figma.currentPage;

  const nodes = page.findAllWithCriteria({ types: [type] }) as SceneNode[];

  return {
    nodes: nodes.slice(0, 100).map(n => serializeNodeSummary(n)),
    count: nodes.length,
    truncated: nodes.length > 100,
  };
}

async function zoomTo(params: any) {
  const { nodeIds } = params;
  const nodes = nodeIds.map((id: string) => {
    const node = figma.getNodeById(id) as SceneNode;
    if (!node) throw new Error(`Node not found: ${id}`);
    return node;
  });

  figma.viewport.scrollAndZoomIntoView(nodes);
  return { zoomed: nodes.map(n => serializeNodeSummary(n)) };
}
