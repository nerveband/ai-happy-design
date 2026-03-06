# AHD Illustrator

Use this skill when an agent needs to inspect, validate, or automate Adobe Illustrator through `ahd-illustrator`.

## Workflow

1. Run `ahd-illustrator doctor`
2. Run `ahd-illustrator tools --json`
3. Run `ahd-illustrator schema <domain.action> --json`
4. Build the payload as raw JSON
5. Run `ahd-illustrator command <domain.action> --json '<payload>' --dry-run`
6. For multi-step work, use `ahd-illustrator batch --ops ops.json --dry-run`
7. Only remove `--dry-run` once validation is clean and Illustrator is running

## Rules

- Default to raw JSON payloads, not bespoke flag sets.
- Prefer `schema` over guessing field names.
- Keep output constrained with `--fields` or `--output ndjson` when the calling agent has tight context limits.
- Do not write output paths outside the working directory unless a future CLI flag explicitly allows it.
- Treat `PLUGIN_REQUIRED` as a diagnostic bridge capability issue, not a blocker for the current script-first Illustrator command surface.
- Use scratch documents and disposable export paths when validating live Illustrator flows.

## Notable Commands

- `ahd-illustrator tools --json`
- `ahd-illustrator schema --all --llms-txt`
- `ahd-illustrator command text.create --json '{...}' --dry-run`
- `ahd-illustrator batch --ops ops.json --dry-run`
- `ahd-illustrator host status`

## Script-First Surface

The current Phase 1 command surface is script-backed and works without the native plugin or SDK, including:

- `inspect.*`
- `appearance.*`
- `page_item.*`
- `workspace.*`
- `preference.*`
- `view.*`
- `matrix.*`
- `perspective.*`
- `style.*`
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

Known script-first exclusions are:

- `workspace.list`
- `export.for_screens`

Experimental script-first commands are still shipped:

- `document.write_as_library`
- `page_item.bring_in_perspective`
- `trace.preset.store`

These return `EXPERIMENTAL_COMMAND` warnings and may still fail on Illustrator 30.2.1 even though they are exposed through `tools` and `schema`.

Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
