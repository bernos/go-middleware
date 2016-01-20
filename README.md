# go-middleware
Composable http middleware for golang
[![Build Status](https://travis-ci.org/bernos/go-middleware.svg)](https://travis-ci.org/bernos/go-middleware)

# Examples

```golang
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
