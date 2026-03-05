package commoncli

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Warning is a machine-readable non-fatal issue surfaced to agents.
type Warning struct {
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Fix     any    `json:"fix,omitempty"`
}

// CommandError is the stable error payload returned by the Illustrator CLI.
type CommandError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Envelope is the common response structure for single command execution.
type Envelope struct {
	OK        bool          `json:"ok"`
	RequestID string        `json:"requestId"`
	Command   string        `json:"command,omitempty"`
	Result    any           `json:"result,omitempty"`
	Warnings  []Warning     `json:"warnings,omitempty"`
	TimingMs  int64         `json:"timingMs,omitempty"`
	Error     *CommandError `json:"error,omitempty"`
	Retryable bool          `json:"retryable,omitempty"`
}

// BatchSummary summarizes a batch execution.
type BatchSummary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

// BatchStep represents a single batch result.
type BatchStep struct {
	Index    int       `json:"index"`
	Name     string    `json:"name,omitempty"`
	Command  string    `json:"command,omitempty"`
	OK       bool      `json:"ok"`
	Result   any       `json:"result,omitempty"`
	Warnings []Warning `json:"warnings,omitempty"`
	Error    any       `json:"error,omitempty"`
}

// BatchEnvelope is the stable batch response shape.
type BatchEnvelope struct {
	OK        bool          `json:"ok"`
	RequestID string        `json:"requestId"`
	Summary   BatchSummary  `json:"summary"`
	Steps     []BatchStep   `json:"steps"`
	TimingMs  int64         `json:"timingMs"`
	Error     *CommandError `json:"error,omitempty"`
	Retryable bool          `json:"retryable,omitempty"`
}

// BatchOp is the machine-facing batch input format.
type BatchOp struct {
	Name    string         `json:"name,omitempty"`
	Command string         `json:"command"`
	Params  map[string]any `json:"params,omitempty"`
}

// NewRequestID creates a stable request identifier.
func NewRequestID() string {
	return uuid.NewString()
}

// SuccessEnvelope builds a standard success response.
func SuccessEnvelope(command string, result any, warnings []Warning, started time.Time) Envelope {
	return Envelope{
		OK:        true,
		RequestID: NewRequestID(),
		Command:   strings.TrimSpace(command),
		Result:    result,
		Warnings:  warnings,
		TimingMs:  time.Since(started).Milliseconds(),
	}
}

// ErrorEnvelope builds a standard error response.
func ErrorEnvelope(command, code, message string, details any, retryable bool, warnings []Warning, started time.Time) Envelope {
	return Envelope{
		OK:        false,
		RequestID: NewRequestID(),
		Command:   strings.TrimSpace(command),
		Warnings:  warnings,
		TimingMs:  time.Since(started).Milliseconds(),
		Error: &CommandError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Retryable: retryable,
	}
}

// BatchSuccess builds a stable batch success response.
func BatchSuccess(steps []BatchStep, started time.Time) BatchEnvelope {
	summary := BatchSummary{Total: len(steps)}
	for _, step := range steps {
		if step.OK {
			summary.Succeeded++
		} else {
			summary.Failed++
		}
	}
	return BatchEnvelope{
		OK:        summary.Failed == 0,
		RequestID: NewRequestID(),
		Summary:   summary,
		Steps:     steps,
		TimingMs:  time.Since(started).Milliseconds(),
	}
}

// BatchFailure builds a stable batch-level error response.
func BatchFailure(code, message string, details any, retryable bool, started time.Time) BatchEnvelope {
	return BatchEnvelope{
		OK:        false,
		RequestID: NewRequestID(),
		Summary:   BatchSummary{},
		Steps:     []BatchStep{},
		TimingMs:  time.Since(started).Milliseconds(),
		Error: &CommandError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Retryable: retryable,
	}
}
