---
title: "v0.6.0 — Site Primitives"
date: 2026-03-18
description: Wiki links, date/draft frontmatter, page listings, RSS feed, and date-aware nav sorting.
---

# v0.6.0 — Site Primitives

New primitives for building richer content sites — blogs, knowledge bases, changelogs — without a separate mode.

## Wiki links

Link between pages using `[[Page Title]]` syntax:

```md
See [[Configuration]] for all options.
```

Links resolve by page title (case-insensitive). Unknown targets render as plain text.

## Date and draft frontmatter

Pages can now have `date` and `draft` fields:

```yaml
---
title: My Post
date: 2026-03-18
draft: true
---
```

Draft pages are excluded from the build entirely — no output, no nav, no search index, no feed.

## RSS feed

Enable in `config.toml`:

```toml
[feed]
enabled = true
link = "https://your-site.com"
```

Only pages with `date` are included, sorted newest first. See [[Configuration]] for all feed options.

## Page listings

Templates and shortcodes can access all pages via `{{ .Pages }}` and `{{ .SectionPages "section" }}`. See [[Layouts]] and [[Shortcodes]] for details.

## Date-aware nav sorting

Sections where pages have `date` frontmatter automatically sort newest first in the sidebar. No configuration needed.
