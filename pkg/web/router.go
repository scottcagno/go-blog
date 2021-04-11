package web

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
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
	return len(s.entries)
}

func (s *ServeMux) Less(i, j int) bool {
	return s.entries[i].pattern < s.entries[j].pattern
}

func (s *ServeMux) Swap(i, j int) {
	s.entries[j], s.entries[i] = s.entries[i], s.entries[j]
}

type ServeMux struct {
	mu      sync.Mutex
	bp      sync.Pool
	stdout  *log.Logger
	stderr  *log.Logger
	entries []muxEntry
	t       *template.Template
}

func NewServeMux() *ServeMux {
	s := &ServeMux{
		bp: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		stdout:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		stderr:  log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		entries: make([]muxEntry, 0),
		t:       nil,
	}
	s.Get("/favicon.ico", http.NotFoundHandler())
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

func (s *ServeMux) newMuxEntry(method string, pattern string, handler http.Handler) {
	entry := muxEntry{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	s.entries = append(s.entries, entry)
	sort.Sort(s)
}

func (s *ServeMux) Forward(oldpattern string, newpattern string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newMuxEntry(http.MethodGet, oldpattern, http.RedirectHandler(newpattern, http.StatusTemporaryRedirect))
}

func (s *ServeMux) Handle(pattern string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newMuxEntry(http.MethodGet, pattern, handler)
	s.newMuxEntry(http.MethodPost, pattern, handler)
}

func (s *ServeMux) Get(pattern string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newMuxEntry(http.MethodGet, pattern, handler)
}

func (s *ServeMux) Post(pattern string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newMuxEntry(http.MethodPost, pattern, handler)
}

func (s *ServeMux) Put(pattern string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newMuxEntry(http.MethodPut, pattern, handler)
}

func (s *ServeMux) Delete(pattern string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newMuxEntry(http.MethodDelete, pattern, handler)
}

func (s *ServeMux) Static(pattern string, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	s.newMuxEntry(http.MethodGet, pattern, staticHandler)
}

func (s *ServeMux) GetEntries() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var entries []string
	for _, entry := range s.entries {
		entries = append(entries, fmt.Sprintln(entry))
	}
	return entries
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := http.NotFoundHandler()
	for _, entry := range s.entries {
		if entry.method != r.Method {
			continue
		}
		if entry.pattern != r.RequestURI {
			continue
		}
		h = entry.handler
		break
	}
	h.ServeHTTP(w, r)
}
