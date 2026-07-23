package parser

import (
	"fmt"
	"strings"
)

// NEIR v2.0 — Conditional Modules
// Modules can have a `when` condition evaluated against environment variables.

type Condition struct {
	Env    string `json:"env,omitempty"`
	Equals string `json:"equals,omitempty"`
	Not    bool   `json:"not,omitempty"`
}

type ConditionalModule struct {
	Module    Module    `json:"module"`
	Condition Condition `json:"when"`
}

// EvaluateCondition checks whether a condition is satisfied given a set of env vars.
func EvaluateCondition(cond Condition, envVars map[string]string) bool {
	if cond.Env == "" {
		return true
	}
	val, ok := envVars[cond.Env]
	if !ok {
		return cond.Not
	}
	if cond.Equals == "" {
		return !cond.Not
	}
	match := val == cond.Equals
	if cond.Not {
		return !match
	}
	return match
}

// FilterConditionalModules returns only modules whose conditions are satisfied.
func FilterConditionalModules(modules []ConditionalModule, envVars map[string]string) []Module {
	var result []Module
	for _, cm := range modules {
		if EvaluateCondition(cm.Condition, envVars) {
			result = append(result, cm.Module)
		}
	}
	return result
}

// NEIR v2.0 — Environment Profiles
// A profile overrides module/service config per environment.

type EnvironmentProfile struct {
	Name     string            `json:"name"`
	Modules  []ModuleOverride  `json:"modules,omitempty"`
	Services []ServiceOverride `json:"services,omitempty"`
}

type ModuleOverride struct {
	Name         string   `json:"name"`
	Path         string   `json:"path,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	Enabled      *bool    `json:"enabled,omitempty"`
}

type ServiceOverride struct {
	Name string `json:"name"`
	Port *int   `json:"port,omitempty"`
}

// ApplyProfiles applies environment-specific overrides to modules and services.
func ApplyProfiles(modules []Module, services []Service, profiles []EnvironmentProfile, activeProfile string) ([]Module, []Service) {
	if activeProfile == "" {
		return modules, services
	}

	for _, p := range profiles {
		if p.Name != activeProfile {
			continue
		}

		moduleMap := make(map[string]int)
		for i, m := range modules {
			moduleMap[m.Name] = i
		}
		for _, mo := range p.Modules {
			if idx, ok := moduleMap[mo.Name]; ok {
				if mo.Path != "" {
					modules[idx].Path = mo.Path
				}
				if len(mo.Dependencies) > 0 {
					modules[idx].Dependencies = mo.Dependencies
				}
				if mo.Enabled != nil && !*mo.Enabled {
					modules = append(modules[:idx], modules[idx+1:]...)
				}
			}
		}

		svcMap := make(map[string]int)
		for i, s := range services {
			svcMap[s.Name] = i
		}
		for _, so := range p.Services {
			if idx, ok := svcMap[so.Name]; ok {
				if so.Port != nil {
					services[idx].Port = *so.Port
				}
			}
		}
	}

	return modules, services
}

// NEIR v2.0 — Module Inheritance
// Modules can `extend` a base module, inheriting its path and dependencies.

type InheritedModule struct {
	Module
	Extend string `json:"extend,omitempty"`
}

// ResolveInheritance resolves module inheritance chains by matching name substrings.
func ResolveInheritance(modules []Module) []Module {
	index := make(map[string]int)
	for i, m := range modules {
		index[m.Name] = i
	}

	resolved := make([]Module, len(modules))
	copy(resolved, modules)

	visited := make(map[string]bool)
	resolve := func(name string) {
		visited[name] = true
		idx, ok := index[name]
		if !ok {
			return
		}
		m := resolved[idx]
		if m.Path == "" && len(m.Dependencies) == 0 {
			for _, other := range modules {
				if other.Name != m.Name && strings.Contains(m.Name, other.Name) {
					if other.Path != "" {
						resolved[idx].Path = other.Path
					}
					if len(other.Dependencies) > 0 {
						resolved[idx].Dependencies = other.Dependencies
					}
					break
				}
			}
		}
	}

	for _, m := range modules {
		if !visited[m.Name] {
			resolve(m.Name)
		}
	}

	return resolved
}

// ParseInheritedModules parses modules that may contain an `extend` field.
func ParseInheritedModules(rawModules []any) []InheritedModule {
	var result []InheritedModule
	for _, raw := range rawModules {
		mod, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		im := InheritedModule{
			Module: extractModule(mod),
		}
		if extend, ok := mod["extend"].(string); ok {
			im.Extend = extend
		}
		result = append(result, im)
	}
	return result
}

// ResolveInheritedModules resolves `extend` references between InheritedModules.
func ResolveInheritedModules(inherited []InheritedModule) []Module {
	baseMap := make(map[string]Module)
	for _, im := range inherited {
		baseMap[im.Name] = im.Module
	}

	result := make([]Module, 0, len(inherited))
	for _, im := range inherited {
		m := im.Module
		if im.Extend != "" {
			if base, ok := baseMap[im.Extend]; ok {
				if m.Path == "" {
					m.Path = base.Path
				}
				if len(m.Dependencies) == 0 {
					m.Dependencies = base.Dependencies
				}
				if m.Description == "" {
					m.Description = base.Description
				}
			}
		}
		result = append(result, m)
	}
	return result
}

// ResolveInheritanceChains resolves multi-level extend chains iteratively.
func ResolveInheritanceChains(inherited []InheritedModule) ([]Module, error) {
	baseMap := make(map[string]InheritedModule)
	for _, im := range inherited {
		baseMap[im.Name] = im
	}

	result := make(map[string]Module)
	depth := 0
	maxDepth := 10

	for depth < maxDepth {
		changed := false
		for _, im := range inherited {
			if _, resolved := result[im.Name]; resolved {
				continue
			}
			if im.Extend == "" {
				result[im.Name] = im.Module
				changed = true
				continue
			}
			if base, ok := result[im.Extend]; ok {
				m := im.Module
				if m.Path == "" {
					m.Path = base.Path
				}
				if len(m.Dependencies) == 0 {
					m.Dependencies = base.Dependencies
				}
				if m.Description == "" {
					m.Description = base.Description
				}
				result[im.Name] = m
				changed = true
			}
		}
		if !changed {
			break
		}
		depth++
	}

	if len(result) != len(inherited) {
		return nil, fmt.Errorf("could not resolve all inheritance chains (possible circular reference)")
	}

	output := make([]Module, 0, len(inherited))
	for _, im := range inherited {
		output = append(output, result[im.Name])
	}
	return output, nil
}
