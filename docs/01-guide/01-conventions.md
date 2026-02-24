---
title: Conventions
description: How moat maps files to pages
---

# Conventions

moat is convention-based. No config files — the directory structure is the config.

## Directory structure

```
docs/
├── _layout.html          # Required. Go template.
├── _static/              # Copied to output as-is
├── index.md              # → /
├── quickstart.md         # → /quickstart/
├── guide/
│   ├── 01-intro.md       # → /guide/intro/
│   └── 02-advanced.md    # → /guide/advanced/
└── reference/
    └── api.md            # → /reference/api/
```

## Rules

- `_layout.html` is required — it's the Go template that wraps every page
- `_static/` is copied to the output directory as-is (CSS, images, etc.)
- Files and directories prefixed with `_` or `.` are skipped
- `index.md` at any level becomes the directory's root page
- All other `.md` files get clean URLs: `file.md` → `/file/`

## Number prefixes

Prefix files and directories with `01-`, `02-`, etc. to control ordering:

```
01-guide/
  01-getting-started.md
  02-configuration.md
  03-deployment.md
02-reference/
  01-api.md
  02-cli.md
```

The prefixes are stripped from both URLs and display names:

| File | URL | Nav label |
|------|-----|-----------|
| `01-guide/01-getting-started.md` | `/guide/getting-started/` | Getting Started |
| `02-reference/01-api.md` | `/reference/api/` | Api |

## Clean URLs

Every page gets a clean URL by writing to `path/index.html`:

| Source | Output | URL |
|--------|--------|-----|
| `index.md` | `index.html` | `/` |
| `quickstart.md` | `quickstart/index.html` | `/quickstart/` |
| `guide/01-intro.md` | `guide/intro/index.html` | `/guide/intro/` |
