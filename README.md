# moat

Markdown + [oat](https://oat.ink). A static site generator in one Go binary.

**[Documentation](https://oddship.github.io/moat/)**

```bash
go install github.com/oddship/moat@latest
mkdir docs
echo "# Hello" > docs/index.md
moat build docs/ _site/
moat serve _site/
```

That's it. No layout files needed — moat has a built-in oat layout with sidebar nav, dark mode, and syntax highlighting.

## Install

```bash
go install github.com/oddship/moat@latest
```

Or grab a binary from [releases](https://github.com/oddship/moat/releases).

## Features

- **Zero config** — built-in oat layout, just write markdown and build
- **Convention-based** — directory structure is the config, number prefixes control ordering
- **Syntax highlighting** — 70 Chroma themes with automatic light/dark mode
- **Layout inheritance** — base layout with `{{ block }}`/`{{ define }}` variants
- **Shortcodes** — reusable components inside markdown (`{{< note >}}...{{< /note >}}`)
- **Config file** — optional `config.toml` for site name, base path, highlight themes, sidebar links
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

[[links]]
title = "GitHub"
url = "https://github.com/you/project"

[extra]
footer = '&copy; <a href="https://example.com">You</a>'
```

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
- [chroma](https://github.com/alecthomas/chroma) — Syntax highlighting
- [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML frontmatter
- [toml](https://github.com/BurntSushi/toml) — Config file

## Why oat?

[oat](https://oat.ink) is an ultra-lightweight (~8KB) semantic HTML/CSS UI library — buttons, forms, cards, sidebar layouts, dark mode, all with zero dependencies. moat ships a built-in layout designed for oat, but you can use any CSS — just provide your own `_layout.html`.

## License

MIT
