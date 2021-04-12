package web

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)
import "html/template"

type muxEntry struct {
	method  string
	pattern string
	handler http.Handler
}

func (m muxEntry) String() string {
	return fmt.Sprintf("%s %s", m.method, m.pattern)
}

func (s *ServeMux) Len() int {
	return len(s.es)
}

func (s *ServeMux) Less(i, j int) bool {
	return s.es[i].pattern < s.es[j].pattern
}

func (s *ServeMux) Swap(i, j int) {
	s.es[j], s.es[i] = s.es[i], s.es[j]
}

func (s *ServeMux) Search(x string) int {
	return sort.Search(len(s.es), func(i int) bool {
		return s.es[i].pattern >= x
	})
}

type ServeMux struct {
	mu     sync.Mutex
	bp     sync.Pool
	stdout *log.Logger
	stderr *log.Logger
	em     map[string]muxEntry
	es     []muxEntry
	t      *template.Template
}

func NewServeMux() *ServeMux {
	s := &ServeMux{
		bp: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		stdout: log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		stderr: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		em:     make(map[string]muxEntry),
		es:     make([]muxEntry, 0),
		t:      nil,
	}
	s.Get("/favicon.ico", http.NotFoundHandler())
	s.Get("/render", s.renderer())
	return s
}

func (s *ServeMux) WithLogging(stdout, stderr io.Writer) *ServeMux {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stdout = log.New(stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	s.stdout = log.New(stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	return s
}

func (s *ServeMux) WithTemplates(pattern string) *ServeMux {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.t = template.Must(template.New("*").Funcs(fm).ParseGlob(pattern))
	return s
}

func (s *ServeMux) Handle(method string, pattern string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := s.em[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}
	entry := muxEntry{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	s.em[pattern] = entry
	if pattern[len(pattern)-1] == '/' {
		s.es = appendSorted(s.es, entry)
	}
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].pattern) < len(e.pattern)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // Move shorter entries down
	es[i] = e
	return es
}

func (s *ServeMux) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	s.Handle(method, pattern, http.HandlerFunc(handler))
}

func (s *ServeMux) Forward(oldpattern string, newpattern string) {
	s.Handle(http.MethodGet, oldpattern, http.RedirectHandler(newpattern, http.StatusTemporaryRedirect))
}

func (s *ServeMux) Get(pattern string, handler http.Handler) {
	s.Handle(http.MethodGet, pattern, handler)
}

func (s *ServeMux) Post(pattern string, handler http.Handler) {
	s.Handle(http.MethodPost, pattern, handler)
}

func (s *ServeMux) Put(pattern string, handler http.Handler) {
	s.Handle(http.MethodPut, pattern, handler)
}

func (s *ServeMux) Delete(pattern string, handler http.Handler) {
	s.Handle(http.MethodDelete, pattern, handler)
}

func (s *ServeMux) Static(pattern string, path string) {
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	s.Handle(http.MethodGet, pattern, staticHandler)
}

func (s *ServeMux) GetEntries() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var entries []string
	for _, entry := range s.em {
		entries = append(entries, fmt.Sprintf("%s %s\n", entry.method, entry.pattern))
	}
	return entries
}

func (s *ServeMux) match(path string) (string, string, http.Handler) {
	e, ok := s.em[path]
	if ok {
		return e.method, e.pattern, e.handler
	}
	for _, e = range s.es {
		if strings.HasPrefix(path, e.pattern) {
			return e.method, e.pattern, e.handler
		}
	}
	return "", "", nil
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, _, h := s.match(r.URL.Path)
	if m != r.Method || h == nil {
		h = http.NotFoundHandler()
	}
	h = s.requestLogger(h)
	h.ServeHTTP(w, r)
}

func (s *ServeMux) renderer() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("view")
		if name == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}
		buffer := s.bp.Get().(*bytes.Buffer)
		buffer.Reset()
		data := struct {
			Data interface{}
		}{
			Data: r.URL.Query(),
		}
		err := s.t.ExecuteTemplate(buffer, name, data)
		if err != nil {
			s.bp.Put(buffer)
			s.stderr.Printf("Error while executing template (%s): %v\n", name, err)
			http.Redirect(w, r, "/error/404", http.StatusTemporaryRedirect)
			return
		}
		_, err = buffer.WriteTo(w)
		if err != nil {
			s.stderr.Printf("Error while writing (Render) to ResponseWriter: %v\n", err)
		}
		s.bp.Put(buffer)
		return
	}
	return http.HandlerFunc(fn)
}

func (s *ServeMux) ContentType(w http.ResponseWriter, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ct := mime.TypeByExtension(content)
	if ct == "" {
		s.stderr.Printf("Error, incompatible content type!\n")
		return
	}
	w.Header().Set("Content-Type", ct)
	return
}
