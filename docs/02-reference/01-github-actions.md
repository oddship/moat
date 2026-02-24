---
title: GitHub Actions
description: Deploy docs to GitHub Pages with one workflow file
---

# GitHub Actions

moat provides a reusable workflow for deploying docs to GitHub Pages. One file, zero config.

## Usage

Create `.github/workflows/docs.yml` in your repo:

```yaml
name: Docs
on:
  push:
    branches: [main]
    paths: ['docs/**']
  workflow_dispatch:

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

Push to main. Your site deploys to GitHub Pages.

## Setup

1. Go to your repo's **Settings → Pages**
2. Set **Source** to **GitHub Actions**
3. Add the workflow file above
4. Push

## Inputs

| Input | Default | Description |
|-------|---------|-------------|
| `docs_dir` | `docs` | Path to your docs source directory |
| `output_dir` | `_site` | Build output directory |
| `moat_version` | `latest` | moat release version (e.g. `v0.2.0`) |
| `go_version` | `1.24` | Go version for building moat |
| `from_source` | `false` | Build moat from source instead of `go install` |

## Example with all options

```yaml
jobs:
  docs:
    uses: oddship/moat/.github/workflows/build-docs.yml@main
    with:
      docs_dir: website
      output_dir: dist
      moat_version: v0.2.0
```

## What it does

The workflow:

1. Checks out your repo
2. Installs Go and builds moat from source
3. Runs `moat build` with your docs directory
4. Reads `config.toml` from the docs directory for site name, base path, etc.
5. Uploads the built site as a GitHub Pages artifact
6. Deploys to GitHub Pages

{{< note type="info" >}}
The workflow builds moat from source rather than downloading a binary. This means you always get the latest features on the `@main` branch, and you can pin to a specific commit or tag.
{{< /note >}}

## Custom workflow

If you need more control, build your own:

```yaml
name: Docs
on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go install github.com/oddship/moat@latest
      - run: moat build docs/ _site/
      - uses: actions/upload-pages-artifact@v3
        with:
          path: _site/

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/deploy-pages@v4
        id: deployment
```
