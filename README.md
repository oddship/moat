# moat

Markdown + [oat](https://github.com/knadh/oat). A static site generator in one binary.

```
moat build docs/ _site/
moat serve _site/
```

## What it does

Reads markdown files from a directory, converts them to HTML, wraps them in a
layout template styled with oat CSS, generates a sidebar nav, and writes
static HTML files. That's it.

## Install

```bash
go install github.com/oddship/moat@latest
```

Or build from source:

```bash
git clone https://github.com/oddship/moat.git
cd moat && go build .
```

## Usage

### Build

```bash
moat build <source-dir> <output-dir>
```

Source directory should contain:
- `_layout.html` — Go template for page layout (required)
- `*.md` files — your content
- `_static/` — static assets copied as-is (optional)

### Serve

```bash
moat serve <dir> [--port 3000]
```

Simple static file server for local preview. Default port: 8080.

## Conventions

### Directory structure

```
docs/
├── _layout.html          # Required. Go template.
├── _static/              # Copied to output as-is
├── index.md              # → /index.html
├── quickstart.md         # → /quickstart/index.html
├── guide/
│   ├── 01-intro.md       # → /guide/intro/index.html
│   └── 02-advanced.md    # → /guide/advanced/index.html
└── reference/
    └── api.md            # → /reference/api/index.html
```

### Frontmatter

Optional YAML frontmatter:

```yaml
---
title: Getting Started
description: How to set up your project
url: getting-started
---
```

- **title**: Page title. If omitted, derived from filename (`02-agents.md` → "Agents").
- **description**: Meta description for the page.
- **url**: Override the URL path. If omitted, derived from file path with number prefixes stripped.

### Navigation

Nav is generated automatically from the directory structure:
- Root `.md` files appear as top-level links
- Directories become collapsible sections
- Files are sorted alphabetically within sections
- Number prefixes (`01-`, `02-`) control order but are stripped from display and URLs
- `index.md` is excluded from nav listings

### Layout template

`_layout.html` is a Go template with these variables:

```html
{{ .Title }}          <!-- Page title -->
{{ .Description }}    <!-- Page description -->
{{ .Content }}        <!-- Rendered HTML content -->
{{ .Nav }}            <!-- Generated navigation HTML -->
{{ .CurrentPath }}    <!-- Current page URL path -->
{{ .SiteName }}       <!-- Site name (from flag or "Site") -->
```

### Clean URLs

All pages get clean URLs:
- `quickstart.md` → `/quickstart/` (served from `/quickstart/index.html`)
- `guide/01-intro.md` → `/guide/intro/`
- `index.md` → `/` (root index)

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) — Markdown parser (with GFM tables, strikethrough, autolinks)
- [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML frontmatter

## License

MIT
