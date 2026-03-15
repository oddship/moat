package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func boolPtr(v bool) *bool {
	return &v
}

func TestSearchEnabledDefaultsToTrue(t *testing.T) {
	var cfg Config
	if !cfg.SearchEnabled() {
		t.Fatal("expected search to default to enabled")
	}
}

func TestSearchEnabledCanBeDisabled(t *testing.T) {
	cfg := Config{Search: SearchConfig{Enabled: boolPtr(false)}}
	if cfg.SearchEnabled() {
		t.Fatal("expected search to be disabled")
	}
}

func TestBuildSearchIndexUsesBasePathAndTitleFallback(t *testing.T) {
	// page.HTML is what buildSearchIndex reads — simulate rendered output
	index := buildSearchIndex([]Page{{
		RelPath: "01-guide/02-config.md",
		Frontmatter: Frontmatter{
			Description: "Site-level configuration",
		},
		HTML: []byte("<h1>Config</h1>\n<p>Hello world</p>\n"),
	}}, "/docs")

	if len(index.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(index.Entries))
	}

	entry := index.Entries[0]
	if entry.URL != "/docs/guide/config/" {
		t.Fatalf("expected URL %q, got %q", "/docs/guide/config/", entry.URL)
	}
	if entry.Title != "Config" {
		t.Fatalf("expected fallback title %q, got %q", "Config", entry.Title)
	}
	if entry.Description != "Site-level configuration" {
		t.Fatalf("expected description %q, got %q", "Site-level configuration", entry.Description)
	}
	if entry.Text != "Config Hello world" {
		t.Fatalf("expected text %q, got %q", "Config Hello world", entry.Text)
	}
}

func TestBuildSearchIndexUsesFrontmatterTitleWhenPresent(t *testing.T) {
	index := buildSearchIndex([]Page{{
		RelPath: "01-guide/02-config.md",
		Frontmatter: Frontmatter{
			Title: "My Custom Title",
		},
		HTML: []byte("<p>Hello world</p>"),
	}}, "")

	if len(index.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(index.Entries))
	}
	if index.Entries[0].Title != "My Custom Title" {
		t.Fatalf("expected custom title %q, got %q", "My Custom Title", index.Entries[0].Title)
	}
}

func TestExtractSearchTextStripsHTMLTags(t *testing.T) {
	html := []byte(`<h1>Title</h1>
<p>Some <strong>bold</strong> and <em>italic</em> text.</p>
<a href="https://example.com">a link</a>
<pre><code class="language-go">fmt.Println("hello")</code></pre>
<div role="alert" data-variant="info"><p>A shortcode note</p></div>`)

	got := extractSearchText(html)

	// Should not contain any HTML tags
	if strings.Contains(got, "<") || strings.Contains(got, ">") {
		t.Errorf("expected HTML tags stripped, got %q", got)
	}

	// Should preserve visible text content
	for _, want := range []string{"Title", "bold", "italic", "text", "a link", "hello", "A shortcode note"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q to be preserved, got %q", want, got)
		}
	}
}

func TestExtractSearchTextCapsLength(t *testing.T) {
	// Build HTML content longer than maxSearchTextLen
	var b strings.Builder
	for i := 0; i < 500; i++ {
		b.WriteString("<p>word</p> ")
	}

	got := extractSearchText([]byte(b.String()))
	if len(got) > maxSearchTextLen {
		t.Fatalf("expected text capped at %d, got %d chars", maxSearchTextLen, len(got))
	}
	if got[len(got)-1] == ' ' {
		t.Fatal("trailing space after truncation")
	}
}

func TestWriteSearchIndexWithNoPagesProducesEmptyEntriesArray(t *testing.T) {
	dst := t.TempDir()

	if err := writeSearchIndex(dst, buildSearchIndex(nil, "")); err != nil {
		t.Fatalf("writeSearchIndex: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, searchIndexFilename))
	if err != nil {
		t.Fatalf("reading search index: %v", err)
	}

	if string(data) != "{\"entries\":[]}\n" {
		t.Fatalf("expected empty entries JSON, got %q", string(data))
	}
}

func TestBuildWritesSearchIndexByDefault(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	content := "---\ntitle: Home\ndescription: Welcome page\n---\n\n# Hello\n\nSearch me\n"
	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("writing source page: %v", err)
	}

	if err := Build(src, dst, Config{SiteName: "Site"}); err != nil {
		t.Fatalf("Build: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, searchIndexFilename))
	if err != nil {
		t.Fatalf("reading search index: %v", err)
	}

	var index SearchIndex
	if err := json.Unmarshal(data, &index); err != nil {
		t.Fatalf("unmarshal search index: %v", err)
	}

	if len(index.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(index.Entries))
	}

	entry := index.Entries[0]
	if entry.URL != "/" {
		t.Fatalf("expected URL %q, got %q", "/", entry.URL)
	}
	if entry.Title != "Home" {
		t.Fatalf("expected title %q, got %q", "Home", entry.Title)
	}
	if entry.Description != "Welcome page" {
		t.Fatalf("expected description %q, got %q", "Welcome page", entry.Description)
	}
	// Text comes from rendered HTML — should contain words without tags or markdown
	if !strings.Contains(entry.Text, "Hello") || !strings.Contains(entry.Text, "Search me") {
		t.Fatalf("expected rendered text content, got %q", entry.Text)
	}
	if strings.Contains(entry.Text, "<") || strings.Contains(entry.Text, "#") {
		t.Fatalf("expected no HTML tags or markdown markers, got %q", entry.Text)
	}
}

func TestBuildSkipsSearchIndexWhenDisabled(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Home\n"), 0o644); err != nil {
		t.Fatalf("writing source page: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dst, searchIndexFilename), []byte("stale"), 0o644); err != nil {
		t.Fatalf("writing stale search index: %v", err)
	}

	cfg := Config{SiteName: "Site", Search: SearchConfig{Enabled: boolPtr(false)}}
	if err := Build(src, dst, cfg); err != nil {
		t.Fatalf("Build: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, searchIndexFilename)); !os.IsNotExist(err) {
		t.Fatalf("expected no search index file, got err=%v", err)
	}
}
