---
title: Custom Parsers
layout: doc
nav_section: "Reference"
nav_weight: 4
tags:
  - reference
  - parsers
---

Fuego supports two ways to define content parsers: declarative (config-only) and compiled (Go code). Both produce the same universal AST.

## Universal AST

All parsers produce `[]Node`:

```go
type Node struct {
    Type       string
    Attributes map[string]any
    Content    string
    Children   []Node
    Raw        bool // pass Content through as raw HTML, unwrapped and unescaped
}
```

The engine never interprets node types. Your templates decide how `question`, `front`, `answer`, or any other type renders.

## Declarative Parsers

Define parsers with ordered regex rules in `config.yaml`. No Go code needed:

```yaml
parsers:
  trivia:
    rules:
      - match: '^\?\s*(.+)$'
        emit:
          type: question
          content: "$1"
      - match: '^\[([A-Z])\]\s+(.+)$'
        emit:
          type: answer
          content: "$2"
          attributes:
            letter: "$1"
```

This parses `.trivia` files line-by-line. First matching rule wins per line. Capture groups (`$1`, `$2`) are substituted into content and attributes.

## Compiled Parsers

Implement the `core.Parser` interface for full control. The parser receives the entire raw file and owns envelope extraction:

```go
type Parser interface {
    Type() string
    Parse(raw []byte) (Envelope, []Node, error)
}
```

Use `core.SplitFrontmatter(raw)` if your format carries YAML frontmatter, or the convenience wrappers `core.WithYAMLFrontmatter(...)` / `core.WithNoEnvelope(...)` for the common cases. To claim files by name rather than bare extension — extensionless files like `Dockerfile`, or compound suffixes like `*.adr.md` — also implement `core.FilenameParser` with `Filenames() []string`. Patterns match the base filename, and a filename claim beats a bare-extension claim (see Parser Precedence below).

**Keep envelope values JSON-shaped** — maps, slices, strings, numbers, bools (`map[string]any`, `[]map[string]string`, `[]string`, …) rather than your own structs or pointers. Anything a YAML or JSON decode could produce is fine. Pages whose envelopes hold other types still build normally, but are skipped by the incremental-build cache (a warning names them), so they re-parse on every `serve` rebuild.

Register it in `main.go`:

```go
eng := engine.New()
eng.Register(&MyCustomParser{})
eng.Run(os.Args)
```

## Tree Parsers

A plain parser turns one file into one page. A parser that *also* implements `core.TreeParser` can expand one rich artifact — an OpenAPI spec, a database schema — into a whole section: an index page plus a tree of real pages.

```go
type TreeParser interface {
    Parser
    ParseTree(raw []byte) (*PageTree, error)
}

type PageTree struct {
    Envelope Envelope
    Nodes    []Node
    Children map[string]*PageTree
}
```

The engine detects `TreeParser` by interface assertion at PARSE — no registration change, and plain parsers are untouched. When a parser implements it, the engine calls `ParseTree` (`Parse` is not called for that file) and expands the returned tree:

- The **root** `PageTree`'s envelope and nodes become the source file's own routed page, exactly as a plain parser's output would.
- Each entry in **`Children`** becomes a real page. Keys are **relative slug paths** — one or more `/`-separated segments like `"tags/billing/get-invoice"` — and nesting a `PageTree` under a key extends the path.
- A child's **URL** is the root's routed URL joined with its slug-path segments, composed *after* the root goes through normal three-tier routing — so `slug`, route patterns, and the index-file convention on the root all still apply. If `content/api.openapi.yaml` routes to `/api/`, the child `"operations/get-invoice"` lands at `/api/operations/get-invoice/`.
- Children carry their **own envelopes**, so taxonomies, collections, and pagination see them like any hand-written page.

```go
func (p *APIParser) ParseTree(raw []byte) (*core.PageTree, error) {
    spec, err := decodeSpec(raw)
    if err != nil {
        return nil, err
    }
    tree := &core.PageTree{
        Envelope: core.Envelope{"title": spec.Title},
        Nodes:    overviewNodes(spec),
        Children: map[string]*core.PageTree{},
    }
    for _, op := range spec.Operations {
        tree.Children["operations/"+op.Slug] = &core.PageTree{
            Envelope: core.Envelope{"title": op.Summary, "tags": op.Tags},
            Nodes:    operationNodes(op),
        }
    }
    return tree, nil
}
```

Conventions and behavior to know:

- **Keep child envelopes JSON-shaped** (same rule as above) so the whole tree stays cache-eligible. On incremental builds the tree is cached as one unit under the source file's content hash: one child holding a non-JSON-shaped value drops the *whole file* from the cache, not just that child.
- A child envelope may set `layout` and `type` like any page; a **missing layout falls back to the base template** silently.
- Two children whose slug paths compose to the same URL, or a child colliding with another page, fail the build through the engine's normal URL-collision detection — there is no tree-local check.
- In the [site manifest](docs/concepts/build-pipeline/), every page of a tree lists the shared source artifact as its `source_path`.

## Error Positions

Return a `core.ParseError` to point build errors at a source line:

```go
return nil, nil, &core.ParseError{
    Line: lineNo,
    Err:  fmt.Errorf("unexpected token %q", tok),
}
```

The build output then reads `[PARSE] error content/quiz.trivia:7: ...`. Reporting positions is optional — plain errors keep working. Frontmatter errors from `core.SplitFrontmatter` carry file-relative line numbers automatically.

## Parser Precedence

There are no built-in parsers — the engine ships with zero format opinions. When multiple parsers exist for the same file extension:

1. **Compiled** (registered via `eng.Register()`) — wins
2. **Declarative** (defined in `config.yaml`)

A file is discovered as content only if a registered parser **claims** it, and the same claim decides which parser parses it. Claims are resolved by specificity:

1. **Filename patterns beat bare extensions** — `guide.adr.md` goes to the parser claiming `*.adr.md`, not to the `md` parser.
2. **The longest matching pattern wins** — `*.adr.md` beats `*.md`.
3. **Equal-length ties resolve by parser precedence** — compiled parsers above declarative ones.

If no pattern matches, the bare extension is looked up; if nothing claims the file, it is copied as a static asset.

A parser claims by exactly one kind: if it declares filename patterns (`FilenameParser`), those patterns are its **complete** claim set — its `Type()` is not implicitly claimed as an extension. A parser without patterns claims its `Type()` as a bare extension.

## Markdown (opt-in)

Markdown is a first-party, opt-in parser — not a built-in. Register it like any other compiled parser:

```go
import "github.com/gofuego/fuego/parsers/markdown"

eng.Register(markdown.Parser())
```

It uses goldmark with GitHub-Flavored Markdown (tables, strikethrough, autolinks, task lists) and emits a single node with `Raw: true` containing the rendered HTML.

By default it claims the bare `md` extension. For a repo whose markdown files don't match — or to claim only specific files — override the claim with `markdown.WithPatterns(...)`; the patterns replace the default claim entirely:

```go
eng.Register(markdown.Parser(markdown.WithPatterns("*.markdown")))
eng.Register(markdown.Parser(markdown.WithPatterns("README.md"))) // only README.md
```

The parser's full output contract (claims, envelope keys, node types) is documented in [`parsers/markdown/schema.md`](https://github.com/gofuego/fuego/blob/develop/parsers/markdown/schema.md), following the fuego-formats schema convention.
