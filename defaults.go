package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:embed
var embeddedFS embed.FS

// Init scaffolds a docs directory with starter content.
// Layouts come from embedded defaults (same ones used at build time).
// Also writes a sample config.toml and starter markdown files.
func Init(dir string) error {
	dir, _ = filepath.Abs(dir)

	if _, err := os.Stat(filepath.Join(dir, "_layout.html")); err == nil {
		return fmt.Errorf("%s already has a _layout.html — use moat init on a new directory", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "index.md")); err == nil {
		return fmt.Errorf("%s already has an index.md — use moat init on a new directory", dir)
	}

	// Write embedded layouts + config
	count := 0
	err := fs.WalkDir(embeddedFS, "embed", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel("embed", path)
		out := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		data, err := embeddedFS.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Printf("  %s\n", rel)
		count++
		return os.WriteFile(out, data, 0o644)
	})
	if err != nil {
		return fmt.Errorf("writing files: %w", err)
	}

	// Write starter content
	starters := map[string]string{
		"index.md":                       "---\ntitle: Home\nlayout: landing\n---\n\n# My Project\n\n<h3 class=\"tagline text-light\">Your project description here.</h3>\n\n<div class=\"hstack\">\n  <a href=\"guide/getting-started/\" class=\"button\">Get started</a>\n</div>\n",
		"01-guide/01-getting-started.md": "---\ntitle: Getting Started\n---\n\n# Getting started\n\nWelcome to your docs. Edit this file at `docs/01-guide/01-getting-started.md`.\n",
	}
	for rel, content := range starters {
		out := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(out, []byte(content), 0o644); err != nil {
			return err
		}
		fmt.Printf("  %s\n", rel)
		count++
	}

	fmt.Printf("Initialized %d files in %s\n", count, dir)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  moat build %s/ _site/\n", filepath.Base(dir))
	fmt.Printf("  moat serve _site/\n")
	return nil
}

// readEmbeddedFile returns an embedded file's content.
func readEmbeddedFile(path string) ([]byte, error) {
	return embeddedFS.ReadFile("embed/" + path)
}

// listEmbeddedDir returns entries in an embedded directory.
func listEmbeddedDir(dir string) ([]fs.DirEntry, error) {
	path := "embed"
	if dir != "" {
		path = "embed/" + dir
	}
	return embeddedFS.ReadDir(path)
}
