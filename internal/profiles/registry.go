package profiles

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed profiles.json
var builtinProfilesJSON []byte

type Profile struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Industry    string            `json:"industry"`
	Version     string            `json:"version"`
	Modules     []ModuleTemplate  `json:"modules"`
	Services    []ServiceTemplate `json:"services"`
	Architecture ArchTemplate     `json:"architecture"`
	Security    SecurityTemplate  `json:"security"`
	Deployment  DeployTemplate    `json:"deployment"`
	Testing     TestTemplate      `json:"testing"`
	Tags        []string          `json:"tags"`
}

type ModuleTemplate struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type ServiceTemplate struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Port     int    `json:"port"`
	Description string `json:"description,omitempty"`
}

type ArchTemplate struct {
	Pattern    string   `json:"pattern"`
	Principles []string `json:"principles,omitempty"`
}

type SecurityTemplate struct {
	Authentication string   `json:"authentication,omitempty"`
	Authorization  string   `json:"authorization,omitempty"`
	Roles          []string `json:"roles,omitempty"`
	Encryption     bool     `json:"encryption,omitempty"`
}

type DeployTemplate struct {
	Strategy     string   `json:"strategy"`
	Environments []string `json:"environments,omitempty"`
}

type TestTemplate struct {
	Strategy   string   `json:"strategy"`
	Coverage   string   `json:"coverage,omitempty"`
	Frameworks []string `json:"frameworks,omitempty"`
}

type Registry struct {
	profiles map[string]*Profile
}

func NewRegistry() *Registry {
	r := &Registry{
		profiles: make(map[string]*Profile),
	}
	r.loadBuiltin()
	return r
}

func (r *Registry) loadBuiltin() {
	var profiles []Profile
	if err := json.Unmarshal(builtinProfilesJSON, &profiles); err != nil {
		return
	}
	for i := range profiles {
		r.profiles[profiles[i].ID] = &profiles[i]
	}
}

func (r *Registry) Register(p *Profile) {
	r.profiles[p.ID] = p
}

func (r *Registry) Get(id string) (*Profile, bool) {
	p, ok := r.profiles[id]
	return p, ok
}

func (r *Registry) List() []Profile {
	result := make([]Profile, 0, len(r.profiles))
	for _, p := range r.profiles {
		result = append(result, *p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (r *Registry) Search(query string) []Profile {
	query = strings.ToLower(query)
	var result []Profile
	for _, p := range r.profiles {
		if strings.Contains(strings.ToLower(p.Name), query) ||
			strings.Contains(strings.ToLower(p.Description), query) ||
			strings.Contains(strings.ToLower(p.Industry), query) ||
			strings.Contains(strings.ToLower(p.ID), query) {
			result = append(result, *p)
		}
	}
	return result
}

func (r *Registry) ByIndustry(industry string) []Profile {
	var result []Profile
	for _, p := range r.profiles {
		if strings.EqualFold(p.Industry, industry) {
			result = append(result, *p)
		}
	}
	return result
}

func (r *Registry) ToSpecYAML(p *Profile) string {
	var sb strings.Builder

	slug := strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
	sb.WriteString(fmt.Sprintf("project: %s\n", slug))
	sb.WriteString(fmt.Sprintf("description: %s\n\n", p.Description))

	sb.WriteString("modules:\n")
	for _, m := range p.Modules {
		sb.WriteString(fmt.Sprintf("  - name: %s\n    path: %s\n    description: %s\n", m.Name, m.Path, m.Description))
		if len(m.Dependencies) > 0 {
			sb.WriteString("    dependencies:\n")
			for _, d := range m.Dependencies {
				sb.WriteString(fmt.Sprintf("      - %s\n", d))
			}
		}
	}

	sb.WriteString("\nservices:\n")
	for _, s := range p.Services {
		sb.WriteString(fmt.Sprintf("  - name: %s\n    kind: %s\n    port: %d\n", s.Name, s.Kind, s.Port))
		if s.Description != "" {
			sb.WriteString(fmt.Sprintf("    description: %s\n", s.Description))
		}
	}

	sb.WriteString(fmt.Sprintf("\narchitecture:\n  pattern: %s\n", p.Architecture.Pattern))
	if len(p.Architecture.Principles) > 0 {
		sb.WriteString("  principles:\n")
		for _, pr := range p.Architecture.Principles {
			sb.WriteString(fmt.Sprintf("    - %s\n", pr))
		}
	}

	if p.Security.Authentication != "" {
		sb.WriteString(fmt.Sprintf("\nsecurity:\n  authentication:\n    method: %s\n", p.Security.Authentication))
		if p.Security.Authorization != "" {
			sb.WriteString(fmt.Sprintf("  authorization:\n    model: %s\n", p.Security.Authorization))
		}
		if len(p.Security.Roles) > 0 {
			sb.WriteString("    roles:\n")
			for _, role := range p.Security.Roles {
				sb.WriteString(fmt.Sprintf("      - %s\n", role))
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\ndeployment:\n  strategy: %s\n", p.Deployment.Strategy))
	if len(p.Deployment.Environments) > 0 {
		sb.WriteString("  environments:\n")
		for _, env := range p.Deployment.Environments {
			sb.WriteString(fmt.Sprintf("    - %s\n", env))
		}
	}

	sb.WriteString(fmt.Sprintf("\ntesting:\n  strategy: %s\n", p.Testing.Strategy))
	if p.Testing.Coverage != "" {
		sb.WriteString(fmt.Sprintf("  coverage: %s\n", p.Testing.Coverage))
	}
	if len(p.Testing.Frameworks) > 0 {
		sb.WriteString("  frameworks:\n")
		for _, f := range p.Testing.Frameworks {
			sb.WriteString(fmt.Sprintf("    - %s\n", f))
		}
	}

	return sb.String()
}
