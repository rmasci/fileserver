package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rmasci/fileserver"
)

func main() {
	var d fileserver.Directory
	d.Px = 15
	// Need to have this set to the uri you want the user to use: http://yourserver.dom.com/download
	d.BaseURI = "/download/"
	port := flag.String("p", "8100", "port to serve on")
	flag.StringVar(&d.Srv, "d", "/var/www/html", "the directory of static file to host")
	flag.Parse()

	//http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.Handle(d.BaseURI, http.StripPrefix("/", http.HandlerFunc(d.Fileserver)))
	log.Printf("Serving %s on HTTP port: %s\n", d.Srv, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
