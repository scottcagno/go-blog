package servemux

import (
	"log"
	"net/http"
	"sync"
)

// TODO: Choose name, I would like a clear API
// Server
// ServeMux
// WebServer
// HttpServer
// HTTPServer

type ServeMux struct {
	routes         []route
	stdout, stderr *log.Logger
	muxer          *http.ServeMux
	sync.RWMutex
}

type route struct {
	meth string
	patt string
	http.Handler
}

func NewServeMux(out, err *log.Logger) *ServeMux {
	return &ServeMux{
		routes: make([]route, 0),
		stdout: out,
		stderr: err,
		muxer:  http.NewServeMux(),
	}
}

func (s *ServeMux) Handle(pattern string, handler http.Handler) {

	s.muxer.Handle(pattern, handler)
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.muxer.ServeHTTP(w, r)
}
