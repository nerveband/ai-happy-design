package commonschema

import (
	"fmt"
	"sort"
	"strings"
)

var registry []Command

// Register adds a command schema to the global registry.
func Register(command Command) {
	registry = append(registry, command)
}

// All returns the registered schemas sorted by command name.
func All() []Command {
	out := make([]Command, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// Lookup finds a command schema by its domain.action name.
func Lookup(name string) *Command {
	needle := strings.ToLower(strings.TrimSpace(name))
	for i := range registry {
		if strings.ToLower(registry[i].Name) == needle {
			return &registry[i]
		}
		for _, alias := range registry[i].Aliases {
			if strings.ToLower(alias) == needle {
				return &registry[i]
			}
		}
	}
	return nil
}

// LookupParam resolves a parameter by canonical name or alias.
func LookupParam(command *Command, name string) *Param {
	if command == nil {
		return nil
	}
	needle := strings.ToLower(strings.TrimSpace(name))
	for i := range command.Params {
		if strings.ToLower(command.Params[i].Name) == needle {
			return &command.Params[i]
		}
		for _, alias := range command.Params[i].Aliases {
			if strings.ToLower(alias) == needle {
				return &command.Params[i]
			}
		}
	}
	return nil
}

// Domains returns the distinct registered domains.
func Domains() []string {
	set := map[string]struct{}{}
	for _, command := range registry {
		set[command.Domain] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for domain := range set {
		out = append(out, domain)
	}
	sort.Strings(out)
	return out
}

// LLMSDoc returns a markdown reference for all registered commands.
func LLMSDoc(productName, binaryName string) string {
	lines := []string{
		fmt.Sprintf("# %s Command Reference", productName),
		"",
		fmt.Sprintf("Auto-generated from CLI schemas. Use `%s schema <command> --json` for machine-readable output.", binaryName),
		"",
	}
	for _, command := range All() {
		lines = append(lines, fmt.Sprintf("## %s", command.Name))
		lines = append(lines, "")
		lines = append(lines, command.Description)
		lines = append(lines, "")
		lines = append(lines, "| Param | Type | Required | Notes |")
		lines = append(lines, "| --- | --- | --- | --- |")
		for _, param := range command.Params {
			notes := param.Description
			if len(param.Enum) > 0 {
				notes += " enum=" + strings.Join(param.Enum, "/")
			}
			if param.Pattern != "" {
				notes += " pattern=" + param.Pattern
			}
			required := "no"
			if param.Required {
				required = "yes"
			}
			lines = append(lines, fmt.Sprintf("| %s | %s | %s | %s |", param.Name, param.Type, required, notes))
		}
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
