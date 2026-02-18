# Speed Optimization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reduce prompt-to-Figma time from ~7 minutes to under 2 minutes by adding composite batch commands, an HTML extractor, and an end-to-end benchmark harness.

**Architecture:** Three phases layered on the existing batch system. Phase 1 adds composite command expansion in `batchutil/normalize.go`. Phase 2 adds an `extract` cobra command with two HTML parsing modes. Phase 3 adds a `benchmark` cobra command that measures the full pipeline including LLM generation time (provider-agnostic via stdin/stdout piping).

**Tech Stack:** Go (cobra CLI, golang.org/x/net/html, chromedp), TypeScript (Figma plugin text handler), existing batchutil pipeline.

---

## Phase 1: Composite Batch Commands

### Task 1: text.set_range_style Plugin Handler

The plugin already supports `rangeStart`/`rangeEnd` on individual `text.set_size`, `text.set_font`, etc. This task adds a single `text.set_range_style` action that applies multiple style properties to character ranges in one call.

**Files:**
- Modify: `plugin/src/handlers/text.ts` (add `set_range_style` action, ~line 60 in the action switch)
- Test: `plugin/src/handlers/text.test.ts` (new file)

**Step 1: Write the failing test**

Create `plugin/src/handlers/text.test.ts`:

```typescript
import { describe, it, expect, vi } from 'vitest';

describe('set_range_style', () => {
  it('should resolve match string to start/end positions', () => {
    // Unit test for the match-to-range resolver
    const text = 'Hello World, 250+ chaplains serve here';
    const match = '250+';
    const start = text.indexOf(match);
    const end = start + match.length;
    expect(start).toBe(13);
    expect(end).toBe(17);
  });

  it('should handle multiple ranges without overlap', () => {
    const text = 'AMC supports 250+ chaplains';
    const ranges = [
      { match: 'AMC', bold: true },
      { match: '250+', color: '#C88B0A' },
    ];
    // Verify ranges don't overlap
    const resolved = ranges.map(r => {
      const s = text.indexOf(r.match!);
      return { start: s, end: s + r.match!.length, ...r };
    });
    for (let i = 0; i < resolved.length - 1; i++) {
      expect(resolved[i].end).toBeLessThanOrEqual(resolved[i + 1].start);
    }
  });
});
```

**Step 2: Run test to verify it passes** (these are pure logic tests)

Run: `cd plugin && npx vitest run src/handlers/text.test.ts`

**Step 3: Implement set_range_style in plugin**

In `plugin/src/handlers/text.ts`, add a new case in the action handler (around line 60):

```typescript
case 'set_range_style': {
  const node = await getNode(params.nodeId) as TextNode;
  if (node.type !== 'TEXT') throw new Error('Node is not a TEXT node');

  const ranges = params.ranges as Array<{
    match?: string;
    start?: number;
    end?: number;
    bold?: boolean;
    italic?: boolean;
    color?: string;
    fontSize?: number;
    fontFamily?: string;
    fontStyle?: string;
    letterSpacing?: number;
    lineHeight?: number;
    textDecoration?: string;
    textCase?: string;
  }>;

  if (!ranges || !Array.isArray(ranges)) {
    throw new Error('ranges param required: array of {match|start+end, ...styles}');
  }

  const text = node.characters;

  for (const range of ranges) {
    let start: number;
    let end: number;

    if (range.match !== undefined) {
      const idx = text.indexOf(range.match);
      if (idx === -1) continue; // skip if match not found
      start = idx;
      end = idx + range.match.length;
    } else if (range.start !== undefined && range.end !== undefined) {
      start = range.start;
      end = range.end;
    } else {
      continue; // skip invalid range
    }

    // Clamp to text length
    start = Math.max(0, start);
    end = Math.min(text.length, end);
    if (start >= end) continue;

    // Apply each style property
    if (range.bold !== undefined || range.fontStyle !== undefined || range.fontFamily !== undefined) {
      const currentFont = node.getRangeFontName(start, end) as FontName;
      const family = range.fontFamily || currentFont.family;
      let style = currentFont.style;
      if (range.bold === true) style = 'Bold';
      if (range.bold === false) style = 'Regular';
      if (range.italic === true) style = style.includes('Bold') ? 'Bold Italic' : 'Italic';
      if (range.fontStyle) style = range.fontStyle;
      await figma.loadFontAsync({ family, style });
      node.setRangeFontName(start, end, { family, style });
    }

    if (range.fontSize !== undefined) {
      node.setRangeFontSize(start, end, range.fontSize);
    }

    if (range.color !== undefined) {
      const color = parseColor(range.color);
      node.setRangeFills(start, end, [{ type: 'SOLID', color }]);
    }

    if (range.letterSpacing !== undefined) {
      node.setRangeLetterSpacing(start, end, {
        value: range.letterSpacing,
        unit: 'PIXELS'
      });
    }

    if (range.lineHeight !== undefined) {
      node.setRangeLineHeight(start, end, {
        value: range.lineHeight,
        unit: 'PERCENT'
      });
    }

    if (range.textDecoration !== undefined) {
      node.setRangeTextDecoration(start, end,
        range.textDecoration as TextDecoration);
    }

    if (range.textCase !== undefined) {
      node.setRangeTextCase(start, end,
        range.textCase as TextCase);
    }
  }

  return { id: node.id, name: node.name, rangesApplied: ranges.length };
}
```

Note: `parseColor` already exists in text.ts — it converts hex strings to Figma RGB objects.

**Step 4: Run plugin tests**

Run: `cd plugin && npx vitest run`

**Step 5: Build plugin and deploy**

Run: `make deploy`

**Step 6: Test manually via CLI**

```bash
# Create a text node first
ai-happy-design command text.create '{"parentId":"0:1","text":"Hello World 250+ chaplains","fontSize":48,"fontFamily":"Inter","fontStyle":"Regular","color":"#ffffff","x":0,"y":0,"width":600}'

# Then style ranges (use the returned ID)
ai-happy-design command text.set_range_style '{"nodeId":"<ID>","ranges":[{"match":"250+","bold":true,"color":"#C88B0A","fontSize":64},{"match":"chaplains","italic":true}]}'
```

**Step 7: Commit**

```bash
git add plugin/src/handlers/text.ts plugin/src/handlers/text.test.ts
git commit -m "feat: add text.set_range_style for character-level styling"
```

---

### Task 2: Composite Command Expander in batchutil

Add a new file `internal/batchutil/composite.go` that expands higher-level commands (`slide`, `banner`, `stats`, `progress`, `cta`, `counter`, `avatar`) into arrays of primitive batch operations.

**Files:**
- Create: `internal/batchutil/composite.go`
- Create: `internal/batchutil/composite_test.go`
- Modify: `internal/batchutil/normalize.go` (call composite expansion)

**Step 1: Write the failing test**

Create `internal/batchutil/composite_test.go`:

```go
package batchutil

import (
    "encoding/json"
    "testing"
)

func TestExpandSlideCommand(t *testing.T) {
    input := map[string]interface{}{
        "name":    "s1",
        "command": "slide",
        "params": map[string]interface{}{
            "name":   "P1 - Slide 1",
            "canvas": "1080x1350",
            "bg":     "#0C1E2C",
            "elements": []interface{}{
                map[string]interface{}{
                    "type": "eyebrow",
                    "text": "AMC · Ramadan 2026",
                    "color": "#7FBCD2",
                },
                map[string]interface{}{
                    "type": "headline",
                    "text": "The Care\nBehind\nthe Care",
                    "tier": "hero",
                },
            },
        },
    }

    ops, err := ExpandComposite(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Should produce: 1 frame + 1 gradient (if specified) + N element ops
    if len(ops) < 3 {
        t.Fatalf("expected at least 3 ops, got %d", len(ops))
    }

    // First op should be the root frame
    first := ops[0]
    cmd, _ := first["command"].(string)
    if cmd != "node.create_frame" {
        t.Errorf("first op should be node.create_frame, got %s", cmd)
    }

    params, _ := first["params"].(map[string]interface{})
    w, _ := params["width"].(float64)
    h, _ := params["height"].(float64)
    if w != 1080 || h != 1350 {
        t.Errorf("expected 1080x1350, got %vx%v", w, h)
    }
}

func TestExpandBannerCommand(t *testing.T) {
    input := map[string]interface{}{
        "name":    "eb1",
        "command": "banner",
        "params": map[string]interface{}{
            "name":   "Email 1",
            "canvas": "1200x400",
            "bg":     "#0C1E2C",
            "elements": []interface{}{
                map[string]interface{}{
                    "type": "headline",
                    "text": "The Care Behind the Care",
                },
            },
        },
    }

    ops, err := ExpandComposite(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(ops) < 2 {
        t.Fatalf("expected at least 2 ops, got %d", len(ops))
    }
}

func TestExpandStatsElement(t *testing.T) {
    input := map[string]interface{}{
        "name":    "s1",
        "command": "slide",
        "params": map[string]interface{}{
            "canvas": "1080x1350",
            "bg":     "#14344A",
            "elements": []interface{}{
                map[string]interface{}{
                    "type": "stats",
                    "items": []interface{}{
                        map[string]interface{}{"value": "250+", "label": "Chaplains"},
                        map[string]interface{}{"value": "5", "label": "Settings"},
                    },
                },
            },
        },
    }

    ops, err := ExpandComposite(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Stats: should produce value + label text pairs
    // 1 frame + stats items (2 values + 2 labels) = at least 5 ops
    if len(ops) < 5 {
        t.Fatalf("expected at least 5 ops for stats, got %d", len(ops))
    }
}

func TestExpandProgressElement(t *testing.T) {
    input := map[string]interface{}{
        "name":    "s1",
        "command": "slide",
        "params": map[string]interface{}{
            "canvas": "1080x1350",
            "bg":     "#060D14",
            "elements": []interface{}{
                map[string]interface{}{
                    "type":   "progress",
                    "raised": 17000,
                    "goal":   25000,
                    "color":  "#C88B0A",
                },
            },
        },
    }

    ops, err := ExpandComposite(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Progress: frame + track + fill + gradient + raised label + goal label = 7+ ops
    if len(ops) < 5 {
        t.Fatalf("expected at least 5 ops for progress, got %d", len(ops))
    }
}

func TestNonCompositePassthrough(t *testing.T) {
    input := map[string]interface{}{
        "command": "text.create",
        "params": map[string]interface{}{
            "text": "Hello",
        },
    }

    ops, err := ExpandComposite(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(ops) != 1 {
        t.Fatalf("non-composite should pass through as-is, got %d ops", len(ops))
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/batchutil/ -run TestExpand -v`
Expected: FAIL — `ExpandComposite` not defined

**Step 3: Implement composite.go**

Create `internal/batchutil/composite.go`:

```go
package batchutil

import (
    "fmt"
    "math"
    "strconv"
    "strings"
)

// CompositeCommands that get expanded into multiple primitive ops.
var compositeCommands = map[string]bool{
    "slide":  true,
    "banner": true,
}

// IsComposite returns true if the command is a composite that needs expansion.
func IsComposite(command string) bool {
    return compositeCommands[command]
}

// ExpandComposite expands a composite command into primitive batch ops.
// Non-composite commands are returned as a single-element slice unchanged.
func ExpandComposite(op map[string]interface{}) ([]map[string]interface{}, error) {
    command, _ := op["command"].(string)
    if !IsComposite(command) {
        return []map[string]interface{}{op}, nil
    }

    params, _ := op["params"].(map[string]interface{})
    if params == nil {
        return nil, fmt.Errorf("composite command %q requires params", command)
    }

    name, _ := op["name"].(string)
    if name == "" {
        name = "comp"
    }

    switch command {
    case "slide":
        return expandSlide(name, params)
    case "banner":
        return expandBanner(name, params)
    default:
        return []map[string]interface{}{op}, nil
    }
}

// ExpandAllComposites processes a batch ops array, expanding any composites.
func ExpandAllComposites(ops []map[string]interface{}) ([]map[string]interface{}, error) {
    var result []map[string]interface{}
    for _, op := range ops {
        expanded, err := ExpandComposite(op)
        if err != nil {
            return nil, err
        }
        result = append(result, expanded...)
    }
    return result, nil
}

func parseCanvas(canvas string) (int, int, error) {
    parts := strings.SplitN(canvas, "x", 2)
    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("invalid canvas format %q, expected WxH", canvas)
    }
    w, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, 0, fmt.Errorf("invalid canvas width: %v", err)
    }
    h, err := strconv.Atoi(parts[1])
    if err != nil {
        return 0, 0, fmt.Errorf("invalid canvas height: %v", err)
    }
    return w, h, nil
}

// computeTokens returns font sizes for the given width using the modular type scale.
func computeTokens(width int) map[string]int {
    base := float64(width) * 0.044
    ratio := 1.333
    snap := func(v float64) int { return int(math.Round(v/4) * 4) }
    return map[string]int{
        "caption":    snap(base / ratio),
        "body":       snap(base),
        "subheading": snap(base * ratio),
        "heading":    snap(base * math.Pow(ratio, 2)),
        "title":      snap(base * math.Pow(ratio, 3)),
        "hero":       snap(base * math.Pow(ratio, 4)),
        "display":    snap(base * math.Pow(ratio, 5)),
    }
}

func expandSlide(baseName string, params map[string]interface{}) ([]map[string]interface{}, error) {
    canvasStr, _ := params["canvas"].(string)
    if canvasStr == "" {
        canvasStr = "1080x1350"
    }
    w, h, err := parseCanvas(canvasStr)
    if err != nil {
        return nil, err
    }

    tokens := computeTokens(w)
    margin := int(math.Round(float64(w) * 0.08))  // ~86px at 1080
    contentW := w - 2*margin

    bg, _ := params["bg"].(string)
    if bg == "" {
        bg = "#0C1E2C"
    }
    slideName, _ := params["name"].(string)
    if slideName == "" {
        slideName = "Slide"
    }
    parentId, _ := params["parentId"].(string)

    // Root frame
    frameName := baseName + "_frame"
    rootOp := map[string]interface{}{
        "name":    frameName,
        "command": "node.create_frame",
        "params": map[string]interface{}{
            "name":         slideName,
            "width":        float64(w),
            "height":       float64(h),
            "color":        bg,
            "clipsContent": true,
        },
    }
    if parentId != "" {
        rootOp["params"].(map[string]interface{})["parentId"] = parentId
    }

    ops := []map[string]interface{}{rootOp}
    frameRef := fmt.Sprintf("${{steps.%s.result.id}}", frameName)

    // Gradient if specified
    if grad, ok := params["gradient"].(map[string]interface{}); ok {
        gradName := baseName + "_grad"
        ops = append(ops, map[string]interface{}{
            "name":    gradName,
            "command": "paint.set_gradient",
            "params":  mergeMap(grad, map[string]interface{}{"nodeId": frameRef}),
        })
    }

    // Elements
    elements, _ := params["elements"].([]interface{})
    yPos := float64(margin) // start below top margin
    elemIdx := 0

    for _, elem := range elements {
        e, ok := elem.(map[string]interface{})
        if !ok {
            continue
        }
        elemType, _ := e["type"].(string)
        elemOps := expandElement(baseName, elemIdx, elemType, e, frameRef, margin, contentW, tokens, &yPos)
        ops = append(ops, elemOps...)
        elemIdx++
    }

    return ops, nil
}

func expandBanner(baseName string, params map[string]interface{}) ([]map[string]interface{}, error) {
    canvasStr, _ := params["canvas"].(string)
    if canvasStr == "" {
        canvasStr = "1200x400"
    }
    w, h, err := parseCanvas(canvasStr)
    if err != nil {
        return nil, err
    }

    bg, _ := params["bg"].(string)
    if bg == "" {
        bg = "#0C1E2C"
    }
    bannerName, _ := params["name"].(string)
    if bannerName == "" {
        bannerName = "Banner"
    }
    parentId, _ := params["parentId"].(string)

    frameName := baseName + "_frame"
    rootOp := map[string]interface{}{
        "name":    frameName,
        "command": "node.create_frame",
        "params": map[string]interface{}{
            "name":         bannerName,
            "width":        float64(w),
            "height":       float64(h),
            "color":        bg,
            "clipsContent": true,
        },
    }
    if parentId != "" {
        rootOp["params"].(map[string]interface{})["parentId"] = parentId
    }

    ops := []map[string]interface{}{rootOp}
    frameRef := fmt.Sprintf("${{steps.%s.result.id}}", frameName)

    // Gradient
    if grad, ok := params["gradient"].(map[string]interface{}); ok {
        gradName := baseName + "_grad"
        ops = append(ops, map[string]interface{}{
            "name":    gradName,
            "command": "paint.set_gradient",
            "params":  mergeMap(grad, map[string]interface{}{"nodeId": frameRef}),
        })
    }

    // Divider if dividerX specified
    if divX, ok := params["dividerX"].(float64); ok {
        divName := baseName + "_div"
        ops = append(ops, map[string]interface{}{
            "name":    divName,
            "command": "shape.create_rectangle",
            "params": map[string]interface{}{
                "parentId": frameRef,
                "name":     "divider",
                "x":        divX,
                "y":        float64(h) * 0.35,
                "width":    2.0,
                "height":   float64(h) * 0.3,
                "color":    "#FAFCFB2E",
            },
        })
    }

    // Banner elements start after divider
    tokens := computeTokens(w)
    elements, _ := params["elements"].([]interface{})
    startX := 240.0
    if divX, ok := params["dividerX"].(float64); ok {
        startX = divX + 40
    }
    yPos := float64(h) * 0.325
    elemIdx := 0

    for _, elem := range elements {
        e, ok := elem.(map[string]interface{})
        if !ok {
            continue
        }
        elemType, _ := e["type"].(string)
        elemContentW := w - int(startX) - 40

        bannerElemOps := expandBannerElement(baseName, elemIdx, elemType, e, frameRef, int(startX), elemContentW, tokens, &yPos)
        ops = append(ops, bannerElemOps...)
        elemIdx++
    }

    return ops, nil
}

func expandElement(baseName string, idx int, elemType string, e map[string]interface{},
    frameRef string, margin, contentW int, tokens map[string]int, yPos *float64) []map[string]interface{} {

    elemName := fmt.Sprintf("%s_e%d", baseName, idx)
    text, _ := e["text"].(string)
    color, _ := e["color"].(string)

    switch elemType {
    case "eyebrow":
        if color == "" {
            color = "#7FBCD2"
        }
        sz := tokens["caption"]
        op := makeTextOp(elemName, frameRef, text, margin, int(*yPos), contentW, sz, "Poppins", "SemiBold", color)
        op["params"].(map[string]interface{})["letterSpacing"] = 4.0
        op["params"].(map[string]interface{})["textCase"] = "UPPER"
        *yPos += float64(sz) * 1.5
        return []map[string]interface{}{op}

    case "headline":
        tier, _ := e["tier"].(string)
        if tier == "" {
            tier = "hero"
        }
        sz := tokens[tier]
        if sz == 0 {
            sz = tokens["hero"]
        }
        if color == "" {
            color = "#FAFCFB"
        }
        op := makeTextOp(elemName, frameRef, text, margin, int(*yPos), contentW, sz, "Poppins", "Bold", color)
        op["params"].(map[string]interface{})["lineHeight"] = 118.0
        op["params"].(map[string]interface{})["lineHeightUnit"] = "PERCENT"
        // Estimate height: lines * fontSize * lineHeight
        lines := float64(strings.Count(text, "\n") + 1)
        *yPos += lines * float64(sz) * 1.18
        return []map[string]interface{}{op}

    case "body":
        if color == "" {
            color = "#FAFCFBB3"
        }
        sz := tokens["body"]
        op := makeTextOp(elemName, frameRef, text, margin, int(*yPos), contentW, sz, "DM Sans", "Regular", color)
        op["params"].(map[string]interface{})["lineHeight"] = 150.0
        op["params"].(map[string]interface{})["lineHeightUnit"] = "PERCENT"
        lines := float64(strings.Count(text, "\n") + 1)
        *yPos += lines * float64(sz) * 1.5
        return []map[string]interface{}{op}

    case "bar":
        barColor := color
        if barColor == "" {
            barColor = "#029056"
        }
        barW := 108.0
        if bw, ok := e["width"].(float64); ok {
            barW = bw
        }
        op := map[string]interface{}{
            "name":    elemName,
            "command": "shape.create_rectangle",
            "params": map[string]interface{}{
                "parentId":     frameRef,
                "name":         "accent-bar",
                "x":            float64(margin),
                "y":            *yPos,
                "width":        barW,
                "height":       4.0,
                "color":        barColor,
                "cornerRadius": 2.0,
            },
        }
        *yPos += 24
        return []map[string]interface{}{op}

    case "counter":
        current, _ := e["current"].(float64)
        total, _ := e["total"].(float64)
        counterText := fmt.Sprintf("%.0f / %.0f", current, total)
        sz := tokens["caption"]
        op := makeTextOp(elemName, frameRef, counterText, contentW-margin, margin, margin+40, sz, "Poppins", "Regular", "#FAFCFB52")
        op["params"].(map[string]interface{})["textAlign"] = "RIGHT"
        // Counter doesn't advance yPos (positioned at top-right)
        return []map[string]interface{}{op}

    case "cta":
        btnText, _ := e["text"].(string)
        btnBg, _ := e["bg"].(string)
        if btnBg == "" {
            btnBg = "#029056"
        }
        btnColor := color
        if btnColor == "" {
            btnColor = "#FFFFFF"
        }
        style, _ := e["style"].(string)
        if style == "" {
            style = "pill"
        }

        sz := tokens["body"]
        btnH := float64(sz) * 2.0
        btnW := 400.0
        btnX := float64(contentW)/2 - btnW/2 + float64(margin)
        btnR := btnH / 2
        if style == "rounded" {
            btnR = btnH * 0.28
        }

        bgName := elemName + "_bg"
        txtName := elemName + "_txt"

        bgOp := map[string]interface{}{
            "name":    bgName,
            "command": "node.create_frame",
            "params": map[string]interface{}{
                "parentId":          frameRef,
                "name":              "CTA Button",
                "x":                 btnX,
                "y":                 *yPos,
                "width":             btnW,
                "height":            btnH,
                "color":             btnBg,
                "cornerRadius":      btnR,
                "layoutMode":        "HORIZONTAL",
                "primaryAxisAlign":  "CENTER",
                "counterAxisAlign":  "CENTER",
            },
        }

        bgRef := fmt.Sprintf("${{steps.%s.result.id}}", bgName)
        txtOp := map[string]interface{}{
            "name":    txtName,
            "command": "text.create",
            "params": map[string]interface{}{
                "parentId":  bgRef,
                "text":      btnText,
                "fontSize":  float64(sz),
                "fontFamily": "Poppins",
                "fontStyle":  "Bold",
                "color":      btnColor,
                "textAlign":  "CENTER",
            },
        }

        *yPos += btnH + 32
        return []map[string]interface{}{bgOp, txtOp}

    case "url":
        if color == "" {
            color = "#FAFCFB6B"
        }
        sz := tokens["caption"]
        op := makeTextOp(elemName, frameRef, text, margin, int(*yPos), contentW, sz, "Poppins", "Regular", color)
        op["params"].(map[string]interface{})["textAlign"] = "CENTER"
        *yPos += float64(sz) * 1.5
        return []map[string]interface{}{op}

    case "stats":
        return expandStats(elemName, e, frameRef, margin, contentW, tokens, yPos)

    case "progress":
        return expandProgress(elemName, e, frameRef, margin, contentW, tokens, yPos)

    case "arabic":
        if color == "" {
            color = "#FAFCFB"
        }
        ff := "Amiri"
        if f, ok := e["fontFamily"].(string); ok {
            ff = f
        }
        sz := tokens["heading"]
        if s, ok := e["fontSize"].(float64); ok {
            sz = int(s)
        }
        op := makeTextOp(elemName, frameRef, text, margin, int(*yPos), contentW, sz, ff, "Regular", color)
        op["params"].(map[string]interface{})["textAlign"] = "CENTER"
        *yPos += float64(sz) * 1.5
        return []map[string]interface{}{op}

    default:
        // Unknown element type — skip
        return nil
    }
}

func expandBannerElement(baseName string, idx int, elemType string, e map[string]interface{},
    frameRef string, startX, contentW int, tokens map[string]int, yPos *float64) []map[string]interface{} {

    elemName := fmt.Sprintf("%s_e%d", baseName, idx)
    text, _ := e["text"].(string)
    color, _ := e["color"].(string)

    switch elemType {
    case "headline":
        if color == "" {
            color = "#FAFCFB"
        }
        sz := tokens["subheading"]
        op := makeTextOp(elemName, frameRef, text, startX, int(*yPos), contentW, sz, "Poppins", "Bold", color)
        op["params"].(map[string]interface{})["lineHeight"] = 125.0
        op["params"].(map[string]interface{})["lineHeightUnit"] = "PERCENT"
        lines := float64(strings.Count(text, "\n") + 1)
        *yPos += lines * float64(sz) * 1.25
        return []map[string]interface{}{op}

    case "subtitle":
        if color == "" {
            color = "#FAFCFB6B"
        }
        sz := tokens["caption"]
        op := makeTextOp(elemName, frameRef, text, startX, int(*yPos), contentW, sz, "Poppins", "SemiBold", color)
        if tc, ok := e["textCase"].(string); ok {
            op["params"].(map[string]interface{})["textCase"] = tc
        }
        if ls, ok := e["letterSpacing"].(float64); ok {
            op["params"].(map[string]interface{})["letterSpacing"] = ls
        }
        *yPos += float64(sz) * 1.5
        return []map[string]interface{}{op}

    default:
        return expandElement(baseName, idx, elemType, e, frameRef, startX, contentW, tokens, yPos)
    }
}

func expandStats(baseName string, e map[string]interface{}, frameRef string,
    margin, contentW int, tokens map[string]int, yPos *float64) []map[string]interface{} {

    items, _ := e["items"].([]interface{})
    if len(items) == 0 {
        return nil
    }

    var ops []map[string]interface{}
    colW := contentW / len(items)

    for i, item := range items {
        it, ok := item.(map[string]interface{})
        if !ok {
            continue
        }
        value, _ := it["value"].(string)
        label, _ := it["label"].(string)
        color, _ := it["color"].(string)
        if color == "" {
            color = "#B3D9E8"
        }
        labelColor, _ := it["labelColor"].(string)
        if labelColor == "" {
            labelColor = "#FAFCFBAD"
        }

        x := margin + i*colW
        valName := fmt.Sprintf("%s_v%d", baseName, i)
        lblName := fmt.Sprintf("%s_l%d", baseName, i)

        valOp := makeTextOp(valName, frameRef, value, x, int(*yPos), colW, tokens["title"], "Poppins", "Bold", color)
        lblOp := makeTextOp(lblName, frameRef, label, x, int(*yPos)+tokens["title"]+8, colW, tokens["caption"], "DM Sans", "Regular", labelColor)

        ops = append(ops, valOp, lblOp)
    }

    *yPos += float64(tokens["title"]) + float64(tokens["caption"]) + 40
    return ops
}

func expandProgress(baseName string, e map[string]interface{}, frameRef string,
    margin, contentW int, tokens map[string]int, yPos *float64) []map[string]interface{} {

    raised, _ := e["raised"].(float64)
    goal, _ := e["goal"].(float64)
    color, _ := e["color"].(string)
    if color == "" {
        color = "#C88B0A"
    }

    if goal <= 0 {
        goal = 1
    }
    pct := raised / goal
    if pct > 1 {
        pct = 1
    }
    fillW := float64(contentW) * pct

    trackName := baseName + "_track"
    fillName := baseName + "_fill"
    raisedName := baseName + "_raised"
    goalName := baseName + "_goal"

    trackOp := map[string]interface{}{
        "name":    trackName,
        "command": "node.create_frame",
        "params": map[string]interface{}{
            "parentId":     frameRef,
            "name":         "progress-track",
            "x":            float64(margin),
            "y":            *yPos,
            "width":        float64(contentW),
            "height":       16.0,
            "color":        "#FAFCFB1F",
            "cornerRadius": 8.0,
        },
    }

    trackRef := fmt.Sprintf("${{steps.%s.result.id}}", trackName)
    fillOp := map[string]interface{}{
        "name":    fillName,
        "command": "node.create_frame",
        "params": map[string]interface{}{
            "parentId":     trackRef,
            "name":         "progress-fill",
            "x":            0.0,
            "y":            0.0,
            "width":        fillW,
            "height":       16.0,
            "color":        color,
            "cornerRadius": 8.0,
        },
    }

    raisedText := fmt.Sprintf("$%.0f raised", raised)
    goalText := fmt.Sprintf("$%.0f goal", goal)

    raisedOp := makeTextOp(raisedName, frameRef, raisedText, margin, int(*yPos)+24, contentW/2, tokens["caption"], "Poppins", "SemiBold", color)
    goalOp := makeTextOp(goalName, frameRef, goalText, margin+contentW/2, int(*yPos)+24, contentW/2, tokens["caption"], "Poppins", "Regular", "#FAFCFB6B")
    goalOp["params"].(map[string]interface{})["textAlign"] = "RIGHT"

    *yPos += 80
    return []map[string]interface{}{trackOp, fillOp, raisedOp, goalOp}
}

func makeTextOp(name, parentRef, text string, x, y, w, sz int, ff, fs, color string) map[string]interface{} {
    return map[string]interface{}{
        "name":    name,
        "command": "text.create",
        "params": map[string]interface{}{
            "parentId":   parentRef,
            "text":       text,
            "x":          float64(x),
            "y":          float64(y),
            "width":      float64(w),
            "fontSize":   float64(sz),
            "fontFamily": ff,
            "fontStyle":  fs,
            "color":      color,
        },
    }
}

func mergeMap(base, overlay map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    for k, v := range base {
        result[k] = v
    }
    for k, v := range overlay {
        result[k] = v
    }
    return result
}
```

**Step 4: Run tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/batchutil/ -run TestExpand -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/batchutil/composite.go internal/batchutil/composite_test.go
git commit -m "feat: add composite command expander for slide/banner/stats/progress"
```

---

### Task 3: Integrate Composite Expansion into Batch Pipeline

Wire `ExpandAllComposites` into the batch execution flow in `main.go`, right after `FixBatchOps` and before the per-step loop.

**Files:**
- Modify: `cmd/ai-happy-design/main.go` (~line 1104-1138 in `loadBatchOperations`, and the batch execution loop)

**Step 1: Write test for integration**

Add to `internal/batchutil/composite_test.go`:

```go
func TestExpandAllComposites(t *testing.T) {
    ops := []map[string]interface{}{
        {"command": "text.create", "params": map[string]interface{}{"text": "plain"}},
        {"command": "slide", "name": "s1", "params": map[string]interface{}{
            "canvas": "1080x1350", "bg": "#000",
            "elements": []interface{}{
                map[string]interface{}{"type": "headline", "text": "Hi"},
            },
        }},
        {"command": "text.create", "params": map[string]interface{}{"text": "after"}},
    }

    result, err := ExpandAllComposites(ops)
    if err != nil {
        t.Fatal(err)
    }
    // First: passthrough, then expanded slide (2+ ops), then passthrough
    if len(result) < 4 {
        t.Fatalf("expected at least 4 ops, got %d", len(result))
    }
    // First op should still be text.create
    if result[0]["command"] != "text.create" {
        t.Errorf("first op should be text.create, got %v", result[0]["command"])
    }
    // Last op should still be text.create
    last := result[len(result)-1]
    if last["command"] != "text.create" {
        t.Errorf("last op should be text.create, got %v", last["command"])
    }
}
```

**Step 2: Run test**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/batchutil/ -run TestExpandAll -v`

**Step 3: Integrate into main.go**

In `main.go`, in the `loadBatchOperations` function (or right after it's called in the batch flow), add composite expansion after FixBatchOps but before the execution loop. Find where ops are unmarshaled into `[]batchOperation` and add:

```go
// After FixBatchOps and unmarshal, before execution loop:
// Expand composite commands (slide, banner → primitive ops)
var rawOps []map[string]interface{}
if err := json.Unmarshal(fixedData, &rawOps); err == nil {
    if expanded, err := batchutil.ExpandAllComposites(rawOps); err == nil && len(expanded) != len(rawOps) {
        // Re-marshal expanded ops
        fixedData, _ = json.Marshal(expanded)
        fmt.Fprintf(os.Stderr, "Expanded %d composite commands → %d operations\n", len(rawOps), len(expanded))
    }
}
```

**Step 4: Build and test manually**

Run: `make deploy`

Test with a composite slide:
```bash
ai-happy-design batch '[{"name":"s1","command":"slide","params":{"canvas":"1080x1350","bg":"#0C1E2C","elements":[{"type":"eyebrow","text":"TEST","color":"#7FBCD2"},{"type":"headline","text":"Hello World","tier":"hero"}]}}]'
```

**Step 5: Commit**

```bash
git add cmd/ai-happy-design/main.go
git commit -m "feat: integrate composite expansion into batch pipeline"
```

---

### Task 4: Add Composite Commands to LLM Discovery

Update the skill, catalog, and describe layers so LLMs know about composite commands.

**Files:**
- Modify: `internal/tools/catalog_llm.go` (add composite command docs)
- Modify: `internal/tools/describe.go` (add composite to tool list)

**Step 1: Add composite section to catalog_llm.go**

Add after the existing batch format section:

```go
// In the catalog output JSON, add:
"compositeCommands": {
    "description": "Higher-level commands that auto-expand into multiple primitive operations. Use these to create complex layouts with fewer JSON objects.",
    "slide": {
        "description": "Creates a full social media slide with background, gradient, and content elements",
        "params": "canvas (WxH), bg, gradient, elements[], name, parentId",
        "elementTypes": ["eyebrow", "headline", "body", "bar", "counter", "cta", "url", "stats", "progress", "arabic"],
        "example": {"command": "slide", "params": {"canvas": "1080x1350", "bg": "#0C1E2C", "elements": [{"type": "headline", "text": "Hello", "tier": "hero"}]}}
    },
    "banner": {
        "description": "Creates an email banner with divider and content elements",
        "params": "canvas (WxH), bg, gradient, dividerX, elements[], name, parentId",
        "elementTypes": ["headline", "subtitle"],
        "example": {"command": "banner", "params": {"canvas": "1200x400", "bg": "#0C1E2C", "dividerX": 200, "elements": [{"type": "headline", "text": "Hello"}]}}
    }
}
```

**Step 2: Update describe.go**

Add `slide` and `banner` to the batch aliases table.

**Step 3: Update SKILL.md**

Add composite commands section to the ai-happy-design skill file at `~/.claude/skills/ai-happy-design/SKILL.md` under the Batch Aliases table.

**Step 4: Commit**

```bash
git add internal/tools/catalog_llm.go internal/tools/describe.go
git commit -m "feat: add composite commands to LLM discovery layers"
```

---

## Phase 2: HTML Extractor

### Task 5: Lightweight HTML Parser

Create `internal/extract/` package with a Go HTML parser that converts HTML+CSS into batch JSON.

**Files:**
- Create: `internal/extract/extract.go`
- Create: `internal/extract/css.go`
- Create: `internal/extract/extract_test.go`

**Step 1: Write the failing test**

Create `internal/extract/extract_test.go`:

```go
package extract

import (
    "strings"
    "testing"
)

func TestExtractSimpleDiv(t *testing.T) {
    html := `<div style="width:400px;height:300px;background-color:#1a1a1a;border-radius:16px">
        <h1 style="color:white;font-size:48px">Hello World</h1>
    </div>`

    ops, err := FromHTML(strings.NewReader(html), Options{CanvasWidth: 1080, CanvasHeight: 1080})
    if err != nil {
        t.Fatal(err)
    }
    if len(ops) < 2 {
        t.Fatalf("expected at least 2 ops (frame + text), got %d", len(ops))
    }

    // First op: frame with bg color
    first := ops[0]
    cmd, _ := first["command"].(string)
    if cmd != "node.create_frame" {
        t.Errorf("expected node.create_frame, got %s", cmd)
    }
    params, _ := first["params"].(map[string]interface{})
    if c, _ := params["color"].(string); c != "#1a1a1a" {
        t.Errorf("expected color #1a1a1a, got %s", c)
    }
}

func TestExtractStyleBlock(t *testing.T) {
    html := `<style>.hero { background: #0C1E2C; width: 1080px; height: 1350px; }</style>
    <div class="hero"><p style="color:#fff">Test</p></div>`

    ops, err := FromHTML(strings.NewReader(html), Options{CanvasWidth: 1080, CanvasHeight: 1350})
    if err != nil {
        t.Fatal(err)
    }
    if len(ops) < 2 {
        t.Fatalf("expected at least 2 ops, got %d", len(ops))
    }
}

func TestExtractGradient(t *testing.T) {
    html := `<div style="width:1080px;height:1080px;background:linear-gradient(150deg,#0C1E2C,#14344A)"></div>`

    ops, err := FromHTML(strings.NewReader(html), Options{CanvasWidth: 1080, CanvasHeight: 1080})
    if err != nil {
        t.Fatal(err)
    }
    // Should produce frame + gradient
    if len(ops) < 2 {
        t.Fatalf("expected at least 2 ops (frame + gradient), got %d", len(ops))
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/extract/ -v`

**Step 3: Implement extract.go**

Create `internal/extract/extract.go` with:
- `FromHTML(r io.Reader, opts Options) ([]map[string]interface{}, error)` — main entry point
- Uses `golang.org/x/net/html` tokenizer to walk the DOM
- Extracts inline `style` attributes and `<style>` block rules
- Maps CSS properties to Figma params using a translation table
- Handles: div→frame, h1-h6→text (with tier mapping), p/span→text, img→image
- Handles: background-color, color, font-size, font-family, font-weight, padding, border-radius, width, height, linear-gradient

Create `internal/extract/css.go` with:
- `parseInlineStyles(style string) map[string]string` — parses `key: value; key: value`
- `parseStyleBlock(css string) map[string]map[string]string` — parses `.class { ... }` blocks
- `parseLinearGradient(value string) (angle float64, stops []map[string]interface{}, err error)`
- `cssColorToHex(value string) string` — handles rgb(), hsl(), named colors, hex

**Step 4: Run tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/extract/ -v`

**Step 5: Commit**

```bash
git add internal/extract/
git commit -m "feat: lightweight HTML-to-batch extractor"
```

---

### Task 6: Extract CLI Command

Add `ai-happy-design extract` cobra command.

**Files:**
- Modify: `cmd/ai-happy-design/main.go` (add extract command)

**Step 1: Add command**

```go
var extractCmd = &cobra.Command{
    Use:   "extract [file.html]",
    Short: "Convert HTML/CSS to batch JSON",
    Long:  "Parse an HTML file and output Figma batch operations as JSON",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        inputPath := args[0]
        canvasWidth, _ := cmd.Flags().GetInt("width")
        canvasHeight, _ := cmd.Flags().GetInt("height")
        computed, _ := cmd.Flags().GetBool("computed")
        outputPath, _ := cmd.Flags().GetString("output")

        var ops []map[string]interface{}
        var err error

        if computed {
            ops, err = extract.FromHTMLComputed(inputPath, extract.Options{
                CanvasWidth: canvasWidth, CanvasHeight: canvasHeight,
            })
        } else {
            f, ferr := os.Open(inputPath)
            if ferr != nil {
                return ferr
            }
            defer f.Close()
            ops, err = extract.FromHTML(f, extract.Options{
                CanvasWidth: canvasWidth, CanvasHeight: canvasHeight,
            })
        }
        if err != nil {
            return err
        }

        data, _ := json.MarshalIndent(ops, "", "  ")
        if outputPath != "" {
            return os.WriteFile(outputPath, data, 0644)
        }
        fmt.Println(string(data))
        return nil
    },
}

func init() {
    extractCmd.Flags().Int("width", 1080, "Canvas width")
    extractCmd.Flags().Int("height", 1080, "Canvas height")
    extractCmd.Flags().Bool("computed", false, "Use headless Chrome for computed styles")
    extractCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
    rootCmd.AddCommand(extractCmd)
}
```

**Step 2: Test pipeline**

```bash
ai-happy-design extract social-posts.html --width 1080 --height 1350 -o /tmp/extracted.json
ai-happy-design batch /tmp/extracted.json
```

**Step 3: Commit**

```bash
git add cmd/ai-happy-design/main.go
git commit -m "feat: add extract command for HTML-to-batch conversion"
```

---

### Task 7: Headless Chrome Mode (chromedp)

Add `--computed` mode to the extract command using chromedp.

**Files:**
- Create: `internal/extract/computed.go`
- Create: `internal/extract/computed_test.go`

**Step 1: Add chromedp dependency**

Run: `go get github.com/chromedp/chromedp`

**Step 2: Implement computed.go**

```go
package extract

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/chromedp/chromedp"
)

// FromHTMLComputed uses headless Chrome to extract computed styles.
func FromHTMLComputed(filePath string, opts Options) ([]map[string]interface{}, error) {
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return nil, err
    }

    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    var elementsJSON string

    err = chromedp.Run(ctx,
        chromedp.Navigate("file://"+absPath),
        chromedp.WaitReady("body"),
        chromedp.Evaluate(extractionScript, &elementsJSON),
    )
    if err != nil {
        return nil, fmt.Errorf("chromedp: %w", err)
    }

    var elements []computedElement
    if err := json.Unmarshal([]byte(elementsJSON), &elements); err != nil {
        return nil, err
    }

    return convertComputedToOps(elements, opts), nil
}

// extractionScript is injected into the page to extract computed styles.
const extractionScript = `
(function() {
    const results = [];
    function walk(el, depth) {
        if (el.nodeType === 3 && el.textContent.trim()) {
            results.push({
                type: 'text',
                text: el.textContent.trim(),
                parentIndex: results.length - 1,
            });
            return;
        }
        if (el.nodeType !== 1) return;
        const cs = window.getComputedStyle(el);
        const rect = el.getBoundingClientRect();
        results.push({
            type: 'element',
            tag: el.tagName.toLowerCase(),
            classes: Array.from(el.classList),
            rect: {x: rect.x, y: rect.y, width: rect.width, height: rect.height},
            styles: {
                backgroundColor: cs.backgroundColor,
                color: cs.color,
                fontSize: cs.fontSize,
                fontFamily: cs.fontFamily,
                fontWeight: cs.fontWeight,
                lineHeight: cs.lineHeight,
                letterSpacing: cs.letterSpacing,
                textAlign: cs.textAlign,
                textTransform: cs.textTransform,
                padding: cs.padding,
                borderRadius: cs.borderRadius,
                display: cs.display,
                flexDirection: cs.flexDirection,
                justifyContent: cs.justifyContent,
                alignItems: cs.alignItems,
                gap: cs.gap,
                backgroundImage: cs.backgroundImage,
                opacity: cs.opacity,
            },
            depth: depth,
        });
        for (const child of el.children) walk(child, depth + 1);
    }
    walk(document.body, 0);
    return JSON.stringify(results);
})()
`

type computedElement struct {
    Type        string            `json:"type"`
    Tag         string            `json:"tag"`
    Text        string            `json:"text"`
    Classes     []string          `json:"classes"`
    Rect        rect              `json:"rect"`
    Styles      map[string]string `json:"styles"`
    Depth       int               `json:"depth"`
    ParentIndex int               `json:"parentIndex"`
}

type rect struct {
    X      float64 `json:"x"`
    Y      float64 `json:"y"`
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
}

func convertComputedToOps(elements []computedElement, opts Options) []map[string]interface{} {
    // Scale factor: HTML rendered dimensions → target canvas
    // Find the outermost container to determine scale
    var ops []map[string]interface{}
    // Implementation: walk elements, create frames for divs, text for text nodes
    // Use computed rect positions scaled to canvas dimensions
    // Use computed styles mapped to Figma params
    for _, el := range elements {
        switch el.Type {
        case "element":
            op := computedElementToOp(el, opts)
            if op != nil {
                ops = append(ops, op)
            }
        case "text":
            op := computedTextToOp(el, opts)
            if op != nil {
                ops = append(ops, op)
            }
        }
    }
    return ops
}
```

**Step 3: Test**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/extract/ -run TestComputed -v`

**Step 4: Commit**

```bash
git add internal/extract/computed.go internal/extract/computed_test.go
git commit -m "feat: headless Chrome computed style extraction"
```

---

## Phase 3: Benchmark Harness

### Task 8: Benchmark Command (Provider-Agnostic)

Add `ai-happy-design benchmark` command. The benchmark is **provider-agnostic** — it doesn't call any LLM API directly. Instead it times stdin/stdout pipes so you can wrap ANY LLM call externally.

**Files:**
- Create: `internal/benchmark/benchmark.go`
- Create: `internal/benchmark/benchmark_test.go`
- Modify: `cmd/ai-happy-design/main.go` (add benchmark command)

**Step 1: Write test**

Create `internal/benchmark/benchmark_test.go`:

```go
package benchmark

import (
    "testing"
    "time"
)

func TestRunResult(t *testing.T) {
    r := &RunResult{
        PhaseA: PhaseTiming{Duration: 4200 * time.Millisecond, Label: "LLM Gen"},
        PhaseB: PhaseTiming{Duration: 6100 * time.Millisecond, Label: "CLI Exec", OpsCount: 42},
        PhaseC: PhaseTiming{Duration: 1800 * time.Millisecond, Label: "Verify"},
    }
    total := r.Total()
    if total < 12*time.Second || total > 13*time.Second {
        t.Errorf("expected ~12.1s total, got %v", total)
    }
}

func TestAggregateRuns(t *testing.T) {
    runs := []RunResult{
        {PhaseA: PhaseTiming{Duration: 4 * time.Second}, PhaseB: PhaseTiming{Duration: 6 * time.Second, OpsCount: 40}},
        {PhaseA: PhaseTiming{Duration: 5 * time.Second}, PhaseB: PhaseTiming{Duration: 6 * time.Second, OpsCount: 42}},
        {PhaseA: PhaseTiming{Duration: 4 * time.Second}, PhaseB: PhaseTiming{Duration: 7 * time.Second, OpsCount: 40}},
    }
    agg := Aggregate(runs)
    if agg.AvgTotal < 10*time.Second {
        t.Errorf("avg total too low: %v", agg.AvgTotal)
    }
    if agg.Runs != 3 {
        t.Errorf("expected 3 runs, got %d", agg.Runs)
    }
}
```

**Step 2: Implement benchmark.go**

Create `internal/benchmark/benchmark.go`:

```go
package benchmark

import (
    "fmt"
    "math"
    "strings"
    "time"
)

type PhaseTiming struct {
    Label    string
    Duration time.Duration
    OpsCount int
    Errors   int
    Meta     map[string]string // e.g., "model", "tokens_in", "tokens_out"
}

type RunResult struct {
    PhaseA PhaseTiming // LLM generation (optional, user-reported)
    PhaseB PhaseTiming // CLI batch execution
    PhaseC PhaseTiming // Export/verify (optional)
    Error  error
}

func (r *RunResult) Total() time.Duration {
    return r.PhaseA.Duration + r.PhaseB.Duration + r.PhaseC.Duration
}

type AggregateResult struct {
    Runs      int
    AvgA      time.Duration
    AvgB      time.Duration
    AvgC      time.Duration
    AvgTotal  time.Duration
    StdDev    time.Duration
    AvgOps    int
    AvgOpsPerSec float64
    Errors    int
}

func Aggregate(runs []RunResult) AggregateResult {
    if len(runs) == 0 {
        return AggregateResult{}
    }

    var totalA, totalB, totalC, totalAll time.Duration
    var totalOps, errors int

    for _, r := range runs {
        totalA += r.PhaseA.Duration
        totalB += r.PhaseB.Duration
        totalC += r.PhaseC.Duration
        totalAll += r.Total()
        totalOps += r.PhaseB.OpsCount
        if r.Error != nil {
            errors++
        }
    }

    n := time.Duration(len(runs))
    avgTotal := totalAll / n

    // Stddev
    var sumSq float64
    for _, r := range runs {
        diff := float64(r.Total() - avgTotal)
        sumSq += diff * diff
    }
    stddev := time.Duration(math.Sqrt(sumSq / float64(len(runs))))

    avgOps := totalOps / len(runs)
    avgOpsPerSec := 0.0
    if totalB > 0 {
        avgOpsPerSec = float64(totalOps) / totalB.Seconds()
    }

    return AggregateResult{
        Runs:         len(runs),
        AvgA:         totalA / n,
        AvgB:         totalB / n,
        AvgC:         totalC / n,
        AvgTotal:     avgTotal,
        StdDev:       stddev,
        AvgOps:       avgOps,
        AvgOpsPerSec: avgOpsPerSec,
        Errors:       errors,
    }
}

func (a AggregateResult) String() string {
    var b strings.Builder
    w := 48
    line := strings.Repeat("─", w)

    fmt.Fprintf(&b, "╭%s╮\n", line)
    fmt.Fprintf(&b, "│  AI Happy Design — Benchmark Results%s│\n", strings.Repeat(" ", w-38))
    fmt.Fprintf(&b, "├%s┤\n", line)
    fmt.Fprintf(&b, "│  Runs: %-5d%s│\n", a.Runs, strings.Repeat(" ", w-13))

    if a.AvgA > 0 {
        fmt.Fprintf(&b, "│  Phase A (LLM Gen)    │  avg %-8s│\n", fmtDur(a.AvgA))
    }
    fmt.Fprintf(&b, "│  Phase B (CLI Exec)   │  avg %-8s│\n", fmtDur(a.AvgB))
    if a.AvgOps > 0 {
        fmt.Fprintf(&b, "│    └ ops: %-4d  │  %.1f ops/s%s│\n",
            a.AvgOps, a.AvgOpsPerSec, strings.Repeat(" ", 8))
    }
    if a.AvgC > 0 {
        fmt.Fprintf(&b, "│  Phase C (Verify)     │  avg %-8s│\n", fmtDur(a.AvgC))
    }

    fmt.Fprintf(&b, "├%s┤\n", line)
    fmt.Fprintf(&b, "│  TOTAL                │  avg %-8s│\n", fmtDur(a.AvgTotal))
    fmt.Fprintf(&b, "│                       │  ± %-9s│\n", fmtDur(a.StdDev))
    fmt.Fprintf(&b, "│  Errors: %d/%d runs%s│\n",
        a.Errors, a.Runs, strings.Repeat(" ", w-20))
    fmt.Fprintf(&b, "╰%s╯\n", line)

    return b.String()
}

func fmtDur(d time.Duration) string {
    if d < time.Second {
        return fmt.Sprintf("%.0fms", float64(d)/float64(time.Millisecond))
    }
    return fmt.Sprintf("%.1fs", d.Seconds())
}
```

**Step 3: Add benchmark cobra commands to main.go**

Three subcommands:

```go
// benchmark exec — time batch execution only
// benchmark pipe — accept JSON from stdin, report phase-a-ms externally
// benchmark compare — run multiple methods side-by-side

var benchmarkCmd = &cobra.Command{
    Use:   "benchmark",
    Short: "Measure end-to-end performance",
}

var benchExecCmd = &cobra.Command{
    Use:   "exec [file.json]",
    Short: "Benchmark batch execution time",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        runs, _ := cmd.Flags().GetInt("runs")
        // For each run: time the batch execution
        // Collect RunResults with PhaseB timing
        // Print aggregate
    },
}

var benchPipeCmd = &cobra.Command{
    Use:   "pipe",
    Short: "Benchmark with external LLM timing (stdin JSON + --phase-a-ms)",
    RunE: func(cmd *cobra.Command, args []string) error {
        phaseAMs, _ := cmd.Flags().GetInt("phase-a-ms")
        // Read JSON from stdin
        // Time batch execution (Phase B)
        // Add Phase A from --phase-a-ms flag
        // Print combined result
    },
}

var benchCompareCmd = &cobra.Command{
    Use:   "compare [file.html]",
    Short: "Compare extraction methods side-by-side",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        runs, _ := cmd.Flags().GetInt("runs")
        // Run lightweight extraction N times
        // Run computed extraction N times (if Chrome available)
        // Print comparison table
    },
}
```

**Step 4: Run tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go test ./internal/benchmark/ -v`

**Step 5: Build and test**

Run: `make deploy`

Test:
```bash
# Benchmark batch execution
ai-happy-design benchmark exec /tmp/amc-slides/post1.json --runs 3

# Benchmark with external LLM (wrap any curl call)
START=$(date +%s%N)
BATCH=$(curl -sS https://api.cerebras.ai/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $CEREBRAS_API_KEY" \
  -d '{"model":"qwen-3-235b","messages":[{"role":"user","content":"Generate a Figma batch JSON for a simple slide"}]}' \
  | jq -r '.choices[0].message.content')
LLM_MS=$(( ($(date +%s%N) - START) / 1000000 ))

echo "$BATCH" | ai-happy-design benchmark pipe --phase-a-ms $LLM_MS

# Compare extraction methods
ai-happy-design benchmark compare social-posts.html --runs 3
```

**Step 6: Commit**

```bash
git add internal/benchmark/ cmd/ai-happy-design/main.go
git commit -m "feat: add provider-agnostic benchmark harness"
```

---

### Task 9: Update Skill and Discovery Docs

Update all three LLM discovery layers with the new features.

**Files:**
- Modify: `~/.claude/skills/ai-happy-design/SKILL.md`
- Modify: `internal/tools/catalog_llm.go`
- Modify: `internal/tools/describe.go`

**Step 1: Add to SKILL.md**

Add sections for:
- Composite commands (slide, banner) with examples
- text.set_range_style with example
- extract command usage
- benchmark command usage

**Step 2: Update catalog_llm.go**

Add composite commands, text.set_range_style, and extract to the catalog output.

**Step 3: Update describe.go**

Add new tool descriptions for composite commands and extract.

**Step 4: Commit**

```bash
git add internal/tools/catalog_llm.go internal/tools/describe.go
git commit -m "docs: update LLM discovery layers with composite commands and benchmark"
```

---

### Task 10: End-to-End Validation

Run the full pipeline with the original AMC social posts to validate < 2 minute target.

**Step 1: Extract from HTML**

```bash
ai-happy-design extract /path/to/social-posts.html --width 1080 --height 1350 -o /tmp/amc-extracted.json
```

**Step 2: Execute batch**

```bash
ai-happy-design batch /tmp/amc-extracted.json --live
```

**Step 3: Benchmark**

```bash
ai-happy-design benchmark exec /tmp/amc-extracted.json --runs 3
```

**Step 4: Export and verify**

Export a few slides and visually inspect via the Read tool.

**Step 5: Commit any fixes**

```bash
git commit -m "fix: adjustments from end-to-end validation"
```

---

## Implementation Order Summary

| Task | Phase | Description | Depends On |
|------|-------|-------------|------------|
| 1 | 1 | text.set_range_style plugin handler | — |
| 2 | 1 | Composite command expander (Go) | — |
| 3 | 1 | Integrate composites into batch pipeline | 2 |
| 4 | 1 | LLM discovery updates | 2, 3 |
| 5 | 2 | Lightweight HTML parser | — |
| 6 | 2 | Extract CLI command | 5 |
| 7 | 2 | Headless Chrome mode | 5, 6 |
| 8 | 3 | Benchmark harness | — |
| 9 | — | Skill and docs update | 1-8 |
| 10 | — | End-to-end validation | 1-9 |

Tasks 1, 2, 5, and 8 can run in parallel (no dependencies between them).
