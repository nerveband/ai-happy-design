package commands_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/illustrator/commands"
	_ "github.com/nerveband/ai-happy-design/internal/illustrator/schema"
	illustratorvalidate "github.com/nerveband/ai-happy-design/internal/illustrator/validate"
)

func TestAcceptanceFixturesValidateInDryRunMode(t *testing.T) {
	t.Parallel()

	fixtures := []string{
		"testdata/acceptance/multi_artboard_poster.json",
		"testdata/acceptance/inspect_summary.json",
	}
	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(filepath.Base(fixture), func(t *testing.T) {
			t.Parallel()
			runAcceptanceFixture(t, fixture)
		})
	}
}

func runAcceptanceFixture(t *testing.T, fixture string) {
	t.Helper()

	data, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var ops []commoncli.BatchOp
	if err := json.Unmarshal(data, &ops); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}

	var steps []commoncli.BatchStep
	for index, op := range ops {
		resolved, err := commands.ResolveBatchOp(op, steps)
		if err != nil {
			t.Fatalf("resolve op %d (%s): %v", index, op.Name, err)
		}
		schema := commonschema.Lookup(resolved.Command)
		if schema == nil {
			t.Fatalf("missing schema for %s", resolved.Command)
		}
		validation := illustratorvalidate.ValidateCommand(schema, resolved.Params)
		if len(validation.Errors) != 0 {
			t.Fatalf("fixture %s step %d errors: %+v", fixture, index, validation.Errors)
		}
		steps = append(steps, commoncli.BatchStep{
			Name:    resolved.Name,
			Command: resolved.Command,
			OK:      true,
			Result: map[string]any{
				"validated":      true,
				"params":         validation.Params,
				"pluginRequired": schema.PluginRequired,
				"mutating":       schema.Mutating,
			},
		})
	}
}
