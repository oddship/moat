package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Page represents a single markdown page.
type Page struct {
	RelPath     string      // Relative path from source root, e.g. "guide/02-agents.md"
	Frontmatter Frontmatter // Parsed YAML frontmatter
	Body        []byte      // Markdown body (without frontmatter)
	HTML        []byte      // Rendered HTML
}

// TemplateData is passed to the layout template.
type TemplateData struct {
	Title       string
	Description string
	Content     template.HTML
	Nav         template.HTML
	CurrentPath string
	SiteName    string
}

// Build reads markdown from src, renders HTML, and writes to dst.
func Build(src, dst, siteName string) error {
	src, _ = filepath.Abs(src)

	if siteName == "" {
		siteName = "Site"
	}
	dst, _ = filepath.Abs(dst)

	// Load layout template
	layoutPath := filepath.Join(src, "_layout.html")
	layoutBytes, err := os.ReadFile(layoutPath)
	if err != nil {
		return fmt.Errorf("reading _layout.html: %w (expected at %s)", err, layoutPath)
	}
	tmpl, err := template.New("layout").Parse(string(layoutBytes))
	if err != nil {
		return fmt.Errorf("parsing _layout.html: %w", err)
	}

	// Discover and parse markdown files
	var pages []Page
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip underscore-prefixed dirs/files (except root _layout.html handled above)
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
		html, err := RenderMarkdown(body)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", relPath, err)
		}

		pages = append(pages, Page{
			RelPath:     relPath,
			Frontmatter: fm,
			Body:        body,
			HTML:        html,
		})
		return nil
	})
	if err != nil {
		return fmt.Errorf("walking source: %w", err)
	}

	fmt.Printf("Found %d pages\n", len(pages))

	// Build navigation
	nav := BuildNav(pages)

	// Render each page
	for _, page := range pages {
		currentPath := pageURL(page)
		outPath := outputPathFromURL(dst, currentPath)

		navHTML := RenderNav(nav, currentPath)

		title := page.Frontmatter.Title
		if title == "" {
			title = TitleFromFilename(filepath.Base(page.RelPath))
		}

		data := TemplateData{
			Title:       title,
			Description: page.Frontmatter.Description,
			Content:     template.HTML(page.HTML),
			Nav:         template.HTML(navHTML),
			CurrentPath: currentPath,
			SiteName:    siteName,
		}

		if err := renderToFile(tmpl, data, outPath); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
		fmt.Printf("  %s → %s\n", page.RelPath, outPath)
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
	// Strip leading/trailing slashes
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
