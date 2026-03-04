package designlint

import (
	"testing"
)

func TestTextTooSmall(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "node.create_frame", "params": map[string]interface{}{"width": 1080.0, "height": 1350.0, "name": "Slide"}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Hi", "parentId": "1:23", "fontSize": 14.0}},
	}
	result := Check(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "TEXT_TOO_SMALL" {
			found = true
			if w.Fix.(float64) < 30 {
				t.Errorf("expected fix >= caption tier, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected TEXT_TOO_SMALL warning")
	}
}

func TestContrastRatio(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "node.create_frame", "params": map[string]interface{}{"width": 1080.0, "height": 1350.0, "color": "#777777"}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Hi", "parentId": "1:23", "color": "#666666"}},
	}
	result := Check(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "LOW_CONTRAST" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected LOW_CONTRAST warning")
	}
}

func TestRadiusOverflow(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "shape.create_rectangle", "params": map[string]interface{}{
			"parentId": "1:23", "width": 40.0, "height": 40.0, "cornerRadius": 50.0,
		}},
	}
	result := Check(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "RADIUS_OVERFLOW" {
			found = true
			if w.Fix.(float64) != 20.0 {
				t.Errorf("expected fix=20, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected RADIUS_OVERFLOW warning")
	}
}

func TestScoreComputed(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "node.create_frame", "params": map[string]interface{}{"width": 1080.0, "height": 1350.0}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Title", "fontSize": 112.0}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Body", "fontSize": 48.0}},
	}
	result := Check(ops)
	if result.Score.Overall == 0 {
		t.Fatal("expected non-zero design score")
	}
}
