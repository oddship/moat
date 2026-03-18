package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBlockShortcodeInnerContentResolvesWikilinks(t *testing.T) {
	dir := t.TempDir()
	scDir := filepath.Join(dir, "_shortcodes")
	if err := os.MkdirAll(scDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scDir, "note.html"), []byte(`<div class="note">{{ .Inner }}</div>`), 0o644); err != nil {
		t.Fatal(err)
	}

	reg, err := loadShortcodes(dir)
	if err != nil {
		t.Fatal(err)
	}

	pages := []Page{
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
		{RelPath: "index.md", Frontmatter: Frontmatter{Title: "Home"}},
	}
	resolver := newPageResolver(pages, "")
	page := &TemplateData{Title: "Home", CurrentPath: "/", Pages: buildPageMeta(pages, "")}

	source := []byte("{{< note >}}See [[About]].{{< /note >}}")
	out, err := reg.ProcessShortcodes(source, page, resolver)
	if err != nil {
		t.Fatal(err)
	}

	html, err := RenderMarkdownWithResolver(out, resolver)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(html), `href="/about/"`) {
		t.Fatalf("expected wikilink inside shortcode content to resolve, got: %s", html)
	}
}

func TestSectionPagesFiltering(t *testing.T) {
	pages := buildPageMeta([]Page{
		{RelPath: "index.md", Frontmatter: Frontmatter{Title: "Home"}},
		{RelPath: "01-guide/01-intro.md", Frontmatter: Frontmatter{Title: "Intro"}},
		{RelPath: "01-guide/02-config.md", Frontmatter: Frontmatter{Title: "Config"}},
		{RelPath: "posts/hello.md", Frontmatter: Frontmatter{Title: "Hello", Date: "2026-03-18"}},
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
	}, "")

	sc := ShortcodeContext{
		Page: &TemplateData{
			CurrentPath: "/guide/intro/",
			Pages:       pages,
		},
	}

	// Filter by section
	guide := sc.SectionPages("guide")
	if len(guide) != 1 {
		t.Fatalf("expected 1 guide page (excluding current), got %d", len(guide))
	}
	if guide[0].Title != "Config" {
		t.Errorf("expected Config, got %s", guide[0].Title)
	}

	// Empty section = all pages except current (root index.md excluded by buildPageMeta)
	all := sc.SectionPages("")
	if len(all) != 3 {
		t.Fatalf("expected 3 pages (all minus current, no root index), got %d", len(all))
	}

	// Nonexistent section
	none := sc.SectionPages("nonexistent")
	if len(none) != 0 {
		t.Fatalf("expected 0 pages for nonexistent section, got %d", len(none))
	}

	// Nil page
	nilSC := ShortcodeContext{Page: nil}
	if nilSC.SectionPages("guide") != nil {
		t.Error("expected nil for nil page")
	}
}
