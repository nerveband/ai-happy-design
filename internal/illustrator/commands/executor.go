package commands

import (
	"fmt"
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

// NewExecutor returns a command executor.
func NewExecutor() *Executor {
	adapter := host.NewAdapter()
	return &Executor{
		host:   adapter,
		bridge: bridge.NewClient(adapter),
	}
}

type executionPlan struct {
	mode     string
	script   string
	selector string
	timeout  time.Duration
}

// Execute runs a validated command request.
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

	var (
		response *bridge.Response
		cmdErr   *commoncli.CommandError
	)
	if plan.mode == "selector" {
		response, cmdErr = e.bridge.ExecuteSelector(plan.selector, bridge.Request{
			V:         "1.0",
			Command:   request.Command.Name,
			Params:    request.Params,
			DryRun:    request.DryRun,
			TimeoutMs: int(plan.timeout / time.Millisecond),
		}, plan.timeout)
	} else {
		response, cmdErr = e.bridge.ExecuteScript(plan.script, plan.timeout)
	}
	if cmdErr != nil {
		return nil, nil, cmdErr
	}
	warnings := make([]commoncli.Warning, 0, len(response.Warnings))
	for _, warning := range response.Warnings {
		warnings = append(warnings, commoncli.Warning{
			Code:    "HOST_WARNING",
			Message: warning,
		})
	}
	return response.Result, warnings, nil
}

func buildPlan(request Request) (executionPlan, error) {
	switch request.Command.Name {
	case "app.info":
		return executionPlan{mode: "script", timeout: 10 * time.Second, script: `(function () {
  return {
    name: app.name,
    version: app.version,
    documents: app.documents.length,
    activeDocument: app.documents.length ? app.activeDocument.name : null
  };
}())`}, nil
	case "app.version":
		return executionPlan{mode: "script", timeout: 10 * time.Second, script: `(function () {
  return { version: app.version };
}())`}, nil
	case "document.list":
		return executionPlan{mode: "script", timeout: 10 * time.Second, script: `(function () {
  var docs = [];
  for (var i = 0; i < app.documents.length; i++) {
    docs.push({
      index: i,
      name: app.documents[i].name
    });
  }
  return docs;
}())`}, nil
	case "document.info":
		return executionPlan{mode: "script", timeout: 10 * time.Second, script: `(function () {
  if (app.documents.length === 0) {
    return {};
  }
  var doc = app.activeDocument;
  return {
    name: doc.name,
    artboards: doc.artboards.length,
    layers: doc.layers.length
  };
}())`}, nil
	default:
		return executionPlan{}, fmt.Errorf("%s execution is not wired yet", request.Command.Name)
	}
}
