# Learnings

## [LRN-20260217-001] best_practice

**Logged**: 2026-02-17T17:30:00Z
**Priority**: critical
**Status**: resolved
**Area**: docs

### Summary
New features must be surfaced in ALL THREE discoverability layers or LLMs won't find them.

### Details
Implemented 3 new features (design_system.analyze, compact tree, export to temp dir) but none were discoverable by LLMs. The tool code worked fine, but:
1. `describe.go` (toolDescriptions) — export.image didn't mention auto-save behavior
2. `catalog_llm.go` (LLM catalog) — no workflow hints for design_system or compact tree
3. `SKILL.md` (skill file) — none of the 3 features were mentioned

A model calling `tools --llm --json` or reading the skill would have zero awareness of these features.

### Suggested Action
Checklist for every new feature:
1. Add/update entry in `describe.go` toolDescriptions (with Params: for buildActionSpec)
2. Add workflow hint in `catalog_llm.go` if the feature has a "when to use" pattern
3. Update `SKILL.md` workflow, key commands, and pitfalls sections

### Resolution
- **Resolved**: 2026-02-17T17:30:00Z
- **Commit**: a92ef08
- **Notes**: Updated all three layers. Export description mentions temp dir. Catalog has 3 new workflow hints. Skill has updated workflow, commands, and pitfalls.

### Metadata
- Source: conversation (LLM discoverability audit)
- Related Files: internal/tools/describe.go, internal/tools/catalog_llm.go, SKILL.md
- Tags: discoverability, llm, documentation, skill

---

## [LRN-20260217-002] knowledge_gap

**Logged**: 2026-02-17T17:30:00Z
**Priority**: high
**Status**: promoted
**Area**: docs

### Summary
Skill files should target 2,500-4,000 tokens for cross-model compatibility. Adding 1-5 usage examples improves accuracy 72% to 90%.

### Details
Research from Anthropic's Advanced Tool Use guide and OpenAI's GPT-4.1 guide:
- Small models (Haiku, 4o-mini): 2,000-4,000 tokens, max 15 tools
- Medium models (Sonnet, 4o): 3,000-6,000 tokens, 20+ tools OK
- Large models (Opus, o3): 5,000-10,000 tokens, 50+ tools OK
- System prompt + tools should be 5-10% of context window
- Three-reminder pattern boosts SWE-bench by ~20% (OpenAI)
- Tool descriptions: 1-2 sentences, front-load important info, state when NOT to use
- Weaker models: limit to 10-15 tools, use enums, add "call X before Y" sequencing

Our skill was ~3,200 tokens (sweet spot). Trimmed creative section (~400 tokens) and added new features (~100 tokens) = nearly token-neutral.

### Suggested Action
When creating/updating skills, check token count and keep within 2,500-4,000 range. Trim verbose sections to reference files.

### Metadata
- Source: web_research (Anthropic tool use guide, OpenAI GPT-4.1 guide, Merge MCP guide)
- Related Files: SKILL.md
- Tags: tokens, skill-design, cross-model, steering
- Promoted: MEMORY.md

---

## [LRN-20260217-003] best_practice

**Logged**: 2026-02-17T17:30:00Z
**Priority**: high
**Status**: resolved
**Area**: backend

### Summary
Use os.TempDir() in Go for OS-agnostic temp paths, never hardcode /tmp.

### Details
Export feature initially hardcoded `/tmp/ahd-export-*` paths. This works on macOS/Linux but would fail on Windows (which uses `C:\Users\...\AppData\Local\Temp`). Go's `os.TempDir()` returns the correct platform temp directory.

Applied to both `internal/tools/export.go` (MCP) and `cmd/ai-happy-design/main.go` (CLI).

### Suggested Action
Always use `os.TempDir()` for temporary file paths. Grep for hardcoded `/tmp/` in Go code.

### Resolution
- **Resolved**: 2026-02-17T17:30:00Z
- **Commit**: a92ef08
- **Notes**: Replaced all hardcoded /tmp/ with os.TempDir() in export paths.

### Metadata
- Source: user_feedback
- Related Files: internal/tools/export.go, cmd/ai-happy-design/main.go
- Tags: cross-platform, temp-files, go

---

## [LRN-20260217-004] best_practice

**Logged**: 2026-02-17T17:30:00Z
**Priority**: medium
**Status**: resolved
**Area**: frontend

### Summary
Figma's QuickJS runtime doesn't support spread syntax in push — use loops instead.

### Details
`result.push(...serializeNodeCompact(...))` fails silently or throws in Figma's plugin sandbox (QuickJS engine). The spread operator in function call context is not supported. Must use explicit loop:
```typescript
var items = serializeNodeCompact(child, maxDepth, node.id, currentDepth + 1);
for (var i = 0; i < items.length; i++) {
  result.push(items[i]);
}
```

This was caught during compact tree implementation. The esbuild `es2015` target doesn't downlevel spread-in-call because it's technically ES2015, but QuickJS doesn't implement it.

### Suggested Action
Never use `array.push(...otherArray)` in Figma plugin code. Use loops. Add this to plugin coding guidelines.

### Resolution
- **Resolved**: 2026-02-17T17:30:00Z
- **Commit**: 75087e8
- **Notes**: Fixed in serializeCompact.ts. All existing handlers already use loops.

### Metadata
- Source: error (runtime failure in Figma plugin)
- Related Files: plugin/src/utils/serializeCompact.ts
- Tags: quickjs, figma-plugin, spread-syntax, compatibility
- See Also: esbuild $& replacement bug (MEMORY.md)

---

## [LRN-20260217-005] best_practice

**Logged**: 2026-02-17T17:30:00Z
**Priority**: medium
**Status**: resolved
**Area**: backend

### Summary
Flat array serialization is 3-5x more token-efficient than nested tree for structural discovery.

### Details
When exploring Figma document structure, the full nested tree with all properties consumes excessive tokens. A compact flat array format:
```json
[{"id":"1:2","type":"FRAME","name":"Card","x":0,"y":0,"w":400,"h":300,"childCount":3,"parentId":"0:1","depth":0}]
```
is 3-5x more token-efficient than the full recursive tree while still providing all structural information needed for discovery (IDs, positions, sizes, hierarchy).

This pattern applies broadly: when LLMs only need to discover what exists before making targeted edits, compact summaries save significant tokens.

### Suggested Action
For any API response intended for LLM consumption, offer a compact/summary mode that strips properties not needed for the specific use case.

### Resolution
- **Resolved**: 2026-02-17T17:30:00Z
- **Commit**: 91fa3dc
- **Notes**: Added compact:true param to node.get_tree. Returns flat array instead of nested tree.

### Metadata
- Source: conversation (Figma MCP comparison research)
- Related Files: plugin/src/utils/serializeCompact.ts, internal/tools/node.go
- Tags: token-efficiency, llm-optimization, serialization

---
