package main

import (
	"net/http/httptest"
	"testing"
)

var hosts = [...]string{"bobkidbob.com", "bobkidbob.co.uk", "bobkidbob.info", "bobkidbob.net", "bobkidbob.org"}

func TestForceHTTPS(t *testing.T) {
	for i, l := 0, len(hosts); i < l; i++ {
		h := hosts[i]
		forceHTTPS(t, h+"/")
		forceHTTPS(t, "www."+h+"/")
	}
}

func TestNaked(t *testing.T) {
	for i, l := 0, len(hosts); i < l; i++ {
		h := hosts[i]
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "https://"+h, nil)
		handler(w, r)
		res := w.Result()
		if e, g := "308 Permanent Redirect", res.Status; e != g {
			t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
		}
		l, err := res.Location()
		if err != nil {
			t.Error(err)
		}
		if e, g := "https://www."+h, l.String(); e != g {
			t.Errorf("Expected location header \"%v\", but got \"%v\"", e, g)
		}
	}
}

func TestAliases(t *testing.T) {
	for i, l := 1, len(hosts); i < l; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "https://www."+hosts[i]+"/", nil)
		handler(w, r)
		res := w.Result()
		if e, g := "307 Temporary Redirect", res.Status; e != g {
			t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
		}
		l, err := res.Location()
		if err != nil {
			t.Error(err)
		}
		if e, g := "https://www.bobkidbob.com/", l.String(); e != g {
			t.Errorf("Expected location header \"%v\", but got \"%v\"", e, g)
		}
	}
}

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://www.bobkidbob.com/", nil)
	handler(w, r)
	res := w.Result()
	if e, g := "303 See Other", res.Status; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
	l, err := res.Location()
	if err != nil {
		t.Error(err)
	}
	if e, g := "https://www.linkedin.com/in/bobkidbob/", l.String(); e != g {
		t.Errorf("Expected location header \"%v\", but got \"%v\"", e, g)
	}
}

func forceHTTPS(t *testing.T, urn string) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://"+urn, nil)
	handler(w, r)
	res := w.Result()
	if e, g := "308 Permanent Redirect", res.Status; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
	l, err := res.Location()
	if err != nil {
		t.Error(err)
	}
	if e, g := "https://"+urn, l.String(); e != g {
		t.Errorf("Expected location header \"%v\", but got \"%v\"", e, g)
	}
}
