package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Plugin represents a custom metric plugin
type Plugin struct {
	Name        string
	Path        string
	Type        string // "shell" or "lua"
	RefreshRate int    // seconds
	LastRun     time.Time
	LastOutput  string
	Enabled     bool
}

// PluginOutput represents the output from a plugin
type PluginOutput struct {
	Title   string            `json:"title"`
	Value   string            `json:"value"`
	Metrics map[string]string `json:"metrics,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// Manager manages all plugins
type Manager struct {
	Plugins    []*Plugin
	PluginsDir string
}

// NewManager creates a new plugin manager
func NewManager(pluginsDir string) *Manager {
	return &Manager{
		Plugins:    make([]*Plugin, 0),
		PluginsDir: pluginsDir,
	}
}

// LoadPlugins loads all plugins from the plugins directory
func (m *Manager) LoadPlugins() error {
	// Create plugins directory if it doesn't exist
	if err := os.MkdirAll(m.PluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	// Find all plugin files
	files, err := os.ReadDir(m.PluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		path := filepath.Join(m.PluginsDir, name)

		var pluginType string
		if strings.HasSuffix(name, ".sh") {
			pluginType = "shell"
		} else if strings.HasSuffix(name, ".lua") {
			pluginType = "lua"
		} else {
			continue // Skip non-plugin files
		}

		plugin := &Plugin{
			Name:        strings.TrimSuffix(name, filepath.Ext(name)),
			Path:        path,
			Type:        pluginType,
			RefreshRate: 5, // Default 5 seconds
			Enabled:     true,
		}

		m.Plugins = append(m.Plugins, plugin)
	}

	return nil
}

// ExecutePlugin runs a plugin and returns its output
func (m *Manager) ExecutePlugin(plugin *Plugin) (*PluginOutput, error) {
	if !plugin.Enabled {
		return nil, fmt.Errorf("plugin is disabled")
	}

	// Check if we should run based on refresh rate
	if time.Since(plugin.LastRun) < time.Duration(plugin.RefreshRate)*time.Second {
		// Return cached output
		if plugin.LastOutput != "" {
			var output PluginOutput
			if err := json.Unmarshal([]byte(plugin.LastOutput), &output); err == nil {
				return &output, nil
			}
		}
	}

	var cmd *exec.Cmd
	switch plugin.Type {
	case "shell":
		cmd = exec.Command("/bin/bash", plugin.Path)
	case "lua":
		cmd = exec.Command("lua", plugin.Path)
	default:
		return nil, fmt.Errorf("unsupported plugin type: %s", plugin.Type)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &PluginOutput{
			Title: plugin.Name,
			Error: fmt.Sprintf("Execution failed: %v", err),
		}, err
	}

	// Try to parse as JSON
	var pluginOutput PluginOutput
	if err := json.Unmarshal(output, &pluginOutput); err != nil {
		// If not JSON, treat as plain text
		pluginOutput = PluginOutput{
			Title: plugin.Name,
			Value: string(output),
		}
	}

	// Cache the output
	plugin.LastRun = time.Now()
	if jsonOutput, err := json.Marshal(pluginOutput); err == nil {
		plugin.LastOutput = string(jsonOutput)
	}

	return &pluginOutput, nil
}

// ExecuteAll runs all enabled plugins
func (m *Manager) ExecuteAll() map[string]*PluginOutput {
	results := make(map[string]*PluginOutput)

	for _, plugin := range m.Plugins {
		if !plugin.Enabled {
			continue
		}

		output, err := m.ExecutePlugin(plugin)
		if err != nil {
			output = &PluginOutput{
				Title: plugin.Name,
				Error: err.Error(),
			}
		}
		results[plugin.Name] = output
	}

	return results
}

// CreateExamplePlugins creates example plugin files
func CreateExamplePlugins(pluginsDir string) error {
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return err
	}

	// Example shell plugin
	shellExample := `#!/bin/bash
# Example OSM Plugin - System Uptime
# Output must be valid JSON

UPTIME=$(uptime -p)
cat <<EOF
{
  "title": "System Uptime",
  "value": "$UPTIME"
}
EOF
`
	shellPath := filepath.Join(pluginsDir, "uptime.sh")
	if err := os.WriteFile(shellPath, []byte(shellExample), 0755); err != nil {
		return err
	}

	// Example shell plugin with metrics
	metricsExample := `#!/bin/bash
# Example OSM Plugin - Custom Metrics
# Output JSON with multiple metrics

cat <<EOF
{
  "title": "Custom Metrics",
  "value": "System Info",
  "metrics": {
    "Kernel": "$(uname -r)",
    "Shell": "$SHELL",
    "User": "$USER"
  }
}
EOF
`
	metricsPath := filepath.Join(pluginsDir, "system_info.sh")
	if err := os.WriteFile(metricsPath, []byte(metricsExample), 0755); err != nil {
		return err
	}

	return nil
}
