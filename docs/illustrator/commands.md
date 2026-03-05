# Illustrator CLI Surface

## Core Commands

```text
ahd-illustrator tools [--json] [--llm]
ahd-illustrator schema [domain.action] [--json] [--all --llms-txt]
ahd-illustrator command <domain.action> --json '<payload>' [--dry-run] [--fields '<mask>'] [--output json|ndjson|text]
ahd-illustrator batch --ops <file|json> [--dry-run] [--strict] [--output json|ndjson]
ahd-illustrator host status|open|quit
ahd-illustrator doctor
ahd-illustrator examples [category]
```

## v0.1 Domains

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

## Stable Envelope

```json
{
  "ok": true,
  "requestId": "uuid",
  "command": "text.create",
  "result": {},
  "warnings": [],
  "timingMs": 0
}
```

Errors and batch responses follow the same top-level stability rules and remain machine-first by default.
