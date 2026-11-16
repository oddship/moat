# Proposed support for Oatbase documentation

The Moat homepage uses the built-in Oatbase foundation, including its agent
prompt, copy behavior, and optional SVG styling.

Before moving Oatbase's documentation shell to Moat, consider literal shortcode
text support. UI foundation configuration is implemented as described below.

## Literal text in shortcodes

Expose the unrendered shortcode body as `RawInner`, while preserving `Inner` as
rendered Markdown. Add an explicit opt-in path for inserting a shortcode's HTML
output after Markdown rendering, so a generated `pre` or `textarea` cannot be
reinterpreted as a list when its text contains blank lines.

The distinction matters for prompts, shell commands, and configuration snippets.
In the homepage prototype, a multiline prompt inside nested raw HTML became an
HTML ordered list and the browser closed its `code` node before the steps. Copying
then returned only the introductory sentence. Encoding newlines as HTML entities
preserves the text with today's API, but is awkward for authors to maintain.

The exact option name should be reviewed for compatibility with existing shortcode
arguments. Do not change the rendering rules of existing shortcodes by default.

Acceptance: clipboard text matches the authored source for blank lines, lists,
backticks, quotes, HTML-like strings, and Unicode; existing rendered-inner and
wiki-link behavior remains unchanged; literal HTML escaping uses html/template.

## UI foundation assets (implemented)

Moat's built-in layout uses Oatbase 0.2.0 as its complete foundation. `[ui]`
accepts local or HTTP(S) stylesheet and script URLs and explicit script disabling.
Local asset URLs use `base_path`. The homepage needs no additional Oat or Oatbase
assets. See the Configuration guide for self-hosting instructions.

## Documentation migration boundary

Moat should own static pages, navigation, search, Markdown exports, and publication.
Oatbase should own the component registry, examples, styles, and enhancement code.
A build script can generate component Markdown from the registry. Mount examples
in isolated demo documents using their declared imports, and expose their exact
HTML and JavaScript as copyable source. Preserve old hash links with a small
redirect adapter. Do not transplant the current documentation router into Moat.
