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
	return &responseWriter{ResponseWriter: w}
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
			if 100 <= wrapd.status && wrapd.status <= 399 {
				stdout.Printf("%s %s %d (%v)\n", r.Method, r.URL.EscapedPath(), wrapd.status, time.Since(start))
				return
			}
			if 400 <= wrapd.status && wrapd.status <= 599 {
				stdout.Printf("%s %s %d (%v)\n", r.Method, r.URL.EscapedPath(), wrapd.status, time.Since(start))
				return
			}
		}
		return http.HandlerFunc(fn)
	}
}
