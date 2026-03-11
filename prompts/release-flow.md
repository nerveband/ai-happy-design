# Release Flow

Step-by-step release process for ai-happy-design.

## 1. Pre-Flight
```bash
# Build and test everything
go build ./...
go test ./...
cd plugin && npm run check && cd ..
make deploy

# Review prompts/pre-release-check.md
```

## 2. Version Decision
- **Patch** (v1.2.3 -> v1.2.4): Bug fixes, documentation updates
- **Minor** (v1.2.3 -> v1.3.0): New features, new commands, new flags
- **Major** (v1.2.3 -> v2.0.0): Breaking changes to CLI interface or protocol

Check latest: `git tag --sort=-v:refname | head -1`

## 3. Commit & Tag
```bash
git add -A
git commit -m "release: v<version> -- <summary>"
git tag v<version>
git push origin main
git push origin v<version>
```

## 4. Build Release
```bash
GITHUB_TOKEN=$(gh auth token) goreleaser release --clean
```

## 5. Verify
```bash
ai-happy-design upgrade
ai-happy-design --version
ai-happy-design tools --llm --json | head -5
```

## 6. Post-Release
```bash
# Update skill if features changed
cp skills/ai-happy-design/SKILL.md ~/.claude/skills/ai-happy-design/SKILL.md
skillshare sync

# Update site documentation if needed
cd /Users/nerveband/Documents/GitHub/ai-happy-design-site
netlify deploy --prod --dir=.
```
