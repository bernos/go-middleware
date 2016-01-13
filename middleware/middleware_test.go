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
	runTest(Compose(a, b), "AB", t)
}

func TestComposeAll(t *testing.T) {
	runTest(ComposeAll(a, b, c), "ABC", t)
}

func TestSingleComposeAll(t *testing.T) {
	runTest(ComposeAll(a), "A", t)
}

func TestEmptyComposeAll(t *testing.T) {
	runTest(ComposeAll(), "", t)
}

func TestCompositionChain(t *testing.T) {
	runTest(ComposeAll().Compose(a).Compose(b), "AB", t)
}

func TestIdCompose(t *testing.T) {
	runTest(Id().Compose(a).Compose(b), "AB", t)
}

func TestAssociativity(t *testing.T) {
	assertEqual(Compose(a, Compose(b, c)), Compose(Compose(a, b), c), t)
}

func TestComposeAllComposeAll(t *testing.T) {
	abc := ComposeAll(a, b, c)
	cba := ComposeAll(c, b, a)
	abccba := ComposeAll(abc, cba)

	runTest(abccba, "ABCCBA", t)
}

func TestFold(t *testing.T) {
	runTest(fold(Compose, a, []MiddlewareFunc{b, c}), "ABC", t)
}

func runTest(h http.Handler, expected string, t *testing.T) {
	actual := getResponse(h)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func assertEqual(a http.Handler, b http.Handler, t *testing.T) {
	ra := getResponse(a)
	rb := getResponse(b)

	if ra != rb {
		t.Errorf("Response %s does not equal response %s", ra, rb)
	}
}

func getResponse(h http.Handler) string {
	r, _ := http.NewRequest("GET", "http://example.com/", nil)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	return w.Body.String()
}
