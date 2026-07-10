package architecture

import "testing"

func TestPatternConstants(t *testing.T) {
	tests := []struct {
		constant Pattern
		expected string
	}{
		{PatternLayered, "layered"},
		{PatternClean, "clean"},
		{PatternHexagonal, "hexagonal"},
		{PatternMicrokernel, "microkernel"},
		{PatternEventDriven, "event-driven"},
		{PatternCQRS, "cqrs"},
		{PatternMonolith, "monolith"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("Pattern %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var a Architecture
	if a.Pattern != "" {
		t.Errorf("expected empty Pattern, got %q", a.Pattern)
	}
	if a.Style != "" {
		t.Errorf("expected empty Style, got %q", a.Style)
	}
	if a.Layers != nil {
		t.Errorf("expected nil Layers, got %v", a.Layers)
	}
	if a.Principles != nil {
		t.Errorf("expected nil Principles, got %v", a.Principles)
	}

	var l Layer
	if l.Name != "" {
		t.Errorf("expected empty Name, got %q", l.Name)
	}
	if l.Modules != nil {
		t.Errorf("expected nil Modules, got %v", l.Modules)
	}
}

func TestInitialization(t *testing.T) {
	a := Architecture{
		Pattern: PatternHexagonal,
		Style:   "ports-and-adapters",
		Layers: []Layer{
			{Name: "domain", Modules: []string{"core"}},
			{Name: "adapter", Modules: []string{"http", "db"}},
		},
		Principles: []string{"separation of concerns", "dependency inversion"},
		Attributes: map[string]string{"team": "platform"},
	}

	if a.Pattern != PatternHexagonal {
		t.Errorf("expected Pattern %q, got %q", PatternHexagonal, a.Pattern)
	}
	if len(a.Layers) != 2 {
		t.Errorf("expected 2 layers, got %d", len(a.Layers))
	}
	if a.Layers[0].Modules[0] != "core" {
		t.Errorf("expected first module 'core', got %q", a.Layers[0].Modules[0])
	}
	if a.Principles[0] != "separation of concerns" {
		t.Errorf("unexpected principles: %v", a.Principles)
	}
}
