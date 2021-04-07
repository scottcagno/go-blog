package main

import (
	"fmt"
	"github.com/scottcagno/go-blog/internal/user"
	"github.com/scottcagno/go-blog/pkg/logging"
	m "github.com/scottcagno/go-blog/pkg/middleware"
	"github.com/scottcagno/go-blog/pkg/web"
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

	mux.Handle("/", http.NotFoundHandler())            //m.Get(internal.IndexHandler(tmpl)))
	mux.Handle("/favicon.ico", http.NotFoundHandler()) //m.GetOrPost(internal.FaviconHandler))

	u := user.NewUserService(tmpl)
	mux.Handle("/user", m.Get(u.UserHandler()))
	mux.Handle("/login", m.GetOrPost(u.LoginHandler()))
	mux.Handle("/logout", m.Get(u.LogoutHandler()))
	mux.Handle("/home", m.Get(u.HomeHandler()))

	mux.Handle("/chained2", HandlerZero(HandlerOne(HandlerTwo(http.HandlerFunc(HandlerThree)))))

	configuredRoute := HandlerZero(http.HandlerFunc(HandlerThree))
	mux.Handle("/chained1", configuredRoute)

	mux.Handle("/error/", http.NotFoundHandler()) //m.GetOrPost(internal.ErrorHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))))

	test := &Test{0}
	mux.Handle("/foo", test.Foo(http.HandlerFunc(HandlerThree)))

	tools.HandleSignalInterrupt()
	logOut.Printf("Server started, listening on %s\n", LISTENING_PORT)

	//loggingHandler := m.RequestLogger(logOut, logErr)
	//err := http.ListenAndServe(LISTENING_PORT, loggingHandler(mux))

	server := web.NewWebServer(os.Stdout, os.Stderr, nil)
	err := server.ListenAndServe(LISTENING_PORT, nil) //m.WithLogging(logOut, logErr)(mux))

	logErr.Fatalf("Encountered error: %v\n", err)
}

type Test struct {
	count int
}

func (t *Test) Foo(next http.Handler) http.Handler {
	foo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.count++
		fmt.Printf("In the Foo handler, counter is (%d)\n", t.count)
		_, _ = w.Write([]byte("handler Foo"))
		// call next
		next.ServeHTTP(w, r)
	})
	return t.Bar(foo)
}

func (t *Test) Bar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.count++
		fmt.Printf("In the Bar handler, counter is (%d)\n", t.count)
		_, _ = w.Write([]byte("handler Bar"))
		// call next
		next.ServeHTTP(w, r)
	})
}

func HandlerZero(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In handler ZERO")
		value := r.URL.Query().Get("pass")
		if value != "true" {
			http.Error(w, "bad access code", http.StatusUnauthorized)
			return
		}
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
	_, _ = w.Write([]byte("handler three"))
	return
}
