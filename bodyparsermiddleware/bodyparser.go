package bodyparsermiddleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bernos/go-middleware/middleware"
)

type key int

const (
	ctxKey key = 0

	defaultMaxMemory = 32 << 20 // 32 mb
)

type Validatable interface {
	Validate() error
}

// Decoder decodes a request body
type Decoder interface {
	Decode(interface{}) error
}

// BodyParser uses Decoder to parse a request body
type BodyParser func(Decoder) (interface{}, error)

// FormParser parses a form into a struct
type FormParser func(*http.Request) (interface{}, error)

// NewContext adds a parsed request body to a context
func NewContext(parent context.Context, body interface{}) context.Context {
	return context.WithValue(parent, ctxKey, body)
}

// FromContext retrieves the parsed request body from the context
func FromContext(ctx context.Context) (interface{}, bool) {
	body := ctx.Value(ctxKey)
	return body, body != nil
}

// UpdateRequest adds a parsed request body to a context
func UpdateRequest(r *http.Request, body interface{}) *http.Request {
	ctx := NewContext(r.Context(), body)
	return r.WithContext(ctx)
}

// FromRequest returns the parsed body from the request, and a bool indicating
// whether any resource was found
func FromRequest(r *http.Request) (interface{}, bool) {
	return FromContext(r.Context())
}

func ParseJSONBody(parser BodyParser, options ...func(*options)) middleware.Middleware {
	options = append(options, WithJSONDecoder())
	return ParseBody(parser, options...)
}

// ParseBody creates an http middleware from a BodyParser. The middleware will use the
// BodyParser to parse the request body and add it to the request context. If the BodyParser
// returns an error then the configured error handler will be called
func ParseBody(parser BodyParser, options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			shouldContinue := true

			if r.Body != nil {
				decoder := cfg.decoder(r)
				body, err := parser(decoder)

				if err != nil {
					shouldContinue = cfg.errorHandler(err, w, r)
				} else {
					err := cfg.validator(body)

					if err != nil {
						shouldContinue = cfg.errorHandler(err, w, r)
					}
				}

				r = UpdateRequest(r, body)
			}

			if shouldContinue {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// ParseForm creates an http middleware from a FormParser. The middleware will use the
// FormParser to parse the form in the request and add it to the request context. If the FormParser
// returns an error then the configured error handler will be called
func ParseForm(parser FormParser, options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shouldContinue := true

			err := r.ParseMultipartForm(cfg.maxMemory)

			if err != nil {
				shouldContinue = cfg.errorHandler(err, w, r)
			} else {
				body, err := parser(r)

				if err != nil {
					shouldContinue = cfg.errorHandler(err, w, r)
				} else {
					err := cfg.validator(body)

					if err != nil {
						shouldContinue = cfg.errorHandler(err, w, r)
					}
				}

				r = UpdateRequest(r, body)
			}

			if shouldContinue {
				next.ServeHTTP(w, r)
			}
		})
	}
}

type options struct {
	decoder      func(*http.Request) Decoder
	errorHandler func(error, http.ResponseWriter, *http.Request) bool
	maxMemory    int64
	validator    func(interface{}) error
}

func defaultOptions() *options {
	return &options{
		decoder:      jsonDecoder,
		errorHandler: defaultErrorHandler,
		maxMemory:    defaultMaxMemory,
		validator:    defaultValidator,
	}
}

func WithMaxMemory(x int64) func(*options) {
	return func(o *options) {
		o.maxMemory = x
	}
}

func WithJSONDecoder() func(*options) {
	return func(o *options) {
		o.decoder = jsonDecoder
	}
}

func WithErrorHandler(h func(error, http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.errorHandler = h
	}
}

func jsonDecoder(r *http.Request) Decoder {
	return json.NewDecoder(r.Body)
}

func defaultErrorHandler(err error, w http.ResponseWriter, r *http.Request) bool {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return false
}

func defaultValidator(x interface{}) error {
	v, ok := x.(Validatable)

	if ok {
		return v.Validate()
	}

	return nil
}
