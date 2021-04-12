package web

import (
	"log"
	"net/http"
	"sort"
	"sync"
)

type Router2 struct {
	stdout *log.Logger
	stderr *log.Logger
	sm     *http.ServeMux
}

func NewRouter2(stdout, stderr *log.Logger) *Router2 {
	return &Router2{
		stdout: stdout,
		stderr: stderr,
		sm:     new(http.ServeMux),
	}
}

func (rt *Router2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do stuff
	rt.sm.ServeHTTP(w, r)
}

func (rt *Router2) Forward(pattern string, url string) {
	rt.Handle(pattern, http.RedirectHandler(url, http.StatusTemporaryRedirect))
}

func (rt *Router2) Handle(pattern string, handler http.Handler) {
	rt.sm.Handle(pattern, handler)
}

func (rt *Router2) HandleStatic(pattern string, path string) {
	rt.sm.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(path))))
}

type route struct {
	pattern string
	handler http.Handler
}

type Router struct {
	stdout  *log.Logger
	stderr  *log.Logger
	entries map[string]route
	routes  []route
	mu      sync.Mutex
}

func NewRouter(stdout, stderr *log.Logger) *Router {
	return &Router{
		stdout:  stdout,
		stderr:  stderr,
		entries: make(map[string]route),
		routes:  make([]route, 0),
	}
}

func (rt *Router) Handle(pattern string, handler http.Handler) {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := rt.entries[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if rt.entries == nil {
		rt.entries = make(map[string]route)
	}
	rte := route{
		pattern: pattern,
		handler: handler,
	}
	rt.entries[pattern] = rte
	if pattern[len(pattern)-1] == '/' {
		rt.routes = appendSorted(rt.routes, rte)
	}
}

func appendSorted(routes []route, rte route) []route {
	n := len(routes)
	i := sort.Search(n, func(i int) bool {
		return len(routes[i].pattern) < len(rte.pattern)
	})
	if i == n {
		return append(routes, rte)
	}
	// we now know that i points at where we want to insert
	routes = append(routes, route{}) // try to grow the slice in place, any entry works.
	copy(routes[i+1:], routes[i:])   // Move shorter es down
	routes[i] = rte
	return routes
}

// HandleFunc registers the handler function for the given pattern.
func (rt *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	rt.Handle(pattern, http.HandlerFunc(handler))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do stuff
	rt.ServeHTTP(w, r)
}
