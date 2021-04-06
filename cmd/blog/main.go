package main

import (
	"fmt"
	"github.com/scottcagno/go-blog/internal"
	"github.com/scottcagno/go-blog/internal/user"
	"github.com/scottcagno/go-blog/pkg/logging"
	m "github.com/scottcagno/go-blog/pkg/middleware"
	"github.com/scottcagno/go-blog/pkg/web/templates"
	"github.com/scottcagno/go-blog/tools"
	"net/http"
	"os"
)

const (
	STATIC_PATH    = "web/static/"
	LISTENING_PORT = ":9090"
)

func init() {
	tools.CreateDirIfNotExist(STATIC_PATH)
}

func main() {

	logTo := os.Stdout
	logOut := logging.NewStdOutLogger(logTo)
	logErr := logging.NewStdErrLogger(logTo)
	tmpl := templates.NewTemplateCache("web/templates/*.html", logErr)

	mux := http.NewServeMux()

	mux.Handle("/", m.Get(internal.IndexHandler(tmpl)))
	mux.Handle("/favicon.ico", http.NotFoundHandler()) //m.GetOrPost(internal.FaviconHandler))

	u := user.NewUserService(tmpl)
	mux.Handle("/user", m.Get(u.UserHandler()))
	mux.Handle("/login", m.GetOrPost(u.LoginHandler()))
	mux.Handle("/logout", m.Get(u.LogoutHandler()))
	mux.Handle("/home", m.Get(u.HomeHandler()))

	chained1 := m.ChainedMiddleware(http.HandlerFunc(HandlerThree), HandlerTwo, HandlerOne, HandlerZero)

	mux.Handle("/chained1", chained1)
	mux.Handle("/chained2", HandlerZero(HandlerOne(HandlerTwo(http.HandlerFunc(HandlerThree)))))

	mux.Handle("/error/", http.NotFoundHandler()) //m.GetOrPost(internal.ErrorHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))))

	tools.HandleSignalInterrupt()
	logOut.Printf("Server started, listening on %s\n", LISTENING_PORT)
	loggingHandler := m.RequestLogger(logOut, logErr)
	err := http.ListenAndServe(LISTENING_PORT, loggingHandler(mux))
	logErr.Fatalf("Encountered error: %v\n", err)
}

func HandlerZero(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In handler ZERO")
		next.ServeHTTP(w, r)
	})
}

func HandlerOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In handler ONE")
		next.ServeHTTP(w, r)
	})
}

func HandlerTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In handler TWO")
		next.ServeHTTP(w, r)
	})
}

func HandlerThree(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In handler THREE")
	return
}
