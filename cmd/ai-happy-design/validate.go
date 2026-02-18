package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var validateFixFlag bool

var validateCmd = &cobra.Command{
	Use:   "validate [file.json or '-' for stdin]",
	Short: "Validate batch JSON against the ai-happy-design schema",
	Long: `Checks batch JSON for common schema errors before sending to Figma.

Detects:
  - Using 'type' instead of 'command'
  - Missing 'params' field
  - Design properties at top level instead of inside params
  - Broken ${{steps.X.result.id}} references (X not defined as a step name)

Use --fix to auto-correct common issues in-place before validating:
  - Strips markdown fences (models add these even when told not to)
  - Renames "type" to "command"
  - Hoists top-level design props (x, y, color, etc.) into "params"

Exit code 0 = valid, 1 = validation errors found.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromFile := len(args) > 0 && args[0] != "-"

		var data []byte
		var err error
		if fromFile {
			data, err = os.ReadFile(args[0])
		} else {
			data, err = io.ReadAll(os.Stdin)
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		if validateFixFlag {
			fixed, fixes, fixErr := fixBatchOps(data)
			if fixErr != nil {
				fmt.Fprintf(os.Stderr, "✗ Could not parse input: %v\n", fixErr)
				os.Exit(1)
				return nil
			}

			if len(fixes) > 0 {
				fmt.Fprintf(os.Stderr, "✎ Applied %d fix(es):\n", len(fixes))
				for _, f := range fixes {
					fmt.Fprintf(os.Stderr, "  %s\n", f)
				}
				if fromFile {
					if writeErr := os.WriteFile(args[0], fixed, 0644); writeErr != nil {
						return fmt.Errorf("write error: %w", writeErr)
					}
					fmt.Fprintf(os.Stderr, "  → written back to %s\n", args[0])
				} else {
					// stdin mode: emit fixed JSON to stdout
					fmt.Println(string(fixed))
					return nil
				}
				fmt.Fprintln(os.Stderr, "")
			}
			data = fixed
		}

		errs := validateBatchOps(data)
		if len(errs) == 0 {
			fmt.Println("✓ Valid — 0 errors found")
			return nil
		}
		fmt.Fprintf(os.Stderr, "✗ %d error(s) found:\n\n", len(errs))
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "  "+e)
		}
		os.Exit(1)
		return nil
	},
}

// knownTopLevelProps are design properties that belong inside "params", not at op root.
var knownTopLevelProps = []string{
	"x", "y", "width", "height", "w", "h",
	"color", "fillColor", "bg",
	"fontSize", "fontFamily", "fontStyle", "sz", "ff", "fs",
	"parentId", "pid",
	"cornerRadius", "r",
	"layoutMode", "itemSpacing", "padding", "opacity",
	"text", "imageData", "stroke", "strokeWidth",
}

var interpolationRef = regexp.MustCompile(`\$\{\{steps\.([a-zA-Z0-9_]+)\.result\.[a-z]+\}\}`)

// stripMarkdownFences removes ```json ... ``` or ``` ... ``` wrappers that
// models add even when explicitly told not to.
func stripMarkdownFences(data []byte) []byte {
	s := strings.TrimSpace(string(data))
	if !strings.HasPrefix(s, "```") {
		return data
	}
	lines := strings.Split(s, "\n")
	var filtered []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if i == 0 && strings.HasPrefix(trimmed, "```") {
			continue // opening fence
		}
		if i == len(lines)-1 && trimmed == "```" {
			continue // closing fence
		}
		filtered = append(filtered, line)
	}
	return []byte(strings.Join(filtered, "\n"))
}

// fixBatchOps applies auto-corrections to common model output drift.
// Returns: fixed JSON bytes, list of human-readable fix descriptions, error.
func fixBatchOps(data []byte) ([]byte, []string, error) {
	data = stripMarkdownFences(data)

	// Unwrap {"ops": [...]} or any single-key dict wrapping an array
	// Models often output {"ops": [...]} instead of bare [...]
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, err
	}
	var fixes []string
	if obj, isObj := raw.(map[string]interface{}); isObj {
		// Find the first array value
		for k, v := range obj {
			if arr, isArr := v.([]interface{}); isArr {
				fixes = append(fixes, fmt.Sprintf("unwrapped dict key %q to get ops array", k))
				b, _ := json.Marshal(arr)
				data = b
				break
			}
		}
	}

	var ops []map[string]interface{}
	if err := json.Unmarshal(data, &ops); err != nil {
		return nil, fixes, err
	}
	for i, op := range ops {
		label := fmt.Sprintf("op[%d]", i)
		if name, ok := op["name"].(string); ok && name != "" {
			label = fmt.Sprintf("op[%d] %q", i, name)
		}

		// Fix "type" → "command"
		if typeVal, hasType := op["type"]; hasType {
			if _, hasCmd := op["command"]; !hasCmd {
				op["command"] = typeVal
				fixes = append(fixes, label+`: renamed "type" to "command"`)
			}
			delete(op, "type")
		}

		// Ensure params exists
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			params = map[string]interface{}{}
		}

		// Hoist known top-level design props into params
		var hoisted []string
		for _, prop := range knownTopLevelProps {
			if val, ok := op[prop]; ok {
				params[prop] = val
				delete(op, prop)
				hoisted = append(hoisted, prop)
			}
		}
		if len(hoisted) > 0 {
			fixes = append(fixes, fmt.Sprintf(`%s: moved %s into "params"`, label, strings.Join(hoisted, ", ")))
		}

		op["params"] = params
		ops[i] = op
	}

	fixed, err := json.MarshalIndent(ops, "", "  ")
	if err != nil {
		return nil, fixes, err
	}
	// Preserve ${{...}} interpolation — json.Marshal escapes < > & by default
	fixed = bytes.ReplaceAll(fixed, []byte(`\u0026`), []byte(`&`))
	fixed = bytes.ReplaceAll(fixed, []byte(`\u003c`), []byte(`<`))
	fixed = bytes.ReplaceAll(fixed, []byte(`\u003e`), []byte(`>`))
	return fixed, fixes, nil
}

func validateBatchOps(data []byte) []string {
	var ops []map[string]interface{}
	if err := json.Unmarshal(data, &ops); err != nil {
		return []string{fmt.Sprintf("Invalid JSON: %v", err)}
	}
	if len(ops) == 0 {
		return []string{"Batch is empty"}
	}

	var errs []string
	definedNames := map[string]bool{}

	for i, op := range ops {
		label := fmt.Sprintf("op[%d]", i)
		if name, ok := op["name"].(string); ok && name != "" {
			definedNames[name] = true
			label = fmt.Sprintf("op[%d] %q", i, name)
		}

		if _, hasType := op["type"]; hasType {
			errs = append(errs, label+`: use "command" not "type"`)
		}
		if cmd, ok := op["command"].(string); !ok || cmd == "" {
			errs = append(errs, label+`: missing "command" field`)
		}
		if p, hasParams := op["params"]; !hasParams {
			errs = append(errs, label+`: missing "params" field`)
		} else if _, ok := p.(map[string]interface{}); !ok {
			errs = append(errs, label+`: "params" must be an object`)
		}
		for _, prop := range knownTopLevelProps {
			if _, ok := op[prop]; ok {
				errs = append(errs, fmt.Sprintf(`%s: %q must be inside "params", not at top level`, label, prop))
			}
		}
	}

	raw, _ := json.Marshal(ops)
	for _, m := range interpolationRef.FindAllSubmatch(raw, -1) {
		refName := string(m[1])
		if !definedNames[refName] {
			errs = append(errs, fmt.Sprintf(`interpolation ${{steps.%s.result.id}} references undefined step name`, refName))
		}
	}

	return errs
}

func init() {
	validateCmd.Flags().BoolVar(&validateFixFlag, "fix", false, "Auto-fix common issues: markdown fences, type→command, top-level params")
	rootCmd.AddCommand(validateCmd)
}
