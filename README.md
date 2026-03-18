<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/_static/logo-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="docs/_static/logo-light.svg">
    <img src="docs/_static/logo-light.svg" alt="moat" height="48">
  </picture>
</p>

<p align="center">Markdown + <a href="https://oat.ink">oat</a>. A static site generator in one Go binary.</p>

**[Documentation](https://oddship.github.io/moat/)**

```bash
go install github.com/oddship/moat@latest
mkdir docs
echo "# Hello" > docs/index.md
moat build docs/ _site/
moat serve _site/
```

That's it. No layout files needed — moat has a built-in oat layout with sidebar nav, built-in search, dark mode, and syntax highlighting.

## Install

```bash
go install github.com/oddship/moat@latest
```

Or grab a binary from [releases](https://github.com/oddship/moat/releases).

## Features

- **Zero config** — built-in oat layout, just write markdown and build
- **Convention-based** — directory structure is the config, number prefixes control ordering
- **Built-in search** — client-side `_search.json` index with oat-native sidebar search UI
- **Wiki links** — `[[Page Title]]` resolves to internal page URLs
- **Site primitives** — `date`, `draft`, page listings via templates/shortcodes, and optional RSS feed
- **Syntax highlighting** — 70 Chroma themes with automatic light/dark mode
- **Layout inheritance** — base layout with `{{ block }}`/`{{ define }}` variants
- **Shortcodes** — reusable components inside markdown (`{{< note >}}...{{< /note >}}`)
- **Config file** — optional `config.toml` for site name, base path, highlight themes, navigation links, search, and feed settings
- **GitHub Actions** — reusable workflow for one-line GitHub Pages deployment
- **Single binary** — no Node.js, no npm, just Go

## Quick start

```bash
# Zero-config — just markdown
mkdir docs
echo "# Hello" > docs/index.md
moat build docs/ _site/

# Or scaffold with layouts you can customize
moat init docs
moat build docs/ _site/
```

## Directory structure

```
docs/
├── config.toml               # Site config (optional)
├── _layout.html              # Custom layout (optional — built-in used if absent)
├── _layout.landing.html      # Layout variant (optional)
├── _shortcodes/              # Shortcode templates
│   └── note.html
├── _static/                  # Copied as-is
├── index.md                  # → /
└── 01-guide/
    ├── 01-intro.md           # → /guide/intro/
    └── 02-config.md          # → /guide/config/
```

## Configuration

Optional `config.toml` in your docs directory:

```toml
site_name = "My Project"
base_path = "/my-project"

[highlight]
light = "github"
dark  = "github-dark"

[search]
enabled = false

[feed]
enabled = true
link = "https://docs.example.com"

[[links]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"

[[topnav]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"

[extra]
footer = '&copy; <a href="https://example.com">You</a>'
```

By default, `moat build` also emits `_search.json`, and the built-in oat layout wires up a sidebar search box automatically. Disable it with:

```toml
[search]
enabled = false
```

Enable RSS feed generation with:

```toml
[feed]
enabled = true
link = "https://docs.example.com"
# optional
# title = "My Site Feed"
```

Only pages with `date: YYYY-MM-DD` are included in `feed.xml`, newest first.

Wiki links are also supported:

```md
See [[Getting Started]] for the overview.
```

Pages and shortcodes can use site-level page data for listings via templates (`{{ .Pages }}`) or shortcode context (`{{ .Page.Pages }}` / `{{ .SectionPages "guide" }}`).

## GitHub Pages

Deploy with moat's reusable workflow — one file, zero config:

```yaml
# .github/workflows/docs.yml
name: Docs
on:
  push:
    branches: [main]
    paths: ['docs/**']

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  docs:
    uses: oddship/moat/.github/workflows/build-docs.yml@main
    with:
      docs_dir: docs
```

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) — Markdown with GFM
- [goldmark/wikilink](https://pkg.go.dev/go.abhg.dev/goldmark/wikilink) — `[[Wiki Links]]`
- [chroma](https://github.com/alecthomas/chroma) — Syntax highlighting
- [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML frontmatter
- [toml](https://github.com/BurntSushi/toml) — Config file

## Why oat?

[oat](https://oat.ink) is an ultra-lightweight (~8KB) semantic HTML/CSS UI library — buttons, forms, cards, sidebar layouts, dark mode, all with zero dependencies. moat ships a built-in layout designed for oat, but you can use any CSS — just provide your own `_layout.html`.

## License

MIT
