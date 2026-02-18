import { describe, it, expect } from 'vitest';
import { resolveRange } from './text';

describe('resolveRange', () => {
  const text = 'Hello, World! Welcome to Figma.';

  // --- match-based resolution ---

  it('resolves a match to correct start/end', () => {
    const result = resolveRange({ match: 'World' }, text);
    expect(result).toEqual({ start: 7, end: 12 });
  });

  it('returns null when match is not found', () => {
    const result = resolveRange({ match: 'Universe' }, text);
    expect(result).toBeNull();
  });

  it('resolves first occurrence when match appears multiple times', () => {
    const repeated = 'abc abc abc';
    const result = resolveRange({ match: 'abc' }, repeated);
    expect(result).toEqual({ start: 0, end: 3 });
  });

  it('resolves match at end of string', () => {
    const result = resolveRange({ match: 'Figma.' }, text);
    expect(result).toEqual({ start: 25, end: 31 });
  });

  it('resolves match at start of string', () => {
    const result = resolveRange({ match: 'Hello' }, text);
    expect(result).toEqual({ start: 0, end: 5 });
  });

  it('resolves single-character match', () => {
    const result = resolveRange({ match: ',' }, text);
    expect(result).toEqual({ start: 5, end: 6 });
  });

  it('returns null for empty match string', () => {
    // indexOf('') returns 0, but start===end so null
    const result = resolveRange({ match: '' }, text);
    expect(result).toBeNull();
  });

  // --- start/end based resolution ---

  it('resolves explicit start/end', () => {
    const result = resolveRange({ start: 0, end: 5 }, text);
    expect(result).toEqual({ start: 0, end: 5 });
  });

  it('clamps end beyond text length', () => {
    const result = resolveRange({ start: 25, end: 100 }, text);
    expect(result).toEqual({ start: 25, end: 31 });
  });

  it('clamps negative start to 0', () => {
    const result = resolveRange({ start: -5, end: 5 }, text);
    expect(result).toEqual({ start: 0, end: 5 });
  });

  it('returns null when start equals end after clamping', () => {
    const result = resolveRange({ start: 50, end: 60 }, text);
    expect(result).toBeNull();
  });

  it('returns null when start >= end', () => {
    const result = resolveRange({ start: 10, end: 5 }, text);
    expect(result).toBeNull();
  });

  // --- edge cases ---

  it('returns null when neither match nor start/end provided', () => {
    const result = resolveRange({}, text);
    expect(result).toBeNull();
  });

  it('returns null when only start is provided (no end)', () => {
    const result = resolveRange({ start: 0 }, text);
    expect(result).toBeNull();
  });

  it('handles empty text', () => {
    const result = resolveRange({ match: 'anything' }, '');
    expect(result).toBeNull();
  });

  it('handles empty text with start/end', () => {
    const result = resolveRange({ start: 0, end: 5 }, '');
    expect(result).toBeNull();
  });

  it('match takes precedence (match field is used when both match and start/end given)', () => {
    const result = resolveRange({ match: 'World', start: 0, end: 1 }, text);
    // match is checked first since typeof match === 'string'
    expect(result).toEqual({ start: 7, end: 12 });
  });
});
