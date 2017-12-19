// This small program is just a small web server created in static mode
// in order to provide the smallest docker image possible

package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	// Def of flags
	portPtr     = flag.Int("p", 80, "The listening port")
	path        = flag.String("static", "/srv/http", "The path for the static files")
	routePrefix = flag.String("routePrefix", "", "The route path prefix for the static files")
	headerFlag  = flag.String("appendHeader", "", "HTTP response header, specified as `HeaderName:Value` that should be added to all responses.")
)

func parseHeaderFlag(headerFlag string) (string, string) {
	if len(headerFlag) == 0 {
		return "", ""
	}
	pieces := strings.SplitN(headerFlag, ":", 2)
	if len(pieces) == 1 {
		return pieces[0], ""
	}
	return pieces[0], pieces[1]
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		http.ServeFile(w, r, *path+"/_version.json")
		return
	}

	http.ServeFile(w, r, *path+"/_version.html")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("INCOMING REQUEST: %v %v %v %v [HEADERS]: %v", r.Proto, r.Method, r.Host, r.URL, r.Header)
	fileServer := http.FileServer(http.Dir(*path))
	// Extra headers.
	if len(*headerFlag) > 0 {
		header, headerValue := parseHeaderFlag(*headerFlag)
		if len(header) > 0 && len(headerValue) > 0 {
			w.Header().Set(header, headerValue)
			fileServer.ServeHTTP(w, r)
		} else {
			log.Println("appendHeader misconfigured; ignoring.")
		}
		return
	}
	fileServer.ServeHTTP(w, r)
}

func main() {

	flag.Parse()

	port := ":" + strconv.FormatInt(int64(*portPtr), 10)

	http.HandleFunc("/", handler)
	http.HandleFunc("/_version", versionHandler)
	if *routePrefix != "" {
		prefix := filepath.Clean(*routePrefix)
		if prefix[0:1] != "/" {
			prefix = "/" + prefix
		}
		if prefix[len(prefix)-1:] != "/" {
			prefix = prefix + "/"
		}
		http.HandleFunc(prefix+"_version", versionHandler)
		http.Handle(prefix, http.StripPrefix(prefix[:len(prefix)-1], http.FileServer(http.Dir(*path))))
	}

	log.Printf("Listening at 0.0.0.0%v...", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}
