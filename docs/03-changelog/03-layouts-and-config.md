---
title: "v0.4.0 — Layouts, Shortcodes & Config"
date: 2026-02-25
description: Layout inheritance, shortcodes, config.toml, syntax highlighting, base_path, and moat init.
---

# v0.4.0 — Layouts, Shortcodes & Config

Major feature release adding customization and extensibility.

## Layout inheritance

Base layout with `{{ block }}`/`{{ define }}` variants. Create `_layout.html` to override the built-in, or add `_layout.landing.html` for named variants. See [[Layouts]] for details.

## Shortcodes

Reusable HTML components inside markdown via `_shortcodes/` directory. Block and self-closing syntax. See [[Shortcodes]] for usage and examples.

## Config file

Optional `config.toml` for site name, base path, syntax highlighting themes, and sidebar links. See [[Configuration]] for all options.

## Syntax highlighting

70 Chroma themes with automatic light/dark mode via CSS classes. Configure theme pairs in `config.toml`.

## Other additions

- `--base-path` flag for GitHub project pages
- `moat init` to scaffold a new docs directory with built-in layouts
- Sidebar links via `[[links]]` config
- Embedded default layouts (no `_layout.html` required)
