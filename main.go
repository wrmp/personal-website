package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"google.golang.org/appengine"
)

var dev = appengine.IsDevAppServer()
var pages struct {
	MethodNotAllowed []byte
	NotFound         []byte
	NotImplemented   []byte
}
var methods = map[string]func(route, http.ResponseWriter, *http.Request){
	http.MethodGet:     route.Get,
	http.MethodHead:    route.Head,
	http.MethodPost:    route.Post,
	http.MethodPut:     route.Put,
	http.MethodPatch:   route.Patch,
	http.MethodDelete:  route.Delete,
	http.MethodConnect: route.Connect,
	http.MethodOptions: route.Options,
	http.MethodTrace:   route.Trace,
}
var routes = map[string]route{
	"/": &home{defaultRoute{Allow: http.MethodGet + ", " + http.MethodHead}},
}

func main() {
	var err error
	pages.MethodNotAllowed, err = ioutil.ReadFile("frontend/method-not-allowed.html")
	handleError(err)
	pages.NotFound, err = ioutil.ReadFile("frontend/not-found.html")
	handleError(err)
	pages.NotImplemented, err = ioutil.ReadFile("frontend/not-implemented.html")
	handleError(err)
	http.HandleFunc("/", handler)
	appengine.Main()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if dev {
		serve(w, r, r.URL)
		return
	}
	url := r.URL
	if url.Scheme != "https" {
		permanentRedirect(w, r, url, url.Hostname())
		return
	}
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	hostname := url.Hostname()
	switch hostname {
	case "www.bobkidbob.com":
		serve(w, r, url)
	case "bobkidbob.com", "bobkidbob.co.uk", "bobkidbob.info", "bobkidbob.net", "bobkidbob.org":
		permanentRedirect(w, r, url, "www."+hostname)
	default:
		http.Redirect(w, r, "https://www.bobkidbob.com"+pathAndQuery(url), http.StatusTemporaryRedirect)
	}
}

func serve(w http.ResponseWriter, r *http.Request, url *url.URL) {
	if method, exists := methods[r.Method]; exists {
		if rou, exists := routes[url.Path]; exists {
			method(rou, w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write(pages.NotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write(pages.NotImplemented)
	}
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(pages.MethodNotAllowed)
}

func pathAndQuery(url *url.URL) string {
	s := url.EscapedPath()
	if url.RawQuery != "" {
		s += "?" + url.RawQuery
	}
	return s
}

func permanentRedirect(w http.ResponseWriter, r *http.Request, url *url.URL, host string) {
	http.Redirect(w, r, "https://"+host+pathAndQuery(url), http.StatusPermanentRedirect)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
