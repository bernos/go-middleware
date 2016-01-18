package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func RunTest(h http.Handler, expected string, t *testing.T) {
	actual := GetResponse(h)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func AssertEqual(a http.Handler, b http.Handler, t *testing.T) {
	ra := GetResponse(a)
	rb := GetResponse(b)

	if ra != rb {
		t.Errorf("Response %s does not equal response %s", ra, rb)
	}
}

func GetResponse(h http.Handler) string {
	r, _ := http.NewRequest("GET", "http://example.com/", nil)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	return w.Body.String()
}
