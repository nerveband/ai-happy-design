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

## Recommended Workflow

1. `ahd-illustrator doctor`
2. `ahd-illustrator tools --json`
3. `ahd-illustrator schema <domain.action> --json`
4. `ahd-illustrator command <domain.action> --json '{...}' --dry-run`
5. `ahd-illustrator batch --ops ops.json --dry-run`
6. Only remove `--dry-run` once the payload is stable and Illustrator is running

## Start Here

1. Read [architecture.md](architecture.md)
2. Use [commands.md](commands.md) for the public CLI surface and output contract
3. Follow [plugin-build.md](plugin-build.md) if you need the plugin capability path
4. Review [release-notes-v0.1.md](release-notes-v0.1.md) for the shipped surface area

## Current Caveats

- Scratch-document validation is the recommended live test path so existing user artwork is not modified.
- All output paths are sandboxed to the current working directory unless the CLI later adds an explicit override.
- The native bridge skeleton is still buildable in CMake, but live Illustrator SDK wiring remains a later step for future native-only capabilities.

Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
