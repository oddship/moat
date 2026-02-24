---
title: moat
description: A static site generator in one binary
---

# moat

Markdown + [oat](https://github.com/knadh/oat). A static site generator in one binary.

```
moat build docs/ _site/
moat serve _site/
```

Reads markdown files from a directory, converts them to HTML, wraps them in a layout template, generates a sidebar nav, and writes static HTML. That's it.

## Install

```bash
go install github.com/oddship/moat@latest
```

Or grab a binary from [releases](https://github.com/oddship/moat/releases).

## Quick start

```
mkdir docs
```

Create `docs/_layout.html`:

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="https://unpkg.com/@knadh/oat/oat.min.css">
</head>
<body>
  <nav>{{ .Nav }}</nav>
  <main>{{ .Content }}</main>
</body>
</html>
```

Create `docs/index.md`:

```markdown
# Hello

This is my site.
```

Build and serve:

```bash
moat build docs/ _site/
moat serve _site/
```

## What you get

- **Clean URLs** — `guide/01-intro.md` becomes `/guide/intro/`
- **Auto nav** — sidebar generated from directory structure
- **Frontmatter** — title, description, custom URLs
- **Sections** — directories become collapsible nav groups
- **Number prefixes** — `01-`, `02-` control order, stripped from URLs
- **Static assets** — `_static/` directory copied as-is
- **Base path** — `--base-path /repo` for GitHub project pages
- **Zero config** — no config files, just conventions

## Why

Static site generators are either too simple (no nav, no sections) or too complex (config files, themes, plugins, build chains). moat sits in between: one binary, one layout file, conventions over configuration.

It pairs well with [oat](https://github.com/knadh/oat) for styling but works with any CSS framework — the layout template is plain HTML.
