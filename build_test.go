package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildPageMetaSortOrder(t *testing.T) {
	pages := []Page{
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
		{RelPath: "posts/older.md", Frontmatter: Frontmatter{Title: "Older", Date: "2026-01-01"}},
		{RelPath: "posts/newer.md", Frontmatter: Frontmatter{Title: "Newer", Date: "2026-03-18"}},
		{RelPath: "zebra.md", Frontmatter: Frontmatter{Title: "Zebra"}},
	}

	metas := buildPageMeta(pages, "")

	// Expected order: dated desc (Newer, Older), then undated alpha (About, Zebra)
	expected := []string{"Newer", "Older", "About", "Zebra"}
	if len(metas) != len(expected) {
		t.Fatalf("expected %d metas, got %d", len(expected), len(metas))
	}
	for i, want := range expected {
		if metas[i].Title != want {
			t.Errorf("metas[%d].Title = %q, want %q", i, metas[i].Title, want)
		}
	}
}

func TestBuildPageMetaSection(t *testing.T) {
	pages := []Page{
		{RelPath: "index.md", Frontmatter: Frontmatter{Title: "Home"}},
		{RelPath: "01-guide/01-intro.md", Frontmatter: Frontmatter{Title: "Intro"}},
		{RelPath: "about.md", Frontmatter: Frontmatter{Title: "About"}},
	}

	metas := buildPageMeta(pages, "/site")

	// Root index.md should be excluded (matches nav behavior)
	if len(metas) != 2 {
		t.Fatalf("expected 2 metas (no root index), got %d", len(metas))
	}

	sectionMap := map[string]string{}
	for _, m := range metas {
		sectionMap[m.Title] = m.Section
	}

	if _, ok := sectionMap["Home"]; ok {
		t.Error("root index.md should not appear in Pages")
	}
	if sectionMap["Intro"] != "guide" {
		t.Errorf("Intro section = %q, want guide", sectionMap["Intro"])
	}
	if sectionMap["About"] != "" {
		t.Errorf("About section = %q, want empty", sectionMap["About"])
	}

	// Check basePath applied
	for _, m := range metas {
		if m.Title == "Intro" && m.URL != "/site/guide/intro/" {
			t.Errorf("Intro URL = %q, want /site/guide/intro/", m.URL)
		}
	}
}

func TestBuildPageMetaSortsMixedDateFormatsChronologically(t *testing.T) {
	metas := buildPageMeta([]Page{
		{RelPath: "morning.md", Frontmatter: Frontmatter{Title: "Morning", Date: "2026-03-18T09:00"}},
		{RelPath: "evening.md", Frontmatter: Frontmatter{Title: "Evening", Date: "2026-03-18 21:00"}},
	}, "")

	if metas[0].Title != "Evening" || metas[1].Title != "Morning" {
		t.Fatalf("expected chronological order Evening, Morning; got %q, %q", metas[0].Title, metas[1].Title)
	}
}

func TestOutputPathRejectsTraversal(t *testing.T) {
	if _, err := outputPathFromURL("/tmp/site", "/../outside/"); err == nil {
		t.Fatal("expected traversal URL to be rejected")
	}
}

func TestBuildCleansDestination(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	stale := filepath.Join(dst, "old", "index.html")
	if err := os.MkdirAll(filepath.Dir(stale), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Home\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Build(src, dst, Config{}); err != nil {
		t.Fatalf("Build: %v", err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("expected stale output to be removed, got err=%v", err)
	}
}

func TestBuildBuiltinLayoutUsesRootRelativeURLs(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Home\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	enabled := true
	cfg := Config{
		Favicon: "_static/favicon.svg",
		Logo:    "_static/logo.svg",
		Feed:    FeedConfig{Enabled: &enabled},
	}

	if err := Build(src, dst, cfg); err != nil {
		t.Fatalf("Build: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(data)
	for _, want := range []string{
		`href="/"`,
		`href="/_static/favicon.svg"`,
		`href="/_syntax.css"`,
		`href="/feed.xml"`,
		`src="/_static/logo.svg"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("built page does not contain %s", want)
		}
	}
	if strings.Contains(html, `href="//`) || strings.Contains(html, `src="//`) {
		t.Error("built page contains protocol-relative asset URLs")
	}
}
