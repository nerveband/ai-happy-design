package bridge

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
	"github.com/nerveband/ai-happy-design/internal/illustrator/host"
)

const (
	pluginName          = "AHDIllustrator"
	pluginProbeSelector = "ahd.version"
)

const pluginRemediation = "Build and install the AHD Illustrator plugin bridge for Illustrator 2026, then retry the selector-backed command."

type scriptRunner interface {
	ExecuteJavaScript(script string, timeout time.Duration) (string, error)
}

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

// PluginStatus describes whether the optional plugin bridge can be reached.
type PluginStatus struct {
	Reachable   bool   `json:"reachable"`
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Selector    string `json:"selector,omitempty"`
	ProbeID     string `json:"probeId,omitempty"`
	Version     string `json:"version,omitempty"`
	Remediation string `json:"remediation,omitempty"`
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
	host scriptRunner
}

// NewClient returns a bridge client bound to a host adapter.
func NewClient(adapter *host.Adapter) *Client {
	return &Client{host: adapter}
}

// ProbePlugin checks whether the optional sendScriptMessage bridge is live.
func (c *Client) ProbePlugin(timeout time.Duration) PluginStatus {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	request := Request{
		V:         "1.0",
		ID:        uuid.NewString(),
		Command:   "plugin.probe",
		Params:    map[string]any{},
		TimeoutMs: int(timeout / time.Millisecond),
	}
	response, err := c.execute(runtimePayload{
		ID:       request.ID,
		Mode:     "selector",
		Selector: pluginProbeSelector,
		Request:  request,
	}, timeout, true)
	if err != nil {
		code := err.Code
		message := err.Message
		if strings.Contains(strings.ToLower(message), "timed out") {
			code = "PLUGIN_REQUIRED"
			message = "AHD Illustrator plugin bridge is not installed or not responding"
		}
		return PluginStatus{
			Reachable:   false,
			Code:        code,
			Message:     message,
			Selector:    pluginProbeSelector,
			ProbeID:     request.ID,
			Remediation: pluginRemediation,
		}
	}

	status := PluginStatus{
		Reachable: true,
		Selector:  pluginProbeSelector,
		ProbeID:   request.ID,
	}
	switch result := response.Result.(type) {
	case map[string]any:
		if version, ok := result["version"].(string); ok {
			status.Version = version
		}
	case string:
		status.Version = result
	}
	return status
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
	probe := c.ProbePlugin(minDuration(timeout, 5*time.Second))
	if !probe.Reachable {
		return nil, pluginRequiredError(selector, probe.Message, probe)
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
	script, err := buildRuntimeScript(payload)
	if err != nil {
		return nil, &commoncli.CommandError{
			Code:    "HOST_EXEC_ERROR",
			Message: "failed to encode bridge payload",
			Details: err.Error(),
		}
	}

	raw, execErr := c.host.ExecuteJavaScript(script, timeout)
	if execErr != nil {
		return nil, classifyHostExecutionError(execErr, selectorMode, payload.Selector)
	}

	return parseBridgeResponse(raw, selectorMode, payload.Selector)
}

func buildRuntimeScript(payload runtimePayload) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(bridgeTemplate, "__AHD_PAYLOAD__", string(body)), nil
}

func parseBridgeResponse(raw string, selectorMode bool, selector string) (*Response, *commoncli.CommandError) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		if selectorMode {
			return nil, pluginRequiredError(selector, "plugin selector returned an empty response", nil)
		}
		return nil, &commoncli.CommandError{
			Code:    "HOST_EXEC_ERROR",
			Message: "host bridge returned an empty response",
		}
	}

	var response Response
	if err := json.Unmarshal([]byte(trimmed), &response); err != nil {
		return nil, &commoncli.CommandError{
			Code:    "HOST_EXEC_ERROR",
			Message: "failed to decode bridge response",
			Details: map[string]any{"raw": trimmed, "error": err.Error()},
		}
	}
	if response.OK {
		return &response, nil
	}

	if selectorMode && isPluginUnavailable(response.Error) {
		return &response, pluginRequiredError(selector, response.Error, map[string]any{
			"selector": selector,
			"raw":      trimmed,
		})
	}

	code := "HOST_EXEC_ERROR"
	if selectorMode && strings.Contains(strings.ToUpper(response.Error), "TIMEOUT") {
		code = "PLUGIN_TIMEOUT"
	}
	return &response, &commoncli.CommandError{
		Code:    code,
		Message: response.Error,
		Details: map[string]any{
			"selector": selector,
			"details":  response.Details,
		},
	}
}

func classifyHostExecutionError(execErr error, selectorMode bool, selector string) *commoncli.CommandError {
	message := strings.TrimSpace(execErr.Error())
	if selectorMode && isPluginUnavailable(message) {
		return pluginRequiredError(selector, message, nil)
	}
	code := "HOST_EXEC_ERROR"
	if selectorMode && strings.Contains(strings.ToUpper(message), "TIMEOUT") {
		code = "PLUGIN_TIMEOUT"
	}
	return &commoncli.CommandError{
		Code:    code,
		Message: message,
	}
}

func pluginRequiredError(selector, message string, details any) *commoncli.CommandError {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "AHD Illustrator plugin bridge is not installed or not responding"
	}
	return &commoncli.CommandError{
		Code:    "PLUGIN_REQUIRED",
		Message: msg,
		Details: map[string]any{
			"selector":    selector,
			"remediation": pluginRemediation,
			"probe":       details,
		},
	}
}

func isPluginUnavailable(message string) bool {
	needle := strings.ToLower(strings.TrimSpace(message))
	if needle == "" {
		return false
	}
	snippets := []string{
		"1344357988",
		"df!p",
		"sendscriptmessage",
		"empty plugin response",
		"no messaging plug-in",
		"no messaging plugin",
	}
	for _, snippet := range snippets {
		if strings.Contains(needle, strings.ToLower(snippet)) {
			return true
		}
	}
	return false
}

func minDuration(a, b time.Duration) time.Duration {
	if a <= 0 {
		return b
	}
	if a < b {
		return a
	}
	return b
}

var bridgeTemplate = `(function () {
  function ahdQuote(value) {
    var str = String(value);
    var out = '"';
    var i;
    var ch;
    var code;
    var hex;
    for (i = 0; i < str.length; i++) {
      ch = str.charAt(i);
      code = str.charCodeAt(i);
      if (ch === "\\") {
        out += "\\\\";
      } else if (ch === "\"") {
        out += "\\\"";
      } else if (ch === "\b") {
        out += "\\b";
      } else if (ch === "\f") {
        out += "\\f";
      } else if (ch === "\n") {
        out += "\\n";
      } else if (ch === "\r") {
        out += "\\r";
      } else if (ch === "\t") {
        out += "\\t";
      } else if (code < 32) {
        hex = code.toString(16);
        out += "\\u" + ("0000" + hex).slice(-4);
      } else {
        out += ch;
      }
    }
    return out + '"';
  }

  function ahdIsArray(value) {
    return Object.prototype.toString.call(value) === "[object Array]";
  }

  function ahdStringify(value, seen) {
    var i;
    var keys;
    var parts;
    var item;
    if (!seen) {
      seen = [];
    }
    if (value === null || typeof value === "undefined") {
      return "null";
    }
    if (typeof value === "string") {
      return ahdQuote(value);
    }
    if (typeof value === "number") {
      return isFinite(value) ? String(value) : "null";
    }
    if (typeof value === "boolean") {
      return value ? "true" : "false";
    }
    if (typeof value === "function") {
      return "null";
    }
    for (i = 0; i < seen.length; i++) {
      if (seen[i] === value) {
        return ahdQuote("[Circular]");
      }
    }
    seen.push(value);
    if (ahdIsArray(value)) {
      parts = [];
      for (i = 0; i < value.length; i++) {
        parts.push(ahdStringify(value[i], seen));
      }
      seen.pop();
      return "[" + parts.join(",") + "]";
    }
    parts = [];
    keys = [];
    for (item in value) {
      if (Object.prototype.hasOwnProperty.call(value, item)) {
        keys.push(item);
      }
    }
    keys.sort();
    for (i = 0; i < keys.length; i++) {
      item = keys[i];
      if (typeof value[item] === "undefined" || typeof value[item] === "function") {
        continue;
      }
      parts.push(ahdQuote(item) + ":" + ahdStringify(value[item], seen));
    }
    seen.pop();
    return "{" + parts.join(",") + "}";
  }

  function ahdSuccess(id, result) {
    return ahdStringify({
      v: "1.0",
      id: id,
      ok: true,
      result: typeof result === "undefined" ? {} : result,
      warnings: []
    });
  }

  function ahdError(id, message) {
    return "{\"v\":\"1.0\",\"id\":" + ahdQuote(id || "") + ",\"ok\":false,\"error\":" + ahdQuote(String(message || "unknown bridge error")) + ",\"warnings\":[]}";
  }

  var payload = __AHD_PAYLOAD__;
  try {
    if (payload.mode === "script") {
      return ahdSuccess(payload.id, eval(payload.script));
    }

    var requestBody = ahdStringify(payload.request || {});
    var raw = app.sendScriptMessage("` + pluginName + `", payload.selector, requestBody);
    if (!raw || raw === "") {
      return ahdError((payload.request && payload.request.id) || payload.id, "empty plugin response");
    }
    return raw;
  } catch (err) {
    return ahdError((payload.request && payload.request.id) || payload.id, err);
  }
}());`
