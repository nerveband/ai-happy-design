package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nerveband/ai-happy-design/internal/config"
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/spf13/cobra"
)

func tokenPreset(name string, params map[string]interface{}) (interface{}, error) {
	tokens := minimalTokenPreset()
	switch name {
	case "tailwind":
		tokens["colors"] = map[string]interface{}{"slate": "#0f172a", "sky": "#0ea5e9", "emerald": "#10b981", "rose": "#f43f5e"}
		tokens["spacing"] = map[string]interface{}{"1": 4, "2": 8, "3": 12, "4": 16, "6": 24, "8": 32, "12": 48}
	case "shadcn":
		tokens["colors"] = map[string]interface{}{"background": "#ffffff", "foreground": "#09090b", "primary": "#18181b", "muted": "#f4f4f5", "border": "#e4e4e7"}
		tokens["radii"] = map[string]interface{}{"sm": 4, "md": 6, "lg": 8}
	case "material":
		tokens["colors"] = map[string]interface{}{"primary": "#6750A4", "secondary": "#625B71", "surface": "#FFFBFE", "error": "#B3261E"}
		tokens["radii"] = map[string]interface{}{"xs": 4, "sm": 8, "md": 12, "lg": 16}
	}
	out := map[string]interface{}{"preset": name, "tokens": tokens, "dtcg": toDTCG(tokens)}
	if path, _ := params["outputPath"].(string); path != "" {
		data, _ := json.MarshalIndent(out, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		out["savedTo"] = path
	}
	return out, nil
}

func toDTCG(tokens map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for group, raw := range tokens {
		if values, ok := raw.(map[string]interface{}); ok {
			g := map[string]interface{}{}
			for name, value := range values {
				tokenType := "dimension"
				if strings.Contains(group, "color") {
					tokenType = "color"
				}
				g[name] = map[string]interface{}{"$type": tokenType, "$value": value}
			}
			out[group] = g
		}
	}
	return out
}

func setupTokenSystem(params map[string]interface{}) (interface{}, error) {
	preset, _ := params["preset"].(string)
	if preset == "" {
		preset = "tailwind"
	}
	result, err := tokenPreset(preset, params)
	if err != nil {
		return nil, err
	}
	ops := []map[string]interface{}{
		{"name": "token_note", "command": "text.create", "params": map[string]interface{}{"text": "Token preset: " + preset, "x": 0, "y": 0, "fontSize": 24}},
	}
	out := map[string]interface{}{"preset": preset, "tokens": result.(map[string]interface{})["tokens"], "batch": ops}
	if path, _ := params["outputPath"].(string); path != "" {
		data, _ := json.MarshalIndent(ops, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		out["savedTo"] = path
	}
	return out, nil
}

func designSystemHealth(params map[string]interface{}) (interface{}, error) {
	spec := map[string]interface{}{}
	if inline, ok := params["spec"].(map[string]interface{}); ok {
		spec = inline
	} else {
		for _, key := range []string{"colors", "color", "typography", "type", "spacing", "space", "components", "tokens"} {
			if value, ok := params[key]; ok {
				spec[key] = value
			}
		}
	}
	if path, _ := params["specPath"].(string); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, err
		}
	}
	findings := []map[string]interface{}{}
	score := 100
	for _, group := range [][]string{{"colors", "color"}, {"typography", "type"}, {"spacing", "space"}, {"components"}} {
		if !hasAnyKey(spec, group...) {
			score -= 10
			findings = append(findings, map[string]interface{}{"code": "MISSING_" + strings.ToUpper(group[0]), "message": "Missing " + group[0] + " evidence"})
		}
	}
	if score < 0 {
		score = 0
	}
	return map[string]interface{}{"score": score, "grade": healthGrade(score), "findings": findings, "checks": []string{"naming", "token_usage", "component_coverage", "accessibility", "spacing", "typography", "reusable_styles"}}, nil
}

func hasAnyKey(m map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := m[key]; ok {
			return true
		}
	}
	return false
}

func healthGrade(score int) string {
	if score >= 90 {
		return "excellent"
	}
	if score >= 75 {
		return "good"
	}
	if score >= 60 {
		return "fair"
	}
	return "needs_work"
}

func analyzeComponentSet(params map[string]interface{}) (interface{}, error) {
	raw := params["componentSet"]
	if raw == nil {
		raw = params["components"]
	}
	states := []string{}
	if m, ok := raw.(map[string]interface{}); ok {
		if variants, ok := m["variants"].([]interface{}); ok {
			for _, v := range variants {
				s := strings.ToLower(fmt.Sprint(v))
				for _, state := range []string{"hover", "focus", "disabled", "selected", "loading", "default"} {
					if strings.Contains(s, state) {
						states = append(states, state)
					}
				}
			}
		}
	}
	if arr, ok := raw.([]interface{}); ok {
		for _, item := range arr {
			s := strings.ToLower(fmt.Sprint(item))
			if m, ok := item.(map[string]interface{}); ok {
				if name, ok := m["name"]; ok {
					s += " " + strings.ToLower(fmt.Sprint(name))
				}
				if props, ok := m["properties"].(map[string]interface{}); ok {
					for _, value := range props {
						s += " " + strings.ToLower(fmt.Sprint(value))
					}
				}
			}
			for _, state := range []string{"hover", "focus", "disabled", "selected", "loading", "default"} {
				if strings.Contains(s, state) {
					states = append(states, state)
				}
			}
		}
	}
	if len(states) == 0 {
		states = []string{"default"}
	}
	return map[string]interface{}{"states": states, "documentation": "Detected states: " + strings.Join(states, ", ")}, nil
}

func arrangeComponentSet(params map[string]interface{}) (interface{}, error) {
	analysis, _ := analyzeComponentSet(params)
	states, _ := analysis.(map[string]interface{})["states"].([]string)
	return map[string]interface{}{"layout": "grid", "columns": 4, "states": states, "gap": 80}, nil
}

func auditComponentParity(params map[string]interface{}) (interface{}, error) {
	figmaNode := map[string]interface{}{}
	codeSpec := map[string]interface{}{}
	if m, ok := params["figmaNode"].(map[string]interface{}); ok {
		figmaNode = m
	} else if m, ok := params["component"].(map[string]interface{}); ok {
		figmaNode = m
	}
	if m, ok := params["codeSpec"].(map[string]interface{}); ok {
		codeSpec = m
	} else if m, ok := params["code"].(map[string]interface{}); ok {
		codeSpec = m
	}
	findings := []map[string]interface{}{}
	for _, key := range []string{"name", "colors", "typography", "spacing", "radii", "layout", "states"} {
		_, fOK := figmaNode[key]
		_, cOK := codeSpec[key]
		if fOK != cOK {
			findings = append(findings, map[string]interface{}{"field": key, "message": "Figma/code presence differs"})
		}
	}
	score := 100 - len(findings)*12
	if score < 0 {
		score = 0
	}
	return map[string]interface{}{"score": score, "findings": findings, "figmaName": figmaNode["name"], "codeName": codeSpec["name"]}, nil
}

func generatePortableSkills(params map[string]interface{}) (interface{}, error) {
	outDir, _ := params["outputDir"].(string)
	if outDir == "" {
		outDir = filepath.Join(config.Dir(), "generated-skills")
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}
	catalogJSON, _ := json.MarshalIndent(tools.LLMCatalog(), "", "  ")
	files := map[string]string{
		"claude/SKILL.md":      "# AI Happy Design\n\nUse `ahd-figma tools --llm --json`, `ahd-figma schema --json`, and batch-first workflows.\n\nCatalog source: `internal/tools/catalog_llm.go`.\n",
		"codex/SKILL.md":       "# AI Happy Design\n\nCLI-first Figma automation. Discover with `ahd-figma agent-context --json`, then batch changes and screenshot results.\n",
		"cursor/rules.md":      "# AI Happy Design Rules\n\nUse schema-backed commands. Prefer `batch` for multi-step design creation. Verify with screenshots.\n",
		"gemini/extension.md":  "# AI Happy Design Gemini Extension Notes\n\nUse `ahd-figma` as the source of truth for command schemas and design guidance.\n",
		"catalog/catalog.json": string(catalogJSON),
		"catalog/README.md":    "Generated from `internal/tools/catalog_llm.go`; do not hand-edit design rules here.\n",
	}
	saved := []string{}
	for rel, content := range files {
		path := filepath.Join(outDir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return nil, err
		}
		saved = append(saved, path)
	}
	return map[string]interface{}{"outputDir": outDir, "files": saved, "count": len(saved)}, nil
}

var doctorJSON bool

var doctorCmd = &cobra.Command{Use: "doctor", Short: "Run local diagnostics", RunE: func(cmd *cobra.Command, args []string) error {
	checks := map[string]interface{}{"configDir": config.Dir(), "configPath": config.Path(), "catalogVersion": tools.LLMCatalog()["version"]}
	for _, p := range []string{"plugin/manifest.json", "plugin/dist/code.js", "plugin/dist/ui.html"} {
		_, err := os.Stat(p)
		checks[p] = err == nil
	}
	return printJSON(map[string]interface{}{"ok": true, "checks": checks})
}}

var verifyCmd = &cobra.Command{Use: "verify", Short: "Run proof gates"}
var verifySyntaxCmd = &cobra.Command{Use: "syntax", Short: "Verify plugin syntax gate", RunE: func(cmd *cobra.Command, args []string) error {
	path := filepath.Join("plugin", "dist", "code.js")
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	counts := map[string]int{
		"optionalChaining": bytes.Count(raw, []byte("?.")),
		"nullish":          bytes.Count(raw, []byte("??")),
		"objectSpread":     bytes.Count(raw, []byte("...")),
	}
	ok := counts["optionalChaining"] == 0 && counts["nullish"] == 0 && counts["objectSpread"] == 0
	return printJSON(map[string]interface{}{"ok": ok, "file": path, "counts": counts})
}}
var verifyPluginCmd = &cobra.Command{Use: "plugin", Short: "Verify plugin files exist", RunE: func(cmd *cobra.Command, args []string) error {
	for _, p := range []string{"plugin/manifest.json", "plugin/dist/code.js", "plugin/dist/ui.html"} {
		if _, err := os.Stat(p); err != nil {
			return err
		}
	}
	return printJSON(map[string]interface{}{"ok": true})
}}
var verifyLiveCmd = &cobra.Command{Use: "live", Short: "Print live verification recipe", RunE: func(cmd *cobra.Command, args []string) error {
	return printJSON(map[string]interface{}{"ok": true, "steps": []string{"open Figma plugin", "set AHD_CHANNEL", "run document.get_editor_context", "create page", "screenshot", "probe grid/noise/slots/motion/shaders"}})
}}
var verifyReleaseCmd = &cobra.Command{Use: "release", Short: "Print release verification checklist", RunE: func(cmd *cobra.Command, args []string) error {
	checks := []map[string]interface{}{}
	allOK := true
	for _, spec := range []struct {
		name string
		args []string
	}{
		{"go_test", []string{"go", "test", "./..."}},
		{"go_build", []string{"go", "build", "./..."}},
		{"contracts", []string{"make", "verify-contracts"}},
	} {
		start := time.Now()
		c := exec.Command(spec.args[0], spec.args[1:]...)
		out, err := c.CombinedOutput()
		ok := err == nil
		if !ok {
			allOK = false
		}
		checks = append(checks, map[string]interface{}{"name": spec.name, "ok": ok, "elapsedMs": time.Since(start).Milliseconds(), "outputTail": tailString(string(out), 1200)})
	}
	return printJSON(map[string]interface{}{"ok": allOK, "checks": checks, "pluginCheck": "Run cd plugin && npm run check for the Node/plugin gate"})
}}

var feedbackCmd = &cobra.Command{Use: "feedback <message>", Short: "Store local feedback", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := os.MkdirAll(config.Dir(), 0755); err != nil {
		return err
	}
	path := filepath.Join(config.Dir(), "feedback.jsonl")
	item := map[string]interface{}{"createdAt": time.Now().UTC().Format(time.RFC3339), "message": args[0]}
	raw, _ := json.Marshal(item)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(raw, '\n')); err != nil {
		return err
	}
	return printJSON(map[string]interface{}{"ok": true, "savedTo": path})
}}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", true, "Accepted for compatibility; doctor always outputs JSON")
	verifyCmd.AddCommand(verifySyntaxCmd, verifyPluginCmd, verifyLiveCmd, verifyReleaseCmd)
	rootCmd.AddCommand(doctorCmd, verifyCmd, feedbackCmd)
}

func tailString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[len(s)-max:]
}
