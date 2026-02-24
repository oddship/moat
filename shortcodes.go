package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ShortcodeContext is passed to shortcode templates.
type ShortcodeContext struct {
	Inner template.HTML  // Rendered inner content (markdown → HTML)
	Args  map[string]string // Named arguments
	Page  *TemplateData     // Parent page data
}

// Get returns a named argument or empty string.
func (sc ShortcodeContext) Get(key string) string {
	return sc.Args[key]
}

// shortcodeRegistry holds parsed shortcode templates.
type shortcodeRegistry struct {
	templates map[string]*template.Template
}

// loadShortcodes discovers shortcode templates from _shortcodes/ directory,
// falling back to embedded defaults for any not provided.
func loadShortcodes(src string) (*shortcodeRegistry, error) {
	reg := &shortcodeRegistry{templates: make(map[string]*template.Template)}

	// Load from source directory
	dir := filepath.Join(src, "_shortcodes")
	entries, err := os.ReadDir(dir)
	if err == nil {
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, ".html") {
				continue
			}

			scName := strings.TrimSuffix(name, ".html")
			path := filepath.Join(dir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("reading shortcode %s: %w", name, err)
			}

			tmpl, err := template.New(scName).Parse(string(data))
			if err != nil {
				return nil, fmt.Errorf("parsing shortcode %s: %w", name, err)
			}

			reg.templates[scName] = tmpl
			fmt.Printf("  Shortcode: %s\n", scName)
		}
	}

	return reg, nil
}

// Opening tag: {{< name key="val" >}}
var reShortcodeOpen = regexp.MustCompile(`\{\{<\s*(\w+)((?:\s+\w+="[^"]*")*)\s*>\}\}`)

// Closing tag: {{< /name >}}
var reShortcodeClose = regexp.MustCompile(`\{\{<\s*/(\w+)\s*>\}\}`)

// Self-closing shortcode: {{< name key="val" />}}
var reShortcodeSelf = regexp.MustCompile(`\{\{<\s*(\w+)((?:\s+\w+="[^"]*")*)\s*/>\}\}`)

// Argument parser: key="value"
var reArgs = regexp.MustCompile(`(\w+)="([^"]*)"`)

// ProcessShortcodes replaces shortcode calls in markdown source with rendered HTML.
// Must be called BEFORE markdown rendering.
func (reg *shortcodeRegistry) ProcessShortcodes(source []byte, page *TemplateData) ([]byte, error) {
	if len(reg.templates) == 0 {
		return source, nil
	}

	result := source

	// Process self-closing shortcodes first (simpler, no nesting concerns)
	var lastErr error
	result = reShortcodeSelf.ReplaceAllFunc(result, func(match []byte) []byte {
		if lastErr != nil {
			return match
		}

		parts := reShortcodeSelf.FindSubmatch(match)
		name := string(parts[1])
		argsStr := string(parts[2])

		tmpl, ok := reg.templates[name]
		if !ok {
			lastErr = fmt.Errorf("unknown shortcode: %s", name)
			return match
		}

		ctx := ShortcodeContext{
			Args: parseArgs(argsStr),
			Page: page,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, ctx); err != nil {
			lastErr = fmt.Errorf("executing shortcode %s: %w", name, err)
			return match
		}

		return buf.Bytes()
	})

	if lastErr != nil {
		return nil, lastErr
	}

	// Process block shortcodes by finding matched open/close pairs
	result, err := reg.processBlockShortcodes(result, page)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// processBlockShortcodes finds matched {{< name >}}...{{< /name >}} pairs
// and replaces them with rendered shortcode output.
func (reg *shortcodeRegistry) processBlockShortcodes(source []byte, page *TemplateData) ([]byte, error) {
	s := string(source)

	for {
		// Find the first opening tag
		openLoc := reShortcodeOpen.FindStringIndex(s)
		if openLoc == nil {
			break
		}

		openMatch := reShortcodeOpen.FindStringSubmatch(s[openLoc[0]:])
		name := openMatch[1]
		argsStr := openMatch[2]

		// Find the matching closing tag for this name
		closePattern := regexp.MustCompile(`\{\{<\s*/` + regexp.QuoteMeta(name) + `\s*>\}\}`)
		rest := s[openLoc[1]:]
		closeLoc := closePattern.FindStringIndex(rest)
		if closeLoc == nil {
			return nil, fmt.Errorf("shortcode %q opened but never closed", name)
		}

		inner := rest[:closeLoc[0]]
		fullEnd := openLoc[1] + closeLoc[1]

		tmpl, ok := reg.templates[name]
		if !ok {
			return nil, fmt.Errorf("unknown shortcode: %s", name)
		}

		// Render inner content as markdown
		innerHTML, err := RenderMarkdown([]byte(inner))
		if err != nil {
			return nil, fmt.Errorf("rendering inner content for shortcode %s: %w", name, err)
		}

		ctx := ShortcodeContext{
			Inner: template.HTML(innerHTML),
			Args:  parseArgs(argsStr),
			Page:  page,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, ctx); err != nil {
			return nil, fmt.Errorf("executing shortcode %s: %w", name, err)
		}

		// Replace the full shortcode call with rendered output
		s = s[:openLoc[0]] + buf.String() + s[fullEnd:]
	}

	return []byte(s), nil
}

func parseArgs(s string) map[string]string {
	args := make(map[string]string)
	matches := reArgs.FindAllStringSubmatch(s, -1)
	for _, m := range matches {
		args[m[1]] = m[2]
	}
	return args
}
