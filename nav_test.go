package main

import (
	"strings"
	"testing"
)

func TestDefaultURLPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"01-guide/02-agents.md", "/guide/agents/"},
		{"quickstart.md", "/quickstart/"},
		{"1-intro.md", "/intro/"},
		{"001-advanced.md", "/advanced/"},
		{"guide/index.md", "/guide/"},
		{"index.md", "/"},
		{"no-prefix.md", "/no-prefix/"},
		{"10-config/03-options.md", "/config/options/"},
	}

	for _, tt := range tests {
		got := defaultURLPath(tt.input)
		if got != tt.want {
			t.Errorf("defaultURLPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildNavSections(t *testing.T) {
	pages := []Page{
		{RelPath: "index.md", Frontmatter: Frontmatter{Title: "Home"}},
		{RelPath: "01-guide/01-quickstart.md", Frontmatter: Frontmatter{Title: "Quick Start"}},
		{RelPath: "01-guide/02-config.md", Frontmatter: Frontmatter{Title: "Config"}},
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
	}

	nav := BuildNav(pages)

	// Should have root page "About" and section "Guide"
	if len(nav) != 2 {
		t.Fatalf("expected 2 nav items, got %d", len(nav))
	}

	// First should be root-level "About"
	if nav[0].Title != "About" {
		t.Errorf("first nav item = %q, want About", nav[0].Title)
	}

	// Second should be section "Guide" with 2 children
	if nav[1].Title != "Guide" {
		t.Errorf("section title = %q, want Guide", nav[1].Title)
	}
	if len(nav[1].Children) != 2 {
		t.Errorf("section children = %d, want 2", len(nav[1].Children))
	}
}

func TestBuildNavDateSortedSection(t *testing.T) {
	pages := []Page{
		{RelPath: "index.md", Frontmatter: Frontmatter{Title: "Home"}},
		{RelPath: "posts/older.md", Frontmatter: Frontmatter{Title: "Older Post", Date: "2026-01-01"}},
		{RelPath: "posts/newest.md", Frontmatter: Frontmatter{Title: "Newest Post", Date: "2026-03-18"}},
		{RelPath: "posts/middle.md", Frontmatter: Frontmatter{Title: "Middle Post", Date: "2026-02-15"}},
	}

	nav := BuildNav(pages)

	if len(nav) != 1 {
		t.Fatalf("expected 1 nav section, got %d", len(nav))
	}

	section := nav[0]
	if section.Title != "Posts" {
		t.Errorf("section title = %q, want Posts", section.Title)
	}
	if len(section.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(section.Children))
	}

	// Should be reverse chronological
	if section.Children[0].Title != "Newest Post" {
		t.Errorf("first child = %q, want Newest Post", section.Children[0].Title)
	}
	if section.Children[1].Title != "Middle Post" {
		t.Errorf("second child = %q, want Middle Post", section.Children[1].Title)
	}
	if section.Children[2].Title != "Older Post" {
		t.Errorf("third child = %q, want Older Post", section.Children[2].Title)
	}
}

func TestBuildNavMixedDatedUndatedSection(t *testing.T) {
	pages := []Page{
		{RelPath: "posts/undated.md", Frontmatter: Frontmatter{Title: "Undated Post"}},
		{RelPath: "posts/newer.md", Frontmatter: Frontmatter{Title: "Newer Post", Date: "2026-03-18"}},
		{RelPath: "posts/older.md", Frontmatter: Frontmatter{Title: "Older Post", Date: "2026-01-01"}},
	}

	nav := BuildNav(pages)
	if len(nav) != 1 {
		t.Fatalf("expected 1 section, got %d", len(nav))
	}

	children := nav[0].Children
	if len(children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(children))
	}

	// Dated pages first (newest), then undated (sorted by path)
	if children[0].Title != "Newer Post" {
		t.Errorf("first = %q, want Newer Post", children[0].Title)
	}
	if children[1].Title != "Older Post" {
		t.Errorf("second = %q, want Older Post", children[1].Title)
	}
	if children[2].Title != "Undated Post" {
		t.Errorf("third = %q, want Undated Post", children[2].Title)
	}
}

func TestRenderNavEscapesHTML(t *testing.T) {
	items := []NavItem{
		{Title: `<script>alert("xss")</script>`, Path: "/evil/"},
	}
	links := []LinkConfig{
		{Title: `Bob & "friends"`, URL: `https://example.com/?a=1&b=2`},
	}

	html := RenderNav(items, "/other/", "", links)

	// Should NOT contain raw < or unescaped &
	if strings.Contains(html, "<script>") {
		t.Error("RenderNav did not escape <script> in nav item title")
	}
	if strings.Contains(html, `"friends"`) && !strings.Contains(html, `&#34;friends&#34;`) && !strings.Contains(html, `&quot;friends&quot;`) {
		t.Error("RenderNav did not escape quotes in link title")
	}
	if strings.Contains(html, `?a=1&b=2`) && !strings.Contains(html, `?a=1&amp;b=2`) {
		t.Error("RenderNav did not escape & in URL")
	}
}

func TestRenderNavAriaCurrent(t *testing.T) {
	items := []NavItem{
		{Title: "Home", Path: "/"},
		{Title: "Guide", Children: []NavItem{
			{Title: "Config", Path: "/guide/config/"},
			{Title: "Setup", Path: "/guide/setup/"},
		}},
	}

	html := RenderNav(items, "/guide/config/", "", nil)

	if !strings.Contains(html, `aria-current="page"`) {
		t.Error("expected aria-current on active page")
	}
	// Only one aria-current
	if strings.Count(html, `aria-current="page"`) != 1 {
		t.Errorf("expected exactly 1 aria-current, got %d", strings.Count(html, `aria-current="page"`))
	}
}

func TestBuildNavChildrenFalse(t *testing.T) {
	pages := []Page{
		{RelPath: "index.md", Frontmatter: Frontmatter{Title: "Home"}},
		{RelPath: "01-guide/01-intro.md", Frontmatter: Frontmatter{Title: "Intro"}},
		{RelPath: "01-guide/02-config.md", Frontmatter: Frontmatter{Title: "Config"}},
		{RelPath: "02-changelog/index.md", Frontmatter: Frontmatter{
			Title: "Changelog",
			Extra: map[string]any{"nav_children": false},
		}},
		{RelPath: "02-changelog/v0.1.0.md", Frontmatter: Frontmatter{Title: "v0.1.0", Date: "2026-01-01"}},
		{RelPath: "02-changelog/v0.2.0.md", Frontmatter: Frontmatter{Title: "v0.2.0", Date: "2026-02-01"}},
	}

	nav := BuildNav(pages)

	// Should have Guide (with children) and Changelog (flat link, no children)
	if len(nav) != 2 {
		t.Fatalf("expected 2 nav items, got %d", len(nav))
	}

	guide := nav[0]
	if guide.Title != "Guide" || len(guide.Children) != 2 {
		t.Errorf("guide: title=%q children=%d, want Guide/2", guide.Title, len(guide.Children))
	}

	changelog := nav[1]
	if changelog.Title != "Changelog" {
		t.Errorf("changelog title = %q, want Changelog", changelog.Title)
	}
	if len(changelog.Children) != 0 {
		t.Errorf("changelog should have no children (nav_children: false), got %d", len(changelog.Children))
	}
	if changelog.Path != "/changelog/" {
		t.Errorf("changelog path = %q, want /changelog/", changelog.Path)
	}
}
