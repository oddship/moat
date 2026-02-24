---
title: CLI Reference
description: Command-line interface for moat
---

# CLI Reference

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
