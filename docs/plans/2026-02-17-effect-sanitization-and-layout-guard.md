# Effect Sanitization & Layout Guard Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix two bugs — (1) effect append operations fail on nodes with beta effects (GLASS/NOISE/TEXTURE) because read-back objects contain internal properties Figma rejects on write-back, and (2) LLMs incorrectly use `layoutPositioning: "ABSOLUTE"` on children of non-auto-layout frames because catalog guidance conflates the concept with manual x/y positioning.

**Architecture:** Create sanitizer utilities that strip internal-only properties from Figma node arrays (effects, fills) before writing them back. Add a layout guard that silently skips `layoutPositioning` when the parent lacks auto-layout. Fix catalog terminology to prevent LLM confusion.

**Tech Stack:** TypeScript (plugin handlers/utils), Go (catalog_llm.go), vitest (new test infrastructure)

---

### Task 1: Set up vitest for plugin tests

**Files:**
- Modify: `plugin/package.json`
- Create: `plugin/vitest.config.ts`

**Step 1: Install vitest**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npm install --save-dev vitest`

**Step 2: Create vitest config**

Create `plugin/vitest.config.ts`:
```typescript
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: ['src/**/*.test.ts'],
    environment: 'node',
  },
});
```

**Step 3: Add test script to package.json**

Add to the `"scripts"` object in `plugin/package.json`:
```json
"test": "vitest run",
"test:watch": "vitest"
```

**Step 4: Verify vitest runs (no tests yet)**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run`
Expected: exits cleanly with "no test files found" or similar

**Step 5: Commit**

```bash
git add plugin/package.json plugin/vitest.config.ts plugin/package-lock.json
git commit -m "chore: add vitest test infrastructure for plugin"
```

---

### Task 2: Create `sanitizeEffects` utility with tests (TDD)

**Files:**
- Create: `plugin/src/utils/sanitizeEffects.ts`
- Create: `plugin/src/utils/sanitizeEffects.test.ts`

**Step 1: Write the failing tests**

Create `plugin/src/utils/sanitizeEffects.test.ts`:
```typescript
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
      // internal property that Figma might add on read-back
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
```

**Step 2: Run tests to verify they fail**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/sanitizeEffects.test.ts`
Expected: FAIL — module not found

**Step 3: Write the implementation**

Create `plugin/src/utils/sanitizeEffects.ts`:
```typescript
// Known-safe property lists per effect type.
// When reading node.effects and writing back, Figma's internal representation
// may include extra properties (boundVariables, _internal, etc.) that cause
// "Invalid format for effects" when re-assigned. This sanitizer strips them.

const SHADOW_KEYS = ['type', 'color', 'offset', 'radius', 'spread', 'visible', 'blendMode'] as const;
const BLUR_KEYS = ['type', 'blurType', 'radius', 'visible'] as const;
const NOISE_KEYS = ['type', 'noiseType', 'color', 'noiseSize', 'density', 'visible', 'secondaryColor', 'opacity', 'blendMode'] as const;
const TEXTURE_KEYS = ['type', 'noiseSize', 'radius', 'clipToShape', 'visible'] as const;
const GLASS_KEYS = ['type', 'visible', 'lightIntensity', 'lightAngle', 'refraction', 'depth', 'dispersion', 'radius'] as const;

function pick(obj: any, keys: readonly string[]): any {
  const out: any = {};
  for (const k of keys) {
    if (obj[k] !== undefined) out[k] = obj[k];
  }
  return out;
}

export function sanitizeEffect(e: any): any {
  switch (e.type) {
    case 'DROP_SHADOW':
    case 'INNER_SHADOW':
      return pick(e, SHADOW_KEYS);
    case 'LAYER_BLUR':
    case 'BACKGROUND_BLUR':
      return pick(e, BLUR_KEYS);
    case 'NOISE':
      return pick(e, NOISE_KEYS);
    case 'TEXTURE':
      return pick(e, TEXTURE_KEYS);
    case 'GLASS':
      return pick(e, GLASS_KEYS);
    default:
      // Unknown future type — pass through as-is for forward compatibility.
      // This is better than throwing, which would break addShadow on nodes
      // with effect types we haven't seen yet.
      return { ...e };
  }
}

export function sanitizeEffects(effects: readonly any[]): any[] {
  return effects.map(sanitizeEffect);
}
```

**Step 4: Run tests to verify they pass**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/sanitizeEffects.test.ts`
Expected: all tests PASS

**Step 5: Commit**

```bash
git add plugin/src/utils/sanitizeEffects.ts plugin/src/utils/sanitizeEffects.test.ts
git commit -m "feat: add sanitizeEffects utility to strip internal properties from read-back effects"
```

---

### Task 3: Create `sanitizeFills` utility with tests (TDD)

**Files:**
- Create: `plugin/src/utils/sanitizeFills.ts`
- Create: `plugin/src/utils/sanitizeFills.test.ts`

**Step 1: Write the failing tests**

Create `plugin/src/utils/sanitizeFills.test.ts`:
```typescript
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
```

**Step 2: Run tests to verify they fail**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/sanitizeFills.test.ts`
Expected: FAIL — module not found

**Step 3: Write the implementation**

Create `plugin/src/utils/sanitizeFills.ts`:
```typescript
// Known-safe property lists per fill type.
// Figma's internal read-back may include boundVariables, imageRef, and other
// properties that cause errors when re-assigned via node.fills = [...].

const SOLID_KEYS = ['type', 'color', 'opacity', 'visible', 'blendMode'] as const;
const GRADIENT_KEYS = ['type', 'gradientStops', 'gradientTransform', 'opacity', 'visible', 'blendMode'] as const;
const IMAGE_KEYS = ['type', 'imageHash', 'scaleMode', 'opacity', 'visible', 'blendMode',
  'imageTransform', 'scalingFactor', 'rotation', 'filters'] as const;

function pick(obj: any, keys: readonly string[]): any {
  const out: any = {};
  for (const k of keys) {
    if (obj[k] !== undefined) out[k] = obj[k];
  }
  return out;
}

export function sanitizeFill(f: any): any {
  switch (f.type) {
    case 'SOLID':
      return pick(f, SOLID_KEYS);
    case 'GRADIENT_LINEAR':
    case 'GRADIENT_RADIAL':
    case 'GRADIENT_ANGULAR':
    case 'GRADIENT_DIAMOND':
      return pick(f, GRADIENT_KEYS);
    case 'IMAGE':
      return pick(f, IMAGE_KEYS);
    default:
      return { ...f };
  }
}

export function sanitizeFills(fills: readonly any[]): any[] {
  return fills.map(sanitizeFill);
}
```

**Step 4: Run tests to verify they pass**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/sanitizeFills.test.ts`
Expected: all tests PASS

**Step 5: Commit**

```bash
git add plugin/src/utils/sanitizeFills.ts plugin/src/utils/sanitizeFills.test.ts
git commit -m "feat: add sanitizeFills utility to strip internal properties from read-back fills"
```

---

### Task 4: Create layout guard utility with tests (TDD)

**Files:**
- Create: `plugin/src/utils/layoutGuard.ts`
- Create: `plugin/src/utils/layoutGuard.test.ts`

**Step 1: Write the failing tests**

Create `plugin/src/utils/layoutGuard.test.ts`:
```typescript
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
```

**Step 2: Run tests to verify they fail**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/layoutGuard.test.ts`
Expected: FAIL — module not found

**Step 3: Write the implementation**

Create `plugin/src/utils/layoutGuard.ts`:
```typescript
// Check if a parent node has auto-layout enabled.
// layoutPositioning: "ABSOLUTE" is only valid inside auto-layout frames.
// Setting it on a child of a non-auto-layout frame causes a Figma error.

export function parentHasAutoLayout(parent: any): boolean {
  if (!parent) return false;
  if (!('layoutMode' in parent)) return false;
  return parent.layoutMode !== 'NONE';
}
```

**Step 4: Run tests to verify they pass**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/layoutGuard.test.ts`
Expected: all tests PASS

**Step 5: Commit**

```bash
git add plugin/src/utils/layoutGuard.ts plugin/src/utils/layoutGuard.test.ts
git commit -m "feat: add parentHasAutoLayout guard for layoutPositioning safety"
```

---

### Task 5: Apply `sanitizeEffects` to all read-modify-write paths in `effect.ts`

**Files:**
- Modify: `plugin/src/handlers/effect.ts`

**Step 1: Add import**

At the top of `plugin/src/handlers/effect.ts` (line 1, after the existing import), add:
```typescript
import { sanitizeEffects } from '../utils/sanitizeEffects';
```

**Step 2: Replace `node.effects.slice()` in `addShadow` (line 103)**

Change:
```typescript
  const nextEffects = node.effects.slice();
```
To:
```typescript
  const nextEffects = sanitizeEffects(node.effects);
```

**Step 3: Replace in `addBlur` (line 118)**

Change:
```typescript
  const nextEffects = node.effects.slice();
```
To:
```typescript
  const nextEffects = sanitizeEffects(node.effects);
```

**Step 4: Replace in `removeEffect` (line 144)**

Change:
```typescript
  const effects = node.effects.slice();
```
To:
```typescript
  const effects = sanitizeEffects(node.effects);
```

**Step 5: Replace in `addNoise` (lines 192 and 203)**

Change both occurrences:
```typescript
  var nextEffects = node.effects.slice();
```
and:
```typescript
    var fallbackEffects = node.effects.slice();
```
To:
```typescript
  var nextEffects = sanitizeEffects(node.effects);
```
and:
```typescript
    var fallbackEffects = sanitizeEffects(node.effects);
```

**Step 6: Replace in `addTexture` (line 229)**

Change:
```typescript
  var nextEffects = node.effects.slice();
```
To:
```typescript
  var nextEffects = sanitizeEffects(node.effects);
```

**Step 7: Replace in `addNativeGlass` (line 254)**

Change:
```typescript
  var nextEffects = node.effects.slice();
```
To:
```typescript
  var nextEffects = sanitizeEffects(node.effects);
```

**Step 8: Replace in `applyGlass` (lines 299 and 335)**

Change both occurrences:
```typescript
  var nextEffects = node.effects.slice();
```
and:
```typescript
  var simEffects: any[] = node.effects.slice();
```
To:
```typescript
  var nextEffects = sanitizeEffects(node.effects);
```
and:
```typescript
  var simEffects: any[] = sanitizeEffects(node.effects);
```

**Step 9: Update `buildEffect` default case (line 386)**

Change the default case from throwing to passing through:
```typescript
    default:
      // Unknown/beta effect type — pass through as-is for forward compatibility
      return { ...e };
```

**Step 10: Verify plugin builds**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx tsc -p tsconfig.json --noEmit`
Expected: no type errors (we use `any` types throughout)

**Step 11: Run all tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run`
Expected: all tests PASS

**Step 12: Commit**

```bash
git add plugin/src/handlers/effect.ts
git commit -m "fix: sanitize effects before write-back to prevent beta effect serialization errors"
```

---

### Task 6: Apply `sanitizeFills` to read-modify-write paths in `paint.ts`

**Files:**
- Modify: `plugin/src/handlers/paint.ts`

**Step 1: Add import**

At the top of `plugin/src/handlers/paint.ts` (line 1, after the existing import), add:
```typescript
import { sanitizeFills } from '../utils/sanitizeFills';
```

**Step 2: Replace in `addFill` (line 219)**

Change:
```typescript
  const fills = (node.fills as Paint[]).slice();
```
To:
```typescript
  const fills = sanitizeFills(node.fills as Paint[]);
```

**Step 3: Replace in `removeFill` (line 229)**

Change:
```typescript
  const fills = (node.fills as Paint[]).slice();
```
To:
```typescript
  const fills = sanitizeFills(node.fills as Paint[]);
```

**Step 4: Verify plugin builds**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx tsc -p tsconfig.json --noEmit`
Expected: no type errors

**Step 5: Commit**

```bash
git add plugin/src/handlers/paint.ts
git commit -m "fix: sanitize fills before write-back to prevent internal property errors"
```

---

### Task 7: Add layout guard to `shape.ts`

**Files:**
- Modify: `plugin/src/handlers/shape.ts:36-42`

**Step 1: Add import**

At the top of `plugin/src/handlers/shape.ts`, add:
```typescript
import { parentHasAutoLayout } from '../utils/layoutGuard';
```

**Step 2: Guard `layoutPositioning` in `applyLayoutProps` (line 41)**

Change line 41:
```typescript
  if (params.layoutPositioning !== undefined) (node as any).layoutPositioning = params.layoutPositioning;
```
To:
```typescript
  if (params.layoutPositioning !== undefined) {
    // layoutPositioning: "ABSOLUTE" is only valid inside auto-layout frames.
    // Silently skip if parent has no auto-layout to prevent Figma errors.
    if (params.layoutPositioning !== 'ABSOLUTE' || parentHasAutoLayout(node.parent)) {
      (node as any).layoutPositioning = params.layoutPositioning;
    }
  }
```

This guard logic:
- Always allows setting `layoutPositioning` to values other than `ABSOLUTE` (e.g. `AUTO`)
- Only sets `ABSOLUTE` when the parent actually has auto-layout
- Silently skips `ABSOLUTE` on non-auto-layout parents (no error, no broken design)

**Step 3: Verify plugin builds**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx tsc -p tsconfig.json --noEmit`
Expected: no type errors

**Step 4: Commit**

```bash
git add plugin/src/handlers/shape.ts
git commit -m "fix: guard layoutPositioning:ABSOLUTE to prevent errors on non-auto-layout parents"
```

---

### Task 8: Fix catalog terminology in `catalog_llm.go`

**Files:**
- Modify: `internal/tools/catalog_llm.go:278` (decision tree)
- Modify: `internal/tools/catalog_llm.go:387-407` (absolutePositioning section)
- Modify: `internal/tools/catalog_llm.go:695` (composition tips)

**Step 1: Fix the decision tree (line 278)**

Change:
```go
"_decisionTree": "DEFAULT: Use auto-layout (flexbox). Think in rows and columns. Figma auto-layout IS CSS flexbox — you already know it. Create frames with layoutMode, itemSpacing, padding, and alignment in a SINGLE create_frame call. Use absolute positioning ONLY for decorative overlays (layoutPositioning: ABSOLUTE). After creating, run layout.check_overlaps to verify no elements overlap.",
```
To:
```go
"_decisionTree": "DEFAULT: Use auto-layout (flexbox). Think in rows and columns. Figma auto-layout IS CSS flexbox — you already know it. Create frames with layoutMode, itemSpacing, padding, and alignment in a SINGLE create_frame call. For decorative overlays inside auto-layout frames, use layoutPositioning:ABSOLUTE on the child to exempt it from the flow. After creating, run layout.check_overlaps to verify no elements overlap.",
```

**Step 2: Rename and rewrite the absolutePositioning section (lines 387-407)**

Replace the entire `"absolutePositioning"` key and its value with `"manualPositioning"`:
```go
"manualPositioning": map[string]interface{}{
    "_overview":      "For non-auto-layout frames, children are positioned by x/y coordinates relative to the parent. No special property is needed — this is Figma's default behavior. Do NOT set layoutPositioning:ABSOLUTE on children of non-auto-layout frames (it will error).",
    "rule":           "Non-auto-layout frame + children with x/y = manual positioning. Auto-layout frame + layoutPositioning:ABSOLUTE = exempt child from flow (decorative overlays only).",
    "WARNING":        "layoutPositioning:ABSOLUTE is ONLY valid inside auto-layout parents. Setting it on a child of a frame with layoutMode:NONE causes an error. If the parent has no auto-layout, just use x/y — no extra property needed.",
    "colorShortcut":  "node.create_frame, shape.create_rectangle, and text.create all accept a 'color' param (hex string like '#0F0F23' or rgb object). The param name is 'color', NOT 'fillColor' or 'backgroundColor'. This is a shortcut that avoids a separate paint.set_solid call.",
    "howItWorks": []string{
        "1. Create root frame (NO auto-layout): node.create_frame {name, x, y, width, height, color: '#0F0F23'}",
        "2. Add children with parentId and x/y: text.create {text, parentId, x, y, width, fontSize, color: '#FFFFFF'}",
        "3. Children are positioned by x/y automatically — NO layoutPositioning needed.",
        "4. For badges/buttons that need centered text: use auto-layout on the badge frame itself.",
        "5. x/y are RELATIVE to the parent frame, not the page.",
        "6. Text width: use the 'width' parameter on text.create for text wrapping control.",
    },
    "coordinatePlanning": []string{
        "Plan your layout TOP-DOWN. Write down the y-coordinate of each major element.",
        "Account for text height: height ≈ fontSize * (lineHeight/100) * numLines.",
        "Leave 48-96px gaps between major sections. Use compute_tokens for spacing values.",
        "Cards: use consistent heights for siblings (e.g., all 200px). Place with exact x/y.",
        "Horizontal card rows: cardWidth = (contentWidth - (N-1) * gap) / N, then x[i] = sidePadding + i * (cardWidth + gap).",
    },
    "hybridExample": "Root frame (manual x/y) → Badge (auto-layout for centering) → Hero text (x/y, width param) → Subtitle (x/y, width param) → Cards (x/y) → CTA button (auto-layout for centering)",
},
```

**Step 3: Fix composition tip (line 695)**

Change:
```go
"Use layoutPositioning:ABSOLUTE ONLY for decorative overlays (circles, stripes, gradients).",
```
To:
```go
"layoutPositioning:ABSOLUTE is ONLY for decorative overlays inside AUTO-LAYOUT frames. Never use it on children of non-auto-layout frames — they already position by x/y.",
```

**Step 4: Run Go tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/tools/...`
Expected: all tests PASS

**Step 5: Commit**

```bash
git add internal/tools/catalog_llm.go
git commit -m "fix: clarify catalog terminology — rename absolutePositioning, warn about auto-layout prerequisite"
```

---

### Task 9: Build, sign, and verify

**Do NOT execute this task until all previous tasks pass.**

**Step 1: Run all plugin tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run`
Expected: all tests PASS

**Step 2: Run all Go tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./...`
Expected: all tests PASS

**Step 3: Build and deploy**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && make deploy`

**Step 4: Verify the plugin loads**

Reopen the Figma plugin to load the new `code.js`.

**Step 5: Final commit (if any unstaged changes)**

```bash
git status
```
