package handlers

import (
	"fmt"
	"github.com/scottcagno/go-blog/pkg/templates"
	"net/http"
)

var FaviconHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return
})

func IndexHandler(t *templates.TemplateCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		t.Render(w, r, "index.html", nil)
	})
}

var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from LOGIN")
	return
})

var LogoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from LOGOUT")
	return
})

var HomeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from HOME")
	return
})

var ErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from HOME")
	return
})
