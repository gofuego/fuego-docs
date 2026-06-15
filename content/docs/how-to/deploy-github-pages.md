---
title: Deploy to GitHub Pages
layout: doc
nav_section: "Guides"
nav_weight: 5
tags:
  - how-to
  - deployment
---

## Prerequisites

- A GitHub repository with your Fuego site
- GitHub Pages enabled in repo settings (source: "GitHub Actions")

## 1. Set the Base URL

In `config.yaml`, set `base_url` to your repository name:

```yaml
site:
  name: "My Site"
  base_url: "/my-repo"
```

GitHub Pages serves your site at `https://username.github.io/my-repo/`, so all links and asset paths need this prefix.

## 2. Make Links Base-Aware

Under a subpath, links must include the base URL or they break. In templates,
prefix page URLs with `{{.Site.BaseURL}}`; in content, use base-relative links
(no leading slash). See [Linking](docs/templates/#linking).

```html
<link rel="stylesheet" href="{{.Site.BaseURL}}/style.css">
<a href="{{.Site.BaseURL}}/">Home</a>
```

`content/index.md` routes to `/`, so there's a real `index.html` at the site
root — no redirect file is needed.

## 3. Create the Workflow

Create `.github/workflows/docs.yml`:

```yaml
name: Deploy Site

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - run: go run . build --strict-links   # fail the deploy on a broken internal link
      - uses: actions/upload-pages-artifact@v3
        with:
          path: build

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
```

## 4. Enable GitHub Pages

Go to your repository's Settings > Pages and set the source to **"GitHub Actions"**.

Push to `main` and the workflow will build and deploy your site automatically.
