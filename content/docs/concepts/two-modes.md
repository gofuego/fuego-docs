---
title: Two Modes of Use
layout: doc
nav_section: "Concepts"
nav_weight: 1
tags:
  - concepts
  - architecture
---

Fuego is a render-engine **framework**: the pipeline — discovery, parsing,
routing, taxonomies, rendering — is fixed, and you decide where the engine
lives. In practice every Fuego deployment is one of two shapes.

## Mode 1 — a site configured in place

The classic static-site shape: a project directory owns its content and drives
the engine directly.

```
mysite/
  config.yaml      # routes, taxonomies, collections, parsers
  main.go          # registers parsers, packs, hooks
  content/         # the content, written for this site
  theme/           # templates
```

`fuego init` scaffolds this shape ([Getting Started](docs/getting-started/)):
`config.yaml` carries the declarative concerns, `main.go` the code-level ones,
and `fuego build` writes the site to `build/`. Pick
[format modules](https://github.com/gofuego/fuego-formats) at scaffold time
with `--formats`, or start from a [format pack](docs/concepts/format-packs/)
with `--pack`.

**When it fits:** docs sites, blogs, any site whose content is written *for*
the site. This site is mode 1 — content in `content/`, theme from
[fuego-doctheme](https://github.com/gofuego/fuego-doctheme).

## Mode 2 — a compiled, specialized engine

Invert the ownership: the content is an existing repository that knows nothing
about Fuego. You model the repo's artifacts — which files matter, what
structure they carry, how they should render — as a
[format pack](docs/concepts/format-packs/), wrap it in a thin CLI, and point
the compiled binary at the repo:

```bash
fuego-systheme -site-name "My System" /path/to/repo build
```

The engine treats the repository as its content directory and writes the site
somewhere else. Nothing in the repo changes — no config file, no theme
directory, no cache, no generated content. Because the repo's own artifacts
are the only input, the rendered site can never drift from them: re-render on
push and it reflects exactly what's committed.

[The Specialized Engine](docs/concepts/specialized-engine/) covers the pattern
in depth; the [ecosystem page](docs/ecosystem/) lists the shipped specialized
engines — decision records, infrastructure, Claude Code workspaces, whole
systems — each with a live demo.

## Same engine, same skills

Both modes run the same pipeline with the same parsers, themes, hooks, and
config semantics — a pack built for an in-place site becomes a specialized
engine by adding a small CLI over the
[programmatic API](docs/concepts/embedding/). Choose per project:

| | In-place site | Specialized engine |
|---|---|---|
| Content | written for the site | already in the repo |
| Configuration | `config.yaml` in the project | compiled in (pack defaults + CLI flags) |
| Who runs it | the site's own `main.go` | anyone, on any matching repo |
| Repo footprint | the project owns its repo | zero — non-invasive |
