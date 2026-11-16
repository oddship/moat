# Proposed support for Oatbase documentation

The Moat homepage's agent prompt works with the existing landing-layout hook and
focused Oatbase assets. It does not require a Moat core change.

Before moving Oatbase's documentation shell to Moat, consider these two additions.
This document proposes them; neither changes Moat's runtime or shortcode API.

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

## Configurable UI foundation assets

Allow the built-in layout's Oat stylesheet and script URLs to be replaced with
local, vendored assets. Keep the current defaults when configuration is absent.
Apply Moat's base-path handling to local URLs. Disabling the default script should
be explicit and documented because some built-in controls depend on it.

The current `head` block can add Oatbase extensions, as the homepage does. A complete
Oatbase bundle already contains Oat. Replacing the foundation through configuration
would let Oatbase's docs reuse Moat's navigation and search without loading Oat
twice or copying the entire built-in layout merely to change its asset URLs.

Acceptance: project Pages paths and custom domains both load local assets; the
configured UI script loads once; default sites retain their existing controls;
stylesheet and script attributes remain correctly escaped.

## Documentation migration boundary

Moat should own static pages, navigation, search, Markdown exports, and publication.
Oatbase should own the component registry, examples, styles, and enhancement code.
A build script can generate component Markdown from the registry. Mount examples
in isolated demo documents using their declared imports, and expose their exact
HTML and JavaScript as copyable source. Preserve old hash links with a small
redirect adapter. Do not transplant the current documentation router into Moat.
