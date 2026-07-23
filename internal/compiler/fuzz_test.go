package compiler

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
)

func FuzzBuildProjectContext(f *testing.F) {
	f.Add("test-proj", "1.0.0", "A test project")
	f.Add("", "", "")
	f.Add("my-app", "2.0.0", "Production app")

	f.Fuzz(func(t *testing.T, name, version, description string) {
		neir := &model.NEIR{
			Project: &project.Project{
				Name:        name,
				Version:     version,
				Description: description,
			},
		}
		ctx := buildProjectContext(neir)
		if neir.Project.Name != "" && ctx == "" {
			t.Error("expected non-empty context for non-empty project")
		}
	})
}

func FuzzCompileAll(f *testing.F) {
	f.Add("test-proj", "1.0.0")
	f.Add("my-app", "")

	f.Fuzz(func(t *testing.T, name, version string) {
		c := New()
		c.Register(&stubAdapter{target: TargetCopilot})
		c.Register(&stubAdapter{target: TargetClaude})
		c.Register(&errAdapter{target: TargetCursor})

		neir := &model.NEIR{
			Project: &project.Project{
				Name:    name,
				Version: version,
			},
		}
		results := c.CompileAll(neir)
		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
		for target, out := range results {
			if out == nil {
				t.Errorf("nil result for target %s", target)
			}
			if out != nil && out.Target != target {
				t.Errorf("target mismatch: expected %s, got %s", target, out.Target)
			}
		}
	})
}

func FuzzCompileSingle(f *testing.F) {
	f.Add("test-proj", "1.0.0")
	f.Add("my-app", "")

	f.Fuzz(func(t *testing.T, name, version string) {
		c := New()
		c.Register(&stubAdapter{target: TargetCopilot})

		neir := &model.NEIR{
			Project: &project.Project{
				Name:    name,
				Version: version,
			},
		}
		out, err := c.Compile(neir, TargetCopilot)
		if err != nil {
			t.Fatal(err)
		}
		if out == nil {
			t.Fatal("output should not be nil")
		}
		if out.Target != TargetCopilot {
			t.Errorf("target mismatch: expected %s, got %s", TargetCopilot, out.Target)
		}

		_, err = c.Compile(neir, TargetCursor)
		if err == nil {
			t.Error("expected error for unknown target")
		}
	})
}

func FuzzTargets(f *testing.F) {
	f.Add("copilot")
	f.Add("claude")
	f.Add("cursor")

	f.Fuzz(func(t *testing.T, target string) {
		c := New()
		c.Register(&stubAdapter{target: Target(target)})
		targets := c.Targets()
		if len(targets) != 1 {
			t.Errorf("expected 1 target, got %d", len(targets))
		}
	})
}
