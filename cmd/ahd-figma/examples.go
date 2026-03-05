package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var examplesCmd = &cobra.Command{
	Use:   "examples [category]",
	Short: "Show example batch JSON payloads",
	Long: `Show ready-to-use batch JSON examples by category.
Run without arguments to list available categories.
Run with a category name to see the example payload.

The output is valid JSON that can be piped directly to batch:
  ahd-figma examples carousel > /tmp/ops.json && ahd-figma batch /tmp/ops.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listExampleCategories()
		}
		return showExample(args[0])
	},
}

func listExampleCategories() error {
	categories := make([]string, 0, len(examplePayloads))
	for k := range examplePayloads {
		categories = append(categories, k)
	}
	sort.Strings(categories)

	fmt.Fprintln(os.Stderr, "Available example categories:")
	fmt.Fprintln(os.Stderr, "")
	for _, cat := range categories {
		desc := exampleDescriptions[cat]
		fmt.Fprintf(os.Stderr, "  %-12s  %s\n", cat, desc)
	}
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage: ahd-figma examples <category>")
	fmt.Fprintln(os.Stderr, "Pipe:  ahd-figma examples carousel > /tmp/ops.json && ahd-figma batch /tmp/ops.json")
	return nil
}

func showExample(category string) error {
	category = strings.ToLower(strings.TrimSpace(category))
	payload, ok := examplePayloads[category]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown category %q. Run 'ahd-figma examples' to list categories.\n", category)
		return fmt.Errorf("unknown example category %q", category)
	}

	// Pretty-print the JSON
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

var exampleDescriptions = map[string]string{
	"slide":      "Single slide composite with eyebrow, headline, body, bar, CTA",
	"carousel":   "3-slide Instagram carousel with strong hierarchy and semantic naming",
	"banner":     "Banner composite with divider and headline/subtitle",
	"effects":    "Shadows, glass, gradient overlay, noise, masking",
	"batch":      "Raw batch starter using token aliases (sz/padding/w/r)",
	"newsletter": "Email newsletter layout with sections",
}

var examplePayloads = map[string]interface{}{
	"slide": []interface{}{
		map[string]interface{}{
			"name":    "s1",
			"command": "slide",
			"params": map[string]interface{}{
				"canvas": "1080x1350",
				"bg":     "#0C1E2C",
				"gradient": map[string]interface{}{
					"type": "LINEAR",
					"stops": []interface{}{
						map[string]interface{}{"color": "#0C1E2C", "position": 0},
						map[string]interface{}{"color": "#1A3A5C", "position": 1},
					},
				},
				"elements": []interface{}{
					map[string]interface{}{"type": "counter", "current": "1", "total": "5", "color": "#7FBCD2"},
					map[string]interface{}{"type": "eyebrow", "text": "CATEGORY · TOPIC", "color": "#7FBCD2"},
					map[string]interface{}{"type": "headline", "text": "Your Main Headline Goes Here", "tier": "hero", "color": "#FFFFFF"},
					map[string]interface{}{"type": "bar", "color": "#FFD600"},
					map[string]interface{}{"type": "body", "text": "Supporting paragraph text that explains the key message. Keep it concise and impactful.", "color": "#B0C4D8"},
					map[string]interface{}{"type": "cta", "text": "Learn More", "bgColor": "#FFD600", "textColor": "#0C1E2C"},
					map[string]interface{}{"type": "url", "text": "yoursite.com", "color": "#7FBCD2"},
				},
			},
		},
	},

	"carousel": []interface{}{
		map[string]interface{}{
			"name":    "s1",
			"command": "slide",
			"params": map[string]interface{}{
				"canvas": "1080x1350",
				"bg":     "#111827",
				"gradient": map[string]interface{}{
					"type": "LINEAR",
					"stops": []interface{}{
						map[string]interface{}{"color": "#0B1220", "position": 0},
						map[string]interface{}{"color": "#1F2A44", "position": 1},
					},
				},
				"elements": []interface{}{
					map[string]interface{}{"type": "eyebrow", "text": "SLIDE 1 OF 3", "color": "#9CA3AF"},
					map[string]interface{}{"type": "headline", "text": "The Hook", "tier": "hero", "color": "#FFFFFF"},
					map[string]interface{}{"type": "bar", "color": "#3B82F6"},
					map[string]interface{}{"type": "body", "text": "Open with a bold statement that grabs attention immediately.", "color": "#D1D5DB"},
				},
			},
		},
		map[string]interface{}{
			"name":    "s2",
			"command": "slide",
			"params": map[string]interface{}{
				"canvas": "1080x1350",
				"bg":     "#111827",
				"elements": []interface{}{
					map[string]interface{}{"type": "eyebrow", "text": "SLIDE 2 OF 3", "color": "#9CA3AF"},
					map[string]interface{}{"type": "headline", "text": "The Details", "tier": "title", "color": "#FFFFFF"},
					map[string]interface{}{"type": "body", "text": "Support your claim with crisp proof points.", "color": "#D1D5DB"},
					map[string]interface{}{"type": "stats", "items": []interface{}{
						map[string]interface{}{"value": "10x", "label": "Faster"},
						map[string]interface{}{"value": "99%", "label": "Accuracy"},
						map[string]interface{}{"value": "24/7", "label": "Available"},
					}, "valueColor": "#3B82F6", "labelColor": "#9CA3AF"},
				},
			},
		},
		map[string]interface{}{
			"name":    "s3",
			"command": "slide",
			"params": map[string]interface{}{
				"canvas": "1080x1350",
				"bg":     "#111827",
				"elements": []interface{}{
					map[string]interface{}{"type": "eyebrow", "text": "SLIDE 3 OF 3", "color": "#9CA3AF"},
					map[string]interface{}{"type": "headline", "text": "Take Action", "tier": "hero", "color": "#FFFFFF"},
					map[string]interface{}{"type": "body", "text": "Ready to get started? Click the link below.", "color": "#D1D5DB"},
					map[string]interface{}{"type": "cta", "text": "Get Started", "bgColor": "#3B82F6", "textColor": "#FFFFFF"},
					map[string]interface{}{"type": "url", "text": "@yourhandle", "color": "#9CA3AF"},
				},
			},
		},
	},

	"banner": []interface{}{
		map[string]interface{}{
			"name":    "b1",
			"command": "banner",
			"params": map[string]interface{}{
				"canvas":       "1200x628",
				"bg":           "#FFFFFF",
				"dividerX":     480,
				"dividerColor": "#E5E7EB",
				"elements": []interface{}{
					map[string]interface{}{"type": "headline", "text": "Product Launch", "tier": "heading", "color": "#111827"},
					map[string]interface{}{"type": "subtitle", "text": "Everything you need to know about our latest release.", "color": "#6B7280"},
				},
			},
		},
	},

	"effects": []interface{}{
		map[string]interface{}{
			"_comment": "Card with layered shadows, glass effect, gradient overlay, and noise",
			"name":     "bg",
			"command":  "frame",
			"params": map[string]interface{}{
				"name": "Effects Demo", "w": 800, "h": 600, "bg": "#0F172A",
			},
		},
		map[string]interface{}{
			"name":    "card",
			"command": "frame",
			"params": map[string]interface{}{
				"name": "Glass Card", "pid": "$bg", "x": 100, "y": 100, "w": 360, "h": 400,
				"bg": "#FFFFFF20", "r": 24,
				"stroke": "#FFFFFF18", "sw": 1,
				"layoutMode": "VERTICAL", "padding": 32, "itemSpacing": 16,
			},
		},
		map[string]interface{}{
			"command": "glass",
			"params":  map[string]interface{}{"nodeId": "$card", "intensity": "medium", "tint": "#FFFFFF"},
		},
		map[string]interface{}{
			"command": "shadow",
			"params":  map[string]interface{}{"nodeId": "$card", "offsetY": 4, "radius": 12, "color": "#0000001A"},
		},
		map[string]interface{}{
			"command": "shadow",
			"params":  map[string]interface{}{"nodeId": "$card", "offsetY": 16, "radius": 48, "color": "#00000033"},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "Glass Card", "pid": "$card", "sz": 32, "fontStyle": "Bold", "color": "#FFFFFF",
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "With layered shadows and glass morphism.", "pid": "$card", "sz": 16, "color": "#FFFFFFAA", "w": 296,
			},
		},
		map[string]interface{}{
			"command": "noise",
			"params":  map[string]interface{}{"nodeId": "$bg", "noiseType": "monotone", "color": "#FFFFFF", "density": 0.15},
		},
	},

	"batch": []interface{}{
		map[string]interface{}{
			"_comment": "Raw batch starter with token aliases and semantic names",
			"name":     "root",
			"command":  "frame",
			"params": map[string]interface{}{
				"name": "Social Post", "w": 1080, "h": 1080, "bg": "#111827", "clipsContent": true,
			},
		},
		map[string]interface{}{
			"name":    "content",
			"command": "frame",
			"params": map[string]interface{}{
				"name": "Content", "pid": "$root", "w": 1080, "h": 1080,
				"noFill": true, "layoutMode": "VERTICAL", "itemSpacing": "section", "padding": "side",
				"primaryAxisAlign": "CENTER", "counterAxisAlign": "CENTER",
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"name": "Hero Title", "text": "Hello World", "pid": "$content", "sz": "hero", "fontStyle": "Bold",
				"color": "#ffffff", "textAlign": "CENTER", "w": "content", "lh": 110,
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"name": "Body Copy", "text": "Built with AHD Figma", "pid": "$content", "sz": "body",
				"color": "#9CA3AF", "textAlign": "CENTER", "w": "content", "lh": 150,
			},
		},
		map[string]interface{}{
			"name":    "cta",
			"command": "frame",
			"params": map[string]interface{}{
				"name": "CTA Button", "pid": "$content", "bg": "#3B82F6", "r": "button",
				"layoutMode": "HORIZONTAL", "paddingLeft": "card", "paddingRight": "card",
				"paddingTop": "item", "paddingBottom": "item",
				"primaryAxisSizing": "AUTO", "counterAxisSizing": "AUTO",
				"primaryAxisAlign": "CENTER", "counterAxisAlign": "CENTER",
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"name": "CTA Text", "text": "Get Started", "pid": "$cta", "sz": "cta",
				"fontStyle": "Bold", "color": "#FFFFFF",
			},
		},
	},

	"newsletter": []interface{}{
		map[string]interface{}{
			"name":    "email",
			"command": "frame",
			"params": map[string]interface{}{
				"name": "Newsletter", "w": 600, "h": 900, "bg": "#FFFFFF",
				"layoutMode": "VERTICAL", "padding": 40, "itemSpacing": 24,
				"primaryAxisAlign": "MIN", "counterAxisAlign": "CENTER",
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "Your Weekly Update", "pid": "$email", "sz": 28, "fontStyle": "Bold",
				"color": "#111827", "w": 520,
			},
		},
		map[string]interface{}{
			"name":    "divider",
			"command": "rect",
			"params": map[string]interface{}{
				"pid": "$email", "w": 520, "h": 1, "bg": "#E5E7EB",
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "Hi there! Here's what's new this week.", "pid": "$email", "sz": 16,
				"color": "#374151", "w": 520, "lh": 150,
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "Feature Highlight", "pid": "$email", "sz": 20, "fontStyle": "Bold",
				"color": "#111827", "w": 520,
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "We shipped a major update that makes everything faster and more reliable. Check it out in the app.", "pid": "$email", "sz": 16,
				"color": "#374151", "w": 520, "lh": 150,
			},
		},
		map[string]interface{}{
			"name":    "cta_btn",
			"command": "frame",
			"params": map[string]interface{}{
				"pid": "$email", "bg": "#3B82F6", "r": 8,
				"layoutMode": "HORIZONTAL", "paddingLeft": 24, "paddingRight": 24,
				"paddingTop": 12, "paddingBottom": 12,
				"primaryAxisSizing": "AUTO", "counterAxisSizing": "AUTO",
			},
		},
		map[string]interface{}{
			"command": "text",
			"params": map[string]interface{}{
				"text": "Read More", "pid": "$cta_btn", "sz": 14, "fontStyle": "Bold", "color": "#FFFFFF",
			},
		},
	},
}

func init() {
	rootCmd.AddCommand(examplesCmd)
}
