package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
	"github.com/nerveband/ai-happy-design/internal/validate"
	"github.com/spf13/cobra"
)

var validateFixFlag bool

var validateCmd = &cobra.Command{
	Use:   "validate [file.json or operations-json or '-' for stdin]",
	Short: "Dry-run: validate batch JSON without executing (structural + schema checks)",
	Long: `Runs the full validation pipeline without sending any commands to Figma.

Detects:
  - Using 'type' instead of 'command'
  - Missing 'params' field
  - Design properties at top level instead of inside params
  - Broken ${{steps.X.result.id}} references
  - Unknown commands (with fuzzy suggestions)
  - Invalid enum values (with auto-correction)
  - Out-of-bounds numbers (with clamping)
  - Named CSS colors (converted to hex)
  - Missing dependencies (e.g. lineHeightUnit for lineHeight)

Use --fix to auto-correct common issues in-place before validating.

Exit code 0 = valid, 1 = validation errors found.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromFile := len(args) > 0 && args[0] != "-"

		var data []byte
		var err error
		if fromFile {
			data, err = os.ReadFile(args[0])
		} else if len(args) > 0 && args[0] == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else if len(args) > 0 {
			// Inline JSON
			data = []byte(args[0])
		} else {
			data, err = io.ReadAll(os.Stdin)
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		if validateFixFlag {
			fixed, fixes, fixErr := batchutil.FixBatchOps(data)
			if fixErr != nil {
				fmt.Fprintf(os.Stderr, "Could not parse input: %v\n", fixErr)
				os.Exit(1)
				return nil
			}

			if len(fixes) > 0 {
				fmt.Fprintf(os.Stderr, "Applied %d fix(es):\n", len(fixes))
				for _, f := range fixes {
					fmt.Fprintf(os.Stderr, "  %s\n", f)
				}
				if fromFile {
					if writeErr := os.WriteFile(args[0], fixed, 0644); writeErr != nil {
						return fmt.Errorf("write error: %w", writeErr)
					}
					fmt.Fprintf(os.Stderr, "  written back to %s\n", args[0])
				} else {
					fmt.Println(string(fixed))
					return nil
				}
				fmt.Fprintln(os.Stderr, "")
			}
			data = fixed
		}

		data = unwrapOperationsPayload(data)

		// Structural validation (existing checks)
		structuralErrs := validateBatchOps(data)

		// Schema validation (new)
		var ops []map[string]interface{}
		if jsonErr := json.Unmarshal(data, &ops); jsonErr != nil {
			out := map[string]interface{}{
				"ok":     false,
				"errors": []string{fmt.Sprintf("Invalid JSON: %v", jsonErr)},
			}
			j, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(j))
			os.Exit(1)
			return nil
		}

		// Convert to validation format
		validationOps := make([]map[string]interface{}, len(ops))
		for i, op := range ops {
			validationOps[i] = map[string]interface{}{
				"command": op["command"],
				"params":  op["params"],
				"name":    op["name"],
			}
		}

		schemaResult := validate.ValidateBatch(validationOps)

		out := map[string]interface{}{
			"ok":    len(structuralErrs) == 0 && schemaResult.Blocked == 0,
			"total": len(ops),
		}

		if len(structuralErrs) > 0 {
			out["structuralErrors"] = structuralErrs
		}

		out["preValidation"] = map[string]interface{}{
			"schema": map[string]interface{}{
				"warnings": schemaResult.Warnings,
				"fixed":    schemaResult.Fixed,
				"blocked":  schemaResult.Blocked,
			},
		}

		j, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(j))

		if len(structuralErrs) > 0 || schemaResult.Blocked > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func unwrapOperationsPayload(data []byte) []byte {
	var wrapped struct {
		Operations json.RawMessage `json:"operations"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Operations) > 0 {
		return wrapped.Operations
	}
	return data
}

var interpolationRef = regexp.MustCompile(`\$\{\{steps\.([a-zA-Z0-9_]+)\.result\.[a-z]+\}\}`)

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
		for _, prop := range batchutil.KnownTopLevelProps {
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
	validateCmd.Flags().BoolVar(&validateFixFlag, "fix", false, "Auto-fix common issues: markdown fences, type->command, top-level params")
	rootCmd.AddCommand(validateCmd)
}
