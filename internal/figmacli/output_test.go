package figmacli

import (
	"context"
	"encoding/json"
	"testing"
)

func TestExecuteCommandContractShape(t *testing.T) {
	out, err := ExecuteCommand(context.Background(), CommandOptions{Command: "design.compute_tokens", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if out["command"] != "design.compute_tokens" || out["dryRun"] != true {
		t.Fatalf("unexpected output: %#v", out)
	}
}

func TestMarshalJSONL(t *testing.T) {
	data, err := Marshal("jsonl", map[string]interface{}{"ok": true})
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["ok"] != true {
		t.Fatalf("unexpected jsonl: %s", data)
	}
}
