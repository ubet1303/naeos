package compliance

import (
	"strings"
	"testing"
)

func TestSOC2Report(t *testing.T) {
	t.Parallel()
	checks := map[string]ControlStatus{
		"CC6.1": StatusPass,
		"CC6.2": StatusPass,
		"CC6.3": StatusFail,
		"CC7.1": StatusPartial,
	}

	report := GenerateSOC2Report(checks)

	if report.Framework != FrameworkSOC2 {
		t.Errorf("expected SOC2, got %s", report.Framework)
	}
	if len(report.Controls) == 0 {
		t.Fatal("expected controls")
	}
	if report.Summary.Pass != 2 {
		t.Errorf("expected 2 pass, got %d", report.Summary.Pass)
	}
	if report.Summary.Fail != 1 {
		t.Errorf("expected 1 fail, got %d", report.Summary.Fail)
	}
}

func TestHIPAAReport(t *testing.T) {
	t.Parallel()
	checks := map[string]ControlStatus{
		"§164.308(a)(1)": StatusPass,
		"§164.312(a)(1)": StatusPass,
	}

	report := GenerateHIPAAReport(checks)

	if report.Framework != FrameworkHIPAA {
		t.Errorf("expected HIPAA, got %s", report.Framework)
	}
	if report.Summary.Pass != 2 {
		t.Errorf("expected 2 pass, got %d", report.Summary.Pass)
	}
}

func TestFormatReport(t *testing.T) {
	t.Parallel()
	report := GenerateSOC2Report(map[string]ControlStatus{
		"CC6.1": StatusPass,
		"CC6.2": StatusFail,
	})

	output := FormatReport(report)
	if !strings.Contains(output, "SOC2 Compliance Report") {
		t.Error("expected report header")
	}
	if !strings.Contains(output, "Score:") {
		t.Error("expected score")
	}
}

func TestScore(t *testing.T) {
	t.Parallel()
	checks := map[string]ControlStatus{
		"CC6.1": StatusPass,
		"CC6.2": StatusPass,
		"CC6.3": StatusPass,
		"CC7.1": StatusPass,
		"CC7.2": StatusPass,
		"A1.2":  StatusPass,
		"C1.1":  StatusPass,
		"CC8.1": StatusPass,
	}

	report := GenerateSOC2Report(checks)
	if report.Summary.Score != 100.0 {
		t.Errorf("expected score 100, got %.1f", report.Summary.Score)
	}
}
