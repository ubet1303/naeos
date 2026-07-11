package marketplace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(t.TempDir())
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSearchDefaultEntries(t *testing.T) {
	client := NewClient(t.TempDir())
	results, err := client.Search(SearchFilter{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) < 3 {
		t.Errorf("expected at least 3 default entries, got %d", len(results))
	}
}

func TestSearchWithQuery(t *testing.T) {
	client := NewClient(t.TempDir())
	results, err := client.Search(SearchFilter{Query: "go", Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected results for 'go' query")
	}
}

func TestSearchWithTags(t *testing.T) {
	client := NewClient(t.TempDir())
	results, err := client.Search(SearchFilter{Tags: []string{"python"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected results for 'python' tag")
	}
}

func TestGet(t *testing.T) {
	client := NewClient(t.TempDir())
	entry, err := client.Get("go-http-api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Name != "go-http-api" {
		t.Errorf("expected go-http-api, got %s", entry.Name)
	}
}

func TestGetNotFound(t *testing.T) {
	client := NewClient(t.TempDir())
	_, err := client.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent entry")
	}
}

func TestPublish(t *testing.T) {
	cacheDir := t.TempDir()
	client := NewClient(cacheDir)
	entry := RegistryEntry{Name: "test-spec", Version: "1.0.0", Description: "Test"}
	if err := client.Publish(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it was published
	data, err := os.ReadFile(filepath.Join(cacheDir, "registry.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty registry file")
	}
}
