package fileserver

/*
Todo: Use Go templates -- add <copy link> <wget>
Todo: add buttons for <link><wget link>
*/

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cgi"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pyk/byten"
	"github.com/rmasci/script"
	"github.com/yuin/goldmark"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Directory struct {
	// Render this as the default page.
	Default string
	// True renders the default on top, false is on bottom. Default is on bottom
	DefPlace  bool
	Srv       string
	Px        int
	BaseURI   string
	Lgout     *log.Logger
	Header    string
	Directory string
	Template  string
	PHP       string
	Perl      string
	Python    string
	Bash      string
	// name, date, size
	Sort      string
	Direction string
	// location either in embed or elsewhere. I'll put a default in.
}

func New(lgOut *log.Logger, srv string, opts map[string]string) Directory {
	var err error
	d := Directory{
		Lgout:     lgOut,
		Px:        15,
		Srv:       srv,
		Header:    "MyFiles",
		Default:   "index.html",
		Sort:      "time",
		Direction: "dec",
	}
	d.Perl, err = exec.LookPath("perl")
	if err != nil {
		d.Perl = "notfound"
	}
	d.PHP, err = exec.LookPath("php-cgi")
	if err != nil {
		d.PHP = "notfound"
	}
	d.Python, err = exec.LookPath("python")
	if err != nil {
		d.Python = "notfound"
	}
	d.Bash, err = exec.LookPath("bash")
	if err != nil {
		d.Bash = "notfound"
	}

	for key, value := range opts {
		switch key {
		case "Px":
			if px, err := strconv.Atoi(value); err == nil {
				d.Px = px
			}
		case "Header":
			d.Header = value
		}
	}

	return d
}
func queryPath(uri string) (string, string) {
	parsedUrl, _ := url.Parse(uri)
	return parsedUrl.Path, parsedUrl.RawQuery
}

func (d *Directory) Fileserver(w http.ResponseWriter, r *http.Request) {
	//var upDir string
	// This will hold the default page, if the d.Default file is found.
	var cgiRun cgi.Handler //server.CgiServer()

	var sty string
	var err error
	timeFormat := "2006-01-02 15:04:05"
	reqDir := strings.Trim(r.RequestURI, "/")
	reqCount := len(strings.Split(strings.Trim(d.BaseURI, "/"), "/"))
	blackFile := blackFile(d.Px)
	blackFolder := blackFolder(d.Px)
	reqDirA := strings.Split(reqDir, "/")[reqCount:]
	srv, query := queryPath(fmt.Sprintf("%v/%v", d.Srv, strings.TrimRight(strings.Join(reqDirA, "/"), "/")))
	query, err = url.QueryUnescape(query)
	if err != nil {
		d.Lgout.Printf("Error unescaping query: %v", err)
		http.Error(w, "Error processing query", http.StatusBadRequest)
		return
	}
	dir, err := os.Stat(srv)
	if err != nil {
		w.WriteHeader(404)
		http.Error(w, fmt.Sprintf("%v", err), 404)
		return
	}
	if dir.IsDir() {
		fslist, err := filepath.Glob(srv + "/*")
		if err != nil {
			w.WriteHeader(404)
			http.Error(w, fmt.Sprintf("%v", err), 404)
			return
		}
		fslist = d.sortFSlist(fslist)

		//fmt.Fprintln(w, "<html>\n\t<head>\n\t\t<style>table, th, td {border: 0px;padding: 0px;} tr:nth-child(odd) {background-color: #E0E0E0;}\n\t\t</style>\n\t</head>")
		fmt.Fprintf(w, "<html><head>%s</head><body>\n", getHeader())
		// Detect if there is a default file in there.
		if d.DefPlace && d.Default != "" {
			for _, f := range fslist {
				if filepath.Base(f) == d.Default {
					d.Lgout.Printf("found Default, displaying it%s", filepath.Base(f))
					b, err := os.ReadFile(srv + "/" + f)
					if err != nil {
						fmt.Fprintf(w, "<b>Error Reading %s</b>%v", err)
					}
					if filepath.Ext(srv) == ".md" {
						mdHtml, err := mdToHTML(b)
						if err != nil {
							d.Lgout.Println("Error converting markdown to html", err)
							w.Write(b)
						}
						w.Write(mdHtml)
					} else {
						w.Write(b)
					}
					fmt.Fprintf(w, "<hr><br>")
				}
			}
		}
		fmt.Fprintln(w, `<div class="container"`)
		if d.Header != "" {
			fmt.Fprintf(w, "\t<body>\n\t<br>\n\t<br>\n\t<br>\n\t<h1>%v</h1>\n", d.Header)
		} else {
			fmt.Fprintln(w, "<body><br><br><br>")
		}
		if len(reqDirA) > 0 {
			// if r.RequestURI has a '/' on the end, it will not remove the last directory
			upDir := path.Dir(strings.TrimRight(r.RequestURI, "/"))

			fmt.Fprintf(w, "\t<a href='%v/'>Parent Directory</a>\n", upDir)
		}

		fmt.Fprintln(w, `<hr><table class="table table-hover">`)
		fmt.Fprintln(w, `<thead class="thead-dark">`)
		fmt.Fprintln(w, `<tr><th style="width: 3%" ></th><th>Name</th><th>Size</th><th>Date</th><th></th></tr>`)
		for i, f := range fslist {
			var ico string
			fstat, err := os.Stat(f)
			if err != nil {
				continue
			}
			if fstat.IsDir() {
				ico = blackFolder
			} else {
				ico = blackFile
			}
			//Odd or even:
			if i%2 == 0 {
				sty = fmt.Sprint(`class="table-default"`)
			} else {
				sty = fmt.Sprint(`class="table-success"`)
			}

			//link := fmt.Sprintf("<a href='%v/%v'>%v</a>", r.URL, fstat.Name(), fstat.Name())
			link := fmt.Sprintf("<a href='%v/%v'>%v</a>", r.RequestURI, fstat.Name(), fstat.Name())
			// 'File Permissions' 'Link to file', 'File name' 'file size' 'file modifiied'
			link = strings.ReplaceAll(link, "//", "/")
			fullLink := fmt.Sprintf("%v/%v/%v", r.Host, r.RequestURI, fstat.Name())
			fullLink = "https://" + strings.ReplaceAll(fullLink, "//", "/")
			out := fmt.Sprintf(`
          <tr %s>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
			<td><button class="btn btn-outline-primary btn-sm rounded" onclick="copyToClipboard('%v')">Copy Link</button>
				<button class="btn btn-outline-success btn-sm rounded" onclick="copyToClipboard('wget --output-document=%v %v')">wget</button>
			</td>
		  </tr>`, sty, ico, link, byten.Size(fstat.Size()), fstat.ModTime().Format(timeFormat), fullLink, fstat.Name(), fullLink)
			fmt.Fprintf(w, "%v", out)

		}
		fmt.Fprintln(w, "\n\t\t\t</table></div>\n\t\t<hr>\n\t")
		if !d.DefPlace && d.Default != "" {
			for _, f := range fslist {
				if filepath.Base(f) == d.Default {
					d.Lgout.Printf("found Default, displaying it %s", filepath.Base(f))
					b, err := os.ReadFile(f)
					if err != nil {
						fmt.Fprintf(w, "<b>Error Reading %s</b>%v", err)
					}
					switch filepath.Ext(srv) {
					case ".md":
						mdHtml, err := mdToHTML(b)
						if err != nil {
							d.Lgout.Println("Error converting markdown to html", err)
							w.Write(b)
						}
						w.Write(mdHtml)

					case ".php": //, ".pl", ".py", ".sh":
						cgiRun.Path = "/usr/bin/php-cgi"
						cgiRun.Dir = filepath.Dir(f)
						cgiRun.Env = []string{"SCRIPT_FILENAME=" + f}
						d.Lgout.Println("php")
						d.Lgout.Println(cgiRun)
						cgiRun.ServeHTTP(w, r)
					case ".pl":
						cgiRun.Path = "/usr/bin/perl"
						cgiRun.Dir = filepath.Dir(f)
						cgiRun.Env = []string{"SCRIPT_FILENAME=" + f}
						cgiRun.ServeHTTP(w, r)
					case ".py":
						cgiRun.Path = "/usr/bin/python"
						cgiRun.Dir = filepath.Dir(f)
						cgiRun.Env = []string{"SCRIPT_FILENAME=" + f}
						cgiRun.ServeHTTP(w, r)
					case ".sh":
						cgiRun.Path = "/bin/bash"
						cgiRun.Dir = filepath.Dir(f)
						cgiRun.Env = []string{"SCRIPT_FILENAME=" + f}
						cgiRun.ServeHTTP(w, r)

						//if php, err := exec.LookPath("php"); err != nil {
						//	d.Lgout.Println("php is not in the executable path")
						//	w.Write(b)
						//} else {
						//	// Retrieve posted variables
						//	postedVars := r.PostForm.Encode()
						//	cmd := fmt.Sprintf("%v %v", php, f)
						//	// Execute PHP script with posted variables
						//	phpHtml, err := script.Echo(postedVars).Exec(cmd).Bytes()
						//	if err != nil {
						//		d.Lgout.Println("Error executing php", err)
						//		w.Write(b)
						//	}
						//	w.Write(phpHtml)

					default:
						d.Lgout.Println("default")
						w.Write(b)
					}
					fmt.Fprintf(w, "<hr><br>")
				}
			}
		}
		fmt.Fprintln(w, "</body>\n</html>")
		return
	} else {
		d.Lgout.Println("Open", srv)
		// Detect Content Type
		openFile, err := os.Open(srv)
		defer openFile.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 404)
			return
		}
		FileHeader := make([]byte, 512)
		openFile.Read(FileHeader)
		FileContentType := http.DetectContentType(FileHeader)
		d.Lgout.Printf("Filename: %v, Mimetype: %v\n", srv, FileContentType)
		//

		if !strings.Contains(FileContentType, "executable") {
			source, err := os.ReadFile(srv)
			if err != nil {
				d.Lgout.Println("Could not read", srv)
				fmt.Fprintln(w, "Could not read", srv)
				return
			} else {
				d.Lgout.Println("ext", filepath.Ext(srv))
			}

			switch filepath.Ext(srv) {
			case ".md":
				mdHtml, err := mdToHTML(source)
				if err != nil {
					d.Lgout.Println("Error converting markdown to html", err)
					source, err := os.ReadFile(srv)
					if err != nil {
						d.Lgout.Println("Could not read", srv)
						fmt.Fprintln(w, "Could not read", srv)
						w.Write(source)
						return
					} else {
						d.Lgout.Println("extension", filepath.Ext(srv))
					}
					w.Write(source)
				}
				w.Write(mdHtml)
			case ".php": //, ".pl", ".py", ".sh":
				//out, err := script.Echo("query").Exec(fmt.Sprintf("/usr/bin/php %s", srv)).Bytes()

				//cmd := exec.Command("/usr/bin/php-cgi", []string{"-q", srv, query}...)
				//cmd.Dir = filepath.Dir(srv)
				//out, err := cmd.Output()
				if d.PHP == "notfound" || d.PHP == "" {
					var err error
					d.PHP, err = exec.LookPath("php-cgi")
					if err != nil {
						d.Lgout.Println("php-cgi not found.")
						w.Write(source)
						break
					}
				}
				cmd := fmt.Sprintf("%s -q %s %s")
				out, err := script.Exec(cmd).Bytes()
				if err != nil {
					d.Lgout.Printf("Error executing php %v", err)
					w.Write([]byte(fmt.Sprintf("Error executing php %v", err)))
				} else {
					d.Lgout.Printf("Exec PHP: /usr/bin/php %s -- %s", srv, query)
					w.Write(out)
				}

			case ".pl":
				if d.Perl == "notfound" || d.Perl == "" {
					var err error
					d.Perl, err = exec.LookPath("perl")
					if err != nil {
						d.Lgout.Println("perl not found.")
						w.Write(source)
						break
					}
				}
				cmd := fmt.Sprintf("%s %s %s", d.Perl, srv, query)
				out, err := script.Exec(cmd).Bytes()
				if err != nil {
					d.Lgout.Printf("Error executing perl %v", err)
					w.Write([]byte(fmt.Sprintf("Error executing perl %v", err)))
				} else {
					d.Lgout.Printf("Exec Perl: %s %s -- %s", d.Perl, srv, query)
					w.Write(out)
				}

			case ".py":
				if d.Python == "notfound" || d.Python == "" {
					var err error
					d.Python, err = exec.LookPath("python")
					if err != nil {
						d.Lgout.Println("python not found.")
						w.Write(source)
						break
					}
				}
				cmd := fmt.Sprintf("%s %s %s", d.Python, srv, query)
				out, err := script.Exec(cmd).Bytes()
				if err != nil {
					d.Lgout.Printf("Error executing python %v", err)
					w.Write([]byte(fmt.Sprintf("Error executing python %v", err)))
				} else {
					d.Lgout.Printf("Exec Python: %s %s -- %s", d.Python, srv, query)
					w.Write(out)
				}
			case ".sh":
				if d.Bash == "notfound" || d.Bash == "" {
					var err error
					d.Bash, err = exec.LookPath("bash")
					if err != nil {
						d.Lgout.Println("bash not found.")
						w.Write(source)
						break
					}
				}
				cmd := fmt.Sprintf("%s %s %s", d.Bash, srv, query)
				out, err := script.Exec(cmd).Bytes()
				if err != nil {
					d.Lgout.Printf("Error executing bash %v", err)
					w.Write([]byte(fmt.Sprintf("Error executing bash %v", err)))
				} else {
					d.Lgout.Printf("Exec Bash: %s %s -- %s", d.Bash, srv, query)
					w.Write(out)
				}
			default:
				w.Write(source)
			}
			return
		} else {
			w.Header().Set("Content-Disposition", "attachment; filename="+dir.Name())
			w.Header().Set("Content-Type", FileContentType)
			w.Header().Set("Content-Length", strconv.FormatInt(dir.Size(), 10))
			openFile.Seek(0, 0)
			io.Copy(w, openFile) //'Copy' the file to the client
			http.Error(w, "File Sent Successfully", 200)
			return
		}
	}
	http.Error(w, "File Not Sent.", 500)
	return
}

func blackFolder(px int) (ico string) {
	img := "data:image/png;base64,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"
	//img := "iVBORw0KGgoAAAANSUhEUgAAAOEAAADhCAMAAAAJbSJIAAADXnpUWHRSYXcgcHJvZmlsZSB0eXBlIGV4aWYAAHja7ZdpcjQnDIb/c4ocAUkIwXFYq3KDHD8vDO7xjO2vHDtVWcrNNItaCNCjhh43/vh9ut9wUU7RBbUUc4weV8ghc0El+dt1K8mHne9rnBLtB7m7HjBEglJuzTiOfoFc7x0sHHl9lDtrx046hu6G9yVr5FU/eukYEr7J6bRdPv1KeLWcc3M7Zo/x53YwOKMr7Ak7HkLikYc1itzugpuQswR+kYgocpUPfOeu6pPzrtqT73w5cnl0hfPxKMQnHx056ZNcrmH4YUZ0H/nhQanXEG98N2dPc47b6kqI8FR0Z1EvS9k1KFa4Una3iGS4FXXbKSMlLLGBWAfNitQcZWJ4e1KgToUmjV02aphi4MGGkrmxbFkS48xtwwgr0WSTLN1JApUGagIxX3OhPW7e4zVKGLkTNJlgjNDjTXLvCb+SLkNzrtAl8unyFebFK6YxjUVu5dACEJrHp7r9u5N7FTf+FVgBQd1uTlhg8fVmoirdY0s2Z4Ge+uD87dUg68cAXISxFZMhAQEfSZQieWM2IvgxgU/BzFfYVxAgVe7kJtiIRMBJvMZGH6Oty8o3MbYWgFCJYkCTpQBWCIr4sZAQQ0VFg1PVqKZJs5YoMUSNMVpce1QxsWBq0cySZStJUkiaYrKUUk4lcxZsYZpjNpdTzrkUDFpguqB3gUYplavUULXGajXVXEtD+LTQtMVmLbXcSucuHa9/j91cTz33MmgglEYYOuKwkUYeZSLWpswwdcZpM808y0XtUH2kRk/kfk2NDrVFLGw9u1OD2OzFBK3tRBczEONAIG6LAAKaFzOfKARe5BYzn3ltVQxqpAtOp0UMBMMg1kkXuzu5X3JzGv4SN/6InFvo/g5ybqE75N5ye4daL/tEkQ1ovYXLp14mNjYojFQ4lXUmfbl03zXAOfiW2Fwa62jcT4AHqfR16il+75acYkTY4A310o36SJDnRtghaVTEUv6476dK9yyA8/Kworvt/afLbQix0Wb8p539fzI0wnsB4r5J/cfQj6EfQ/9hQzg+JGb81RF8zxflr5lz35yPn6XT+hvhuuDbEXPCd93rQwmtbx5Hny+jSMeWKbx2THaYDL58ha3SxBT/LZv/xOcK/ty6PwEP9GlLdg0WNgAAAHJQTFRFAAEBAAAA////e3t7uLi4x8fHSUlJ8fHx9fX1+/v7l5eX19fX+fn5cnJyzMzMvb296enpTk5OQEBAnZ2dCAgIioqKWFhYYGBgLCwsEhISICAgbW1tNjY25OTkkZGRZ2dn2trarq6uoqKiMTExg4ODIyMjh9TkVgAAAAF0Uk5TAEDm2GYAAAABYktHRACIBR1IAAAACXBIWXMAAC4jAAAuIwF4pT92AAAAB3RJTUUH5AEZASM2IaFKXgAAAUlJREFUeNrt3dttxSAAREG2B37ov9EUkIeubpwY1nM6GGFZiJ8dQ5IkSZIkSbqvlTdax/Dym9p9BxhzQfXArYm5qk19M+km5tLqgRsSk3LiSjsx7cKknZh2YghPJ4bwdGIICbd/u/lL4WwH7vGZEhISEhIS3i381/qF3yBT1qwXfj7HpN2YemLSbkw9MfXEWuGsF6ZfmH7hrBemXxjCEiLh2a16YQgJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCeuFpnQI925ZXisRznbgE2Ys7cme/SNtJj5runrUA58wIj/qgWXE8WX1wCLjGOXE8WOr6SLTd5Iv6SRJkiRJknRLH8jRj/KPpoERAAAAAElFTkSuQmCC"
	//ico = fmt.Sprintf("<img =`width:%vpx;height:%vpx;' id='base64image' src=data:image/png;base64, %v` >", px, px, img)
	ico = fmt.Sprintf(`<img height="%v" width="%v" src="%v" >`, px, px, img)
	return ico
}

func blackFile(px int) (ico string) {
	img := "data:image/png;base64,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"
	ico = fmt.Sprintf(`<img height="%v" width="%v" src="%v" >`, px, px, img)
	return ico
}

func toTitle(input string) string {
	return cases.Title(language.English).String(input)
}

func (d *Directory) getHeader() string {
	ret := fmt.Sprintf(`
		<title>%s</title>
		<
        <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.0.0/dist/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
        <script>
            function copyToClipboard(text) {
                navigator.clipboard.writeText(text).then(function() {
                    // Show feedback message
                    var feedback = document.getElementById('clipboard-feedback');
                    feedback.innerText = 'Copied to clipboard!';
                    feedback.style.display = 'block';
                    setTimeout(function() {
                        feedback.style.display = 'none';
                    }, 2000);
                }, function(err) {
                    console.error('Could not copy text: ', err);
                });
            }
        </script>
        <div id="clipboard-feedback" style="display: none; position: fixed; top: 10px; right: 10px; background-color: #28a745; color: white; padding: 10px; border-radius: 5px; z-index: 1000;"></div>
    `, d.Header)
	return ret
}

//func mdToHTML(md []byte) []byte {
//	md = markdown.NormalizeNewlines(md)
//	exts := parser.CommonExtensions // parser.OrderedListStart | parser.NoEmptyLineBeforeBlock
//	p := parser.NewWithExtensions(exts)
//	doc := markdown.Parse(md, p)
//
//	htmlFlags := mdhtml.Smartypants |
//		mdhtml.SmartypantsFractions |
//		mdhtml.SmartypantsDashes |
//		mdhtml.SmartypantsLatexDashes
//	htmlOpts := mdhtml.RendererOptions{
//		Flags: htmlFlags,
//	}
//	renderer := mdhtml.NewRenderer(htmlOpts)
//	html := markdown.Render(doc, renderer)
//	return html
//
//}

func mdToHTML(md []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := goldmark.Convert(md, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sortBy(sort, direction string) string {
	switch sort {
	case "name":
		return "name"
	case "size":
		return "size"
	case "date":
		return "date"
	default:
		return "name"
	}
}
