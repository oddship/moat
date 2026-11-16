---
title: Configuration
description: Site-level configuration via config.toml
---

# Configuration

Static Markdown exports and page actions are enabled by default. See
[[Agent-friendly docs]] for the `[agents]` settings and optional full bundle.

Place a `config.toml` in your docs directory for site-level settings. It's optional — everything has sensible defaults. See [[Conventions]] for directory structure and frontmatter fields.

```toml
site_name = "My Project"
base_path = "/my-project"
footer_text = "© [you](https://example.com) 2026"
disable_moat_citation = false

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
title = "Quickstart"
url = "/guide/quickstart/"

[[topnav]]
title = "Reference"
url = "/reference/github-actions/"

[[topnav_more]]
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
| `footer_text` | Footer copy rendered by the built-in layout; markdown links are allowed |
| `disable_moat_citation` | Removes the built-in `built with oddship/moat` suffix from `footer_text` |
| `search.enabled` | Enable built-in client-side search (defaults to `true`) |
| `feed.enabled` | Generate `feed.xml` (defaults to `false`) |
| `feed.link` | Absolute site URL used for RSS item links (recommended) |
| `feed.title` | Optional RSS title override |
| `[[topnav]]` | Primary links in the top navigation bar |
| `[[topnav_more]]` | Secondary links grouped under the built-in `More` dropdown |

CLI flags `--site-name` and `--base-path` override config values.

Use `--config PATH` to point to a config file outside the docs directory. By default, moat looks for `config.toml` in the source directory.

## Built-in footer

The built-in layout can render a footer line from plain config, without writing raw HTML.

```toml
footer_text = "© [oddship](https://github.com/oddship) 2026"
disable_moat_citation = false
```

That renders as:

- `© oddship 2026 · built with oddship/moat` when citation is enabled
- `© oddship 2026` when `disable_moat_citation = true`

For backward compatibility, the built-in layout still falls back to `[extra].footer` if `footer_text` is not set.

## Built-in search

moat generates a static `_search.json` file during `build` and the built-in Oatbase layout renders modal search in the top navigation automatically.

```toml
[search]
enabled = false
```

Search is enabled by default. Set `enabled = false` to skip generating `_search.json` and hide the built-in search UI. In the built-in layout, press `/` to open search.

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
- The built-in layout typically exposes the feed from the `More` dropdown when feed is enabled

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

The optional `icon` field adds a built-in SVG icon before the link title. Currently supported: `github`, `rss`.

## Top navigation links

Add primary links to the right side of the top navigation bar:

```toml
[[topnav]]
title = "Quickstart"
url = "/guide/quickstart/"

[[topnav]]
title = "Reference"
url = "/reference/github-actions/"
```

Top nav links use the same format as sidebar links (`title`, `url`, optional `icon`). In the built-in layout they appear before the search button and `More` dropdown.

## More dropdown links

Add secondary links under the built-in `More` dropdown:

```toml
[[topnav_more]]
title = "GitHub"
url = "https://github.com/you/project"
icon = "github"
```

`[[topnav_more]]` uses the same link format as `[[links]]` and `[[topnav]]`. When feed generation is enabled, the built-in layout also adds an `RSS feed` item in this dropdown.

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
- `footer` — rendered at the bottom of every page as a backward-compatible HTML fallback

## Version navigation

Use `current_version` and `[[versions]]` to add a version picker to the built-in
sidebar. The titles appear in the configured order, and `current_version`
selects the matching title. Without a current version, the picker starts with
a placeholder. Sites without `[[versions]]` keep their existing navigation.

```toml
current_version = "v2.0"

[[versions]]
title = "v2.0"
url = "/docs/v2.0/"

[[versions]]
title = "v1.0"
url = "/docs/v1.0/"
```

Version URLs must start at the site root (`/docs/v1.0/`) or use an absolute
HTTP(S) URL. They are used directly, without prepending `base_path`, so a page
under `/docs/v2.0/` can navigate to `/docs/v1.0/`. Each version link opens the
configured destination. Moat does not fetch releases or combine version builds;
the publishing workflow manages the content and output directories.

Custom layouts can use `.CurrentVersion` and `.Versions` to render the same
configuration.

## UI assets

The default is `assets = "vendored"`: Moat embeds a lightweight Oatbase foundation
and copies it into `_static/oatbase/` in the generated site. The foundation
contains pinned Oat and Oatbase's existing compatibility fixes. It does not load
extension components. All default UI assets are served locally, including Prose
and Kbd on documentation pages and the focused Copy/Agent Prompt assets on the
Moat homepage. No npm install or frontend build is required.

```toml
[ui]
assets = "vendored" # default; "cdn" selects a pinned CDN
# enabled = false # omit built-in UI asset tags and automatic vendor copies
# disable_script = true # retain CSS; omit foundation and prompt scripts
```

CDN mode loads the same files from an immutable upstream commit on jsDelivr.
The source revision and file checksums are recorded in the embedded vendor
manifest. This foundation is not in npm 0.2.0; the CDN uses the verified source
commit until the next npm release is published. It never follows `latest` or
`main`. CDN mode requires network access in visitors' browsers; builds remain
offline. Unknown asset modes fail the build. Files you provide under `_static/`
are still copied normally in every mode.

To provide your own foundation, set `stylesheet` and `script` under `[ui]`.
Local paths use `base_path`; absolute HTTP(S) URLs are used as given. Omitted or
empty URLs select the chosen mode's defaults. Focused component files still
come from the selected mode, so supply a compatible Oat foundation.

Disabling assets or scripts leaves native content readable. Enhanced dropdowns,
mobile navigation, and copy actions need JavaScript; supply your own behavior
when disabling it. Search and Markdown exports have separate `[search]` and
`[agents]` switches. Custom layouts can use `.UIStylesheet`, `.UIScript`,
`.UIAssetBase`, and `.UIEnabled`, or own their assets entirely.

Do not load Oat separately or combine the foundation with the complete Oatbase
bundle. Optional components should load only on the pages that use them.
