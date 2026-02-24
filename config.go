package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds site-level configuration from config.toml.
type Config struct {
	SiteName  string            `toml:"site_name"`
	BasePath  string            `toml:"base_path"`
	Highlight HighlightConfig   `toml:"highlight"`
	Links     []LinkConfig      `toml:"links"`
	Extra     map[string]any    `toml:"extra"`
}

// LinkConfig is a sidebar link above the nav.
type LinkConfig struct {
	Title string `toml:"title"`
	URL   string `toml:"url"`
}

// HighlightConfig controls syntax highlighting themes.
type HighlightConfig struct {
	Light string `toml:"light"` // Chroma style for light mode (default: "github")
	Dark  string `toml:"dark"`  // Chroma style for dark mode (default: "github-dark")
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
