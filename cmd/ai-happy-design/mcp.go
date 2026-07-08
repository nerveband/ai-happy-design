package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
	"github.com/nerveband/ai-happy-design/internal/schema"
	"github.com/nerveband/ai-happy-design/internal/tools"
)

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server over stdio",
	Long:  "Starts a schema-backed MCP server over stdio for AI editors.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)
		return runMCP(os.Stdin, os.Stdout)
	},
}

func runMCP(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var req mcpRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			writeMCP(writer, mcpResponse{JSONRPC: "2.0", Error: mcpError(-32700, err.Error())})
			continue
		}
		resp, notifyOnly := handleMCPRequest(req)
		if notifyOnly {
			continue
		}
		writeMCP(writer, resp)
	}
	return scanner.Err()
}

func handleMCPRequest(req mcpRequest) (mcpResponse, bool) {
	if req.ID == nil {
		return mcpResponse{}, true
	}
	resp := mcpResponse{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		resp.Result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo":      map[string]interface{}{"name": "ahd-figma", "version": version},
			"capabilities": map[string]interface{}{
				"tools":     map[string]interface{}{},
				"resources": map[string]interface{}{},
				"prompts":   map[string]interface{}{},
			},
		}
	case "tools/list":
		resp.Result = map[string]interface{}{"tools": mcpTools()}
	case "tools/call":
		result, err := mcpCallTool(req.Params)
		if err != nil {
			resp.Error = mcpError(-32000, err.Error())
		} else {
			resp.Result = result
		}
	case "resources/list":
		resp.Result = map[string]interface{}{"resources": []map[string]interface{}{
			{"uri": "ahd://schema", "name": "AHD command schemas", "mimeType": "application/json"},
			{"uri": "ahd://schema/all", "name": "AHD all command schemas", "mimeType": "application/json"},
			{"uri": "ahd://guide", "name": "AHD design guide", "mimeType": "application/json"},
			{"uri": "ahd://guide/design", "name": "AHD design guide", "mimeType": "application/json"},
			{"uri": "ahd://tools", "name": "AHD tool catalog", "mimeType": "application/json"},
			{"uri": "ahd://tools/catalog", "name": "AHD tool catalog", "mimeType": "application/json"},
			{"uri": "ahd://examples/batch", "name": "AHD batch examples", "mimeType": "application/json"},
		}}
	case "resources/read":
		result, err := mcpReadResource(req.Params)
		if err != nil {
			resp.Error = mcpError(-32000, err.Error())
		} else {
			resp.Result = result
		}
	case "prompts/list":
		resp.Result = map[string]interface{}{"prompts": mcpPrompts()}
	case "prompts/get":
		result, err := mcpGetPrompt(req.Params)
		if err != nil {
			resp.Error = mcpError(-32000, err.Error())
		} else {
			resp.Result = result
		}
	default:
		resp.Error = mcpError(-32601, "method not found: "+req.Method)
	}
	return resp, false
}

func mcpTools() []map[string]interface{} {
	out := []map[string]interface{}{
		{
			"name":        "ahd_describe",
			"description": "Return the schema-backed tool catalog or design guide.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"catalog", "guide"}}},
			},
		},
	}
	for _, s := range schema.All {
		props := map[string]interface{}{}
		required := []string{}
		for _, p := range s.Params {
			props[p.Name] = mcpParamSchema(p)
			if p.Required {
				required = append(required, p.Name)
			}
		}
		input := map[string]interface{}{"type": "object", "properties": props}
		if len(required) > 0 {
			input["required"] = required
		}
		out = append(out, map[string]interface{}{
			"name":        mcpToolName(s.Command),
			"description": s.Description,
			"inputSchema": input,
			"annotations": map[string]interface{}{
				"safety":         s.Safety,
				"idempotency":    s.Idempotency,
				"supportsDryRun": s.SupportsDryRun,
				"requiresFigma":  s.RequiresFigma,
				"requiresRelay":  s.RequiresRelay,
				"requiresAuth":   s.RequiresAuth,
			},
		})
	}
	return out
}

func mcpParamSchema(p schema.Param) map[string]interface{} {
	out := map[string]interface{}{"type": p.Type, "description": p.Desc}
	if len(p.Enum) > 0 {
		out["enum"] = p.Enum
	}
	if p.Min != nil {
		out["minimum"] = *p.Min
	}
	if p.Max != nil {
		out["maximum"] = *p.Max
	}
	return out
}

func mcpCallTool(raw json.RawMessage) (interface{}, error) {
	var payload struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	args := payload.Arguments
	if args == nil {
		args = map[string]interface{}{}
	}
	if payload.Name == "ahd_describe" {
		action, _ := args["action"].(string)
		if action == "guide" {
			return mcpTextContent(toJSONText(tools.DesignGuide())), nil
		}
		return mcpTextContent(toJSONText(tools.LLMCatalog())), nil
	}
	command, ok := commandFromMCPToolName(payload.Name)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", payload.Name)
	}
	params := batchutil.NormalizeBatchParams(command, args)
	channel, _ := params["channel"].(string)
	delete(params, "channel")
	if handled, result, err := handleLocalCommand(command, params); handled {
		if err != nil {
			return nil, err
		}
		return mcpTextContent(toJSONText(result)), nil
	}
	if err := ensureRelayIfNeeded(); err != nil {
		return nil, err
	}
	channelKey, err := resolveChannel(channel)
	if err != nil {
		return nil, err
	}
	if err := checkPluginConnected(channelKey); err != nil {
		return nil, err
	}
	client, err := newConnectedClient(channelKey, false)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	result, err := client.SendCommand(command, params)
	if err != nil {
		return nil, err
	}
	var parsed interface{}
	if err := json.Unmarshal(result, &parsed); err == nil {
		return mcpTextContent(toJSONText(parsed)), nil
	}
	return mcpTextContent(string(result)), nil
}

func mcpReadResource(raw json.RawMessage) (interface{}, error) {
	var payload struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	var text string
	switch payload.URI {
	case "ahd://schema", "ahd://schema/all":
		text = toJSONText(schema.All)
	case "ahd://guide", "ahd://guide/design":
		text = toJSONText(tools.DesignGuide())
	case "ahd://tools", "ahd://tools/catalog":
		text = toJSONText(tools.LLMCatalog())
	case "ahd://examples/batch":
		text = toJSONText([]map[string]interface{}{
			{"path": "docs/examples/batch-interpolation.json", "description": "Interpolation and root frame validation example"},
			{"path": "docs/examples/live-acceptance-full-parity.json", "description": "Live Figma acceptance batch"},
			{"path": "docs/examples/html-css-recreation-workflow.json", "description": "Measured HTML/CSS-like recreation workflow"},
		})
	default:
		if strings.HasPrefix(payload.URI, "ahd://schema/") {
			name := strings.TrimPrefix(payload.URI, "ahd://schema/")
			s := schema.Lookup(name)
			if s == nil {
				return nil, fmt.Errorf("unknown schema resource: %s", name)
			}
			text = toJSONText(s)
			break
		}
		return nil, fmt.Errorf("unknown resource: %s", payload.URI)
	}
	return map[string]interface{}{"contents": []map[string]interface{}{{"uri": payload.URI, "mimeType": "application/json", "text": text}}}, nil
}

func mcpToolName(command string) string {
	return strings.ReplaceAll(command, ".", "_")
}

func commandFromMCPToolName(name string) (string, bool) {
	for _, s := range schema.All {
		if mcpToolName(s.Command) == name {
			return s.Command, true
		}
	}
	return "", false
}

func mcpPrompts() []map[string]interface{} {
	prompts := tools.MCPPrompts()
	out := make([]map[string]interface{}, 0, len(prompts))
	for _, prompt := range prompts {
		out = append(out, map[string]interface{}{
			"name":        prompt.Name,
			"description": prompt.Description,
		})
	}
	return out
}

func mcpGetPrompt(raw json.RawMessage) (interface{}, error) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	prompt, ok := tools.GetMCPPrompt(payload.Name)
	if !ok {
		return nil, fmt.Errorf("unknown prompt: %s", payload.Name)
	}
	return map[string]interface{}{
		"description": prompt.Description,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": map[string]interface{}{
					"type": "text",
					"text": prompt.Text,
				},
			},
		},
	}, nil
}

func mcpTextContent(text string) map[string]interface{} {
	return map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": text}}}
}

func toJSONText(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error":%q}`, err.Error())
	}
	return string(data)
}

func mcpError(code int, message string) map[string]interface{} {
	return map[string]interface{}{"code": code, "message": message}
}

func writeMCP(w *bufio.Writer, resp mcpResponse) {
	data, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(data))
	w.Flush()
}
