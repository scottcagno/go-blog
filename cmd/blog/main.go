package main

import (
	h "github.com/scottcagno/go-blog/internal/handlers"
	"github.com/scottcagno/go-blog/pkg/logging"
	"github.com/scottcagno/go-blog/tools"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	STATIC_PATH    = "web/static/"
	LISTENING_PORT = ":9090"
)

var (
	templates *template.Template
	logger    *logging.Logger
)

func init() {
	funcMap := template.FuncMap{}

	base := template.Must(template.ParseFiles("web/templates/base.html"))
	templates = template.Must(base.Funcs(funcMap).ParseGlob("web/templates/*.html"))
	func() {
		if _, err := os.Stat(STATIC_PATH); os.IsNotExist(err) {
			if err := os.MkdirAll(STATIC_PATH, 0655); err != nil {
				log.Fatalf("could not create static file path %q: %v\n", STATIC_PATH, err)
			}
		}
	}()
	logger = logging.NewStdoutLogger()
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("static", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))))
	mux.Handle("/favicon.ico", h.CheckGetOrPost(logger, h.FaviconHandler))
	mux.Handle("/", h.CheckGet(logger, h.IndexHandler(templates)))
	mux.Handle("/login", h.CheckGetOrPost(logger, h.LoginHandler))
	mux.Handle("/logout", h.CheckGet(logger, h.LogoutHandler))
	mux.Handle("/home", h.CheckGet(logger, h.HomeHandler))
	mux.Handle("/error/", h.CheckGetOrPost(logger, h.ErrorHandler))

	tools.HandleSignalInterrupt()
	log.Printf("Server started, listening on %s\n", LISTENING_PORT)
	logger.Error.Fatalf("Encountered error: %v\n", http.ListenAndServe(LISTENING_PORT, mux))
}
