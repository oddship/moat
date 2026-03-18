---
title: Conventions
description: Directory structure, links, and frontmatter rules
---

# Conventions

moat is convention-based. No config files needed — the directory structure is the config. For site-level settings, see [[Configuration]].

## Directory structure

```
docs/
├── _layout.html          # Optional. Overrides built-in layout.
├── _layout.wide.html     # Optional layout variant
├── _shortcodes/          # Optional. Shortcode templates.
│   └── note.html
├── _static/              # Copied to output as-is
├── config.toml           # Optional site config
├── index.md              # → /
├── quickstart.md         # → /quickstart/
├── 01-guide/
│   ├── 01-intro.md       # → /guide/intro/
│   └── 02-advanced.md    # → /guide/advanced/
└── 02-reference/
    └── api.md            # → /reference/api/
```

## Rules

- `_layout.html` is optional — provide it to override the built-in oat layout
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
```

The numbers are stripped from URLs and display names:
- `01-guide/01-getting-started.md` → `/guide/getting-started/` (title: "Getting Started")
- `02-reference/` → section title "Reference"

## Frontmatter

Pages can have optional YAML frontmatter:

```yaml
---
title: Getting Started
description: How to set up your project
url: custom-path
layout: wide
date: 2026-03-18 14:30
draft: false
---
```

| Field | Default | Description |
|-------|---------|-------------|
| `title` | From filename | Page title for `<title>` and nav |
| `description` | — | Meta description |
| `url` | From file path | Override the URL path |
| `layout` | (default) | Use a named layout variant |
| `date` | — | Page date (`YYYY-MM-DD` or `YYYY-MM-DD HH:MM` or full timestamp) |
| `draft` | `false` | Skip the page during build |
| `nav_children` | `true` | Set to `false` on a section's `index.md` to hide children from sidebar |

### Dates and drafts

- Accepted date formats: `YYYY-MM-DD`, `YYYY-MM-DD HH:MM`, `YYYY-MM-DD HH:MM:SS` (also with `T` separator)
- `draft: true` excludes a page from the output, nav, search index, and feed
- Sections with dated pages sort newest first in the sidebar
- `feed.xml` includes only dated pages

### Wiki links

moat supports wiki-style internal links:

```md
See [[Getting Started]] for the intro.
```

Links resolve by page title (case-insensitive). Unknown targets render as plain text rather than broken links.

### Extra fields

Any field not in the table above is available as `{{ .Extra }}` in templates:

```yaml
---
title: Button
weight: 50
icon: box
---
```

Access with `{{ index .Extra "weight" }}` or `{{ .Extra.icon }}`.

### Title derivation

If no `title` is set in frontmatter, moat derives it from the filename:

1. Strip number prefix: `02-agents.md` → `agents.md`
2. Remove `.md` extension
3. Replace hyphens with spaces
4. Title case: `Agents`
