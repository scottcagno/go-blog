package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var FaviconHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return
})

func IndexHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		err := t.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Printf("index handler, error: %v\n", err)
			http.Redirect(w, r, "/error/500", http.StatusInternalServerError)
		}
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
