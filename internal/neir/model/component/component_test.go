package component

import "testing"

func TestComponentKindConstants(t *testing.T) {
	tests := []struct {
		constant ComponentKind
		expected string
	}{
		{KindHandler, "handler"},
		{KindService, "service"},
		{KindRepository, "repository"},
		{KindMiddleware, "middleware"},
		{KindModel, "model"},
		{KindConfig, "config"},
		{KindWorker, "worker"},
		{KindScheduler, "scheduler"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("ComponentKind %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var c Component
	if c.Name != "" {
		t.Errorf("expected empty Name, got %q", c.Name)
	}
	if c.Kind != "" {
		t.Errorf("expected empty Kind, got %q", c.Kind)
	}
	if c.Dependencies != nil {
		t.Errorf("expected nil Dependencies, got %v", c.Dependencies)
	}
	if c.Attributes != nil {
		t.Errorf("expected nil Attributes, got %v", c.Attributes)
	}
}

func TestInitialization(t *testing.T) {
	c := Component{
		Name:         "user-service",
		Kind:         KindService,
		Module:       "internal/users",
		Description:  "Handles user CRUD",
		Dependencies: []string{"db-repo", "cache-middleware"},
		Attributes:   map[string]string{"lang": "go"},
	}

	if c.Kind != KindService {
		t.Errorf("expected Kind %q, got %q", KindService, c.Kind)
	}
	if len(c.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(c.Dependencies))
	}
	if c.Dependencies[0] != "db-repo" {
		t.Errorf("expected first dependency 'db-repo', got %q", c.Dependencies[0])
	}
}
