---
title: "Tutorial: Render a Repo's Artifacts"
layout: doc
nav_section: "Tutorials"
nav_weight: 1
order: 1
tags:
  - tutorial
  - formats
---

Your repository already contains precise, machine-readable artifacts — an
OpenAPI spec, a database schema, decision records, test suites, manifests,
diagrams. This tutorial turns them into a navigable site two ways: the
instant way (a ready-made [specialized engine](docs/concepts/specialized-engine/))
and the owned way (your own project with the matching format modules).

## The instant way: fuego-systheme

If the artifacts follow common naming — `*.openapi.yaml`, `*.dbml`,
`Dockerfile`, `*.k8s.yaml`, `*.adr.md`, `*.mmd`, `*.spec.ts`, `*.md` —
[fuego-systheme](https://github.com/gofuego/fuego-systheme) renders them with
no configuration at all:

```bash
go run github.com/gofuego/fuego-systheme/cmd/fuego-systheme@latest \
  -site-name "My System" /path/to/repo serve
```

Open the dev server and you get a dashboard counting every artifact family, a
file-tree sidebar on every page, the OpenAPI spec expanded into an API section
(index plus a page per tag, operation, and schema), the DBML schema drawn as
an ERD, decisions on a status board — and **nothing written into the repo**:
no config, no theme, no cache.

For CI, build with a deploy subpath and strict link checking:

```bash
go run github.com/gofuego/fuego-systheme/cmd/fuego-systheme@latest \
  -site-name "My System" -base-url /my-repo -strict-links -output build . build
```

No repo handy? Clone
[demo-fuego-systheme](https://github.com/gofuego/demo-fuego-systheme) — an
application repo exercising every format at once — and run the command
against it; the result is deployed at
[gofuego.github.io/demo-fuego-systheme](https://gofuego.github.io/demo-fuego-systheme/).

## The owned way: your own project

When you want your own theme, routes, and taxonomies, scaffold a project with
the [format modules](https://github.com/gofuego/fuego-formats) that match your
artifacts:

```bash
fuego init mysite --formats markdown,openapi,dbml
cd mysite
```

`init` registers each format's parsers in the generated `formats.go` and
materializes each format's **contract** — `schema.md` plus golden fixtures,
exactly what the parser emits — under `docs/formats/`.

Bring the artifacts in: copy them under `content/`, or point `dirs.content`
at the folder that already holds them
([Configuration](docs/configuration/)).

```bash
cp ../api/dispatch.openapi.yaml content/
go run . serve
```

The spec doesn't become one page — it becomes a whole **section**: an index
plus a routed page per tag, operation, and schema, visible to your taxonomies
and stable in its URLs (see tree parsers in
[Custom Parsers](docs/custom-parsers/)). Style it by adding
`theme/renderers/` and `theme/layouts/` templates written against the
contract in `docs/formats/openapi/schema.md` — no parser source required.

Later, adopt more formats with `fuego formats add <name>` and refresh the
contracts after upgrades with `fuego formats sync`
([CLI reference](docs/cli/#formats)).

## Where to go next

- Deploy the result: [GitHub Pages](docs/how-to/deploy-github-pages/) or the
  [reusable workflows](docs/how-to/deploy-with-reusable-workflows/) — rebuild
  on push and the site always matches the repo.
- Turn your theme into a shareable pack, or a specialized engine of your own:
  [Build a Format Pack](docs/tutorials/build-a-format-pack/) and
  [The Specialized Engine](docs/concepts/specialized-engine/).
