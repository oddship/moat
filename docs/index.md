---
title: moat
layout: landing
description: A static site generator in one Go binary. Markdown + oat.
---

# moat

<h3 class="tagline text-light">Markdown + <a href="https://oat.ink">oat</a>. A static site generator in one Go binary.</h3>

<div class="hstack">
  <a href="guide/quickstart/" class="button">Get started</a>
  <a href="https://github.com/oddship/moat" class="button outline">GitHub</a>
</div>

<br>

<div class="features">
<article class="card">
<header><h3>Zero config</h3></header>

Built-in oat layout with sidebar nav, dark mode, and syntax highlighting.
Just write markdown and build. No layout files needed.
</article>

<article class="card">
<header><h3>One binary</h3></header>

Single Go binary. No Node.js, no npm, no Ruby.
`go install` or grab a release binary.
</article>

<article class="card">
<header><h3>Syntax highlighting</h3></header>

70 Chroma themes with automatic light/dark mode.
Configure in `config.toml` — CSS classes, not inline styles.
</article>

<article class="card">
<header><h3>Layouts & shortcodes</h3></header>

Multiple layouts via Go template `block`/`define`. Shortcodes for reusable
components in markdown. Like Zola, without the complexity.
</article>
</div>

<ot-copy class="card" data-agent-prompt data-copy-target="moat-migration-prompt" data-copied="Prompt copied">
  <section class="vstack gap-4">
    <div class="hstack"><span class="badge outline">Agent quickstart</span></div>
    <div class="hstack"><h2>Bring your docs to Moat</h2></div>
    <p>Give your coding agent a migration plan. Review the local result before publishing.</p>
    <div class="hstack"><button type="button" data-copy-button><svg data-icon aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" > <path d="M7 9.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667l0 -8.666" /> <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1" /> </svg><span data-copy-label>Copy agent prompt</span></button></div>
    <small class="text-light">Works with any agent that can read your repository.</small>
  </section>
  <section data-code-block data-wrap aria-label="Migration prompt">
    <header><strong>Prompt</strong><span class="badge outline">Moat + GitHub Pages</span></header>
    <pre tabindex="0" aria-label="Selectable migration prompt"><code id="moat-migration-prompt" data-copy-source>Migrate this repository&#x27;s documentation to Moat and publish it with GitHub Pages.&#10;&#10;1. Read the repository instructions and inspect the current docs, assets, URLs, build, and deployment workflow. Read https://oddship.github.io/moat/llms.txt and follow the linked Moat documentation.&#10;2. Preserve the existing prose, links, images, and stable public URLs. Use Moat&#x27;s built-in layout first. Explain any URL changes or unsupported features before applying them.&#10;3. Install a pinned Moat release or use the project&#x27;s existing version. Put site configuration in the docs directory and set base_path correctly for this repository or its custom domain.&#10;4. Build the site locally. Check navigation, search, Markdown links, assets, code blocks, and narrow layouts. Show me a local preview and summarize the diff.&#10;5. Prepare a GitHub Pages Actions workflow using Moat&#x27;s documented deployment workflow, pinned to a reviewed commit. Do not replace unrelated workflows or expose private content.&#10;6. Use the existing GitHub authentication. Before enabling Pages, changing repository settings, pushing, or publishing, show the exact changes and ask for approval unless I have already authorized those actions.&#10;7. After publication is authorized, wait for the deployment matching the pushed commit and verify the live homepage, a documentation page, and an asset. Report the exact deployment state.</code></pre>
  </section>
</ot-copy>

## Quick start

```bash
go install github.com/oddship/moat@latest
mkdir docs
echo "# Hello" > docs/index.md
moat build docs/ _site/
moat serve _site/
```

## What it does

Reads markdown files from a directory, converts to HTML, wraps in a layout template, generates sidebar nav, writes static HTML. That's it.

```
docs/
├── config.toml           # Site config (optional)
├── _layout.html          # Custom layout (optional)
├── _shortcodes/          # Reusable components
├── _static/              # Copied as-is
├── index.md              # → /
└── 01-guide/
    ├── 01-intro.md       # → /guide/intro/
    └── 02-config.md      # → /guide/config/
```

Ships with a built-in [oat](https://oat.ink) layout — sidebar, dark mode, semantic HTML. Or provide your own `_layout.html` with any CSS you want.
