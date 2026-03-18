package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

// Page represents a single markdown page.
type Page struct {
	RelPath     string      // Relative path from source root, e.g. "guide/02-agents.md"
	Frontmatter Frontmatter // Parsed YAML frontmatter
	Body        []byte      // Markdown body (without frontmatter)
	HTML        []byte      // Rendered HTML (set after shortcode + markdown processing)
}

// PageMeta is a lightweight page summary available to templates and shortcodes.
type PageMeta struct {
	Title       string
	Description string
	URL         string
	Date        string
	Extra       map[string]any
	Section     string // Top-level directory, e.g. "guide" (empty for root pages)
}

// TemplateData is passed to the layout template.
type TemplateData struct {
	Title         string
	Description   string
	Date          string         // Page date from frontmatter (raw string, e.g. "2026-03-18 14:30")
	Content       template.HTML
	Nav           template.HTML
	CurrentPath   string
	SiteName      string
	BasePath      string
	Logo          string         // Path to logo image (relative to BasePath)
	LogoInline    template.HTML  // Inlined SVG content (set when logo is .svg)
	Favicon       string         // Path to favicon (relative to BasePath)
	SearchEnabled bool           // Whether built-in search UI should render
	FeedEnabled   bool           // Whether RSS feed is enabled
	TopNav        []LinkConfig   // Top navigation links
	Extra         map[string]any // Per-page extra frontmatter
	Site          map[string]any // Site-level extra from config.toml [extra]
	Pages         []PageMeta     // All non-draft pages (sorted by date desc, then title)
}

// Build reads markdown from src, renders HTML, and writes to dst.
func Build(src, dst string, cfg Config) error {
	src, _ = filepath.Abs(src)
	dst, _ = filepath.Abs(dst)

	basePath := strings.TrimRight(cfg.BasePath, "/")
	siteName := cfg.SiteName
	searchEnabled := cfg.SearchEnabled()
	if siteName == "" {
		siteName = "Site"
	}

	// Load layout templates (base + named variants)
	layouts, err := loadLayouts(src)
	if err != nil {
		return err
	}

	// Load shortcode templates
	shortcodes, err := loadShortcodes(src)
	if err != nil {
		return err
	}

	// Discover and parse markdown files
	var pages []Page
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() || !strings.HasSuffix(name, ".md") {
			return nil
		}

		relPath, _ := filepath.Rel(src, path)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", relPath, err)
		}

		fm, body := ParseFrontmatter(content)

		// Skip draft pages
		if fm.Draft {
			fmt.Printf("  Skipping draft: %s\n", relPath)
			return nil
		}

		// Warn on malformed dates
		if fm.Date != "" {
			if _, ok := ParseDate(fm.Date); !ok {
				fmt.Printf("  Warning: %s has invalid date %q (expected YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)\n", relPath, fm.Date)
			}
		}

		// Body stored raw — shortcodes processed per-page during render
		pages = append(pages, Page{
			RelPath:     relPath,
			Frontmatter: fm,
			Body:        body,
		})
		return nil
	})
	if err != nil {
		return fmt.Errorf("walking source: %w", err)
	}

	fmt.Printf("Found %d pages\n", len(pages))

	// Build navigation
	nav := BuildNav(pages)

	// Build wikilink resolver from discovered pages
	wikiResolver := newPageResolver(pages, basePath)

	// Build page metadata list for templates and shortcodes
	allPages := buildPageMeta(pages, basePath)

	// Generate syntax highlighting CSS
	if err := writeSyntaxCSS(dst, cfg.Highlight); err != nil {
		return fmt.Errorf("writing syntax CSS: %w", err)
	}

	// Read inline SVG logo if configured
	var logoInline template.HTML
	if cfg.Logo != "" && strings.HasSuffix(strings.ToLower(cfg.Logo), ".svg") {
		logoPath := filepath.Join(src, cfg.Logo)
		if data, err := os.ReadFile(logoPath); err == nil {
			logoInline = template.HTML(data)
		}
	}

	// Render each page
	for i, page := range pages {
		currentPath := pageURL(page)
		prefixedPath := basePath + currentPath
		outPath := outputPathFromURL(dst, currentPath)

		navHTML := RenderNav(nav, prefixedPath, basePath, cfg.Links)

		title := page.Frontmatter.Title
		if title == "" {
			title = TitleFromFilename(filepath.Base(page.RelPath))
		}

		data := TemplateData{
			Title:         title,
			Description:   page.Frontmatter.Description,
			Date:          page.Frontmatter.Date,
			Nav:           template.HTML(navHTML),
			CurrentPath:   prefixedPath,
			SiteName:      siteName,
			BasePath:      basePath,
			Logo:          cfg.Logo,
			LogoInline:    logoInline,
			Favicon:       cfg.Favicon,
			SearchEnabled: searchEnabled,
			FeedEnabled:   cfg.FeedEnabled(),
			TopNav:        cfg.TopNav,
			Extra:         page.Frontmatter.Extra,
			Site:          cfg.Extra,
			Pages:         allPages,
		}

		// Process shortcodes in markdown source (before markdown rendering)
		body, err := shortcodes.ProcessShortcodes(page.Body, &data, wikiResolver)
		if err != nil {
			return fmt.Errorf("processing shortcodes in %s: %w", page.RelPath, err)
		}

		// Render markdown to HTML (with wiki link resolution)
		html, err := RenderMarkdownWithResolver(body, wikiResolver)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", page.RelPath, err)
		}
		pages[i].HTML = html
		data.Content = template.HTML(html)

		// Pick layout: frontmatter "layout: name" → _layout.name.html, default → _layout.html
		layoutName := page.Frontmatter.Layout
		tmpl, ok := layouts[layoutName]
		if !ok {
			if layoutName == "" {
				return fmt.Errorf("missing default layout _layout.html")
			}
			return fmt.Errorf("page %s requests layout %q but _layout.%s.html not found", page.RelPath, layoutName, layoutName)
		}

		if err := renderToFile(tmpl, data, outPath); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
		fmt.Printf("  %s → %s\n", page.RelPath, outPath)
	}

	// Generate or remove the static search index (after rendering so page.HTML is populated)
	if searchEnabled {
		if err := writeSearchIndex(dst, buildSearchIndex(pages, basePath)); err != nil {
			return fmt.Errorf("writing search index: %w", err)
		}
		fmt.Printf("  Generated %s\n", searchIndexFilename)
	} else {
		if err := removeSearchIndex(dst); err != nil {
			return fmt.Errorf("removing search index: %w", err)
		}
		fmt.Printf("  Search disabled (%s skipped)\n", searchIndexFilename)
	}

	// Generate or remove the RSS feed (after rendering so page.HTML is populated)
	if cfg.FeedEnabled() {
		feed := buildFeed(pages, cfg)
		if err := writeFeed(dst, feed); err != nil {
			return fmt.Errorf("writing feed: %w", err)
		}
		fmt.Printf("  Generated %s\n", feedFilename)
	} else {
		if err := removeFeed(dst); err != nil {
			return fmt.Errorf("removing feed: %w", err)
		}
	}

	// Copy _static directory
	staticSrc := filepath.Join(src, "_static")
	staticDst := filepath.Join(dst, "_static")
	if info, err := os.Stat(staticSrc); err == nil && info.IsDir() {
		if err := copyDir(staticSrc, staticDst); err != nil {
			return fmt.Errorf("copying _static: %w", err)
		}
		fmt.Printf("  Copied _static/\n")
	}

	fmt.Printf("Built %d pages → %s\n", len(pages), dst)
	return nil
}

// outputPathFromURL converts a URL path like "/guide/agents/" to a file path.
func outputPathFromURL(dst, urlPath string) string {
	p := strings.Trim(urlPath, "/")
	if p == "" {
		return filepath.Join(dst, "index.html")
	}
	return filepath.Join(dst, p, "index.html")
}

func renderToFile(tmpl *template.Template, data TemplateData, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, 0o644)
	})
}

// buildPageMeta creates a sorted list of PageMeta from all pages.
// Root index.md is excluded (matches nav behavior).
// Pages with dates sort reverse-chronologically first, then undated pages alphabetically.
func buildPageMeta(pages []Page, basePath string) []PageMeta {
	metas := make([]PageMeta, 0, len(pages))
	for _, p := range pages {
		if p.RelPath == "index.md" {
			continue
		}
		title := pageTitle(p)
		section := ""
		dir := filepath.Dir(p.RelPath)
		if dir != "." {
			parts := strings.SplitN(dir, string(filepath.Separator), 2)
			section = reNumPrefix.ReplaceAllString(parts[0], "")
		}
		metas = append(metas, PageMeta{
			Title:       title,
			Description: p.Frontmatter.Description,
			URL:         basePath + pageURL(p),
			Date:        p.Frontmatter.Date,
			Extra:       p.Frontmatter.Extra,
			Section:     section,
		})
	}

	sort.Slice(metas, func(i, j int) bool {
		// Dated pages first, reverse chronological
		if metas[i].Date != "" && metas[j].Date != "" {
			return metas[i].Date > metas[j].Date
		}
		if metas[i].Date != "" {
			return true
		}
		if metas[j].Date != "" {
			return false
		}
		return metas[i].Title < metas[j].Title
	})

	return metas
}

// writeSyntaxCSS generates a combined light/dark syntax highlighting stylesheet.
func writeSyntaxCSS(dst string, hl HighlightConfig) error {
	lightName := hl.Light
	if lightName == "" {
		lightName = "github"
	}
	darkName := hl.Dark
	if darkName == "" {
		darkName = "github-dark"
	}

	path := filepath.Join(dst, "_syntax.css")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	formatter := chromahtml.New(chromahtml.WithClasses(true))

	lightStyle := styles.Get(lightName)
	if lightStyle == nil {
		lightStyle = styles.Fallback
	}
	if err := formatter.WriteCSS(f, lightStyle); err != nil {
		return err
	}

	darkStyle := styles.Get(darkName)
	if darkStyle == nil {
		darkStyle = styles.Fallback
	}
	f.WriteString("\n/* Dark theme */\n")
	f.WriteString("[data-theme=\"dark\"] {\n")
	// Reset all chroma token colors to inherit from the dark base color.
	// This ensures tokens styled in the light theme but absent from the
	// dark theme (e.g. NameOther, Punctuation) don't keep their light colors.
	f.WriteString("  .chroma span { color: inherit; }\n")
	if err := writeScopedCSS(f, formatter, darkStyle); err != nil {
		return err
	}
	f.WriteString("}\n")

	fmt.Printf("  Generated _syntax.css (light: %s, dark: %s)\n", lightName, darkName)
	return nil
}

func writeScopedCSS(f *os.File, formatter *chromahtml.Formatter, style *chroma.Style) error {
	var buf strings.Builder
	if err := formatter.WriteCSS(&buf, style); err != nil {
		return err
	}
	for _, line := range strings.Split(buf.String(), "\n") {
		if strings.TrimSpace(line) != "" {
			f.WriteString("  " + line + "\n")
		}
	}
	return nil
}
