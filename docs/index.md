---
title: moat
layout: landing
description: A static site generator in one Go binary. Markdown + oat.
---

# moat

<h3 class="tagline text-light">Docs are your project's moat. Build them from markdown in one command.</h3>

<div class="hstack">
  <a href="guide/quickstart/" class="button">Get started</a>
  <a href="https://github.com/oddship/moat" class="button outline">GitHub</a>
</div>

<br>

<div class="features">
<article class="card">
<header><h3>Convention over config</h3></header>

Directory structure is the config. Number prefixes control ordering.
No YAML config sprawl.
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
moat build docs/ _site/
moat serve _site/
```

## What it does

Reads markdown files from a directory, converts to HTML, wraps in a layout template, generates sidebar nav, writes static HTML. That's it.

```
docs/
├── config.toml           # Site config (optional)
├── _layout.html          # Go template (required)
├── _shortcodes/          # Reusable components
├── _static/              # Copied as-is
├── index.md              # → /
└── 01-guide/
    ├── 01-intro.md       # → /guide/intro/
    └── 02-config.md      # → /guide/config/
```

Pairs with [oat](https://github.com/knadh/oat) for styling — or use any CSS you want.
