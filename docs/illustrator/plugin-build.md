# Illustrator Plugin Build

The Illustrator bridge plugin lives under `tools/illustrator/plugin-cpp/`.

## Scope

- macOS-only for v0.1
- C++ plugin bridge for `sendScriptMessage` selectors
- selectors: `ahd.capabilities`, `ahd.exec`, `ahd.inspect`, `ahd.version`

## Build Goals

- keep the plugin optional for script-only command coverage
- fail cleanly with `PLUGIN_REQUIRED` when a selector-backed command is requested without the plugin installed
- provide a narrow JSON bridge rather than exposing plugin internals directly

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

## Selectors

- `ahd.capabilities`
- `ahd.exec`
- `ahd.inspect`
- `ahd.version`

The exported bridge function in `src/PluginMain.cpp` is designed to be the single point where Illustrator SDK script-message wiring forwards selector requests into the JSON handlers.
