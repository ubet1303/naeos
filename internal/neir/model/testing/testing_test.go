package testing

import "testing"

func TestTestingStrategyConstants(t *testing.T) {
	tests := []struct {
		constant TestingStrategy
		expected string
	}{
		{StrategyUnit, "unit"},
		{StrategyIntegration, "integration"},
		{StrategyE2E, "e2e"},
		{StrategyContract, "contract"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("TestingStrategy %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var ts Testing
	if ts.Strategy != "" {
		t.Errorf("expected empty Strategy, got %q", ts.Strategy)
	}
	if ts.Frameworks != nil {
		t.Errorf("expected nil Frameworks, got %v", ts.Frameworks)
	}
	if ts.Coverage != nil {
		t.Errorf("expected nil Coverage, got %v", ts.Coverage)
	}
	if ts.Fixtures != nil {
		t.Errorf("expected nil Fixtures, got %v", ts.Fixtures)
	}

	var cov Coverage
	if cov.MinPercent != 0 {
		t.Errorf("expected zero MinPercent, got %f", cov.MinPercent)
	}

	var f Fixture
	if f.Name != "" {
		t.Errorf("expected empty Name, got %q", f.Name)
	}
}

func TestInitialization(t *testing.T) {
	ts := Testing{
		Strategy:   StrategyUnit,
		Frameworks: []string{"testing", "testify"},
		Coverage:   &Coverage{MinPercent: 80.0},
		Fixtures: []Fixture{
			{Name: "users.json", Kind: "json", Path: "testdata/users.json"},
			{Name: "seed.sql", Kind: "sql", Path: "testdata/seed.sql"},
		},
		Attributes: map[string]string{"ci": "github-actions"},
	}

	if ts.Strategy != StrategyUnit {
		t.Errorf("expected Strategy %q, got %q", StrategyUnit, ts.Strategy)
	}
	if ts.Coverage == nil || ts.Coverage.MinPercent != 80.0 {
		t.Errorf("expected Coverage.MinPercent 80.0, got %v", ts.Coverage)
	}
	if len(ts.Fixtures) != 2 {
		t.Errorf("expected 2 fixtures, got %d", len(ts.Fixtures))
	}
	if ts.Fixtures[0].Path != "testdata/users.json" {
		t.Errorf("expected fixture path 'testdata/users.json', got %q", ts.Fixtures[0].Path)
	}
}

func TestCoverageNilPointer(t *testing.T) {
	var ts Testing
	if ts.Coverage != nil {
		t.Error("expected nil Coverage")
	}
	ts.Coverage = &Coverage{MinPercent: 95.5}
	if ts.Coverage == nil || ts.Coverage.MinPercent != 95.5 {
		t.Errorf("expected Coverage.MinPercent 95.5, got %v", ts.Coverage)
	}
}
