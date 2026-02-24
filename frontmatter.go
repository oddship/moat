package main

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter holds YAML metadata from the top of a markdown file.
type Frontmatter struct {
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	URL         string         `yaml:"url"`
	Layout      string         `yaml:"layout"`
	Extra       map[string]any `yaml:"-"` // All other fields
}

// ParseFrontmatter splits a markdown file into frontmatter and body.
// If no frontmatter delimiter (---) is found, returns empty Frontmatter and the full content.
func ParseFrontmatter(content []byte) (Frontmatter, []byte) {
	var fm Frontmatter

	s := string(content)
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "---\r\n") {
		return fm, content
	}

	// Find closing ---
	rest := s[4:] // skip opening "---\n"
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return fm, content
	}

	yamlBlock := rest[:idx]
	body := rest[idx+4:] // skip "\n---"
	// Skip optional newline after closing ---
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	} else if strings.HasPrefix(body, "\r\n") {
		body = body[2:]
	}

	// Parse known fields
	_ = yaml.Unmarshal([]byte(yamlBlock), &fm)

	// Parse all fields into a map for extras
	var raw map[string]any
	_ = yaml.Unmarshal([]byte(yamlBlock), &raw)
	if raw != nil {
		// Remove known fields, keep the rest as Extra
		delete(raw, "title")
		delete(raw, "description")
		delete(raw, "url")
		delete(raw, "layout")
		if len(raw) > 0 {
			fm.Extra = raw
		}
	}

	return fm, []byte(body)
}

// TitleFromFilename derives a title from a filename.
// "02-agents.md" → "Agents", "quickstart.md" → "Quickstart"
func TitleFromFilename(name string) string {
	// Remove .md extension
	name = strings.TrimSuffix(name, ".md")

	// Strip leading number prefix: "02-agents" → "agents"
	if len(name) >= 3 && name[0] >= '0' && name[0] <= '9' && name[1] >= '0' && name[1] <= '9' && name[2] == '-' {
		name = name[3:]
	}

	// Replace hyphens with spaces and title case
	name = strings.ReplaceAll(name, "-", " ")
	return titleCase(name)
}

// TitleFromDir derives a section title from a directory name.
func TitleFromDir(name string) string {
	// Strip leading number prefix: "01-guide" → "guide"
	if len(name) >= 3 && name[0] >= '0' && name[0] <= '9' && name[1] >= '0' && name[1] <= '9' && name[2] == '-' {
		name = name[3:]
	}
	name = strings.ReplaceAll(name, "-", " ")
	return titleCase(name)
}

func titleCase(s string) string {
	var buf bytes.Buffer
	upper := true
	for _, r := range s {
		if upper && r >= 'a' && r <= 'z' {
			buf.WriteRune(r - 32)
			upper = false
		} else {
			buf.WriteRune(r)
			upper = r == ' '
		}
	}
	return buf.String()
}
