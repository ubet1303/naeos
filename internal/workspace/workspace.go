package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Workspace struct {
	Name     string         `json:"name"`
	Root     string         `json:"root"`
	Modules  []ModuleRef    `json:"modules"`
}

type ModuleRef struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	SpecFile string `json:"spec_file,omitempty"`
}

type WorkspaceManager struct {
	rootDir string
}

func NewManager(rootDir string) *WorkspaceManager {
	return &WorkspaceManager{rootDir: rootDir}
}

func (m *WorkspaceManager) configPath() string {
	return filepath.Join(m.rootDir, "naeos.workspace.json")
}

func (m *WorkspaceManager) Init(name string) (*Workspace, error) {
	ws := &Workspace{
		Name: name,
		Root: m.rootDir,
	}
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(m.rootDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(m.configPath(), data, 0o600); err != nil {
		return nil, err
	}
	return ws, nil
}

func (m *WorkspaceManager) Load() (*Workspace, error) {
	data, err := os.ReadFile(m.configPath())
	if err != nil {
		return nil, fmt.Errorf("no workspace found at %s: %w", m.rootDir, err)
	}
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

func (m *WorkspaceManager) AddModule(name, path, specFile string) error {
	ws, err := m.Load()
	if err != nil {
		return err
	}
	for _, mod := range ws.Modules {
		if mod.Name == name {
			return fmt.Errorf("module %s already exists", name)
		}
	}
	ws.Modules = append(ws.Modules, ModuleRef{
		Name:     name,
		Path:     path,
		SpecFile: specFile,
	})
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath(), data, 0o600)
}

func (m *WorkspaceManager) RemoveModule(name string) error {
	ws, err := m.Load()
	if err != nil {
		return err
	}
	for i, mod := range ws.Modules {
		if mod.Name == name {
			ws.Modules = append(ws.Modules[:i], ws.Modules[i+1:]...)
			data, err := json.MarshalIndent(ws, "", "  ")
			if err != nil {
				return err
			}
			return os.WriteFile(m.configPath(), data, 0o600)
		}
	}
	return fmt.Errorf("module %s not found", name)
}

func (m *WorkspaceManager) ListModules() ([]ModuleRef, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	return ws.Modules, nil
}
