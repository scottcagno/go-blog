package main

import (
	h "github.com/scottcagno/go-blog/internal/handlers"
	"github.com/scottcagno/go-blog/internal/logging"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	STATIC_PATH = "internal/templates/static"
)

var (
	templates *template.Template
	logger    *logging.Logger
)

func init() {
	funcMap := template.FuncMap{}
	templates = template.Must(template.New("*").Funcs(funcMap).ParseGlob("internal/templates/*.html"))
	func() {
		if _, err := os.Stat(STATIC_PATH); os.IsNotExist(err) {
			if err := os.MkdirAll(STATIC_PATH, 0655); err != nil {
				log.Fatalf("could not create static file path %q: %v\n", STATIC_PATH, err)
			}
		}
	}()
	logger = logging.NewLogger("blog")
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("static", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))))
	mux.Handle("/favicon.ico", h.CheckGetOrPost(logger, h.FaviconHandler))
	mux.Handle("/", h.CheckGet(logger, h.IndexHandler))
	mux.Handle("/login", h.CheckGetOrPost(logger, h.LoginHandler))
	mux.Handle("/logout", h.CheckGet(logger, h.LogoutHandler))
	mux.Handle("/home", h.CheckGet(logger, h.HomeHandler))

	mux.Handle("/home", h.CheckGetOrPost(logger, h.ErrorHandler))

	logger.Error.Fatalln(http.ListenAndServe("9090", mux))

}
