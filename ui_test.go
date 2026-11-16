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
		{"default", "", "", "https://cdn.jsdelivr.net/npm/@oddship/oatbase@0.2.0/dist/oatbase.min.css", "https://cdn.jsdelivr.net/npm/@oddship/oatbase@0.2.0/dist/oatbase.min.js"},
		{"local root", "", "[ui]\nstylesheet = '_static/ui.css'\nscript = '_static/ui.js'", "/_static/ui.css", "/_static/ui.js"},
		{"local project", "/project", "[ui]\nstylesheet = '_static/ui.css'\nscript = '_static/ui.js'", "/project/_static/ui.css", "/project/_static/ui.js"},
		{"remote", "/project", "[ui]\nstylesheet = 'https://example.org/ui.css'\nscript = 'https://example.org/ui.js'", "https://example.org/ui.css", "https://example.org/ui.js"},
		{"script disabled", "", "[ui]\ndisable_script = true", "https://cdn.jsdelivr.net/npm/@oddship/oatbase@0.2.0/dist/oatbase.min.css", ""},
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
			if !strings.Contains(html, `href="`+tc.css+`"`) {
				t.Fatal("missing stylesheet")
			}
			if tc.js != "" && strings.Count(html, `src="`+tc.js+`"`) != 1 {
				t.Fatal("UI script must load exactly once")
			}
			if tc.js == "" && strings.Contains(html, `<script src=`) {
				t.Fatal("UI script must be omitted")
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
