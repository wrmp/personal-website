package server

import "net/http"

type home struct {
	defaultRoute
}

func (h *home) Get(w http.ResponseWriter, r *http.Request) {
	h.Head(w, r)
}

func (*home) Head(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.linkedin.com/in/bobkidbob/", http.StatusSeeOther)
}
