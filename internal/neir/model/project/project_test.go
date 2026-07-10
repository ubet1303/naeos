package project

import (
	"testing"
)

func TestProjectZeroValue(t *testing.T) {
	p := &Project{}
	if p.Name != "" {
		t.Errorf("zero-value Project.Name = %q, want empty", p.Name)
	}
	if p.Version != "" {
		t.Errorf("zero-value Project.Version = %q, want empty", p.Version)
	}
}

func TestProjectWithFields(t *testing.T) {
	p := &Project{
		Name:        "my-api",
		Version:     "1.0.0",
		Description: "A test API",
		License:     "Apache-2.0",
		Authors:     []string{"alice", "bob"},
		Repository:  "github.com/test/my-api",
		Tags:        []string{"api", "backend"},
		Attributes:  map[string]string{"env": "production"},
	}
	if p.Name != "my-api" {
		t.Errorf("Name = %q, want %q", p.Name, "my-api")
	}
	if p.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", p.Version, "1.0.0")
	}
	if len(p.Authors) != 2 {
		t.Errorf("Authors has %d entries, want 2", len(p.Authors))
	}
	if p.Attributes["env"] != "production" {
		t.Errorf("Attributes[env] = %q, want %q", p.Attributes["env"], "production")
	}
}
