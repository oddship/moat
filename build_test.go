package main

import "testing"

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
