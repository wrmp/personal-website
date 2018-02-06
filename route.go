package main

import "net/http"

type route interface {
	Get(w http.ResponseWriter, r *http.Request)
	Head(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Connect(w http.ResponseWriter, r *http.Request)
	Options(w http.ResponseWriter, r *http.Request)
	Trace(w http.ResponseWriter, r *http.Request)
}

type defaultRoute struct {
	Allow string
}

func (d *defaultRoute) Post(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}

func (d *defaultRoute) Put(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}

func (d *defaultRoute) Patch(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}

func (d *defaultRoute) Delete(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}

func (d *defaultRoute) Connect(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}

func (d *defaultRoute) Options(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}

func (d *defaultRoute) Trace(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w, d.Allow)
}
