package middleware

import (
	"fmt"
	"github.com/bernos/go-middleware/testutils"
	"net/http"
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
	testutils.RunTest(Compose(a, b), "AB", t)
}

func TestComposeAll(t *testing.T) {
	testutils.RunTest(ComposeAll(a, b, c), "ABC", t)
}

func TestSingleComposeAll(t *testing.T) {
	testutils.RunTest(ComposeAll(a), "A", t)
}

func TestEmptyComposeAll(t *testing.T) {
	testutils.RunTest(ComposeAll(), "", t)
}

func TestCompositionChain(t *testing.T) {
	testutils.RunTest(ComposeAll().Compose(a).Compose(b), "AB", t)
}

func TestIdCompose(t *testing.T) {
	testutils.RunTest(Id().Compose(a).Compose(b), "AB", t)
}

func TestAssociativity(t *testing.T) {
	testutils.AssertEqual(Compose(a, Compose(b, c)), Compose(Compose(a, b), c), t)
}

func TestComposeAllComposeAll(t *testing.T) {
	abc := ComposeAll(a, b, c)
	cba := ComposeAll(c, b, a)
	abccba := ComposeAll(abc, cba)

	testutils.RunTest(abccba, "ABCCBA", t)
}

func TestFold(t *testing.T) {
	testutils.RunTest(fold(Compose, a, []Middleware{b, c}), "ABC", t)
}
