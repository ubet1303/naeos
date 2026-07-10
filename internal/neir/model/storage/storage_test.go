package storage

import "testing"

func TestStorageTypeConstants(t *testing.T) {
	tests := []struct {
		constant StorageType
		expected string
	}{
		{TypeSQL, "sql"},
		{TypeNoSQL, "nosql"},
		{TypeFile, "file"},
		{TypeCache, "cache"},
		{TypeQueue, "queue"},
		{TypeBlob, "blob"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("StorageType %v = %q, want %q", tt.constant, string(tt.constant), tt.expected)
		}
	}
}

func TestZeroValue(t *testing.T) {
	var s Storage
	if s.Name != "" {
		t.Errorf("expected empty Name, got %q", s.Name)
	}
	if s.Type != "" {
		t.Errorf("expected empty Type, got %q", s.Type)
	}
	if s.Collections != nil {
		t.Errorf("expected nil Collections, got %v", s.Collections)
	}

	var c Collection
	if c.Name != "" {
		t.Errorf("expected empty Name, got %q", c.Name)
	}
	if c.Schema != nil {
		t.Errorf("expected nil Schema, got %v", c.Schema)
	}
}

func TestInitialization(t *testing.T) {
	s := Storage{
		Name:       "main-db",
		Type:       TypeSQL,
		Provider:   "postgresql",
		Connection: "postgres://localhost:5432/mydb",
		Collections: []Collection{
			{Name: "users", Schema: map[string]string{"id": "uuid", "email": "varchar"}},
			{Name: "orders", Schema: map[string]string{"id": "uuid", "total": "numeric"}},
		},
	}

	if s.Type != TypeSQL {
		t.Errorf("expected Type %q, got %q", TypeSQL, s.Type)
	}
	if len(s.Collections) != 2 {
		t.Errorf("expected 2 collections, got %d", len(s.Collections))
	}
	if s.Collections[0].Schema["email"] != "varchar" {
		t.Errorf("expected email=varchar, got %q", s.Collections[0].Schema["email"])
	}
}

func TestAllStorageTypes(t *testing.T) {
	types := []StorageType{TypeSQL, TypeNoSQL, TypeFile, TypeCache, TypeQueue, TypeBlob}
	seen := map[StorageType]bool{}
	for _, st := range types {
		if seen[st] {
			t.Errorf("duplicate StorageType: %q", st)
		}
		seen[st] = true
	}
}
