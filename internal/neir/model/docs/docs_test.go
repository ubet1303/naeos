package docs

import "testing"

func TestDocKindConstants(t *testing.T) {
	tests := []struct {
		constant DocKind
		expected string
	}{
		{KindGuide, "guide"},
		{KindReference, "reference"},
		{KindADR, "adr"},
		{KindRFC, "rfc"},
		{KindChangelog, "changelog"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("DocKind %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var d Documentation
	if d.Guides != nil {
		t.Errorf("expected nil Guides, got %v", d.Guides)
	}
	if d.References != nil {
		t.Errorf("expected nil References, got %v", d.References)
	}
	if d.ADRs != nil {
		t.Errorf("expected nil ADRs, got %v", d.ADRs)
	}
	if d.RFCs != nil {
		t.Errorf("expected nil RFCs, got %v", d.RFCs)
	}

	var doc Doc
	if doc.Title != "" {
		t.Errorf("expected empty Title, got %q", doc.Title)
	}
	if doc.Kind != "" {
		t.Errorf("expected empty Kind, got %q", doc.Kind)
	}
}

func TestInitialization(t *testing.T) {
	d := Documentation{
		Guides: []Doc{
			{Title: "Getting Started", Path: "docs/getting-started.md", Kind: KindGuide, Summary: "Quick start guide"},
		},
		References: []Doc{
			{Title: "API Reference", Path: "docs/api.md", Kind: KindReference},
		},
		ADRs: []Doc{
			{Title: "Use PostgreSQL", Kind: KindADR, Summary: "Chose PostgreSQL for main storage"},
		},
		RFCs: []Doc{
			{Title: "Auth redesign", Kind: KindRFC},
		},
		Attributes: map[string]string{"version": "2.0"},
	}

	if len(d.Guides) != 1 || d.Guides[0].Title != "Getting Started" {
		t.Errorf("unexpected Guides: %v", d.Guides)
	}
	if d.References[0].Kind != KindReference {
		t.Errorf("expected KindReference, got %q", d.References[0].Kind)
	}
	if d.ADRs[0].Summary != "Chose PostgreSQL for main storage" {
		t.Errorf("unexpected ADR summary: %q", d.ADRs[0].Summary)
	}
}
