package parser

import (
	"testing"
)

func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		name     string
		cond     Condition
		envVars  map[string]string
		expected bool
	}{
		{
			name:     "no condition means always true",
			cond:     Condition{},
			envVars:  map[string]string{},
			expected: true,
		},
		{
			name:     "env set, equals matches",
			cond:     Condition{Env: "APP_ENV", Equals: "production"},
			envVars:  map[string]string{"APP_ENV": "production"},
			expected: true,
		},
		{
			name:     "env set, equals mismatches",
			cond:     Condition{Env: "APP_ENV", Equals: "production"},
			envVars:  map[string]string{"APP_ENV": "staging"},
			expected: false,
		},
		{
			name:     "env not set",
			cond:     Condition{Env: "APP_ENV", Equals: "production"},
			envVars:  map[string]string{},
			expected: false,
		},
		{
			name:     "env exists check",
			cond:     Condition{Env: "FEATURE_X"},
			envVars:  map[string]string{"FEATURE_X": "anything"},
			expected: true,
		},
		{
			name:     "env not exists",
			cond:     Condition{Env: "FEATURE_X"},
			envVars:  map[string]string{},
			expected: false,
		},
		{
			name:     "not equals match",
			cond:     Condition{Env: "APP_ENV", Equals: "production", Not: true},
			envVars:  map[string]string{"APP_ENV": "staging"},
			expected: true,
		},
		{
			name:     "not equals mismatch",
			cond:     Condition{Env: "APP_ENV", Equals: "production", Not: true},
			envVars:  map[string]string{"APP_ENV": "production"},
			expected: false,
		},
		{
			name:     "not exists check",
			cond:     Condition{Env: "FEATURE_X", Not: true},
			envVars:  map[string]string{},
			expected: true,
		},
		{
			name:     "not exists but exists",
			cond:     Condition{Env: "FEATURE_X", Not: true},
			envVars:  map[string]string{"FEATURE_X": "val"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := EvaluateCondition(tt.cond, tt.envVars)
			if got != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilterConditionalModules(t *testing.T) {
	t.Parallel()
	modules := []ConditionalModule{
		{Module: Module{Name: "core", Path: "./core"}},
		{Module: Module{Name: "feature-x", Path: "./feature-x"}, Condition: Condition{Env: "FEATURE_X"}},
		{Module: Module{Name: "prod-only", Path: "./prod"}, Condition: Condition{Env: "APP_ENV", Equals: "production"}},
	}

	t.Run("all enabled when no env vars", func(t *testing.T) {
		t.Parallel()
		result := FilterConditionalModules(modules, map[string]string{})
		if len(result) != 1 {
			t.Errorf("expected 1 module, got %d", len(result))
		}
	})

	t.Run("feature-x enabled", func(t *testing.T) {
		t.Parallel()
		result := FilterConditionalModules(modules, map[string]string{"FEATURE_X": "1"})
		if len(result) != 2 {
			t.Errorf("expected 2 modules, got %d", len(result))
		}
	})

	t.Run("all enabled", func(t *testing.T) {
		t.Parallel()
		result := FilterConditionalModules(modules, map[string]string{
			"FEATURE_X": "1",
			"APP_ENV":   "production",
		})
		if len(result) != 3 {
			t.Errorf("expected 3 modules, got %d", len(result))
		}
	})
}

func TestApplyProfiles(t *testing.T) {
	t.Parallel()
	port8080 := 8080
	port9090 := 9090
	enabled := true

	modules := []Module{
		{Name: "core", Path: "./core"},
		{Name: "web", Path: "./web"},
	}
	services := []Service{
		{Name: "api", Port: 8080},
		{Name: "admin", Port: 9090},
	}

	profiles := []EnvironmentProfile{
		{
			Name: "production",
			Modules: []ModuleOverride{
				{Name: "core", Path: "./core-prod"},
				{Name: "web", Enabled: &enabled},
			},
			Services: []ServiceOverride{
				{Name: "api", Port: &port8080},
			},
		},
		{
			Name: "staging",
			Services: []ServiceOverride{
				{Name: "api", Port: &port9090},
			},
		},
	}

	t.Run("no active profile", func(t *testing.T) {
		t.Parallel()
		m, s := ApplyProfiles(modules, services, profiles, "")
		if len(m) != 2 || len(s) != 2 {
			t.Error("should not modify when no profile active")
		}
	})

	t.Run("production profile", func(t *testing.T) {
		t.Parallel()
		m, s := ApplyProfiles(
			[]Module{{Name: "core", Path: "./core"}, {Name: "web", Path: "./web"}},
			[]Service{{Name: "api", Port: 8080}},
			profiles, "production",
		)
		if m[0].Path != "./core-prod" {
			t.Errorf("expected ./core-prod, got %s", m[0].Path)
		}
		if s[0].Port != 8080 {
			t.Errorf("expected port 8080, got %d", s[0].Port)
		}
	})

	t.Run("staging profile", func(t *testing.T) {
		t.Parallel()
		_, s := ApplyProfiles(
			[]Module{},
			[]Service{{Name: "api", Port: 8080}},
			profiles, "staging",
		)
		if s[0].Port != 9090 {
			t.Errorf("expected port 9090, got %d", s[0].Port)
		}
	})
}

func TestResolveInheritedModules(t *testing.T) {
	t.Parallel()
	inherited := []InheritedModule{
		{Module: Module{Name: "base", Path: "./base", Dependencies: []string{"utils"}}},
		{Module: Module{Name: "app"}, Extend: "base"},
		{Module: Module{Name: "service", Path: "./service"}},
	}

	result := ResolveInheritedModules(inherited)

	baseResult := result[0]
	if baseResult.Path != "./base" {
		t.Errorf("base should keep path, got %s", baseResult.Path)
	}

	appResult := result[1]
	if appResult.Path != "./base" {
		t.Errorf("app should inherit path from base, got %s", appResult.Path)
	}
	if len(appResult.Dependencies) != 1 || appResult.Dependencies[0] != "utils" {
		t.Errorf("app should inherit dependencies from base, got %v", appResult.Dependencies)
	}

	serviceResult := result[2]
	if serviceResult.Path != "./service" {
		t.Errorf("service should keep its own path, got %s", serviceResult.Path)
	}
}

func TestParseInheritedModules(t *testing.T) {
	t.Parallel()
	raw := []any{
		map[string]any{"name": "base", "path": "./base"},
		map[string]any{"name": "app", "extend": "base"},
	}
	result := ParseInheritedModules(raw)
	if len(result) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(result))
	}
	if result[1].Extend != "base" {
		t.Errorf("expected extend=base, got %s", result[1].Extend)
	}
}
