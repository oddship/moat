package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds site-level configuration from config.toml.
type Config struct {
	SiteName            string          `toml:"site_name"`
	BasePath            string          `toml:"base_path"`
	Logo                string          `toml:"logo"`
	Favicon             string          `toml:"favicon"`
	FooterText          string          `toml:"footer_text"`
	DisableMoatCitation bool            `toml:"disable_moat_citation"`
	UI                  UIConfig        `toml:"ui"`
	Highlight           HighlightConfig `toml:"highlight"`
	Links               []LinkConfig    `toml:"links"`
	TopNav              []LinkConfig    `toml:"topnav"`
	TopNavMore          []LinkConfig    `toml:"topnav_more"`
	CurrentVersion      string          `toml:"current_version"`
	Versions            []VersionConfig `toml:"versions"`
	Search              SearchConfig    `toml:"search"`
	Agents              AgentsConfig    `toml:"agents"`
	Feed                FeedConfig      `toml:"feed"`
	Extra               map[string]any  `toml:"extra"`
}

// UIConfig controls the built-in layout's lightweight UI assets.
type UIConfig struct {
	Enabled       *bool  `toml:"enabled"`
	Assets        string `toml:"assets"`
	Stylesheet    string `toml:"stylesheet"`
	Script        string `toml:"script"`
	DisableScript bool   `toml:"disable_script"`
}

func (c Config) UIEnabled() bool { return c.UI.Enabled == nil || *c.UI.Enabled }

func (c Config) validateUI() error {
	if c.UI.Assets != "" && c.UI.Assets != "vendored" && c.UI.Assets != "cdn" {
		return fmt.Errorf("ui.assets must be vendored or cdn")
	}
	return nil
}

// The CDN and embedded copies identify the same immutable upstream revision.
// Switch this to an npm version only after that version is published and verified.
func (c Config) UIAssetBase() string {
	if c.UI.Assets == "cdn" {
		return "https://cdn.jsdelivr.net/gh/oddship/oatbase@0e4b8945c97fec72e050a0ca5d8b6c5b0de38206/dist/"
	}
	return "_static/oatbase/"
}

func (c Config) UIStylesheet() string {
	if !c.UIEnabled() {
		return ""
	}
	if c.UI.Stylesheet != "" {
		return c.UI.Stylesheet
	}
	if c.UI.Assets == "cdn" {
		return c.UIAssetBase() + "foundation.css"
	}
	return c.UIAssetBase() + "foundation.min.css"
}

func (c Config) UIScript() string {
	if !c.UIEnabled() || c.UI.DisableScript {
		return ""
	}
	if c.UI.Script != "" {
		return c.UI.Script
	}
	if c.UI.Assets == "cdn" {
		return c.UIAssetBase() + "foundation.js"
	}
	return c.UIAssetBase() + "foundation.min.js"
}

// VersionConfig links to another independently built version of a site.
// URLs are site-root-relative or absolute, independent of BasePath.
type VersionConfig struct {
	Title string `toml:"title"`
	URL   string `toml:"url"`
}

func (c Config) validateVersions() error {
	titles := make(map[string]bool, len(c.Versions))
	for _, version := range c.Versions {
		if strings.TrimSpace(version.Title) == "" || titles[version.Title] {
			return fmt.Errorf("versions must have nonempty, unique titles")
		}
		titles[version.Title] = true
		u, err := url.Parse(version.URL)
		if err != nil || u.Opaque != "" || strings.Contains(u.Path, "\\") ||
			!((u.Scheme == "" && u.Host == "" && strings.HasPrefix(u.Path, "/") && !strings.HasPrefix(u.Path, "//")) ||
				((u.Scheme == "https" || u.Scheme == "http") && u.Host != "")) {
			return fmt.Errorf("version %q URL must be root-relative or use http(s)", version.Title)
		}
	}
	if c.CurrentVersion != "" && !titles[c.CurrentVersion] {
		return fmt.Errorf("current_version %q is missing from versions", c.CurrentVersion)
	}
	return nil
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

// AgentsConfig controls static Markdown exports and their built-in page actions.
type AgentsConfig struct {
	Enabled *bool `toml:"enabled"`
	Full    bool  `toml:"full"`
}

// AgentsEnabled defaults to true, including on sites with custom layouts.
func (c Config) AgentsEnabled() bool {
	return c.Agents.Enabled == nil || *c.Agents.Enabled
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
