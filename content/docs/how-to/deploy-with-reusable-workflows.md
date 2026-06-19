---
title: Deploy with Reusable Workflows
layout: doc
nav_section: "Guides"
nav_weight: 6
tags:
  - how-to
  - deployment
  - ci
---

Instead of hand-rolling a deploy pipeline (see [Deploy to GitHub
Pages](docs/how-to/deploy-github-pages/)), the **gofuego** org publishes
**reusable GitHub Actions workflows** that build and publish a Fuego site for you.
A caller workflow is a few lines: you say *which base URL* and *where to publish*,
the reusable workflow does the rest. They live in the public
[`gofuego/.github`](https://github.com/gofuego/.github) repo — reference them
directly, or fork them into your own org and pin to a tag.

## The model

`fuego-deploy` builds your site with a `--base-url` and pushes the generated
`build/` output to a **serve branch**:

- `pages_base` → built for GitHub Pages, pushed to **`gh-pages`**
- `studio_base` → built for the [fuego-studio](https://github.com/gofuego/fuego-studio)
  mount path, pushed to **`fuego-pages`**

Pass either or both — one workflow covers a Pages-only site, a studio-only site,
or a dual-hosted one.

> **This is branch-based publishing.** Set your repo's **Settings → Pages →
> Source** to **"Deploy from a branch" → `gh-pages` / `(root)`**. (That differs
> from the hand-rolled guide, which uses the "GitHub Actions" artifact source —
> pick one model, not both.)

## Deploy to GitHub Pages

Add `.github/workflows/deploy.yml`, replacing `my-repo` with your repository name:

```yaml
name: Deploy
on:
  push: { branches: [main] }
  workflow_dispatch:
jobs:
  deploy:
    permissions:
      contents: write          # required — the workflow pushes to gh-pages
    uses: gofuego/.github/.github/workflows/fuego-deploy.yml@main
    with:
      pages_base: /my-repo
      check_links: true        # fail the deploy on a broken internal link
```

On every push to `main`, your site builds and publishes to `gh-pages`.

## Check pull requests before they merge

Pair the deploy with `fuego-check` — it compiles and validates the site (`go run .
validate`) **without** publishing, so a broken build is caught on the PR. Add
`.github/workflows/check.yml`:

```yaml
name: Check
on:
  pull_request:
  push:
    branches-ignore: [main, gh-pages, fuego-pages, fuego-studio-edits]
jobs:
  check:
    uses: gofuego/.github/.github/workflows/fuego-check.yml@main
```

With a `develop`-default / protected-`main` setup, this means: **check runs on
pull requests and `develop` pushes; deploy runs only when you merge to `main`.**

## Dual-host (Pages + fuego-studio)

Hosting the same site on GitHub Pages *and* in fuego-studio? Pass both bases —
each builds with its own mount path and publishes to its own branch:

```yaml
    with:
      pages_base: /my-repo            # → gh-pages
      studio_base: /my-org/my-repo    # → fuego-pages
      check_links: true
```

## A site that lives in a subdirectory

If your Fuego project is a **subset** of a larger repo (e.g. a `docs/` folder
inside a code repo), set `project_dir` — the workflow builds there and publishes
`docs/build`:

```yaml
    with:
      project_dir: docs
      pages_base: /my-repo
```

`fuego-check` takes the same `project_dir` input.

## Other reusable workflows

| Workflow | For | Key inputs |
|----------|-----|------------|
| `fuego-deploy.yml` | building & publishing a site | `pages_base`, `studio_base`, `project_dir`, `check_links` |
| `fuego-check.yml` | validating a site on PRs | `project_dir` |
| `go-ci.yml` | a Go module (a custom parser/pack repo): `build` + `vet` + `test -race` | `working_directory`, `race` |
| `fuego-adr-deploy.yml` | publishing an [ADR](https://github.com/gofuego/fuego-adr) site from `*.adr.md` files | `adr_path`, `pages_base`, `studio_base` |

All default the Go toolchain to `1.25` (Fuego's minimum). Full input reference is in
the [`gofuego/.github` README](https://github.com/gofuego/.github).

## Base-aware links

However you deploy, a site served under a subpath needs base-aware links or they
break. Prefix template URLs with `{{.Site.BaseURL}}` and use base-relative links
(no leading slash) in content — see [Linking](docs/templates/#linking) and
[Check for Broken Links](docs/how-to/check-for-broken-links/).
