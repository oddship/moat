package main

import (
	"fmt"
	"path/filepath"
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
		sort.Slice(pages, func(i, j int) bool {
			return pages[i].RelPath < pages[j].RelPath
		})

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
func RenderNav(items []NavItem, currentPath, basePath string) string {
	var b strings.Builder
	b.WriteString("<nav>\n<ul>\n")

	for _, item := range items {
		if len(item.Children) > 0 {
			// Section with children
			b.WriteString(fmt.Sprintf("  <li>\n    <details open>\n      <summary>%s</summary>\n      <ul>\n", item.Title))
			for _, child := range item.Children {
				aria := ""
				href := basePath + child.Path
				if href == currentPath {
					aria = ` aria-current="page"`
				}
				b.WriteString(fmt.Sprintf("        <li><a href=\"%s\"%s>%s</a></li>\n", href, aria, child.Title))
			}
			b.WriteString("      </ul>\n    </details>\n  </li>\n")
		} else {
			// Top-level page
			aria := ""
			href := basePath + item.Path
			if href == currentPath {
				aria = ` aria-current="page"`
			}
			b.WriteString(fmt.Sprintf("  <li><a href=\"%s\"%s>%s</a></li>\n", href, aria, item.Title))
		}
	}

	b.WriteString("</ul>\n</nav>\n")
	return b.String()
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

// defaultURLPath converts "01-guide/02-agents.md" → "/guide/agents/"
func defaultURLPath(relPath string) string {
	// Remove .md
	p := strings.TrimSuffix(relPath, ".md")

	// Strip number prefixes from each path segment
	parts := strings.Split(p, string(filepath.Separator))
	for i, part := range parts {
		if len(part) >= 3 && part[0] >= '0' && part[0] <= '9' && part[1] >= '0' && part[1] <= '9' && part[2] == '-' {
			parts[i] = part[3:]
		}
	}
	p = strings.Join(parts, "/")

	// index pages → directory path
	if strings.HasSuffix(p, "/index") {
		p = strings.TrimSuffix(p, "index")
	} else if p == "index" {
		return "/"
	}

	return "/" + p + "/"
}
