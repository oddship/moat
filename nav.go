package main

import (
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// NavItem represents a page or section in the navigation.
type NavItem struct {
	Title    string
	Path     string // URL path relative to site root, e.g. "/guide/agents"
	Children []NavItem
}

// BuildNav constructs a navigation tree from a list of pages.
// Pages are grouped by directory, sorted alphabetically within each group.
// Number prefixes (01-, 02-) control ordering but are stripped from display.
func BuildNav(pages []Page) []NavItem {
	// Group by directory
	rootPages := []Page{}
	sections := map[string][]Page{}
	sectionOrder := []string{}

	for _, p := range pages {
		if p.RelPath == "index.md" {
			continue // Skip root index from nav listing
		}

		dir := filepath.Dir(p.RelPath)
		if dir == "." {
			rootPages = append(rootPages, p)
		} else {
			// Use top-level directory only (max 2 levels)
			parts := strings.SplitN(dir, string(filepath.Separator), 2)
			section := parts[0]
			if _, exists := sections[section]; !exists {
				sectionOrder = append(sectionOrder, section)
			}
			sections[section] = append(sections[section], p)
		}
	}

	sort.Strings(sectionOrder)

	var nav []NavItem

	// Root-level pages first (sorted)
	sort.Slice(rootPages, func(i, j int) bool {
		return rootPages[i].RelPath < rootPages[j].RelPath
	})
	for _, p := range rootPages {
		nav = append(nav, NavItem{
			Title: pageTitle(p),
			Path:  pageURL(p),
		})
	}

	// Then sections
	for _, section := range sectionOrder {
		pages := sections[section]

		// If any page has a date, sort reverse-chronologically (newest first).
		// Otherwise, sort alphabetically by path.
		hasDate := false
		for _, p := range pages {
			if p.Frontmatter.Date != "" {
				hasDate = true
				break
			}
		}
		if hasDate {
			sort.Slice(pages, func(i, j int) bool {
				di, dj := pages[i].Frontmatter.Date, pages[j].Frontmatter.Date
				if di != dj {
					return di > dj // reverse chronological
				}
				return pages[i].RelPath < pages[j].RelPath
			})
		} else {
			sort.Slice(pages, func(i, j int) bool {
				return pages[i].RelPath < pages[j].RelPath
			})
		}

		children := []NavItem{}
		for _, p := range pages {
			children = append(children, NavItem{
				Title: pageTitle(p),
				Path:  pageURL(p),
			})
		}

		nav = append(nav, NavItem{
			Title:    TitleFromDir(section),
			Children: children,
		})
	}

	return nav
}

// RenderNav generates HTML for the navigation sidebar.
// Uses oat's sidebar nav patterns: <ul> lists, <details> for sections, aria-current for active.
// basePath is prepended to all href values (e.g. "/moat" for GitHub project pages).
// links are extra items rendered at the top of the nav (e.g. GitHub link).
func RenderNav(items []NavItem, currentPath, basePath string, links []LinkConfig) string {
	var b strings.Builder
	b.WriteString("<nav>\n<ul>\n")

	// Extra links first (e.g. GitHub)
	for _, link := range links {
		icon := linkIcon(link.Icon)
		b.WriteString(fmt.Sprintf("  <li><a href=\"%s\">%s%s</a></li>\n",
			html.EscapeString(link.URL), icon, html.EscapeString(link.Title)))
	}

	for _, item := range items {
		if len(item.Children) > 0 {
			// Section with children
			b.WriteString(fmt.Sprintf("  <li>\n    <details open>\n      <summary>%s</summary>\n      <ul>\n",
				html.EscapeString(item.Title)))
			for _, child := range item.Children {
				aria := ""
				href := basePath + child.Path
				if href == currentPath {
					aria = ` aria-current="page"`
				}
				b.WriteString(fmt.Sprintf("        <li><a href=\"%s\"%s>%s</a></li>\n",
					html.EscapeString(href), aria, html.EscapeString(child.Title)))
			}
			b.WriteString("      </ul>\n    </details>\n  </li>\n")
		} else {
			// Top-level page
			aria := ""
			href := basePath + item.Path
			if href == currentPath {
				aria = ` aria-current="page"`
			}
			b.WriteString(fmt.Sprintf("  <li><a href=\"%s\"%s>%s</a></li>\n",
				html.EscapeString(href), aria, html.EscapeString(item.Title)))
		}
	}

	b.WriteString("</ul>\n</nav>\n")
	return b.String()
}

// Built-in SVG icons for sidebar links.
var builtinIcons = map[string]string{
	"github": `<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" style="vertical-align: middle"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>`,
}

// linkIcon returns the inline SVG for a built-in icon name, or empty string.
func linkIcon(name string) string {
	if svg, ok := builtinIcons[name]; ok {
		return svg + " "
	}
	return ""
}

func pageTitle(p Page) string {
	if p.Frontmatter.Title != "" {
		return p.Frontmatter.Title
	}
	return TitleFromFilename(filepath.Base(p.RelPath))
}

// pageURL returns the URL for a page, using frontmatter url if set.
func pageURL(p Page) string {
	if p.Frontmatter.URL != "" {
		u := p.Frontmatter.URL
		if !strings.HasPrefix(u, "/") {
			u = "/" + u
		}
		if !strings.HasSuffix(u, "/") {
			u = u + "/"
		}
		return u
	}
	return defaultURLPath(p.RelPath)
}

// reNumPrefix matches leading digits followed by a hyphen: "01-", "1-", "001-"
var reNumPrefix = regexp.MustCompile(`^\d+-`)

// defaultURLPath converts "01-guide/02-agents.md" → "/guide/agents/"
func defaultURLPath(relPath string) string {
	// Remove .md
	p := strings.TrimSuffix(relPath, ".md")

	// Strip number prefixes from each path segment
	parts := strings.Split(p, string(filepath.Separator))
	for i, part := range parts {
		parts[i] = reNumPrefix.ReplaceAllString(part, "")
	}
	p = strings.Join(parts, "/")

	// index pages → directory path
	if strings.HasSuffix(p, "/index") {
		p = strings.TrimSuffix(p, "/index")
		return "/" + p + "/"
	} else if p == "index" {
		return "/"
	}

	return "/" + p + "/"
}
