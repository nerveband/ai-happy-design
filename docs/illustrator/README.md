# AHD Illustrator

`ahd-illustrator` is the Illustrator side of the AHD Design monorepo.

## v0.1 Target

- macOS only
- CLI only
- agent-first JSON envelopes
- AppleScript `do javascript` host control with JSX bridge support
- optional C++ plugin bridge for capability and inspection extensions

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

- `inspect.*`, `appearance.set_gradient`, and `appearance.apply_graphic_style` use the plugin bridge path.
- All output paths are sandboxed to the current working directory unless the CLI later adds an explicit override.
- The bridge skeleton is buildable in CMake, but live Illustrator SDK wiring still depends on local Adobe SDK installation.

Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
