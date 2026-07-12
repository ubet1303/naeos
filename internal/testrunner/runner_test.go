package testrunner

import (
	"strings"
	"testing"
)

func TestNewRunnerDefaults(t *testing.T) {
	r := NewRunner(TestConfig{})
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
	if r.config.WorkingDir != "." {
		t.Errorf("expected default working dir '.', got %s", r.config.WorkingDir)
	}
	if r.config.Timeout != 300 {
		t.Errorf("expected default timeout 300, got %d", r.config.Timeout)
	}
}

func TestNewRunnerCustom(t *testing.T) {
	r := NewRunner(TestConfig{
		WorkingDir: "/tmp/test",
		Timeout:    60,
		Languages:  []string{"go"},
		Verbose:    true,
	})
	if r.config.WorkingDir != "/tmp/test" {
		t.Errorf("expected '/tmp/test', got %s", r.config.WorkingDir)
	}
	if r.config.Timeout != 60 {
		t.Errorf("expected 60, got %d", r.config.Timeout)
	}
}

func TestRunLanguageUnsupported(t *testing.T) {
	r := NewRunner(TestConfig{WorkingDir: "."})
	_, err := r.RunLanguage("cobol")
	if err == nil {
		t.Error("expected error for unsupported language")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected 'unsupported' in error, got: %s", err.Error())
	}
}

func TestRunAllEmpty(t *testing.T) {
	r := NewRunner(TestConfig{
		WorkingDir: "/nonexistent-dir",
		Languages:  []string{},
	})
	results, err := r.RunAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty languages, got %d", len(results))
	}
}

func TestFormatResultsAllPass(t *testing.T) {
	results := []TestResult{
		{Language: "go", Passed: true, Tests: 5},
		{Language: "python", Passed: true, Tests: 3},
	}
	output := FormatResults(results)

	if !strings.Contains(output, "All tests passed!") {
		t.Error("expected 'All tests passed!'")
	}
	if !strings.Contains(output, "[PASS] go") {
		t.Error("expected PASS for go")
	}
	if !strings.Contains(output, "[PASS] python") {
		t.Error("expected PASS for python")
	}
}

func TestFormatResultsSomeFail(t *testing.T) {
	results := []TestResult{
		{Language: "go", Passed: true, Tests: 5},
		{Language: "rust", Passed: false, Failures: 2},
	}
	output := FormatResults(results)

	if !strings.Contains(output, "Some tests failed.") {
		t.Error("expected 'Some tests failed.'")
	}
	if !strings.Contains(output, "[FAIL] rust") {
		t.Error("expected FAIL for rust")
	}
	if !strings.Contains(output, "(2 failures)") {
		t.Error("expected failure count")
	}
}

func TestParseGoOutput(t *testing.T) {
	r := NewRunner(TestConfig{})
	result := &TestResult{
		Output: "ok  \tpkg1\t0.01s\nok  \tpkg2\t0.02s\nFAIL\tpkg3\t0.03s",
	}
	r.parseGoOutput(result)

	if result.Tests != 2 {
		t.Errorf("expected 2 tests, got %d", result.Tests)
	}
	if result.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", result.Failures)
	}
}

func TestParseGoOutputEmpty(t *testing.T) {
	r := NewRunner(TestConfig{})
	result := &TestResult{Output: ""}
	r.parseGoOutput(result)

	if result.Tests != 0 {
		t.Errorf("expected 0 tests, got %d", result.Tests)
	}
	if result.Failures != 0 {
		t.Errorf("expected 0 failures, got %d", result.Failures)
	}
}

func TestDetectLanguagesGo(t *testing.T) {
	r := NewRunner(TestConfig{WorkingDir: "../.."})
	langs := r.detectLanguages()
	found := false
	for _, l := range langs {
		if l == "go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'go' to be detected (go.mod exists in project root)")
	}
}
