package main

import (
	"sort"

	"github.com/nerveband/ai-happy-design/internal/schema"
	"github.com/spf13/cobra"
)

var agentContextJSON bool

var agentContextCmd = &cobra.Command{
	Use:   "agent-context",
	Short: "Print compact machine-readable guidance for agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printJSON(buildAgentContextMust())
	},
}

func buildAgentContextMust() map[string]interface{} {
	ctx, err := buildAgentContext()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return ctx
}

func buildAgentContext() (map[string]interface{}, error) {
	groups := schema.GroupedCommands()
	groupNames := make([]string, 0, len(groups))
	for group := range groups {
		groupNames = append(groupNames, group)
	}
	sort.Strings(groupNames)

	safety := map[string]int{}
	for _, s := range schema.All {
		safety[s.Safety]++
	}

	return map[string]interface{}{
		"version": version,
		"commands": map[string]interface{}{
			"groups": groupNames,
			"count":  len(schema.All),
			"schema": "ahd-figma schema --json",
		},
		"recommendedWorkflows": []string{
			"Use batch for creating or changing many nodes in one relay round trip.",
			"Use command for one targeted edit or read.",
			"Use document.screenshot or document.screenshot_selection after visual changes, then inspect the artifact before finalizing.",
		},
		"outputModes": map[string]interface{}{
			"json":   "--json on discovery commands or --output-format json globally",
			"fields": "--fields on command results for compact extraction",
			"jq":     "--jq .path for simple projection",
		},
		"safetyMetadata": map[string]interface{}{
			"counts":          safety,
			"fields":          []string{"safety", "idempotency", "supportsDryRun", "requiresFigma", "requiresRelay", "requiresAuth"},
			"commitUndoParam": "mutating plugin commands commit Figma undo by default; pass commitUndo:false to opt out",
		},
		"exitCodes": map[string]interface{}{
			"0": "success",
			"1": "validation, connection, runtime, or command failure",
		},
		"configPrecedence": []string{"flags", "environment", "config file", "defaults"},
		"artifactDelivery": map[string]interface{}{
			"exports":     "export.* saves binary data and returns savedTo/fileSize metadata unless --base64 is used",
			"screenshots": "document.screenshot and document.screenshot_selection return screenshot metadata suitable for local artifact inspection",
		},
		"docs": map[string]interface{}{
			"guide":  "ahd-figma guide",
			"tools":  "ahd-figma tools --llm --json",
			"schema": "ahd-figma schema <command> --json",
		},
		"examples": []string{
			`ahd-figma batch ops.json --strict-quality`,
			`ahd-figma command document.screenshot '{"nodeId":"1:2","scale":2}'`,
			`ahd-figma command node.modify '{"nodeId":"1:2","name":"Hero Card"}'`,
		},
	}, nil
}

func init() {
	agentContextCmd.Flags().BoolVar(&agentContextJSON, "json", false, "Output JSON")
	rootCmd.AddCommand(agentContextCmd)
}
