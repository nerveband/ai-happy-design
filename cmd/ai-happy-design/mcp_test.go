package main

import (
	"bufio"
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPListsSchemaBackedTools(t *testing.T) {
	in := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}` + "\n")
	var out strings.Builder
	if err := runMCP(in, &out); err != nil {
		t.Fatalf("runMCP: %v", err)
	}
	line, _ := bufio.NewReader(strings.NewReader(out.String())).ReadString('\n')
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		t.Fatalf("response json: %v\n%s", err, line)
	}
	result := resp["result"].(map[string]interface{})
	tools := result["tools"].([]interface{})
	found := false
	for _, raw := range tools {
		tool := raw.(map[string]interface{})
		if tool["name"] == "node_create_frame" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected node_create_frame MCP tool")
	}
}

func TestMCPComputeTokensTool(t *testing.T) {
	in := strings.NewReader(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"design_compute_tokens","arguments":{"width":1080,"height":1350}}}` + "\n")
	var out strings.Builder
	if err := runMCP(in, &out); err != nil {
		t.Fatalf("runMCP: %v", err)
	}
	if !strings.Contains(out.String(), `"content"`) || !strings.Contains(out.String(), `hero`) {
		t.Fatalf("expected token JSON content, got %s", out.String())
	}
}

func TestMCPListsAndGetsPrompts(t *testing.T) {
	in := strings.NewReader(
		`{"jsonrpc":"2.0","id":1,"method":"prompts/list"}` + "\n" +
			`{"jsonrpc":"2.0","id":2,"method":"prompts/get","params":{"name":"visual_verification_strategy"}}` + "\n",
	)
	var out strings.Builder
	if err := runMCP(in, &out); err != nil {
		t.Fatalf("runMCP: %v", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	if !scanner.Scan() {
		t.Fatal("missing prompts/list response")
	}
	if !strings.Contains(scanner.Text(), `"visual_verification_strategy"`) {
		t.Fatalf("expected visual verification prompt in list: %s", scanner.Text())
	}
	if !scanner.Scan() {
		t.Fatal("missing prompts/get response")
	}
	if !strings.Contains(scanner.Text(), `"messages"`) || !strings.Contains(scanner.Text(), "screenshot") {
		t.Fatalf("expected prompt messages with screenshot guidance: %s", scanner.Text())
	}
}

func TestMCPToolsExposeSafetyMetadata(t *testing.T) {
	tools := mcpTools()
	for _, tool := range tools {
		if tool["name"] == "node_create_frame" {
			meta, ok := tool["annotations"].(map[string]interface{})
			if !ok {
				t.Fatalf("node_create_frame missing annotations: %#v", tool)
			}
			if meta["safety"] != "write" || meta["requiresFigma"] != true || meta["requiresRelay"] != true {
				t.Fatalf("unexpected annotations: %#v", meta)
			}
			return
		}
	}
	t.Fatal("node_create_frame not found")
}
