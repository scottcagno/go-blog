package web

import (
	"net/http"
)

type ctxEntry struct {
	name string
	data interface{}
}

func Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	http.Redirect(w, r, "/render?view="+name, http.StatusTemporaryRedirect)
}
