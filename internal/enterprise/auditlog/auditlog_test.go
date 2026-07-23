package auditlog

import (
	"testing"
	"time"
)

func TestMemoryWriter(t *testing.T) {
	t.Parallel()
	w := NewMemoryWriter()
	defer w.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		UserID:    "user-1",
		Action:    ActionCreate,
		Resource:  "spec",
		Success:   true,
	}

	if err := w.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}

	entries := w.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].UserID != "user-1" {
		t.Errorf("expected user-1, got %s", entries[0].UserID)
	}
}

func TestLoggerAutoTimestamp(t *testing.T) {
	t.Parallel()
	w := NewMemoryWriter()
	logger := NewLogger(w)
	defer logger.Close()

	logger.Log(Entry{UserID: "u1", Action: ActionRead})

	entries := w.Entries()
	if len(entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("timestamp should be set automatically")
	}
	if entries[0].Severity != SeverityInfo {
		t.Errorf("expected severity info, got %s", entries[0].Severity)
	}
}

func TestFilterByAction(t *testing.T) {
	t.Parallel()
	entries := []Entry{
		{Action: ActionCreate, UserID: "u1"},
		{Action: ActionRead, UserID: "u2"},
		{Action: ActionCreate, UserID: "u3"},
	}
	result := FilterByAction(entries, ActionCreate)
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestFilterByUser(t *testing.T) {
	t.Parallel()
	entries := []Entry{
		{UserID: "u1", Action: ActionCreate},
		{UserID: "u2", Action: ActionRead},
		{UserID: "u1", Action: ActionDelete},
	}
	result := FilterByUser(entries, "u1")
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestFilterBySeverity(t *testing.T) {
	t.Parallel()
	entries := []Entry{
		{Severity: SeverityInfo},
		{Severity: SeverityWarning},
		{Severity: SeverityError},
		{Severity: SeverityCritical},
	}
	result := FilterBySeverity(entries, SeverityError)
	if len(result) != 2 {
		t.Errorf("expected 2 (error+critical), got %d", len(result))
	}
}

func TestSplunkWriter(t *testing.T) {
	t.Parallel()
	w := NewSplunkWriter()
	defer w.Close()

	w.Write(Entry{
		Timestamp: time.Now().UTC(),
		UserID:    "u1",
		Action:    ActionCreate,
		Resource:  "spec",
		Success:   true,
	})

	events := w.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	event := events[0]
	if event["source"] != "naeos-audit" {
		t.Errorf("expected source=naeos-audit, got %v", event["source"])
	}
	inner, ok := event["event"].(map[string]any)
	if !ok {
		t.Fatal("expected nested event map")
	}
	if inner["user_id"] != "u1" {
		t.Errorf("expected user_id=u1, got %v", inner["user_id"])
	}
}
