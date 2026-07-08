# Agent CLI Audit: ahd-figma

Binary tested: `bin/ahd-figma`
Version: `0.0.0-dev`
Date: `2026-07-08`

Checklist: local `cli-best-practices-audit` 85-point agent CLI audit.
Older baseline: the May 2026 roadmap recorded the external AHD scorecard as `50/50` for the original 50 checks. This refresh uses the newer 85-check rubric.

## Score

| Category | Score |
|----------|-------|
| Discoverability | 7/7 |
| Structured output | 5/6 |
| Input flexibility | 4/5 |
| Safety rails | 5/6 |
| Error handling | 5/7 |
| Context discipline | 5/5 |
| Predictability | 6/7 |
| Agent knowledge | 7/7 |
| Resilience | 5/7 |
| Distribution | 4/5 |
| Three-layer introspection | 5/5 |
| Persistent identity/config | 3/5 |
| Two-way I/O/artifacts | 4/5 |
| Contract/generation discipline | 4/5 |
| Unix composability/restraint | 4/5 |
| API-native payload ergonomics | 4/5 |
| Domain depth/proof gates | 5/5 |
| Total | 78/85 |

## Drift Found And Fixed

1. `bin/ahd-figma status --json` failed with `unknown flag: --json`.
   Fixed by accepting `--json` as a compatibility flag; status already emits JSON.
2. `bin/ahd-figma batch --dry-run docs/examples/instagram-post-v3.json` failed with `unknown flag: --dry-run`.
   Fixed by adding a relay-free dry run path that loads, normalizes, schema-validates, and prints JSON.
3. Several schema entries had no safety metadata, and REST schemas used the older `non-idempotent` spelling.
   Fixed through registry normalization plus enforcement tests.
4. MCP lacked prompts and tool safety annotations.
   Fixed with `prompts/list`, `prompts/get`, and schema-derived annotations.

## Evidence Notes

- `bin/ahd-figma --version` -> `ahd-figma version 0.0.0-dev`.
- `bin/ahd-figma --help` -> lists `agent-context`, `schema`, `tools`, `validate`, `batch`, `mcp`, and other subcommands.
- `bin/ahd-figma schema --json` -> valid JSON, 142990 bytes.
- `bin/ahd-figma tools --llm --json` -> valid JSON catalog backed by `internal/tools/catalog_llm.go`.
- `bin/ahd-figma status --json` -> valid JSON, 134 bytes in this no-relay environment.
- `bin/ahd-figma batch --dry-run docs/examples/instagram-post-v3.json` -> valid JSON, 54700 bytes, no relay required.
- `bin/ahd-figma agent-context --json` -> valid JSON, 2237 bytes.
- MCP `tools/list`, `resources/list`, `prompts/list`, and `prompts/get` -> four valid JSON-RPC responses.
- `bin/ahd-figma command verify.visual -p '{"artifactPath":"/tmp/ahd-agent-context.json","target":"gate"}'` -> valid JSON proof envelope for an existing artifact.

## Remaining Gaps

- Structured error shape is still not canonical across every Cobra path; some errors remain prose plus usage.
- Exit codes are still mostly `1` for broad failure classes.
- Persistent profiles/config inspection is partial; there is no complete `profile list/use/inspect --redacted` surface.
- Durable async jobs ledger and delivery sinks are not implemented yet.
- Full `--dry-run` semantics for every mutating command are schema-backed at the command/batch layer, but not yet implemented as per-command side-effect previews for every plugin write.

## Live Figma Checks

Not run in this pass because no live Figma plugin document/channel was available. The implemented surfaces are ready for a dedicated test page:

- `document.screenshot`
- `document.screenshot_selection`
- central `figma.commitUndo()` behavior and `commitUndo:false`
- Motion capability probes
- Shader capability probes
- Slot capability probes
- grid reorder and auto-flow field readback
- `noiseSizeVector`, `noiseSizeX`, and `noiseSizeY`
