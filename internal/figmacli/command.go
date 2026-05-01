package figmacli

import "context"

type CommandOptions struct {
	BinaryName     string
	Version        string
	Command        string
	Params         map[string]interface{}
	Channel        string
	Live           bool
	DryRun         bool
	Fields         []string
	Limit          int
	TimeoutSeconds int
}

type ResultEnvelope struct {
	OK      bool                   `json:"ok"`
	Command string                 `json:"command,omitempty"`
	Result  interface{}            `json:"result,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Error   *ErrorEnvelope         `json:"error,omitempty"`
}

type ErrorEnvelope struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
	Hint      string `json:"hint,omitempty"`
}

// ExecuteCommand is the shared contract for CLI and MCP command execution.
// The current Cobra command path owns transport wiring; this package carries
// the stable API shape so future binaries cannot drift.
func ExecuteCommand(ctx context.Context, opts CommandOptions) (map[string]interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return map[string]interface{}{
		"ok":      true,
		"command": opts.Command,
		"params":  opts.Params,
		"dryRun":  opts.DryRun,
	}, nil
}
