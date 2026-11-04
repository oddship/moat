package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	content := `current_version = "v2"
[[versions]]
title = "v2"
url = "/docs/v2/"
[[versions]]
title = "v1"
url = "https://archive.example.com/v1/"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CurrentVersion != "v2" || len(cfg.Versions) != 2 || cfg.Versions[1].Title != "v1" {
		t.Fatalf("unexpected version config: %+v", cfg)
	}
	if err := cfg.validateVersions(); err != nil {
		t.Fatal(err)
	}
}

func TestVersionConfigRejectsInvalidOptions(t *testing.T) {
	tests := []Config{
		{CurrentVersion: "missing"},
		{Versions: []VersionConfig{{Title: "", URL: "/v1/"}}},
		{Versions: []VersionConfig{{Title: "v1", URL: "/v1/"}, {Title: "v1", URL: "/other/"}}},
	}
	for _, address := range []string{"javascript:alert(1)", "//example.com/", "v1/", "", "/\\example.com/"} {
		tests = append(tests, Config{Versions: []VersionConfig{{Title: "v1", URL: address}}})
	}
	for _, cfg := range tests {
		if err := cfg.validateVersions(); err == nil {
			t.Errorf("expected invalid configuration to fail: %+v", cfg)
		}
	}
}

func TestBuildVersionPicker(t *testing.T) {
	for _, current := range []string{"", "v2"} {
		t.Run("current="+current, func(t *testing.T) {
			src, dst := t.TempDir(), t.TempDir()
			if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Home\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			cfg := Config{BasePath: "/docs/v2", CurrentVersion: current, Versions: []VersionConfig{
				{Title: "v2", URL: "/docs/v2/"},
				{Title: "v1", URL: "https://archive.example.com/v1/"},
			}}
			if err := Build(src, dst, cfg); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(dst, "index.html"))
			if err != nil {
				t.Fatal(err)
			}
			page := string(data)
			for _, want := range []string{`<label for="moat-version">Version</label>`,
				`value="https://archive.example.com/v1/"`, `window.location.assign(versionSelect.value)`} {
				if !strings.Contains(page, want) {
					t.Errorf("missing %s", want)
				}
			}
			if current == "v2" && !strings.Contains(page, `value="/docs/v2/" selected>v2</option>`) {
				t.Error("current version is not selected")
			}
			if current == "" && !strings.Contains(page, `value="" disabled selected>Choose a version</option>`) {
				t.Error("missing version placeholder")
			}
			if strings.Contains(page, `/docs/v2/docs/`) {
				t.Error("version URLs were prefixed with base_path")
			}
		})
	}
}

func TestBuildWithoutVersionsHasNoPicker(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Home\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Build(src, dst, Config{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `<select id="moat-version"`) {
		t.Error("unversioned site has a version picker")
	}
}
