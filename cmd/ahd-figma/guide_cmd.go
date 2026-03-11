package main

import (
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/spf13/cobra"
)

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Print design intelligence guide for LLMs",
	RunE: func(cmd *cobra.Command, args []string) error {
		catalog := tools.LLMCatalog()
		return printJSON(catalog)
	},
}

func init() {
	rootCmd.AddCommand(guideCmd)
}
