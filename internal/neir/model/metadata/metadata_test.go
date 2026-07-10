package metadata

import (
	"testing"
	"time"
)

func TestZeroValue(t *testing.T) {
	var m Metadata
	if m.NEIRVersion != "" {
		t.Errorf("expected empty NEIRVersion, got %q", m.NEIRVersion)
	}
	if m.SchemaVersion != "" {
		t.Errorf("expected empty SchemaVersion, got %q", m.SchemaVersion)
	}
	if m.ProjectVersion != "" {
		t.Errorf("expected empty ProjectVersion, got %q", m.ProjectVersion)
	}
	if m.CreatedAt != nil {
		t.Errorf("expected nil CreatedAt, got %v", m.CreatedAt)
	}
	if m.ModifiedAt != nil {
		t.Errorf("expected nil ModifiedAt, got %v", m.ModifiedAt)
	}
	if m.Tags != nil {
		t.Errorf("expected nil Tags, got %v", m.Tags)
	}
	if m.Source != nil {
		t.Errorf("expected nil Source, got %v", m.Source)
	}

	var sr SourceRef
	if sr.Kind != "" {
		t.Errorf("expected empty Kind, got %q", sr.Kind)
	}
	if sr.Ref != "" {
		t.Errorf("expected empty Ref, got %q", sr.Ref)
	}
}

func TestInitialization(t *testing.T) {
	now := time.Now()
	m := Metadata{
		NEIRVersion:    "1.0.0",
		SchemaVersion:  "2.0",
		ProjectVersion: "0.1.0",
		CreatedAt:      &now,
		ModifiedAt:     &now,
		Tags:           []string{"v1", "stable"},
		Labels:         map[string]string{"team": "platform"},
		Source:         &SourceRef{Kind: "git", Ref: "https://github.com/org/repo"},
	}

	if m.NEIRVersion != "1.0.0" {
		t.Errorf("expected NEIRVersion '1.0.0', got %q", m.NEIRVersion)
	}
	if m.CreatedAt == nil || !m.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt to equal now")
	}
	if len(m.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(m.Tags))
	}
	if m.Source == nil || m.Source.Kind != "git" {
		t.Errorf("expected Source.Kind 'git', got %v", m.Source)
	}
}

func TestMetadataPointerSemantics(t *testing.T) {
	var m Metadata
	m.CreatedAt = &time.Time{}
	if m.CreatedAt == nil {
		t.Error("expected non-nil CreatedAt after assignment")
	}
	m.Source = &SourceRef{Kind: "local", Ref: "/path/to/spec.json"}
	if m.Source.Ref != "/path/to/spec.json" {
		t.Errorf("expected Source.Ref '/path/to/spec.json', got %q", m.Source.Ref)
	}
}
