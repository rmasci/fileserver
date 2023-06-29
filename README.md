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
```
![alt-text][screenshot] (https://github.com/rmasci/fileserver/Fileserver.png "Screenshot Fileserver.png")
Todo: Document code, document better examples.
