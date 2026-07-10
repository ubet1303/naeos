package ai

import "testing"

func TestZeroValue(t *testing.T) {
	var ai AI
	if ai.Models != nil {
		t.Errorf("expected nil Models, got %v", ai.Models)
	}
	if ai.Prompts != nil {
		t.Errorf("expected nil Prompts, got %v", ai.Prompts)
	}
	if ai.ContextBundles != nil {
		t.Errorf("expected nil ContextBundles, got %v", ai.ContextBundles)
	}
	if ai.Embeddings != nil {
		t.Errorf("expected nil Embeddings, got %v", ai.Embeddings)
	}
	if ai.Attributes != nil {
		t.Errorf("expected nil Attributes, got %v", ai.Attributes)
	}

	var m Model
	if m.Name != "" {
		t.Errorf("expected empty Name, got %q", m.Name)
	}
	if m.Kind != "" {
		t.Errorf("expected empty Kind, got %q", m.Kind)
	}

	var p Prompt
	if p.Name != "" {
		t.Errorf("expected empty Name, got %q", p.Name)
	}
	if p.Template != "" {
		t.Errorf("expected empty Template, got %q", p.Template)
	}

	var cb ContextBundle
	if cb.Name != "" {
		t.Errorf("expected empty Name, got %q", cb.Name)
	}
	if cb.Sources != nil {
		t.Errorf("expected nil Sources, got %v", cb.Sources)
	}

	var e Embedding
	if e.Name != "" {
		t.Errorf("expected empty Name, got %q", e.Name)
	}
	if e.Dimension != 0 {
		t.Errorf("expected zero Dimension, got %d", e.Dimension)
	}
}

func TestInitialization(t *testing.T) {
	ai := AI{
		Models:  []Model{{Name: "gpt-4", Kind: "llm", Version: "v1"}},
		Prompts: []Prompt{{Name: "system", Template: "You are...", Kind: "system"}},
		ContextBundles: []ContextBundle{
			{Name: "docs", Sources: []string{"file1.md", "file2.md"}},
		},
		Embeddings: []Embedding{
			{Name: "ada", Dimension: 1536},
		},
		Attributes: map[string]string{"env": "prod"},
	}

	if len(ai.Models) != 1 || ai.Models[0].Name != "gpt-4" {
		t.Errorf("unexpected Models: %v", ai.Models)
	}
	if len(ai.Prompts) != 1 || ai.Prompts[0].Template != "You are..." {
		t.Errorf("unexpected Prompts: %v", ai.Prompts)
	}
	if ai.ContextBundles[0].Sources[0] != "file1.md" {
		t.Errorf("unexpected Sources: %v", ai.ContextBundles[0].Sources)
	}
	if ai.Embeddings[0].Dimension != 1536 {
		t.Errorf("unexpected Dimension: %d", ai.Embeddings[0].Dimension)
	}
	if ai.Attributes["env"] != "prod" {
		t.Errorf("unexpected Attributes: %v", ai.Attributes)
	}
}

func TestModelNameUniqueness(t *testing.T) {
	models := []Model{
		{Name: "gpt-4"},
		{Name: "claude-3"},
		{Name: "gpt-4"},
	}
	names := map[string]int{}
	for _, m := range models {
		names[m.Name]++
	}
	if names["gpt-4"] != 2 {
		t.Errorf("expected 2 gpt-4 models, got %d", names["gpt-4"])
	}
}
