---
title: The Specialized Engine
layout: doc
nav_section: "Concepts"
nav_weight: 2
tags:
  - concepts
  - architecture
  - packs
---

A **specialized engine** is a compiled binary that renders a repository it
doesn't own. You model a class of repos — the artifacts they contain and how
those artifacts should read as a site — into a
[format pack](docs/concepts/format-packs/), and a thin CLI turns that pack
into a tool anyone can point at a matching repo. This is
[mode 2](docs/concepts/two-modes/) of using Fuego, and it's how every shipped
Fuego tool — [fuego-adr](https://github.com/gofuego/fuego-adr),
[fuego-devops](https://github.com/gofuego/fuego-devops),
[fuego-dotclaude](https://github.com/gofuego/fuego-dotclaude),
[fuego-systheme](https://github.com/gofuego/fuego-systheme) — is built.

## The pattern

Three pieces: a pack, the programmatic API, and a content directory that
points at the target repository.

```go
eng := engine.New()
eng.Use(systheme.Pack()) // parsers + theme + hooks + config defaults

err := eng.Build(ctx, engine.BuildOptions{
    ContentDir: repoPath, // the target repository — read, never written
    OutputDir:  outDir,   // the site goes elsewhere
    SiteName:   siteName,
    BaseURL:    baseURL,
})
```

`ContentDir` doesn't have to be your project's content folder — it can be any
directory, including the root of somebody else's repository. Discovery walks
it, the pack's parsers claim the files they recognize by name
(`*.openapi.yaml`, `Dockerfile`, `*.adr.md`, …), and everything unclaimed is
simply not content. The full `BuildOptions` surface is on
[Embedding Fuego](docs/concepts/embedding/).

## Non-invasive by construction

The engine reads `ContentDir` and writes only to `OutputDir` (plus, for
incremental builds, `CacheDir`). Point both away from the target repo and the
render leaves it untouched: no config file, no theme directory, no cache, no
generated content. fuego-systheme's dev server, for example, keeps its build
cache in a temp directory so even `serve` writes nothing into the scanned
repo.

Because the repo's committed artifacts are the only input, the site cannot
drift from them. Rebuild on every push — the
[reusable deploy workflows](docs/how-to/deploy-with-reusable-workflows/) do —
and the documentation is **always up to date**: there is no second copy to
forget to edit.

## Worked example: fuego-systheme

[fuego-systheme](https://github.com/gofuego/fuego-systheme) is the *system
theme*: it registers every [fuego-formats](https://github.com/gofuego/fuego-formats)
parser (plus the engine's Markdown parser) and renders a repository's
engineering artifacts as one site that reads like the repo itself:

| Artifact | Becomes |
|---|---|
| OpenAPI 3 specs | an API section: index, a page per tag, operation, and schema |
| DBML schemas | a schema index plus a page per table |
| Playwright suites | a page per spec, suite, and test |
| Dockerfiles | a page per build, stage by stage |
| Kubernetes manifests | a page per resource, plus a by-kind taxonomy |
| ADRs | a page per decision, with status and supersession |
| Mermaid diagrams | a page per diagram |
| Markdown docs | prose pages, repo-relative links rewritten to rendered pages |

The output is **visualization-oriented** — the site shows the system rather
than describing it: a repo-structure sidebar (the file tree, artifacts as
links) on every page; a landing page that orients (root modules, the decisions
currently in force, the rendered README); the OpenAPI index shaped like the
yaml devs already scan; and the DBML schema drawn as an entity-relationship
diagram generated from the tables, keys, and refs.

Two engine features carry the pattern. A rich artifact becomes a whole
**section** of routed pages — an API spec expands into an index plus a page
per tag, operation, and schema — via tree parsers
([Custom Parsers](docs/custom-parsers/)). And `BeforeRender` hooks rewrite the
repo's own Markdown links (written for the GitHub view) to their rendered
pages, which is what makes [`--strict-links`](docs/how-to/check-for-broken-links/)
viable over an unmodified repo.

**See it live:** [demo-fuego-systheme](https://github.com/gofuego/demo-fuego-systheme)
— the application repo of a fictional AI service, zero Fuego code committed,
rendered at
[gofuego.github.io/demo-fuego-systheme](https://gofuego.github.io/demo-fuego-systheme/).

## Build your own

1. Put the domain in a pack — parsers, theme, hooks, config defaults
   ([tutorial](docs/tutorials/build-a-format-pack/)).
2. Wrap it in a CLI that maps flags to `engine.BuildOptions` and calls
   `eng.Build` / `eng.Serve` / `eng.Validate`
   ([Embedding Fuego](docs/concepts/embedding/)).
3. Study the shipped engines and their demos on the
   [ecosystem page](docs/ecosystem/).
