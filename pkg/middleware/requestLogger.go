package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: 200, wroteHeader: false}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
	return
}

func RequestLogger(stdout, stderr *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					stderr.Printf("err: %v, trace: %v\n", err, debug.Stack())
				}
			}()
			start := time.Now()
			wrapd := wrapResponseWriter(w)
			next.ServeHTTP(wrapd, r)
			// http status code general information:
			// 100-199='Info', 200-299='Success', 300-399='Redirects', 400-499='Client Error', 500-599='Server Error'
			if 400 <= wrapd.status && wrapd.status <= 599 {
				stderr.Printf("%s %s %d %s (%v)\n", r.Method, r.URL.EscapedPath(), wrapd.status, r.RemoteAddr, time.Since(start))
				return
			}
			stdout.Printf("%s %s %d %s (%v)\n", r.Method, r.URL.EscapedPath(), wrapd.status, r.RemoteAddr, time.Since(start))
			return
		}
		return http.HandlerFunc(fn)
	}
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		h.ServeHTTP(recorder, r)
		log.Printf("Handling request for %s from %s, status: %d", r.URL.Path, r.RemoteAddr, recorder.Status)
	})
}
