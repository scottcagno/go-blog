package handlers

import (
	"fmt"
	"github.com/scottcagno/go-blog/pkg/templates"
	"net/http"
	"strings"
)

var FaviconHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return
})

func EndpointHandler(ep []string, t *templates.TemplateCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i, _ := range ep {
			ep[i] = `<a href="` + ep[i] + `">` + ep[i] + `</a>`
		}
		t.Raw(w, "%s", strings.Join(ep, "<br>"))
	})
}

func IndexHandler(t *templates.TemplateCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		t.Render(w, r, "index.html", nil)
	})
}

var ErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from HOME")
	return
})
