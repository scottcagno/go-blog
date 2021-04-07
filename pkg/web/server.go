package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

type webServerHandler struct {
	ws *WebServer
}

type WebServer struct {
	server *http.Server
	stdout *log.Logger
	stderr *log.Logger
}

func NewWebServer(out, err io.Writer, server *http.Server) *WebServer {
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}
	if server == nil {
		server = &http.Server{}
	}
	s := &WebServer{
		server: server,
		stdout: log.New(out, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		stderr: log.New(err, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
	}
	if server.Handler == nil {
		server.Handler = webServerHandler{s}
	}
	return s
}

func (s *WebServer) LogAccess(format string, v ...interface{}) {
	s.stdout.Printf(format, v...)
}

func (s *WebServer) LogError(format string, v ...interface{}) {
	s.stderr.Printf(format, v...)
}

func (s webServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.ws.stderr.Printf("err: %v, trace: %v\n", err, debug.Stack())
		}
	}()
	lrw := loggingResponseWriter{
		ResponseWriter: w,
		data: &responseData{
			status: 200,
			size:   0,
		},
	}
	handler := s.ws.server.Handler
	handler.ServeHTTP(&lrw, r)
	format, values := ">> %s - - [%s] \"%s %s %s\" %d %d\n", []interface{}{
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.URL.EscapedPath(),
		r.Proto,
		lrw.data.status,
		r.ContentLength,
	}
	if 400 <= lrw.data.status && lrw.data.status <= 599 {
		s.ws.stderr.Printf(format, values...)
		return
	}
	s.ws.stdout.Printf(format, values...)
	return
}

func (s *WebServer) ServeHTTP1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("should be doing logging stuff here most likely")
	s.server.Handler.ServeHTTP(w, r)
}

func (s *WebServer) ListenAndServe(addr string, handler http.Handler) error {
	if s.server.Addr == "" {
		s.server.Addr = addr
	}
	if s.server.Handler == nil {
		s.server.Handler = handler
	}
	return s.server.ListenAndServe()
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.data.status = statusCode
}
