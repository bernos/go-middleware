package bodyparsermiddleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseBody(t *testing.T) {
	type foo struct {
		Name string
		Age  int
	}

	expect := &foo{
		Name: "Brendan",
		Age:  38,
	}

	var parsedBody interface{}

	parserWasCalled := false

	parser := func(d Decoder) (interface{}, error) {
		parserWasCalled = true
		data := &foo{}
		return data, d.Decode(data)
	}

	body, err := json.Marshal(expect)

	if err != nil {
		t.Error(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedBody, _ = FromRequest(r)
	})

	mw := ParseBody(parser)

	req, err := http.NewRequest("GET", "/", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := mw(next)
	handler.ServeHTTP(w, req)

	if !parserWasCalled {
		t.Error("Expected parser to be called")
	}

	if !reflect.DeepEqual(expect, parsedBody) {
		t.Errorf("Want %v, got %v", expect, parsedBody)
	}
}
