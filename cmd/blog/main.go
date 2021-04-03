package main

import (
	h "github.com/scottcagno/go-blog/internal/handlers"
	"github.com/scottcagno/go-blog/pkg/logging"
	"github.com/scottcagno/go-blog/pkg/templates"
	"github.com/scottcagno/go-blog/tools"
	"github.com/scottcagno/net-tools/pkg/logger"
	"log"
	"net/http"
	"os"
)

const (
	STATIC_PATH    = "web/static/"
	LISTENING_PORT = ":9090"
)

var (
	tmpl *templates.TemplateCache
	logr *logging.Logger
)

func init() {

	func() {
		if _, err := os.Stat(STATIC_PATH); os.IsNotExist(err) {
			if err := os.MkdirAll(STATIC_PATH, 0655); err != nil {
				log.Fatalf("could not create static file path %q: %v\n", STATIC_PATH, err)
			}
		}
	}()
	logr = logging.NewLogger(os.Stdout)
	tmpl = templates.NewTemplateCache("web/templates/*.html", logr)
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("static", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))))
	mux.Handle("/favicon.ico", h.CheckGetOrPost(logr, h.FaviconHandler))
	mux.Handle("/", h.CheckGet(logr, h.IndexHandler(tmpl)))
	mux.Handle("/login", h.CheckGetOrPost(logr, h.LoginHandler))
	mux.Handle("/logout", h.CheckGet(logr, h.LogoutHandler))
	mux.Handle("/home", h.CheckGet(logr, h.HomeHandler))
	mux.Handle("/error/", h.CheckGetOrPost(logr, h.ErrorHandler))

	tools.HandleSignalInterrupt()
	log.Printf("Server started, listening on %s\n", LISTENING_PORT)
	logger.Error.Fatalf("Encountered error: %v\n", http.ListenAndServe(LISTENING_PORT, mux))
}
