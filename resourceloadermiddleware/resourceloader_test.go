package resourceloadermiddleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewContext(t *testing.T) {
	want := "foo"
	ctx := NewContext(context.Background(), want)
	got := ctx.Value(ctxKey)

	if ctx.Value(ctxKey) != want {
		t.Errorf("Want %s, got %s", want, got)
	}
}

func TestFromContext(t *testing.T) {
	want := "foo"
	ctx := NewContext(context.Background(), want)
	got, ok := FromContext(ctx)

	if got != want {
		t.Errorf("Want %s, got %s", want, got)
	}

	if !ok {
		t.Error("Expected ok")
	}
}

func TestFromEmptyContext(t *testing.T) {
	got, ok := FromContext(context.Background())

	if ok != false {
		t.Error("Expected ok to be false")
	}

	if got != nil {
		t.Errorf("Expected nil resource, got %v", got)
	}
}

func TestResourceLoaderFunc(t *testing.T) {
	loader := ResourceLoaderFunc(func(r *http.Request) (interface{}, error) {
		return "foo", nil
	})

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	r, err := loader.Load(req)

	if err != nil {
		t.Error("Expected err to be nil")
	}

	if r != "foo" {
		t.Errorf("Want %s, got %s", "foo", r)
	}
}

func TestLoadResource(t *testing.T) {
	var loadedResource interface{}

	loaderWasCalled := false

	loader := ResourceLoaderFunc(func(r *http.Request) (interface{}, error) {
		loaderWasCalled = true
		return "foo", nil
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loadedResource, _ = FromRequest(r)
	})

	mw := LoadResource(loader)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := mw(next)
	handler.ServeHTTP(w, req)

	if !loaderWasCalled {
		t.Error("Expected load to be called")
	}

	if loadedResource != "foo" {
		t.Error("Expected foo")
	}
}
