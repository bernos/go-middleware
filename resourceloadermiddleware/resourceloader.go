package resourceloadermiddleware

import (
	"context"
	"net/http"

	"github.com/bernos/go-middleware/middleware"
)

type key int

const ctxKey key = 0

// ResourceLoader loads a resource.
//
// Load should extract any parameters from the http.Request and attempt to load
// the appropriate resource. The error return parameter should be used to
// indicate that something went wrong attempting to load the resource, and will
// be passed to the configured error handler. If no resource can be found the
// ResourceLoader should return `nil, nil`. This will cause the configured not
// found handler to be run for the request.
type ResourceLoader interface {
	Load(*http.Request) (interface{}, error)
}

// ResourceLoaderFunc is an adapter that allows ordinary functions to be used as
// resource loaders.
type ResourceLoaderFunc func(*http.Request) (interface{}, error)

// Load calls f(r)
func (f ResourceLoaderFunc) Load(r *http.Request) (interface{}, error) {
	return f(r)
}

// NewContext wraps a loading resource in a new context
func NewContext(parent context.Context, resource interface{}) context.Context {
	return context.WithValue(parent, ctxKey, resource)
}

// FromContext retrieves the resource from a context
func FromContext(ctx context.Context) (interface{}, bool) {
	resource := ctx.Value(ctxKey)
	return resource, resource != nil
}

// UpdateRequest adds a resource to the request context
func UpdateRequest(r *http.Request, resource interface{}) *http.Request {
	ctx := NewContext(r.Context(), resource)
	return r.WithContext(ctx)
}

// FromRequest returns the loaded resource from the request, and a bool indicating
// whether any resource was found
func FromRequest(r *http.Request) (interface{}, bool) {
	return FromContext(r.Context())
}

// LoadResource creates http Middleware from a ReasourceLoader. The Middleware
// will pass incoming requests to the ResourceLoader and add the loaded resource
// to the request Context.
func LoadResource(loader ResourceLoader, options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shouldContinue := true
			resource, err := loader.Load(r)

			if err != nil {
				shouldContinue = cfg.errorHandler(err, w, r)
			} else if resource == nil {
				shouldContinue = cfg.notFoundHandler(w, r)
			}

			if shouldContinue {
				next.ServeHTTP(w, UpdateRequest(r, resource))
			}
		})
	}
}
