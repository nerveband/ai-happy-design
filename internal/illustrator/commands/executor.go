package commands

import (
	"strings"
	"time"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/illustrator/bridge"
	"github.com/nerveband/ai-happy-design/internal/illustrator/host"
)

// Request is the normalized execution input.
type Request struct {
	Command *commonschema.Command
	Params  map[string]any
	DryRun  bool
}

// Executor is the core command runner used by the CLI.
type Executor struct {
	host   *host.Adapter
	bridge *bridge.Client
}

type executionPlan struct {
	mode     string
	script   string
	selector string
	timeout  time.Duration
	activate bool
}

// NewExecutor returns a command executor.
func NewExecutor() *Executor {
	adapter := host.NewAdapter()
	return &Executor{
		host:   adapter,
		bridge: bridge.NewClient(adapter),
	}
}

// Execute runs a validated command request.
func (e *Executor) Execute(request Request) (any, []commoncli.Warning, *commoncli.CommandError) {
	if request.Command == nil {
		return nil, nil, &commoncli.CommandError{
			Code:    "UNSUPPORTED_COMMAND",
			Message: "unknown command",
		}
	}
	warnings := commandWarnings(request.Command.Name)
	if request.DryRun {
		return map[string]any{
			"validated":      true,
			"pluginRequired": request.Command.PluginRequired,
			"mutating":       request.Command.Mutating,
			"params":         request.Params,
		}, warnings, nil
	}

	status := e.host.Status()
	if !status.Supported {
		return nil, nil, &commoncli.CommandError{
			Code:    "ILLUSTRATOR_NOT_RUNNING",
			Message: "Illustrator host support is currently macOS-only",
			Details: status,
		}
	}
	if !status.IllustratorRunning {
		return nil, nil, &commoncli.CommandError{
			Code:    "ILLUSTRATOR_NOT_RUNNING",
			Message: "Adobe Illustrator is not running",
			Details: status,
		}
	}

	plan, err := buildPlan(request)
	if err != nil {
		return nil, nil, &commoncli.CommandError{
			Code:    "UNSUPPORTED_COMMAND",
			Message: err.Error(),
			Details: map[string]any{
				"command": request.Command.Name,
			},
		}
	}

	response, cmdErr := e.executePlan(plan, request)
	for attempt := 0; attempt < 2 && shouldRetryScriptError(plan, cmdErr); attempt++ {
		if err := e.host.Activate(); err == nil {
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
		response, cmdErr = e.executePlan(plan, request)
	}
	if cmdErr != nil {
		return nil, nil, cmdErr
	}
	for _, warning := range response.Warnings {
		warnings = append(warnings, commoncli.Warning{
			Code:    "HOST_WARNING",
			Message: warning,
		})
	}
	return response.Result, warnings, nil
}

func (e *Executor) executePlan(plan executionPlan, request Request) (*bridge.Response, *commoncli.CommandError) {
	if plan.activate {
		if err := e.host.Activate(); err != nil {
			return nil, &commoncli.CommandError{
				Code:    "HOST_EXEC_ERROR",
				Message: err.Error(),
			}
		}
	}

	if plan.mode == "selector" {
		return e.bridge.ExecuteSelector(plan.selector, bridge.Request{
			V:         "1.0",
			Command:   request.Command.Name,
			Params:    request.Params,
			DryRun:    request.DryRun,
			TimeoutMs: int(plan.timeout / time.Millisecond),
		}, plan.timeout)
	}

	return e.bridge.ExecuteScript(plan.script, plan.timeout)
}

func shouldRetryScriptError(plan executionPlan, cmdErr *commoncli.CommandError) bool {
	if cmdErr == nil || plan.mode != "script" || plan.activate {
		return false
	}
	if cmdErr.Code != "HOST_EXEC_ERROR" {
		return false
	}
	message := strings.ToLower(cmdErr.Message)
	return strings.Contains(message, "timed out") ||
		strings.Contains(message, "connection is invalid") ||
		strings.Contains(message, "(-609)")
}

func commandWarnings(name string) []commoncli.Warning {
	switch strings.TrimSpace(name) {
	case "document.write_as_library":
		return []commoncli.Warning{{
			Code:    "EXPERIMENTAL_COMMAND",
			Message: "document.write_as_library is exposed as experimental. Illustrator 30.2.1 can reject documented library types at runtime.",
		}}
	case "page_item.bring_in_perspective":
		return []commoncli.Warning{{
			Code:    "EXPERIMENTAL_COMMAND",
			Message: "page_item.bring_in_perspective is exposed as experimental. Illustrator 30.2.1 can reject documented perspective plane values at runtime.",
		}}
	case "trace.preset.store":
		return []commoncli.Warning{{
			Code:    "EXPERIMENTAL_COMMAND",
			Message: "trace.preset.store is exposed as experimental. Illustrator 30.2.1 can return a native error instead of persisting the preset.",
		}}
	default:
		return nil
	}
}
