package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme          string `yaml:"theme"`
	RefreshRate    int    `yaml:"refresh_rate"` // in milliseconds
	DefaultPage    string `yaml:"default_page"`
	RemoteMode     bool   `yaml:"remote_mode"`     // Enable remote monitoring
	RemoteURL      string `yaml:"remote_url"`      // Remote server URL
	RemoteServer   bool   `yaml:"remote_server"`   // Run as server
	RemotePort     string `yaml:"remote_port"`     // Server port
	PluginsEnabled bool   `yaml:"plugins_enabled"` // Enable plugin system
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "osm", "config.yaml")

	// Default config
	cfg := &Config{
		Theme:          "default",
		RefreshRate:    1000,
		DefaultPage:    "procs",
		RemoteMode:     false,
		RemoteURL:      "",
		RemoteServer:   false,
		RemotePort:     "8080",
		PluginsEnabled: false,
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return cfg, nil // Return default if no config exists
	} else if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
