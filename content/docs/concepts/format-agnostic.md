---
title: Why Format-Agnostic?
layout: doc
nav_section: "Concepts"
nav_weight: 3
tags:
  - concepts
  - architecture
---

Most static site generators treat Markdown as the default content format. Hugo, Jekyll, Eleventy — they all parse Markdown first and offer limited escape hatches for other formats. Fuego inverts this: **no format is privileged**.

## Your content is already structured

Look at what a repository actually contains: an OpenAPI spec, a DBML schema, Kubernetes manifests, decision records, Playwright suites, Mermaid diagrams. Each is precise, structured, and machine-readable — an operation has a method, a path, parameters, responses; a table has columns, keys, refs. None of it is prose.

A Markdown-first generator can only publish that content as a hand-written *description* of it — a second copy that starts drifting the moment it's written. A format-agnostic engine renders the artifact itself: the parser understands the format's real structure, and the site can never disagree with the source. That's the foundation of the [specialized engine](docs/concepts/specialized-engine/) pattern — and of docs that stay correct because they're rendered, not rewritten.

## Fuego's Approach

In Fuego, the content format is a first-class concept:

1. **You define the syntax** — via regex rules in config or a Go parser
2. **You define the semantics** — node types like `operation`, `table`, `front` carry meaning
3. **You define the rendering** — template per node type, or CSS on `data-type` attributes

The engine never interprets what an `operation` node means. It just knows how to discover files, dispatch them to parsers, route the results, and render templates. Your format, your rules.

## The Universal AST

All parsers — Markdown, declarative regex, compiled Go — produce the same output:

```go
type Node struct {
    Type       string            // "operation", "front", "html", anything
    Attributes map[string]any    // structured metadata
    Content    string            // text content
    Children   []Node            // nested nodes
    Raw        bool              // emit Content as raw HTML, unwrapped
}
```

This uniformity means the entire pipeline after PARSE is format-agnostic. Routing, taxonomy indexing, collections, rendering, manifest generation — all work the same regardless of whether the content started as `.md`, `.openapi.yaml`, or `.dbml`.

And a file doesn't have to become exactly one page: a **tree parser** returns a whole tree of pages from one artifact, so an API spec expands into an index plus a routed page per tag, operation, and schema — each visible to taxonomies and collections like any other page. See [Custom Parsers](docs/custom-parsers/).

## Markdown is Still First-Class

Format-agnostic doesn't mean anti-Markdown. Fuego ships a first-party Markdown parser with full GFM support (tables, strikethrough, autolinks, task lists) — you opt in by registering it (`eng.Register(markdown.Parser())`), the same way you'd register any parser. It emits a single `Raw: true` node that the renderer passes through as HTML.

You can mix formats freely in one site: Markdown for the prose, `.openapi.yaml` for the API section, `.adr.md` for decisions. Each file dispatches to the parser that claims it — by filename pattern or bare extension, most specific claim first — and the pipeline handles the rest.

## Custom DSLs, too

The same machinery covers content you author yourself. A flashcard has a front and a back; a quiz question has a prompt, choices, and a correct answer — structure Markdown can only fake with conventions. Define a small DSL instead (`.card`, `.trivia`), declaratively [in config](docs/custom-parsers/) or as a compiled parser, and author in a syntax that matches the content. The [quiz-site tutorial](docs/tutorials/build-a-quiz-site/) builds one end to end.

## When This Matters

Format-agnostic design shines when:

- **Your repo already holds structured artifacts** — render the source of truth instead of describing it, in your own site or [non-invasively](docs/concepts/specialized-engine/) across any matching repo
- **Your content has inherent structure** that Markdown can't express — define a DSL for it
- **Client-side interactivity** needs typed data (the JSON embed carries the full AST)
- **Multiple content types coexist** in one site with different rendering needs

If all your content is prose, register the first-party Markdown parser and you have a normal Markdown site. The framework is there for when you need more.
