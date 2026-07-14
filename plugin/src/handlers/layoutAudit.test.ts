import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getNodeById } from '../utils/getNode';
import { auditLayout } from './layoutAudit';

vi.mock('../utils/getNode', () => ({
  getNodeById: vi.fn(),
}));

function createText(overrides: any = {}): any {
  var remove = vi.fn();
  var node: any = {
    id: '1:2',
    type: 'TEXT',
    name: 'Childcare copy',
    x: 10,
    y: 10,
    width: 80,
    height: 20,
    absoluteBoundingBox: { x: 10, y: 10, width: 80, height: 20 },
    textAutoResize: 'NONE',
    characters: 'Long text that needs more vertical room',
    clone: vi.fn(() => ({
      id: 'clone:2',
      type: 'TEXT',
      width: 120,
      height: 64,
      textAutoResize: 'NONE',
      remove: remove,
    })),
    ...overrides,
  };
  node.removeClone = remove;
  return node;
}

function createRoot(children: any[]): any {
  return {
    id: '1:1',
    type: 'FRAME',
    name: 'Page section',
    x: 0,
    y: 0,
    width: 200,
    height: 200,
    absoluteBoundingBox: { x: 0, y: 0, width: 200, height: 200 },
    layoutMode: 'NONE',
    children: children,
  };
}

beforeEach(() => {
  (globalThis as any).figma = { mixed: Symbol('mixed') };
  vi.clearAllMocks();
});

describe('auditLayout', () => {
  it('reports text overflow and sibling overlap without leaving measurement clones', async () => {
    const text = createText();
    const sibling = {
      id: '1:3', type: 'RECTANGLE', name: 'Overlapping box', x: 70, y: 15,
      width: 80, height: 40,
      absoluteBoundingBox: { x: 70, y: 15, width: 80, height: 40 },
    };
    const root = createRoot([text, sibling]);
    vi.mocked(getNodeById).mockResolvedValueOnce(root as any);

    const result = await auditLayout({ nodeId: root.id });

    expect(result.ok).toBe(false);
    expect(result.summary.visited).toBe(3);
    expect(result.issues.map((issue: any) => issue.code)).toEqual(expect.arrayContaining([
      'FIXED_TEXT_OVERFLOW',
      'SIBLING_OVERLAP',
    ]));
    expect(text.removeClone).toHaveBeenCalled();
    expect(root.children).toHaveLength(2);
  });

  it('honors a caller-provided minimum sibling gap', async () => {
    const first = { id: '1:4', type: 'RECTANGLE', name: 'First', x: 10, y: 80, width: 40, height: 20, absoluteBoundingBox: { x: 10, y: 80, width: 40, height: 20 } };
    const second = { id: '1:5', type: 'RECTANGLE', name: 'Second', x: 54, y: 80, width: 40, height: 20, absoluteBoundingBox: { x: 54, y: 80, width: 40, height: 20 } };
    const root = createRoot([first, second]);
    vi.mocked(getNodeById).mockResolvedValueOnce(root as any);

    const result = await auditLayout({ nodeId: root.id, minGap: 16 });

    expect(result.issues.map((issue: any) => issue.code)).toContain('TIGHT_GAP');
    expect(result.issues.find((issue: any) => issue.code === 'TIGHT_GAP').evidence.minimumRecommended).toBe(16);
  });

  it('reports low-confidence measurement failures instead of guessing', async () => {
    const text = createText({
      id: '1:6',
      clone: vi.fn(() => { throw new Error('clone unavailable'); }),
    });
    const root = createRoot([text]);
    vi.mocked(getNodeById).mockResolvedValueOnce(root as any);

    const result = await auditLayout({ nodeId: root.id });
    const issue = result.issues.find((item: any) => item.code === 'TEXT_MEASUREMENT_UNAVAILABLE');

    expect(issue).toMatchObject({ severity: 'warning', confidence: 'low', nodeIds: [text.id] });
    expect(issue.fix).toBeNull();
  });

  it('returns compact actionable findings without geometry evidence', async () => {
    const text = createText();
    const root = createRoot([text]);
    vi.mocked(getNodeById).mockResolvedValueOnce(root as any);

    const result = await auditLayout({ nodeId: root.id, compact: true });
    const issue = result.issues[0];

    expect(result.truncated).toBe(false);
    expect(issue).toMatchObject({ code: 'FIXED_TEXT_OVERFLOW', nodeIds: [text.id] });
    expect(issue.evidence).toBeUndefined();
    expect(issue.fix.strategy).toBeDefined();
  });
});
