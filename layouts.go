package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// loadLayouts discovers and parses layout templates from the source directory.
//
// _layout.html is the base layout (required). It should use {{ block "name" . }}
// to define overridable sections.
//
// _layout.{name}.html are variants that override blocks from the base.
// They contain {{ define "blockname" }}...{{ end }} to replace base blocks.
//
// Returns a map: "" → base template, "name" → variant template.
func loadLayouts(src string) (map[string]*template.Template, error) {
	layouts := make(map[string]*template.Template)

	// Read base layout (required)
	basePath := filepath.Join(src, "_layout.html")
	baseBytes, err := os.ReadFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("reading _layout.html: %w (expected at %s)", err, basePath)
	}

	baseTmpl, err := template.New("layout").Parse(string(baseBytes))
	if err != nil {
		return nil, fmt.Errorf("parsing _layout.html: %w", err)
	}
	layouts[""] = baseTmpl

	// Discover named variants: _layout.{name}.html
	entries, err := os.ReadDir(src)
	if err != nil {
		return nil, fmt.Errorf("reading source directory: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "_layout.") || !strings.HasSuffix(name, ".html") {
			continue
		}
		if name == "_layout.html" {
			continue // Already handled as base
		}

		// Extract variant name: _layout.wide.html → "wide"
		variant := strings.TrimPrefix(name, "_layout.")
		variant = strings.TrimSuffix(variant, ".html")

		variantPath := filepath.Join(src, name)
		variantBytes, err := os.ReadFile(variantPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", name, err)
		}

		// Clone base and parse variant into it.
		// The variant's {{ define "block" }} overrides the base's {{ block "block" }} defaults.
		cloned, err := baseTmpl.Clone()
		if err != nil {
			return nil, fmt.Errorf("cloning base for %s: %w", name, err)
		}
		if _, err := cloned.Parse(string(variantBytes)); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, err)
		}

		layouts[variant] = cloned
		fmt.Printf("  Layout: %s → %s\n", name, variant)
	}

	return layouts, nil
}
