---
title: Development
description: Contributing to moat — build, test, and project structure
---

# Development

Guide for contributing to moat or understanding the codebase.

## Prerequisites

- **Go 1.24+** — `go version`
- **Node.js 18+** — only for e2e tests, not required for building moat itself
- **Chromium** with `--remote-debugging-port=9222` — for running Playwright e2e tests

## Build from source

```bash
git clone https://github.com/oddship/moat.git
cd moat
go build -o moat .
./moat version
```

## Run tests

### Go unit tests

```bash
go test ./...
```

Covers search index generation, config parsing, and build output. No browser needed.

### Playwright e2e tests

The e2e tests connect to an existing Chromium browser via CDP — they don't launch their own. This means they test against the real browser the user is running.

```bash
# 1. Start Chromium with remote debugging (in a separate terminal)
chromium --remote-debugging-port=9222

# 2. Build the docs site without base_path for local testing
moat build docs/ _site/ --site-name moat

# 3. Serve locally
moat serve _site/

# 4. Install test dependencies and run (in another terminal)
npm install
npx playwright test --workers=1
```

Tests cover layout structure, search dialog, keyboard shortcuts, theme toggle, and accessibility attributes.

## Project structure

```
moat/
├── main.go            # CLI entrypoint (build, serve, init, version)
├── build.go           # Build pipeline: walk → parse → render → write
├── config.go          # Config types and TOML parsing
├── search.go          # Search index generation from rendered HTML
├── nav.go             # Sidebar navigation tree + HTML rendering
├── layouts.go         # Template loading (built-in + custom)
├── markdown.go        # Goldmark markdown → HTML rendering
├── frontmatter.go     # YAML frontmatter parsing
├── shortcodes.go      # Shortcode template processing
├── defaults.go        # Title/filename conventions (strip prefixes)
├── serve.go           # Simple static file server
├── search_test.go     # Go unit tests
├── embed/             # Built-in templates (embedded via go:embed)
│   ├── _layout.html         # Base layout (oat sidebar + topnav)
│   ├── _layout.landing.html # Landing page variant
│   └── config.toml          # Default config scaffold
├── e2e/               # Playwright e2e tests
│   ├── fixtures.js    # CDP connection fixture
│   ├── layout.spec.js
│   ├── search.spec.js
│   └── theme.spec.js
├── docs/              # moat's own documentation (built with moat)
└── .github/workflows/ # CI: build docs, deploy to Pages, releases
```

## Build pipeline

The `Build()` function in `build.go` runs this pipeline:

1. **Walk** source directory for `.md` files (skip `_` and `.` prefixed paths)
2. **Parse** frontmatter and store raw markdown body on each `Page`
3. **Build** navigation tree from directory structure
4. **Generate** syntax highlighting CSS (light + dark themes via Chroma)
5. **Render** each page:
   - Process shortcodes (expand `{{</* name */>}}` templates)
   - Render markdown to HTML via Goldmark
   - Store rendered HTML on `Page.HTML`
   - Execute layout template with `TemplateData`
   - Write output HTML file
6. **Generate** search index from rendered HTML (strip tags, cap at 2000 chars)
7. **Copy** `_static/` directory as-is

## Search indexing

The search index (`_search.json`) is built **after** the render loop so it indexes the final rendered HTML, including shortcode output. `extractSearchText()` strips HTML tags with a single regex and caps text at 2000 characters per entry.

The built-in layout includes inline JavaScript that:
- Lazily fetches the search index on first keystroke
- Scores matches: title (100 pts) > description (25 pts) > body (10 pts)
- Renders top 8 results in a `<dialog>` modal
- Supports `/` keyboard shortcut to open search

## Layout system

Layouts are Go `html/template` files. The built-in layout is embedded via `go:embed` and used when no `_layout.html` exists in the source directory.

Custom layouts override the built-in by providing `_layout.html`. Named variants (`_layout.wide.html`) are selected via `layout: wide` in frontmatter.

Template functions available in layouts:

| Function | Description |
|----------|-------------|
| `safeHTML` | Render a string as raw HTML |
| `linkIcon` | Return built-in SVG icon by name (e.g. `"github"`) |

## Adding a built-in icon

Icons for `[[links]]` and `[[topnav]]` config are stored in `nav.go`:

```go
var builtinIcons = map[string]string{
    "github": `<svg ...>...</svg>`,
}
```

To add a new icon, add an entry to the `builtinIcons` map with the SVG markup. Keep SVGs at 16×16, using `fill="currentColor"` and `aria-hidden="true"`.

## Style guidelines

- **No inline styles in the built-in layout.** Use oat utility classes (`hstack`, `gap-2`, `justify-end`, `mt-6`, etc.) or component patterns (`role="search"`, `data-topnav`, etc.).
- **No new CSS framework or build pipeline.** The layout uses oat via CDN.
- **No new Go dependencies for HTML processing.** Use regex-based extraction for search text.
- **Conventional commits** for commit messages (`feat:`, `fix:`, `test:`, `docs:`).

## Releasing

Releases are automated via GitHub Actions. Push a `v*` tag to trigger:

```bash
git tag v0.3.0
git push origin v0.3.0
```

This builds binaries for linux/darwin × amd64/arm64 and creates a GitHub Release.
