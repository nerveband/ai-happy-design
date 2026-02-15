# LLM-Friendly CLI + MCP Notes (2026-02-13)

## Scope
- Repo: `.`
- Goal: make CLI/MCP highly discoverable and composable for LLM agents.

## What “LLM-friendly” means in practice
1. Discoverable capabilities with schemas.
2. Deterministic, machine-readable outputs.
3. Predictable error contracts and correlation IDs.
4. Native support for tool chaining and progress.
5. Pagination/streaming for large outputs.

## Primary-source findings

### MCP protocol requirements that matter
- MCP tools are discoverable via `tools/list` and include `inputSchema` (JSON Schema).
- MCP servers should expose capabilities (`tools`, `resources`, `prompts`) and support `listChanged` notifications when applicable.
- `tools/list`, `resources/list`, and `prompts/list` are pagination-capable.
- Lifecycle requires initialization/capability negotiation first, then normal operation.

### JSON-RPC response discipline
- A response must include either `result` or `error`, never both.
- Correlated request IDs are mandatory for robust tool chains.

### CLI machine interface conventions (mature precedent)
- Stable machine output should be explicit (`--json`).
- Field-select/transform affordances (`--jq`, templates) dramatically improve script/agent interoperability.

### Claude Code MCP ergonomics
- Clients can discover and read MCP resources automatically.
- Prompt/tool/resource discovery is dynamic.
- Large output limits matter; pagination/filtering is preferred for oversized payloads.

## Sources
- MCP Lifecycle: https://modelcontextprotocol.io/specification/latest/basic/lifecycle
- MCP Tools: https://modelcontextprotocol.io/specification/2024-11-05/server/tools
- MCP Resources: https://modelcontextprotocol.io/specification/2025-03-26/server/resources
- MCP Pagination: https://modelcontextprotocol.io/specification/2024-11-05/server/utilities/pagination
- JSON-RPC 2.0: https://www.jsonrpc.org/specification
- GitHub CLI formatting (`--json`, `--jq`): https://cli.github.com/manual/gh_help_formatting
- Claude Code MCP docs: https://code.claude.com/docs/en/mcp

## Current project status vs these patterns

### Already aligned
- `tools --json` capability catalog exists.
- `command` and `batch` support tool chaining patterns.
- Wrapped response contract with `id + result|error` is implemented.
- Presence checks for wrapped result/error avoid falsey-result timeout bugs.
- Progress updates are supported.

### Remaining improvements (recommended)
1. Add `--output` modes for CLI (`json`, `jsonl`, `text`) with strict machine-stable JSON objects.
2. Add optional `--jq` filtering in CLI output path.
3. Add explicit error code taxonomy (`code`, `message`, `details`) for CLI + MCP parity.
4. Add tests for output contract stability (including falsey values and nested errors).
5. Expose richer `describe` metadata for each action (required args, types, examples).
6. Add pagination knobs to scanning/listing operations that can return large payloads.
7. Add a “capability summary” endpoint in relay status for richer agent planning.

## Suggested implementation order
1. Output contract hardening (`--output`, error code schema).
2. `--jq` integration.
3. Action metadata expansion in `tools --json`.
4. Pagination support for large list/scan operations.
5. Contract tests + golden output tests.
