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
