package headermiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderMiddleware(t *testing.T) {
	want := map[string]string{
		"key_one": "value_one",
		"key_two": "value_two",
	}

	r, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	mw := New(want)
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, r)

	for k, v := range want {
		if w.Header().Get(k) != v {
			t.Errorf("Want %s, got %s", v, w.Header().Get(k))
		}
	}
}
