package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type routeTest interface {
	TestGet(t *testing.T)
	TestHead(t *testing.T)
	TestPost(t *testing.T)
	TestPut(t *testing.T)
	TestPatch(t *testing.T)
	TestDelete(t *testing.T)
	TestOptions(t *testing.T)
	TestTrace(t *testing.T)
}

type defaultRouteTest struct {
	Path string
}

func (routeTest *defaultRouteTest) TestPost(t *testing.T) {
	testMethodNotAllowed(t, http.MethodPost, routeTest.Path)
}

func (routeTest *defaultRouteTest) TestPut(t *testing.T) {
	testMethodNotAllowed(t, http.MethodPut, routeTest.Path)
}

func (routeTest *defaultRouteTest) TestPatch(t *testing.T) {
	testMethodNotAllowed(t, http.MethodPatch, routeTest.Path)
}

func (routeTest *defaultRouteTest) TestDelete(t *testing.T) {
	testMethodNotAllowed(t, http.MethodDelete, routeTest.Path)
}

func (routeTest *defaultRouteTest) TestOptions(t *testing.T) {
	testMethodNotAllowed(t, http.MethodOptions, routeTest.Path)
}

func (routeTest *defaultRouteTest) TestTrace(t *testing.T) {
	testMethodNotAllowed(t, http.MethodTrace, routeTest.Path)
}

func testMethodNotAllowed(t *testing.T, method, path string) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, "https://www.bobkidbob.com"+path, nil)
	t.Log(r.URL)
	handler(w, r)
	res := w.Result()
	if e, g := http.StatusMethodNotAllowed, res.StatusCode; e != g {
		t.Errorf("Expected HTTP status \"%v\", but got \"%v\"", e, g)
	}
}
