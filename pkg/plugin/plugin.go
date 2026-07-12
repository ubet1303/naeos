// Deprecated: Use github.com/NAEOS-foundation/naeos/internal/pluginhost instead.
// This package is a thin wrapper around internal/pluginhost for backward compatibility.
package plugin

import (
	"context"

	"github.com/NAEOS-foundation/naeos/internal/pluginhost"
)

// Deprecated: Use pluginhost.Plugin instead.
type Plugin = pluginhost.Plugin

// Deprecated: Use pluginhost.PluginInfo instead.
type PluginInfo = pluginhost.PluginInfo

// Deprecated: Use pluginhost.PluginConfig instead.
type PluginConfig = pluginhost.PluginConfig

// Deprecated: Use pluginhost.SandboxConfig instead.
type SandboxConfig = pluginhost.SandboxConfig

// Deprecated: Use pluginhost.Sandbox instead.
type Sandbox = pluginhost.Sandbox

// Deprecated: Use pluginhost.PluginManager or pluginhost.Manager instead.
type PluginManager = pluginhost.Manager

// Deprecated: Use pluginhost.PluginContext instead.
type Context = pluginhost.PluginContext

// Deprecated: Use pluginhost.NewSandbox instead.
var NewSandbox = pluginhost.NewSandbox

// Deprecated: Use pluginhost.NewManager instead.
var NewManager = pluginhost.NewManager

// Execute runs a plugin action through the unified manager.
func Execute(ctx context.Context, mgr *pluginhost.Manager, name string, input any) (any, error) {
	return mgr.Execute(ctx, name, "execute", map[string]any{"input": input})
}
