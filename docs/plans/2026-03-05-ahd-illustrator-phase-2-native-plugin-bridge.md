# AHD Illustrator Phase 2: Native Plugin Bridge

## Goal
- Turn the current plugin skeleton into a real Illustrator-native `.aip` bridge.
- Keep the plugin optional for users.
- Build the plugin with the Adobe Illustrator SDK as a maintainer-only dependency.
- Preserve the script-first CLI as the primary baseline.
- This phase is implementation-only. It is not an exploration pass, build experiment, or partial bridge scaffold.
- Phase 2 is complete only when the native plugin path is production-real: buildable, installable, loadable in Illustrator, and fully integrated with the CLI.

## Why Phase 2 Comes After Phase 1
- Most near-term CLI value can be delivered through scripting.
- The SDK review shows the native plugin is a materially different build problem:
  - real `.aip` bundle
  - valid PiPL
  - Adobe `Plugin` / `Suites` bootstrap
  - `PluginMain`
  - bundle install flow
- Users should not be exposed to SDK complexity.

## Native Direction Locked By The SDK
- Use Adobe's generic Template plugin skeleton as the build/lifecycle baseline.
- Borrow only the `AIScriptMessage` message handling pattern from the `ScriptMessage` sample.
- Do not use CEP/panel infrastructure as the main CLI transport.
- Keep the transport model:
  - `ahd-illustrator CLI`
  - AppleScript `do javascript`
  - ExtendScript `app.sendScriptMessage(...)`
  - `AHDIllustrator.aip`

## Required Native Work

### Build and packaging
- Replace the current generic shared library target with a real macOS `.aip`.
- Add proper SDK bootstrap:
  - `AllocatePlugin`
  - `FixupReload`
  - `PluginMain`
  - `StartupPlugin`
  - `ShutdownPlugin`
  - `Message`
- Add PiPL/resource generation using Adobe tooling.
- Add `Info.plist` and bundle metadata.
- Add a repo-owned build script that accepts `ILLUSTRATOR_SDK_ROOT`.
- Add install/uninstall helpers for:
  - `/Applications/Adobe Illustrator 2026/Plug-ins.localized/`
  - user-specified Additional Plug-ins Folder

### Selector bridge
- Keep selector contract:
  - `ahd.capabilities`
  - `ahd.version`
  - `ahd.exec`
  - `ahd.inspect`
- Replace stub/echo bridge logic with suite-backed implementations.
- Add clear error mapping for:
  - plugin missing
  - version mismatch
  - unsupported selector
  - native runtime failure
- `ahd.exec` and `ahd.inspect` must be real native implementations, not pass-through placeholders or string echoes.

### Minimum suite surface first
- `SPBlocks`
- `AIUnicodeString`
- `AIArt`
- `AIDocument`
- `AIUser`

### Likely suite expansion after the base bridge is live
- text-native inspection/editing
- gradient/path-style/art-style inspection and mutation
- notifier/undo support
- plugin-group or annotator support if future native-only UX needs it

## Feature-Complete Definition
Phase 2 is feature-complete only when all of the following are true:
- The plugin builds as a real Illustrator `.aip` from the repo using the local SDK.
- The plugin installs cleanly into Illustrator development and user-facing locations.
- Illustrator loads the plugin successfully in a live session.
- The CLI can probe plugin presence, compatibility, and version deterministically.
- `app.sendScriptMessage("AHDIllustrator", selector, payload)` returns real structured JSON for all supported selectors.
- `ahd.capabilities` reports actual native capability coverage.
- `ahd.version` reports real plugin identity/version data.
- `ahd.inspect` performs real suite-backed inspection.
- `ahd.exec` performs real suite-backed execution for native-owned functionality.
- User-facing install/status/uninstall flows exist without requiring end users to know anything about the SDK.
- The script-first baseline continues to work when the plugin is absent.

## Scope

### In scope
- macOS plugin build
- plugin probing and diagnostics
- install/status/uninstall UX
- native selector execution for cases where scripting is weak or insufficient
- production packaging decisions for delivering the prebuilt plugin to users

### Out of scope
- Windows host/runtime implementation
- CEP panel UX
- making the plugin mandatory for standard users
- leaving selector handlers as protocol-only echoes

## Files Likely To Change
- `tools/illustrator/plugin-cpp/*`
- possibly add repo-owned Xcode project or generator under `tools/illustrator/plugin-cpp/`
- `internal/illustrator/bridge/*`
- `cmd/ahd-illustrator/main.go`
- `docs/illustrator/plugin-build.md`
- `docs/illustrator/architecture.md`
- `README.md`

## Constraints
- Do not edit these unrelated dirty files:
  - `internal/designlint/lint.go`
  - `internal/tools/catalog_llm.go`
  - `plugin/src/handlers/effect.ts`
  - `plugin/src/handlers/shape.ts`
- Do not regress the current script-first path.
- Do not require the SDK for end users.
- Do not make the plugin a hidden dependency for current core commands.

## Acceptance Criteria
- The plugin builds as a real `.aip` with the Illustrator SDK.
- Illustrator loads the plugin successfully.
- `ahd-illustrator doctor` can distinguish:
  - not installed
  - installed but incompatible
  - installed and reachable
- `app.sendScriptMessage("AHDIllustrator", ...)` returns deterministic JSON responses through the CLI bridge.
- The plugin can be installed by a user without requiring the SDK.

## Completion Gate
Do not mark Phase 2 done if any of the following remain true:
- the repo still only builds a generic shared library
- the plugin is not a real `.aip`
- the plugin cannot be installed and loaded live in Illustrator
- selector handlers are still stubbed or echo-only
- install/status/uninstall UX is still manual-maintainer-only
- the CLI cannot distinguish plugin-missing from plugin-broken states

## Suggested Work Order
1. Replace the build skeleton with Adobe Template-based plugin lifecycle.
2. Wire `AIScriptMessage` message handling into the existing bridge JSON core.
3. Implement `ahd.version` and `ahd.capabilities`.
4. Implement minimal suite-backed `ahd.inspect`.
5. Implement minimal suite-backed `ahd.exec`.
6. Add plugin build/install/status CLI commands or equivalent scripts.
7. Add live integration coverage with Illustrator running and plugin installed.
8. Document prebuilt distribution flow for users.

## AI Coder Prompt
```text
You are in implementation mode in the ai-happy-design monorepo root.

Work only in:
/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design

Implement Phase 2 from:
/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/docs/plans/2026-03-05-ahd-illustrator-phase-2-native-plugin-bridge.md

Context:
- The repo already has a script-first ahd-illustrator CLI.
- The local Adobe Illustrator 2026 SDK is available under Downloads/adobe-illustrator.
- The current tools/illustrator/plugin-cpp target is only a stub shared library, not a real Illustrator plugin.
- Use Adobe’s Template plugin skeleton as the build/lifecycle baseline.
- Borrow only the AIScriptMessage bridge behavior from Adobe’s ScriptMessage sample.
- Keep the plugin optional for users and preserve the script-first CLI baseline.

Requirements:
1. Implement the real macOS native plugin path end to end without pausing unless blocked.
2. Do not edit unrelated dirty files:
   - internal/designlint/lint.go
   - internal/tools/catalog_llm.go
   - plugin/src/handlers/effect.ts
   - plugin/src/handlers/shape.ts
3. Replace the generic shared library build with a real `.aip` plugin build.
4. Add the proper Illustrator SDK plugin lifecycle and PiPL/resource packaging.
5. Wire the existing JSON bridge logic into `kCallerAIScriptMessage`.
6. Keep selectors:
   - ahd.capabilities
   - ahd.version
   - ahd.exec
   - ahd.inspect
7. Add plugin diagnostics/install docs and any CLI/plugin helper commands needed.
8. Run relevant build/test/live checks and report pass/fail clearly.
9. At the end, print:
   - files changed
   - build/test results
   - plugin status
   - remaining risks

Implementation bar:
- Do not fork Adobe’s sample UI/panel stack unless absolutely necessary.
- Keep the plugin bridge deterministic and machine-readable.
- Treat the SDK as a maintainer-only build input, not a user runtime requirement.
- Do not degrade the existing script-first command path while adding the native bridge.
```
