---
title: "v0.1.0 — Initial Release"
date: 2026-02-25
description: First release — markdown to static HTML with oat layout, syntax highlighting, and GitHub Actions.
---

# v0.1.0 — Initial Release

moat's first release. A static site generator in one Go binary.

## Features

- Markdown to HTML via goldmark with GFM support
- Built-in oat layout with sidebar nav and dark mode
- Convention-based directory structure — number prefixes for ordering
- Syntax highlighting with 70 Chroma themes
- Layout inheritance with Go template `block`/`define`
- [[Shortcodes]] for reusable components inside markdown
- [[GitHub Actions]] reusable workflow for one-file deployment
- Optional `config.toml` for site name, base path, highlighting themes, and sidebar links
- `moat init` to scaffold a new docs directory
