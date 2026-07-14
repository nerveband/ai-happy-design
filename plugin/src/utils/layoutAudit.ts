export type LayoutSeverity = 'error' | 'warning' | 'info';
export type LayoutConfidence = 'high' | 'medium' | 'low';

export interface Bounds {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface GapMeasurement {
  horizontal: number;
  vertical: number;
  distance: number;
}

export interface FixCommand {
  command: string;
  params: Record<string, any>;
}

export interface FixRecommendation {
  strategy: string;
  reason?: string;
  commands: FixCommand[];
}

export interface LayoutIssue {
  code: string;
  severity: LayoutSeverity;
  nodeIds: string[];
  parentId?: string;
  message: string;
  evidence: Record<string, any>;
  confidence: LayoutConfidence;
  fix: FixRecommendation | null;
}

export interface AuditSummary {
  errors: number;
  warnings: number;
  info: number;
  visited: number;
  issues: number;
}

export function intersectionBounds(a: Bounds, b: Bounds): Bounds | null {
  var x = Math.max(a.x, b.x);
  var y = Math.max(a.y, b.y);
  var right = Math.min(a.x + a.width, b.x + b.width);
  var bottom = Math.min(a.y + a.height, b.y + b.height);
  var width = right - x;
  var height = bottom - y;
  if (width <= 0 || height <= 0) return null;
  return { x: x, y: y, width: width, height: height };
}

export function containsBounds(parent: Bounds, child: Bounds): boolean {
  return child.x >= parent.x &&
    child.y >= parent.y &&
    child.x + child.width <= parent.x + parent.width &&
    child.y + child.height <= parent.y + parent.height;
}

export function gapBetweenBounds(a: Bounds, b: Bounds): GapMeasurement {
  var horizontal = 0;
  var vertical = 0;
  if (a.x + a.width < b.x) horizontal = b.x - (a.x + a.width);
  else if (b.x + b.width < a.x) horizontal = a.x - (b.x + b.width);
  if (a.y + a.height < b.y) vertical = b.y - (a.y + a.height);
  else if (b.y + b.height < a.y) vertical = a.y - (b.y + b.height);
  return {
    horizontal: horizontal,
    vertical: vertical,
    distance: Math.max(horizontal, vertical),
  };
}

export function recommendMove(nodeId: string, x: number, y: number, gap: number): FixRecommendation {
  return {
    strategy: 'move_node',
    reason: 'Move the affected node by the minimum measured amount while preserving the requested gap of ' + gap + 'px.',
    commands: [{ command: 'node.move', params: { nodeId: nodeId, x: x, y: y } }],
  };
}

export function recommendResize(nodeId: string, width: number, height: number): FixRecommendation {
  return {
    strategy: 'resize_node',
    reason: 'Resize the node to the smallest measured dimensions that contain its content.',
    commands: [{
      command: 'node.resize',
      params: { nodeId: nodeId, width: Math.max(1, width), height: Math.max(1, height) },
    }],
  };
}

export function emptyAuditSummary(): AuditSummary {
  return { errors: 0, warnings: 0, info: 0, visited: 0, issues: 0 };
}

export function addIssue(summary: AuditSummary, issue: LayoutIssue): void {
  summary.issues++;
  if (issue.severity === 'error') summary.errors++;
  else if (issue.severity === 'warning') summary.warnings++;
  else summary.info++;
}
