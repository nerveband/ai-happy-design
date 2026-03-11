# Pre-Release Checklist

Run this checklist before every release to ensure quality and consistency.

## Build Verification
- [ ] `go build ./...` compiles without errors
- [ ] `go test ./...` all tests pass
- [ ] `cd plugin && npm run check` plugin builds and syntax verification passes
- [ ] `make deploy` completes successfully (build + sign + install + restart)

## Schema & Documentation Sync
- [ ] Every plugin handler action has a corresponding schema in `internal/schema/`
- [ ] `describe.go` descriptions match schema param lists
- [ ] `catalog_llm.go` references match actual command names
- [ ] SKILL.md reflects current features and aliases
- [ ] AGENTS.md is up to date with architecture changes
- [ ] README.md reflects current installation and usage

## Agent DX Scoring (Target: 21/21)
Run `ai-happy-design` through the Agent DX CLI Scale:

1. **Machine-Readable Output (3/3)**: All commands output JSON. Errors are structured JSON on stderr. TTY auto-detection works.
2. **Raw Payload Input (3/3)**: All mutating commands accept raw JSON payload. Convenience aliases work within JSON.
3. **Schema Introspection (3/3)**: `schema --json` lists all schemas. `schema <cmd> --json` returns full typed schema. All commands covered.
4. **Context Window Discipline (3/3)**: `--fields` flag on commands. `compact:true` on tree queries. Pagination on list commands. SKILL.md has field mask guidance.
5. **Input Hardening (3/3)**: Control character rejection. Path traversal prevention. Fuzzy matching for hallucinations. Response sanitization.
6. **Safety Rails (3/3)**: `--dry-run` on all mutating commands. `validate` for batch dry-run. `--strict-quality` quality gate. Response sanitization.
7. **Agent Knowledge Packaging (3/3)**: SKILL.md with guardrails. 3-layer discoverability. Versioned catalog. Reference files.

## Shiptypes Compliance
- [ ] Schema is single source of truth for all command contracts
- [ ] No manual duplicate artifacts that can drift from schemas
- [ ] describe.go is generated from or validated against schemas
- [ ] Coverage test verifies schema<>handler alignment
- [ ] Breaking changes caught at build time (coverage test fails)

## Design Intelligence
- [ ] `guide` returns compact methodology
- [ ] `guide --topic X` returns deep-dive for each topic (typography, color, layout, depth, states, quality)
- [ ] SKILL.md contains auto-delivered design methodology
- [ ] Quality gates in `--lint` catch design issues (contrast, hierarchy, grid alignment)

## Plugin Verification
- [ ] No ES2020+ syntax in dist/code.js (optional chaining, nullish coalescing, object spread)
- [ ] Plugin connects to relay and responds to commands
- [ ] Export at scale 2 works for standard frame sizes
- [ ] All handler actions respond correctly

## Release Steps
1. Verify all checks above pass
2. `git tag v<version>` (check latest with `git tag --sort=-v:refname | head -1`)
3. `git push origin v<version>`
4. `GITHUB_TOKEN=$(gh auth token) goreleaser release --clean`
5. `ai-happy-design upgrade` to verify
6. Update SKILL.md if CLI features changed
7. `skillshare sync` to distribute skill updates
8. Reopen Figma plugin to load new code.js
