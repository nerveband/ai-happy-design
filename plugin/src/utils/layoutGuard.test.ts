import { describe, it, expect } from 'vitest';
import { parentHasAutoLayout } from './layoutGuard';

describe('parentHasAutoLayout', () => {
  it('returns true when parent has VERTICAL layout', () => {
    const parent = { layoutMode: 'VERTICAL' };
    expect(parentHasAutoLayout(parent)).toBe(true);
  });

  it('returns true when parent has HORIZONTAL layout', () => {
    const parent = { layoutMode: 'HORIZONTAL' };
    expect(parentHasAutoLayout(parent)).toBe(true);
  });

  it('returns false when parent has NONE layout', () => {
    const parent = { layoutMode: 'NONE' };
    expect(parentHasAutoLayout(parent)).toBe(false);
  });

  it('returns false when parent has no layoutMode property', () => {
    const parent = {};
    expect(parentHasAutoLayout(parent)).toBe(false);
  });

  it('returns false when parent is null', () => {
    expect(parentHasAutoLayout(null)).toBe(false);
  });

  it('returns false when parent is undefined', () => {
    expect(parentHasAutoLayout(undefined)).toBe(false);
  });
});
