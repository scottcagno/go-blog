package handlers

import (
	"github.com/scottcagno/go-blog/pkg/logging"
	"net/http"
)

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}

func WithLog(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Printf("(%s) %s %s\n", "WithLog", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// CheckGet is middleware to check for get method
func CheckGet(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			code := http.StatusMethodNotAllowed
			mesg := http.StatusText(code)
			logger.Error.Printf("%s %s -> %d %s\n", r.Method, r.URL.Path, code, mesg)
			return
		}
		logger.Info.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// CheckPost is middleware to check for get method
func CheckPost(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			code := http.StatusMethodNotAllowed
			mesg := http.StatusText(code)
			logger.Error.Printf("%s %s -> %d %s\n", r.Method, r.URL.Path, code, mesg)
			return
		}
		logger.Info.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// CheckGetOrPost is middleware to check for get OR post method
func CheckGetOrPost(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			code := http.StatusMethodNotAllowed
			mesg := http.StatusText(code)
			logger.Error.Printf("%s %s -> %d %s\n", r.Method, r.URL.Path, code, mesg)
			return
		}
		logger.Info.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
