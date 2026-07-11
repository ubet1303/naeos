package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	plugin "plugin"
	"strings"
	"sync"
	"time"

	naeoslog "github.com/NAEOS-foundation/naeos/internal/shared/log"
)

type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Path        string `json:"path,omitempty"`
	Enabled     bool   `json:"enabled"`
}

type Plugin interface {
	Name() string
	Version() string
	Initialize(ctx *Context) error
	Execute(input any) (any, error)
	Cleanup() error
}

type Context struct {
	ConfigDir string
	OutputDir string
	Verbose   bool
}

type SandboxConfig struct {
	AllowedDirs  []string      `json:"allowed_dirs,omitempty"`
	ExecTimeout  time.Duration `json:"exec_timeout,omitempty"`
	MaxCalls     int           `json:"max_calls,omitempty"`
}

type Sandbox struct {
	config  SandboxConfig
	mu      sync.Mutex
	callCnt map[string]int
}

func NewSandbox(cfg SandboxConfig) *Sandbox {
	if cfg.ExecTimeout <= 0 {
		cfg.ExecTimeout = 30 * time.Second
	}
	if cfg.MaxCalls <= 0 {
		cfg.MaxCalls = 1000
	}
	return &Sandbox{
		config:  cfg,
		callCnt: make(map[string]int),
	}
}

func (s *Sandbox) ValidatePath(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if len(s.config.AllowedDirs) == 0 {
		return nil
	}
	for _, dir := range s.config.AllowedDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(abs, absDir+string(filepath.Separator)) || abs == absDir {
			return nil
		}
	}
	return fmt.Errorf("plugin path %s is outside allowed directories", path)
}

func (s *Sandbox) CheckRateLimit(pluginName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callCnt[pluginName]++
	if s.callCnt[pluginName] > s.config.MaxCalls {
		return fmt.Errorf("plugin %s exceeded max call limit (%d)", pluginName, s.config.MaxCalls)
	}
	return nil
}

func (s *Sandbox) ExecuteWithTimeout(ctx context.Context, fn func() (any, error)) (any, error) {
	type result struct {
		value any
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		v, err := fn()
		ch <- result{v, err}
	}()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("plugin execution timed out: %w", ctx.Err())
	case r := <-ch:
		return r.value, r.err
	}
}

type PluginManager struct {
	pluginDir string
	plugins   map[string]Plugin
	config    PluginConfig
	sandbox   *Sandbox
}

type PluginConfig struct {
	Plugins []PluginInfo `json:"plugins"`
	Sandbox SandboxConfig `json:"sandbox,omitempty"`
}

func NewManager(pluginDir string) *PluginManager {
	return &PluginManager{
		pluginDir: pluginDir,
		plugins:   make(map[string]Plugin),
		sandbox:   NewSandbox(SandboxConfig{}),
	}
}

func (m *PluginManager) configPath() string {
	return filepath.Join(m.pluginDir, "plugins.json")
}

func (m *PluginManager) LoadConfig() error {
	data, err := os.ReadFile(m.configPath())
	if err != nil {
		if os.IsNotExist(err) {
			m.config = PluginConfig{}
			return nil
		}
		return err
	}
	if err := json.Unmarshal(data, &m.config); err != nil {
		return err
	}
	m.sandbox = NewSandbox(m.config.Sandbox)
	return nil
}

func (m *PluginManager) SaveConfig() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(m.pluginDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(m.configPath(), data, 0o600)
}

func (m *PluginManager) List() []PluginInfo {
	return m.config.Plugins
}

func (m *PluginManager) Install(path string) (*PluginInfo, error) {
	info, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open plugin %s: %w", path, err)
	}

	symName, err := info.Lookup("PluginName")
	if err != nil {
		return nil, fmt.Errorf("plugin %s does not export PluginName: %w", path, err)
	}
	namePtr, ok := symName.(*string)
	if !ok {
		return nil, fmt.Errorf("PluginName is not *string")
	}

	symVersion, err := info.Lookup("PluginVersion")
	version := "0.0.0"
	if err == nil {
		if vPtr, ok := symVersion.(*string); ok {
			version = *vPtr
		}
	}

	symDesc, err := info.Lookup("PluginDescription")
	description := ""
	if err == nil {
		if dPtr, ok := symDesc.(*string); ok {
			description = *dPtr
		}
	}

	pInfo := PluginInfo{
		Name:        *namePtr,
		Version:     version,
		Description: description,
		Path:        path,
		Enabled:     true,
	}

	for i, p := range m.config.Plugins {
		if p.Name == pInfo.Name {
			m.config.Plugins[i] = pInfo
			return &pInfo, m.SaveConfig()
		}
	}

	m.config.Plugins = append(m.config.Plugins, pInfo)
	return &pInfo, m.SaveConfig()
}

func (m *PluginManager) Uninstall(name string) error {
	for i, p := range m.config.Plugins {
		if p.Name == name {
			m.config.Plugins = append(m.config.Plugins[:i], m.config.Plugins[i+1:]...)
			delete(m.plugins, name)
			return m.SaveConfig()
		}
	}
	return fmt.Errorf("plugin %s not found", name)
}

func (m *PluginManager) Enable(name string) error {
	for i, p := range m.config.Plugins {
		if p.Name == name {
			m.config.Plugins[i].Enabled = true
			return m.SaveConfig()
		}
	}
	return fmt.Errorf("plugin %s not found", name)
}

func (m *PluginManager) Disable(name string) error {
	for i, p := range m.config.Plugins {
		if p.Name == name {
			m.config.Plugins[i].Enabled = false
			return m.SaveConfig()
		}
	}
	return fmt.Errorf("plugin %s not found", name)
}

func (m *PluginManager) Get(name string) (Plugin, bool) {
	p, ok := m.plugins[name]
	return p, ok
}

func (m *PluginManager) Execute(ctx context.Context, name string, input any) (any, error) {
	p, ok := m.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not loaded", name)
	}
	if err := m.sandbox.CheckRateLimit(name); err != nil {
		return nil, err
	}
	return m.sandbox.ExecuteWithTimeout(ctx, func() (any, error) {
		return p.Execute(input)
	})
}

func (m *PluginManager) LoadAll(ctx *Context) error {
	for _, pInfo := range m.config.Plugins {
		if !pInfo.Enabled || pInfo.Path == "" {
			continue
		}
		if err := m.sandbox.ValidatePath(pInfo.Path); err != nil {
			naeoslog.Warn("plugin path rejected by sandbox", "plugin", pInfo.Name, "error", err)
			continue
		}
		p, err := m.loadGoPlugin(pInfo.Path)
		if err != nil {
			naeoslog.Warn("failed to load plugin", "plugin", pInfo.Name, "error", err)
			continue
		}
		if err := p.Initialize(ctx); err != nil {
			naeoslog.Warn("failed to initialize plugin", "plugin", pInfo.Name, "error", err)
			continue
		}
		m.plugins[pInfo.Name] = p
	}
	return nil
}

func (m *PluginManager) loadGoPlugin(path string) (Plugin, error) {
	goPlugin, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	sym, err := goPlugin.Lookup("NaeosPlugin")
	if err != nil {
		return nil, fmt.Errorf("plugin does not export NaeosPlugin: %w", err)
	}

	p, ok := sym.(Plugin)
	if !ok {
		return nil, fmt.Errorf("NaeosPlugin does not implement plugin.Plugin interface")
	}

	return p, nil
}

func (m *PluginManager) Cleanup() error {
	var errs []string
	for name, p := range m.plugins {
		if err := p.Cleanup(); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
