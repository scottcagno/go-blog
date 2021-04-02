package handlers

import (
	"fmt"
	"net/http"
)

var FaviconHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return
})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from INDEX")
	return
})

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
