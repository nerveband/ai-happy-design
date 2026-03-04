package tools

import (
	"strings"
	"testing"
)

func TestLLMCatalogDiscoverySection(t *testing.T) {
	catalog := LLMCatalog()
	discovery, ok := catalog["discovery"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected discovery map in catalog")
	}

	for _, key := range []string{"tools", "actions", "examples", "batch"} {
		if got, ok := discovery[key].(string); !ok || strings.TrimSpace(got) == "" {
			t.Fatalf("expected non-empty discovery.%s", key)
		}
	}
}

func TestLLMCatalogLintChecksContract(t *testing.T) {
	catalog := LLMCatalog()
	lintChecks, ok := catalog["lintChecks"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected lintChecks map")
	}
	for _, key := range []string{"absolute_child_non_autolayout", "absolute_overflow", "overflow", "overlap"} {
		if _, exists := lintChecks[key]; !exists {
			t.Fatalf("expected lintChecks[%q]", key)
		}
	}
}

func TestLLMCatalogImageRulesSection(t *testing.T) {
	catalog := LLMCatalog()
	imageRules, ok := findNestedMap(catalog, "imageRules")
	if !ok {
		t.Fatalf("expected imageRules section in catalog")
	}

	for _, key := range []string{"_overview", "methods", "imageData", "scaleModes"} {
		if _, exists := imageRules[key]; !exists {
			t.Fatalf("expected imageRules[%q]", key)
		}
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

func TestDesignGuideIncludesLintChecks(t *testing.T) {
	guide := DesignGuide()
	if _, ok := guide["lintChecks"].(map[string]interface{}); !ok {
		t.Fatalf("expected design guide to include lintChecks")
	}
	if _, ok := guide["designThinking"].(map[string]interface{}); !ok {
		t.Fatalf("expected design guide to include designThinking")
	}
	if _, ok := guide["designPatterns"].(map[string]interface{}); !ok {
		t.Fatalf("expected design guide to include designPatterns")
	}
}

func TestLLMCatalogDesignPatternsPresent(t *testing.T) {
	catalog := LLMCatalog()
	dt, ok := catalog["designThinking"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected designThinking map")
	}
	for _, key := range []string{"cssToFigma", "visualHierarchy", "designDecisions", "layerOrganization"} {
		if _, exists := dt[key]; !exists {
			t.Fatalf("expected designThinking[%q]", key)
		}
	}

	dp, ok := catalog["designPatterns"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected designPatterns map")
	}
	for _, key := range []string{"coordinateSystem", "autoLayout", "sizingSystem", "balance"} {
		if _, exists := dp[key]; !exists {
			t.Fatalf("expected designPatterns[%q]", key)
		}
	}
}

func TestLLMCatalogWorkflowSection(t *testing.T) {
	catalog := LLMCatalog()
	wf, ok := catalog["workflow"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected workflow map")
	}
	for _, key := range []string{"rule", "create", "edit", "verify"} {
		if _, exists := wf[key]; !exists {
			t.Fatalf("expected workflow[%q]", key)
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
