import { describe, expect, it } from 'vitest';
import {
  containsBounds,
  gapBetweenBounds,
  intersectionBounds,
  recommendMove,
  recommendResize,
} from './layoutAudit';

describe('layout audit geometry', () => {
  it('calculates overlap rectangles and treats touching edges as non-overlap', () => {
    expect(intersectionBounds(
      { x: 0, y: 0, width: 100, height: 50 },
      { x: 80, y: 20, width: 40, height: 40 },
    )).toEqual({ x: 80, y: 20, width: 20, height: 30 });

    expect(intersectionBounds(
      { x: 0, y: 0, width: 100, height: 50 },
      { x: 100, y: 0, width: 20, height: 20 },
    )).toBeNull();
  });

  it('checks containment and returns the smallest gap between separated bounds', () => {
    expect(containsBounds(
      { x: 0, y: 0, width: 100, height: 100 },
      { x: 10, y: 20, width: 40, height: 30 },
    )).toBe(true);
    expect(containsBounds(
      { x: 0, y: 0, width: 100, height: 100 },
      { x: 90, y: 20, width: 40, height: 30 },
    )).toBe(false);
    expect(gapBetweenBounds(
      { x: 0, y: 0, width: 40, height: 20 },
      { x: 48, y: 5, width: 20, height: 20 },
    )).toEqual({ horizontal: 8, vertical: 0, distance: 8 });
  });

  it('creates bounded move and resize recommendations', () => {
    expect(recommendMove('1:2', 8, 0, 4)).toMatchObject({
      strategy: 'move_node',
      commands: [{ command: 'node.move', params: { nodeId: '1:2', x: 8, y: 0 } }],
    });
    expect(recommendResize('1:3', 120, 96)).toMatchObject({
      strategy: 'resize_node',
      commands: [{ command: 'node.resize', params: { nodeId: '1:3', width: 120, height: 96 } }],
    });
  });
});
