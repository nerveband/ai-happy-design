package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file.json or '-' for stdin]",
	Short: "Validate batch JSON against the ai-happy-design schema",
	Long: `Checks batch JSON for common schema errors before sending to Figma.

Detects:
  - Using 'type' instead of 'command'
  - Missing 'params' field
  - Design properties at top level instead of inside params
  - Broken ${{steps.X.result.id}} references (X not defined as a step name)

Exit code 0 = valid, 1 = validation errors found.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var data []byte
		var err error

		if len(args) == 0 || args[0] == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(args[0])
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
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

		// "type" instead of "command"
		if _, hasType := op["type"]; hasType {
			errs = append(errs, label+`: use "command" not "type"`)
		}

		// Missing "command"
		if cmd, ok := op["command"].(string); !ok || cmd == "" {
			errs = append(errs, label+`: missing "command" field`)
		}

		// Missing or wrong-type "params"
		if p, hasParams := op["params"]; !hasParams {
			errs = append(errs, label+`: missing "params" field`)
		} else if _, ok := p.(map[string]interface{}); !ok {
			errs = append(errs, label+`: "params" must be an object`)
		}

		// Design props at top level
		for _, prop := range knownTopLevelProps {
			if _, ok := op[prop]; ok {
				errs = append(errs, fmt.Sprintf(`%s: %q must be inside "params", not at top level`, label, prop))
			}
		}
	}

	// Second pass: validate all interpolation references
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
	rootCmd.AddCommand(validateCmd)
}
