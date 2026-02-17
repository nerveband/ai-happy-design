import { describe, it, expect } from 'vitest';
import { sanitizeEffect, sanitizeEffects } from './sanitizeEffects';

describe('sanitizeEffect', () => {
  it('passes through DROP_SHADOW with only documented properties', () => {
    const input = {
      type: 'DROP_SHADOW',
      color: { r: 0, g: 0, b: 0, a: 0.25 },
      offset: { x: 0, y: 4 },
      radius: 4,
      spread: 0,
      visible: true,
      blendMode: 'NORMAL',
      _internal: true,
      boundVariables: { color: 'some-var-id' },
    } as any;

    const result = sanitizeEffect(input);
    expect(result.type).toBe('DROP_SHADOW');
    expect(result.color).toEqual({ r: 0, g: 0, b: 0, a: 0.25 });
    expect(result.offset).toEqual({ x: 0, y: 4 });
    expect(result.radius).toBe(4);
    expect(result.spread).toBe(0);
    expect(result.visible).toBe(true);
    expect(result.blendMode).toBe('NORMAL');
    expect((result as any)._internal).toBeUndefined();
    expect((result as any).boundVariables).toBeUndefined();
  });

  it('passes through INNER_SHADOW with only documented properties', () => {
    const input = {
      type: 'INNER_SHADOW',
      color: { r: 1, g: 0, b: 0, a: 0.5 },
      offset: { x: 2, y: 2 },
      radius: 8,
      spread: 1,
      visible: true,
      blendMode: 'MULTIPLY',
      boundVariables: {},
    } as any;

    const result = sanitizeEffect(input);
    expect(result.type).toBe('INNER_SHADOW');
    expect(result.blendMode).toBe('MULTIPLY');
    expect((result as any).boundVariables).toBeUndefined();
  });

  it('passes through LAYER_BLUR with only documented properties', () => {
    const input = {
      type: 'LAYER_BLUR',
      radius: 10,
      visible: true,
      blurType: 'NORMAL',
      _extra: 'junk',
    } as any;

    const result = sanitizeEffect(input);
    expect(result).toEqual({ type: 'LAYER_BLUR', blurType: 'NORMAL', radius: 10, visible: true });
  });

  it('passes through BACKGROUND_BLUR with only documented properties', () => {
    const input = {
      type: 'BACKGROUND_BLUR',
      radius: 20,
      visible: true,
      blurType: 'NORMAL',
      _something: 123,
    } as any;

    const result = sanitizeEffect(input);
    expect(result).toEqual({ type: 'BACKGROUND_BLUR', blurType: 'NORMAL', radius: 20, visible: true });
  });

  it('passes through NOISE with known-safe properties only', () => {
    const input = {
      type: 'NOISE',
      noiseType: 'MONOTONE',
      color: { r: 1, g: 1, b: 1, a: 0.25 },
      noiseSize: 100,
      density: 0.3,
      visible: true,
      _internalId: 'abc',
    } as any;

    const result = sanitizeEffect(input);
    expect(result.type).toBe('NOISE');
    expect((result as any).noiseType).toBe('MONOTONE');
    expect((result as any).noiseSize).toBe(100);
    expect((result as any).density).toBe(0.3);
    expect((result as any)._internalId).toBeUndefined();
  });

  it('passes through NOISE DUOTONE with secondaryColor', () => {
    const input = {
      type: 'NOISE',
      noiseType: 'DUOTONE',
      color: { r: 0, g: 0, b: 0, a: 1 },
      secondaryColor: { r: 1, g: 1, b: 1, a: 1 },
      noiseSize: 50,
      density: 0.5,
      visible: true,
    } as any;

    const result = sanitizeEffect(input);
    expect((result as any).secondaryColor).toEqual({ r: 1, g: 1, b: 1, a: 1 });
  });

  it('passes through TEXTURE with known-safe properties only', () => {
    const input = {
      type: 'TEXTURE',
      noiseSize: 100,
      radius: 5,
      clipToShape: true,
      visible: true,
      _internal: 'drop-me',
    } as any;

    const result = sanitizeEffect(input);
    expect(result.type).toBe('TEXTURE');
    expect((result as any).noiseSize).toBe(100);
    expect((result as any).radius).toBe(5);
    expect((result as any).clipToShape).toBe(true);
    expect((result as any)._internal).toBeUndefined();
  });

  it('passes through GLASS with known-safe properties only', () => {
    const input = {
      type: 'GLASS',
      visible: true,
      lightIntensity: 0.5,
      lightAngle: 45,
      refraction: 0.5,
      depth: 1,
      dispersion: 0,
      radius: 12,
      _somethingInternal: true,
    } as any;

    const result = sanitizeEffect(input);
    expect(result.type).toBe('GLASS');
    expect((result as any).lightIntensity).toBe(0.5);
    expect((result as any)._somethingInternal).toBeUndefined();
  });

  it('passes through unknown types as-is (forward compatibility)', () => {
    const input = {
      type: 'FUTURE_EFFECT',
      someParam: 42,
      visible: true,
    } as any;

    const result = sanitizeEffect(input);
    expect(result.type).toBe('FUTURE_EFFECT');
    expect((result as any).someParam).toBe(42);
  });
});

describe('sanitizeEffects', () => {
  it('sanitizes an array of mixed effects', () => {
    const input = [
      { type: 'DROP_SHADOW', color: { r: 0, g: 0, b: 0, a: 0.25 }, offset: { x: 0, y: 4 }, radius: 4, spread: 0, visible: true, blendMode: 'NORMAL', _junk: true },
      { type: 'GLASS', visible: true, lightIntensity: 0.5, lightAngle: 45, refraction: 0.5, depth: 1, dispersion: 0, radius: 12, _junk: true },
    ] as any;

    const result = sanitizeEffects(input);
    expect(result).toHaveLength(2);
    expect((result[0] as any)._junk).toBeUndefined();
    expect((result[1] as any)._junk).toBeUndefined();
  });

  it('returns empty array for empty input', () => {
    expect(sanitizeEffects([])).toEqual([]);
  });
});
