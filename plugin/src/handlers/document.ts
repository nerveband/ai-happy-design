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
    case 'find_nodes': return findNodes(params);
    case 'lint':
    case 'check':
    case 'validate': return lintNode(params);
    default: throw new Error('Unknown document action: ' + action + '. Available: get_info, get_selection, set_selection, scan_text, scan_by_type, get_styles, find_by_name, find_by_type, focus, zoom_to, find_free_space, find_nodes, lint');
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
    const sorted = frames.slice().sort((a, b) => a.y - b.y || a.x - b.x);

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
      let globalMaxBottom = 0;
      for (const frame of frames) {
        if (frame.bottom > globalMaxBottom) globalMaxBottom = frame.bottom;
      }
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

async function findNodes(params: any) {
  var query = params.query || '';
  var type = params.type;
  var textContent = params.textContent;
  var pageId = params.pageId;

  var page = pageId ? await getPageNodeById(pageId) : figma.currentPage;
  await page.loadAsync();

  var nodes: SceneNode[];
  if (type) {
    nodes = page.findAllWithCriteria({ types: [type] }) as SceneNode[];
  } else {
    nodes = page.findAll() as SceneNode[];
  }

  var results = nodes;

  // Filter by name query
  if (query) {
    var lowerQuery = query.toLowerCase();
    results = results.filter(function(n) {
      return n.name.toLowerCase().includes(lowerQuery);
    });
  }

  // Filter by text content (only for TEXT nodes)
  if (textContent) {
    var lowerText = textContent.toLowerCase();
    results = results.filter(function(n) {
      if (n.type === 'TEXT') {
        return (n as TextNode).characters.toLowerCase().includes(lowerText);
      }
      return false;
    });
  }

  return {
    nodes: results.slice(0, 100).map(function(n) { return serializeNodeSummary(n); }),
    count: results.length,
    truncated: results.length > 100,
  };
}

interface LintWarning {
  severity: 'error' | 'warning' | 'info';
  type: string;
  nodeId: string;
  nodeName: string;
  message: string;
}

const DEFAULT_NAME_RE = /^(Frame|Rectangle|Ellipse|Text|Group|Line|Polygon|Star|Vector|Component|Instance|Section)\s+\d+$/;

async function lintNode(params: any) {
  var nodeId = params.nodeId;
  if (!nodeId) throw new Error('nodeId is required');

  var root = await getSceneNodeById(nodeId);
  var warnings: LintWarning[] = [];

  await walkNode(root, null, warnings, 0);

  return { warnings: warnings, count: warnings.length };
}

async function walkNode(
  node: SceneNode,
  parent: (SceneNode & ChildrenMixin) | null,
  warnings: LintWarning[],
  depth: number
) {
  if (depth > 10) return;

  // Check default name
  if (DEFAULT_NAME_RE.test(node.name)) {
    warnings.push({
      severity: 'info',
      type: 'default_name',
      nodeId: node.id,
      nodeName: node.name,
      message: 'Node has a default Figma name. Consider renaming for clarity.',
    });
  }

  // Text-specific checks
  if (node.type === 'TEXT') {
    var textNode = node as TextNode;
    var fs = textNode.fontSize;
    if (fs !== figma.mixed && typeof fs === 'number') {
      if (fs < 12) {
        warnings.push({
          severity: 'info',
          type: 'text_too_small',
          nodeId: node.id,
          nodeName: node.name,
          message: 'Font size ' + fs + 'px is below 12px and may be unreadable.',
        });
      }
      // Check if text is too large relative to nearest parent frame
      if (parent && ('layoutMode' in parent)) {
        var parentHeight = parent.height;
        if (parentHeight > 0 && fs > parentHeight * 0.5) {
          warnings.push({
            severity: 'warning',
            type: 'text_too_large',
            nodeId: node.id,
            nodeName: node.name,
            message: 'Font size ' + fs + 'px exceeds 50% of parent frame height (' + parentHeight + 'px).',
          });
        }
      }
    }
  }

  // Frame-specific checks (overflow, overlap, oversized children)
  if ('children' in node) {
    var frameNode = node as SceneNode & ChildrenMixin;
    var isAutoLayout = ('layoutMode' in frameNode) &&
      (frameNode as any).layoutMode !== 'NONE' &&
      (frameNode as any).layoutMode !== undefined;

    var children = frameNode.children;

    if ('width' in frameNode) {
      var pw = (frameNode as any).width as number;
      var ph = (frameNode as any).height as number;

      // Track absolute children explicitly. In non-auto-layout parents, ABSOLUTE is suspicious
      // and should still be linted as normal content.
      var checkableChildren: SceneNode[] = [];
      var autoLayoutAbsoluteChildren: SceneNode[] = [];
      for (var i = 0; i < children.length; i++) {
        var child = children[i];
        var isAbsolute = (child as any).layoutPositioning === 'ABSOLUTE';

        if (isAbsolute && !isAutoLayout) {
          warnings.push({
            severity: 'info',
            type: 'absolute_child_non_autolayout',
            nodeId: child.id,
            nodeName: child.name,
            message: 'Child uses layoutPositioning:ABSOLUTE under non-auto-layout parent "' + frameNode.name + '". This is likely unintended and can hide overlap issues.',
          });
          checkableChildren.push(child);
          continue;
        }

        if (isAbsolute && isAutoLayout) {
          autoLayoutAbsoluteChildren.push(child);
          continue;
        }

        checkableChildren.push(child);
      }

      // Overflow check for normal/manual-positioned children.
      for (var i = 0; i < checkableChildren.length; i++) {
        var c = checkableChildren[i];
        var cx = c.x;
        var cy = c.y;
        var cw = c.width;
        var ch = c.height;

        if (cx < 0 || cy < 0 || cx + cw > pw || cy + ch > ph) {
          warnings.push({
            severity: 'warning',
            type: 'overflow',
            nodeId: c.id,
            nodeName: c.name,
            message: 'Child extends beyond parent "' + frameNode.name + '" bounds (' + pw + 'x' + ph + ').',
          });
        }
      }

      // Oversized child check
      for (var i = 0; i < checkableChildren.length; i++) {
        var c = checkableChildren[i];
        if (c.width > pw * 1.1) {
          warnings.push({
            severity: 'warning',
            type: 'oversized_child',
            nodeId: c.id,
            nodeName: c.name,
            message: 'Child width (' + c.width + 'px) exceeds parent width (' + pw + 'px) by more than 10%.',
          });
        }
        if (c.height > ph * 1.1) {
          warnings.push({
            severity: 'warning',
            type: 'oversized_child',
            nodeId: c.id,
            nodeName: c.name,
            message: 'Child height (' + c.height + 'px) exceeds parent height (' + ph + 'px) by more than 10%.',
          });
        }
      }

      // Overlap check — compare all sibling pairs.
      for (var i = 0; i < checkableChildren.length; i++) {
        for (var j = i + 1; j < checkableChildren.length; j++) {
          var a = checkableChildren[i];
          var b = checkableChildren[j];
          // Bounding box intersection
          var ax1 = a.x, ay1 = a.y, ax2 = a.x + a.width, ay2 = a.y + a.height;
          var bx1 = b.x, by1 = b.y, bx2 = b.x + b.width, by2 = b.y + b.height;
          if (ax1 < bx2 && ax2 > bx1 && ay1 < by2 && ay2 > by1) {
            warnings.push({
              severity: 'warning',
              type: 'overlap',
              nodeId: a.id,
              nodeName: a.name,
              message: 'Overlaps with sibling "' + b.name + '" (' + b.id + ') in parent "' + frameNode.name + '".',
            });
          }
        }
      }

      // Auto-layout parents: absolute children are intentional, but still ensure they stay in-bounds.
      if (isAutoLayout) {
        for (var i = 0; i < autoLayoutAbsoluteChildren.length; i++) {
          var ac = autoLayoutAbsoluteChildren[i];
          var acx = ac.x;
          var acy = ac.y;
          var acw = ac.width;
          var ach = ac.height;
          if (acx < 0 || acy < 0 || acx + acw > pw || acy + ach > ph) {
            warnings.push({
              severity: 'warning',
              type: 'absolute_overflow',
              nodeId: ac.id,
              nodeName: ac.name,
              message: 'Absolute child extends beyond auto-layout parent "' + frameNode.name + '" bounds (' + pw + 'x' + ph + ').',
            });
          }
        }
      }
    }

    // Recurse into children
    for (var i = 0; i < children.length; i++) {
      await walkNode(children[i] as SceneNode, frameNode, warnings, depth + 1);
    }
  }
}
