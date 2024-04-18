package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rmasci/fileserver"
)

var lgOut *log.Logger

func main() {
	html := os.Getenv("PWD") + "/html"
	lgOut := log.New(os.Stdout, "", log.Lshortfile)
	dwnld := fileserver.Directory{
		Lgout:   lgOut,
		Px:      15,
		Srv:     html,
		Header:  "MyFiles",
		Default: "index.html",
	}

	http.HandleFunc("/downloads/", dwnld.Fileserver)
	//redirect to downloads
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/downloads/", http.StatusFound)
	})

	// let user know when server is listening on port
	lgOut.Println("Listening on port 8888")

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
