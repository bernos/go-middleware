# go-middleware
[![Build Status](https://travis-ci.org/bernos/go-middleware.svg?branch=master)](https://travis-ci.org/bernos/go-middleware)&nbsp;[![GoDoc](https://godoc.org/github.com/bernos/go-middleware/middleware?status.svg)](https://godoc.org/github.com/bernos/go-middleware/middleware)

Composable http middleware for golang. Middleware is defined as

```go
type Middleware func(http.Handler) http.Handler
```

This package provides some composition functions to make building middleware "stacks" as simple as

```go
// Compose takes two middleware funcs and returns... a new middleware func!
middlewareOneThenTwo := Compose(middlewareOne, middlewareTwo)

// And stacks can be composed with themselves. We can use ComposeAll to compose more than two middleware funcs
middlewareFive := ComposeAll(middlewareOneThenTwo, middlewareThree, middlewareFour)

// Finally, we wrap a regular http.Handler with our middleware stack and get back a regular http.Handler
finalHandler := middlewareFive(myHandler)
```

# Installation
For the regular, garden variety, `http.Handler` compatible middleware, use

```golang
go get github.com/bernos/go-middleware/middleware
```

To support passing a `Context` from `golang.org/x/net/context` between middleware and handlers, use

```golang
go get github.com/bernos/go-middleware/middlewarec
```

# Examples
## Basic middleware composition
```go
package main

import (
	"fmt"
	"github.com/bernos/go-middleware/middleware"
	"net/http"
)

// Middleware is any function that takes an http.Handler argument and returns
// another http.Handler
func MyFirstMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "First Middleware\n")
		h.ServeHTTP(w, r)
	})
}

func MySecondMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Second Middleware\n")
		h.ServeHTTP(w, r)
	})
}

func MyThirdMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Third Middleware\n")
		h.ServeHTTP(w, r)
	})
}

func MyFourthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Fourth Middleware\n")
		h.ServeHTTP(w, r)
	})
}

func main() {
	// Create a couple of middleware stacks using the Compose() function
	stackOne := middleware.Compose(MyFirstMiddleware, MySecondMiddleware)
	stackTwo := middleware.Compose(MyThirdMiddleware, MyFourthMiddleware)

	// Middleware compositions can also be composed with themselves
	stackThree := middleware.Compose(stackOne, stackTwo)

	// Alternatively, we can use method chaining, if we prefer
	stackFour := middleware.Middleware(MyFirstMiddleware).
		Compose(MySecondMiddleware).
		Compose(MyThirdMiddleware).
		Compose(stackThree)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Handler")
	})

	http.Handle("/one", stackOne(handler))
	http.Handle("/two", stackTwo(handler))
	http.Handle("/three", stackThree(handler))
	http.Handle("/four", stackFour(handler))

	http.ListenAndServe(":8080", nil)
}
```

## Supporting a golang.org/x/net/context.Context

```go
package main

import (
	"fmt"
	"github.com/bernos/go-middleware/middlewarec"
	"golang.org/x/net/context"
	"net/http"
)

// Appends to the value stored in our context
func appendToContextValue(c context.Context, v string) context.Context {
	key := 0
	currentValue, ok := c.Value(key).(string)

	if !ok {
		currentValue = ""
	}

	newValue := fmt.Sprintf("%s%s", currentValue, v)

	return context.WithValue(c, key, newValue)
}

// Context aware middleware funcs need to accept and return the custom middlewarec.Handler type
func MyFirstMiddleware(h middlewarec.Handler) middlewarec.Handler {
	return middlewarec.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTPC(appendToContextValue(c, "First\n"), w, r)
	})
}

func MySecondMiddleware(h middlewarec.Handler) middlewarec.Handler {
	return middlewarec.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTPC(appendToContextValue(c, "Second\n"), w, r)
	})
}

func MyThirdMiddleware(h middlewarec.Handler) middlewarec.Handler {
	return middlewarec.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTPC(appendToContextValue(c, "Third\n"), w, r)
	})
}

func MyFourthMiddleware(h middlewarec.Handler) middlewarec.Handler {
	return middlewarec.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTPC(appendToContextValue(c, "Fourth\n"), w, r)
	})
}

func main() {
	// Create a couple of middleware stacks using the Compose() function
	stackOne := middlewarec.Compose(MyFirstMiddleware, MySecondMiddleware)
	stackTwo := middlewarec.Compose(MyThirdMiddleware, MyFourthMiddleware)

	// Middleware compositions can also be composed with themselves
	stackThree := middlewarec.Compose(stackOne, stackTwo)

	// Alternatively, we can use method chaining, if we prefer
	stackFour := middlewarec.Middleware(MyFirstMiddleware).
		Compose(MySecondMiddleware).
		Compose(MyThirdMiddleware).
		Compose(stackThree)

	// Our actual handler. The custom middlewarec.Handler type accepts a context.Context
	// in addition to the standard http.ResponseWriter and *http.Request params from the
	// regular http.Handler type
	handler := middlewarec.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		valueFromContext := c.Value(0).(string)
		fmt.Fprintf(w, "Handler - %s", valueFromContext)
	})

	// Our context aware handlers can actually be used as a regular http.Handler. 
	// In this case context.Background() will be used as the root context for the request
	http.Handle("/one", stackOne(handler))

	// If we want more control over the context that is sent through our middleware stack, 
	// then we can use AsHttpHandlerWithContext()
	myCustomContext := appendToContextValue(context.Background(), "Custom Value\n")
	http.Handle("/two", stackTwo(handler).AsHttpHandlerWithContext(myCustomContext))

	http.Handle("/three", stackThree(handler))
	http.Handle("/four", stackFour(handler))

	http.ListenAndServe(":8080", nil)
}
```
