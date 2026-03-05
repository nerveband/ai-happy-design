# Documentation

This repo now contains both Figma and Illustrator surfaces.

## Monorepo

- `illustrator/README.md` - entry point for the new `ahd-illustrator` CLI.
- `illustrator/architecture.md` - host, JSX, and plugin bridge design.
- `illustrator/commands.md` - command domains and output contract.
- `illustrator/plugin-build.md` - C++ plugin bridge build notes.
- `illustrator/release-notes-v0.1.md` - initial Illustrator release notes.
- `plans/2026-03-05-ahd-illustrator-monorepo-spec.md` - locked implementation spec.

## Figma

- `getting-started.md` - install/build/run steps for `ahd-figma` and the plugin.
- `prebuild-checklist.md` - prerequisites and prebuild validation before running.
- `llm-integration.md` - connect the Figma MCP server to LLM tools/editors.
- `llm-discovery-playbook.md` - discovery-first flow and prompt patterns for low-failure agent usage.
- `cli-batch-and-payloads.md` - how to send single vs multi-operation payloads.
- `troubleshooting.md` - common issues and fast fixes.

## Technical Notes

- `image-fill-contract.md` - URL/base64 image fill behavior and response contract.
- `pairing-turnkey-design.md` - sticky pairing design and relay/channel strategy.
- `llm-friendly-cli-mcp-notes-2026-02-13.md` - design notes and standards references.

## Archive

- `archive/codex-session-2026-02-13.md` - session log with local-machine debugging context.
- `archive/figma-image-debug.md` - historical image fill debugging notes.
