# fuego-docs

The documentation site for [Fuego](https://github.com/gofuego/fuego) — itself a
Fuego project (dogfooding). Built against the published engine and hosted +
edited in place via [fuego-studio](https://github.com/gofuego/fuego-studio).

## Develop

```bash
go run . serve                                 # dev server with live reload
go run . build --base-url /gofuego/fuego-docs  # production build → build/
```

A gitignored `go.work` builds against a local engine checkout
(`../test-render-eng`) so you can author against unreleased engine features; CI
builds against the published `gofuego/fuego`.
