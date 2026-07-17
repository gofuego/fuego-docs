---
title: "The non-invasive render engine framework"
layout: home
cta:
  - label: "Get started"
    url: "/docs/getting-started/"
  - label: "Explore the ecosystem"
    url: "/docs/ecosystem/"
    ghost: true
---

Your repo is already full of precise, machine-readable artifacts: OpenAPI specs, database schemas, test suites, diagrams, Kubernetes manifests, decision records. Fuego renders them as one navigable site — and a rich artifact doesn't become a page, it becomes a whole **section**: an API spec expands into an index plus a routed page per tag, operation, and schema, visible to your taxonomies and stable in its URLs.

There are [two ways to use it](docs/concepts/two-modes/). **Configure a site in place**: scaffold a project, pick the matching [format modules](https://github.com/gofuego/fuego-formats), add content and a theme, and drive the engine directly:

```bash
go run github.com/gofuego/fuego/cmd/fuego@latest init mysite --formats markdown,openapi,dbml,mermaid
cd mysite && go run . serve
```

Or **compile a [specialized engine](docs/concepts/specialized-engine/)**: model a repo's structure and artifacts as a [format pack](docs/concepts/format-packs/), and ship a purpose-built binary that renders *any* matching repository **non-invasively** — nothing written into the target repo, always up to date because the committed artifacts are the only input. See it live: [an application repo rendered by fuego-systheme](https://gofuego.github.io/demo-fuego-systheme/), one of the [ecosystem's](docs/ecosystem/) theme-and-demo pairs.

Either way, every installed format's **contract** — exactly what its parser emits — is materialized into your project, so you (or the coding agent you hand the theme to) build against it without reading parser source. No format is privileged: Markdown is a first-party parser you opt into, and when no module fits, a custom parser or pack is one small interface away. This very site is built with Fuego.
