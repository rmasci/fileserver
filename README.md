Fileserver offers a more advanced file server than the basic file server that is a part of Go.

In Go to serve a directory you can pass something like this:
```
http.Handle("/", http.FileServer(http.Dir(*directory)))
```
When you go to "/" you get a basic list of files that you can click on.  What this fileserver does is offer more information on the files.

Example:

```
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
		Lgout:  lgOut,
		Px:     15,
		Srv:    html,
		Header: "MyFiles",
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

```
![alt-text][screenshot] (https://github.com/rmasci/fileserver/Fileserver.png "Screenshot Fileserver.png")
Todo: Document code, document better examples.
