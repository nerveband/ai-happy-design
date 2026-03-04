package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nerveband/ai-happy-design/internal/schema"
	"github.com/spf13/cobra"
)

var schemaJSON bool
var schemaAllLLMSTxt bool

var schemaCmd = &cobra.Command{
	Use:   "schema [command.action]",
	Short: "Print parameter schema for a command",
	Long:  "Print the exact parameter schema for a command. LLMs can use this to generate valid JSON.\nUse --all --llms-txt to generate llms-full.txt for aihappydesign.com.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if schemaAllLLMSTxt {
			return printLLMSTxt()
		}

		if len(args) == 0 {
			// List all commands
			for _, s := range schema.All {
				aliases := ""
				if len(s.Aliases) > 0 {
					aliases = " (aliases: " + strings.Join(s.Aliases, ", ") + ")"
				}
				fmt.Printf("%-30s %s%s\n", s.Command, s.Description, aliases)
			}
			return nil
		}

		s := schema.Lookup(args[0])
		if s == nil {
			return fmt.Errorf("unknown command: %s. Run 'ai-happy-design schema' to list all commands.", args[0])
		}

		if schemaJSON {
			out, _ := json.MarshalIndent(s, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		// Human-readable table
		fmt.Printf("## %s\n\n%s\n", s.Command, s.Description)
		if len(s.Aliases) > 0 {
			fmt.Printf("Aliases: %s\n", strings.Join(s.Aliases, ", "))
		}
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "PARAM\tTYPE\tREQUIRED\tDEFAULT\tCONSTRAINTS\tDESCRIPTION\n")
		for _, p := range s.Params {
			req := ""
			if p.Required {
				req = "yes"
			}
			def := ""
			if p.Default != nil {
				def = fmt.Sprintf("%v", p.Default)
			}
			constraints := buildConstraints(p)
			desc := p.Desc
			if len(p.Aliases) > 0 {
				desc += " (alias: " + strings.Join(p.Aliases, "/") + ")"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", p.Name, p.Type, req, def, constraints, desc)
		}
		w.Flush()
		return nil
	},
}

func buildConstraints(p schema.Param) string {
	parts := []string{}
	if p.Min != nil && p.Max != nil {
		parts = append(parts, fmt.Sprintf("%.0f-%.0f", *p.Min, *p.Max))
	} else if p.Min != nil {
		parts = append(parts, fmt.Sprintf(">= %.0f", *p.Min))
	} else if p.Max != nil {
		parts = append(parts, fmt.Sprintf("<= %.0f", *p.Max))
	}
	if len(p.Enum) > 0 {
		parts = append(parts, strings.Join(p.Enum, "/"))
	}
	if p.Pattern != "" {
		parts = append(parts, "pattern: "+p.Pattern)
	}
	if p.SemanticTokens {
		parts = append(parts, "tokens: hero/title/heading/subheading/body/caption")
	}
	return strings.Join(parts, ", ")
}

func printLLMSTxt() error {
	fmt.Println("# AI Happy Design — Full Command Reference")
	fmt.Println()
	fmt.Println("Auto-generated from CLI schemas. Use `ai-happy-design schema <command> --json` for machine-readable format.")
	fmt.Println()

	for _, s := range schema.All {
		fmt.Printf("## %s\n\n", s.Command)
		fmt.Printf("%s\n\n", s.Description)
		if len(s.Aliases) > 0 {
			fmt.Printf("Aliases: `%s`\n\n", strings.Join(s.Aliases, "`, `"))
		}
		fmt.Println("### Parameters\n")
		fmt.Println("| Param | Type | Required | Default | Constraints | Description |")
		fmt.Println("|-------|------|----------|---------|-------------|-------------|")
		for _, p := range s.Params {
			req := ""
			if p.Required {
				req = "yes"
			}
			def := ""
			if p.Default != nil {
				def = fmt.Sprintf("%v", p.Default)
			}
			constraints := buildConstraints(p)
			desc := p.Desc
			if len(p.Aliases) > 0 {
				desc += " (alias: " + strings.Join(p.Aliases, "/") + ")"
			}
			fmt.Printf("| %s | %s | %s | %s | %s | %s |\n", p.Name, p.Type, req, def, constraints, desc)
		}
		fmt.Println()
	}
	return nil
}

func init() {
	schemaCmd.Flags().BoolVar(&schemaJSON, "json", false, "Output in JSON format")
	schemaCmd.Flags().BoolVar(&schemaAllLLMSTxt, "all", false, "Print all schemas (use with --llms-txt)")
	rootCmd.AddCommand(schemaCmd)
}
