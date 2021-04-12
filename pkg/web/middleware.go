package web

import (
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	http.Redirect(w, r, "/render?view="+name, http.StatusTemporaryRedirect)
}
