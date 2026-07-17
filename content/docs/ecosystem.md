---
title: Themes, Packs & Demos
layout: doc
nav_section: "Ecosystem"
nav_weight: 1
tags:
  - ecosystem
  - packs
---

Everything here is built on the engine this site documents. Each tool ships as
an importable [format pack](docs/concepts/format-packs/) plus a thin CLI — a
[specialized engine](docs/concepts/specialized-engine/) you can point at a
matching repo. And each has a **live demo** rendered from the repos of a
fictional company, *Acme Parcel*, so you can judge the output before
installing anything: the four demos are four repos of the same company — its
platform decisions, its infrastructure, its agent workspace, and one of its
applications.

## fuego-systheme — render a whole system

Point it at a repository and it renders the engineering artifacts the repo
already contains — OpenAPI, DBML, Playwright, Dockerfiles, Kubernetes
manifests, ADRs, Mermaid, Markdown — as one site that reads like the repo
itself: file-tree sidebar, an orienting landing page, the OpenAPI index shaped
like the yaml, the DBML schema drawn as an ERD. Registers every
[fuego-formats](https://github.com/gofuego/fuego-formats) parser at once;
nothing is written into the repository.
[Repo](https://github.com/gofuego/fuego-systheme) ·
[demo repo](https://github.com/gofuego/demo-fuego-systheme) (the "Dispatch AI"
application — Go backend, TS frontend, e2e specs, manifests, ADRs, diagrams) ·
**[live site](https://gofuego.github.io/demo-fuego-systheme/)**

## fuego-adr — decision records

Write decisions as `*.adr.md` files, point `fuego-adr` at the folder, and get
a status dashboard, a chronological timeline, a per-file "what decisions
touched this?" index, tag pages, and a page per decision — with section
completeness and supersession links validated for you.
[Repo](https://github.com/gofuego/fuego-adr) ·
[demo repo](https://github.com/gofuego/demo-fuego-adr) (nine records
exercising every status, including a mutual-supersession pair) ·
**[live site](https://gofuego.github.io/demo-fuego-adr/)**

## fuego-devops — infrastructure

Scans a repo for Dockerfiles and Kubernetes manifests — Kustomize- and
Helm-aware, so overlays and charts are documented as they would actually be
applied — and produces a per-namespace overview, an architecture diagram, and
a page per resource. Zero config, nothing written back.
[Repo](https://github.com/gofuego/fuego-devops) ·
[demo repo](https://github.com/gofuego/demo-fuego-devops) (multi-stage builds,
Kustomize overlays, a Helm chart) ·
**[live site](https://gofuego.github.io/demo-fuego-devops/)**

## fuego-dotclaude — Claude Code workspaces

Renders a `.claude/` directory — agents, skills, commands, output styles,
memory, MCP/settings config, installed plugins — as a navigable site with
whole-word cross-links between artifacts, backlinks, and faceted taxonomies.
[Repo](https://github.com/gofuego/fuego-dotclaude) ·
[its docs](https://gofuego.github.io/fuego-dotclaude/) ·
[demo repo](https://github.com/gofuego/demo-fuego-dotclaude) ·
**[live site](https://gofuego.github.io/demo-fuego-dotclaude/)**

## fuego-doctheme — the documentation theme

The shared look of the Fuego documentation sites, packaged as two packs:
`Public()` (user-facing docs — hero home, CTA buttons) and `Blueprint()`
(muted maintainer docs). Dogfooded in [mode 1](docs/concepts/two-modes/): the
site you are reading uses `Public()` and supplies only its topbar, sidebar,
and content.
[Repo](https://github.com/gofuego/fuego-doctheme)

## fuego-formats — the parser library

Not a theme but the substrate: one independently-versioned Go module per
format (OpenAPI, DBML, Mermaid, Playwright, Dockerfile, Kubernetes, ADR),
parsers only. Each ships its **contract** — a `schema.md` plus golden fixtures
stating exactly what the parser emits — which
[`fuego init --formats`](docs/cli/#init) materializes into your project and
fuego-systheme's renderers are written against.
[Repo](https://github.com/gofuego/fuego-formats)
