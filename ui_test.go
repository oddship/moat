package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUIAssets(t *testing.T) {
	for _, tc := range []struct {
		name, base, config, css, js string
	}{
		{"default", "", "", "/_static/oatbase/foundation.min.css", "/_static/oatbase/foundation.min.js"},
		{"disabled", "", "[ui]\nenabled = false", "", ""},
		{"CDN", "/project", "[ui]\nassets = 'cdn'", Config{UI: UIConfig{Assets: "cdn"}}.UIStylesheet(), Config{UI: UIConfig{Assets: "cdn"}}.UIScript()},
		{"vendored project", "/project", "", "/project/_static/oatbase/foundation.min.css", "/project/_static/oatbase/foundation.min.js"},
		{"local root", "", "[ui]\nstylesheet = '_static/ui.css'\nscript = '_static/ui.js'", "/_static/ui.css", "/_static/ui.js"},
		{"local project", "/project", "[ui]\nstylesheet = '_static/ui.css'\nscript = '_static/ui.js'", "/project/_static/ui.css", "/project/_static/ui.js"},
		{"remote", "/project", "[ui]\nstylesheet = 'https://example.org/ui.css'\nscript = 'https://example.org/ui.js'", "https://example.org/ui.css", "https://example.org/ui.js"},
		{"script disabled", "", "[ui]\ndisable_script = true", "/_static/oatbase/foundation.min.css", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src, dst := t.TempDir(), t.TempDir()
			if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Hello"), 0600); err != nil {
				t.Fatal(err)
			}
			configPath := filepath.Join(src, "config.toml")
			if err := os.WriteFile(configPath, []byte(tc.config), 0600); err != nil {
				t.Fatal(err)
			}
			cfg, err := LoadConfig(configPath)
			if err != nil {
				t.Fatal(err)
			}
			cfg.BasePath = tc.base
			if err := Build(src, dst, cfg); err != nil {
				t.Fatal(err)
			}
			b, err := os.ReadFile(filepath.Join(dst, "index.html"))
			if err != nil {
				t.Fatal(err)
			}
			html := string(b)
			if tc.css != "" && !strings.Contains(html, `href="`+tc.css+`"`) {
				t.Fatal("missing stylesheet")
			}
			if tc.js != "" && strings.Count(html, `src="`+tc.js+`"`) != 1 {
				t.Fatal("UI script must load exactly once")
			}
			if tc.js == "" && strings.Contains(html, `<script src=`) {
				t.Fatal("UI script must be omitted")
			}
			_, assetErr := os.Stat(filepath.Join(dst, "_static/oatbase/foundation.min.css"))
			if cfg.UIEnabled() && cfg.UI.Assets != "cdn" && assetErr != nil {
				t.Fatal("missing vendored foundation")
			}
			if (!cfg.UIEnabled() || cfg.UI.Assets == "cdn") && !os.IsNotExist(assetErr) {
				t.Fatal("unwanted vendored assets")
			}
			if !cfg.UIEnabled() && strings.Contains(html, "foundation.min") {
				t.Fatal("UI disabled but foundation included")
			}
			if strings.Contains(html, "@knadh/oat") {
				t.Fatal("duplicate Oat foundation")
			}
			if !strings.Contains(html, `<article data-prose>`) {
				t.Fatal("missing Oatbase prose")
			}
		})
	}
}

func TestRejectUnknownUIAssetMode(t *testing.T) {
	if err := Build(t.TempDir(), t.TempDir(), Config{UI: UIConfig{Assets: "typo"}}); err == nil {
		t.Fatal("unknown asset modes must fail")
	}
}
