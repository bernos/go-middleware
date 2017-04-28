package viewmiddleware

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/bernos/go-middleware/bodyparsermiddleware"
	"github.com/bernos/go-middleware/errormiddleware"
	"github.com/bernos/go-middleware/middleware"
)

type View struct {
	Template *template.Template
	Model    interface{}
}

func RenderView(options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := FromRequest(r)

			if v != nil {
				t := v.Template

				if t == nil {
					t = cfg.defaultTemplate
				}

				if t == nil {
					r = errormiddleware.UpdateRequest(r, fmt.Errorf("No template specified for this view"), http.StatusInternalServerError)
				} else {
					vm := struct {
						Model           interface{}
						ValidationError error
					}{
						Model:           v.Model,
						ValidationError: bodyparsermiddleware.Validate(r),
					}

					err := t.Execute(w, vm)

					if err != nil {
						r = errormiddleware.UpdateRequest(r, err, http.StatusInternalServerError)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func HandlerFunc(fn func(http.ResponseWriter, *http.Request) (*View, error)) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v, err := fn(w, r)

			if err != nil {
				r = errormiddleware.UpdateRequest(r, err, http.StatusInternalServerError)
			}

			if v != nil {
				r = UpdateRequest(r, v)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func templateOrDefault(v *View, r *http.Request, fn func(*http.Request) *template.Template) *template.Template {
	if v == nil || v.Template == nil {
		return fn(r)
	}
	return v.Template
}

func modelOrDefault(v *View, r *http.Request, fn func(*http.Request) interface{}) interface{} {
	if v == nil || v.Model == nil {
		return fn(r)
	}
	return v.Model
}
