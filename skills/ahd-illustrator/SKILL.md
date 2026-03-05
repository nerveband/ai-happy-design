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
- Treat `PLUGIN_REQUIRED` as a bridge capability issue, not a generic command failure.
- Do not write output paths outside the working directory unless a future CLI flag explicitly allows it.

## Notable Commands

- `ahd-illustrator tools --json`
- `ahd-illustrator schema --all --llms-txt`
- `ahd-illustrator command text.create --json '{...}' --dry-run`
- `ahd-illustrator batch --ops ops.json --dry-run`
- `ahd-illustrator host status`

## Plugin Bridge

These commands rely on the plugin bridge path:

- `inspect.tree`
- `inspect.styles`
- `inspect.bounds`
- `inspect.fonts`
- `inspect.summary`
- `appearance.set_gradient`
- `appearance.apply_graphic_style`

Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
