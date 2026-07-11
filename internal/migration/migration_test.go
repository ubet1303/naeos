package migration

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	v, err := ParseVersion("1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("expected 1.2.3, got %v", v)
	}
}

func TestParseVersionInvalid(t *testing.T) {
	_, err := ParseVersion("invalid")
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
}

func TestVersionLess(t *testing.T) {
	v1 := Version{Major: 1, Minor: 0, Patch: 0}
	v2 := Version{Major: 2, Minor: 0, Patch: 0}
	if !v1.Less(v2) {
		t.Error("expected 1.0.0 < 2.0.0")
	}
	if v2.Less(v1) {
		t.Error("expected 2.0.0 not < 1.0.0")
	}
}

func TestVersionString(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3}
	if v.String() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", v.String())
	}
}

func TestPlannerPlan(t *testing.T) {
	planner := NewPlanner()
	plan, err := planner.Plan("0.1.0", "0.3.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan) < 1 {
		t.Error("expected at least 1 migration step")
	}
}

func TestPlannerPlanNoMigrationNeeded(t *testing.T) {
	planner := NewPlanner()
	_, err := planner.Plan("0.3.0", "0.1.0")
	if err == nil {
		t.Fatal("expected error for downgrade")
	}
}

func TestPlannerMigrate(t *testing.T) {
	planner := NewPlanner()
	spec := []byte("project: test\n")
	result, err := planner.Migrate(spec, "0.1.0", "0.3.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}
