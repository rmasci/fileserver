package main

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// This prints AST of parsed markdown document.
// Usage: printast <markdown-file>

func usageAndExit() {
	fmt.Printf("Usage: printast [-to-html] <markdown-file>\n")
	os.Exit(1)
}

func main(d []byte) []byte {
	d = markdown.NormalizeNewlines(d)
	exts := parser.CommonExtensions // parser.OrderedListStart | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(exts)
	doc := markdown.Parse(d, p)

	htmlFlags := mdhtml.Smartypants |
		mdhtml.SmartypantsFractions |
		mdhtml.SmartypantsDashes |
		mdhtml.SmartypantsLatexDashes
	htmlOpts := mdhtml.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := mdhtml.NewRenderer(htmlOpts)
	html := markdown.Render(doc, renderer)
	return html

}
