package main

import (
    "google.golang.org/appengine"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
)

var dev bool
var pages struct {
    MethodNotAllowed []byte
    NotFound []byte
    NotImplemented []byte
}

func init() {
    dev = appengine.IsDevAppServer()
    f, err := ioutil.ReadFile("frontend/method-not-allowed.html")
	if err != nil {
		log.Fatal(err)
	}
    pages.MethodNotAllowed = f
    f, err = ioutil.ReadFile("frontend/not-found.html")
	if err != nil {
		log.Fatal(err)
	}
    pages.NotFound = f
    f, err = ioutil.ReadFile("frontend/not-implemented.html")
	if err != nil {
		log.Fatal(err)
	}
    pages.NotImplemented = f
    http.HandleFunc("/", handler)
}

func main() {
    appengine.Main()
}

func handler(w http.ResponseWriter, r *http.Request) {
    if dev {
        serve(w, r, r.URL)
        return
    }
    url := r.URL
    switch url.Hostname() {
    case "www.bobkidbob.com":
        serve(w, r, url)
    case "bobkidbob.com":
        permanentRedirect(w, r, url, "www.bobkidbob.com")
    case "bobkidbob.co.uk":
        permanentRedirect(w, r, url, "www.bobkidbob.co.uk")
    case "bobkidbob.info":
        permanentRedirect(w, r, url, "www.bobkidbob.info")
    case "bobkidbob.net":
        permanentRedirect(w, r, url, "www.bobkidbob.net")
    case "bobkidbob.org":
        permanentRedirect(w, r, url, "www.bobkidbob.org")
    default:
        http.Redirect(w, r, "https://www.bobkidbob.com" + pathAndQuery(url), http.StatusTemporaryRedirect)
    }
}

func serve(w http.ResponseWriter, r *http.Request, url *url.URL) {
    switch r.Method {
    case http.MethodGet:
        head(w, r, url)
        get(w, url)
    case http.MethodHead:
        head(w, r, url)
    case http.MethodPost:
        methodNotAllowed(w)
    case http.MethodPut:
        methodNotAllowed(w)
    case http.MethodPatch:
        methodNotAllowed(w)
    case http.MethodDelete:
        methodNotAllowed(w)
    case http.MethodConnect:
        methodNotAllowed(w)
    case http.MethodOptions:
        methodNotAllowed(w)
    case http.MethodTrace:
        methodNotAllowed(w)
    default:
        w.WriteHeader(http.StatusNotImplemented)
        w.Write(pages.NotImplemented)
    }
}

func head(w http.ResponseWriter, r *http.Request, url *url.URL) {
    switch url.Path {
    case "/":
        http.Redirect(w, r, "https://www.linkedin.com/in/bobkidbob/", http.StatusSeeOther)
    default:
        w.WriteHeader(http.StatusNotFound)
    }
}

func get(w http.ResponseWriter, url *url.URL) {
    switch url.Path {
    case "/":
        return
    default:
        w.Write(pages.NotFound)
    }
}

func methodNotAllowed(w http.ResponseWriter) {
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
    http.Redirect(w, r, "https://" + host + pathAndQuery(url), http.StatusPermanentRedirect)
}
