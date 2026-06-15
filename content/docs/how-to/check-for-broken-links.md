---
title: Check for Broken Links
layout: doc
nav_section: "Guides"
nav_weight: 6
tags:
  - how-to
  - links
  - ci
---

Fuego can validate every internal link in your built site and fail the build on a
broken one. Because Fuego doesn't rewrite links for you (see [Linking](docs/templates/#linking)),
this is the safety net that catches links pointing at the wrong place.

## Run the checker

```bash
fuego build --check-links      # report broken internal links (build still succeeds)
fuego build --strict-links     # fail the build on a broken internal link
```

`--check-links` reports; `--strict-links` turns each broken link into a build
failure (and implies `--check-links`). A report names the **source file**, the
link as written, and where it resolved:

```
[LINKS] error docs/guide.md: broken link "/docs/setup/" resolves to /docs/setup, which is not a generated page
```

## Run it with your real base URL

The checker resolves each `<a href>` exactly as a browser would — against the
page's `<base href>` and the site `base_url` — and checks that the target exists.
So run it with the **base URL you actually deploy under**, or it can't catch links
that escape the deploy subpath:

```bash
fuego build --base-url /my-repo --strict-links
```

A link like `/docs/setup/` resolves fine at the root but escapes `/my-repo` once
deployed — only a check at the real base URL catches it.

## What it checks (and doesn't)

- **Checks** every `<a href>` in the rendered HTML — from content *and* templates
  — resolving relative and absolute paths through `<base href>`, and verifies the
  target is a generated page or asset.
- **Skips** external links (anything with a scheme, or protocol-relative `//…`),
  `#fragments`, and anchor existence.

It runs after the output is written, so it sees the final HTML as shipped. It is
opt-in: a plain `fuego build` does no link checking and pays no cost for it.

## Wire it into CI

Fail a deploy before it ships a dead link:

```yaml
- run: go run . build --base-url /my-repo --strict-links
```

This also catches the pull requests an in-place editor (like fuego-studio) opens,
so a typo'd link never reaches your published site. See
[Deploy to GitHub Pages](docs/how-to/deploy-github-pages/).
