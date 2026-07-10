package module

import (
	"testing"
)

func TestModuleZeroValue(t *testing.T) {
	m := &Module{}
	if m.Name != "" {
		t.Errorf("zero-value Module.Name = %q, want empty", m.Name)
	}
	if m.Path != "" {
		t.Errorf("zero-value Module.Path = %q, want empty", m.Path)
	}
}

func TestModuleWithFields(t *testing.T) {
	m := &Module{
		Name:         "auth",
		Path:         "./internal/auth",
		Description:  "Authentication module",
		Packages:     []string{"auth", "auth/domain"},
		Dependencies: []string{"user", "crypto"},
		Attributes:   map[string]string{"layer": "domain"},
	}
	if m.Name != "auth" {
		t.Errorf("Name = %q, want %q", m.Name, "auth")
	}
	if m.Path != "./internal/auth" {
		t.Errorf("Path = %q, want %q", m.Path, "./internal/auth")
	}
	if len(m.Packages) != 2 {
		t.Errorf("Packages has %d entries, want 2", len(m.Packages))
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("Dependencies has %d entries, want 2", len(m.Dependencies))
	}
	if m.Attributes["layer"] != "domain" {
		t.Errorf("Attributes[layer] = %q, want %q", m.Attributes["layer"], "domain")
	}
}
