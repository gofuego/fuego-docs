---
title: Templates
layout: doc
nav_section: "Reference"
nav_weight: 3
tags:
  - reference
  - templates
---

Templates use Go's `html/template`. The theme directory structure:

```
theme/
  base.html              # HTML shell (required)
  layouts/
    post.html            # named layouts
    tag.html
  renderers/
    question.html        # per-node-type renderers
  partials/
    nav.html             # reusable snippets, called via {{partial "nav" .}}
```

## Base Template

The base template is required. It defines the HTML shell:

```html
<!DOCTYPE html>
<html>
<head><title>{{.Page.Envelope.title}} | {{.Site.Name}}</title></head>
<body>
{{block "content" .}}<div id="root">{{.Page.Content}}</div>{{end}}
<script type="application/json" id="fuego-data">{{.JSON}}</script>
</body>
</html>
```

## Layout Templates

Layouts override the `"content"` block defined in the base template:

```html
{{define "content"}}
<article class="post">{{.Page.Content}}</article>
{{end}}
```

Set a page's layout via frontmatter (`layout: post`) or config (taxonomies, collections).

## Renderer Templates

Per-node-type renderer templates override the default `<div data-type="...">` rendering. Place them in `theme/renderers/{type}.html`.

For example, `theme/renderers/question.html` controls how `question` nodes render.

## Partials

Partials are reusable template snippets in `theme/partials/`. They are available in the base template, layouts, and renderers via the `partial` function:

```html
{{partial "nav" .}}
{{partial "footer" (dict "year" "2026" "owner" .Site.Name)}}
```

The second argument becomes the partial's data (`.`); omit it to pass nothing. Partials can call other partials. Calling a partial that does not exist fails the page render with an error listing the available partials.

## Template Functions

| Function | Example | Description |
|---|---|---|
| `render` | `{{render .Children}}` | Recursively render a `[]Node` through renderer templates (used inside a renderer template to render child nodes; the page body is already in `.Page.Content`) |
| `safeHTML` | `{{safeHTML .Page.Envelope.snippet}}` | Mark a string as raw HTML (skips escaping) |
| `partial` | `{{partial "nav" .}}` | Execute a template from `theme/partials/` |
| `dict` | `{{dict "k1" "v1" "k2" "v2"}}` | Build a map from key/value pairs, for partial arguments |
| `where` | `{{where .Site.Pages "type" "post"}}` | Filter a slice by key = value |
| `sortBy` | `{{sortBy .Site.Pages "weight" "desc"}}` | Copy-and-sort a slice by key; optional `"asc"`/`"desc"` |
| `limit` | `{{where .Site.Pages "type" "post" \| limit 5}}` | At most n leading elements (collection last, pipeline-friendly) |
| `first` | `{{first .Site.Pages}}` | First element, or nil if empty |
| `dateFormat` | `{{dateFormat "Jan 2, 2006" .Page.Envelope.date}}` | Format a `time.Time` or date string with a Go layout |

Key resolution for `where` and `sortBy`: map keys for maps, exported struct fields (case-insensitive), then the element's `Envelope` map — so `"type"` matches a page ref's `Type` field and `"author"` matches `Envelope["author"]`. `where` compares loosely (values match if their string forms are equal), so YAML numbers and template string literals compare predictably. `sortBy` sorts numbers numerically and everything else by string form; elements missing the key sort first.

`dateFormat` accepts `time.Time` or strings in RFC3339, `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, or `2006-01-02`.

## Template Data

| Field | Description |
|---|---|
| `.Page.Envelope` | Frontmatter map |
| `.Page.Content` | Pre-rendered HTML |
| `.Page.URL` | Resolved URL |
| `.Page.Layout` | Layout name |
| `.Page.Type` | Parser type |
| `.Site.Name` | From config |
| `.Site.BaseURL` | From config |
| `.Site.Pages` | All pages as slim refs (`.URL`, `.Type`, `.Layout`, `.Envelope`), URL-sorted — see [Build a Navigation Menu](docs/how-to/build-a-nav-menu/) |
| `.Paginator` | On paginated listing pages: `.CurrentPage`, `.TotalPages`, `.PrevURL`, `.NextURL` (nil otherwise) — see [Paginate a Collection](docs/how-to/paginate-a-collection/) |
| `.JSON` | Full page data as JSON (computed only if the template references it) |

## Linking

Fuego does **not** rewrite links to account for the site's `base_url`, so links
must be base-aware or they break when the site is deployed under a subpath (e.g.
`/owner/repo` on GitHub Pages). Page URLs (`.URL`, `.Paginator.PrevURL`, a ref's
`.URL`) are root-relative (they start with `/`), so in a template you **prefix
them with the base**:

```html
<a href="{{.Site.BaseURL}}{{.URL}}">{{.Envelope.title}}</a>   <!-- ✓ -->
<a href="{{.URL}}">…</a>                                       <!-- ✗ escapes the base -->
```

Inside a `range`, `.` is the element, so reach the site root with `$`:
`{{$.Site.BaseURL}}{{.URL}}`. The same applies to assets:
`href="{{.Site.BaseURL}}/style.css"`.

In **content** (Markdown and other formats), where you can't call `.Site.BaseURL`,
use **base-relative** links — no leading slash:

```md
[Guide](docs/guide/)     <!-- ✓ resolves against <base href> = the site root -->
[Guide](/docs/guide/)    <!-- ✗ absolute path; escapes the deploy subpath -->
```

This works because `theme/base.html` should set `<base href="{{.Site.BaseURL}}/">`,
against which a no-leading-slash link resolves from the site root under any
`base_url`. Catch violations with [`build --strict-links`](docs/cli/#build).

## JSON Embed

Embedding `{{.JSON}}` in your base template or a layout gives the page its full data (envelope, parsed nodes, URL) as JSON, typically inside `<script type="application/json" id="fuego-data">`. This enables client-side interactivity without a separate API.

The payload is computed only for pages whose template actually references `.JSON` — Fuego inspects templates at load time, so sites that never hydrate pay no serialization cost. There is no configuration for this; using `.JSON` anywhere in a template (including inside `if`/`range` blocks or via `$.JSON`) enables it for pages rendered with that template.
