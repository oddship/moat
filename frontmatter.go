package main

import (
	"bytes"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// dateFormats are the accepted date formats in frontmatter, tried in order.
// All are ISO-style so string comparison works for sorting.
var dateFormats = []string{
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
	"2006-01-02",
}

// ParseDate tries to parse a frontmatter date string.
// Returns the parsed time and true on success, or zero time and false.
func ParseDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, fmt := range dateFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// Frontmatter holds YAML metadata from the top of a markdown file.
type Frontmatter struct {
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	URL         string         `yaml:"url"`
	Layout      string         `yaml:"layout"`
	Date        string         `yaml:"date"`
	Draft       bool           `yaml:"draft"`
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
		delete(raw, "date")
		delete(raw, "draft")
		if len(raw) > 0 {
			fm.Extra = raw
		}
	}

	return fm, []byte(body)
}

// reNumPrefixFM matches leading digits followed by a hyphen: "01-", "1-", "001-"
var reNumPrefixFM = regexp.MustCompile(`^\d+-`)

// TitleFromFilename derives a title from a filename.
// "02-agents.md" → "Agents", "quickstart.md" → "Quickstart", "1-intro.md" → "Intro"
func TitleFromFilename(name string) string {
	// Remove .md extension
	name = strings.TrimSuffix(name, ".md")

	// Strip leading number prefix
	name = reNumPrefixFM.ReplaceAllString(name, "")

	// Replace hyphens with spaces and title case
	name = strings.ReplaceAll(name, "-", " ")
	return titleCase(name)
}

// TitleFromDir derives a section title from a directory name.
func TitleFromDir(name string) string {
	// Strip leading number prefix
	name = reNumPrefixFM.ReplaceAllString(name, "")
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
