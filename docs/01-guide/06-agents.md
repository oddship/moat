---
title: Agent-friendly docs
description: Copy pages as Markdown and publish static llms.txt exports for AI assistants
---

# Agent-friendly docs

Every build publishes a Markdown version of each page and an `llms.txt` catalog.
These are ordinary static files, so they work on GitHub Pages and other static
hosts without a backend or API key.

## Copy or view a page

The built-in layout includes a compact **Copy page** control beside the heading,
with **View Markdown** in its dropdown menu. On small screens, the copy control
uses an icon to leave room for the title. Copy page copies the exported Markdown rather than navigation,
footer text, or the page's HTML layout. Clipboard access requires HTTPS or
localhost and browser permission; if copying fails, View Markdown remains
available for manual copying.

A page at `/guide/quickstart/` has its Markdown export at
`/guide/quickstart/index.md`. The homepage uses `/index.md`. A configured
`base_path` prefixes both the catalog links and page actions, so a project site
uses URLs such as `/moat/guide/quickstart/index.md`.

Exports omit YAML frontmatter, preserve authored Markdown, and resolve wiki
links into standard links to the published pages. Shortcodes expand as they do
on the website, retaining their generated HTML inside the Markdown. This keeps
custom components and generated listings available without converting them to
plain text. Code examples containing wiki-link syntax remain unchanged.

## Documentation catalog

Give an assistant your site's `llms.txt` URL to help it find relevant pages.
The catalog lists published pages with their titles, descriptions, and Markdown
URLs. Links are root-relative, including `base_path`; they resolve against the
origin serving the catalog. The same build controls HTML and exports: drafts
and hidden source files are excluded, and rebuilding removes stale exports.

The built-in layout advertises the page's Markdown through
`rel="alternate" type="text/markdown"` and the catalog through
`rel="describedby"`. Publishing these files makes them available to agents;
whether an assistant discovers them automatically depends on the client.

## Configuration

Exports and page actions are enabled by default. Disable them together:

```toml
[agents]
enabled = false
```

For assistants that prefer all documentation in one download, enable the
optional `llms-full.txt` bundle:

```toml
[agents]
full = true
```

The full bundle includes each published page's title, Markdown URL, and exported
content. It can be large, so it is disabled by default. Each independently built
site version has its own catalog and exports.

## Custom layouts

Custom layouts receive `.MarkdownURL` and `.LLMsURL`; both are empty when
exports are disabled. Use these fields to add your own links and copy controls.
Static files are generated even when your custom layout omits those controls.
