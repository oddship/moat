# moat

Docs are your project's moat. A static site generator in one Go binary.

**[Documentation](https://oddship.github.io/moat/)**

```
moat build docs/ _site/
moat serve _site/
```

Reads markdown, wraps it in a layout template, generates sidebar nav, writes static HTML. No config files — just conventions.

## Install

```bash
go install github.com/oddship/moat@latest
```

Or grab a binary from [releases](https://github.com/oddship/moat/releases).

## Usage

```bash
moat build <src> <dst> [--site-name NAME] [--base-path PATH]
moat serve <dir> [--port PORT]
```

Source directory contains:

```
docs/
├── _layout.html          # Go template (required)
├── _static/              # Copied as-is
├── index.md              # → /
├── quickstart.md         # → /quickstart/
└── 01-guide/
    ├── 01-intro.md       # → /guide/intro/
    └── 02-advanced.md    # → /guide/advanced/
```

Number prefixes (`01-`, `02-`) control ordering but are stripped from URLs and display names.

## Template variables

```html
{{ .Title }}        {{ .Description }}    {{ .Content }}
{{ .Nav }}          {{ .CurrentPath }}    {{ .SiteName }}
{{ .BasePath }}
```

Pairs well with [oat](https://github.com/knadh/oat) for styling — see the [layout guide](https://oddship.github.io/moat/guide/layout/).

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) — Markdown with GFM tables, strikethrough, autolinks
- [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML frontmatter

## License

MIT
