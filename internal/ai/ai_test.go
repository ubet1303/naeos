package ai

import (
	"testing"
)

func TestSuggestEmpty(t *testing.T) {
	svc := NewService()
	_, err := svc.Suggest("")
	if err == nil {
		t.Fatal("expected error for empty spec")
	}
}

func TestSuggestMissingArchitecture(t *testing.T) {
	svc := NewService()
	suggestions, err := svc.Suggest("project: test\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, s := range suggestions {
		if s.Category == "architecture" {
			found = true
		}
	}
	if !found {
		t.Error("expected architecture suggestion")
	}
}

func TestSuggestCompleteSpec(t *testing.T) {
	svc := NewService()
	spec := "project: test\narchitecture:\n  pattern: hexagonal\ndeployment:\n  strategy: rolling\ntesting:\n  strategy: unit\ndescription: A test project\n"
	suggestions, err := svc.Suggest(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) == 0 {
		t.Error("expected at least one suggestion")
	}
}

func TestExplainPipeline(t *testing.T) {
	svc := NewService()
	exp, err := svc.Explain("pipeline", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exp.Content == "" {
		t.Error("expected non-empty content")
	}
	if len(exp.Details) == 0 {
		t.Error("expected non-empty details")
	}
}

func TestExplainNEIR(t *testing.T) {
	svc := NewService()
	exp, err := svc.Explain("neir", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exp.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestExplainEmptyTopic(t *testing.T) {
	svc := NewService()
	_, err := svc.Explain("", "")
	if err == nil {
		t.Fatal("expected error for empty topic")
	}
}

func TestExplainUnknownTopic(t *testing.T) {
	svc := NewService()
	exp, err := svc.Explain("unknown-topic", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exp.Content == "" {
		t.Error("expected non-empty content for unknown topic")
	}
}
