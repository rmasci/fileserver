package main

import (
	"fmt"
	"github.com/rmasci/fileserver"
	"log"
	"net/http"
	"os"
)

var lgOut *log.Logger

func main() {
	dwnld := fileserver.Directory{
		Lgout:  log.New(os.Stdout, "", log.Lshortfile),
		Px:     15,
		Srv:    "/var/tmp/html",
		Header: "MyFiles",
	}

	// Start web server
	http.Handle("/downloads/", dwnld)
	http.HandleFunc("/", wrongUrl)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func wrongUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<a href="http://%v/%vdownloads">Go Here</a>`, r.Host, r.URL.Path)
}
