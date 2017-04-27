package viewmiddleware

import (
	"html/template"
	"net/http"

	"github.com/bernos/go-middleware/bodyparsermiddleware"
	"github.com/bernos/go-middleware/middleware"
)

var defaultTemplate = template.Must(template.New("_default").Parse(`This is the default template`))

type View struct {
	Template *template.Template
	Model    interface{}
}

func RenderView(defaultTemplate *template.Template, options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions(defaultTemplate)

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := FromRequest(r)
			t := templateOrDefault(v, r, cfg.templateProvider)
			m := modelOrDefault(v, r, cfg.viewModelProvider)

			vm := struct {
				Model           interface{}
				Error           error
				ValidationError error
			}{
				Model:           m,
				Error:           nil,
				ValidationError: bodyparsermiddleware.Validate(r),
			}

			shouldContinue := true
			err := t.Execute(w, vm)

			if err != nil {
				shouldContinue = cfg.errorHandler(err, w, r)
			}

			if shouldContinue {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func BuildView(fn func(*http.Request) *View) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, UpdateRequest(r, fn(r)))
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
