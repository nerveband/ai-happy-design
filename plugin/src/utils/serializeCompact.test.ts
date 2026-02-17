import { describe, it, expect } from 'vitest';
import { serializeNodeCompact } from './serializeCompact';

function mockNode(overrides: any = {}): any {
  return {
    id: '1:1',
    type: 'FRAME',
    name: 'TestFrame',
    x: 0,
    y: 0,
    width: 100,
    height: 200,
    children: [],
    ...overrides,
  };
}

describe('serializeNodeCompact', () => {
  it('serializes a leaf node', () => {
    const node = mockNode({ type: 'RECTANGLE', id: '1:2', name: 'Rect' });
    delete node.children;
    const result = serializeNodeCompact(node);
    expect(result).toHaveLength(1);
    expect(result[0]).toEqual({
      id: '1:2', type: 'RECTANGLE', name: 'Rect',
      x: 0, y: 0, w: 100, h: 200,
      childCount: 0, parentId: null, depth: 0,
    });
  });

  it('serializes parent with children', () => {
    const child1 = mockNode({ id: '1:3', name: 'Child1', type: 'TEXT', x: 10, y: 20, width: 80, height: 30 });
    delete child1.children;
    const child2 = mockNode({ id: '1:4', name: 'Child2', type: 'RECTANGLE', x: 10, y: 60, width: 80, height: 40 });
    delete child2.children;
    const parent = mockNode({ id: '1:2', name: 'Parent', children: [child1, child2] });

    const result = serializeNodeCompact(parent);
    expect(result).toHaveLength(3);
    expect(result[0].id).toBe('1:2');
    expect(result[0].childCount).toBe(2);
    expect(result[0].parentId).toBeNull();
    expect(result[1].parentId).toBe('1:2');
    expect(result[2].parentId).toBe('1:2');
  });

  it('respects maxDepth', () => {
    const grandchild = mockNode({ id: '1:5', name: 'GC', type: 'TEXT' });
    delete grandchild.children;
    const child = mockNode({ id: '1:3', name: 'Child', children: [grandchild] });
    const parent = mockNode({ id: '1:2', name: 'Parent', children: [child] });

    const result = serializeNodeCompact(parent, 1);
    expect(result).toHaveLength(2);
    expect(result[1].childCount).toBe(1);
  });

  it('returns flat array with correct depth values', () => {
    const grandchild = mockNode({ id: '1:5', name: 'GC', type: 'TEXT' });
    delete grandchild.children;
    const child = mockNode({ id: '1:3', name: 'Child', children: [grandchild] });
    const parent = mockNode({ id: '1:2', name: 'Parent', children: [child] });

    const result = serializeNodeCompact(parent, 3);
    expect(result.map(n => n.depth)).toEqual([0, 1, 2]);
  });
});
