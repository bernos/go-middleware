package middlewarec

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func appendToContext(c context.Context, v string) context.Context {
	newValue := fmt.Sprintf("%s%s", getFromContext(c), v)
	return context.WithValue(c, 0, newValue)
}

func getFromContext(c context.Context) string {
	v, ok := c.Value(0).(string)

	if !ok {
		return ""
	}

	return v
}

func a(next Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		ctx := appendToContext(c, "A")
		next.ServeHTTP(ctx, w, r)
	})
}

func b(next Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		ctx := appendToContext(c, "B")
		next.ServeHTTP(ctx, w, r)
	})
}

func c(next Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		ctx := appendToContext(c, "C")
		next.ServeHTTP(ctx, w, r)
	})
}

var handler HandlerFunc = HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", getFromContext(c))
})

func TestCompose(t *testing.T) {
	runTest(Compose(a, b)(handler).AsHttpHandler(), "AB", t)
}

func TestComposeAll(t *testing.T) {
	runTest(ComposeAll(a, b, c)(handler).AsHttpHandler(), "ABC", t)
}

func TestSingleComposeAll(t *testing.T) {
	runTest(ComposeAll(a)(handler).AsHttpHandler(), "A", t)
}

func TestCompositionChain(t *testing.T) {
	runTest(ComposeAll().Compose(a).Compose(b).Compose(c)(handler).AsHttpHandler(), "ABC", t)
}

func TestIdCompose(t *testing.T) {
	runTest(Id().Compose(a).Compose(b)(handler).AsHttpHandler(), "AB", t)
}

func TestAssociativity(t *testing.T) {
	ma := Compose(a, Compose(b, c))
	mb := Compose(Compose(a, b), c)

	assertEqual(ma(handler).AsHttpHandler(), mb(handler).AsHttpHandler(), t)
}

func TestComposeAllComposeALl(t *testing.T) {
	abc := ComposeAll(a, b, c)
	cba := ComposeAll(c, b, a)
	abccba := ComposeAll(abc, cba)

	runTest(abccba(handler).AsHttpHandler(), "ABCCBA", t)
}

func TestFold(t *testing.T) {
	runTest(fold(Compose, a, []Middleware{b, c})(handler).AsHttpHandler(), "ABC", t)
}

func runTest(h http.Handler, expected string, t *testing.T) {
	actual := getResponse(h)

	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func assertEqual(a http.Handler, b http.Handler, t *testing.T) {
	ra := getResponse(a)
	rb := getResponse(b)

	if ra != rb {
		t.Errorf("Response %s is not euqal to response %s", ra, rb)
	}
}

func getResponse(h http.Handler) string {
	r, _ := http.NewRequest("GET", "http://example.com", nil)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	return w.Body.String()
}
