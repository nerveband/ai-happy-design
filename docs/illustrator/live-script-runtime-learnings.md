# Illustrator Live Script Runtime Learnings

This document records the differences between the Adobe scripting references, the Illustrator 30.2.1 runtime on this Mac, and the agent-first CLI behavior we normalized in `ahd-illustrator`.

## Sources Reviewed

- `Illustrator AppleScript Scripting Ref (Nov-2025).pdf`
- `Illustrator JavaScript Scripting Reference (Nov-2025).pdf`
- `Illustrator VBScript Reference (Nov-2025).pdf`
- `mac Adobe Illustrator 2026 SDK`
- `windows Adobe Illustrator 2026 SDK`

## Agent-First Rules Applied

- `tools` and `schema` remain the canonical discovery surface.
- Commands take raw JSON payloads first, not bespoke one-off flags.
- Validation rejects malformed identifiers, encoded traversal tricks, query/fragment-like targets, and unsafe output paths.
- Host quirks are normalized in the CLI instead of leaked to the caller whenever possible.

These choices follow the same direction as [Ship Types](https://shiptypes.com/) and Justin Poehnelt's input-hardening guidance for AI-facing CLIs.

## Runtime Facts That Matter

### 1. ExtendScript in Illustrator 30.2.1 does not provide a usable global `JSON`

- The checked-in bridge runtime originally depended on `JSON.stringify`.
- Live runtime validation showed that assumption was wrong for this host.
- The CLI bridge now ships its own ES3-safe serializer so command envelopes still return valid JSON.

### 2. Illustrator host execution is sensitive to frontmost/UI state

- Even non-obviously interactive `do javascript` calls can stall when Illustrator is backgrounded or left in a bad state.
- Clipboard-style commands need frontmost activation.
- Large live integration runs are more reliable when split into smaller restart-isolated families instead of one huge stateful flow.

### 3. `View.zoom` does not behave like the doc wording suggests

- The JavaScript reference describes `100.0` as 100%.
- Live `view.info` on this host returned values like `1.38`.
- The CLI now documents `view.set_zoom.zoom` as a zoom factor rather than a percentage to match observed behavior.

### 4. `document.rearrangeArtboards` is real despite a doc typo

- The runtime exposes `document.rearrangeArtboards`.
- Some reference material also shows the typo `rearrangeArboards`.
- The CLI normalizes this by checking both names and using whichever exists.

### 5. Text conversion methods can return `null` even when the conversion succeeds

- `convertPointObjectToAreaObject()` and `convertAreaObjectToPointObject()` do not reliably return the converted object.
- The CLI treats a `null` return as "conversion happened in place" and falls back to the original text frame.

## Script-First Discoveries Beyond The PDFs

### `SymbolItem.breakLink()` exists live

- The JavaScript reference pages for `SymbolItem` only document `duplicate` and `move`.
- Live reflection against the Illustrator object model showed a real `breakLink()` method on `SymbolItem`.
- `ahd-illustrator` now exposes that through `symbol.break_link`.

### `document.exportForScreens` exists live but is not agent-safe yet

- Live reflection showed `doc.exportForScreens` as a function.
- The scripting PDFs expose the `ExportForScreens` option shape, but do not document a deterministic, CLI-safe invocation pattern.
- Direct live probes from `do javascript` hung or behaved like an interactive/UI path.
- Because the current goal is an agent-first CLI, this surface should not be exposed until it can run deterministically without modal UI.

### Workspace enumeration is not exposed in the scripting references

- The references expose:
  - `saveWorkspace`
  - `switchWorkspace`
  - `resetWorkspace`
  - `deleteWorkspace`
- They do not expose a documented `list workspaces` script API.
- That means `workspace.list` should not be treated as part of the script-first contract unless we deliberately adopt host UI scripting for it.

## Menu-Command Normalizations

Some view operations are only safely expressible through menu toggles, not direct JavaScript setters.

Validated menu command strings from the local SDK/runtime:

- Ruler visibility toggle: `ruler`
- Transparency grid toggle: `TransparencyGrid Menu Item`

The CLI now uses these as normalized setters for:

- `view.set_ruler_visibility`
- `view.set_transparency_grid_visibility`

`document.arrange` also falls into this bucket on Illustrator 30.2.1.

- The JavaScript reference documents a document arrange surface.
- Live runtime probes on this host did not expose `doc.arrange(...)` or `app.arrange(...)`.
- The CLI therefore normalizes `document.arrange` to the stable menu commands:
  - `cascade`
  - `tile`

Those commands are implemented as idempotent "set desired state" operations instead of leaking toggle semantics to the caller.

## Print Runtime Learnings

- The scripting references document `printerList`, `PPDFileList`, and `printPresetsList`.
- On this host, live enumeration of printers and PPDs was unreliable enough to hang the session or raise `RTUO`.
- `print.devices` therefore reports deterministic degraded output:
  - `printers: []`
  - `printerEnumeration: "skipped"`
  - `ppdFiles: []`
  - `ppdEnumeration: "skipped"`
  - `presets: [...]`
- `app.preset_lists` keeps the preset/color/tracing lists, but it degrades printer/PPD enumeration to:
  - `printerEnumeration: "unavailable"`
  - `ppdEnumeration: "unavailable"`

That is an intentional stability choice for the current CLI.

## Export And Artifact Naming Learnings

- TIFF export file names are not stable on this host.
  - Requested `basename.tif` can materialize as:
    - `basename.tif`
    - `basename.tiff`
    - `basename-01.tif`
- Window capture can produce `basename.tif.tiff`.
- PDF and print preset export can ignore the requested file name and emit a preset-derived artifact name instead.

The CLI now snapshots the destination folder and resolves the actual created artifact path before returning it.

## Validation / Hardening Learnings

- Output paths must be normalized under the current working directory for predictable agent behavior.
- Opaque identifiers should reject:
  - URL-like strings
  - query fragments
  - fragment syntax
  - encoded traversal markers
- Cross-field validation matters for Illustrator:
  - `document.new` artboard layout combinations
  - typed preference writes
  - RGB/CMYK swatch and spot creation
  - point-array shape for full path replacement

## What Phase 1 Can Reliably Cover Without Plugin Or SDK

Reliable script-first coverage on this branch includes:

- `app.*`
- `document.*` including PDF/print preset import/export helpers
- `artboard.*`
- `layer.*`
- `selection.*`
- `page_item.*`
- `path.*`
- `text.*`
- `appearance.*`
- `action.*`
- `export.*` except `export.for_screens`
- `inspect.*`
- `preference.*`
- `view.*`
- `matrix.*`
- `perspective.*`
- `workspace.save|switch|reset|delete`
- `swatch.*`
- `spot.*`
- `style.character.*`
- `style.paragraph.*`
- `style.graphic.*`
- `symbol.*`
- `placed.*`
- `raster.*`
- `repeat.*`
- `dataset.*`
- `variable.*`
- `trace.preset.list`
- `capture.*`
- `print.*`

## Open Script-First Boundaries

- `workspace.list`
  - not documented in the scripting references
  - would require host UI scripting or a different integration path
- `export.for_screens`
  - hidden live method exists
  - current `do javascript` behavior is not yet deterministic enough for agent-first CLI exposure
- `document.write_as_library`
  - documented in the JavaScript reference
  - live 30.2.1 rejected both `LibraryType.IllustratorArtwork` and the other exposed `LibraryType` variants with `Invalid enumeration value`
  - excluded from the Phase 1 public CLI surface until the runtime enum mapping is proven
- `page_item.bring_in_perspective`
  - documented in the JavaScript reference
  - live 30.2.1 rejected the exposed `PerspectiveGridPlaneType` values with `Invalid enumeration value`
  - excluded from the Phase 1 public CLI surface until there is a deterministic script-only path
- `trace.preset.store`
  - documented in the JavaScript reference
  - live 30.2.1 returned Illustrator error 100 instead of storing a preset from traced plugin art
  - excluded from the Phase 1 public CLI surface until there is a deterministic script-only path
- Very deep text-range styling surfaces
  - possible in script
  - should be exposed only with tighter schemas so the CLI does not become a vague bag of partially validated text mutations

## Recommendation For Future Contributors

- Treat the local Adobe references as necessary but not sufficient.
- Probe the live runtime before claiming a surface is scriptable in an agent-friendly way.
- Prefer deterministic script-backed commands over UI-driven menu workflows.
- Document every mismatch between references and live behavior as soon as it is discovered.
