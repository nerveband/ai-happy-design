package schema

// Param defines a single parameter for a command schema.
type Param struct {
	Name           string      `json:"name"`
	Type           string      `json:"type"` // "string", "number", "boolean", "array", "object"
	Required       bool        `json:"required,omitempty"`
	Aliases        []string    `json:"aliases,omitempty"`
	Desc           string      `json:"description"`
	Enum           []string    `json:"enum,omitempty"`
	Min            *float64    `json:"minimum,omitempty"`
	Max            *float64    `json:"maximum,omitempty"`
	Default        interface{} `json:"default,omitempty"`
	Pattern        string      `json:"pattern,omitempty"`
	SemanticTokens bool        `json:"semanticTokens,omitempty"` // allows token names like "hero", "body"
	AutoFix        string      `json:"autoFix,omitempty"`        // e.g. "lineHeightUnit:PERCENT"
}

// Schema defines the full schema for a command.action pair.
type Schema struct {
	Command        string   `json:"command"`
	Aliases        []string `json:"aliases,omitempty"`
	Description    string   `json:"description"`
	Params         []Param  `json:"params"`
	Safety         string   `json:"safety,omitempty"`      // "read", "write", "destructive", "local"
	Idempotency    string   `json:"idempotency,omitempty"` // "idempotent", "non_idempotent", "unknown"
	SupportsDryRun bool     `json:"supportsDryRun"`
	RequiresFigma  bool     `json:"requiresFigma,omitempty"`
	RequiresRelay  bool     `json:"requiresRelay,omitempty"`
	RequiresAuth   bool     `json:"requiresAuth,omitempty"`
}

// Ptr returns a pointer to a float64 (helper for Min/Max).
func Ptr(v float64) *float64 { return &v }
