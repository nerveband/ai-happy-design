# AHD Illustrator Agent-First CLI Spec (Monorepo in `ai-happy-design`)

## Summary
This is the implementation-ready spec for adding Illustrator support as a first-class CLI in your existing monorepo root at `ai-happy-design`, while hard-switching naming to `ahd-figma` and `ahd-illustrator` now.

### Locked decisions
| Area | Decision |
|---|---|
| Monorepo strategy | Convert existing `ai-happy-design` repo in place to monorepo root |
| Primary new tool | `ahd-illustrator` |
| Existing tool naming | Hard switch to `ahd-figma` now |
| Stack | Go CLI + C++ Illustrator plugin bridge |
| Platform v0.1 | macOS only |
| Surface v0.1 | CLI only (no MCP in v0.1) |
| Brand | `ahd-illustrator` with non-affiliation trademark disclaimer |
| License | MIT (retain existing) |

## 1. Repo Migration Plan (In-Place Monorepo)
1. Keep existing GitHub repo as public root (no new repo creation now).
2. Convert top-level docs and release flow to “AHD Design Monorepo”.
3. Introduce two binaries:
`ahd-figma` (renamed from current CLI)
`ahd-illustrator` (new)
4. Do not edit currently dirty files unless explicitly required:
`internal/designlint/lint.go`
`internal/tools/catalog_llm.go`
`plugin/src/handlers/effect.ts`
`plugin/src/handlers/shape.ts`
5. Use a clean worktree/branch for implementation to avoid collisions with local dirty state.

## 2. File/Folder Spec
Create or repurpose these paths:

```text
cmd/ahd-figma/main.go
cmd/ahd-illustrator/main.go

internal/commoncli/
internal/commonschema/
internal/commonvalidate/

internal/illustrator/host/
internal/illustrator/bridge/
internal/illustrator/schema/
internal/illustrator/commands/
internal/illustrator/inspect/
internal/illustrator/validate/

tools/illustrator/bridge/ahd_bridge.jsx
tools/illustrator/plugin-cpp/
tools/illustrator/plugin-cpp/CMakeLists.txt
tools/illustrator/plugin-cpp/src/

skills/ahd-illustrator/SKILL.md
docs/illustrator/README.md
docs/illustrator/architecture.md
docs/illustrator/commands.md
docs/illustrator/plugin-build.md

docs/plans/2026-03-05-ahd-illustrator-monorepo-spec.md
```

## 3. Public CLI Interface Spec

### 3.1 Binaries
1. `ahd-figma` (existing functionality under new name)
2. `ahd-illustrator` (new Illustrator automation CLI)

### 3.2 Core commands (`ahd-illustrator`)
1. `ahd-illustrator tools [--json] [--llm]`
2. `ahd-illustrator schema [domain.action] [--json] [--all --llms-txt]`
3. `ahd-illustrator command <domain.action> --json '<payload>' [--dry-run] [--fields '<mask>'] [--output json|ndjson|text]`
4. `ahd-illustrator batch --ops <file|json> [--dry-run] [--strict] [--output json|ndjson]`
5. `ahd-illustrator host status|open|quit`
6. `ahd-illustrator doctor`
7. `ahd-illustrator examples [category]`

### 3.3 Output contract (agent-first)
Every command returns a stable envelope:

```json
{
  "ok": true,
  "requestId": "uuid",
  "command": "text.create",
  "result": {},
  "warnings": [],
  "timingMs": 0
}
```

Error envelope:

```json
{
  "ok": false,
  "requestId": "uuid",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "human-readable",
    "details": {}
  },
  "retryable": false
}
```

Batch envelope:

```json
{
  "ok": true,
  "requestId": "uuid",
  "summary": {"total": 0, "succeeded": 0, "failed": 0},
  "steps": [{"index": 0, "name": "optional", "ok": true, "result": {}}],
  "timingMs": 0
}
```

## 4. Illustrator Execution Pipeline Spec

### 4.1 Transport chain
`ahd-illustrator` -> AppleScript host adapter (`do javascript`) -> JSX bridge runtime -> `app.sendScriptMessage(...)` to C++ plugin when needed.

### 4.2 Why this is locked
Local Illustrator AppleScript dictionary confirms commands for:
`do javascript`, `do script`, `execute menu command`, `load action`, `unload action`, `Select Tool`, and script-message behavior.
This gives maximal command coverage with least friction while avoiding foreground dependency for most flows.

### 4.3 Bridge protocol
Plugin selectors:
1. `ahd.capabilities`
2. `ahd.exec`
3. `ahd.inspect`
4. `ahd.version`

Request payload:

```json
{
  "v": "1.0",
  "id": "uuid",
  "command": "domain.action",
  "params": {},
  "dryRun": false,
  "timeoutMs": 30000
}
```

Response payload:

```json
{
  "v": "1.0",
  "id": "uuid",
  "ok": true,
  "result": {},
  "warnings": []
}
```

If plugin is unavailable:
1. Script-only commands continue working.
2. Plugin-required commands fail with `PLUGIN_REQUIRED` and remediation text.

## 5. v0.1 Command Domain Coverage
Implement these command groups in v0.1:

1. `app.*`: info, version, select_tool, execute_menu, user_interaction_level
2. `document.*`: new, open, save, save_as, close, list, info
3. `artboard.*`: list, create, resize, set_active, fit_to_artwork
4. `layer.*`: list, create, rename, visibility, lock, reorder
5. `selection.*`: get, clear, set_by_ids, select_by_name
6. `path.*`: create_rect, create_ellipse, create_path, transform, duplicate
7. `text.*`: create, set_contents, set_style, outline
8. `appearance.*`: set_fill, set_stroke, set_gradient, apply_graphic_style
9. `action.*`: load, run, unload
10. `export.*`: png, jpg, svg, pdf, ai
11. `inspect.*`: tree, styles, bounds, fonts, summary

## 6. Agent-First Rules (Shiptypes + Justin Principles)
1. Raw payload first: full JSON bodies, not flag-heavy bespoke arguments.
2. Runtime introspection first: `schema` and `tools --llm` are canonical.
3. Context discipline: support `--fields` and NDJSON streaming.
4. Input hardening: reject path traversal, control chars, malformed IDs, encoded tricks.
5. Safety rails: `--dry-run` on all mutating commands.
6. Skill-first packaging: `skills/ahd-illustrator/SKILL.md` plus `llms.txt` generation.
7. Deterministic machine output by default.

## 7. Validation and Security Spec
1. Add schema-driven validation for all commands before execution.
2. Add fuzzy correction only for low-risk fields, never for destructive targets.
3. Sandboxed output paths to CWD unless explicit override is approved by flags.
4. Enforce timeout and cancellation on Illustrator calls.
5. Add structured error codes:
`VALIDATION_ERROR`, `HOST_EXEC_ERROR`, `PLUGIN_REQUIRED`, `PLUGIN_TIMEOUT`, `ILLUSTRATOR_NOT_RUNNING`, `UNSUPPORTED_COMMAND`.

## 8. Testing Spec

### 8.1 Unit tests
1. Schema registry lookup and aliasing.
2. Input validator and hardening rules.
3. JSON envelope stability and error-code mapping.
4. Batch normalization and step interpolation safety.

### 8.2 Integration tests (macOS, tagged)
1. Illustrator launch/status/quit behavior.
2. `do javascript` execution path.
3. `sendScriptMessage` path with plugin installed.
4. Plugin-missing fallback path.
5. Export commands producing expected artifacts.

### 8.3 Acceptance scenarios
1. Agent can create multi-artboard design end-to-end with `batch`.
2. Agent can inspect an existing `.ai` file and return style summary JSON.
3. Agent can rerun same command idempotently with deterministic output shape.
4. No foreground steal in standard non-UI-only command flows.

## 9. README + Initial Docs Commit Spec
First docs-focused commit in monorepo should include:

1. Root `README.md` rewritten to monorepo framing:
`ahd-figma` and `ahd-illustrator`
architecture diagram
quickstart for both CLIs
agent-first design principles
2. New `docs/illustrator/` docs set.
3. New plan file at:
`docs/plans/2026-03-05-ahd-illustrator-monorepo-spec.md`
4. Trademark disclaimer:
“Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.”

## 10. Commit Plan (Execution Mode)
Use this exact commit sequence:

1. `chore(monorepo): introduce ahd binary naming and monorepo docs scaffold`
2. `feat(illustrator): scaffold ahd-illustrator cli, schema, and validation core`
3. `feat(illustrator): add applescript/jsx host bridge and command executor`
4. `feat(illustrator-plugin): add c++ plugin bridge with sendScriptMessage handlers`
5. `feat(illustrator): implement v0.1 command domains and batch engine`
6. `test(illustrator): add unit/integration coverage and acceptance fixtures`
7. `docs(illustrator): publish command reference, skill, and release notes`

## 11. Release Plan
1. Update `.goreleaser.yml` to publish `ahd-figma` and `ahd-illustrator` binaries.
2. Keep MIT license and add repo topics for discoverability.
3. Publish release notes highlighting:
agent-first JSON contract
schema introspection
dry-run safety
plugin capability mode
4. Mark Illustrator support as `v0.1 (macOS)`.

## 12. Assumptions and Defaults
1. No new GitHub repo is created now; existing public `ai-happy-design` remains monorepo root.
2. CLI defaults to machine-readable output; human format is optional.
3. Plugin is required for maximum coverage; script-only fallback is intentionally partial.
4. v0.1 excludes MCP by choice and prioritizes direct CLI agent workflows.
5. Existing dirty files remain untouched unless explicitly requested.

## Sources Used
1. [Ship types, not docs](https://shiptypes.com)
2. [You Need to Rewrite Your CLI for AI Agents](https://justin.poehnelt.com/posts/rewrite-your-cli-for-ai-agents/)
3. [Adobe Illustrator Developer Portal](https://developer.adobe.com/illustrator/)
4. [Illustrator Scripting Language Support](https://ai-scripting.docsforadobe.dev/introduction/scriptingLanguageSupport/)
5. [Illustrator Executing Scripts](https://ai-scripting.docsforadobe.dev/introduction/executingScripts/)
6. [Illustrator JavaScript Application Object](https://ai-scripting.docsforadobe.dev/jsobjref/Application/)
7. Local Illustrator 2026 AppleScript dictionary inspection (`sdef`) confirming command surfaces (`do javascript`, `execute menu command`, `do script`, `Select Tool`, action load/unload, script-message command).
