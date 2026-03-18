---
title: "Built-in Search"
date: 2026-03-15
description: Client-side search with static JSON index, topnav links, and Playwright e2e tests.
---

# Built-in Search

moat now generates a `_search.json` index at build time and the built-in layout includes a search dialog accessible via the `/` key or the search button in the top nav.

## What's included

- Static JSON search index generated during build
- Client-side search with title, description, and body text scoring
- Search dialog with keyboard navigation
- Top navigation links via `[[topnav]]` in [[Configuration]]

## Disable search

```toml
[search]
enabled = false
```

Search is enabled by default. See [[Configuration]] for details.
