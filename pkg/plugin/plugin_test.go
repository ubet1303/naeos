package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	mgr := NewManager(t.TempDir())
	if mgr == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestLoadConfigNonexistent(t *testing.T) {
	mgr := NewManager(t.TempDir())
	err := mgr.LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.LoadConfig()
	mgr.config.Plugins = append(mgr.config.Plugins, PluginInfo{
		Name:    "test-plugin",
		Version: "1.0.0",
		Enabled: true,
	})
	mgr.SaveConfig()

	mgr2 := NewManager(dir)
	mgr2.LoadConfig()
	if len(mgr2.config.Plugins) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(mgr2.config.Plugins))
	}
}

func TestListEmpty(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	plugins := mgr.List()
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestUninstall(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.LoadConfig()
	mgr.config.Plugins = append(mgr.config.Plugins, PluginInfo{Name: "test"})
	mgr.SaveConfig()

	err := mgr.Uninstall("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mgr2 := NewManager(dir)
	mgr2.LoadConfig()
	if len(mgr2.config.Plugins) != 0 {
		t.Errorf("expected 0 plugins after uninstall, got %d", len(mgr2.config.Plugins))
	}
}

func TestUninstallNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	err := mgr.Uninstall("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent plugin")
	}
}

func TestEnableDisable(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	mgr.LoadConfig()
	mgr.config.Plugins = append(mgr.config.Plugins, PluginInfo{Name: "test", Enabled: true})
	mgr.SaveConfig()

	mgr.Disable("test")
	mgr2 := NewManager(dir)
	mgr2.LoadConfig()
	if mgr2.config.Plugins[0].Enabled {
		t.Error("expected plugin to be disabled")
	}

	mgr2.Enable("test")
	mgr3 := NewManager(dir)
	mgr3.LoadConfig()
	if !mgr3.config.Plugins[0].Enabled {
		t.Error("expected plugin to be enabled")
	}
}

func TestEnableNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	err := mgr.Enable("nope")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDisableNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	err := mgr.Disable("nope")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigPath(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManager(dir)
	expected := filepath.Join(dir, "plugins.json")
	if mgr.configPath() != expected {
		t.Errorf("expected %s, got %s", expected, mgr.configPath())
	}
}

func TestGetNonexistent(t *testing.T) {
	mgr := NewManager(t.TempDir())
	_, ok := mgr.Get("nonexistent")
	if ok {
		t.Error("expected not found for nonexistent plugin")
	}
}

func TestLoadAllNoPlugins(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	err := mgr.LoadAll(&Context{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAllDisabledPlugin(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	mgr.config.Plugins = append(mgr.config.Plugins, PluginInfo{Name: "x", Enabled: false, Path: "/fake"})
	err := mgr.LoadAll(&Context{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAllEmptyPath(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.LoadConfig()
	mgr.config.Plugins = append(mgr.config.Plugins, PluginInfo{Name: "x", Enabled: true, Path: ""})
	err := mgr.LoadAll(&Context{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAllSandboxReject(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.config.Sandbox = SandboxConfig{AllowedDirs: []string{"/allowed"}}
	mgr.sandbox = NewSandbox(mgr.config.Sandbox)
	mgr.config.Plugins = append(mgr.config.Plugins, PluginInfo{Name: "x", Enabled: true, Path: "/not-allowed/plugin.so"})
	err := mgr.LoadAll(&Context{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteNotLoaded(t *testing.T) {
	mgr := NewManager(t.TempDir())
	_, err := mgr.Execute(context.Background(), "nope", nil)
	if err == nil {
		t.Fatal("expected error for not loaded plugin")
	}
}

func TestCleanupEmpty(t *testing.T) {
	mgr := NewManager(t.TempDir())
	err := mgr.Cleanup()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCleanupError(t *testing.T) {
	mgr := NewManager(t.TempDir())
	mgr.plugins["bad"] = &failingPlugin{cleanupErr: fmt.Errorf("cleanup failed")}
	err := mgr.Cleanup()
	if err == nil {
		t.Fatal("expected error from failing plugin cleanup")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "plugins.json"), []byte("{bad json"), 0o600)
	mgr := NewManager(dir)
	err := mgr.LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadConfigWithSandbox(t *testing.T) {
	dir := t.TempDir()
	cfg := PluginConfig{
		Sandbox: SandboxConfig{
			AllowedDirs: []string{"/tmp"},
			ExecTimeout: 5 * time.Second,
			MaxCalls:    10,
		},
	}
	data, _ := json.Marshal(cfg)
	os.WriteFile(filepath.Join(dir, "plugins.json"), data, 0o600)

	mgr := NewManager(dir)
	err := mgr.LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mgr.sandbox == nil {
		t.Fatal("expected sandbox to be set")
	}
}

func TestSaveConfigMkdirFail(t *testing.T) {
	mgr := NewManager("/nonexistent/path/that/does/not/exist")
	err := mgr.SaveConfig()
	if err == nil {
		t.Fatal("expected error saving to nonexistent dir")
	}
}

// --- Sandbox Tests ---

func TestNewSandboxDefaults(t *testing.T) {
	s := NewSandbox(SandboxConfig{})
	if s.config.ExecTimeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", s.config.ExecTimeout)
	}
	if s.config.MaxCalls != 1000 {
		t.Errorf("expected 1000 max calls, got %d", s.config.MaxCalls)
	}
}

func TestSandboxValidatePathNoRestrictions(t *testing.T) {
	s := NewSandbox(SandboxConfig{})
	err := s.ValidatePath("/any/path")
	if err != nil {
		t.Errorf("expected no error with no restrictions: %v", err)
	}
}

func TestSandboxValidatePathAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	s := NewSandbox(SandboxConfig{AllowedDirs: []string{tmpDir}})
	err := s.ValidatePath(filepath.Join(tmpDir, "plugin.so"))
	if err != nil {
		t.Errorf("expected allowed path: %v", err)
	}
}

func TestSandboxValidatePathExactMatch(t *testing.T) {
	tmpDir := t.TempDir()
	s := NewSandbox(SandboxConfig{AllowedDirs: []string{tmpDir}})
	err := s.ValidatePath(tmpDir)
	if err != nil {
		t.Errorf("expected exact dir match: %v", err)
	}
}

func TestSandboxValidatePathDenied(t *testing.T) {
	s := NewSandbox(SandboxConfig{AllowedDirs: []string{"/allowed"}})
	err := s.ValidatePath("/not-allowed/file.so")
	if err == nil {
		t.Fatal("expected error for denied path")
	}
}

func TestSandboxCheckRateLimit(t *testing.T) {
	s := NewSandbox(SandboxConfig{MaxCalls: 3})
	for i := 0; i < 3; i++ {
		if err := s.CheckRateLimit("p"); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i+1, err)
		}
	}
	if err := s.CheckRateLimit("p"); err == nil {
		t.Fatal("expected rate limit exceeded")
	}
}

func TestSandboxCheckRateLimitSeparatePlugins(t *testing.T) {
	s := NewSandbox(SandboxConfig{MaxCalls: 1})
	if err := s.CheckRateLimit("a"); err != nil {
		t.Fatal("unexpected error")
	}
	if err := s.CheckRateLimit("b"); err != nil {
		t.Fatal("separate plugin should have separate counter")
	}
}

func TestExecuteWithTimeoutSuccess(t *testing.T) {
	s := NewSandbox(SandboxConfig{})
	result, err := s.ExecuteWithTimeout(context.Background(), func() (any, error) {
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected 'ok', got %v", result)
	}
}

func TestExecuteWithTimeoutError(t *testing.T) {
	s := NewSandbox(SandboxConfig{})
	_, err := s.ExecuteWithTimeout(context.Background(), func() (any, error) {
		return nil, fmt.Errorf("plugin error")
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExecuteWithTimeoutCancellation(t *testing.T) {
	s := NewSandbox(SandboxConfig{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel

	_, err := s.ExecuteWithTimeout(ctx, func() (any, error) {
		time.Sleep(10 * time.Second)
		return "ok", nil
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

// stubPlugin is a minimal Plugin implementation for testing
type stubPlugin struct {
	name    string
	version string
}

func (p *stubPlugin) Name() string                          { return p.name }
func (p *stubPlugin) Version() string                       { return p.version }
func (p *stubPlugin) Initialize(ctx *Context) error         { return nil }
func (p *stubPlugin) Execute(input any) (any, error)        { return nil, nil }
func (p *stubPlugin) Cleanup() error                        { return nil }

type failingPlugin struct {
	cleanupErr error
}

func (p *failingPlugin) Name() string                      { return "failing" }
func (p *failingPlugin) Version() string                   { return "0.0.0" }
func (p *failingPlugin) Initialize(ctx *Context) error     { return nil }
func (p *failingPlugin) Execute(input any) (any, error)    { return nil, nil }
func (p *failingPlugin) Cleanup() error                    { return p.cleanupErr }
