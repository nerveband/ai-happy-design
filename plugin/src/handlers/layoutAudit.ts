import { getNodeById } from '../utils/getNode';
import {
  addIssue,
  AuditSummary,
  Bounds,
  containsBounds,
  emptyAuditSummary,
  gapBetweenBounds,
  intersectionBounds,
  LayoutIssue,
  recommendMove,
  recommendResize,
} from '../utils/layoutAudit';

interface AuditOptions {
  depth: number;
  maxNodes: number;
  minGap: number;
  compact: boolean;
}

interface TextMeasurement {
  naturalWidth?: number;
  naturalHeight?: number;
  confidence: 'high' | 'medium' | 'low';
  error?: string;
}

interface Snapshot {
  node: any;
  bounds: Bounds;
  parent: Snapshot | null;
  depth: number;
}

export async function auditLayout(params: any): Promise<any> {
  var nodeId = params && params.nodeId;
  if (!nodeId) throw new Error('layout.audit requires nodeId');

  var root = await getNodeById(nodeId);
  var options = normalizeOptions(params || {});
  var issues: LayoutIssue[] = [];
  var summary = emptyAuditSummary();
  var truncated = false;

  function add(issue: LayoutIssue): void {
    issues.push(issue);
    addIssue(summary, issue);
  }

  async function walk(node: any, parent: Snapshot | null, depth: number): Promise<void> {
    if (summary.visited >= options.maxNodes || depth > options.depth) {
      truncated = true;
      return;
    }

    var snapshot: Snapshot = {
      node: node,
      bounds: readBounds(node, parent ? parent.bounds : null),
      parent: parent,
      depth: depth,
    };
    summary.visited++;

    if (node.type === 'TEXT') {
      await auditText(snapshot, add);
    }

    var children = getChildren(node);
    if (children.length > 0) {
      auditContainer(snapshot, children, add, options.minGap);
      if (isManualContainer(node, children)) {
        add({
          code: 'MANUAL_LAYOUT_RISK',
          severity: 'info',
          nodeIds: children.map(function(child: any) { return child.id; }),
          parentId: node.id,
          message: 'This container uses manual positioning for multiple content children; text edits may cause collisions or clipping.',
          evidence: { layoutMode: 'NONE', childCount: children.length },
          confidence: 'medium',
          fix: null,
        });
      }
    }

    for (var i = 0; i < children.length; i++) {
      await walk(children[i], snapshot, depth + 1);
    }
  }

  await walk(root, null, 0);

  return {
    ok: summary.errors === 0,
    nodeId: nodeId,
    summary: summary,
    issues: options.compact ? issues.map(compactIssue) : issues,
    visitedCount: summary.visited,
    truncated: truncated,
  };
}

function normalizeOptions(params: any): AuditOptions {
  var depth = typeof params.depth === 'number' && params.depth >= 0 ? Math.floor(params.depth) : 10;
  var maxNodes = typeof params.maxNodes === 'number' && params.maxNodes > 0 ? Math.floor(params.maxNodes) : 1000;
  var minGap = typeof params.minGap === 'number' && params.minGap >= 0 ? params.minGap : 4;
  return {
    depth: Math.min(depth, 50),
    maxNodes: Math.min(maxNodes, 10000),
    minGap: minGap,
    compact: params.compact === true,
  };
}

function getChildren(node: any): any[] {
  if (!node || !Array.isArray(node.children)) return [];
  return node.children;
}

function readBounds(node: any, parentBounds: Bounds | null): Bounds {
  var absolute = node && node.absoluteBoundingBox;
  if (absolute && isFiniteNumber(absolute.x) && isFiniteNumber(absolute.y) &&
      isFiniteNumber(absolute.width) && isFiniteNumber(absolute.height)) {
    return {
      x: absolute.x,
      y: absolute.y,
      width: Math.max(0, absolute.width),
      height: Math.max(0, absolute.height),
    };
  }

  var x = numberOrZero(node && node.x);
  var y = numberOrZero(node && node.y);
  if (parentBounds) {
    x += parentBounds.x;
    y += parentBounds.y;
  }
  return {
    x: x,
    y: y,
    width: Math.max(0, numberOrZero(node && node.width)),
    height: Math.max(0, numberOrZero(node && node.height)),
  };
}

function isFiniteNumber(value: any): boolean {
  return typeof value === 'number' && isFinite(value);
}

function numberOrZero(value: any): number {
  return isFiniteNumber(value) ? value : 0;
}

function auditContainer(parent: Snapshot, children: any[], add: (issue: LayoutIssue) => void, minGap: number): void {
  var childSnapshots: Snapshot[] = children.map(function(child: any) {
    return {
      node: child,
      bounds: readBounds(child, parent.bounds),
      parent: parent,
      depth: parent.depth + 1,
    };
  });

  for (var i = 0; i < childSnapshots.length; i++) {
    var child = childSnapshots[i];
    if (!containsBounds(parent.bounds, child.bounds)) {
      var absolute = isAbsolutePositioned(child.node);
      var move = clampMoveInside(parent.bounds, child.bounds, child.node, 0);
      add({
        code: 'PARENT_OVERFLOW',
        severity: absolute ? 'warning' : 'error',
        nodeIds: [child.node.id],
        parentId: parent.node.id,
        message: 'Child "' + child.node.name + '" extends outside parent "' + parent.node.name + '" bounds.',
        evidence: {
          parent: parent.bounds,
          child: child.bounds,
          absolute: absolute,
        },
        confidence: 'high',
        fix: move,
      });
    }
  }

  for (var a = 0; a < childSnapshots.length; a++) {
    for (var b = a + 1; b < childSnapshots.length; b++) {
      var left = childSnapshots[a];
      var right = childSnapshots[b];
      var overlap = intersectionBounds(left.bounds, right.bounds);
      if (overlap) {
        var intentional = isAbsolutePositioned(left.node) || isAbsolutePositioned(right.node);
        add({
          code: intentional ? 'INTENTIONAL_OVERLAP' : 'SIBLING_OVERLAP',
          severity: intentional ? 'info' : 'error',
          nodeIds: [left.node.id, right.node.id],
          parentId: parent.node.id,
          message: intentional
            ? 'Absolute-positioned children overlap; verify this is intentional.'
            : 'Sibling nodes overlap by ' + overlap.width + 'x' + overlap.height + 'px.',
          evidence: { overlap: overlap, left: left.bounds, right: right.bounds },
          confidence: 'high',
          fix: intentional ? null : recommendMoveForOverlap(right, overlap, parent),
        });
      } else if (isAligned(left.bounds, right.bounds)) {
        var gap = gapBetweenBounds(left.bounds, right.bounds);
        if (gap.distance > 0 && gap.distance < minGap) {
          add({
            code: 'TIGHT_GAP',
            severity: 'warning',
            nodeIds: [left.node.id, right.node.id],
            parentId: parent.node.id,
            message: 'Neighboring nodes have only ' + gap.distance + 'px of separation.',
            evidence: { gap: gap, minimumRecommended: minGap },
            confidence: 'high',
            fix: recommendMoveForGap(right, left, right, gap, minGap),
          });
        }
      }
    }
  }
}

function isAligned(a: Bounds, b: Bounds): boolean {
  var verticalOverlap = Math.min(a.y + a.height, b.y + b.height) - Math.max(a.y, b.y);
  var horizontalOverlap = Math.min(a.x + a.width, b.x + b.width) - Math.max(a.x, b.x);
  return verticalOverlap > 0 || horizontalOverlap > 0;
}

function recommendMoveForOverlap(snapshot: Snapshot, overlap: Bounds, parent: Snapshot): any {
  var node = snapshot.node;
  if (parent.node.layoutMode && parent.node.layoutMode !== 'NONE') return null;
  var moveX = overlap.width <= overlap.height ? overlap.width + 4 : 0;
  var moveY = moveX === 0 ? overlap.height + 4 : 0;
  return recommendMove(node.id, numberOrZero(node.x) + moveX, numberOrZero(node.y) + moveY, 4);
}

function recommendMoveForGap(right: Snapshot, left: Snapshot, rightSnapshot: Snapshot, gap: any, minGap: number): any {
  var node = rightSnapshot.node;
  if (gap.horizontal > 0) {
    return recommendMove(node.id, numberOrZero(node.x) + (minGap - gap.horizontal), numberOrZero(node.y), minGap);
  }
  return recommendMove(node.id, numberOrZero(node.x), numberOrZero(node.y) + (minGap - gap.vertical), minGap);
}

function clampMoveInside(parent: Bounds, child: Bounds, node: any, gap: number): any {
  var dx = 0;
  var dy = 0;
  if (child.x < parent.x + gap) dx = parent.x + gap - child.x;
  else if (child.x + child.width > parent.x + parent.width - gap) dx = parent.x + parent.width - gap - (child.x + child.width);
  if (child.y < parent.y + gap) dy = parent.y + gap - child.y;
  else if (child.y + child.height > parent.y + parent.height - gap) dy = parent.y + parent.height - gap - (child.y + child.height);
  if (dx === 0 && dy === 0) return null;
  return recommendMove(node.id, numberOrZero(node.x) + dx, numberOrZero(node.y) + dy, gap);
}

function isAbsolutePositioned(node: any): boolean {
  return node && node.layoutPositioning === 'ABSOLUTE';
}

async function auditText(snapshot: Snapshot, add: (issue: LayoutIssue) => void): Promise<void> {
  var node = snapshot.node;
  if (node.textAutoResize !== 'NONE' || typeof node.clone !== 'function') return;

  var measurement = await measureText(node);
  if (measurement.error) {
    add({
      code: 'TEXT_MEASUREMENT_UNAVAILABLE',
      severity: 'warning',
      nodeIds: [node.id],
      parentId: snapshot.parent ? snapshot.parent.node.id : undefined,
      message: 'Figma could not measure the natural text bounds; do not guess a repair.',
      evidence: { error: measurement.error, bounds: snapshot.bounds },
      confidence: 'low',
      fix: null,
    });
    return;
  }
  if (measurement.naturalHeight !== undefined && measurement.naturalHeight > numberOrZero(node.height) + 0.5) {
    add({
      code: 'FIXED_TEXT_OVERFLOW',
      severity: 'error',
      nodeIds: [node.id],
      parentId: snapshot.parent ? snapshot.parent.node.id : undefined,
      message: 'Text requires ' + measurement.naturalHeight + 'px of height but its fixed box is only ' + numberOrZero(node.height) + 'px.',
      evidence: {
        bounds: snapshot.bounds,
        actualHeight: numberOrZero(node.height),
        naturalHeight: measurement.naturalHeight,
        delta: measurement.naturalHeight - numberOrZero(node.height),
      },
      confidence: measurement.confidence,
      fix: recommendResize(node.id, numberOrZero(node.width), measurement.naturalHeight),
    });
  }
  if (measurement.naturalWidth !== undefined && measurement.naturalWidth > numberOrZero(node.width) + 0.5) {
    add({
      code: 'TEXT_CLIPPED',
      severity: 'error',
      nodeIds: [node.id],
      parentId: snapshot.parent ? snapshot.parent.node.id : undefined,
      message: 'Text requires ' + measurement.naturalWidth + 'px of width but its fixed box is only ' + numberOrZero(node.width) + 'px.',
      evidence: {
        bounds: snapshot.bounds,
        actualWidth: numberOrZero(node.width),
        naturalWidth: measurement.naturalWidth,
        delta: measurement.naturalWidth - numberOrZero(node.width),
      },
      confidence: measurement.confidence,
      fix: recommendResize(node.id, measurement.naturalWidth, numberOrZero(node.height)),
    });
  }
}

async function measureText(node: any): Promise<TextMeasurement> {
  var result: TextMeasurement = { confidence: 'high' };
  var clones: any[] = [];
  try {
    var heightClone = node.clone();
    clones.push(heightClone);
    if ('visible' in heightClone) heightClone.visible = false;
    if ('textAutoResize' in heightClone) heightClone.textAutoResize = 'HEIGHT';
    result.naturalHeight = numberOrZero(heightClone.height);

    var widthClone = node.clone();
    clones.push(widthClone);
    if ('visible' in widthClone) widthClone.visible = false;
    if ('textAutoResize' in widthClone) widthClone.textAutoResize = 'WIDTH_AND_HEIGHT';
    result.naturalWidth = numberOrZero(widthClone.width);
  } catch (error: any) {
    result.confidence = 'low';
    result.error = error && error.message ? error.message : String(error);
  } finally {
    for (var i = 0; i < clones.length; i++) {
      try {
        if (clones[i] && typeof clones[i].remove === 'function') clones[i].remove();
      } catch (_error) {
        result.confidence = 'low';
      }
    }
  }
  return result;
}

function isManualContainer(node: any, children: any[]): boolean {
  if (!node || node.layoutMode !== 'NONE' || children.length < 2) return false;
  var contentChildren = children.filter(function(child: any) {
    return child.type === 'TEXT' || child.type === 'FRAME' || child.type === 'GROUP' || child.type === 'INSTANCE';
  });
  return contentChildren.length >= 2;
}

function compactIssue(issue: LayoutIssue): any {
  return {
    code: issue.code,
    severity: issue.severity,
    nodeIds: issue.nodeIds,
    parentId: issue.parentId,
    message: issue.message,
    confidence: issue.confidence,
    fix: issue.fix,
  };
}
