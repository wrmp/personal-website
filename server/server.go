package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

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
	http.MethodOptions: route.Options,
	http.MethodTrace:   route.Trace,
}
var routes = map[string]route{
	"/": newHome(),
}

func init() {
	loadPage(&pages.MethodNotAllowed, "method-not-allowed.html")
	loadPage(&pages.NotFound, "not-found.html")
	loadPage(&pages.NotImplemented, "not-implemented.html")
	http.HandleFunc("/", handler)
}

func loadPage(page *[]byte, file string) {
	path, err := filepath.Abs(filepath.Join("frontend", file))
	handleError(err)
	*page, err = ioutil.ReadFile(path)
	handleError(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	if r.Header.Get("X-Forwarded-Proto") != "https" {
		permanentRedirect(w, r, url, r.Host)
		return
	}
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	hostname := r.Host
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
		var rou route
		if rou, exists = routes[url.Path]; exists {
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
