package profiles

import (
	"strings"
	"testing"
)

func TestRegistryLoadBuiltin(t *testing.T) {
	reg := NewRegistry()
	list := reg.List()
	if len(list) < 5 {
		t.Errorf("expected at least 5 builtin profiles, got %d", len(list))
	}
}

func TestRegistryGet(t *testing.T) {
	reg := NewRegistry()

	p, ok := reg.Get("saas")
	if !ok {
		t.Fatal("expected to find saas profile")
	}
	if p.Name != "SaaS Application" {
		t.Errorf("expected SaaS Application, got %q", p.Name)
	}
	if p.Industry != "technology" {
		t.Errorf("expected technology industry, got %q", p.Industry)
	}
}

func TestRegistryGetNotFound(t *testing.T) {
	reg := NewRegistry()
	_, ok := reg.Get("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRegistryRegister(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&Profile{
		ID:          "custom",
		Name:        "Custom Profile",
		Description: "A custom profile",
		Industry:    "custom",
	})

	p, ok := reg.Get("custom")
	if !ok {
		t.Fatal("expected to find custom profile")
	}
	if p.Name != "Custom Profile" {
		t.Errorf("expected Custom Profile, got %q", p.Name)
	}
}

func TestRegistrySearch(t *testing.T) {
	reg := NewRegistry()
	results := reg.Search("saas")
	if len(results) == 0 {
		t.Fatal("expected to find saas profile")
	}

	found := false
	for _, p := range results {
		if strings.Contains(strings.ToLower(p.Name), "saas") {
			found = true
		}
	}
	if !found {
		t.Error("expected saas in search results")
	}
}

func TestRegistryByIndustry(t *testing.T) {
	reg := NewRegistry()
	aiProfiles := reg.ByIndustry("artificial-intelligence")
	if len(aiProfiles) == 0 {
		t.Fatal("expected at least 1 AI profile")
	}
	for _, p := range aiProfiles {
		if p.Industry != "artificial-intelligence" {
			t.Errorf("expected ai industry, got %q", p.Industry)
		}
	}
}

func TestRegistryByIndustryEmpty(t *testing.T) {
	reg := NewRegistry()
	results := reg.ByIndustry("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestToSpecYAML(t *testing.T) {
	reg := NewRegistry()
	p, ok := reg.Get("saas")
	if !ok {
		t.Fatal("expected saas profile")
	}

	spec := reg.ToSpecYAML(p)
	if spec == "" {
		t.Error("expected non-empty spec")
	}
	if !strings.Contains(spec, "modules:") {
		t.Error("expected modules section")
	}
	if !strings.Contains(spec, "services:") {
		t.Error("expected services section")
	}
	if !strings.Contains(spec, "architecture:") {
		t.Error("expected architecture section")
	}
	if !strings.Contains(spec, "testing:") {
		t.Error("expected testing section")
	}
}

func TestProfileFields(t *testing.T) {
	reg := NewRegistry()

	tests := []struct {
		id       string
		modules  int
		services int
	}{
		{"saas", 6, 2},
		{"ai-agent", 7, 3},
		{"fintech", 7, 3},
		{"healthcare", 7, 2},
		{"government", 7, 3},
	}

	for _, tt := range tests {
		p, ok := reg.Get(tt.id)
		if !ok {
			t.Errorf("profile %q not found", tt.id)
			continue
		}
		if len(p.Modules) != tt.modules {
			t.Errorf("profile %s: expected %d modules, got %d", tt.id, tt.modules, len(p.Modules))
		}
		if len(p.Services) != tt.services {
			t.Errorf("profile %s: expected %d services, got %d", tt.id, tt.services, len(p.Services))
		}
	}
}
