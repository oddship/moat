---
title: Layout Template
description: The Go template that wraps every page
---

# Layout template

`_layout.html` is a Go template. moat passes these variables to it:

## Template variables

| Variable | Type | Description |
|----------|------|-------------|
| `{{ .Title }}` | string | Page title |
| `{{ .Description }}` | string | Page description (from frontmatter) |
| `{{ .Content }}` | HTML | Rendered markdown content |
| `{{ .Nav }}` | HTML | Generated navigation sidebar |
| `{{ .CurrentPath }}` | string | Current page URL path (e.g. `/guide/intro/`) |
| `{{ .SiteName }}` | string | Site name (from `--site-name` flag, default "Site") |

## Minimal example

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Title }} — {{ .SiteName }}</title>
</head>
<body>
  <nav>{{ .Nav }}</nav>
  <main>{{ .Content }}</main>
</body>
</html>
```

## With oat

moat pairs well with [oat CSS](https://github.com/knadh/oat):

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Title }} — {{ .SiteName }}</title>
  <link rel="stylesheet" href="https://unpkg.com/@knadh/oat/oat.min.css">
  <script src="https://unpkg.com/@knadh/oat/oat.min.js" defer></script>
</head>
<body data-sidebar-layout>
  <nav data-topnav>
    <a href="/">{{ .SiteName }}</a>
  </nav>
  <aside data-sidebar>
    {{ .Nav }}
  </aside>
  <main class="wrap">
    <div class="container">
      <article>{{ .Content }}</article>
    </div>
  </main>
</body>
</html>
```

## Navigation HTML

The `{{ .Nav }}` variable outputs a `<nav>` with nested `<ul>` lists:

- Top-level pages are direct `<li>` items
- Directories become `<details>` with a `<summary>` (collapsible sections)
- The current page gets `aria-current="page"`
- Active sections are automatically opened

## Static assets

Put CSS, images, and other static files in `_static/`:

```
docs/
├── _layout.html
├── _static/
│   ├── style.css
│   └── logo.png
└── index.md
```

Reference them with absolute paths in your layout:

```html
<link rel="stylesheet" href="/_static/style.css">
<img src="/_static/logo.png">
```
