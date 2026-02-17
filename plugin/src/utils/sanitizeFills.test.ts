import { describe, it, expect } from 'vitest';
import { sanitizeFill, sanitizeFills } from './sanitizeFills';

describe('sanitizeFill', () => {
  it('passes through SOLID with only documented properties', () => {
    const input = {
      type: 'SOLID',
      color: { r: 1, g: 0, b: 0 },
      opacity: 0.5,
      visible: true,
      boundVariables: { color: 'var-id' },
      _internal: true,
    } as any;

    const result = sanitizeFill(input);
    expect(result.type).toBe('SOLID');
    expect(result.color).toEqual({ r: 1, g: 0, b: 0 });
    expect(result.opacity).toBe(0.5);
    expect((result as any).boundVariables).toBeUndefined();
    expect((result as any)._internal).toBeUndefined();
  });

  it('passes through GRADIENT_LINEAR with stops and transform', () => {
    const input = {
      type: 'GRADIENT_LINEAR',
      gradientStops: [
        { position: 0, color: { r: 0, g: 0, b: 0, a: 1 } },
        { position: 1, color: { r: 1, g: 1, b: 1, a: 1 } },
      ],
      gradientTransform: [[1, 0, 0], [0, 1, 0]],
      visible: true,
      opacity: 1,
      _junk: 'remove',
    } as any;

    const result = sanitizeFill(input);
    expect(result.type).toBe('GRADIENT_LINEAR');
    expect(result.gradientStops).toHaveLength(2);
    expect((result as any)._junk).toBeUndefined();
  });

  it('handles all gradient types', () => {
    for (const gradType of ['GRADIENT_LINEAR', 'GRADIENT_RADIAL', 'GRADIENT_ANGULAR', 'GRADIENT_DIAMOND']) {
      const input = { type: gradType, gradientStops: [], gradientTransform: [[1,0,0],[0,1,0]], visible: true } as any;
      expect(sanitizeFill(input).type).toBe(gradType);
    }
  });

  it('passes through IMAGE with only documented properties', () => {
    const input = {
      type: 'IMAGE',
      imageHash: 'abc123',
      scaleMode: 'FILL',
      visible: true,
      opacity: 1,
      imageRef: 'internal-ref',
      _bound: {},
    } as any;

    const result = sanitizeFill(input);
    expect(result.type).toBe('IMAGE');
    expect(result.imageHash).toBe('abc123');
    expect(result.scaleMode).toBe('FILL');
    expect((result as any).imageRef).toBeUndefined();
    expect((result as any)._bound).toBeUndefined();
  });

  it('passes through unknown fill types as-is', () => {
    const input = { type: 'VIDEO', videoHash: 'xyz', visible: true } as any;
    const result = sanitizeFill(input);
    expect(result.type).toBe('VIDEO');
    expect((result as any).videoHash).toBe('xyz');
  });
});

describe('sanitizeFills', () => {
  it('sanitizes an array of mixed fills', () => {
    const input = [
      { type: 'SOLID', color: { r: 1, g: 0, b: 0 }, opacity: 1, _junk: true },
      { type: 'IMAGE', imageHash: 'abc', scaleMode: 'FILL', _junk: true },
    ] as any;

    const result = sanitizeFills(input);
    expect(result).toHaveLength(2);
    expect((result[0] as any)._junk).toBeUndefined();
    expect((result[1] as any)._junk).toBeUndefined();
  });

  it('returns empty array for empty input', () => {
    expect(sanitizeFills([])).toEqual([]);
  });
});
