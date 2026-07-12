package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatalf("failed to create dashboard: %v", err)
	}
	if d == nil {
		t.Error("expected non-nil dashboard")
	}
}

func TestServeHTTP(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatalf("failed to create dashboard: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	d.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html, got %s", ct)
	}
}

func TestGetStatsDefaults(t *testing.T) {
	statsMu.Lock()
	globalStats = Stats{}
	statsMu.Unlock()

	s := GetStats()
	if s.Projects != 0 {
		t.Errorf("expected 0 projects, got %d", s.Projects)
	}
	if s.Artifacts != 0 {
		t.Errorf("expected 0 artifacts, got %d", s.Artifacts)
	}
	if s.Pipelines != 0 {
		t.Errorf("expected 0 pipelines, got %d", s.Pipelines)
	}
}

func TestRecordPipelineRun(t *testing.T) {
	statsMu.Lock()
	globalStats = Stats{}
	statsMu.Unlock()

	RecordPipelineRun()
	RecordPipelineRun()

	s := GetStats()
	if s.Pipelines != 2 {
		t.Errorf("expected 2 pipelines, got %d", s.Pipelines)
	}
	if s.LastRun == "" {
		t.Error("expected non-empty last run")
	}
}

func TestSetProjects(t *testing.T) {
	SetProjects(5)
	s := GetStats()
	if s.Projects != 5 {
		t.Errorf("expected 5 projects, got %d", s.Projects)
	}
}

func TestSetArtifacts(t *testing.T) {
	SetArtifacts(10)
	s := GetStats()
	if s.Artifacts != 10 {
		t.Errorf("expected 10 artifacts, got %d", s.Artifacts)
	}
}

func TestStatsPersistence(t *testing.T) {
	statsMu.Lock()
	globalStats = Stats{}
	statsFile = ""
	statsMu.Unlock()

	dir := t.TempDir()
	path := filepath.Join(dir, "stats.json")

	SetStatsFile(path)

	SetProjects(3)
	SetArtifacts(7)
	RecordPipelineRun()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read stats file: %v", err)
	}

	var loaded Stats
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to unmarshal stats: %v", err)
	}
	if loaded.Projects != 3 {
		t.Errorf("expected 3 projects in file, got %d", loaded.Projects)
	}
	if loaded.Artifacts != 7 {
		t.Errorf("expected 7 artifacts in file, got %d", loaded.Artifacts)
	}
	if loaded.Pipelines != 1 {
		t.Errorf("expected 1 pipeline in file, got %d", loaded.Pipelines)
	}

	statsMu.Lock()
	globalStats = Stats{}
	statsFile = ""
	statsMu.Unlock()
}

func TestStatsLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats.json")

	data := Stats{Projects: 11, Artifacts: 22, Pipelines: 5, LastRun: "2026-01-01T00:00:00Z"}
	b, _ := json.Marshal(data)
	os.WriteFile(path, b, 0o644)

	SetStatsFile(path)
	s := GetStats()
	if s.Projects != 11 {
		t.Errorf("expected 11 projects from file, got %d", s.Projects)
	}
	if s.Artifacts != 22 {
		t.Errorf("expected 22 artifacts from file, got %d", s.Artifacts)
	}

	statsMu.Lock()
	globalStats = Stats{}
	statsFile = ""
	statsMu.Unlock()
}
