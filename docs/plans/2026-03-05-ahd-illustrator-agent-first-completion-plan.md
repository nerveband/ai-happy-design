# AHD Illustrator Agent-First Completion Plan

## Summary
- Current state: `ahd-illustrator` already delivers the original v0.1 command set on macOS through a script-first transport (`CLI -> AppleScript -> ExtendScript`), with stable JSON envelopes, schema introspection, and input hardening.
- The local Adobe references show that Illustrator's script surface is much broader than the current CLI, and the native plugin path is still incomplete. The current `tools/illustrator/plugin-cpp/` target is a bridge stub, not a real Illustrator `.aip`.
- Completion strategy: keep the CLI script-first for the majority of user-facing work, expand the typed command surface to cover the broader scripting model, and ship the native plugin as an optional accelerator and suite-backed extension. Users should never need the SDK. Maintainers need the SDK only to build precompiled plugins.

## Source Material Studied
- `Downloads/adobe-illustrator/Illustrator JavaScript Scripting Reference (Nov-2025).pdf`
- `Downloads/adobe-illustrator/Illustrator AppleScript Scripting Ref (Nov-2025).pdf`
- `Downloads/adobe-illustrator/Illustrator VBScript Reference (Nov-2025).pdf`
- `Downloads/adobe-illustrator/mac Adobe Illustrator 2026 SDK/ReadMe.txt`
- `Downloads/adobe-illustrator/mac Adobe Illustrator 2026 SDK/docs/guides/getting-started-guide.pdf`
- `Downloads/adobe-illustrator/mac Adobe Illustrator 2026 SDK/docs/guides/porting-guide.pdf`
- `Downloads/adobe-illustrator/mac Adobe Illustrator 2026 SDK/samplecode/ScriptMessage/*`
- `Downloads/adobe-illustrator/mac Adobe Illustrator 2026 SDK/samplecode/common/mac/AIPluginCommon.xcconfig`
- `Downloads/adobe-illustrator/windows Adobe Illustrator 2026 SDK/ReadMe.txt`
- `Downloads/adobe-illustrator/windows Adobe Illustrator 2026 SDK/samplecode/ScriptMessage/*`

## Key Findings

### 1. Script coverage is much broader than the current CLI
- The JavaScript and AppleScript references expose far more than the current v0.1 domains.
- Important additional scriptable areas include:
  - `brush.*`
  - `character style.*`
  - `paragraph style.*`
  - `dataset.*`
  - `variable.*`
  - `view.*`
  - `preferences.*`
  - `matrix.*`
  - `spot.*`
  - `swatch.*`
  - `symbol.*`
  - `placed.*`
  - `raster.*`
  - `trace.*`
  - `repeat.grid.*`
  - `repeat.radial.*`
  - `repeat.symmetry.*`
  - richer `export.*` formats including GIF, PNG8, TIFF, WebP, Photoshop, AutoCAD, and Export For Screens
  - `image capture`
  - `print.*`
- The AppleScript command table also confirms scriptable operations like copy, cut, paste, redo, undo, workspace switching, preference getters/setters, style import, variable import/export, rasterize, relink, and perspective grid operations.

### 2. The native plugin gap is real
- The Adobe SDK sample for this use case is `samplecode/ScriptMessage`.
- Adobe expects a real Illustrator plugin:
  - product type `.aip`
  - `Plugin` / `Suites` bootstrap from `samplecode/common/source`
  - `PluginMain` entry point
  - valid PiPL/resource generation
  - bundle packaging recognized by Illustrator
- Only plugins with a valid PiPL and entry point are considered loadable by Illustrator.
- The current repo plugin code is not yet that:
  - `tools/illustrator/plugin-cpp/CMakeLists.txt` builds a generic shared library
  - `tools/illustrator/plugin-cpp/src/PluginMain.cpp` exports helper functions but not a true Illustrator SDK entry point
  - `tools/illustrator/plugin-cpp/src/BridgeHandlers.cpp` is still echo/stub logic

### 3. macOS and Windows SDKs are structurally aligned
- Both SDKs include:
  - `illustratorapi/`
  - `samplecode/common/`
  - `samplecode/ScriptMessage/`
  - `tools/pipl/`
- The Windows SDK confirms the same `ScriptMessage` pattern and `.aip` target via Visual Studio/MSBuild.
- The VBScript reference broadly mirrors the JavaScript object model, which makes a future Windows host path plausible once the macOS path is solid.

### 4. User distribution should remain plugin-optional
- Illustrator users should be able to install Illustrator and the `ahd-illustrator` CLI only.
- The plugin should be optional, prebuilt, and installable by users without the SDK.
- The SDK is only a maintainer build dependency.

### 5. The scripting model has normalization traps that belong in the CLI
- Canonical CLI units should remain points for geometry, degrees for rotation, and percentages for matrix scale.
- Illustrator coordinate systems are inconsistent across document origin, artboard origin, item position, and translation. The CLI must normalize these differences rather than exposing raw host quirks.
- Selection representations differ across scripting hosts:
  - JavaScript can return `null`
  - AppleScript uses `{}`
  - VBScript uses `Empty`
  - text selections can be insertion points or text ranges rather than page items
- Export behavior is format-specific:
  - most export calls append extensions automatically
  - Photoshop/PSD requires the extension to be explicitly present
- Clipboard operations in AppleScript require Illustrator to be frontmost.
- Tracing is asynchronous and needs redraw-aware follow-up handling.
- Color payload validation should remain document-color-space aware, because Illustrator will auto-convert and can silently lose fidelity.

## Agent-First Rules That Stay Locked
- Raw JSON payloads remain the primary interface. No drift back to bespoke flag-heavy commands.
- `tools`, `schema`, and `schema --all --llms-txt` remain the canonical discovery surfaces.
- Stable envelopes remain mandatory for every command and batch step.
- `--dry-run`, `--fields`, and NDJSON stay first-class.
- Input hardening stays strict:
  - safe-path enforcement
  - opaque identifier validation
  - no control chars
  - no encoded traversal tricks
  - no fuzzy correction for destructive or ambiguous targets
- The plugin must never become a hidden runtime requirement for standard document, artboard, layer, selection, path, text, appearance, export, or inspect flows.

## Current Coverage
- Already implemented in `internal/illustrator/schema/v0.go` and `internal/illustrator/commands/builders.go`:
  - `app.*`
  - `document.*`
  - `artboard.*`
  - `layer.*`
  - `selection.*`
  - `path.*`
  - `text.*`
  - `appearance.*`
  - `action.*`
  - `export.*`
  - `inspect.*`
- Live script validation already exists for core document creation, open/save/export, text, path, gradient, and inspection flows.

## Remaining Risks
- The native plugin is still not a loadable Illustrator `.aip`.
- Live coverage for `action.*` still needs a real fixture `.aia` set.
- `app.execute_menu` and UI-sensitive operations are still higher-risk than pure document scripting.
- The current schema only covers the original v0.1 surface, while the Adobe references expose significantly more scriptable capability.
- The current plugin build path is not yet maintainable for shipping precompiled plugins to users.

## Completion Architecture

### Script-first transport
- `ahd-illustrator` -> `internal/illustrator/host/applescript.go`
- `host` -> `do javascript`
- `do javascript` -> `internal/illustrator/commands/builders.go`
- Result -> stable JSON envelope

### Optional native transport
- `ahd-illustrator` -> AppleScript `do javascript`
- ExtendScript bridge -> `app.sendScriptMessage(...)`
- Native `.aip` -> suite-backed selector handlers
- Result -> bridge JSON -> stable CLI envelope

### Responsibility split
- Script path owns:
  - standard document automation
  - layout and art creation
  - appearance and export
  - most inspection and reporting
  - user-safe defaults
- Native path owns:
  - suite-backed operations not reliably available via scripting
  - deeper art metadata and plugin-specific inspection
  - optional performance improvements
  - future notifiers, event streaming, plugin art, and advanced suite-only features

## Phase Plan

## Phase 1: Expand Script-First Command Coverage
Goal: close the biggest feature completeness gap without waiting on native plugin work.

### New domains to add first
- `preference.get`
- `preference.set`
- `preference.delete`
- `view.info`
- `view.set_mode`
- `view.set_ruler_visibility`
- `view.set_transparency_grid_visibility`
- `matrix.identity`
- `matrix.rotation`
- `matrix.scale`
- `matrix.translation`
- `matrix.concatenate`
- `matrix.invert`
- `matrix.equals`
- `matrix.singular`
- `swatch.list`
- `swatch.create`
- `swatch.delete`
- `spot.list`
- `spot.create`
- `spot.delete`
- `style.character.list`
- `style.character.apply`
- `style.character.import`
- `style.paragraph.list`
- `style.paragraph.apply`
- `style.paragraph.import`
- `symbol.list`
- `symbol.create`
- `symbol.place`
- `symbol.break_link`
- `placed.list`
- `placed.place`
- `placed.embed`
- `placed.relink`
- `placed.trace`
- `raster.list`
- `raster.rasterize`
- `raster.trace`
- `raster.release_tracing`
- `repeat.grid.create`
- `repeat.grid.update`
- `repeat.radial.create`
- `repeat.radial.update`
- `repeat.symmetry.create`
- `repeat.symmetry.update`
- `dataset.list`
- `dataset.apply`
- `dataset.import`
- `dataset.export`
- `variable.list`
- `variable.import`
- `variable.export`
- `capture.image`
- `print.presets`
- `print.run`
- `export.gif`
- `export.png8`
- `export.tiff`
- `export.webp`
- `export.photoshop`
- `export.autocad`
- `export.eps`
- `export.fxg`
- `export.for_screens`

### Script-first hardening tasks
- Add typed nested schema support where Illustrator options are structured rather than flat.
- Add enum normalization only for low-risk fields, never for destructive targets.
- Add stronger file-format-specific validators:
  - export option compatibility
  - path extension normalization
  - artboard selector normalization
  - preset name validation
- Add command-specific live integration tests for:
  - `action.*`
  - `execute_menu`
  - style import/apply
  - placed/raster/trace flows
  - repeat flows
  - expanded export formats

### Code areas
- `internal/illustrator/schema/v0.go`
- `internal/illustrator/commands/builders.go`
- `internal/illustrator/commands/executor.go`
- `internal/illustrator/commands/live_integration_test.go`
- `internal/commonvalidate/*`
- `internal/commonschema/*`

## Phase 2: Upgrade Introspection and LLM Discovery
Goal: make the CLI discoverable enough that agents can reliably use the expanded surface.

### Work
- Generate richer `tools --llm` descriptions from schema metadata rather than hand-maintained text.
- Add grouped examples for every domain.
- Add `schema --all --llms-txt` coverage for every newly added command family.
- Add result field masks for large inspect/export responses.
- Add command tags in schema metadata:
  - `mutating`
  - `script-only`
  - `plugin-accelerated`
  - `plugin-required`
  - `ui-sensitive`

### Code areas
- `internal/commoncli/*`
- `internal/commonschema/*`
- `cmd/ahd-illustrator/main.go`
- `skills/ahd-illustrator/SKILL.md`
- `docs/illustrator/*`

## Phase 3: Turn the Native Plugin Into a Real `.aip`
Goal: replace the bridge skeleton with a loadable, script-message-capable Illustrator plugin on macOS.

### Locked build direction
- Production macOS plugin build should follow Adobe’s supported shape:
  - `.aip` bundle
  - Xcode 15.2
  - C++17
  - `Plugin.cpp`, `Suites.cpp`, and related common sources from SDK sample patterns
  - PiPL/resource generation through Adobe’s `tools/pipl/create_pipl.py`
- CMake can remain for isolated bridge logic tests, but it should not remain the sole production build path if it fights the SDK conventions.
- The best baseline is Adobe’s generic Template plug-in skeleton, not the full `ScriptMessage` sample target. `ScriptMessage` should be used for selector/message behavior only, because its sample target also pulls in panel and PlugPlug UI wiring that AHD does not need for a headless agent bridge.

### Required implementation work
- Replace the current stub entrypoint with a real SDK bootstrap:
  - `AllocatePlugin`
  - `FixupReload`
  - `PluginMain`
  - `Message` handling for `kCallerAIScriptMessage`
- Use the SDK Template sample as the target skeleton and borrow only the `AIScriptMessage` handling shape from the SDK `ScriptMessage` sample.
- Implement selector handlers:
  - `ahd.capabilities`
  - `ahd.version`
  - `ahd.exec`
  - `ahd.inspect`
- Replace echo logic with suite-backed execution and inspection where scripting is weak or unavailable.
- Add stable error mapping from native failures into CLI envelopes.
- Add a plugin probe that distinguishes:
  - plugin missing
  - plugin version mismatch
  - selector unsupported
  - plugin runtime error
- Start with the minimum suite surface required for the bridge:
  - `SPBlocks`
  - `AIUnicodeString`
  - `AIArt`
  - `AIDocument`
  - `AIUser`
- Expand into text, gradient, art-style, notifier, and plugin-group suites only after the basic bridge is live and stable.

### macOS build/package work
- Add a repo-owned mac plugin project or generator that produces:
  - `.aip`
  - valid PiPL
  - `Info.plist`
  - any required resources
- Add a build script that accepts `ILLUSTRATOR_SDK_ROOT`.
- Add an install script that copies the built `.aip` into:
  - `/Applications/Adobe Illustrator 2026/Plug-ins.localized/`
  - `~/Library/Application Support/Adobe/Adobe Illustrator CC 2026/<localeCode>/Plug-ins`
  - or a user-specified Additional Plug-ins Folder
- During development, prefer Illustrator’s Additional Plug-ins Folder preference to avoid writing directly into the app bundle on every build.

### Code areas
- `tools/illustrator/plugin-cpp/*`
- new mac build/install wrapper scripts under `tools/illustrator/plugin-cpp/` or `scripts/`
- `internal/illustrator/bridge/*`
- `cmd/ahd-illustrator/main.go`

## Phase 4: Package the Plugin for Designers
Goal: make plugin usage easy for end users and keep the SDK invisible.

### Work
- Publish precompiled macOS plugin artifacts.
- Add a one-command install path:
  - `ahd-illustrator plugin install`
  - `ahd-illustrator plugin status`
  - `ahd-illustrator plugin uninstall`
- Keep `doctor` authoritative for:
  - Illustrator discovery
  - version compatibility
  - plugin presence
  - plugin compatibility
- Document that users do not need the SDK.

### Distribution target
- CLI binary plus optional prebuilt `.aip`
- no SDK requirement for designers
- no manual Xcode build requirement for designers
- Native plugin install should be reducible to one user-facing command or installer flow; users should not need to understand PiPL, Xcode, or SDK layout.

## Phase 5: Windows Parity
Goal: prepare a real path to Windows support after macOS is solid.

### Work
- Add a Windows host adapter using the documented Illustrator automation surface.
- Build the native plugin with Visual Studio 2022 using the Windows `ScriptMessage` sample pattern.
- Keep the command schema identical across macOS and Windows wherever possible.
- Mark unsupported commands with deterministic capability metadata instead of silent omission.

### Why this is feasible
- The Windows SDK mirrors the `ScriptMessage` sample structure.
- The VBScript reference mirrors the JavaScript object model closely enough to keep the domain model shared.

## Testing Strategy

### Unit
- schema registration and aliasing
- nested option validation
- file-format-specific hardening
- deterministic envelope mapping
- plugin capability classification

### Live script integration
- scratch-doc workflows for every new high-value domain
- export fixture generation for every supported format
- `action.*` fixture validation with a checked-in test action set
- selection, trace, repeat, symbol, placed, and style flows

### Native integration
- loadable `.aip` smoke test
- `sendScriptMessage` selector coverage
- plugin version negotiation
- missing plugin and bad plugin error classification

### Packaging
- plugin install into Illustrator plugin folder
- plugin install via Additional Plug-ins Folder
- CLI `doctor` after install and after uninstall

## Definition of Done
- The CLI exposes the expanded Illustrator scripting surface through typed, discoverable, schema-backed commands.
- The majority of design automation remains script-first and works with no SDK and no plugin.
- The macOS native plugin is a real `.aip`, live-loadable in Illustrator 2026, and optional for users.
- The CLI can install, detect, and report plugin state deterministically.
- The docs and skill surface make the full command set discoverable for AI agents.
- The Windows path is explicitly designed, documented, and unblocked by the macOS architecture.

## Risk Retirement Order
1. Finish script-first command expansion for the high-value documented surfaces.
2. Add live fixtures for `action.*`, `execute_menu`, and format-specific export behavior.
3. Replace the native bridge skeleton with a real macOS `.aip`.
4. Package the plugin for non-technical users.
5. Add Windows host/plugin parity once the macOS architecture is stable.

## Repo Constraints
- Do not touch these unrelated dirty files during Illustrator work:
  - `internal/designlint/lint.go`
  - `internal/tools/catalog_llm.go`
  - `plugin/src/handlers/effect.ts`
  - `plugin/src/handlers/shape.ts`

## Immediate Next Step
- Execute Phase 1 first. It retires the biggest real gap: Adobe’s script surface is materially larger than the current CLI, and expanding it keeps the product usable without forcing plugin installation or SDK access on users.
