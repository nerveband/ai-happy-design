package main

import (
	"encoding/json"
	"fmt"

	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/spf13/cobra"
)

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Print design intelligence guide for LLMs",
	RunE: func(cmd *cobra.Command, args []string) error {
		catalog := tools.LLMCatalog()
		j, err := json.MarshalIndent(catalog, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(guideCmd)
}
