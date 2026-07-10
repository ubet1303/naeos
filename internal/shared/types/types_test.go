package types

import (
	"testing"
)

func TestID(t *testing.T) {
	var id ID = "TEST-001"
	if string(id) != "TEST-001" {
		t.Errorf("ID = %q, want %q", id, "TEST-001")
	}
}

func TestReference(t *testing.T) {
	ref := Reference{ID: "REF-001", Name: "my-ref"}
	if string(ref.ID) != "REF-001" {
		t.Errorf("Reference.ID = %q, want %q", ref.ID, "REF-001")
	}
	if ref.Name != "my-ref" {
		t.Errorf("Reference.Name = %q, want %q", ref.Name, "my-ref")
	}
}

func TestErrorInfo(t *testing.T) {
	e := ErrorInfo{Code: "E001", Message: "something went wrong"}
	if e.Code != "E001" {
		t.Errorf("ErrorInfo.Code = %q, want %q", e.Code, "E001")
	}
	if e.Message != "something went wrong" {
		t.Errorf("ErrorInfo.Message = %q, want %q", e.Message, "something went wrong")
	}
}

func TestArtifact(t *testing.T) {
	a := Artifact{Path: "output/main.go", Content: []byte("package main")}
	if a.Path != "output/main.go" {
		t.Errorf("Artifact.Path = %q, want %q", a.Path, "output/main.go")
	}
	if string(a.Content) != "package main" {
		t.Errorf("Artifact.Content = %q, want %q", a.Content, "package main")
	}
}

func TestTask(t *testing.T) {
	task := Task{
		ID:           "T-001",
		Name:         "build",
		Dependencies: []string{"T-000"},
		Priority:     1,
	}
	if task.ID != "T-001" {
		t.Errorf("Task.ID = %q, want %q", task.ID, "T-001")
	}
	if len(task.Dependencies) != 1 {
		t.Errorf("Task.Dependencies has %d entries, want 1", len(task.Dependencies))
	}
	if task.Priority != 1 {
		t.Errorf("Task.Priority = %d, want 1", task.Priority)
	}
}

func TestPolicyRule(t *testing.T) {
	rule := PolicyRule{
		RuleID:    "PR-001",
		Condition: "exists:modules",
		Priority:  1,
		Action:    "block",
		Scope:     "project",
	}
	if rule.RuleID != "PR-001" {
		t.Errorf("PolicyRule.RuleID = %q, want %q", rule.RuleID, "PR-001")
	}
	if rule.Action != "block" {
		t.Errorf("PolicyRule.Action = %q, want %q", rule.Action, "block")
	}
}

func TestKnowledgeEntry(t *testing.T) {
	entry := KnowledgeEntry{
		Topic:     "architecture",
		Component: "kernel",
		Version:   "1.0.0",
		Rationale: "Decision record",
	}
	if entry.Topic != "architecture" {
		t.Errorf("KnowledgeEntry.Topic = %q, want %q", entry.Topic, "architecture")
	}
}

func TestTelemetryEvent(t *testing.T) {
	event := TelemetryEvent{
		Name:      "pipeline.run",
		Timestamp: 1234567890,
		Payload:   map[string]any{"artifacts": 10},
	}
	if event.Name != "pipeline.run" {
		t.Errorf("TelemetryEvent.Name = %q, want %q", event.Name, "pipeline.run")
	}
	if event.Payload["artifacts"] != 10 {
		t.Errorf("TelemetryEvent.Payload[artifacts] = %v, want 10", event.Payload["artifacts"])
	}
}

func TestValidationResult(t *testing.T) {
	result := ValidationResult{
		Valid:  true,
		Errors: []ErrorInfo{},
	}
	if !result.Valid {
		t.Error("ValidationResult.Valid should be true")
	}
	if len(result.Errors) != 0 {
		t.Errorf("ValidationResult.Errors has %d entries, want 0", len(result.Errors))
	}
}

func TestReviewResult(t *testing.T) {
	result := ReviewResult{
		Approved: true,
		Comments: []string{"looks good"},
	}
	if !result.Approved {
		t.Error("ReviewResult.Approved should be true")
	}
	if len(result.Comments) != 1 {
		t.Errorf("ReviewResult.Comments has %d entries, want 1", len(result.Comments))
	}
}

func TestSpecDocument(t *testing.T) {
	doc := SpecDocument{
		Raw:     "project: test",
		Project: "test",
		Modules: []ModuleDef{{Name: "core", Path: "./internal/core"}},
		Services: []ServiceDef{{Name: "api", Kind: "http", Port: 8080}},
	}
	if doc.Project != "test" {
		t.Errorf("SpecDocument.Project = %q, want %q", doc.Project, "test")
	}
	if len(doc.Modules) != 1 {
		t.Errorf("SpecDocument.Modules has %d entries, want 1", len(doc.Modules))
	}
	if len(doc.Services) != 1 {
		t.Errorf("SpecDocument.Services has %d entries, want 1", len(doc.Services))
	}
}
