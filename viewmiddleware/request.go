package viewmiddleware

import (
	"context"
	"net/http"
)

type key int

const (
	viewModelKey key = 0
	templateKey  key = 1
	viewKey      key = 2
)

// NewContext adds a view to context
func NewContext(parent context.Context, view *View) context.Context {
	return context.WithValue(parent, viewKey, view)
}

// FromContext retrieves the view from context
func FromContext(ctx context.Context) *View {
	view, ok := ctx.Value(viewKey).(*View)

	if ok {
		return view
	}

	return nil
}

// UpdateRequest sets the view for a request
func UpdateRequest(r *http.Request, view *View) *http.Request {
	return r.WithContext(NewContext(r.Context(), view))
}

// FromRequest retrieves the view from a request
func FromRequest(r *http.Request) *View {
	return FromContext(r.Context())
}

// func GetViewModel(r *http.Request, defaultModel interface{}) interface{} {
// 	model := r.Context().Value(viewModelKey)

// 	if model == nil {
// 		model = defaultModel
// 	}

// 	return model
// }

// func GetTemplate(r *http.Request, defaultTemplate *template.Template) *template.Template {
// 	t, ok := r.Context().Value(templateKey).(*template.Template)

// 	if !ok {
// 		t = defaultTemplate
// 	}

// 	return t
// }

// type Request struct {
// 	*http.Request
// }

// func NewRequest(r *http.Request) *Request {
// 	return &Request{r}
// }

// func (r *Request) WithViewModel(m interface{}) *Request {
// 	return &Request{RequestWithViewModel(r.HTTPRequest(), m)}
// }

// func (r *Request) WithTemplate(t *template.Template) *Request {
// 	return &Request{RequestWithTemplate(r.HTTPRequest(), t)}
// }

// func (r *Request) GetViewModel(defaultValue interface{}) interface{} {
// 	return GetViewModel(r.HTTPRequest(), defaultValue)
// }

// func (r *Request) GetTemplate(defaultTemplate *template.Template) *template.Template {
// 	return GetTemplate(r.HTTPRequest(), defaultTemplate)
// }

// func (r *Request) HTTPRequest() *http.Request {
// 	return r.Request
// }
