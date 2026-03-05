# Illustrator Plugin Build

The Illustrator bridge plugin lives under `tools/illustrator/plugin-cpp/`.

## Scope

- macOS-only for v0.1
- C++ plugin bridge for `sendScriptMessage` selectors
- selectors: `ahd.capabilities`, `ahd.exec`, `ahd.inspect`, `ahd.version`
- live validation target: Adobe Illustrator 30.2.1 on macOS

## Build Goals

- keep the plugin optional for script-only command coverage
- fail cleanly with `PLUGIN_REQUIRED` when a selector-backed command is requested without the plugin installed
- provide a narrow JSON bridge rather than exposing plugin internals directly

## Current Live State

- Script-backed commands are live-validated without the plugin.
- The current inspect, gradient, and graphic-style command surface is script-backed and does not require the plugin.
- `host status` and `doctor` still probe `sendScriptMessage` so future native-only capabilities can be diagnosed cleanly.
- The Illustrator app path currently resolves to `/Applications/Adobe Illustrator 2026/Adobe Illustrator.app`.
- The local SDK bundle on this machine is `~/Downloads/adobe-illustrator/mac Adobe Illustrator 2026 SDK`.

## Local CMake Build

```bash
cmake -S tools/illustrator/plugin-cpp -B /tmp/ahd-illustrator-plugin-build
cmake --build /tmp/ahd-illustrator-plugin-build
```

This builds the standalone bridge skeleton as `libahd_illustrator_plugin_bridge`.

## With the Illustrator SDK

The Adobe Illustrator 2026 SDK README currently calls for:

- Xcode 15.2
- the `macosx` SDK
- macOS 12.0 deployment target

The SDK itself is Xcode-oriented rather than CMake-oriented. Adobe ships sample `.xcodeproj` and `.vcxproj` files, not upstream `CMakeLists.txt`.

When you have the Adobe Illustrator SDK locally:

```bash
cmake -S tools/illustrator/plugin-cpp \
  -B /tmp/ahd-illustrator-plugin-build \
  -DAHD_ILLUSTRATOR_USE_SDK=ON \
  -DILLUSTRATOR_SDK_ROOT=/path/to/IllustratorSDK
cmake --build /tmp/ahd-illustrator-plugin-build
```

To turn the current bridge skeleton into a real Illustrator plugin, the repo still needs:

- a proper `.aip` bundle target instead of only a generic shared library
- `PluginMain` bootstrap wiring for Illustrator's plugin entrypoint
- PiPL generation and resource packaging
- SDK include search paths matching Adobe's shared sample config:
  - `illustratorapi/ate`
  - `illustratorapi/illustrator`
  - `illustratorapi/illustrator/actions`
  - `illustratorapi/pica_sp`
  - `samplecode/common/includes`
- install packaging for Illustrator's plug-in folder

The closest Adobe sample for `sendScriptMessage` is `samplecode/ScriptMessage/`. That sample shows the real message path through `kCallerAIScriptMessage` into plugin-specific handlers.

After wiring the SDK-specific plugin target, install the built plugin into Illustrator's non-Adobe plug-ins directory, restart Illustrator, then confirm the bridge with this probe:

```applescript
tell application "Adobe Illustrator" to do javascript "app.sendScriptMessage(\"AHDIllustrator\", \"ahd.version\", \"{}\")"
```

`ahd-illustrator host status` and `ahd-illustrator doctor` should then report a reachable plugin probe instead of `PLUGIN_REQUIRED`.

On this machine, Illustrator's plug-ins folder includes a `Plugins.txt` note saying third-party plug-ins belong in:

```text
/Applications/Adobe Illustrator 2026/Plug-ins.localized/
```

## Selectors

- `ahd.capabilities`
- `ahd.exec`
- `ahd.inspect`
- `ahd.version`

The exported bridge function in `src/PluginMain.cpp` is designed to be the single point where Illustrator SDK script-message wiring forwards selector requests into the JSON handlers.
