package main

import (
	"strings"
	"testing"
)

func TestRenderMarkdownBasic(t *testing.T) {
	out, err := RenderMarkdown([]byte("**bold**"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "<strong>bold</strong>") {
		t.Errorf("expected <strong>, got %s", out)
	}
}

func TestRenderMarkdownWikilink(t *testing.T) {
	pages := []Page{
		{RelPath: "01-guide/01-getting-started.md", Frontmatter: Frontmatter{Title: "Getting Started"}},
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
	}
	resolver := newPageResolver(pages, "")

	// Known page
	out, err := RenderMarkdownWithResolver([]byte("See [[Getting Started]] for more."), resolver)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, `href="/guide/getting-started/"`) {
		t.Errorf("expected link to /guide/getting-started/, got: %s", html)
	}
	if !strings.Contains(html, ">Getting Started</a>") {
		t.Errorf("expected link text 'Getting Started', got: %s", html)
	}

	// Case insensitive
	out, err = RenderMarkdownWithResolver([]byte("See [[about]]"), resolver)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `href="/about/"`) {
		t.Errorf("expected case-insensitive match, got: %s", out)
	}

	// Unknown page — should render as plain text (no link)
	out, err = RenderMarkdownWithResolver([]byte("See [[Nonexistent Page]]"), resolver)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "<a") {
		t.Errorf("unknown page should not create a link, got: %s", out)
	}
}

func TestRenderMarkdownWikilinkWithFragment(t *testing.T) {
	pages := []Page{
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
	}
	resolver := newPageResolver(pages, "")

	out, err := RenderMarkdownWithResolver([]byte("[[About#team]]"), resolver)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `href="/about/#team"`) {
		t.Errorf("expected fragment in link, got: %s", out)
	}
}

func TestWikilinkDuplicateTitleFirstWins(t *testing.T) {
	pages := []Page{
		{RelPath: "01-guide/about.md", Frontmatter: Frontmatter{Title: "About"}},
		{RelPath: "02-reference/about.md", Frontmatter: Frontmatter{Title: "About"}},
	}
	resolver := newPageResolver(pages, "")

	out, err := RenderMarkdownWithResolver([]byte("[[About]]"), resolver)
	if err != nil {
		t.Fatal(err)
	}
	// First page wins
	if !strings.Contains(string(out), `href="/guide/about/"`) {
		t.Errorf("expected first page to win, got: %s", out)
	}
}

func TestRenderMarkdownWikilinkWithBasePath(t *testing.T) {
	pages := []Page{
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
	}
	resolver := newPageResolver(pages, "/mysite")

	out, err := RenderMarkdownWithResolver([]byte("[[About]]"), resolver)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `href="/mysite/about/"`) {
		t.Errorf("expected basePath in link, got: %s", out)
	}
}
