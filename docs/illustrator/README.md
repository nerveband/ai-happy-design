# AHD Illustrator

`ahd-illustrator` is the Illustrator side of the AHD Design monorepo.

## v0.1 Target

- macOS only
- CLI only
- agent-first JSON envelopes
- AppleScript `do javascript` host control with JSX bridge support
- optional C++ plugin bridge for capability and inspection extensions

## Live Runtime Notes

- Live validation currently targets Adobe Illustrator 30.2.1 on macOS.
- `host status` and `doctor` resolve the app via bundle id and surface the installed app path plus version.
- The JSX bridge does not rely on `JSON.stringify`; Illustrator's ExtendScript runtime in this build does not provide a global `JSON` object.
- The native plugin probe remains visible in diagnostics, but the current inspect and appearance command surface is script-backed and works without the plugin installed.
- `app.info` exposes runtime discovery fields that agents can use directly, including `scriptingVersion` and the installed startup preset list.
- Multi-artboard document creation is script-backed through Illustrator's real `documents.add(...)` and `documents.addDocument(...)` surfaces rather than post-hoc placeholder logic.
- Ref/runtime mismatches and live host discoveries are tracked in [live-script-runtime-learnings.md](live-script-runtime-learnings.md).

## Without Plugin Or SDK

- End users do not need the Adobe Illustrator SDK.
- End users do not need the native AHD plugin for the current CLI surface.
- With only Illustrator installed, users can run discovery plus the full current command surface:
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
  - `page_item.*`
  - `workspace.*`
  - `preference.*`
  - `view.*`
  - `matrix.*`
  - `perspective.*`
  - `style.character.*`
  - `style.paragraph.*`
  - `style.graphic.*`
  - `swatch.*`
  - `spot.*`
  - `symbol.*`
  - `placed.*`
  - `raster.*`
  - `repeat.*`
  - `dataset.*`
  - `variable.*`
  - `trace.preset.list`
  - `capture.*`
  - `print.*`
- The SDK is only needed if a maintainer wants to build the native plugin from source.
- The plugin is only for future native-only capabilities or deeper integration beyond the current script-backed surface.
- Known script-first exclusions are documented explicitly:
  - `workspace.list` is not exposed by the scripting references.
  - `export.for_screens` exists live but is not yet deterministic enough for agent-first CLI exposure.
  - `document.write_as_library` is documented in the references, but the live 30.2.1 runtime rejected every enumerated `LibraryType` variant we tested.
  - `page_item.bring_in_perspective` is documented in the references, but the live 30.2.1 runtime rejected every plane enumeration we tested.
  - `trace.preset.store` is documented in the references, but the live 30.2.1 runtime returned a deterministic Illustrator error instead of storing the preset.

## Recommended Workflow

1. `ahd-illustrator doctor`
2. `ahd-illustrator tools --json`
3. `ahd-illustrator schema <domain.action> --json`
4. `ahd-illustrator command <domain.action> --json '{...}' --dry-run`
5. `ahd-illustrator batch --ops ops.json --dry-run`
6. Only remove `--dry-run` once the payload is stable and Illustrator is running

## Input Hardening

- Output paths are sandboxed to the current working directory.
- Opaque identifiers such as `itemId`, `layerId`, `artboardId`, `styleName`, and action names reject URL, query, fragment, and encoded traversal syntax.
- Nested payloads such as gradient stops and path points are schema-validated before execution.
- `document.new` applies cross-field validation so invalid artboard-count/layout combinations fail fast instead of hanging Illustrator.

## Start Here

1. Read [architecture.md](architecture.md)
2. Use [commands.md](commands.md) for the public CLI surface and output contract
3. Follow [plugin-build.md](plugin-build.md) if you need the plugin capability path
4. Review [release-notes-v0.1.md](release-notes-v0.1.md) for the shipped surface area

## Current Caveats

- Scratch-document validation is the recommended live test path so existing user artwork is not modified.
- All output paths are sandboxed to the current working directory unless the CLI later adds an explicit override.
- Artboard-targeted SVG export is script-backed, but Illustrator 30.2.1 expects a 1-based numeric artboard range at runtime even though the reference text describes name-based ranges. The CLI normalizes this internally from `artboardId`.
- Export and capture artifact names can drift from the requested path. The CLI normalizes TIFF, window-capture, and preset-export output paths to the actual created file before returning them.
- The native bridge skeleton is still buildable in CMake, but live Illustrator SDK wiring remains a later step for future native-only capabilities.

Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
