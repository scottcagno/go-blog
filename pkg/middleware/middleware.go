package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func ChainedMiddleware(wrapped http.Handler, m ...Middleware) http.Handler {
	// loop in reverse to preserve order
	for i := len(m) - 1; i > -1; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

func NewExampleMiddleware(someThing string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Logic here

			// Call the next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Middleware is an example middleware handler
func MiddlewareExample(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}

// ChainedMiddleware is an example middleware chaining handler
func ChainedMiddlewareWorse(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}

// Get is middleware to check for get method
func Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Post is middleware to check for get method
func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetOrPost is middleware to check for get OR post method
func GetOrPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
