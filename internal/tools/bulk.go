package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/batchutil"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// BulkOperation represents a single operation in a bulk execute request.
type BulkOperation struct {
	Name    string                 `json:"name,omitempty"`
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

// RegisterBulkTool registers the "bulk" tool for executing multiple operations in sequence.
func RegisterBulkTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("bulk",
		mcp.WithDescription("Execute multiple Figma operations in sequence with retries and optional interpolation. Compact aliases: frame/rect/text/fill/stroke/gradient/shadow/blur/glass/noise/modify/mask/find. Shorthands: pid=parentId, w=width, h=height, sz=fontSize, ff=fontFamily, lh=lineHeight(auto-PERCENT), bg=color, r=cornerRadius, fillColor=color. Step names must be snake_case (auto-sanitized). Always use semantic names for layers."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("execute")),
		mcp.WithString("operations", mcp.Required(), mcp.Description("JSON array of operations: [{\"command\": \"...\", \"params\": {...}}, ...]")),
		mcp.WithBoolean("failFast", mcp.Description("Stop at first failed operation (default false)")),
		mcp.WithNumber("retries", mcp.Description("Retry count per operation after first attempt (default 1)")),
		mcp.WithNumber("retryDelayMs", mcp.Description("Delay between retries in milliseconds (default 250)")),
		mcp.WithBoolean("interpolate", mcp.Description("Enable placeholders like ${{steps.0.result.id}} from prior steps (default true)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "execute":
			opsStr, errResult := requireStringArg(args, "operations")
			if errResult != nil {
				return errResult, nil
			}

			var ops []BulkOperation
			if err := json.Unmarshal([]byte(opsStr), &ops); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("invalid operations JSON: %v", err)), nil
			}

			if len(ops) == 0 {
				return mcp.NewToolResultError("operations array is empty"), nil
			}

			failFast := getBoolArg(args, "failFast", false)
			interpolate := getBoolArg(args, "interpolate", true)
			retries := int(getFloat64Arg(args, "retries", 1))
			retryDelayMs := int(getFloat64Arg(args, "retryDelayMs", 250))
			if retries < 0 {
				return mcp.NewToolResultError("retries must be >= 0"), nil
			}
			if retryDelayMs < 0 {
				return mcp.NewToolResultError("retryDelayMs must be >= 0"), nil
			}
			retryDelay := time.Duration(retryDelayMs) * time.Millisecond

			results := make([]map[string]interface{}, 0, len(ops))
			states := make([]batchutil.StepState, 0, len(ops))
			succeeded := 0
			failed := 0
			retriesUsed := 0
			stoppedEarly := false
			batchStart := time.Now()

			for i, op := range ops {
				opStart := time.Now()
				op.Name = batchutil.SanitizeStepName(op.Name)
				params := batchutil.NormalizeBatchParams(op.Command, op.Params)

				if interpolate {
					interpolatedParams, err := batchutil.InterpolateParams(params, states)
					if err != nil {
						entry := map[string]interface{}{
							"index":     i,
							"name":      op.Name,
							"command":   op.Command,
							"ok":        false,
							"attempts":  0,
							"error":     fmt.Sprintf("interpolation error: %v", err),
							"elapsedMs": int(time.Since(opStart).Milliseconds()),
						}
						results = append(results, entry)
						states = append(states, batchutil.StepState{
							Index:   i,
							Name:    op.Name,
							Command: op.Command,
							OK:      false,
							Error:   entry["error"].(string),
						})
						failed++
						if failFast {
							stoppedEarly = true
							break
						}
						continue
					}
					params = interpolatedParams
				}

				maxAttempts := retries + 1
				attempts := 0
				var result interface{}
				var sendErr error
				for attempts = 1; attempts <= maxAttempts; attempts++ {
					result, sendErr = commander.SendCommand(op.Command, params)
					if sendErr == nil {
						break
					}
					if attempts < maxAttempts && retryDelay > 0 {
						time.Sleep(retryDelay)
					}
				}
				if attempts > 1 {
					retriesUsed += attempts - 1
				}

				if sendErr != nil {
					entry := map[string]interface{}{
						"index":     i,
						"name":      op.Name,
						"command":   op.Command,
						"ok":        false,
						"attempts":  attempts,
						"error":     sendErr.Error(),
						"elapsedMs": int(time.Since(opStart).Milliseconds()),
					}
					results = append(results, entry)
					states = append(states, batchutil.StepState{
						Index:   i,
						Name:    op.Name,
						Command: op.Command,
						OK:      false,
						Error:   sendErr.Error(),
					})
					failed++
					if failFast {
						stoppedEarly = true
						break
					}
					continue
				}

				entry := map[string]interface{}{
					"index":     i,
					"name":      op.Name,
					"command":   op.Command,
					"ok":        true,
					"attempts":  attempts,
					"result":    result,
					"elapsedMs": int(time.Since(opStart).Milliseconds()),
				}
				results = append(results, entry)
				states = append(states, batchutil.StepState{
					Index:   i,
					Name:    op.Name,
					Command: op.Command,
					OK:      true,
					Result:  result,
				})
				succeeded++
			}

			processed := len(results)
			pending := len(ops) - processed
			totalElapsed := time.Since(batchStart)
			totalMs := int(totalElapsed.Milliseconds())
			opsPerSec := float64(0)
			avgMs := float64(0)
			if processed > 0 {
				opsPerSec = float64(processed) / totalElapsed.Seconds()
				avgMs = float64(totalMs) / float64(processed)
			}
			out := map[string]interface{}{
				"ok": failed == 0 && pending == 0,
				"summary": map[string]interface{}{
					"total":         len(ops),
					"processed":     processed,
					"succeeded":     succeeded,
					"failed":        failed,
					"pending":       pending,
					"retriesUsed":   retriesUsed,
					"failFast":      failFast,
					"interpolation": interpolate,
				},
				"timing": map[string]interface{}{
					"totalMs":   totalMs,
					"avgMs":     int(avgMs),
					"opsPerSec": int(opsPerSec),
				},
				"stoppedEarly": stoppedEarly,
				"steps":        results,
			}

			text, err := formatResult(out)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to format results: %v", err)), nil
			}
			return mcp.NewToolResultText(text), nil

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown bulk action: %s", action)), nil
		}
	})
}
