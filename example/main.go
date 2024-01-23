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
	dwnld := fileserver.Directory{
		Lgout:  log.New(os.Stdout, "", log.Lshortfile),
		Px:     15,
		Srv:    html,
		Header: "MyFiles",
	}

	http.HandleFunc("/downloads/", dwnld.Fileserver)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
