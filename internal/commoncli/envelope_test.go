package commoncli

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSuccessEnvelopeMarshalsStableKeys(t *testing.T) {
	t.Parallel()

	envelope := SuccessEnvelope("text.create", map[string]any{"id": "abc"}, []Warning{{Code: "WARN", Message: "ok"}}, time.Now())
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	for _, key := range []string{"ok", "requestId", "command", "result", "timingMs"} {
		if _, exists := decoded[key]; !exists {
			t.Fatalf("expected key %q in envelope: %s", key, string(data))
		}
	}
}

func TestBatchSuccessSummaryCountsSteps(t *testing.T) {
	t.Parallel()

	batch := BatchSuccess([]BatchStep{
		{Index: 0, OK: true},
		{Index: 1, OK: false},
	}, time.Now())

	if batch.Summary.Total != 2 || batch.Summary.Succeeded != 1 || batch.Summary.Failed != 1 {
		t.Fatalf("unexpected summary: %+v", batch.Summary)
	}
}
