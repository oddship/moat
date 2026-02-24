---
title: Configuration
description: Site-level configuration via config.toml
---

# Configuration

Place a `config.toml` in your docs directory for site-level settings. It's optional — everything can be set via CLI flags too.

```toml
site_name = "My Project"
base_path = "/my-project"

[highlight]
light = "github"
dark  = "github-dark"

[extra]
tagline = "Build something great"
repo = "https://github.com/you/project"
```

## Fields

| Field | Description |
|-------|-------------|
| `site_name` | Site name, available as `{{ .SiteName }}` in templates |
| `base_path` | URL prefix for GitHub project pages (e.g. `/my-project`) |

CLI flags `--site-name` and `--base-path` override config values.

Use `--config PATH` to point to a config file outside the docs directory. By default, moat looks for `config.toml` in the source directory.

## Syntax highlighting

moat uses [Chroma](https://github.com/alecthomas/chroma) for syntax highlighting with CSS classes, so light and dark themes work automatically.

```toml
[highlight]
light = "github"
dark  = "github-dark"
```

The build generates `_syntax.css` in the output directory. Include it in your layout:

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

Here's how highlighted code looks across languages:

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

```bash
# Build and preview
moat build docs/ _site/
moat serve _site/ --port 8080
```

## Site extras

The `[extra]` section holds arbitrary key-value pairs, available as `{{ .Site }}` in templates:

```toml
[extra]
tagline = "docs are your project's moat"
repo = "https://github.com/oddship/moat"
```

Access in your layout:

```html
{{ if .Site.tagline }}
  <p>{{ index .Site "tagline" }}</p>
{{ end }}
```
