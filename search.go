package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const maxSearchTextLen = 2000

const searchIndexFilename = "_search.json"

// SearchIndex is the static search payload emitted at build time.
type SearchIndex struct {
	Entries []SearchEntry `json:"entries"`
}

// SearchEntry is a single searchable page.
type SearchEntry struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Text        string `json:"text"`
}

// buildSearchIndex creates the search payload from rendered pages.
// Must be called after the render loop so that page.HTML is populated.
func buildSearchIndex(pages []Page, basePath string) SearchIndex {
	basePath = strings.TrimRight(basePath, "/")
	entries := make([]SearchEntry, 0, len(pages))

	for _, page := range pages {
		entries = append(entries, SearchEntry{
			URL:         basePath + pageURL(page),
			Title:       pageTitle(page),
			Description: page.Frontmatter.Description,
			Text:        extractSearchText(page.HTML),
		})
	}

	return SearchIndex{Entries: entries}
}

func writeSearchIndex(dst string, index SearchIndex) error {
	data, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("marshaling search index: %w", err)
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dst, searchIndexFilename)
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func removeSearchIndex(dst string) error {
	path := filepath.Join(dst, searchIndexFilename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing search index: %w", err)
	}
	return nil
}

var reHTMLTag = regexp.MustCompile(`<[^>]+>`)

// extractSearchText strips HTML tags from rendered page content,
// collapses whitespace, and caps length for a compact search index.
func extractSearchText(html []byte) string {
	s := reHTMLTag.ReplaceAllString(string(html), " ")
	s = strings.Join(strings.Fields(s), " ")

	if len(s) > maxSearchTextLen {
		s = s[:maxSearchTextLen]
		if i := strings.LastIndex(s, " "); i > maxSearchTextLen-100 {
			s = s[:i]
		}
	}

	return s
}
