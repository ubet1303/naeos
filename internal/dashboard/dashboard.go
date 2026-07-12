package dashboard

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed templates/*
var templatesFS embed.FS

type Dashboard struct {
	templates *template.Template
}

func New() (*Dashboard, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Dashboard{
		templates: tmpl,
	}, nil
}

func (d *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	d.templates.ExecuteTemplate(w, "index.html", nil)
}

type Stats struct {
	Projects  int    `json:"projects"`
	Artifacts int    `json:"artifacts"`
	Pipelines int    `json:"pipelines"`
	LastRun   string `json:"last_run"`
}

var (
	globalStats  Stats
	statsMu      sync.RWMutex
	statsFile    string
)

func GetStats() *Stats {
	statsMu.RLock()
	defer statsMu.RUnlock()
	s := globalStats
	return &s
}

func RecordPipelineRun() {
	statsMu.Lock()
	defer statsMu.Unlock()
	globalStats.Pipelines++
	globalStats.LastRun = time.Now().Format(time.RFC3339)
	persistStats()
}

func SetProjects(n int) {
	statsMu.Lock()
	defer statsMu.Unlock()
	globalStats.Projects = n
	persistStats()
}

func SetArtifacts(n int) {
	statsMu.Lock()
	defer statsMu.Unlock()
	globalStats.Artifacts = n
	persistStats()
}

func SetStatsFile(path string) {
	statsMu.Lock()
	defer statsMu.Unlock()
	statsFile = path
	loadStats()
}

func persistStats() {
	if statsFile == "" {
		return
	}
	data, err := json.MarshalIndent(globalStats, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(statsFile, data, 0o644)
}

func loadStats() {
	if statsFile == "" {
		return
	}
	data, err := os.ReadFile(statsFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &globalStats)
}
