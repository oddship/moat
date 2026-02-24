package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// loadLayouts discovers and parses layout templates from the source directory,
// falling back to embedded defaults if no _layout.html is found.
//
// _layout.html is the base layout. It should use {{ block "name" . }}
// to define overridable sections.
//
// _layout.{name}.html are variants that override blocks from the base.
// They contain {{ define "blockname" }}...{{ end }} to replace base blocks.
//
// Returns a map: "" → base template, "name" → variant template.
func loadLayouts(src string) (map[string]*template.Template, error) {
	layouts := make(map[string]*template.Template)

	// Try to read base layout from source directory
	basePath := filepath.Join(src, "_layout.html")
	baseBytes, err := os.ReadFile(basePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading _layout.html: %w", err)
		}
		// Fall back to embedded default
		baseBytes, err = readEmbeddedFile("_layout.html")
		if err != nil {
			return nil, fmt.Errorf("no _layout.html found and embedded default missing: %w", err)
		}
		fmt.Printf("  Using built-in layout (create _layout.html to customize)\n")
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}

	baseTmpl, err := template.New("layout").Funcs(funcMap).Parse(string(baseBytes))
	if err != nil {
		return nil, fmt.Errorf("parsing _layout.html: %w", err)
	}
	layouts[""] = baseTmpl

	// Discover named variants from source directory
	if err := discoverVariants(src, baseTmpl, layouts); err != nil {
		return nil, err
	}

	// Also load embedded variants that aren't overridden by source
	if err := discoverEmbeddedVariants(baseTmpl, layouts); err != nil {
		return nil, err
	}

	return layouts, nil
}

// discoverVariants finds _layout.{name}.html files in the source directory.
func discoverVariants(src string, baseTmpl *template.Template, layouts map[string]*template.Template) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("reading source directory: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "_layout.") || !strings.HasSuffix(name, ".html") || name == "_layout.html" {
			continue
		}

		variant := strings.TrimPrefix(name, "_layout.")
		variant = strings.TrimSuffix(variant, ".html")

		variantBytes, err := os.ReadFile(filepath.Join(src, name))
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}

		cloned, err := baseTmpl.Clone()
		if err != nil {
			return fmt.Errorf("cloning base for %s: %w", name, err)
		}
		if _, err := cloned.Parse(string(variantBytes)); err != nil {
			return fmt.Errorf("parsing %s: %w", name, err)
		}

		layouts[variant] = cloned
		fmt.Printf("  Layout: %s → %s\n", name, variant)
	}

	return nil
}

// discoverEmbeddedVariants loads embedded layout variants not already in layouts.
func discoverEmbeddedVariants(baseTmpl *template.Template, layouts map[string]*template.Template) error {
	entries, err := listEmbeddedDir("")
	if err != nil {
		return nil // No embedded files
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "_layout.") || !strings.HasSuffix(name, ".html") || name == "_layout.html" {
			continue
		}

		variant := strings.TrimPrefix(name, "_layout.")
		variant = strings.TrimSuffix(variant, ".html")

		// Skip if already loaded from source
		if _, exists := layouts[variant]; exists {
			continue
		}

		variantBytes, err := readEmbeddedFile(name)
		if err != nil {
			continue
		}

		cloned, err := baseTmpl.Clone()
		if err != nil {
			continue
		}
		if _, err := cloned.Parse(string(variantBytes)); err != nil {
			continue
		}

		layouts[variant] = cloned
		fmt.Printf("  Layout: %s → %s (built-in)\n", name, variant)
	}

	return nil
}
