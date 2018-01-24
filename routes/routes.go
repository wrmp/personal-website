package routes

import (
    "net/http"
    "net/url"
)

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    url := r.URL
    switch url.Hostname() {
    case "www.bobkidbob.com":
        switch r.Method {
        case http.MethodHead:
            head(w, r, url)
        case http.MethodGet:
            head(w, r, url)
        default:
            w.WriteHeader(http.StatusNotImplemented)
        }
    case "bobkidbob.com":
        l := "https://www.bobkidbob.com" + pathAndQuery(url)
        http.Redirect(w, r, l, http.StatusMovedPermanently)
    case "bobkidbob.co.uk":
        l := "https://www.bobkidbob.co.uk" + pathAndQuery(url)
        http.Redirect(w, r, l, http.StatusMovedPermanently)
    case "bobkidbob.info":
        l := "https://www.bobkidbob.info" + pathAndQuery(url)
        http.Redirect(w, r, l, http.StatusMovedPermanently)
    case "bobkidbob.net":
        l := "https://www.bobkidbob.net" + pathAndQuery(url)
        http.Redirect(w, r, l, http.StatusMovedPermanently)
    case "bobkidbob.org":
        l := "https://www.bobkidbob.org" + pathAndQuery(url)
        http.Redirect(w, r, l, http.StatusMovedPermanently)
    default:
        l := "https://www.bobkidbob.com" + pathAndQuery(url)
        http.Redirect(w, r, l, http.StatusFound)
    }
}

func head(w http.ResponseWriter, r *http.Request, url *url.URL) {
    switch url.Path {
    case "/":
        http.Redirect(w, r, "https://www.linkedin.com/in/bobkidbob/", http.StatusFound)
    default:
        w.WriteHeader(http.StatusNotFound)
    }
}

func pathAndQuery(url *url.URL) string {
    s := url.EscapedPath()
    if url.RawQuery != "" {
        s += "?" + url.RawQuery
    }
    return s
}
