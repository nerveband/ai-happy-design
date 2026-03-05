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

Detailed build and install steps are filled in as the bridge implementation lands.
