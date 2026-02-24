---
title: CLI Reference
description: Command-line interface for moat
---

# CLI Reference

## `moat init`

Scaffold a docs directory with layouts, config, and starter content.

```bash
moat init [dir]
```

Default directory is `docs`. Creates:

- `_layout.html` — copy of the built-in base layout
- `_layout.landing.html` — landing page variant
- `config.toml` — sample config with all options commented
- `index.md` — landing page
- `01-guide/01-getting-started.md` — starter content

Won't overwrite existing files — run on a new directory.

```bash
moat init              # creates docs/
moat init my-docs      # creates my-docs/
```

## `moat build`

Build a static site from markdown source.

```bash
moat build <src> <dst> [flags]
```

| Flag | Description |
|------|-------------|
| `--config PATH` | Config file (default: `<src>/config.toml`) |
| `--site-name NAME` | Site name for templates |
| `--base-path PATH` | URL prefix for GitHub project pages |

```bash
# Basic build
moat build docs/ _site/

# With GitHub project pages prefix
moat build docs/ _site/ --base-path /my-project

# Custom config location
moat build docs/ _site/ --config site.toml
```

CLI flags override values from `config.toml`.

If no `_layout.html` exists in the source directory, moat uses its built-in oat layout.

## `moat serve`

Serve a built site for local preview.

```bash
moat serve <dir> [--port PORT]
```

```bash
moat serve _site/
moat serve _site/ --port 3000
```

{{< note type="info" >}}
`moat serve` is a simple static file server for previewing builds. For development with live reload, use a tool like [browser-sync](https://browsersync.io/) or rebuild on file change with `watchexec`.
{{< /note >}}

## `moat version`

Print the moat version.

```bash
moat version
```
