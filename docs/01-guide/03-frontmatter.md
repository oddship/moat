---
title: Frontmatter
description: YAML frontmatter for page metadata
---

# Frontmatter

Pages can have optional YAML frontmatter:

```yaml
---
title: Getting Started
description: How to set up your project
url: getting-started
---
```

## Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | Derived from filename | Page title, used in `<title>` and nav |
| `description` | string | — | Meta description for the page |
| `url` | string | Derived from file path | Override the URL path |

Any other fields are passed through as `{{ .Extra }}` in templates:

```yaml
---
title: Button
weight: 50
icon: box
---
```

Access extra fields with `{{ index .Extra "weight" }}` or `{{ .Extra.icon }}`.

## Title derivation

If no `title` is set, moat derives it from the filename:

1. Strip number prefix: `02-agents.md` → `agents.md`
2. Remove extension: `agents`
3. Replace hyphens with spaces: `agents`
4. Title case: `Agents`

## URL override

The `url` field lets you set a custom path:

```yaml
---
title: FAQ
url: frequently-asked-questions
---
```

This page would be served at `/frequently-asked-questions/` regardless of its filename.

Without `url`, the path is derived from the file's location with number prefixes stripped.
