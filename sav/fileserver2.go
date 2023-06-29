package main

import (
	"io"
	"net/http"
	"os"

	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
)

func main() {
	// Serve the "/assets/gfm.css" file.
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))

	var w io.Writer = os.Stdout // It can be an http.ResponseWriter.
	markdown := []byte("# GitHub Flavored Markdown\n\nHello.")

	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)

}
