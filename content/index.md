---
title: "The meta-engine for static sites"
layout: home
cta:
  - label: "Get started"
    url: "/docs/getting-started/"
  - label: "Bring your own format"
    url: "/docs/custom-parsers/"
    ghost: true
---

Most static site generators bake in Markdown. Fuego is **format-agnostic**: a parser turns *any* file into pages, and the engine handles the rest — discovery, routing, taxonomies, pagination, feeds, and a dev server.

Your repo is already full of precise, machine-readable artifacts: OpenAPI specs, database schemas, test suites, diagrams, Kubernetes manifests, decision records. Pick the matching [format modules](https://github.com/gofuego/fuego-formats) and Fuego renders them as one navigable site. A rich artifact doesn't become a page — it becomes a whole **section**: an API spec expands into an index plus a routed page per tag, operation, and schema, visible to your taxonomies and stable in its URLs.

```bash
go run github.com/gofuego/fuego/cmd/fuego@latest init mysite --formats markdown,openapi,dbml,mermaid
cd mysite && go run . serve
```

Every installed format's **contract** — exactly what its parser emits — is materialized into your project, so you (or the coding agent you hand the theme to) build against it without reading parser source. Markdown is a first-party parser you opt into, and when no module fits, a custom parser or [format pack](docs/concepts/format-packs/) is one small interface away. This very site is built with Fuego.
