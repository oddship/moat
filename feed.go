package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const feedFilename = "feed.xml"

// FeedConfig controls RSS feed generation.
type FeedConfig struct {
	Enabled *bool  `toml:"enabled"`
	Title   string `toml:"title"`
	Link    string `toml:"link"`
}

// FeedEnabled returns the effective feed setting.
// Feed defaults to disabled when omitted from config.toml.
func (c Config) FeedEnabled() bool {
	if c.Feed.Enabled == nil {
		return false
	}
	return *c.Feed.Enabled
}

// RSS 2.0 structures
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description,omitempty"`
	BuildDate   string    `xml:"lastBuildDate"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description,omitempty"`
	PubDate     string `xml:"pubDate,omitempty"`
	GUID        string `xml:"guid"`
}

// buildFeed creates an RSS 2.0 feed from rendered pages.
// Only pages with a valid YYYY-MM-DD date are included, sorted newest first.
func buildFeed(pages []Page, cfg Config) rssFeed {
	basePath := strings.TrimRight(cfg.BasePath, "/")
	siteLink := cfg.Feed.Link
	if siteLink == "" {
		siteLink = "/"
	}
	siteLink = strings.TrimRight(siteLink, "/")

	feedTitle := cfg.Feed.Title
	if feedTitle == "" {
		feedTitle = cfg.SiteName
	}
	if feedTitle == "" {
		feedTitle = "Site"
	}

	type feedPage struct {
		page Page
		date time.Time
	}
	var datedPages []feedPage
	for _, page := range pages {
		if page.Frontmatter.Date == "" {
			continue
		}
		t, err := time.Parse("2006-01-02", page.Frontmatter.Date)
		if err != nil {
			continue
		}
		datedPages = append(datedPages, feedPage{page: page, date: t})
	}

	sort.Slice(datedPages, func(i, j int) bool {
		if !datedPages[i].date.Equal(datedPages[j].date) {
			return datedPages[i].date.After(datedPages[j].date)
		}
		return datedPages[i].page.RelPath < datedPages[j].page.RelPath
	})

	items := make([]rssItem, 0, len(datedPages))
	for _, fp := range datedPages {
		page := fp.page
		url := pageURL(page)
		fullURL := siteLink + basePath + url

		title := pageTitle(page)
		desc := page.Frontmatter.Description
		if desc == "" {
			desc = extractSearchText(page.HTML)
			runes := []rune(desc)
			if len(runes) > 300 {
				runes = runes[:300]
				s := string(runes)
				if i := strings.LastIndex(s, " "); i > 200 {
					s = s[:i]
				}
				desc = s + "…"
			}
		}

		items = append(items, rssItem{
			Title:       title,
			Link:        fullURL,
			Description: desc,
			PubDate:     fp.date.Format(time.RFC1123Z),
			GUID:        fullURL,
		})
	}

	channelDesc := ""
	if tagline, ok := cfg.Extra["tagline"].(string); ok {
		channelDesc = tagline
	}
	if channelDesc == "" {
		channelDesc = feedTitle
	}

	return rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       feedTitle,
			Link:        siteLink,
			Description: channelDesc,
			BuildDate:   time.Now().Format(time.RFC1123Z),
			Items:       items,
		},
	}
}

func removeFeed(dst string) error {
	path := filepath.Join(dst, feedFilename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing feed: %w", err)
	}
	return nil
}

func writeFeed(dst string, feed rssFeed) error {
	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling feed: %w", err)
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	header := []byte(xml.Header)
	path := filepath.Join(dst, feedFilename)
	return os.WriteFile(path, append(header, data...), 0o644)
}
