package main

import (
	h "github.com/scottcagno/go-blog/internal/handlers"
	"github.com/scottcagno/go-blog/internal/user"
	"github.com/scottcagno/go-blog/pkg/logging"
	m "github.com/scottcagno/go-blog/pkg/middleware"
	"github.com/scottcagno/go-blog/pkg/templates"
	"github.com/scottcagno/go-blog/tools"
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

	var ep = []string{
		"/static/",
		"/favicon.ico",
		"/",
		"/endpoints",
		"/user",
		"/login",
		"/logout",
		"/home",
		"/error/",
	}

	mux := http.NewServeMux()

	mux.Handle("/", m.CheckGet(logr, h.IndexHandler(tmpl)))
	mux.Handle("/favicon.ico", m.CheckGetOrPost(logr, h.FaviconHandler))
	mux.Handle("/endpoints", m.CheckGet(logr, h.EndpointHandler(ep, tmpl)))

	u := user.NewUserService(tmpl)
	mux.Handle("/user", m.CheckGet(logr, u.UserHandler()))
	mux.Handle("/login", m.CheckGetOrPost(logr, u.LoginHandler()))
	mux.Handle("/logout", m.CheckGet(logr, u.LogoutHandler()))
	mux.Handle("/home", m.CheckGet(logr, u.HomeHandler()))

	mux.Handle("/error/", m.CheckGetOrPost(logr, h.ErrorHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))))

	tools.HandleSignalInterrupt()
	log.Printf("Server started, listening on %s\n", LISTENING_PORT)
	logr.Error.Fatalf("Encountered error: %v\n", http.ListenAndServe(LISTENING_PORT, mux))
}
