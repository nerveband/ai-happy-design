package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/commonvalidate"
	"github.com/nerveband/ai-happy-design/internal/illustrator/commands"
	illustratorhost "github.com/nerveband/ai-happy-design/internal/illustrator/host"
	_ "github.com/nerveband/ai-happy-design/internal/illustrator/schema"
	illustratorvalidate "github.com/nerveband/ai-happy-design/internal/illustrator/validate"
)

var version = "0.0.0-dev"

var (
	toolsJSON     bool
	toolsLLM      bool
	schemaJSON    bool
	schemaAll     bool
	schemaLLMSTxt bool
	commandJSON   string
	commandDryRun bool
	commandFields string
	commandOutput string
	batchOps      string
	batchDryRun   bool
	batchStrict   bool
	batchOutput   string
)

var rootCmd = &cobra.Command{
	Use:   "ahd-illustrator",
	Short: "AHD Illustrator - agent-first CLI for Adobe Illustrator",
	Long: `A macOS-first Illustrator CLI designed for AI agents.

Discovery-first workflow:
  1) ahd-illustrator tools --json
  2) ahd-illustrator schema <domain.action> --json
  3) ahd-illustrator command <domain.action> --json '{...}' --dry-run`,
	Version: version,
}

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List available Illustrator tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		all := commonschema.All()
		if toolsJSON {
			return commoncli.WriteJSON(all)
		}
		if toolsLLM {
			for _, schema := range all {
				fmt.Printf("%s\t%s\n", schema.Name, schema.Description)
			}
			return nil
		}
		for _, schema := range all {
			fmt.Printf("%-28s %s\n", schema.Name, schema.Description)
		}
		return nil
	},
}

var schemaCmd = &cobra.Command{
	Use:   "schema [domain.action]",
	Short: "Print Illustrator command schemas",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if schemaAll && schemaLLMSTxt {
			fmt.Println(commonschema.LLMSDoc("AHD Illustrator", "ahd-illustrator"))
			return nil
		}
		if len(args) == 0 {
			if schemaJSON {
				return commoncli.WriteJSON(commonschema.All())
			}
			for _, item := range commonschema.All() {
				fmt.Printf("%-28s %s\n", item.Name, item.Description)
			}
			return nil
		}
		schema := commonschema.Lookup(args[0])
		if schema == nil {
			return commoncli.WriteJSON(commoncli.ErrorEnvelope(args[0], "UNSUPPORTED_COMMAND", "unknown command", nil, false, nil, time.Now()))
		}
		if schemaJSON {
			return commoncli.WriteJSON(schema)
		}
		fmt.Printf("%s\n%s\n", schema.Name, schema.Description)
		for _, param := range schema.Params {
			required := ""
			if param.Required {
				required = " required"
			}
			fmt.Printf("- %s (%s%s): %s\n", param.Name, param.Type, required, param.Description)
		}
		return nil
	},
}

var commandCmd = &cobra.Command{
	Use:   "command <domain.action>",
	Short: "Execute a single Illustrator command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		commandName := args[0]
		params, err := parseObject(commandJSON)
		if err != nil {
			envelope := commoncli.ErrorEnvelope(commandName, "VALIDATION_ERROR", err.Error(), nil, false, nil, started)
			return writeEnvelope(commandOutput, commandFields, envelope)
		}

		schema := commonschema.Lookup(commandName)
		validation := illustratorvalidate.ValidateCommand(schema, params)
		warnings := toWarnings(validation.Warnings)
		if len(validation.Errors) > 0 {
			envelope := commoncli.ErrorEnvelope(commandName, "VALIDATION_ERROR", "command payload failed validation", map[string]any{
				"errors": validation.Errors,
				"params": params,
			}, false, warnings, started)
			return writeEnvelope(commandOutput, commandFields, envelope)
		}

		executor := commands.NewExecutor()
		result, extraWarnings, execErr := executor.Execute(commands.Request{
			Command: schema,
			Params:  validation.Params,
			DryRun:  commandDryRun,
		})
		warnings = append(warnings, extraWarnings...)
		if execErr != nil {
			envelope := commoncli.ErrorEnvelope(commandName, execErr.Code, execErr.Message, execErr.Details, false, warnings, started)
			return writeEnvelope(commandOutput, commandFields, envelope)
		}

		envelope := commoncli.SuccessEnvelope(commandName, result, warnings, started)
		return writeEnvelope(commandOutput, commandFields, envelope)
	},
}

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Validate or execute a batch of Illustrator commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		ops, err := parseOps(batchOps)
		if err != nil {
			return writeBatch(batchOutput, commoncli.BatchFailure("VALIDATION_ERROR", err.Error(), nil, false, started))
		}

		executor := commands.NewExecutor()
		steps := make([]commoncli.BatchStep, 0, len(ops))
		for index, op := range ops {
			resolvedOp, resolveErr := commands.ResolveBatchOp(op, steps)
			if resolveErr != nil {
				steps = append(steps, commoncli.BatchStep{
					Index:   index,
					Name:    op.Name,
					Command: op.Command,
					OK:      false,
					Error: map[string]any{
						"code":    "VALIDATION_ERROR",
						"message": resolveErr.Error(),
					},
				})
				if batchStrict {
					break
				}
				continue
			}

			schema := commonschema.Lookup(resolvedOp.Command)
			validation := illustratorvalidate.ValidateCommand(schema, resolvedOp.Params)
			warnings := toWarnings(validation.Warnings)
			if len(validation.Errors) > 0 {
				steps = append(steps, commoncli.BatchStep{
					Index:    index,
					Name:     resolvedOp.Name,
					Command:  resolvedOp.Command,
					OK:       false,
					Warnings: warnings,
					Error:    validation.Errors,
				})
				if batchStrict {
					break
				}
				continue
			}

			result, extraWarnings, execErr := executor.Execute(commands.Request{
				Command: schema,
				Params:  validation.Params,
				DryRun:  batchDryRun,
			})
			warnings = append(warnings, extraWarnings...)
			step := commoncli.BatchStep{
				Index:    index,
				Name:     resolvedOp.Name,
				Command:  resolvedOp.Command,
				OK:       execErr == nil,
				Result:   result,
				Warnings: warnings,
			}
			if execErr != nil {
				step.Error = execErr
				if batchStrict {
					steps = append(steps, step)
					break
				}
			}
			steps = append(steps, step)
		}

		return writeBatch(batchOutput, commoncli.BatchSuccess(steps, started))
	},
}

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Interact with the local Illustrator host",
}

var hostStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Report basic host readiness",
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		status := illustratorhost.NewAdapter().Status()
		return writeEnvelope("json", "", commoncli.SuccessEnvelope("host.status", status, nil, started))
	},
}

var hostOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open Illustrator (wired in the host bridge phase)",
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		adapter := illustratorhost.NewAdapter()
		if err := adapter.Open(); err != nil {
			envelope := commoncli.ErrorEnvelope("host.open", "HOST_EXEC_ERROR", err.Error(), nil, false, nil, started)
			return writeEnvelope("json", "", envelope)
		}
		envelope := commoncli.SuccessEnvelope("host.open", map[string]any{"opened": true}, nil, started)
		return writeEnvelope("json", "", envelope)
	},
}

var hostQuitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Quit Illustrator (wired in the host bridge phase)",
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		adapter := illustratorhost.NewAdapter()
		if err := adapter.Quit(); err != nil {
			envelope := commoncli.ErrorEnvelope("host.quit", "HOST_EXEC_ERROR", err.Error(), nil, false, nil, started)
			return writeEnvelope("json", "", envelope)
		}
		envelope := commoncli.SuccessEnvelope("host.quit", map[string]any{"quit": true}, nil, started)
		return writeEnvelope("json", "", envelope)
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Inspect local prerequisites for ahd-illustrator",
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		status := illustratorhost.NewAdapter().Status()
		result := map[string]any{
			"platform":  runtime.GOOS,
			"arch":      runtime.GOARCH,
			"supported": status.Supported,
			"host":      status,
			"schemas":   len(commonschema.All()),
			"domains":   commonschema.Domains(),
		}
		return writeEnvelope("json", "", commoncli.SuccessEnvelope("doctor", result, nil, started))
	},
}

var examplesCmd = &cobra.Command{
	Use:   "examples [category]",
	Short: "Print example batch payloads",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		category := "poster"
		if len(args) == 1 {
			category = strings.TrimSpace(args[0])
		}
		ops := exampleOps(category)
		if ops == nil {
			return commoncli.WriteJSON(commoncli.ErrorEnvelope("examples", "VALIDATION_ERROR", "unknown example category", map[string]any{"category": category}, false, nil, time.Now()))
		}
		return commoncli.WriteJSON(ops)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	toolsCmd.Flags().BoolVar(&toolsJSON, "json", false, "Output tool schemas as JSON")
	toolsCmd.Flags().BoolVar(&toolsLLM, "llm", false, "Print terse tool descriptions for LLM prompts")

	schemaCmd.Flags().BoolVar(&schemaJSON, "json", false, "Output schema as JSON")
	schemaCmd.Flags().BoolVar(&schemaAll, "all", false, "Print all schemas")
	schemaCmd.Flags().BoolVar(&schemaLLMSTxt, "llms-txt", false, "Print llms.txt style output")

	commandCmd.Flags().StringVar(&commandJSON, "json", "{}", "JSON payload for the command")
	commandCmd.Flags().BoolVar(&commandDryRun, "dry-run", false, "Validate without executing")
	commandCmd.Flags().StringVar(&commandFields, "fields", "", "Comma-separated top-level envelope fields to keep")
	commandCmd.Flags().StringVar(&commandOutput, "output", "json", "Output format: json|ndjson|text")

	batchCmd.Flags().StringVar(&batchOps, "ops", "", "Batch JSON or a path to a JSON file")
	batchCmd.Flags().BoolVar(&batchDryRun, "dry-run", false, "Validate without executing")
	batchCmd.Flags().BoolVar(&batchStrict, "strict", false, "Stop on the first failing step")
	batchCmd.Flags().StringVar(&batchOutput, "output", "json", "Output format: json|ndjson")
	_ = batchCmd.MarkFlagRequired("ops")

	hostCmd.AddCommand(hostStatusCmd, hostOpenCmd, hostQuitCmd)
	rootCmd.AddCommand(toolsCmd, schemaCmd, commandCmd, batchCmd, hostCmd, doctorCmd, examplesCmd)
}

func parseObject(payload string) (map[string]any, error) {
	trimmed := strings.TrimSpace(payload)
	if trimmed == "" {
		return map[string]any{}, nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(trimmed), &out); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}

func parseOps(raw string) ([]commoncli.BatchOp, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, fmt.Errorf("batch payload is empty")
	}
	data := []byte(trimmed)
	if !strings.HasPrefix(trimmed, "[") {
		fileData, err := os.ReadFile(trimmed)
		if err != nil {
			return nil, fmt.Errorf("read batch ops: %w", err)
		}
		data = fileData
	}
	var ops []commoncli.BatchOp
	if err := json.Unmarshal(data, &ops); err != nil {
		return nil, fmt.Errorf("invalid batch JSON: %w", err)
	}
	for i := range ops {
		if ops[i].Params == nil {
			ops[i].Params = map[string]any{}
		}
	}
	return ops, nil
}

func toWarnings(issues []commonvalidate.Issue) []commoncli.Warning {
	if len(issues) == 0 {
		return nil
	}
	out := make([]commoncli.Warning, 0, len(issues))
	for _, issue := range issues {
		out = append(out, commoncli.Warning{
			Code:    issue.Code,
			Field:   issue.Field,
			Message: issue.Message,
			Fix:     issue.Fix,
		})
	}
	return out
}

func writeEnvelope(format, fields string, envelope commoncli.Envelope) error {
	payload := map[string]any{
		"ok":        envelope.OK,
		"requestId": envelope.RequestID,
		"command":   envelope.Command,
		"result":    envelope.Result,
		"warnings":  envelope.Warnings,
		"timingMs":  envelope.TimingMs,
	}
	if envelope.Error != nil {
		payload["error"] = envelope.Error
		payload["retryable"] = envelope.Retryable
	}
	payload = commoncli.FilterFields(payload, fields).(map[string]any)
	switch format {
	case "text":
		return commoncli.WriteText(envelope)
	case "ndjson":
		return commoncli.WriteNDJSON([]any{payload})
	default:
		return commoncli.WriteJSON(payload)
	}
}

func writeBatch(format string, batch commoncli.BatchEnvelope) error {
	switch format {
	case "ndjson":
		records := make([]any, 0, len(batch.Steps)+1)
		for _, step := range batch.Steps {
			records = append(records, step)
		}
		records = append(records, map[string]any{
			"summary":   batch.Summary,
			"ok":        batch.OK,
			"requestId": batch.RequestID,
			"timingMs":  batch.TimingMs,
		})
		return commoncli.WriteNDJSON(records)
	default:
		return commoncli.WriteJSON(batch)
	}
}

func exampleOps(category string) []commoncli.BatchOp {
	switch category {
	case "poster":
		return []commoncli.BatchOp{
			{Name: "new_doc", Command: "document.new", Params: map[string]any{"width": 1440, "height": 1024, "artboards": 2}},
			{Name: "headline", Command: "text.create", Params: map[string]any{"name": "Headline", "contents": "AHD Illustrator", "left": 96, "top": 128}},
			{Name: "shape", Command: "path.create_rect", Params: map[string]any{"name": "Hero Block", "left": 80, "top": 760, "width": 520, "height": 320}},
		}
	case "export":
		return []commoncli.BatchOp{
			{Name: "save_as", Command: "document.save_as", Params: map[string]any{"filePath": "out/illustration.ai"}},
			{Name: "png", Command: "export.png", Params: map[string]any{"outputPath": "out/illustration.png", "scale": 2}},
		}
	default:
		return nil
	}
}
