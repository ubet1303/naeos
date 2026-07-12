// Deprecated: Use github.com/NAEOS-foundation/naeos/internal/pluginhost instead.
// This package is a thin wrapper around internal/pluginhost for backward compatibility.
package pluginsdk

import (
	"github.com/NAEOS-foundation/naeos/internal/pluginhost"
)

// Deprecated: Use pluginhost.Plugin instead.
type Plugin = pluginhost.Plugin

// Deprecated: Use pluginhost.PluginContext instead.
type PluginContext = pluginhost.PluginContext

// Deprecated: Use pluginhost.Logger instead.
type Logger = pluginhost.Logger

// Deprecated: Use pluginhost.MetricsCollector instead.
type MetricsCollector = pluginhost.MetricsCollector

// Deprecated: Use pluginhost.EventEmitter instead.
type EventEmitter = pluginhost.EventEmitter

// Deprecated: Use pluginhost.Manager instead.
type Manager = pluginhost.Manager

// Deprecated: Use pluginhost.Manifest instead.
type Manifest = pluginhost.Manifest

// Deprecated: Use pluginhost.ActionManifest instead.
type ActionManifest = pluginhost.ActionManifest

// Deprecated: Use pluginhost.ConfigField instead.
type ConfigField = pluginhost.ConfigField

// Deprecated: Use pluginhost.BasePlugin instead.
type BasePlugin = pluginhost.BasePlugin

// Deprecated: Use pluginhost.PluginState instead.
type PluginState = pluginhost.PluginState

// Deprecated: Use pluginhost.PluginInfo instead.
type PluginInfo = pluginhost.PluginInfo

// Deprecated: Use pluginhost.StateCreated instead.
const StateCreated = pluginhost.StateCreated

// Deprecated: Use pluginhost.StateInitialized instead.
const StateInitialized = pluginhost.StateInitialized

// Deprecated: Use pluginhost.StateRunning instead.
const StateRunning = pluginhost.StateRunning

// Deprecated: Use pluginhost.StateStopped instead.
const StateStopped = pluginhost.StateStopped

// Deprecated: Use pluginhost.StateError instead.
const StateError = pluginhost.StateError

// Deprecated: Use pluginhost.NewManager instead.
var NewManager = pluginhost.NewManager

// Deprecated: Use pluginhost.NewSimpleLogger instead.
var NewSimpleLogger = pluginhost.NewSimpleLogger

// Deprecated: Use pluginhost.NewSimpleMetrics instead.
var NewSimpleMetrics = pluginhost.NewSimpleMetrics

// Deprecated: Use pluginhost.NewSimpleEventEmitter instead.
var NewSimpleEventEmitter = pluginhost.NewSimpleEventEmitter
