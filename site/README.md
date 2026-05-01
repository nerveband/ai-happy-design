# AI Happy Design Site

Static marketing and LLM documentation for `aihappydesign.com`.

This folder is now managed inside the main `ai-happy-design` monorepo. Netlify should publish the `site/` directory from the root `netlify.toml`.

Key files:

- `index.html` - landing page
- `styles.css` - site styles
- `script.js` - small interaction layer
- `llms.txt` - compact AI-agent reference
- `llms-full.txt` - generated full command reference from `ahd-figma schema --all`

Regenerate the full command reference after schema changes:

```bash
ahd-figma schema --all > site/llms-full.txt
```
