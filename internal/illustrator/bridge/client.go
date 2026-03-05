package bridge

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
	"github.com/nerveband/ai-happy-design/internal/illustrator/host"
)

const pluginName = "AHDIllustrator"

// Request is the selector payload sent to the bridge.
type Request struct {
	V         string         `json:"v"`
	ID        string         `json:"id"`
	Command   string         `json:"command"`
	Params    map[string]any `json:"params"`
	DryRun    bool           `json:"dryRun"`
	TimeoutMs int            `json:"timeoutMs"`
}

// Response is the structured reply from either script mode or plugin mode.
type Response struct {
	V        string          `json:"v"`
	ID       string          `json:"id"`
	OK       bool            `json:"ok"`
	Result   any             `json:"result,omitempty"`
	Warnings []string        `json:"warnings,omitempty"`
	Error    string          `json:"error,omitempty"`
	Details  json.RawMessage `json:"details,omitempty"`
}

type runtimePayload struct {
	ID       string  `json:"id"`
	Mode     string  `json:"mode"`
	Selector string  `json:"selector,omitempty"`
	Request  Request `json:"request,omitempty"`
	Script   string  `json:"script,omitempty"`
}

// Client wraps the host adapter with the shared JSX bridge runtime.
type Client struct {
	host *host.Adapter
}

// NewClient returns a bridge client bound to a host adapter.
func NewClient(adapter *host.Adapter) *Client {
	return &Client{host: adapter}
}

// ExecuteScript runs a raw ExtendScript snippet through the bridge wrapper.
func (c *Client) ExecuteScript(script string, timeout time.Duration) (*Response, *commoncli.CommandError) {
	payload := runtimePayload{
		ID:     uuid.NewString(),
		Mode:   "script",
		Script: script,
	}
	return c.execute(payload, timeout, false)
}

// ExecuteSelector sends a selector request to the optional C++ plugin bridge.
func (c *Client) ExecuteSelector(selector string, request Request, timeout time.Duration) (*Response, *commoncli.CommandError) {
	if request.V == "" {
		request.V = "1.0"
	}
	if request.ID == "" {
		request.ID = uuid.NewString()
	}
	payload := runtimePayload{
		ID:       request.ID,
		Mode:     "selector",
		Selector: selector,
		Request:  request,
	}
	return c.execute(payload, timeout, true)
}

func (c *Client) execute(payload runtimePayload, timeout time.Duration, selectorMode bool) (*Response, *commoncli.CommandError) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &commoncli.CommandError{
			Code:    "HOST_EXEC_ERROR",
			Message: "failed to encode bridge payload",
			Details: err.Error(),
		}
	}
	script := strings.ReplaceAll(bridgeTemplate, "__AHD_PAYLOAD__", string(body))
	raw, execErr := c.host.ExecuteJavaScript(script, timeout)
	if execErr != nil {
		code := "HOST_EXEC_ERROR"
		if strings.Contains(execErr.Error(), "PLUGIN_TIMEOUT") {
			code = "PLUGIN_TIMEOUT"
		}
		return nil, &commoncli.CommandError{
			Code:    code,
			Message: execErr.Error(),
		}
	}
	if strings.TrimSpace(raw) == "" && selectorMode {
		return nil, &commoncli.CommandError{
			Code:    "PLUGIN_REQUIRED",
			Message: "plugin selector returned an empty response",
			Details: map[string]any{"selector": payload.Selector},
		}
	}
	var response Response
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return nil, &commoncli.CommandError{
			Code:    "HOST_EXEC_ERROR",
			Message: "failed to decode bridge response",
			Details: map[string]any{"raw": raw, "error": err.Error()},
		}
	}
	if !response.OK {
		code := "HOST_EXEC_ERROR"
		if selectorMode {
			code = "PLUGIN_REQUIRED"
		}
		return &response, &commoncli.CommandError{
			Code:    code,
			Message: response.Error,
			Details: map[string]any{"selector": payload.Selector},
		}
	}
	return &response, nil
}

var bridgeTemplate = `(function () {
  function stringify(value) {
    try {
      return JSON.stringify(value);
    } catch (err) {
      return JSON.stringify({
        v: "1.0",
        id: (value && value.id) || "",
        ok: false,
        error: "failed to stringify bridge payload: " + err
      });
    }
  }

  var payload = __AHD_PAYLOAD__;
  try {
    if (payload.mode === "script") {
      var result = eval(payload.script);
      return stringify({
        v: "1.0",
        id: payload.id,
        ok: true,
        result: result || {},
        warnings: []
      });
    }

    var raw = app.sendScriptMessage("` + pluginName + `", payload.selector, stringify(payload.request));
    if (!raw || raw === "") {
      return stringify({
        v: "1.0",
        id: payload.request.id,
        ok: false,
        error: "empty plugin response",
        warnings: []
      });
    }
    return raw;
  } catch (err) {
    return stringify({
      v: "1.0",
      id: (payload.request && payload.request.id) || payload.id,
      ok: false,
      error: String(err),
      warnings: []
    });
  }
}());`
