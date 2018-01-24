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
    default:
        l := "https://www.bobkidbob.com" + url.EscapedPath()
        if url.RawQuery != "" {
            l += "?" + url.RawQuery
        }
        http.Redirect(w, r, l, http.StatusMovedPermanently)
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
