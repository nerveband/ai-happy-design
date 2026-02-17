# Figma MCP-Inspired Features Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add three features inspired by Figma's official MCP: export-to-temp-file, compact tree output, and a design system analyzer tool.

**Architecture:** Feature 1 is Go-only (CLI already saves exports, just change default path to /tmp). Feature 2 adds a `compact` param to node.get_tree with a new flat serializer in the plugin. Feature 3 is a new Go tool that orchestrates three existing plugin commands.

**Tech Stack:** Go (MCP server + CLI), TypeScript (Figma plugin)

---

### Task 1: Export to Temp File — Plugin test for compact serializer

**Files:**
- Create: `plugin/src/utils/serializeCompact.ts`
- Create: `plugin/src/utils/serializeCompact.test.ts`

**Step 1: Write the compact serializer**

```typescript
// plugin/src/utils/serializeCompact.ts

export interface CompactNode {
  id: string;
  type: string;
  name: string;
  x: number;
  y: number;
  w: number;
  h: number;
  childCount: number;
  parentId: string | null;
  depth: number;
}

export function serializeNodeCompact(
  node: SceneNode | PageNode,
  maxDepth: number = 3,
  parentId: string | null = null,
  currentDepth: number = 0
): CompactNode[] {
  const result: CompactNode[] = [];

  const entry: CompactNode = {
    id: node.id,
    type: node.type,
    name: node.name,
    x: 'x' in node ? (node as any).x : 0,
    y: 'y' in node ? (node as any).y : 0,
    w: 'width' in node ? (node as any).width : 0,
    h: 'height' in node ? (node as any).height : 0,
    childCount: 'children' in node ? (node as any).children.length : 0,
    parentId: parentId,
    depth: currentDepth,
  };
  result.push(entry);

  if ('children' in node && currentDepth < maxDepth) {
    for (const child of (node as any).children as SceneNode[]) {
      result.push(...serializeNodeCompact(child, maxDepth, node.id, currentDepth + 1));
    }
  }

  return result;
}
```

**Step 2: Write the test**

```typescript
// plugin/src/utils/serializeCompact.test.ts
import { describe, it, expect } from 'vitest';
import { serializeNodeCompact, CompactNode } from './serializeCompact';

function mockNode(overrides: any = {}): any {
  return {
    id: '1:1',
    type: 'FRAME',
    name: 'TestFrame',
    x: 0,
    y: 0,
    width: 100,
    height: 200,
    children: [],
    ...overrides,
  };
}

describe('serializeNodeCompact', () => {
  it('serializes a leaf node', () => {
    const node = mockNode({ children: undefined, type: 'RECTANGLE', id: '1:2', name: 'Rect' });
    delete node.children;
    const result = serializeNodeCompact(node);
    expect(result).toHaveLength(1);
    expect(result[0]).toEqual({
      id: '1:2', type: 'RECTANGLE', name: 'Rect',
      x: 0, y: 0, w: 100, h: 200,
      childCount: 0, parentId: null, depth: 0,
    });
  });

  it('serializes parent with children', () => {
    const child1 = mockNode({ id: '1:3', name: 'Child1', type: 'TEXT', x: 10, y: 20, width: 80, height: 30 });
    delete child1.children;
    const child2 = mockNode({ id: '1:4', name: 'Child2', type: 'RECTANGLE', x: 10, y: 60, width: 80, height: 40 });
    delete child2.children;
    const parent = mockNode({ id: '1:2', name: 'Parent', children: [child1, child2] });

    const result = serializeNodeCompact(parent);
    expect(result).toHaveLength(3);
    expect(result[0].id).toBe('1:2');
    expect(result[0].childCount).toBe(2);
    expect(result[0].parentId).toBeNull();
    expect(result[1].parentId).toBe('1:2');
    expect(result[2].parentId).toBe('1:2');
  });

  it('respects maxDepth', () => {
    const grandchild = mockNode({ id: '1:5', name: 'GC', type: 'TEXT' });
    delete grandchild.children;
    const child = mockNode({ id: '1:3', name: 'Child', children: [grandchild] });
    const parent = mockNode({ id: '1:2', name: 'Parent', children: [child] });

    const result = serializeNodeCompact(parent, 1);
    expect(result).toHaveLength(2); // parent + child, no grandchild
    expect(result[1].childCount).toBe(1); // child reports it has children even though not traversed
  });

  it('returns flat array with correct depth values', () => {
    const grandchild = mockNode({ id: '1:5', name: 'GC', type: 'TEXT' });
    delete grandchild.children;
    const child = mockNode({ id: '1:3', name: 'Child', children: [grandchild] });
    const parent = mockNode({ id: '1:2', name: 'Parent', children: [child] });

    const result = serializeNodeCompact(parent, 3);
    expect(result.map(n => n.depth)).toEqual([0, 1, 2]);
  });
});
```

**Step 3: Run test to verify it passes**

Run: `cd /tmp && cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run src/utils/serializeCompact.test.ts`
Expected: PASS (4 tests)

**Step 4: Commit**

```bash
git add plugin/src/utils/serializeCompact.ts plugin/src/utils/serializeCompact.test.ts
git commit -m "feat: add compact node serializer for token-efficient tree scanning"
```

---

### Task 2: Wire compact mode into plugin node handler

**Files:**
- Modify: `plugin/src/handlers/node.ts:295-305` (getTree function)

**Step 1: Import and use compact serializer in getTree**

Add import at top of `plugin/src/handlers/node.ts`:
```typescript
import { serializeNodeCompact } from '../utils/serializeCompact';
```

Modify the `getTree` function (lines 295-305):
```typescript
async function getTree(params: any) {
  const { nodeId, depth, compact } = params;
  const node = await getNodeById(nodeId);

  // If it's a page node, load it first so children are accessible
  if (node.type === 'PAGE') {
    await (node as PageNode).loadAsync();
  }

  if (compact) {
    return serializeNodeCompact(node as SceneNode, depth ?? 3);
  }

  return serializeNode(node as SceneNode, 0, depth ?? 3);
}
```

**Step 2: Add compact param to Go MCP tool definition**

Modify `internal/tools/node.go` — add after line 36 (depth param):
```go
mcp.WithBoolean("compact", mcp.Description("Return flat summary array instead of full tree. Each entry: {id, type, name, x, y, w, h, childCount, parentId, depth}. Much more token-efficient for structural discovery.")),
```

Modify the `get_tree` case (lines 72-80) to pass compact through:
```go
case "get_tree":
    nodeId, errResult := requireStringArg(args, "nodeId")
    if errResult != nil {
        return errResult, nil
    }
    return sendCommand(commander, "get_node_tree", map[string]interface{}{
        "nodeId":  nodeId,
        "depth":   getFloat64Arg(args, "depth", 3),
        "compact": getBoolArg(args, "compact", false),
    })
```

**Step 3: Update tool descriptions in describe.go**

Find the `"node"` section in `toolDescriptions` and update the `get_tree` entry to mention compact mode.

**Step 4: Build and test manually**

Run: `cd /tmp && "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/ai-happy-design" command node.get_tree '{"nodeId":"0:1","depth":1,"compact":true}'`
Expected: Flat JSON array with `{id, type, name, x, y, w, h, childCount, parentId, depth}` entries

**Step 5: Commit**

```bash
git add plugin/src/handlers/node.ts internal/tools/node.go internal/tools/describe.go
git commit -m "feat: add compact mode to node.get_tree for token-efficient scanning"
```

---

### Task 3: Export to temp file — Change CLI default save path

**Files:**
- Modify: `cmd/ai-happy-design/main.go:316-342` (export auto-save logic)

**Step 1: Change auto-generated path to use /tmp**

In `cmd/ai-happy-design/main.go`, the current auto-save logic (lines 316-342) saves to the current directory with just `name + ext`. Change it to save to `/tmp/ahd-export-<name>-<timestamp><ext>`.

Replace lines 316-342 with:
```go
			outPath := commandOutput
			if outPath == "" {
				// Auto-generate filename in /tmp for LLM-friendly access
				name, _ := parsed["name"].(string)
				if name == "" {
					name = "export"
				}
				// Sanitize name for filename
				name = strings.Map(func(r rune) rune {
					if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' || r == ' ' {
						return '-'
					}
					return r
				}, name)
				ext := ".png"
				switch strings.ToUpper(format) {
				case "JPG":
					ext = ".jpg"
				case "SVG":
					ext = ".svg"
				case "PDF":
					ext = ".pdf"
				case "JSON":
					ext = ".json"
				}
				outPath = fmt.Sprintf("/tmp/ahd-export-%s-%d%s", name, time.Now().Unix(), ext)
			}
```

Also ensure `"time"` is in the imports at the top of the file.

**Step 2: Update MCP export to include path in metadata**

In `internal/tools/export.go`, modify `sendExportCommand` (lines 57-123) to also write a temp file and include the path:

After line 82 (after the empty data check), add temp file write:
```go
	// Write to temp file for CLI/LLM access
	var tempPath string
	{
		nodeIDSafe := strings.ReplaceAll(nodeID, ":", "-")
		ext := ".png"
		switch format {
		case "JPG":
			ext = ".jpg"
		case "PDF":
			ext = ".pdf"
		}
		tempPath = fmt.Sprintf("/tmp/ahd-export-%s-%d%s", nodeIDSafe, time.Now().Unix(), ext)

		decoded, decErr := base64.StdEncoding.DecodeString(base64Data)
		if decErr == nil {
			_ = os.WriteFile(tempPath, decoded, 0644)
		}
	}
```

Then include `tempPath` in the metadata map (after line 96):
```go
	meta["path"] = tempPath
```

Add `"encoding/base64"`, `"os"`, `"strings"`, and `"time"` to imports.

**Step 3: Test CLI export**

Run: `cd /tmp && ai-happy-design command export.image '{"nodeId":"0:1"}'`
Expected: JSON output with `"savedTo": "/tmp/ahd-export-..."` pointing to a temp file

**Step 4: Commit**

```bash
git add cmd/ai-happy-design/main.go internal/tools/export.go
git commit -m "feat: export saves to /tmp by default for LLM-friendly file access"
```

---

### Task 4: Design System Analyzer — Create the Go tool

**Files:**
- Create: `internal/tools/design_system.go`

**Step 1: Write the tool**

```go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterDesignSystemTool registers the "design_system" tool for analyzing
// and generating design rules from the current Figma file.
func RegisterDesignSystemTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("design_system",
		mcp.WithDescription("Analyze the current Figma file's styles, variables, and components to generate design consistency rules. Use before creating designs in existing files."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("analyze")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		action := getStringArg(req.GetArguments(), "action", "")

		switch action {
		case "analyze":
			return analyzeDesignSystem(commander)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown design_system action: %s", action)), nil
		}
	})
}

func analyzeDesignSystem(commander *figma.Commander) (*mcp.CallToolResult, error) {
	// 1. Fetch styles
	stylesRaw, err := commander.SendCommand("get_all_styles", map[string]interface{}{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get styles: %s", err)), nil
	}

	// 2. Fetch variables
	varsRaw, err := commander.SendCommand("get_all_variables", map[string]interface{}{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get variables: %s", err)), nil
	}

	// 3. Fetch components
	compsRaw, err := commander.SendCommand("get_local_components", map[string]interface{}{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get components: %s", err)), nil
	}

	// Parse responses
	styles := toMap(stylesRaw)
	vars := toMap(varsRaw)
	comps := toMap(compsRaw)

	// Build categorized output
	result := buildDesignRules(styles, vars, comps)

	out, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(out)), nil
}

func buildDesignRules(styles, vars, comps map[string]interface{}) map[string]interface{} {
	rules := map[string]interface{}{}

	// --- Colors ---
	colorSection := map[string]interface{}{}
	paintStyles := extractArray(styles, "paintStyles")
	if len(paintStyles) > 0 {
		colorSection["styles"] = paintStyles
	}

	colorVars := filterVariablesByType(vars, "COLOR")
	if len(colorVars) > 0 {
		colorSection["variables"] = colorVars
	}

	if len(paintStyles) > 0 || len(colorVars) > 0 {
		colorSection["rule"] = "Use these existing colors. Do not introduce new hex values — apply paint styles by ID or reference color variables."
		rules["colors"] = colorSection
	}

	// --- Typography ---
	typoSection := map[string]interface{}{}
	textStyles := extractArray(styles, "textStyles")
	if len(textStyles) > 0 {
		typoSection["styles"] = textStyles
		typoSection["rule"] = "Apply text styles by ID when available. Match font families and weights from existing styles."
		rules["typography"] = typoSection
	}

	// --- Spacing ---
	spacingVars := filterVariablesByType(vars, "FLOAT")
	if len(spacingVars) > 0 {
		rules["spacing"] = map[string]interface{}{
			"variables": spacingVars,
			"rule":      "Use spacing variables for padding, gaps, and margins. Snap to the existing scale.",
		}
	}

	// --- Effects ---
	effectStyles := extractArray(styles, "effectStyles")
	if len(effectStyles) > 0 {
		rules["effects"] = map[string]interface{}{
			"styles": effectStyles,
			"rule":   "Apply effect styles by ID for shadows, blurs, and other effects instead of creating new ones.",
		}
	}

	// --- Components ---
	components := extractArray(comps, "components")
	if len(components) > 0 {
		rules["components"] = map[string]interface{}{
			"available": components,
			"rule":      "Instantiate existing components (component.create_instance) instead of rebuilding from scratch.",
		}
	}

	// --- Summary ---
	summary := []string{}
	if n := len(paintStyles); n > 0 {
		summary = append(summary, fmt.Sprintf("%d paint styles", n))
	}
	if n := len(textStyles); n > 0 {
		summary = append(summary, fmt.Sprintf("%d text styles", n))
	}
	if n := len(effectStyles); n > 0 {
		summary = append(summary, fmt.Sprintf("%d effect styles", n))
	}
	if allVars := extractArray(vars, "variables"); len(allVars) > 0 {
		summary = append(summary, fmt.Sprintf("%d variables", len(allVars)))
	}
	if n := len(components); n > 0 {
		summary = append(summary, fmt.Sprintf("%d components", n))
	}

	if len(summary) > 0 {
		rules["summary"] = fmt.Sprintf("This file has %s. Follow the rules above for consistency.", strings.Join(summary, ", "))
	} else {
		rules["summary"] = "This file has no styles, variables, or components defined yet. You may create designs freely."
	}

	return rules
}

// toMap converts an interface{} to map[string]interface{}, returning empty map on failure.
func toMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	// Try via JSON round-trip for typed structs
	b, err := json.Marshal(v)
	if err != nil {
		return map[string]interface{}{}
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]interface{}{}
	}
	return m
}

// extractArray extracts a []interface{} from a map by key.
func extractArray(m map[string]interface{}, key string) []interface{} {
	if arr, ok := m[key].([]interface{}); ok {
		return arr
	}
	return nil
}

// filterVariablesByType returns variables matching the given resolvedType.
func filterVariablesByType(vars map[string]interface{}, typeName string) []interface{} {
	allVars := extractArray(vars, "variables")
	var result []interface{}
	for _, v := range allVars {
		if vm, ok := v.(map[string]interface{}); ok {
			if rt, _ := vm["resolvedType"].(string); rt == typeName {
				result = append(result, v)
			}
		}
	}
	return result
}
```

**Step 2: Register the tool**

In `internal/tools/registry.go`, add after line 28 (`RegisterDesignTool(s)`):
```go
	RegisterDesignSystemTool(s, commander)
```

**Step 3: Add to tool descriptions**

In `internal/tools/describe.go`, add to the `toolDescriptions` map:
```go
	"design_system": {
		"analyze": "Analyze the current Figma file's styles, variables, and components. Returns categorized rules for colors, typography, spacing, effects, and components. Use before creating designs in existing files to maintain consistency.",
	},
```

**Step 4: Add to CLI local action routing**

In `cmd/ai-happy-design/main.go`, find the `validActions` map (around line 654) and add:
```go
	"design_system": {"analyze"},
```

**Step 5: Add to LLM catalog**

In `internal/tools/catalog_llm.go`, add a `design_system` entry in the `workflow` section or `tools` map with guidance like:
```
"design_system": "Call design_system.analyze BEFORE creating designs in existing Figma files. It returns the file's existing styles, variables, and components with rules for reusing them. This ensures design consistency."
```

**Step 6: Commit**

```bash
git add internal/tools/design_system.go internal/tools/registry.go internal/tools/describe.go internal/tools/catalog_llm.go cmd/ai-happy-design/main.go
git commit -m "feat: add design_system.analyze tool for design consistency rules"
```

---

### Task 5: Build, deploy, and verify all three features

**Files:**
- None (build + test)

**Step 1: Build plugin**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npm run build`
Expected: Build succeeds

**Step 2: Run plugin tests**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin" && npx vitest run`
Expected: All tests pass (existing 24 + new 4 = 28)

**Step 3: Deploy**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && make deploy`
Expected: Plugin built, Go binary compiled, signed, installed, relay restarted

**Step 4: Verify compact tree**

Run: `cd /tmp && ai-happy-design command node.get_tree '{"nodeId":"0:1","depth":1,"compact":true}'`
Expected: Flat JSON array

**Step 5: Verify export to /tmp**

Run: `cd /tmp && ai-happy-design command export.image '{"nodeId":"0:1"}'`
Expected: JSON with `"savedTo": "/tmp/ahd-export-..."`, file exists

**Step 6: Verify design system analyzer**

Run: `cd /tmp && ai-happy-design command design_system.analyze`
Expected: JSON with colors/typography/spacing/components/summary sections

**Step 7: Final commit if any fixups needed**

```bash
git add -A && git commit -m "fix: post-integration adjustments for MCP-inspired features"
```
