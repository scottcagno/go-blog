package internal

import (
	"github.com/scottcagno/go-blog/pkg/web/templates"
	"net/http"
)

func IndexHandler(t *templates.TemplateCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		t.Render(w, r, "index.html", nil)
	})
}
