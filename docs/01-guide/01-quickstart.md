---
title: Quick Start
description: Get a docs site running in 5 minutes
---

# Quick start

## Install

```bash
go install github.com/oddship/moat@latest
```

Or grab a binary from [releases](https://github.com/oddship/moat/releases).

## Create a site

The fastest way — moat has built-in layouts so you don't need to provide any:

```bash
mkdir docs
echo "# Hello" > docs/index.md
moat build docs/ _site/
moat serve _site/
```

Open [http://localhost:8080](http://localhost:8080). You get a full oat-styled site with sidebar nav, dark mode toggle, and syntax highlighting.

## Scaffold with `moat init`

For a more complete starting point:

```bash
moat init docs
moat build docs/ _site/
moat serve _site/
```

This creates layouts, a sample config, and starter pages you can edit:

```
docs/
├── _layout.html          # Base layout (customize or delete to use built-in)
├── _layout.landing.html  # Landing page variant
├── config.toml           # Site config
├── index.md              # Home page
└── 01-guide/
    └── 01-getting-started.md
```

## Add more pages

```bash
mkdir docs/01-guide
echo "# Getting Started" > docs/01-guide/01-getting-started.md
echo "# Configuration"   > docs/01-guide/02-configuration.md
```

Number prefixes (`01-`, `02-`) control ordering but are stripped from URLs and nav:

- `01-guide/01-getting-started.md` → `/guide/getting-started/`
- `01-guide/02-configuration.md` → `/guide/configuration/`

Directories become collapsible sections in the sidebar automatically.

## Deploy to GitHub Pages

Add this to your repo — or use moat's [reusable workflow](../reference/github-actions/):

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

That's it. Push to main and your site deploys.
