---
title: Layouts
description: Built-in defaults with optional customization via block inheritance
---

# Layouts

moat ships with a built-in layout that uses [oat](https://oat.ink) for styling — sidebar nav, dark mode toggle, responsive topnav. You don't need to create any layout files to get started.

To customize, create `_layout.html` in your docs directory. It completely replaces the built-in layout.

## Built-in layout

The built-in layout provides:

- oat CSS from CDN
- Sidebar navigation with collapsible sections
- Dark/light theme toggle
- Responsive topnav with sidebar toggle on mobile
- Syntax highlighting CSS (`_syntax.css`)
- `[[links]]` from config rendered above the page nav
- Footer from `[extra].footer` in config (supports HTML)
- Landing page variant for `layout: landing` pages

Run `moat init docs` to get a copy of the built-in layout you can edit.

## Custom layout

Create `_layout.html` with Go template syntax. Use `{{ block "name" . }}` to define sections that variants can override:

```html
<!DOCTYPE html>
<html>
<head>
  {{ block "title" . }}<title>{{ .Title }}</title>{{ end }}
  <link rel="stylesheet" href="https://unpkg.com/@knadh/oat/oat.min.css">
  <link rel="stylesheet" href="{{ .BasePath }}/_syntax.css">
  {{ block "head" . }}{{ end }}
</head>
<body>
  <nav>{{ .Nav }}</nav>
  <main>
    {{ block "content" . }}
    <article>{{ .Content }}</article>
    {{ end }}
  </main>
</body>
</html>
```

The `{{ block "name" . }}...{{ end }}` sections have default content that variants can replace.

## Named variants

Create `_layout.{name}.html` files that override blocks from the base. Only redefine what you need — everything else comes from the base layout.

`_layout.landing.html`:

```html
{{ define "title" }}{{ .SiteName }} — Welcome{{ end }}

{{ define "head" }}
<style>
  .hero { margin-bottom: 3rem; }
</style>
{{ end }}

{{ define "content" }}
<section class="hero">{{ .Content }}</section>
{{ end }}
```

Use it in a page's frontmatter:

```yaml
---
title: Home
layout: landing
---
```

{{< note type="info" >}}
Variant files only contain `{{ define }}` blocks. The base layout provides everything else — nav, head, scripts, etc. Change the base once, all variants inherit.
{{< /note >}}

If you don't provide a custom `_layout.html`, the built-in layout is used as the base, and you can still add variants that override its blocks.

## Template variables

| Variable | Type | Description |
|----------|------|-------------|
| `{{ .Title }}` | string | Page title |
| `{{ .Description }}` | string | Page description |
| `{{ .Content }}` | HTML | Rendered markdown content |
| `{{ .Nav }}` | HTML | Generated navigation sidebar |
| `{{ .CurrentPath }}` | string | Current page URL path |
| `{{ .SiteName }}` | string | Site name from config or CLI |
| `{{ .BasePath }}` | string | URL prefix (e.g. `/moat`) |
| `{{ .Extra }}` | map | Extra frontmatter from the page |
| `{{ .Site }}` | map | Site-level `[extra]` from config |

### Template functions

| Function | Description |
|----------|-------------|
| `safeHTML` | Renders a string as raw HTML (use for trusted config values like footer) |

## Navigation HTML

`{{ .Nav }}` outputs a `<nav>` with nested `<ul>` lists:

- `[[links]]` from config are rendered first as `<li>` items
- Top-level pages are direct `<li>` items
- Directories become collapsible `<details>` sections
- The current page gets `aria-current="page"`

## Static assets

Put CSS, images, and other static files in `_static/`:

```html
<link rel="stylesheet" href="/_static/style.css">
<img src="/_static/logo.png">
```

The `_static/` directory is copied to the output as-is during build.

## How inheritance works

Go templates use `block` and `define`:

1. The base layout uses `{{ block "content" . }}default markup{{ end }}`
2. A variant uses `{{ define "content" }}replacement markup{{ end }}`
3. moat clones the base, parses the variant into the clone — `define` overrides `block`
4. If a variant doesn't define a block, the base's default is used

This is Go's native template mechanism — no custom template engine.
