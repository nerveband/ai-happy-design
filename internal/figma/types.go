package figma

// Color represents an RGBA color in Figma (0-1 range).
type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

// GradientStop represents a stop in a gradient fill.
type GradientStop struct {
	Position float64 `json:"position"`
	Color    Color   `json:"color"`
}

// NodeInfo contains basic information about a Figma node.
type NodeInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Vector2D represents a 2D point or size.
type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Effect represents a visual effect (shadow, blur, etc.).
type Effect struct {
	Type    string  `json:"type"`
	Visible bool    `json:"visible"`
	Radius  float64 `json:"radius,omitempty"`
	Color   *Color  `json:"color,omitempty"`
	OffsetX float64 `json:"offsetX,omitempty"`
	OffsetY float64 `json:"offsetY,omitempty"`
	Spread  float64 `json:"spread,omitempty"`
}

// LayoutConstraints represents auto-layout constraints.
type LayoutConstraints struct {
	Horizontal string `json:"horizontal,omitempty"`
	Vertical   string `json:"vertical,omitempty"`
}

// ExportSettings defines how a node should be exported.
type ExportSettings struct {
	Format     string  `json:"format"`
	Scale      float64 `json:"scale,omitempty"`
	Constraint string  `json:"constraint,omitempty"`
}
