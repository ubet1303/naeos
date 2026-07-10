package deployment

import "testing"

func TestStrategyConstants(t *testing.T) {
	tests := []struct {
		constant Strategy
		expected string
	}{
		{StrategyRolling, "rolling"},
		{StrategyBlueGreen, "blue-green"},
		{StrategyCanary, "canary"},
		{StrategyRecreate, "recreate"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("Strategy %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var d Deployment
	if d.Target != "" {
		t.Errorf("expected empty Target, got %q", d.Target)
	}
	if d.Strategy != "" {
		t.Errorf("expected empty Strategy, got %q", d.Strategy)
	}
	if d.Environments != nil {
		t.Errorf("expected nil Environments, got %v", d.Environments)
	}
	if d.Scaling != nil {
		t.Errorf("expected nil Scaling, got %v", d.Scaling)
	}

	var env Environment
	if env.Name != "" {
		t.Errorf("expected empty Name, got %q", env.Name)
	}
	if env.Config != nil {
		t.Errorf("expected nil Config, got %v", env.Config)
	}

	var s Scaling
	if s.Min != 0 {
		t.Errorf("expected zero Min, got %d", s.Min)
	}
	if s.Max != 0 {
		t.Errorf("expected zero Max, got %d", s.Max)
	}
}

func TestInitialization(t *testing.T) {
	d := Deployment{
		Target:   "k8s",
		Strategy: StrategyCanary,
		Environments: []Environment{
			{Name: "staging", Kind: "pre-prod", Config: map[string]string{"replicas": "2"}},
			{Name: "production", Kind: "prod", Config: map[string]string{"replicas": "5"}},
		},
		Scaling: &Scaling{Min: 2, Max: 10, Replicas: 3},
	}

	if d.Strategy != StrategyCanary {
		t.Errorf("expected Strategy %q, got %q", StrategyCanary, d.Strategy)
	}
	if len(d.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(d.Environments))
	}
	if d.Scaling == nil || d.Scaling.Max != 10 {
		t.Errorf("expected Scaling.Max 10, got %v", d.Scaling)
	}
}

func TestScalingNilPointer(t *testing.T) {
	var d Deployment
	if d.Scaling != nil {
		t.Error("expected nil Scaling")
	}
	d.Scaling = &Scaling{Min: 1, Max: 5}
	if d.Scaling.Min != 1 || d.Scaling.Max != 5 {
		t.Errorf("unexpected Scaling values: %+v", d.Scaling)
	}
}
