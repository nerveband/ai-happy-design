import { serializeNode, serializeNodeSummary } from '../utils/serialize';
import { getSceneNodeById, getPageNodeById } from '../utils/getNode';

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
    case 'find_free_space':
    case 'free_space': return findFreeSpace(params);
    default: throw new Error('Unknown document action: ' + action + '. Available: get_info, get_selection, set_selection, scan_text, scan_by_type, get_styles, find_by_name, find_by_type, focus, zoom_to, find_free_space');
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

  const nodes = await Promise.all(nodeIds.map((id: string) => getSceneNodeById(id)));
  figma.currentPage.selection = nodes;
  return { selected: nodes.map(n => serializeNodeSummary(n)), count: nodes.length };
}

async function scanText(params: any) {
  var pageId = params.pageId;
  var query = params.query;

  var page = pageId ? await getPageNodeById(pageId) : figma.currentPage;
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

  const nodes = await Promise.all(nodeIds.map((id: string) => getSceneNodeById(id)));

  figma.viewport.scrollAndZoomIntoView(nodes);
  return { zoomed: nodes.map(n => serializeNodeSummary(n)) };
}

async function findFreeSpace(params: any) {
  const requestedWidth = params.width || 1080;
  const requestedHeight = params.height || 1080;
  const gap = params.gap || 100;

  const page = figma.currentPage;
  await page.loadAsync();

  const children = page.children;

  // Collect bounding boxes of all top-level nodes
  const frames: Array<{id: string; name: string; x: number; y: number; width: number; height: number; right: number; bottom: number}> = [];
  for (const child of children) {
    frames.push({
      id: child.id,
      name: child.name,
      x: child.x,
      y: child.y,
      width: child.width,
      height: child.height,
      right: child.x + child.width,
      bottom: child.y + child.height,
    });
  }

  let suggestedX = 0;
  let suggestedY = 0;

  if (frames.length === 0) {
    // Empty page — place at origin
    suggestedX = 0;
    suggestedY = 0;
  } else {
    // Strategy: place to the right of the rightmost element, aligned to top of row.
    // First, find all frames and their right edges.
    // We try to place in the same row (similar y) if possible,
    // otherwise start a new row below everything.

    // Sort frames by y then x for row detection
    const sorted = [...frames].sort((a, b) => a.y - b.y || a.x - b.x);

    // Group into rows: frames whose y-ranges overlap significantly
    const rows: Array<{minY: number; maxBottom: number; maxRight: number; frames: typeof frames}> = [];
    for (const f of sorted) {
      let placed = false;
      for (const row of rows) {
        // Check if this frame overlaps vertically with the row
        const overlapTop = Math.max(f.y, row.minY);
        const overlapBottom = Math.min(f.y + f.height, row.maxBottom);
        const overlap = overlapBottom - overlapTop;
        const minHeight = Math.min(f.height, row.maxBottom - row.minY);
        if (minHeight > 0 && overlap / minHeight > 0.3) {
          // Belongs to this row
          row.maxRight = Math.max(row.maxRight, f.right);
          row.maxBottom = Math.max(row.maxBottom, f.bottom);
          row.minY = Math.min(row.minY, f.y);
          row.frames.push(f);
          placed = true;
          break;
        }
      }
      if (!placed) {
        rows.push({
          minY: f.y,
          maxBottom: f.bottom,
          maxRight: f.right,
          frames: [f],
        });
      }
    }

    // Try to fit in the last row (rightmost placement)
    const lastRow = rows[rows.length - 1];
    suggestedX = lastRow.maxRight + gap;
    suggestedY = lastRow.minY;

    // If placing in this row would make the row excessively wide (>5 frames wide),
    // start a new row below everything instead
    if (lastRow.frames.length >= 5) {
      const globalMaxBottom = Math.max(...frames.map(f => f.bottom));
      suggestedX = sorted[0].x; // Align with leftmost frame
      suggestedY = globalMaxBottom + gap;
    }
  }

  // Snap to 8px grid
  suggestedX = Math.round(suggestedX / 8) * 8;
  suggestedY = Math.round(suggestedY / 8) * 8;

  return {
    x: suggestedX,
    y: suggestedY,
    width: requestedWidth,
    height: requestedHeight,
    gap: gap,
    existingFrames: frames.map(f => ({
      id: f.id,
      name: f.name,
      x: f.x,
      y: f.y,
      width: f.width,
      height: f.height,
    })),
    existingCount: frames.length,
  };
}
