package contract

import (
	"fmt"
	"sort"

	"github.com/nerveband/ai-happy-design/internal/schema"
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

type Finding struct {
	Code    string
	Command string
	Detail  string
}

type CommandSurface struct {
	Command     string
	Domain      string
	Action      string
	Aliases     []string
	ReadOnly    bool
	Destructive bool
	Source      string
}

func FromSchemas() []CommandSurface {
	out := make([]CommandSurface, 0, len(schema.All))
	for _, s := range schema.All {
		domain, action, _ := schema.SplitCommand(s.Command)
		out = append(out, CommandSurface{
			Command:     s.Command,
			Domain:      domain,
			Action:      action,
			Aliases:     append([]string{}, s.Aliases...),
			ReadOnly:    s.Safety == "read" || s.Safety == "local",
			Destructive: s.Safety == "destructive",
			Source:      "schema",
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Command < out[j].Command })
	return out
}

func FromToolCatalog() []CommandSurface {
	catalog := tools.ToolCatalog()
	out := make([]CommandSurface, 0)
	for domain, actions := range catalog {
		for action := range actions {
			out = append(out, CommandSurface{
				Command: domain + "." + action,
				Domain:  domain,
				Action:  action,
				Source:  "tool_catalog",
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Command < out[j].Command })
	return out
}

func Diff(schemaSurface, catalogSurface []CommandSurface) []Finding {
	schemaCommands := make(map[string]CommandSurface, len(schemaSurface))
	catalogCommands := make(map[string]CommandSurface, len(catalogSurface))
	for _, surface := range schemaSurface {
		schemaCommands[surface.Command] = surface
	}
	for _, surface := range catalogSurface {
		catalogCommands[surface.Command] = surface
	}
	var findings []Finding
	for command := range schemaCommands {
		if _, ok := catalogCommands[command]; !ok {
			findings = append(findings, Finding{Code: "MISSING_TOOL_CATALOG_ENTRY", Command: command})
		}
	}
	for command := range catalogCommands {
		if _, ok := schemaCommands[command]; !ok {
			findings = append(findings, Finding{Code: "TOOL_CATALOG_WITHOUT_SCHEMA", Command: command})
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Code == findings[j].Code {
			return findings[i].Command < findings[j].Command
		}
		return findings[i].Code < findings[j].Code
	})
	return findings
}

func DiscoveryFindings() []Finding {
	findings := Diff(FromSchemas(), FromToolCatalog())
	schemaCommands := commandSetFromSchemas()

	for command := range schemaCommands {
		domain, action, err := ws.ResolveCommandRoute(command, nil)
		if err != nil {
			findings = append(findings, Finding{Code: "UNROUTABLE_SCHEMA_COMMAND", Command: command, Detail: err.Error()})
			continue
		}
		expectedDomain, expectedAction, ok := schema.SplitCommand(command)
		if !ok {
			findings = append(findings, Finding{Code: "INVALID_SCHEMA_COMMAND", Command: command})
			continue
		}
		if domain != expectedDomain || action != expectedAction {
			findings = append(findings, Finding{
				Code:    "ROUTE_MISMATCH",
				Command: command,
				Detail:  fmt.Sprintf("expected %s.%s got %s.%s", expectedDomain, expectedAction, domain, action),
			})
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Code == findings[j].Code {
			return findings[i].Command < findings[j].Command
		}
		return findings[i].Code < findings[j].Code
	})
	return findings
}

func commandSetFromSchemas() map[string]struct{} {
	out := make(map[string]struct{}, len(schema.All))
	for _, s := range schema.All {
		out[s.Command] = struct{}{}
	}
	return out
}

func commandSetFromTools() map[string]struct{} {
	catalog := tools.ToolCatalog()
	out := make(map[string]struct{})
	for domain, actions := range catalog {
		for action := range actions {
			out[domain+"."+action] = struct{}{}
		}
	}
	return out
}
