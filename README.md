Fileserver offers a file server that gives you information on the files in the directory such as size, text, binary, directory.  The file server also has a built in logger that can log to a file or to stdout. It also allows you to specify a directory that is outside of the path your binary is running in. 

In the example, when you go to "http://localhost:8888/downloads" you get a basic list of files that you can click on, see screenshot below.  

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
![alt text](https://github.com/rmasci/fileserver/blob/main/Fileserver.png?raw=true)
Todo: Document code, document better examples.
