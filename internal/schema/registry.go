package schema

import "strings"

// All registered schemas. Populated by init() in each schema file.
var All []Schema

// Register adds a schema to the global registry.
func Register(s Schema) { All = append(All, s) }

// Lookup finds a schema by command name or alias. Returns nil if not found.
func Lookup(command string) *Schema {
	cmd := strings.ToLower(command)
	for i := range All {
		if strings.ToLower(All[i].Command) == cmd {
			return &All[i]
		}
		for _, alias := range All[i].Aliases {
			if strings.ToLower(alias) == cmd {
				return &All[i]
			}
		}
	}
	return nil
}

// LookupParam finds a param by name or alias within a schema.
func LookupParam(s *Schema, name string) *Param {
	lower := strings.ToLower(name)
	for i := range s.Params {
		if strings.ToLower(s.Params[i].Name) == lower {
			return &s.Params[i]
		}
		for _, alias := range s.Params[i].Aliases {
			if strings.ToLower(alias) == lower {
				return &s.Params[i]
			}
		}
	}
	return nil
}

// Commands returns all registered command names.
func Commands() []string {
	out := make([]string, len(All))
	for i, s := range All {
		out[i] = s.Command
	}
	return out
}
