# npm Wrapper Investigation

Date: 2026-07-08

Goal: reduce install friction while preserving the single Go binary model.

Recommended shape:

- Publish a tiny `ahd-figma` npm package.
- On first run, detect platform and architecture.
- Download the matching GitHub Release asset for `ai-happy-design`.
- Cache under npm/user cache, not the repo.
- Forward all CLI args to the Go binary.
- Verify checksum from the release manifest before execution.
- Keep the Go binary as the only implementation; npm is a transport wrapper.

Non-goals:

- Do not reimplement CLI behavior in Node.
- Do not require npm for existing Homebrew/manual binary users.
- Do not auto-update silently during command execution.

Open work:

- Add goreleaser checksums to release artifacts.
- Decide package name ownership.
- Add smoke tests for macOS arm64/x64, Linux x64/arm64, Windows x64.
