---
title: CLI Reference
layout: doc
nav_section: "Reference"
nav_weight: 2
tags:
  - reference
  - cli
---

## Install

```bash
go install github.com/gofuego/fuego/cmd/fuego@latest
```

Requires Go 1.25+. The binary is placed in `$GOPATH/bin` (usually `~/go/bin`). Ensure this directory is in your `PATH`.

Alternatively, run any command without installing:

```bash
go run github.com/gofuego/fuego/cmd/fuego@latest <command>
```

## Commands

### build

Build the static site. Runs the full pipeline and writes output to `build/`.

```bash
fuego build
fuego build --incremental            # reuse cached parses for unchanged content
fuego build --base-url /owner/repo   # override the deploy subpath without editing config
fuego build --check-links            # report internal links that don't resolve
fuego build --strict-links           # fail the build on a broken internal link
```

With `--incremental`, Fuego keeps a build cache and re-parses only the content files whose bytes changed since the last build; deleted pages have their output removed. A change to the engine binary, the resolved config, or the theme invalidates the cache and triggers a full, clean rebuild — so incremental output is always identical to a clean build. The dev server (`serve`) uses incremental builds automatically. See [The Build Pipeline](docs/concepts/build-pipeline/#incremental-builds).

`--base-url` overrides the site's `base_url` (the deploy subpath, e.g. `/owner/repo`) for that build, so a deploy workflow can set it per-target without a separate config file; pass an empty value (`--base-url ""`) to build for the root. It's also available on `serve`.

`--check-links` resolves every internal `<a href>` against the page's `<base href>` and the site base URL, and reports links that don't land on a generated page; `--strict-links` makes such a link fail the build (for CI). Run it with `--base-url` set to the real deploy path so it catches links that escape the deployment base. See [Check for Broken Links](docs/how-to/check-for-broken-links/).

### serve

Start a dev server with file watching and live rebuild.

```bash
fuego serve
```

Watches `content/` and `theme/` for changes. When a file changes, the site is rebuilt and served at `http://localhost:8080` (configurable via `dev.port` in config).

### validate

Check config and content for errors without producing output. Useful as a CI gate.

```bash
fuego validate
```

Runs the pipeline through INDEX (discovery, parsing, routing, collision detection) without rendering. Exit code 0 on success, 1 on any error. Because it doesn't render, `validate` can't check links — use `build --strict-links` for that.

### list

Print all pages as a table of TYPE, SOURCE, and URL.

```bash
fuego list
```

### config

Print the fully resolved configuration — your `config.yaml` deep-merged with every registered pack's defaults — annotated with per-key provenance (`# user` or `# pack: name`).

```bash
fuego config
```

Useful for answering "why is this value what it is?" when format packs contribute config defaults. Output is deterministic, so it is safe to diff. See [Config Merging](docs/config-merging/).

### init

Scaffold a new Fuego project.

```bash
fuego init mysite
```

Creates a working project with a `.card` flashcard DSL, theme, and sample content.

**Start from a format pack** with `--pack`:

```bash
fuego init mysite --pack github.com/gofuego/fuego-adr/adr
```

This scaffolds the project, wires the pack into `main.go`
(`eng.Use(adr.Pack())`), and runs `go get` to install it. The convention is
that a pack module's package exports a `Pack() core.Pack` function; the package
name defaults to the module path's last segment. If the package name differs,
pass `--pack-symbol <name>`. Init never compiles or runs the pack — it only
downloads it — so installing a pack can't execute third-party code.

## Global Flags

| Flag | Default | Description |
|---|---|---|
| `--config` | `config.yaml` | Path to configuration file |
| `--version` | | Print the version and exit (reports the installed module version) |

## Error Handling

Three severity levels:

- **Warning** — logged, build continues
- **LocalFatal** — page skipped, build continues
- **GlobalFatal** — build fails immediately

`validate` catches config errors, parse failures, and URL collisions before you build.
