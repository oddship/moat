package main

import (
	"bytes"
	"fmt"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/wikilink"
)

// newMarkdown creates a goldmark instance with optional wikilink resolver.
// If resolver is nil, wikilinks are still parsed but resolved with the default
// resolver (appends .html).
func newMarkdown(resolver wikilink.Resolver) goldmark.Markdown {
	wlExt := &wikilink.Extender{}
	if resolver != nil {
		wlExt.Resolver = resolver
	}

	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
			wlExt,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
}

// RenderMarkdown converts markdown bytes to HTML.
// Uses the default goldmark instance (no page-aware wikilink resolution).
func RenderMarkdown(source []byte) ([]byte, error) {
	return renderMarkdownWith(newMarkdown(nil), source)
}

// RenderMarkdownWithResolver converts markdown bytes to HTML using
// a page-aware wikilink resolver.
func RenderMarkdownWithResolver(source []byte, resolver wikilink.Resolver) ([]byte, error) {
	return renderMarkdownWith(newMarkdown(resolver), source)
}

func renderMarkdownWith(md goldmark.Markdown, source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// pageResolver resolves [[wiki links]] to page URLs using a title→URL map.
// Matching is case-insensitive and ignores leading/trailing whitespace.
type pageResolver struct {
	pages map[string]string // lowercase title → URL path
}

// newPageResolver builds a resolver from a list of pages.
// If multiple pages share the same title (case-insensitive), the first one wins
// and a warning is printed.
func newPageResolver(pages []Page, basePath string) *pageResolver {
	m := make(map[string]string, len(pages))
	for _, p := range pages {
		title := pageTitle(p)
		url := basePath + pageURL(p)
		key := strings.ToLower(title)
		if existing, ok := m[key]; ok {
			fmt.Printf("  Warning: duplicate wiki link target %q (%s shadows %s)\n", title, existing, url)
			continue
		}
		m[key] = url
	}
	return &pageResolver{pages: m}
}

func (r *pageResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	target := strings.TrimSpace(string(n.Target))
	key := strings.ToLower(target)

	dest, ok := r.pages[key]
	if !ok {
		// Unknown page — render as plain text (nil destination)
		return nil, nil
	}

	if len(n.Fragment) > 0 {
		dest += "#" + string(n.Fragment)
	}

	return []byte(dest), nil
}
