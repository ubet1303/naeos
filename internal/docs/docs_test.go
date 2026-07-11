package docs

import (
	"testing"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestGenerateAPIDocs(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	endpoints := []Endpoint{
		{Method: "GET", Path: "/health", Description: "Health check"},
		{Method: "POST", Path: "/api/v1/users", Description: "Create user"},
	}
	content := gen.GenerateAPIDocs(endpoints)
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !contains(content, "/health") {
		t.Error("expected /health endpoint in docs")
	}
}

func TestGenerateArchitectureDiagram(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	content := gen.GenerateArchitectureDiagram(
		[]string{"api", "worker"},
		[]string{"core", "auth"},
	)
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !contains(content, "mermaid") {
		t.Error("expected mermaid diagram")
	}
}

func TestGenerateProjectDocs(t *testing.T) {
	gen := NewGenerator("test-project", []ArtifactRef{
		{Path: "main.go", Size: 100, Type: "go"},
	})
	content := gen.GenerateProjectDocs()
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !contains(content, "test-project") {
		t.Error("expected project name in docs")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
