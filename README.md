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

	"github.com/rmasci/fileserver"
)

var lgOut *log.Logger

func main() {
	dwnld := fileserver.Directory{
		Lgout:  log.New(os.Stdout, "", log.Lshortfile),
		Px:     15,
		Srv:    "/var/tmp/html",
		Header: "MyFiles",
	}

	http.Handle("/downloads/", dwnld)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
```
Screenshot:<br>
<img="Fileserver.png"><br>
Todo: Document code, document better examples.
