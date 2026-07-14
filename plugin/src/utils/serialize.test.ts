import { describe, it, expect, beforeEach } from 'vitest';
import { serializeNode, serializeForPostMessage } from './serialize';

function baseNode(overrides: any = {}): any {
  return {
    id: '1:1',
    type: 'FRAME',
    name: 'Test frame',
    visible: true,
    locked: false,
    cornerRadius: 8,
    ...overrides,
  };
}

describe('serializeNode', () => {
  beforeEach(() => {
    (globalThis as any).figma = { mixed: Symbol('mixed') };
  });

  it('sanitizes nested symbols, unsupported values, and cycles globally', () => {
    const mixed = (globalThis as any).figma.mixed;
    const value: any = {
      mixed,
      nested: { value: mixed },
      big: (globalThis as any).BigInt(42),
      missing: undefined,
      callback: () => undefined,
    };
    value.self = value;

    expect(serializeForPostMessage(value)).toEqual({
      mixed: 'mixed',
      nested: { value: 'mixed' },
      big: '42',
      self: '[Circular]',
    });
  });

  it('converts Figma mixed symbols to a postMessage-safe marker', () => {
    const mixed = (globalThis as any).figma.mixed;
    const frame = serializeNode(baseNode({ cornerRadius: mixed }));
    const text = serializeNode(baseNode({
      type: 'TEXT',
      characters: 'Mixed text',
      fontSize: mixed,
      fontName: mixed,
      textAlignHorizontal: 'LEFT',
      textAlignVertical: 'TOP',
    }));

    expect(frame.cornerRadius).toBe('mixed');
    expect(text.fontSize).toBe('mixed');
    expect(text.fontName).toBe('mixed');
  });
});
