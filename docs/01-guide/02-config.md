---
title: Configuration
description: Site-level configuration via config.toml
---

# Configuration

Place a `config.toml` in your docs directory for site-level settings. It's optional — everything has sensible defaults. See [[Conventions]] for directory structure and frontmatter fields.

```toml
site_name = "My Project"
base_path = "/my-project"

[highlight]
light = "github"
dark  = "github-dark"

[search]
enabled = false

[feed]
enabled = true
link = "https://docs.example.com"

[[links]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"

[[topnav]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"

[extra]
tagline = "Build something great"
footer = '&copy; <a href="https://example.com">You</a>'
```

## Fields

| Field | Description |
|-------|-------------|
| `site_name` | Site name, available as `{{ .SiteName }}` in templates |
| `base_path` | URL prefix for GitHub project pages (e.g. `/my-project`) |
| `search.enabled` | Enable built-in client-side search (defaults to `true`) |
| `feed.enabled` | Generate `feed.xml` (defaults to `false`) |
| `feed.link` | Absolute site URL used for RSS item links (recommended) |
| `feed.title` | Optional RSS title override |
| `[[topnav]]` | Links in the top navigation bar (same format as `[[links]]`) |

CLI flags `--site-name` and `--base-path` override config values.

Use `--config PATH` to point to a config file outside the docs directory. By default, moat looks for `config.toml` in the source directory.

## Built-in search

moat generates a static `_search.json` file during `build` and the built-in oat layout renders a sidebar search box automatically.

```toml
[search]
enabled = false
```

Search is enabled by default. Set `enabled = false` to skip generating `_search.json` and hide the built-in search UI.

If you use a custom layout, you can ignore the generated index, consume `_search.json` yourself, or disable search entirely.

## RSS feed

moat can also generate a static `feed.xml` during `build`.

```toml
[feed]
enabled = true
link = "https://docs.example.com"
# optional
# title = "My Site Feed"
```

Feed generation is disabled by default.

Notes:
- Only pages with a valid `date` in frontmatter are included
- Items are sorted newest first
- `feed.link` should be your full site URL including base path (e.g. `https://example.com/moat`)
- If `[extra].tagline` is set, it becomes the feed description

## Sidebar links

Add links above the page navigation in the sidebar:

```toml
[[links]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"

[[links]]
title = "Discord"
url = "https://discord.gg/your-server"
```

Links appear at the top of the sidebar nav, before the auto-generated page tree.

The optional `icon` field adds a built-in SVG icon before the link title. Currently supported: `github`.

## Top navigation links

Add links to the right side of the top navigation bar:

```toml
[[topnav]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"
```

Top nav links use the same format as sidebar links (`title`, `url`, optional `icon`). They appear alongside the search button in the built-in layout.

## Syntax highlighting

moat uses [Chroma](https://github.com/alecthomas/chroma) for syntax highlighting with CSS classes, so light and dark themes work automatically.

```toml
[highlight]
light = "github"
dark  = "github-dark"
```

The build generates `_syntax.css` in the output directory. The built-in layout includes it automatically. If you use a custom layout, add:

```html
<link rel="stylesheet" href="{{ .BasePath }}/_syntax.css">
```

Dark theme styles are scoped under `[data-theme="dark"]`, so they activate when the user switches themes.

### Theme pairs

| Light | Dark | Family |
|-------|------|--------|
| `github` | `github-dark` | GitHub |
| `catppuccin-latte` | `catppuccin-mocha` | Catppuccin |
| `tokyonight-day` | `tokyonight-night` | Tokyo Night |
| `gruvbox-light` | `gruvbox` | Gruvbox |
| `rose-pine-dawn` | `rose-pine-moon` | Rosé Pine |
| `solarized-light` | `solarized-dark` | Solarized |
| `xcode` | `xcode-dark` | Xcode |
| `modus-operandi` | `modus-vivendi` | Modus |

See all 70 themes at the [Chroma style gallery](https://xyproto.github.io/splash/docs/).

### Example

```go
package main

import "fmt"

func main() {
    name := "moat"
    fmt.Printf("Hello from %s!\n", name)
}
```

```javascript
async function buildSite(src, dst) {
  const pages = await discoverPages(src);
  for (const page of pages) {
    const html = render(page);
    await fs.writeFile(outputPath(dst, page), html);
  }
}
```

## Site extras

The `[extra]` section holds arbitrary key-value pairs, available as `{{ .Site }}` in templates:

```toml
[extra]
tagline = "docs are your project's moat"
footer = '&copy; <a href="https://github.com/oddship">oddship</a>'
```

Access in templates:

```
{{ if index .Site "tagline" }}
  <p>{{ index .Site "tagline" }}</p>
{{ end }}
```

The built-in layout uses these extras automatically:
- `tagline` — appended to the site name in the landing page title
- `footer` — rendered at the bottom of every page (supports HTML via `safeHTML`)
