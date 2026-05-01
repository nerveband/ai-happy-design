import { beforeEach, describe, it, expect, vi } from 'vitest';
import { handleText, resolveRange } from './text';

vi.mock('../utils/fonts', () => ({
  loadFont: vi.fn(async () => undefined),
  loadNodeFonts: vi.fn(async () => undefined),
  resolveFontFamily: vi.fn((family: string) => family),
}));

vi.mock('../utils/getNode', () => ({
  getTextNodeById: vi.fn(),
  getParentById: vi.fn(async () => undefined),
}));

vi.mock('../utils/stableId', () => ({
  resolveStableId: vi.fn(async (node: any) => node.id),
}));

let nextNodeId = 1;

function createMockPage() {
  return {
    id: '0:1',
    type: 'PAGE',
    children: [] as any[],
    appendChild(node: any) {
      this.children.push(node);
      node.parent = this;
    },
  };
}

function createMockTextNode() {
  var characters = '';
  var node: any = {
    id: '1:' + nextNodeId++,
    type: 'TEXT',
    name: 'Text',
    x: 0,
    y: 0,
    width: 0,
    height: 0,
    visible: true,
    fontName: { family: 'Inter', style: 'Regular' },
    fontSize: 16,
    textAutoResize: 'WIDTH_AND_HEIGHT',
    resize(width: number) {
      this.width = width;
      recalculateTextSize(this);
    },
    remove() {
      if (!this.parent) return;
      var index = this.parent.children.indexOf(this);
      if (index >= 0) this.parent.children.splice(index, 1);
      this.parent = undefined;
    },
  };
  Object.defineProperty(node, 'characters', {
    get() {
      return characters;
    },
    set(value: string) {
      characters = value;
      recalculateTextSize(node);
    },
  });
  recalculateTextSize(node);
  return node;
}

function recalculateTextSize(node: any) {
  var fontSize = typeof node.fontSize === 'number' ? node.fontSize : 16;
  var content = String(node.characters || '');
  var naturalWidth = Math.max(1, content.length) * fontSize * 0.5;
  var width = node.width || naturalWidth;
  var lines = Math.max(1, Math.ceil(naturalWidth / Math.max(1, width)));
  node.width = width;
  node.height = lines * fontSize * 1.2;
}

beforeEach(() => {
  nextNodeId = 1;
  var currentPage = createMockPage();
  (globalThis as any).figma = {
    currentPage,
    mixed: Symbol('mixed'),
    createText: vi.fn(createMockTextNode),
  };
});

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

describe('high-level text commands', () => {
  it('measures text with the requested width and font size', async () => {
    const result = await handleText('measure', {
      text: 'Sponsor package',
      width: 120,
      fontFamily: 'Inter',
      fontSize: 20,
    });

    expect(result).toMatchObject({
      width: 120,
      fontSize: 20,
      lineCount: 1,
    });
    expect(result.height).toBeGreaterThan(0);
    expect((globalThis as any).figma.currentPage.children).toHaveLength(0);
  });

  it('fits text by returning the largest size that stays inside the box', async () => {
    const result = await handleText('fit_box', {
      text: 'A longer sponsorship headline',
      width: 160,
      height: 40,
      minFontSize: 10,
      maxFontSize: 30,
    });

    expect(result.fits).toBe(true);
    expect(result.fontSize).toBeGreaterThanOrEqual(10);
    expect(result.fontSize).toBeLessThan(30);
    expect(result.height).toBeLessThanOrEqual(40.5);
    expect((globalThis as any).figma.currentPage.children).toHaveLength(0);
  });

  it('creates a rich block from heading, price, benefits, and eligibility text', async () => {
    const result = await handleText('create_rich_block', {
      name: 'Gold sponsor',
      x: 24,
      y: 32,
      width: 280,
      heading: 'Gold Sponsor',
      price: '$5,000',
      benefits: ['Logo placement', 'Reserved table'],
      eligibility: 'Limited availability',
    });

    expect(result.name).toBe('Gold sponsor');
    expect(result.x).toBe(24);
    expect(result.y).toBe(32);
    expect(result.width).toBe(280);
    expect(result.height).toBeGreaterThan(0);
    expect(result.children.map((child: any) => child.name)).toEqual([
      'Gold sponsor heading',
      'Gold sponsor price',
      'Gold sponsor body',
      'Gold sponsor note',
    ]);
    expect((globalThis as any).figma.currentPage.children).toHaveLength(4);
  });
});
