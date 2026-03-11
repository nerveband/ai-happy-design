package main

import (
	"fmt"

	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/spf13/cobra"
)

var guideTopic string

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Design intelligence guide for LLMs",
	Long:  "Print design methodology, patterns, and quality gates. Use --topic to get a specific deep-dive.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result interface{}
		switch guideTopic {
		case "typography":
			result = tools.GuideTypography()
		case "color":
			result = tools.GuideColor()
		case "layout":
			result = tools.GuideLayout()
		case "depth":
			result = tools.GuideDepth()
		case "states":
			result = tools.GuideStates()
		case "quality":
			result = tools.GuideQuality()
		case "all":
			result = tools.GuideAll()
		case "":
			result = tools.DesignGuide()
		default:
			return fmt.Errorf("unknown topic %q — available: typography, color, layout, depth, states, quality, all", guideTopic)
		}
		return printJSON(result)
	},
}

func init() {
	guideCmd.Flags().StringVar(&guideTopic, "topic", "", "Deep-dive topic: typography, color, layout, depth, states, quality, all")
	rootCmd.AddCommand(guideCmd)
}
