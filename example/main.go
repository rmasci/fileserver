package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rmasci/fileserver"
)

var lgOut *log.Logger
var nginx string

func main() {
	lgOut = log.New(os.Stdout, "", log.Lshortfile)
	nginx = "nginx"
	lgOut.Println(pattern("downloads"))
	dwnld := fileserver.Directory{
		Px:      15,
		Srv:     "/var/www/html",
		BaseURI: pattern("downloads"),
		Lgout:   lgOut,
		Header:  "MyFiles",
	}

	// Start web server
	http.Handle(pattern("downloads"), http.StripPrefix(pattern("downloads/"), http.HandlerFunc(dwnld.Fileserver)))
	http.HandleFunc("/", printpath)
	err := http.ListenAndServe(":8789", nil)
	if err != nil {
		fmt.Println("Error")
	}
}
func pattern(path string) string {
	path = strings.Trim(path, "/")
	nginx = strings.Trim(nginx, "/")
	retstr := fmt.Sprintf("/%s/%s/", nginx, path)
	retstr = strings.ReplaceAll(retstr, "//", "/")
	return retstr
}

func printpath(w http.ResponseWriter, r *http.Request) {
	lgOut.Println(r.RequestURI)
	fmt.Fprintf(w, "<html><h1>URI</h1><hr>%v</html>", r.RequestURI)
}
