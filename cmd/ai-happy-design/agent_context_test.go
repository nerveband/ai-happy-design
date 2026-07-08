package main

import (
	"encoding/json"
	"testing"
)

func TestAgentContextJSON(t *testing.T) {
	ctx, err := buildAgentContext()
	if err != nil {
		t.Fatalf("agent-context --json: %v", err)
	}
	out, err := json.Marshal(ctx)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, string(out))
	}
	if payload["version"] == "" {
		t.Fatalf("missing version: %#v", payload)
	}
	if payload["recommendedWorkflows"] == nil || payload["safetyMetadata"] == nil || payload["artifactDelivery"] == nil {
		t.Fatalf("missing agent context sections: %#v", payload)
	}
}
