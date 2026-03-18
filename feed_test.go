package main

import (
	"testing"
)

func TestBuildFeedIncludesOnlyDatedPagesNewestFirst(t *testing.T) {
	pages := []Page{
		{
			RelPath:     "posts/older.md",
			Frontmatter: Frontmatter{Title: "Older Post", Date: "2026-01-01", Description: "Older"},
			HTML:        []byte("<p>Older content</p>"),
		},
		{
			RelPath:     "about.md",
			Frontmatter: Frontmatter{Title: "About"},
			HTML:        []byte("<p>About page</p>"),
		},
		{
			RelPath:     "posts/newer.md",
			Frontmatter: Frontmatter{Title: "Newer Post", Date: "2026-03-18", Description: "Newer"},
			HTML:        []byte("<p>Newer content</p>"),
		},
	}

	cfg := Config{
		SiteName: "Test Site",
		Feed: FeedConfig{
			Link: "https://example.com",
		},
		Extra: map[string]any{
			"tagline": "A test site",
		},
	}

	feed := buildFeed(pages, cfg)

	if feed.Channel.Title != "Test Site" {
		t.Errorf("channel title = %q, want Test Site", feed.Channel.Title)
	}
	if feed.Channel.Description != "A test site" {
		t.Errorf("channel description = %q, want A test site", feed.Channel.Description)
	}
	if len(feed.Channel.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(feed.Channel.Items))
	}

	if feed.Channel.Items[0].Title != "Newer Post" {
		t.Errorf("first item title = %q, want Newer Post", feed.Channel.Items[0].Title)
	}
	if feed.Channel.Items[1].Title != "Older Post" {
		t.Errorf("second item title = %q, want Older Post", feed.Channel.Items[1].Title)
	}
	if feed.Channel.Items[0].Link != "https://example.com/posts/newer/" {
		t.Errorf("first item link = %q, want https://example.com/posts/newer/", feed.Channel.Items[0].Link)
	}
	if feed.Channel.Items[0].PubDate == "" || feed.Channel.Items[1].PubDate == "" {
		t.Error("expected pubDate for all feed items")
	}
}

func TestBuildFeedSkipsInvalidDates(t *testing.T) {
	pages := []Page{
		{RelPath: "posts/bad.md", Frontmatter: Frontmatter{Title: "Bad", Date: "18-03-2026"}, HTML: []byte("<p>x</p>")},
		{RelPath: "posts/good.md", Frontmatter: Frontmatter{Title: "Good", Date: "2026-03-18"}, HTML: []byte("<p>x</p>")},
	}

	feed := buildFeed(pages, Config{SiteName: "My Site", Feed: FeedConfig{Link: "https://example.com"}})
	if len(feed.Channel.Items) != 1 {
		t.Fatalf("expected 1 feed item, got %d", len(feed.Channel.Items))
	}
	if feed.Channel.Items[0].Title != "Good" {
		t.Errorf("item title = %q, want Good", feed.Channel.Items[0].Title)
	}
}

func TestBuildFeedLinkIncludesBasePath(t *testing.T) {
	// feed.link should be the full site root — basePath is NOT appended again
	pages := []Page{
		{RelPath: "posts/hello.md", Frontmatter: Frontmatter{Title: "Hello", Date: "2026-03-18"}, HTML: []byte("<p>hi</p>")},
	}

	cfg := Config{
		SiteName: "Test",
		BasePath: "/moat",
		Feed:     FeedConfig{Link: "https://oddship.github.io/moat"},
	}

	feed := buildFeed(pages, cfg)
	got := feed.Channel.Items[0].Link
	want := "https://oddship.github.io/moat/posts/hello/"
	if got != want {
		t.Errorf("feed item link = %q, want %q (no double base path)", got, want)
	}
}

func TestBuildFeedEmptySiteNameFallback(t *testing.T) {
	feed := buildFeed(nil, Config{Feed: FeedConfig{Link: "https://example.com"}})
	if feed.Channel.Title != "Site" {
		t.Errorf("expected fallback title 'Site', got %q", feed.Channel.Title)
	}
}

func TestBuildFeedCustomTitle(t *testing.T) {
	cfg := Config{
		SiteName: "My Site",
		Feed: FeedConfig{
			Title: "My Site Feed",
			Link:  "https://example.com",
		},
	}

	feed := buildFeed(nil, cfg)
	if feed.Channel.Title != "My Site Feed" {
		t.Errorf("channel title = %q, want My Site Feed", feed.Channel.Title)
	}
}
