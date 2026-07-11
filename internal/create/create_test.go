package create

import (
	"bufio"
	"strings"
	"testing"
)

func TestProjectConfigToSpec(t *testing.T) {
	cfg := &ProjectConfig{
		Name:         "my-project",
		ModulePath:   "./internal/core",
		Language:     "go",
		Architecture: "hexagonal",
		Deployment:   "rolling",
		Port:         8080,
		EnableTesting: true,
	}

	spec := cfg.ToSpec()
	if spec == "" {
		t.Error("expected non-empty spec")
	}
	if !contains(spec, "my-project") {
		t.Error("expected project name in spec")
	}
	if !contains(spec, "hexagonal") {
		t.Error("expected architecture in spec")
	}
	if !contains(spec, "8080") {
		t.Error("expected port in spec")
	}
	if !contains(spec, "rolling") {
		t.Error("expected deployment in spec")
	}
	if !contains(spec, "testing:") {
		t.Error("expected testing section when enabled")
	}
}

func TestProjectConfigToSpecNoTesting(t *testing.T) {
	cfg := &ProjectConfig{
		Name:          "test",
		ModulePath:    "./test",
		Language:      "go",
		Architecture:  "layered",
		Deployment:    "rolling",
		Port:          3000,
		EnableTesting: false,
	}

	spec := cfg.ToSpec()
	if contains(spec, "testing:") {
		t.Error("expected no testing section when disabled")
	}
}

func TestProjectConfigToSpecNoDescription(t *testing.T) {
	cfg := &ProjectConfig{
		Name:       "minimal",
		ModulePath: "./min",
		Port:       80,
	}
	spec := cfg.ToSpec()
	if contains(spec, "description:") {
		t.Error("expected no description when empty")
	}
	if !contains(spec, "minimal") {
		t.Error("expected project name")
	}
}

func newTestWizard(input string) *Wizard {
	return &Wizard{
		reader: bufio.NewReader(strings.NewReader(input)),
	}
}

func TestWizardAskDefault(t *testing.T) {
	w := newTestWizard("\n")
	val := w.askDefault("test", "fallback")
	if val != "fallback" {
		t.Errorf("expected fallback, got %q", val)
	}
}

func TestWizardAskDefaultCustom(t *testing.T) {
	w := newTestWizard("custom\n")
	val := w.askDefault("test", "fallback")
	if val != "custom" {
		t.Errorf("expected custom, got %q", val)
	}
}

func TestWizardAskRequired(t *testing.T) {
	w := newTestWizard("hello\n")
	val := w.askRequired("test")
	if val != "hello" {
		t.Errorf("expected hello, got %q", val)
	}
}

func TestWizardAskRequiredEmptyThenValue(t *testing.T) {
	w := newTestWizard("\nreal\n")
	val := w.askRequired("test")
	if val != "real" {
		t.Errorf("expected real, got %q", val)
	}
}

func TestWizardAskChoiceDefault(t *testing.T) {
	w := newTestWizard("\n")
	val := w.askChoice("pick", []string{"a", "b", "c"}, "b")
	if val != "b" {
		t.Errorf("expected b (default), got %q", val)
	}
}

func TestWizardAskChoiceNumber(t *testing.T) {
	w := newTestWizard("3\n")
	val := w.askChoice("pick", []string{"a", "b", "c"}, "a")
	if val != "c" {
		t.Errorf("expected c, got %q", val)
	}
}

func TestWizardAskChoiceInvalidNumber(t *testing.T) {
	w := newTestWizard("99\n")
	val := w.askChoice("pick", []string{"a", "b"}, "a")
	if val != "a" {
		t.Errorf("expected default a for invalid input, got %q", val)
	}
}

func TestWizardAskChoiceInvalidText(t *testing.T) {
	w := newTestWizard("xyz\n")
	val := w.askChoice("pick", []string{"a", "b"}, "a")
	if val != "a" {
		t.Errorf("expected default a for text input, got %q", val)
	}
}

func TestWizardAskIntDefault(t *testing.T) {
	w := newTestWizard("\n")
	val := w.askInt("port", 8080)
	if val != 8080 {
		t.Errorf("expected 8080, got %d", val)
	}
}

func TestWizardAskIntCustom(t *testing.T) {
	w := newTestWizard("3000\n")
	val := w.askInt("port", 8080)
	if val != 3000 {
		t.Errorf("expected 3000, got %d", val)
	}
}

func TestWizardAskIntInvalid(t *testing.T) {
	w := newTestWizard("abc\n")
	val := w.askInt("port", 8080)
	if val != 8080 {
		t.Errorf("expected default 8080, got %d", val)
	}
}

func TestWizardAskYesNoDefaultYes(t *testing.T) {
	w := newTestWizard("\n")
	val := w.askYesNo("ok?", true)
	if !val {
		t.Error("expected true (default)")
	}
}

func TestWizardAskYesNoDefaultNo(t *testing.T) {
	w := newTestWizard("\n")
	val := w.askYesNo("ok?", false)
	if val {
		t.Error("expected false (default)")
	}
}

func TestWizardAskYesNoYes(t *testing.T) {
	w := newTestWizard("y\n")
	val := w.askYesNo("ok?", false)
	if !val {
		t.Error("expected true for 'y'")
	}
}

func TestWizardAskYesNoYesFull(t *testing.T) {
	w := newTestWizard("yes\n")
	val := w.askYesNo("ok?", false)
	if !val {
		t.Error("expected true for 'yes'")
	}
}

func TestWizardAskYesNoNo(t *testing.T) {
	w := newTestWizard("n\n")
	val := w.askYesNo("ok?", true)
	if val {
		t.Error("expected false for 'n'")
	}
}

func TestWizardAskYesNoInvalid(t *testing.T) {
	w := newTestWizard("maybe\n")
	val := w.askYesNo("ok?", true)
	if val {
		t.Error("expected false for non-yes input")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
