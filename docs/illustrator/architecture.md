# Illustrator Architecture

## Execution Chain

`ahd-illustrator` uses a layered transport model:

1. Go CLI entrypoint
2. Schema validation and input hardening
3. AppleScript host adapter on macOS
4. JSX bridge runtime with an ES3-safe serializer
5. Optional `app.sendScriptMessage(...)` handoff into the C++ plugin bridge

## Why This Shape

- AppleScript exposes `do javascript`, `do script`, `execute menu command`, and action loading without requiring a foreground-only UI workflow.
- JSX remains the broadest low-friction scripting layer for v0.1.
- Illustrator 30.2.1 on macOS does not expose a global `JSON` object inside ExtendScript, so the bridge ships its own serializer for the machine envelope.
- The plugin bridge is reserved for capabilities and deep inspection that are awkward or unreliable via pure scripting.
- Runtime discovery stays close to the executable surface: `tools`, `schema`, and `app.info` expose the typeable contract and the live preset/runtime state that agents need.
- Input hardening follows an agent-first rule set: reject suspicious paths and opaque identifiers up front, then validate nested JSON payload structure and command-specific invariants before any host call.

## Runtime Modes

- Script-only: commands that can be fulfilled entirely through ExtendScript.
- Plugin-capable: commands that use `sendScriptMessage` when the bridge is installed and respond to the `ahd.version` probe.
- Dry-run: validate, normalize, and emit the stable machine envelope without touching Illustrator.

## Script-First Behaviors Verified Live

- `document.new` maps directly to Illustrator's multi-artboard document APIs, including color space, layout, spacing, and preset-backed creation.
- `document.save_as` performs explicit format-aware saves for `ai` and `pdf`.
- `path.create_rect` uses `roundedRectangle(...)` for non-zero corner radii.
- `export.png`, `export.jpg`, and `export.svg` apply artboard targeting in the script layer rather than treating `artboardId` as inert schema.

## Shared Monorepo Packages

- `internal/commoncli`: response envelopes, request IDs, output writers
- `internal/commonschema`: schema registry and JSON export
- `internal/commonvalidate`: low-risk fuzzing, path hardening, and validation
- `internal/illustrator/...`: host, schema, inspect, validate, and command executors
