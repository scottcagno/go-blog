package main

import (
	"github.com/scottcagno/go-blog/pkg/logging"
	"github.com/scottcagno/go-blog/pkg/web"
	"github.com/scottcagno/go-blog/pkg/web/templates"
	"github.com/scottcagno/go-blog/tools"
	"net/http"
	"os"
	"strconv"
)

func main() {

	// set up loggers and template cache
	stdout, stderr := logging.NewLogger(os.Stdout, os.Stderr)
	t := templates.NewTemplateCache("web/templates/*.html", stderr)

	// set up routes
	mux := web.NewServeMux().WithLogging(stdout.Writer(), stderr.Writer())

	// handle not found
	mux.Get("/favicon.ico", http.NotFoundHandler())

	// forward, for testing purposes
	mux.Forward("/", "/user")

	// handle user model, auto generating html form
	mux.Get("/user", HandleIndex(t))

	// handle errors
	mux.Get("/error", HandleError())

	// handle static content
	mux.Static("/static/", "web/static/")

	tools.HandleSignalInterrupt()
	stderr.Fatalln(http.ListenAndServe(":8080", mux))

}

type User struct {
	ID       int    `html:"number"`
	Name     string `html:"text"`
	Email    string `html:"email"`
	Password string `html:"password"`
	IsActive bool   `html:"checkbox"`
}

func HandleError() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		val := r.URL.Query().Get("code")
		if val == "" {
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		code, err := strconv.Atoi(val)
		if err != nil {
			status := http.StatusExpectationFailed
			http.Error(w, http.StatusText(status), status)
			return
		}
		status := http.StatusText(code)
		if status == "" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Error(w, status, code)
		return
	}
	return http.HandlerFunc(fn)
}

func HandleIndex(t *templates.TemplateCache) http.Handler {
	user := User{
		ID:       34,
		Name:     "Jon Doe",
		Email:    "jdoe@ex.com",
		Password: "none",
		IsActive: true,
	}
	model, err := tools.MakeModel(&user, "html")
	if err != nil {
		code := http.StatusNotAcceptable
		return http.RedirectHandler("/error?code="+http.StatusText(code), code)
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"User":  user,
			"Model": model,
		}
		t.Render(w, r, "user-model.html", data)
	}
	return http.HandlerFunc(fn)
}
