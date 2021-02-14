package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
		if e, g := http.StatusPermanentRedirect, res.StatusCode; e != g {
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
		if e, g := http.StatusTemporaryRedirect, res.StatusCode; e != g {
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
	testRoute(t, newHomeTest("/"))
}

// Test if non-existent page returns HTTP 404 Not Found error.
func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "https://www.bobkidbob.com/thispagedoesnotexist", nil)
	handler(w, r)
	res := w.Result()
	if e, g := http.StatusNotFound, res.StatusCode; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
}

// Test if unimplemented HTTP methods return 501 Not Implemented error.
func TestNotImplemented(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("FAKEMETHOD", "https://www.bobkidbob.com/", nil)
	handler(w, r)
	res := w.Result()
	if e, g := http.StatusNotImplemented, res.StatusCode; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
}

// Test if a URN redirects to HTTPS.
func forceHTTPS(t *testing.T, urn string) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://"+urn, nil)
	handler(w, r)
	res := w.Result()
	if e, g := http.StatusPermanentRedirect, res.StatusCode; e != g {
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

func TestPathAndQuery(t *testing.T) {
	if e, g := "/test?foo=1&bar=2", pathAndQuery(&url.URL{Path: "/test", RawQuery: "foo=1&bar=2"}); e != g {
		t.Errorf("Expected path and query \"%v\", but got \"%v\"", e, g)
	}
}

func testRoute(t *testing.T, route routeTest) {
	t.Run("Get", route.TestGet)
	t.Run("Head", route.TestHead)
	t.Run("Post", route.TestPost)
	t.Run("Put", route.TestPut)
	t.Run("Patch", route.TestPatch)
	t.Run("Delete", route.TestDelete)
	t.Run("Options", route.TestOptions)
	t.Run("Trace", route.TestTrace)
}
