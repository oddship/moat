---
title: Shortcodes
description: Reusable components inside markdown
---

# Shortcodes

Shortcodes are reusable HTML components you can call from inside markdown. They're Go templates stored in `_shortcodes/`.

## Creating a shortcode

Create `_shortcodes/note.html`:

```html
<div role="alert"{{ if .Get "type" }} data-variant="{{ .Get "type" }}"{{ end }}>
{{ .Inner }}
</div>
```

## Using shortcodes

### Block shortcodes

Wrap content between opening and closing tags:

```text
{{</* note type="warning" */>}}
This is a **warning** message. Markdown works inside.
{{</* /note */>}}
```

Result:

{{< note type="warning" >}}
This is a **warning** message. Markdown works inside.
{{< /note >}}

### Self-closing shortcodes

For shortcodes without inner content:

```text
{{</* badge text="New" type="success" /*/>}}
```

Result: {{< badge text="New" type="success" />}}

## Shortcode context

Templates receive a `ShortcodeContext` with:

| Field | Description |
|-------|-------------|
| `{{ .Inner }}` | Rendered inner content (markdown → HTML) |
| `{{ .Get "key" }}` | Get a named argument |
| `{{ .Args }}` | Map of all arguments |
| `{{ .Page }}` | Parent page's template data |

## Examples

### Note / alert

`_shortcodes/note.html`:

```html
<div role="alert"{{ if .Get "type" }} data-variant="{{ .Get "type" }}"{{ end }}>
{{ .Inner }}
</div>
```

{{< note type="info" >}}
This is an **info** note. Good for tips and context.
{{< /note >}}

{{< note type="warning" >}}
Watch out — this is a **warning**.
{{< /note >}}

{{< note type="error" >}}
Something went **wrong**. This is an error alert.
{{< /note >}}

### Collapsible details

`_shortcodes/details.html`:

```html
<details{{ if .Get "open" }} open{{ end }}>
<summary>{{ .Get "summary" }}</summary>
{{ .Inner }}
</details>
```

{{< details summary="Click to expand" >}}
This content is hidden by default. Markdown renders here too.

- Item one
- Item two
- Item three
{{< /details >}}

{{< details summary="Another section" open="true" >}}
This one starts open because of `open="true"`.
{{< /details >}}

## Processing order

1. Parse frontmatter from markdown source
2. Find and extract shortcode calls
3. Render inner content of block shortcodes as markdown
4. Execute shortcode templates with rendered inner + arguments
5. Splice shortcode output back into the document
6. Render the full document as markdown

This means shortcode output becomes part of the markdown document — you can mix shortcodes and markdown freely.

## Syntax reference

```text
Block:         {{</* name key="value" */>}}...content...{{</* /name */>}}
Self-closing:  {{</* name key="value" /*/>}}
Arguments:     key="value" pairs (quoted strings only)
```
