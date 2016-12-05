package loggingmiddleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	r := newRequest(t, "GET", "/", nil)

	test := func(info RequestInfo) {
		if info.Request != r {
			t.Errorf("want %#v, got %#v", r, info.Request)
		}
	}

	runTest(t, test, r, nil)
}

func TestStatusCode(t *testing.T) {
	r := newRequest(t, "GET", "/", nil)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})

	test := func(info RequestInfo) {
		if info.Status != http.StatusBadGateway {
			t.Errorf("want %d, got %d", http.StatusBadGateway, info.Status)
		}
	}

	runTest(t, test, r, h)
}

func TestLatency(t *testing.T) {
	r := newRequest(t, "GET", "/", nil)

	test := func(info RequestInfo) {
		if info.Latency == 0 {
			t.Errorf("Expected latency, but got %t", info.Latency)
		}
	}

	runTest(t, test, r, nil)
}

func runTest(t *testing.T, fn func(RequestInfo), r *http.Request, h http.Handler) {
	var info RequestInfo

	log := func(i RequestInfo) {
		info = i
	}

	mw := New(WithLogger(log))
	w := httptest.NewRecorder()

	if h == nil {
		mw.ServeHTTP(w, r)
	} else {
		mw(h).ServeHTTP(w, r)
	}

	fn(info)
}

func newRequest(t *testing.T, method string, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, url, body)

	if err != nil {
		t.Fatalf("Failed to create request. %s", err)
	}

	return r
}
