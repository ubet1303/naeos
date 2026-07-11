package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	w := NewWatcher(time.Second, func(path string) {})
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
	if w.interval != time.Second {
		t.Errorf("expected interval 1s, got %v", w.interval)
	}
}

func TestNewWatcherDefaultInterval(t *testing.T) {
	w := NewWatcher(0, func(path string) {})
	if w.interval != 500*time.Millisecond {
		t.Errorf("expected default interval 500ms, got %v", w.interval)
	}
}

func TestAddDirectory(t *testing.T) {
	dir := t.TempDir()
	w := NewWatcher(time.Second, func(path string) {})
	if err := w.AddDirectory(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.directories) != 1 {
		t.Errorf("expected 1 directory, got %d", len(w.directories))
	}
}

func TestAddDirectoryNotFound(t *testing.T) {
	w := NewWatcher(time.Second, func(path string) {})
	err := w.AddDirectory("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}

func TestAddFileAsDirectory(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(tmpFile, []byte("test"), 0o600)
	w := NewWatcher(time.Second, func(path string) {})
	err := w.AddDirectory(tmpFile)
	if err == nil {
		t.Fatal("expected error for file instead of directory")
	}
}

func TestStartStop(t *testing.T) {
	w := NewWatcher(time.Second, func(path string) {})
	if w.IsRunning() {
		t.Fatal("expected not running initially")
	}
	if err := w.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.IsRunning() {
		t.Fatal("expected running after start")
	}
	w.Stop()
	if w.IsRunning() {
		t.Fatal("expected not running after stop")
	}
}

func TestStartTwice(t *testing.T) {
	w := NewWatcher(time.Second, func(path string) {})
	w.Start()
	err := w.Start()
	if err == nil {
		t.Fatal("expected error on double start")
	}
	w.Stop()
}

func TestSnapshot(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("content1"), 0o600)
	os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("content2"), 0o600)

	w := NewWatcher(time.Second, func(path string) {})
	w.AddDirectory(dir)

	snap, err := w.Snapshot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap) < 2 {
		t.Errorf("expected at least 2 entries, got %d", len(snap))
	}
}

func TestDetectChanges(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("content1"), 0o600)

	w := NewWatcher(time.Second, func(path string) {})
	w.AddDirectory(dir)

	snap, _ := w.Snapshot()

	os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("content2"), 0o600)

	changes := w.DetectChanges(snap)
	if len(changes) == 0 {
		t.Fatal("expected at least 1 change")
	}

	found := false
	for _, c := range changes {
		if c.EventType == "created" {
			found = true
		}
	}
	if !found {
		t.Error("expected to find a 'created' event")
	}
}
