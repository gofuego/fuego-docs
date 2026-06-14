package main

import (
	"fmt"
	"os"

	"github.com/gofuego/fuego-docs/docspack"
	doctheme "github.com/gofuego/fuego-doctheme"
	"github.com/gofuego/fuego/engine"
	"github.com/gofuego/fuego/parsers/markdown"
)

func main() {
	eng := engine.New()
	eng.Register(markdown.Parser())

	// Dogfood the pack system twice over: the docs pack supplies the tags
	// taxonomy and paginated tutorials collection as config defaults, and the
	// shared theme pack supplies the look. This site keeps only its topbar,
	// sidebar, and feed/sitemap outputs in theme/.
	eng.Use(docspack.Pack())
	eng.Use(doctheme.Public())

	if err := eng.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fuego: %v\n", err)
		os.Exit(1)
	}
}
