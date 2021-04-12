package main

import (
	"fmt"
	"github.com/scottcagno/go-blog/pkg/web"
	"log"
	"net/http"
	"strings"
)

func main() {
	mux := web.NewServeMux()
	mux.WithTemplates("web/templates/*.html")
	mux.Get("/one", getOne())
	mux.Get("/one/", getAllOne())
	mux.Get("/one/two", getOneTwo())
	mux.Get("/two", getTwo())
	mux.Get("/two/", getAllTwo())
	mux.Post("/foo", http.NotFoundHandler())
	mux.Get("/test", testController())
	mux.Get("/index/controller", indexController())
	mux.Get("/testing", getTesting())

	mux.Get("/endpoints", getEndpoints(mux.GetEntries()))

	mux.Static("/static/", "web/static/")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
	}
}

func getTesting() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("view")
		fmt.Fprintf(w, "path: %q, name: %q\n", r.URL.Path, name)
	}
	return http.HandlerFunc(fn)
}

func indexController() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//data := r.URL.Query().Get("data")
		data := "This is my index controller data, for now"
		web.Render(w, r, "index.html", data)
	}
	return http.HandlerFunc(fn)
}

func testController() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := "testController()"
		fmt.Fprintf(w, "%s hit!\n%s %s\n", s, r.Method, r.RequestURI)
	}
	return http.HandlerFunc(fn)
}

func getOne() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := "getOne()"
		fmt.Fprintf(w, "%s hit!\n%s %s\n", s, r.Method, r.RequestURI)
	}
	return http.HandlerFunc(fn)
}

func getAllOne() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := "getAllOne()"
		fmt.Fprintf(w, "%s hit!\n%s %s\n", s, r.Method, r.RequestURI)
	}
	return http.HandlerFunc(fn)
}

func getOneTwo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := "getOneTwo()"
		fmt.Fprintf(w, "%s hit!\n%s %s\n", s, r.Method, r.RequestURI)
	}
	return http.HandlerFunc(fn)
}

func getTwo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := "getTwo()"
		fmt.Fprintf(w, "%s hit!\n%s %s\n", s, r.Method, r.RequestURI)
	}
	return http.HandlerFunc(fn)
}

func getAllTwo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := "getAllTwo()"
		fmt.Fprintf(w, "%s hit!\n%s %s\n", s, r.Method, r.RequestURI)
	}
	return http.HandlerFunc(fn)
}

func getEndpoints(entries []string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		var sb strings.Builder
		sb.WriteString("<h1>Endpoints</h1><hr>")
		for i := 0; i < len(entries); i++ {
			sb.WriteString(fmt.Sprintf("%s\t", entries[i]))
			if strings.HasPrefix(entries[i], http.MethodGet) {
				uri := strings.Split(entries[i], " ")[1]
				sb.WriteString(fmt.Sprintf(` -> <a href="%s">%s</a>`, uri, uri))
			}
			sb.WriteString("<br>")
		}
		fmt.Fprintf(w, "%s", sb.String())
	}
	return http.HandlerFunc(fn)
}
