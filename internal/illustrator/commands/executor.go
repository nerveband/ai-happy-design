package commands

import (
	"fmt"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
	"github.com/nerveband/ai-happy-design/internal/commonschema"
)

// Request is the normalized execution input.
type Request struct {
	Command *commonschema.Command
	Params  map[string]any
	DryRun  bool
}

// Executor is the core command runner used by the CLI.
type Executor struct{}

// NewExecutor returns a command executor.
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute runs a validated command request. Non-dry-run execution is wired up in later phases.
func (e *Executor) Execute(request Request) (any, []commoncli.Warning, *commoncli.CommandError) {
	if request.Command == nil {
		return nil, nil, &commoncli.CommandError{
			Code:    "UNSUPPORTED_COMMAND",
			Message: "unknown command",
		}
	}
	if request.DryRun {
		return map[string]any{
			"validated":      true,
			"pluginRequired": request.Command.PluginRequired,
			"mutating":       request.Command.Mutating,
			"params":         request.Params,
		}, nil, nil
	}
	return nil, nil, &commoncli.CommandError{
		Code:    "UNSUPPORTED_COMMAND",
		Message: fmt.Sprintf("%s execution is not wired yet", request.Command.Name),
		Details: map[string]any{
			"command": request.Command.Name,
			"dryRun":  "use --dry-run until the host bridge is added",
		},
	}
}
