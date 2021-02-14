package server

import "net/http"

type home struct {
	*defaultRoute
}

func newHome() *home {
	return &home{&defaultRoute{Allow: http.MethodGet + ", " + http.MethodHead}}
}

func (*home) Get(w http.ResponseWriter, r *http.Request) {
	homeRedirect(w, r)
}

func (*home) Head(w http.ResponseWriter, r *http.Request) {
	homeRedirect(w, r)
}

func homeRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.linkedin.com/in/wrmp/", http.StatusSeeOther)
}
