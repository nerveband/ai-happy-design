# Cerebras Curl Skill — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a skill that teaches Claude how to call the Cerebras API via raw curl commands for fast JSON generation and general-purpose LLM calls.

**Architecture:** Single SKILL.md file with no scripts or dependencies. Claude constructs curl commands on the fly. Works on zsh, bash, and fish.

**Tech Stack:** curl, jq (fallback: python3), Cerebras REST API (OpenAI-compatible)

---

### Task 1: Create skill directory

**Files:**
- Create: `/Users/nerveband/.config/skillshare/skills/cerebras/SKILL.md`

**Step 1: Create the directory**

Run: `mkdir -p /Users/nerveband/.config/skillshare/skills/cerebras`

**Step 2: Create the symlink in ~/.claude/skills/**

Run: `ln -sf /Users/nerveband/.config/skillshare/skills/cerebras /Users/nerveband/.claude/skills/cerebras`

**Step 3: Verify**

Run: `ls -la /Users/nerveband/.claude/skills/cerebras`
Expected: symlink pointing to skillshare directory

---

### Task 2: Write SKILL.md

**Files:**
- Create: `/Users/nerveband/.config/skillshare/skills/cerebras/SKILL.md`

**Step 1: Write the full SKILL.md**

The skill must include these sections:

1. **Frontmatter** — name, description with trigger phrases
2. **API Basics** — endpoint, auth, default model
3. **Simple Prompt** — curl one-liner to get a response
4. **Structured JSON Output** — response_format with json_schema
5. **List Models** — public endpoint, no auth needed
6. **Model Selection Guide** — when to pick which
7. **Streaming** — stream:true pattern
8. **JSON Parsing** — jq preferred, python3 fallback
9. **Gotchas** — tools+response_format conflict, token limits

Key content requirements:
- Default model: `qwen-3-235b-a22b-instruct-2507`
- API key from `$CEREBRAS_API_KEY` env var or user-supplied
- All curl examples must be copy-pasteable
- jq extraction: `.choices[0].message.content`
- python3 fallback: `python3 -c "import sys,json; print(json.loads(sys.stdin.read())['choices'][0]['message']['content'])"`
- Model list endpoint: `https://api.cerebras.ai/public/v1/models` (no auth)
- Warn: `response_format` and `tools` cannot coexist in same request

**Step 2: Verify the file exists and is well-formed**

Run: `wc -l /Users/nerveband/.config/skillshare/skills/cerebras/SKILL.md`
Expected: 150-250 lines

---

### Task 3: Test the skill

**Step 1: Verify skill appears in Claude Code skill list**

Check that `cerebras` appears when listing skills. The symlink must be in place.

**Step 2: Test a simple curl call**

Run the simple prompt pattern from the skill against the live API to verify it works:
```bash
curl -sS https://api.cerebras.ai/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $CEREBRAS_API_KEY" \
  -d '{"model":"qwen-3-235b-a22b-instruct-2507","messages":[{"role":"user","content":"Say hello in JSON format"}]}' | jq '.choices[0].message.content'
```

**Step 3: Test structured output**

Run the structured output pattern with a simple schema to verify JSON schema enforcement works.

**Step 4: Test model listing**

Run: `curl -sS https://api.cerebras.ai/public/v1/models | jq '.data[].id'`
Expected: List of model IDs

---

### Task 4: Commit

**Step 1: Stage and commit**

```bash
git add docs/plans/2026-02-17-cerebras-curl-skill-design.md
git add docs/plans/2026-02-17-cerebras-curl-skill.md
git commit -m "feat: add cerebras curl skill for fast LLM inference via curl"
```

Note: The SKILL.md lives outside this repo (in skillshare), so it won't be in this commit.
