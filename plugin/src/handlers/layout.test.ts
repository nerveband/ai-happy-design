import { beforeEach, describe, expect, it, vi } from 'vitest';
import { handleLayout } from './layout';

vi.mock('../utils/fonts', () => ({
  loadFont: vi.fn(async () => undefined),
  resolveFontFamily: vi.fn((family: string) => family),
}));

vi.mock('../utils/getNode', () => ({
  getNodeById: vi.fn(),
  getSceneNodeById: vi.fn(),
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

function createMockFrame() {
  return {
    id: '1:' + nextNodeId++,
    type: 'FRAME',
    name: 'Frame',
    x: 0,
    y: 0,
    width: 0,
    height: 0,
    children: [] as any[],
    fills: [] as any[],
    strokes: [] as any[],
    appendChild(node: any) {
      this.children.push(node);
      node.parent = this;
    },
    resize(width: number, height: number) {
      this.width = width;
      this.height = height;
    },
  };
}

function createMockText() {
  return {
    id: '1:' + nextNodeId++,
    type: 'TEXT',
    name: 'Text',
    width: 0,
    height: 24,
    fontSize: 16,
    characters: '',
    resize(width: number) {
      this.width = width;
      this.height = Math.max(24, Math.ceil(String(this.characters).length / 24) * this.fontSize * 1.2);
    },
  };
}

beforeEach(() => {
  nextNodeId = 1;
  (globalThis as any).figma = {
    currentPage: createMockPage(),
    createFrame: vi.fn(createMockFrame),
    createText: vi.fn(createMockText),
  };
});

describe('high-level layout commands', () => {
  it('creates a pricing grid with positioned card frames and semantic text layers', async () => {
    const result = await handleLayout('pricing_grid', {
      x: 12,
      y: 20,
      width: 600,
      columns: 2,
      gap: 20,
      cards: [
        { tier: 'Gold', price: '$5,000', benefits: ['Logo', 'Table'], eligibility: 'Sponsors' },
        { tier: 'Silver', price: '$2,500', benefits: ['Logo'] },
      ],
    });

    expect(result.count).toBe(2);
    expect(result.rows).toBe(1);
    expect(result.cards[0]).toMatchObject({ name: 'Gold', x: 12, y: 20, width: 290 });
    expect(result.cards[1]).toMatchObject({ name: 'Silver', x: 322, y: 20, width: 290 });

    const pageChildren = (globalThis as any).figma.currentPage.children;
    expect(pageChildren).toHaveLength(2);
    expect(pageChildren[0].children.map((child: any) => child.name)).toEqual([
      'Gold title',
      'Gold price',
      'Gold benefits',
      'Gold eligibility',
    ]);
    expect(pageChildren[1].children.map((child: any) => child.name)).toEqual([
      'Silver title',
      'Silver price',
      'Silver benefits',
    ]);
  });
});
