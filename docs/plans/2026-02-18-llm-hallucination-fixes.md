# LLM Hallucination Fixes — Compound Aliases + SKILL.md Audit

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix the root causes that make LLMs invent wrong commands (like `document.list_pages` instead of `page.get_all`) by (a) making the wrong commands work anyway and (b) surfacing the right commands in the skill file.

**Architecture:** Two-layer fix. Go layer: add a `compoundAliases` map in `command_routing.go` that intercepts common LLM hallucinations (dotted commands like `document.list_pages`) before they reach the plugin. SKILL.md layer: add a Page Management section, Domain Quick Reference table, and expanded pitfalls so LLMs guess right in the first place.

**Tech Stack:** Go (command routing), Markdown (SKILL.md)

---

## Background: Why `document.list_pages` Failed

The routing code in `internal/ws/command_routing.go` has two paths:
1. If command contains `.` → split into `domain.action` directly (e.g., `document.list_pages` → domain=`document`, action=`list_pages`)
2. Otherwise → check `legacyCommandRoutes` map (e.g., `list_pages` → `page.get_all`)

The legacy map has `list_pages` as an alias for `page.get_all`. But the LLM wrote `document.list_pages`, which hits path #1, sending `list_pages` to the `document` handler — which doesn't know that action.

**Fix**: Add a `compoundAliases` map checked BEFORE the dot-split.

---

## Audit: What's Missing from SKILL.md

Cross-referencing `describe.go` vs `SKILL.md`:

| Domain | Gap | LLM Likely Guesses |
|--------|-----|-------------------|
| `page.*` | Entirely absent | `document.list_pages`, `document.get_pages` |
| `layer.*` | Entirely absent | `node.group`, `node.bring_to_front` |
| `design_system.analyze` | Not in Key CLI Commands | `document.analyze`, `document.get_design_system` |
| `node.get_tree compact:true` | Not in Key CLI Commands | N/A — just unknown |
| `document.*` scan/select | Not in Key CLI Commands | `document.scan_all` |

---

## Task 1: Add Compound Alias Routing in Go

**Files:**
- Modify: `internal/ws/command_routing.go`

No tests exist for command routing. Manual test after.

**Step 1: Add the `compoundAliases` map**

Add this block to `command_routing.go` after the `legacyCommandRoutes` var block (around line 203):

```go
// compoundAliases intercepts common LLM hallucinations for dotted commands
// (e.g. "document.list_pages") before they reach the wrong domain handler.
// These are checked BEFORE splitting on "." in resolveCommandRoute.
var compoundAliases = map[string]commandRoute{
	// Page operations — LLMs often guess document.* for these
	"document.list_pages":   {Domain: "page", Action: "get_all"},
	"document.get_pages":    {Domain: "page", Action: "get_all"},
	"document.pages":        {Domain: "page", Action: "get_all"},
	"document.create_page":  {Domain: "page", Action: "create"},
	"document.delete_page":  {Domain: "page", Action: "delete"},
	"document.rename_page":  {Domain: "page", Action: "rename"},
	"document.switch_page":  {Domain: "page", Action: "set_current"},
	"document.set_page":     {Domain: "page", Action: "set_current"},
	// Layer operations — LLMs often guess node.* for these
	"node.group":           {Domain: "layer", Action: "group"},
	"node.ungroup":         {Domain: "layer", Action: "ungroup"},
	"node.bring_to_front":  {Domain: "layer", Action: "bring_to_front"},
	"node.send_to_back":    {Domain: "layer", Action: "send_to_back"},
	"node.bring_forward":   {Domain: "layer", Action: "bring_forward"},
	"node.send_backward":   {Domain: "layer", Action: "send_backward"},
	"node.move_to_parent":  {Domain: "layer", Action: "insert_child"},
	"node.reorder":         {Domain: "layer", Action: "set_order"},
	// Design system — LLMs might try document.*
	"document.analyze":            {Domain: "design_system", Action: "analyze"},
	"document.get_design_system":  {Domain: "design_system", Action: "analyze"},
	"document.design_system":      {Domain: "design_system", Action: "analyze"},
}
```

**Step 2: Update `resolveCommandRoute` to check compound aliases first**

Current (line 205-207):
```go
func resolveCommandRoute(command string, params map[string]interface{}) (string, string, error) {
	if dot := strings.Index(command, "."); dot > 0 && dot < len(command)-1 {
		return command[:dot], command[dot+1:], nil
	}
```

Change to:
```go
func resolveCommandRoute(command string, params map[string]interface{}) (string, string, error) {
	// Check compound aliases first — catches LLM hallucinations like "document.list_pages"
	if route, ok := compoundAliases[command]; ok {
		return route.Domain, route.Action, nil
	}

	if dot := strings.Index(command, "."); dot > 0 && dot < len(command)-1 {
		return command[:dot], command[dot+1:], nil
	}
```

**Step 3: Build and verify it compiles**

```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design"
go build ./...
```
Expected: no errors.

**Step 4: Manual smoke test**

```bash
AHD_CHANNEL=golden-waffle-25 ai-happy-design command document.list_pages '{}'
```
Expected: returns JSON list of pages (same as `page.get_all`).

**Step 5: Commit**

```bash
git add internal/ws/command_routing.go
git commit -m "fix: add compound alias routing to catch LLM hallucinations

LLMs commonly call document.list_pages, node.group, etc. which hit the
wrong domain handler after the dot-split. Add compoundAliases map checked
before dot-split so these route correctly.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Update SKILL.md — Page Management + Domain Quick Reference

**Files:**
- Modify: `skills/ai-happy-design/SKILL.md`

**Step 1: Add a Page Management section**

Insert before the `## Relay & Port Configuration` section (around line 373):

```markdown
## Page Management

Pages are separate from document operations. Use `page.*` commands — NOT `document.*`.

```bash
ai-happy-design command page.get_all '{}'                          # list all pages
ai-happy-design command page.create '{"name": "Slide 2"}'         # new page
ai-happy-design command page.set_current '{"pageId": "3:0"}'      # switch to page
ai-happy-design command page.rename '{"pageId": "3:0", "name": "Final"}'
ai-happy-design command page.duplicate '{"pageId": "3:0"}'
ai-happy-design command page.delete '{"pageId": "3:0"}'
```

**Common mistake**: `document.list_pages` → use `page.get_all`
```

**Step 2: Add `design_system.analyze` and compact tree to Key CLI Commands**

In the `## Key CLI Commands` section, after the `export.image` line, add:

```bash
ai-happy-design command design_system.analyze '{}'                         # scan file for existing styles/components
ai-happy-design command node.get_tree '{"nodeId": "...", "compact": true}' # compact flat tree (3-5x fewer tokens)
ai-happy-design command document.get_info '{}'
ai-happy-design command document.get_selection '{}'
ai-happy-design command document.scan_by_type '{"nodeType": "FRAME"}'
```

**Step 3: Add domain confusion to Common Pitfalls**

At the end of the `## Common Pitfalls` section, add:

```markdown
- **Wrong domain for pages**: `page.get_all` not `document.list_pages` — page ops always use `page.*`
- **Wrong domain for layer ops**: `layer.group` / `layer.bring_to_front` / `layer.send_to_back` — NOT `node.*`
- **Design system scan**: `design_system.analyze` not `document.analyze` or `document.get_design_system`
- **Compact tree for discovery**: `node.get_tree {"nodeId":"...", "compact":true}` — 3-5x fewer tokens
```

**Step 4: Verify SKILL.md line count stays reasonable**

```bash
wc -l "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/skills/ai-happy-design/SKILL.md"
```
Target: under ~470 lines (currently 423, adding ~40 lines).

**Step 5: Commit**

```bash
git add skills/ai-happy-design/SKILL.md
git commit -m "docs: add page management, layer ops, design_system to SKILL.md

Surfaces page.get_all, page.create, page.set_current etc. which were
entirely absent. Adds domain confusion pitfalls so LLMs guess the right
domain (page.* not document.*, layer.* not node.*).

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Deploy and Verify

**Step 1: Deploy (required after any Go change)**

```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design"
make deploy
```

**Step 2: Test all new compound aliases work**

```bash
AHD_CHANNEL=golden-waffle-25 ai-happy-design command document.list_pages '{}'
AHD_CHANNEL=golden-waffle-25 ai-happy-design command document.get_pages '{}'
# Both should return page list, not error
```

**Step 3: Verify skill sync (if using skillshare)**

```bash
skillshare sync 2>/dev/null || echo "skillshare not configured, skipping"
```

---

## Summary of Changes

| File | Change |
|------|--------|
| `internal/ws/command_routing.go` | Add `compoundAliases` map + pre-check in `resolveCommandRoute` |
| `skills/ai-happy-design/SKILL.md` | Add Page Management section, expand Key CLI Commands, add domain pitfalls |

**No plugin changes needed** — the fix is purely in Go routing layer.
