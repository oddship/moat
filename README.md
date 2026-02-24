# moat

Markdown + [oat](https://oat.ink). A static site generator in one Go binary.

**[Documentation](https://oddship.github.io/moat/)**

```bash
moat build docs/ _site/
moat serve _site/
```

Reads markdown, wraps it in a Go template layout, generates sidebar nav, writes static HTML. Ships with layouts designed for [oat](https://github.com/knadh/oat) — an ultra-lightweight semantic HTML/CSS UI library.

## Install

```bash
go install github.com/oddship/moat@latest
```

Or grab a binary from [releases](https://github.com/oddship/moat/releases).

## Features

- **oat-ready** — example layouts use [oat](https://oat.ink) for sidebar, topnav, dark mode out of the box
- **Convention-based** — directory structure is the config, number prefixes control ordering
- **Syntax highlighting** — 70 Chroma themes with automatic light/dark mode
- **Layout inheritance** — base layout with `{{ block }}`/`{{ define }}` variants
- **Shortcodes** — reusable components inside markdown (`{{< note >}}...{{< /note >}}`)
- **Config file** — optional `config.toml` for site name, base path, highlight themes
- **GitHub Actions** — reusable workflow for one-line GitHub Pages deployment
- **Single binary** — no Node.js, no npm, just Go

## Quick start

```bash
mkdir docs
cat > docs/_layout.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
    <link rel="stylesheet" href="https://unpkg.com/@knadh/oat/oat.min.css">
    <link rel="stylesheet" href="{{ .BasePath }}/_syntax.css">
</head>
<body>
    <nav>{{ .Nav }}</nav>
    <main>{{ block "content" . }}<article>{{ .Content }}</article>{{ end }}</main>
</body>
</html>
EOF

echo "# Hello" > docs/index.md
moat build docs/ _site/
moat serve _site/
```

## Directory structure

```
docs/
├── config.toml               # Site config (optional)
├── _layout.html              # Base layout (required)
├── _layout.landing.html      # Layout variant (optional)
├── _shortcodes/              # Shortcode templates
│   └── note.html
├── _static/                  # Copied as-is
├── index.md                  # → /
└── 01-guide/
    ├── 01-intro.md           # → /guide/intro/
    └── 02-config.md          # → /guide/config/
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

[oat](https://oat.ink) is an ultra-lightweight (~8KB) semantic HTML/CSS UI library — buttons, forms, cards, sidebar layouts, dark mode, all with zero dependencies. moat doesn't require oat, but the example layouts use it and the nav HTML matches oat's sidebar patterns. Include the CSS and you get a site that looks like [oat.ink](https://oat.ink) with just markdown files. Or use any CSS you prefer.

## License

MIT
