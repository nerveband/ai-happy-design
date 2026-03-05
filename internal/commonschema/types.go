package commonschema

// Param defines a single command parameter.
type Param struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Required     bool     `json:"required,omitempty"`
	Aliases      []string `json:"aliases,omitempty"`
	Enum         []string `json:"enum,omitempty"`
	Pattern      string   `json:"pattern,omitempty"`
	Default      any      `json:"default,omitempty"`
	Minimum      *float64 `json:"minimum,omitempty"`
	Maximum      *float64 `json:"maximum,omitempty"`
	SafePath     bool     `json:"safePath,omitempty"`
	LowRiskFuzzy bool     `json:"lowRiskFuzzy,omitempty"`
}

// Command defines the public schema for a CLI command.
type Command struct {
	Name           string   `json:"command"`
	Description    string   `json:"description"`
	Domain         string   `json:"domain"`
	Mutating       bool     `json:"mutating,omitempty"`
	PluginRequired bool     `json:"pluginRequired,omitempty"`
	Examples       []string `json:"examples,omitempty"`
	Params         []Param  `json:"params"`
}

// Ptr returns a float64 pointer.
func Ptr(v float64) *float64 {
	return &v
}
