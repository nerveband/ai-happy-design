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
- Selector-backed commands currently fail with `PLUGIN_REQUIRED` on this machine because `sendScriptMessage` is not yet backed by an installed AHD bridge plugin.
- The Illustrator app path currently resolves to `/Applications/Adobe Illustrator 2026/Adobe Illustrator.app`.

## Local CMake Build

```bash
cmake -S tools/illustrator/plugin-cpp -B /tmp/ahd-illustrator-plugin-build
cmake --build /tmp/ahd-illustrator-plugin-build
```

This builds the standalone bridge skeleton as `libahd_illustrator_plugin_bridge`.

## With the Illustrator SDK

When you have the Adobe Illustrator SDK locally:

```bash
cmake -S tools/illustrator/plugin-cpp \
  -B /tmp/ahd-illustrator-plugin-build \
  -DAHD_ILLUSTRATOR_USE_SDK=ON \
  -DILLUSTRATOR_SDK_ROOT=/path/to/IllustratorSDK
cmake --build /tmp/ahd-illustrator-plugin-build
```

After wiring the SDK-specific plugin target, install the built plugin into the Illustrator 2026 plug-ins directory used by your local SDK/toolchain, restart Illustrator, then confirm the bridge with this probe:

```applescript
tell application "Adobe Illustrator" to do javascript "app.sendScriptMessage(\"AHDIllustrator\", \"ahd.version\", \"{}\")"
```

`ahd-illustrator host status` and `ahd-illustrator doctor` should then report a reachable plugin probe instead of `PLUGIN_REQUIRED`.

## Selectors

- `ahd.capabilities`
- `ahd.exec`
- `ahd.inspect`
- `ahd.version`

The exported bridge function in `src/PluginMain.cpp` is designed to be the single point where Illustrator SDK script-message wiring forwards selector requests into the JSON handlers.
