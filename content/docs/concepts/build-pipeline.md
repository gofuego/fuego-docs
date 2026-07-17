---
title: The Build Pipeline
layout: doc
nav_section: "Concepts"
nav_weight: 4
tags:
  - concepts
  - architecture
---

Every Fuego command runs the same pipeline. Understanding the phases helps you predict behavior, debug issues, and write effective hooks.

## Phase Diagram

```
PREBUILD       →  shell hook (npm, tailwind, etc.)
INIT           →  merge declarative + compiled parsers
DISCOVER       →  walk content dir, apply ignore patterns
PARSE          →  run parsers: envelope + AST (concurrent)
  ↓ AfterParse hooks
ROUTE          →  resolve URLs, detect collisions
INDEX          →  build taxonomies + collections, paginate
  ↓ Index hooks          (add virtual pages → collision re-check)
  ↓ BeforeRender hooks
RENDER         →  execute templates (concurrent)
OUTPUTS        →  render theme/outputs/ (feeds, sitemaps)
MANIFEST       →  write site-manifest.json
STATIC         →  copy public/ and colocated assets
```

## Phase Details

### PREBUILD

Runs the shell command from `config.yaml`'s `prebuild` field. This is for external tooling — Tailwind CSS compilation, npm scripts, asset preprocessing. Runs before any Fuego logic.

### INIT

Merges parser sources in priority order — there are no built-in parsers:

1. **Declarative** (regex rules from config) — lower priority
2. **Compiled** (Go code via `eng.Register()`, including pack parsers via `eng.Use()`) — higher priority

If two parsers target the same file extension, the higher-priority one wins. Markdown is a first-party compiled parser you opt into with `eng.Register(markdown.Parser())`.

### DISCOVER

Walks the `content/` directory and classifies each file:

- **Content files** — claimed by a registered parser, by filename pattern (`Dockerfile`, `*.adr.md`) or bare extension (`.md`, `.trivia`, `.card`). Patterns are checked before extensions and the most specific claim wins, so `guide.adr.md` belongs to the parser claiming `*.adr.md`, not the `md` parser.
- **Binary assets** — images, fonts, etc. — copied to output alongside their routed content
- **Ignored files** — matched by `ignore` glob patterns in config

### PARSE

For each content file, in parallel:

1. Dispatch the raw file to the parser that claimed it at DISCOVER
2. The parser returns the envelope (metadata) and `[]Node` — the universal AST. Envelope extraction is the parser's job: YAML frontmatter is one option (`core.WithYAMLFrontmatter`), not a requirement

A parser that also implements `core.TreeParser` expands one file into a **tree of real pages**: the engine calls `ParseTree` instead of `Parse` and turns the returned root plus children into an index page and one page per child, each with its own envelope — so a single OpenAPI spec or database schema becomes a whole routed section. See [Custom Parsers](docs/custom-parsers/).

Concurrency is bounded by `runtime.NumCPU()` via `errgroup`.

### ROUTE

Assigns a URL to each page using three-tier priority:

1. Frontmatter `slug` field — overrides the slug segment
2. Config `routes[type]` pattern — pattern expansion with `{dir}`, `{slug}`, `{filename}`
3. Filesystem mirror — the default. An `index` file routes to its directory's
   root: `content/index.md` → `/`, `content/blog/index.md` → `/blog/`.

Pages expanded from a tree parser are routed in a second pass: the root goes through the three tiers as normal, then each child's URL is the root's resolved URL joined with the child's slug-path segments.

After resolution, checks for URL collisions. A collision is a `GlobalFatal` error that stops the build.

### INDEX

Generates virtual pages for taxonomies and collections:

- **Taxonomy term pages** — one per unique term value (e.g., `/tags/go/`)
- **Taxonomy index pages** — list all terms (e.g., `/tags/`); generated whenever `index_path` is configured, even when no page carries the field yet
- **Collection pages** — glob-matched, sorted listing pages

Listings with `page_size` set are split into numbered pages (`/blog/`, `/blog/page/2/`, …), each carrying a `.Paginator`. **Index hooks** run here too — the supported place to add your own virtual pages, since their URLs go through the collision re-check that catches conflicts between virtual and real pages.

### RENDER

For each page, in parallel:

1. Pre-render nodes to HTML using the default renderer (or per-type renderer templates)
2. Build the template data (`.Page`, `.Site`, `.Paginator`, `.JSON`)
3. Execute the base template with the selected layout
4. Write `{url}/index.html` to the output directory

Pages marked `Skip` by a hook are excluded from RENDER and the manifest.

### OUTPUTS

Renders every file under `theme/outputs/` as a text template fed with `.Site`, writing non-HTML site assets — RSS feeds, sitemaps, `robots.txt`, search indexes — to matching output paths. See [Add an RSS Feed and Sitemap](docs/how-to/add-feeds-and-sitemaps/).

### MANIFEST

Writes `site-manifest.json` — a JSON index of all pages, taxonomy terms, and collection membership. Each page entry records its `url`, `type`, `layout`, `title`, `summary`, `output_path` (the generated file, e.g. `blog/post/index.html`), `source_path` (the source file relative to the content directory — *omitted* for virtual pages, which are non-editable), and the flattened `envelope`. Pages expanded from a tree-parsed artifact all list the shared artifact as their `source_path`, so several entries may map to one source file — root and children stay distinguishable by `url`/`output_path`; the top-level `content_root` is the content directory relative to the repository root. This is useful for client-side search and navigation, and for mapping a served URL back to the source file it was built from — what an external host or in-place editor needs.

### STATIC

Copies static assets to the output root, in precedence order: each registered pack's `static/` subtree first, then the user's `public/` directory (so user files win on conflict), then content-colocated binary assets at their routed paths.

## Partial Execution

`validate` runs through INDEX without RENDER — catching config errors, parse failures, and collisions without producing output. `list` also runs through INDEX and prints the page table, virtual pages included. This is controlled by `pipeline.RunUntil(phase)`.

## Incremental Builds

`fuego build --incremental` (and the dev server, always) keeps an on-disk cache of parsed pages so unchanged content isn't re-parsed on every rebuild.

The cache is keyed by a header — a hash of the **engine binary**, the **resolved config**, and the **theme tree** — plus a per-file **content hash**:

- If the header matches, content files whose hash is unchanged skip PARSE and are restored from cache; changed and new files are parsed normally; deleted files have their output removed (orphan cleanup via a manifest-style diff).
- If the header doesn't match — you rebuilt the engine, edited `config.yaml`, or touched a template — the whole cache is discarded and the build falls back to a full, clean rebuild.

ROUTE and INDEX always run over the full page set (they are cheap, `O(pages)` map work). RENDER is **narrowed** to the pages whose output can actually differ: the changed/new pages, every page expanded from a changed tree-parsed artifact, every virtual page (taxonomy, collection, pagination — they aggregate content), and any page whose template reads `.Site.Pages` directly or through a partial. Pages with a site-blind layout that didn't change keep their existing HTML. On a 10k-page site, a single-file edit rebuilds in a fraction of a full build.

Despite the narrowing, an incremental build produces **byte-identical output to a clean build** of the same inputs. That guarantee is enforced in CI: every fixture is built clean and then incrementally under each mutation (edit, add, delete, theme touch, config touch) and the output trees are compared byte-for-byte. A corrupt or version-mismatched cache is treated as a miss, never an error.
