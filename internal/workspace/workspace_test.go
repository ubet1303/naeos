package workspace

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	mgr := NewManager(t.TempDir())
	if mgr == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestInitWorkspace(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	ws, err := mgr.Init("test-workspace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Name != "test-workspace" {
		t.Errorf("expected test-workspace, got %s", ws.Name)
	}
}

func TestLoadWorkspace(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.Init("test-workspace")

	ws, err := mgr.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Name != "test-workspace" {
		t.Errorf("expected test-workspace, got %s", ws.Name)
	}
}

func TestLoadNonexistentWorkspace(t *testing.T) {
	mgr := NewManager(t.TempDir())
	_, err := mgr.Load()
	if err == nil {
		t.Fatal("expected error for nonexistent workspace")
	}
}

func TestAddModule(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.Init("test-workspace")

	err := mgr.AddModule("core", "./internal/core", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modules, _ := mgr.ListModules()
	if len(modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(modules))
	}
}

func TestAddDuplicateModule(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.Init("test-workspace")
	mgr.AddModule("core", "./internal/core", "")

	err := mgr.AddModule("core", "./internal/core2", "")
	if err == nil {
		t.Fatal("expected error for duplicate module")
	}
}

func TestRemoveModule(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.Init("test-workspace")
	mgr.AddModule("core", "./internal/core", "")

	err := mgr.RemoveModule("core")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modules, _ := mgr.ListModules()
	if len(modules) != 0 {
		t.Errorf("expected 0 modules, got %d", len(modules))
	}
}

func TestRemoveNonexistentModule(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.Init("test-workspace")

	err := mgr.RemoveModule("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent module")
	}
}

func TestListModules(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.Init("test-workspace")
	mgr.AddModule("core", "./internal/core", "")
	mgr.AddModule("auth", "./internal/auth", "")

	modules, err := mgr.ListModules()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modules) != 2 {
		t.Errorf("expected 2 modules, got %d", len(modules))
	}
}
