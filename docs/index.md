---
title: moat
layout: landing
description: A static site generator in one Go binary. Markdown + oat.
---

# moat

<h3 class="tagline text-light">Markdown + <a href="https://oat.ink">oat</a>. A static site generator in one Go binary.</h3>

<div class="hstack">
  <a href="guide/quickstart/" class="button">Get started</a>
  <a href="https://github.com/oddship/moat" class="button outline">GitHub</a>
</div>

<br>

<div class="features">
<article class="card">
<header><h3>Zero config</h3></header>

Built-in oat layout with sidebar nav, dark mode, and syntax highlighting.
Just write markdown and build. No layout files needed.
</article>

<article class="card">
<header><h3>One binary</h3></header>

Single Go binary. No Node.js, no npm, no Ruby.
`go install` or grab a release binary.
</article>

<article class="card">
<header><h3>Syntax highlighting</h3></header>

70 Chroma themes with automatic light/dark mode.
Configure in `config.toml` — CSS classes, not inline styles.
</article>

<article class="card">
<header><h3>Layouts & shortcodes</h3></header>

Multiple layouts via Go template `block`/`define`. Shortcodes for reusable
components in markdown. Like Zola, without the complexity.
</article>
</div>

## Quick start

```bash
go install github.com/oddship/moat@latest
mkdir docs
echo "# Hello" > docs/index.md
moat build docs/ _site/
moat serve _site/
```

## What it does

Reads markdown files from a directory, converts to HTML, wraps in a layout template, generates sidebar nav, writes static HTML. That's it.

```
docs/
├── config.toml           # Site config (optional)
├── _layout.html          # Custom layout (optional)
├── _shortcodes/          # Reusable components
├── _static/              # Copied as-is
├── index.md              # → /
└── 01-guide/
    ├── 01-intro.md       # → /guide/intro/
    └── 02-config.md      # → /guide/config/
```

Ships with a built-in [oat](https://oat.ink) layout — sidebar, dark mode, semantic HTML. Or provide your own `_layout.html` with any CSS you want.
