---
title: Getting Started
layout: doc
nav_section: "Getting Started"
nav_weight: 1
tags:
  - getting-started
  - tutorial
---

## Install

Fuego requires Go 1.25 or later.

```bash
go install github.com/gofuego/fuego/cmd/fuego@latest
```

This adds the `fuego` binary to your `$GOPATH/bin` (usually `~/go/bin`). Make sure it's in your PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```

You can also run without installing:

```bash
go run github.com/gofuego/fuego/cmd/fuego@latest init mysite
```

## Scaffold a Project

```bash
fuego init mysite
```

This scaffolds a working project with a `.card` flashcard DSL, a Markdown homepage, styled templates, and a dev server ready to go:

```
mysite/
  CLAUDE.md          # agent-friendly project guide
  config.yaml        # site config, parsers, collections
  main.go            # engine entry point (registers Markdown)
  content/
    index.md         # Markdown homepage
    cards/           # sample .card DSL collection (paginated)
  theme/
    base.html        # HTML shell
    layouts/         # named layouts (home, card, listing)
    partials/        # nav.html, driven by .Site.Pages
    renderers/       # per-node-type rendering (front, back, page-ref)
    outputs/         # sitemap.xml + rss.xml (non-HTML outputs)
  public/
    style.css        # starter stylesheet
```

`content/index.md` is the homepage: an `index` file routes to its directory's
root, so it becomes `/` (and `content/blog/index.md` would become `/blog/`).

### Start from a format pack

If a [format pack](docs/concepts/format-packs/) already provides the content
type you want — Markdown ADRs, Kubernetes diagrams, flashcards — scaffold with
it pre-installed:

```bash
fuego init mysite --pack github.com/gofuego/fuego-adr/adr
```

This wires `eng.Use(adr.Pack())` into `main.go` and runs `go get` to install
it. Add the pack's content under `content/` and `go run . build`. See the
[CLI reference](docs/cli/#init) for `--pack-symbol`.

## Build

```bash
cd mysite
fuego build
```

Output is written to `build/` by default. If you didn't install the CLI globally, use `go run . build` instead.

## Dev Server

```bash
fuego serve
```

Starts a local server at `http://localhost:8080` with file watching. Edit any content or theme file and the site rebuilds automatically.

## Project Structure

Every Fuego site has the same layout:

- **config.yaml** — site metadata, parser definitions, routes, taxonomies, collections, packs
- **main.go** — Go entry point. Register parsers, install packs (`eng.Use`), and add hooks here
- **content/** — your content files (anything claimed by a registered parser, by extension or filename pattern)
- **theme/** — HTML templates: `base.html`, `layouts/`, `renderers/`, `partials/`, and `outputs/`
- **public/** — static assets copied to the output root
- **build/** — generated output (gitignored)

## Content Files

Most content formats use YAML frontmatter:

```
---
title: My Page
layout: card
tags:
  - example
---
Your content here, in whatever format the parser expects.
```

The frontmatter becomes the page envelope (accessible in templates as `.Page.Envelope`). Everything below `---` is the raw payload the claiming parser turns into nodes. (Envelope extraction is the parser's job — frontmatter is how the first-party parsers do it, not an engine requirement.)

## Deployment

Set `base_url` in `config.yaml` to your deploy path:

- **Root domain** — `base_url: ""`
- **GitHub Pages subpath** — `base_url: "/my-repo"`

Because a site can deploy under a subpath, **internal links must be
base-aware** — Fuego does not rewrite them for you:

- **In templates**, prefix page URLs with the base: `{{.Site.BaseURL}}{{.URL}}`.
- **In content** (Markdown, etc.), use **base-relative** links — no leading
  slash: `[Guide](docs/guide/)`, not `[Guide](/docs/guide/)`. The theme's
  `<base href="{{.Site.BaseURL}}/">` resolves them from the site root.

A leading-slash link like `/docs/guide/` is absolute and **escapes the deploy
subpath** (it points at the domain root). Run [`build --strict-links`](docs/cli/#build)
to catch any that slip through. See [Linking](docs/templates/#linking) for the full rule.
