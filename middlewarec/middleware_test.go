package middlewarec

import (
	"fmt"
	"github.com/bernos/go-middleware/testutils"
	"golang.org/x/net/context"
	"net/http"
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
		next.ServeHTTPC(ctx, w, r)
	})
}

func b(next Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		ctx := appendToContext(c, "B")
		next.ServeHTTPC(ctx, w, r)
	})
}

func c(next Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		ctx := appendToContext(c, "C")
		next.ServeHTTPC(ctx, w, r)
	})
}

var handler HandlerFunc = HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", getFromContext(c))
})

func TestCompose(t *testing.T) {
	testutils.RunTest(Compose(a, b)(handler).AsHttpHandler(), "AB", t)
}

func TestComposeAll(t *testing.T) {
	testutils.RunTest(ComposeAll(a, b, c)(handler).AsHttpHandler(), "ABC", t)
}

func TestSingleComposeAll(t *testing.T) {
	testutils.RunTest(ComposeAll(a)(handler).AsHttpHandler(), "A", t)
}

func TestCompositionChain(t *testing.T) {
	testutils.RunTest(ComposeAll().Compose(a).Compose(b).Compose(c)(handler).AsHttpHandler(), "ABC", t)
}

func TestIdCompose(t *testing.T) {
	testutils.RunTest(Id().Compose(a).Compose(b)(handler).AsHttpHandler(), "AB", t)
}

func TestAssociativity(t *testing.T) {
	ma := Compose(a, Compose(b, c))
	mb := Compose(Compose(a, b), c)

	testutils.AssertEqual(ma(handler).AsHttpHandler(), mb(handler).AsHttpHandler(), t)
}

func TestComposeAllComposeALl(t *testing.T) {
	abc := ComposeAll(a, b, c)
	cba := ComposeAll(c, b, a)
	abccba := ComposeAll(abc, cba)

	testutils.RunTest(abccba(handler).AsHttpHandler(), "ABCCBA", t)
}

func TestFold(t *testing.T) {
	testutils.RunTest(fold(Compose, a, []Middleware{b, c})(handler).AsHttpHandler(), "ABC", t)
}

func TestComposeWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), 0, "Z")

	testutils.RunTest(ComposeAll(a, b, c)(handler).AsHttpHandlerWithContext(ctx), "ZABC", t)

}
