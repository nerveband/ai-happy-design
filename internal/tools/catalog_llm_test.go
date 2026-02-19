package tools

import (
	"fmt"
	"strings"
	"testing"
)

func TestLLMCatalogDiscoveryIncludesImagePrepAndQualityHints(t *testing.T) {
	catalog := LLMCatalog()
	discovery, ok := catalog["discovery"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected discovery map in catalog")
	}

	if got, ok := discovery["batchHelp"].(string); !ok || strings.TrimSpace(got) == "" {
		t.Fatalf("expected non-empty discovery.batchHelp")
	}

	quickStart, ok := discovery["quickStart"].([]string)
	if !ok {
		t.Fatalf("expected discovery.quickStart as []string")
	}
	if !containsString(quickStart, "ai-happy-design batch --help") {
		t.Fatalf("expected quickStart to include batch help, got: %v", quickStart)
	}

	hint := fmt.Sprint(discovery["batchOutputFields"])
	if !strings.Contains(hint, "output.imagePrep") || !strings.Contains(hint, "summary.qualityGate") {
		t.Fatalf("expected batch output hint to mention imagePrep and qualityGate, got: %q", hint)
	}
}

func TestLLMCatalogBatchObservabilityContract(t *testing.T) {
	catalog := LLMCatalog()
	obs, ok := catalog["batchObservability"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected batchObservability section")
	}

	fields, ok := obs["fields"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected batchObservability.fields map")
	}
	for _, key := range []string{"summary", "timing", "imagePrep"} {
		if _, exists := fields[key]; !exists {
			t.Fatalf("expected batchObservability.fields[%q]", key)
		}
	}
}

func TestLLMCatalogImageRulesExposeBatchPrepPipeline(t *testing.T) {
	catalog := LLMCatalog()
	imageRules, ok := findNestedMap(catalog, "imageRules")
	if !ok {
		t.Fatalf("expected imageRules section in catalog")
	}

	pipeline, ok := imageRules["batchPrepPipeline"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected imageRules.batchPrepPipeline section")
	}

	reporting := fmt.Sprint(pipeline["reporting"])
	if !strings.Contains(reporting, "output.imagePrep") {
		t.Fatalf("expected batch prep reporting hint to mention output.imagePrep, got %q", reporting)
	}
}

func TestSetupInstructionsIncludeCLIDiscoverySequence(t *testing.T) {
	text := SetupInstructions()
	for _, fragment := range []string{
		"ai-happy-design tools --llm --json",
		"ai-happy-design actions paint",
		"ai-happy-design batch --help",
	} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("expected setup instructions to include %q", fragment)
		}
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func findNestedMap(node interface{}, key string) (map[string]interface{}, bool) {
	switch typed := node.(type) {
	case map[string]interface{}:
		if value, exists := typed[key]; exists {
			if out, ok := value.(map[string]interface{}); ok {
				return out, true
			}
		}
		for _, value := range typed {
			if out, ok := findNestedMap(value, key); ok {
				return out, true
			}
		}
	case []interface{}:
		for _, value := range typed {
			if out, ok := findNestedMap(value, key); ok {
				return out, true
			}
		}
	}
	return nil, false
}
