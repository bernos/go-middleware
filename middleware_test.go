package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "done")
	})
)

func a(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "A")
		n.ServeHTTP(w, r)
	})
}

func b(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "B")
		n.ServeHTTP(w, r)
	})
}

func c(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "C")
		n.ServeHTTP(w, r)
	})
}

func TestCompose(t *testing.T) {
	runTest(compose(a, b), "AB", t)
}

func TestVariadicCompose(t *testing.T) {
	runTest(Compose(a, b, c), "ABC", t)
}

func TestSingleCompose(t *testing.T) {
	runTest(Compose(a), "A", t)
}

func TestEmptyCompose(t *testing.T) {
	runTest(Compose(), "", t)
}

func TestCompositionChain(t *testing.T) {
	runTest(Compose().Compose(a).Compose(b), "AB", t)
}

func TestFold(t *testing.T) {
	runTest(fold(compose, a, []MiddlewareFunc{b, c}), "ABC", t)
}

func TestReduce(t *testing.T) {
	runTest(reduce(compose, []MiddlewareFunc{a, b}), "AB", t)
}

func TestEmptyReduce(t *testing.T) {
	runTest(reduce(compose, []MiddlewareFunc{}), "", t)
}

func runTest(h http.Handler, expected string, t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com/", nil)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	actual := w.Body.String()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
