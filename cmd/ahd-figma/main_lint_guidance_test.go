package main

import (
	"strings"
	"testing"
)

func TestLintGuidanceIncludesAbsoluteChildFixHint(t *testing.T) {
	summary := lintSummary{
		Issues: 1,
		ByType: map[string]int{
			"absolute_child_non_autolayout": 1,
		},
	}

	guidance := lintGuidance(summary)
	joined := strings.ToLower(strings.Join(guidance, " | "))
	if !strings.Contains(joined, "absolute") || !strings.Contains(joined, "auto-layout") {
		t.Fatalf("expected absolute/auto-layout hint, got %v", guidance)
	}
}

func TestLintGuidanceIncludesBannerHeadlineHint(t *testing.T) {
	summary := lintSummary{
		Issues: 1,
		ByType: map[string]int{
			"overlap": 1,
		},
		Samples: []lintIssueSample{
			{
				Type:     "overlap",
				NodeName: "Headline — Help AMC raise",
				Message:  "Overlaps with sibling \"Subtitle\" in parent \"Banner — Hero\".",
			},
		},
	}

	guidance := lintGuidance(summary)
	joined := strings.ToLower(strings.Join(guidance, " | "))
	if !strings.Contains(joined, "adaptive headline sizing") {
		t.Fatalf("expected banner overlap adaptive sizing hint, got %v", guidance)
	}
}
