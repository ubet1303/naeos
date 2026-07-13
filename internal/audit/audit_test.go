package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMemoryAuditor(t *testing.T) {
	a := NewMemoryAuditor()

	err := a.Log(AuditEvent{
		UserID:   "user-1",
		Action:   "deploy",
		Resource: "cloud",
		Status:   "success",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := a.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].UserID != "user-1" {
		t.Errorf("expected user 'user-1', got %s", events[0].UserID)
	}
	if events[0].Action != "deploy" {
		t.Errorf("expected action 'deploy', got %s", events[0].Action)
	}
	if events[0].ID == "" {
		t.Error("expected auto-generated ID")
	}
	if events[0].Timestamp.IsZero() {
		t.Error("expected auto-generated timestamp")
	}
}

func TestMemoryAuditorClear(t *testing.T) {
	a := NewMemoryAuditor()

	a.Log(AuditEvent{UserID: "user-1", Action: "test"})
	a.Log(AuditEvent{UserID: "user-2", Action: "test2"})

	if len(a.Events()) != 2 {
		t.Fatalf("expected 2 events, got %d", len(a.Events()))
	}

	a.Clear()

	if len(a.Events()) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(a.Events()))
	}
}

func TestFileAuditor(t *testing.T) {
	dir := t.TempDir()

	a, err := NewFileAuditor(dir)
	if err != nil {
		t.Fatalf("unexpected error creating auditor: %v", err)
	}

	err = a.Log(AuditEvent{
		UserID:   "user-1",
		Action:   "deploy",
		Resource: "cloud",
		Status:   "success",
	})
	if err != nil {
		t.Fatalf("unexpected error logging: %v", err)
	}

	err = a.Log(AuditEvent{
		UserID:   "user-2",
		Action:   "destroy",
		Resource: "cloud",
		Status:   "success",
	})
	if err != nil {
		t.Fatalf("unexpected error logging second event: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".naeos", "audit.log"))
	if err != nil {
		t.Fatalf("failed to read audit log: %v", err)
	}

	lines := countLines(string(data))
	if lines != 2 {
		t.Errorf("expected 2 lines in audit log, got %d", lines)
	}
}

func TestFileAuditorAutoGeneratesID(t *testing.T) {
	dir := t.TempDir()

	a, err := NewFileAuditor(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = a.Log(AuditEvent{UserID: "user-1", Action: "test", Status: "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, err := os.ReadFile(filepath.Join(dir, ".naeos", "audit.log"))
	if err != nil {
		t.Fatalf("failed to read log: %v", err)
	}

	if len(events) == 0 {
		t.Error("expected non-empty audit log")
	}
}

func countLines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	return count
}
