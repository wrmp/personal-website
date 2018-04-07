package server

import (
	"net/http/httptest"
	"testing"
)

var hosts = [...]string{"bobkidbob.com", "bobkidbob.co.uk", "bobkidbob.info", "bobkidbob.net", "bobkidbob.org"}

// Test if HTTP requests are redirected to HTTPS.
func TestForceHTTPS(t *testing.T) {
	for i, l := 0, len(hosts); i < l; i++ {
		h := hosts[i]
		forceHTTPS(t, h+"/")
		forceHTTPS(t, "www."+h+"/")
	}
}

// Test if naked domain (e.g. bobkidbob.com) is redirected to www subdomain (e.g. www.bobkidbob.com).
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

// Test if all alias domains redirect to bobkidbob.com.
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

// Test homepage.
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

// Test if non-existent page returns HTTP 404 Not Found error.
func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "https://www.bobkidbob.com/thispagedoesnotexist", nil)
	handler(w, r)
	res := w.Result()
	if e, g := "404 Not Found", res.Status; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
}

// Test if unimplemented HTTP methods return 501 Not Implemented error.
func TestNotImplemented(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("FAKEMETHOD", "https://www.bobkidbob.com/", nil)
	handler(w, r)
	res := w.Result()
	if e, g := "501 Not Implemented", res.Status; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
}

// Test development mode.
func TestDev(t *testing.T) {
	dev = true
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost/", nil)
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

// Test if a URN redirects to HTTPS.
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
