package figma

// ProgressUpdate represents a progress notification from the Figma plugin.
type ProgressUpdate struct {
	CommandID string  `json:"commandId"`
	Progress  float64 `json:"progress"` // 0.0 to 1.0
	Message   string  `json:"message,omitempty"`
	Done      bool    `json:"done"`
}
