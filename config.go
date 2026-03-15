package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds site-level configuration from config.toml.
type Config struct {
	SiteName  string          `toml:"site_name"`
	BasePath  string          `toml:"base_path"`
	Logo      string          `toml:"logo"`
	Favicon   string          `toml:"favicon"`
	Highlight HighlightConfig `toml:"highlight"`
	Links     []LinkConfig    `toml:"links"`
	Search    SearchConfig    `toml:"search"`
	Extra     map[string]any  `toml:"extra"`
}

// LinkConfig is a sidebar link above the nav.
type LinkConfig struct {
	Title string `toml:"title"`
	URL   string `toml:"url"`
	Icon  string `toml:"icon"` // Optional built-in icon name (e.g. "github")
}

// HighlightConfig controls syntax highlighting themes.
type HighlightConfig struct {
	Light string `toml:"light"` // Chroma style for light mode (default: "github")
	Dark  string `toml:"dark"`  // Chroma style for dark mode (default: "github-dark")
}

// SearchConfig controls built-in client-side search.
type SearchConfig struct {
	Enabled *bool `toml:"enabled"`
}

// SearchEnabled returns the effective search setting.
// Search defaults to enabled when omitted from config.toml.
func (c Config) SearchEnabled() bool {
	if c.Search.Enabled == nil {
		return true
	}
	return *c.Search.Enabled
}

// LoadConfig reads a TOML config file. Returns zero Config if path is empty or file doesn't exist.
func LoadConfig(path string) (Config, error) {
	var cfg Config
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return cfg, nil
}
