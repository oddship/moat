package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// exportWikiParser records source spans while Goldmark determines which wiki
// links are real links. Code examples, escapes, and raw HTML remain untouched.
type exportWikiParser struct {
	wikilink.Parser
	spans map[*wikilink.Node]text.Segment
}

func (p *exportWikiParser) Parse(parent ast.Node, reader text.Reader, ctx parser.Context) ast.Node {
	line, segment := reader.PeekLine()
	n := p.Parser.Parse(parent, reader, ctx)
	if link, ok := n.(*wikilink.Node); ok {
		p.spans[link] = text.NewSegment(segment.Start, segment.Start+bytes.Index(line, []byte("]]"))+2)
	}
	return n
}

func markdownLabel(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]").Replace(s)
}

// exportMarkdown preserves the authored Markdown and expanded shortcode HTML.
// Only actual wiki links are rewritten; this avoids an HTML-to-Markdown round
// trip that could lose code, custom components, or table structure.
func exportMarkdown(body []byte, resolver *pageResolver) []byte {
	p := &exportWikiParser{spans: make(map[*wikilink.Node]text.Segment)}
	md := newMarkdown(resolver)
	md.Parser().AddOptions(parser.WithInlineParsers(util.Prioritized(p, 198)))
	doc := md.Parser().Parse(text.NewReader(body))
	type replacement struct {
		span text.Segment
		text string
	}
	var replacements []replacement
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		link, ok := n.(*wikilink.Node)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}
		span, ok := p.spans[link]
		if !ok {
			return ast.WalkContinue, nil
		}
		destination, _ := resolver.ResolveWikilink(link)
		label := markdownLabel(string(link.Text(body)))
		value := label
		if len(destination) > 0 {
			dest := string(destination)
			if u, err := url.Parse(dest); err == nil {
				dest = u.String()
			}
			dest = strings.NewReplacer("<", "%3C", ">", "%3E").Replace(dest)
			value = "[" + label + "](<" + dest + ">)"
			if link.Embed {
				value = "!" + value
			}
		}
		replacements = append(replacements, replacement{span, value})
		return ast.WalkSkipChildren, nil
	})
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].span.Start < replacements[j].span.Start })
	var out bytes.Buffer
	position := 0
	for _, r := range replacements {
		out.Write(body[position:r.span.Start])
		out.WriteString(r.text)
		position = r.span.Stop
	}
	out.Write(body[position:])
	if out.Len() == 0 || body[len(body)-1] != '\n' {
		out.WriteByte('\n')
	}
	return out.Bytes()
}

func writeAgentCatalogs(dst string, pages []Page, cfg Config) error {
	name := strings.Join(strings.Fields(cfg.SiteName), " ")
	if name == "" {
		name = "Site"
	}
	var index, full strings.Builder
	fmt.Fprintf(&index, "# %s\n\n## Docs\n\n", name)
	fmt.Fprintf(&full, "# %s\n", name)
	for _, page := range pages {
		markdownURL := siteURL(cfg.BasePath, pageURL(page)+"index.md")
		// Angle brackets permit spaces and parentheses in configured URLs.
		u := (&url.URL{Path: markdownURL}).String()
		fmt.Fprintf(&index, "- [%s](<%s>)", markdownLabel(pageTitle(page)), u)
		if page.Frontmatter.Description != "" {
			fmt.Fprintf(&index, ": %s", strings.Join(strings.Fields(page.Frontmatter.Description), " "))
		}
		index.WriteByte('\n')
		if cfg.Agents.Full {
			fmt.Fprintf(&full, "\n---\n\n# %s\n\nSource: <%s>\n\n", strings.Join(strings.Fields(pageTitle(page)), " "), u)
			full.Write(page.Markdown)
		}
	}
	if err := os.WriteFile(filepath.Join(dst, "llms.txt"), []byte(index.String()), 0o644); err != nil {
		return err
	}
	if cfg.Agents.Full {
		return os.WriteFile(filepath.Join(dst, "llms-full.txt"), []byte(full.String()), 0o644)
	}
	return nil
}
