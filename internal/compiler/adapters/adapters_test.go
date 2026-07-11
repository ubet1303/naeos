package adapters

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/compiler"
	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/architecture"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/service"
)

func testNEIR() *model.NEIR {
	return &model.NEIR{
		Project: &project.Project{
			Name:        "test-project",
			Description: "A test project",
			Version:     "1.0.0",
		},
		Architecture: &architecture.Architecture{
			Pattern:    "hexagonal",
			Principles: []string{"DI", "SRP"},
		},
		Modules: []module.Module{
			{Name: "core", Path: "./core", Description: "Core module", Dependencies: []string{}},
			{Name: "api", Path: "./api", Description: "API module", Dependencies: []string{"core"}},
		},
		Services: []service.Service{
			{Name: "http-api", Kind: "http", Port: 8080},
		},
	}
}

func TestCopilotAdapter(t *testing.T) {
	a := NewCopilotAdapter()
	if a.Target() != compiler.TargetCopilot {
		t.Errorf("expected copilot target, got %s", a.Target())
	}

	out, err := a.Compile(testNEIR())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(out.Files))
	}
	if out.Target != compiler.TargetCopilot {
		t.Errorf("expected copilot target in output")
	}
}

func TestCopilotAdapterNil(t *testing.T) {
	a := NewCopilotAdapter()
	_, err := a.Compile(nil)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestClaudeAdapter(t *testing.T) {
	a := NewClaudeAdapter()
	if a.Target() != compiler.TargetClaude {
		t.Errorf("expected claude target, got %s", a.Target())
	}

	out, err := a.Compile(testNEIR())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(out.Files))
	}
}

func TestClaudeAdapterNil(t *testing.T) {
	a := NewClaudeAdapter()
	_, err := a.Compile(nil)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestCursorAdapter(t *testing.T) {
	a := NewCursorAdapter()
	if a.Target() != compiler.TargetCursor {
		t.Errorf("expected cursor target, got %s", a.Target())
	}

	out, err := a.Compile(testNEIR())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(out.Files))
	}
}

func TestCursorAdapterNil(t *testing.T) {
	a := NewCursorAdapter()
	_, err := a.Compile(nil)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestGeminiAdapter(t *testing.T) {
	a := NewGeminiAdapter()
	if a.Target() != compiler.TargetGemini {
		t.Errorf("expected gemini target, got %s", a.Target())
	}

	out, err := a.Compile(testNEIR())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(out.Files))
	}
}

func TestGeminiAdapterNil(t *testing.T) {
	a := NewGeminiAdapter()
	_, err := a.Compile(nil)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestCodexAdapter(t *testing.T) {
	a := NewCodexAdapter()
	if a.Target() != compiler.TargetCodex {
		t.Errorf("expected codex target, got %s", a.Target())
	}

	out, err := a.Compile(testNEIR())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(out.Files))
	}
}

func TestCodexAdapterNil(t *testing.T) {
	a := NewCodexAdapter()
	_, err := a.Compile(nil)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestOpenCodeAdapter(t *testing.T) {
	a := NewOpenCodeAdapter()
	if a.Target() != compiler.TargetOpenCode {
		t.Errorf("expected opencode target, got %s", a.Target())
	}

	out, err := a.Compile(testNEIR())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(out.Files))
	}
}

func TestOpenCodeAdapterNil(t *testing.T) {
	a := NewOpenCodeAdapter()
	_, err := a.Compile(nil)
	if err == nil {
		t.Fatal("expected error for nil NEIR")
	}
}

func TestAllAdaptersContent(t *testing.T) {
	adapters := []compiler.Adapter{
		NewCopilotAdapter(),
		NewClaudeAdapter(),
		NewCursorAdapter(),
		NewGeminiAdapter(),
		NewCodexAdapter(),
		NewOpenCodeAdapter(),
	}

	neir := testNEIR()
	for _, a := range adapters {
		out, err := a.Compile(neir)
		if err != nil {
			t.Fatalf("adapter %s: %v", a.Target(), err)
		}
		if out.Summary == "" {
			t.Errorf("adapter %s: empty summary", a.Target())
		}
		for _, f := range out.Files {
			if f.Content == "" {
				t.Errorf("adapter %s: empty file content for %s", a.Target(), f.Path)
			}
		}
	}
}
