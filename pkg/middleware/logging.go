package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

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

func WithLogging(stdout, stderr *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					stderr.Printf("err: %v, trace: %v\n", err, debug.Stack())
				}
			}()
			lrw := loggingResponseWriter{
				ResponseWriter: w,
				data: &responseData{
					status: 200,
					size:   0,
				},
			}
			next.ServeHTTP(&lrw, r)
			format, values := "# %s - - [%s] \"%s %s %s\" %d %d\n", []interface{}{
				r.RemoteAddr,
				time.Now().Format(time.RFC1123Z),
				r.Method,
				r.URL.EscapedPath(),
				r.Proto,
				lrw.data.status,
				r.ContentLength,
			}
			if 400 <= lrw.data.status && lrw.data.status <= 599 {
				stderr.Printf(format, values...)
				return
			}
			stdout.Printf(format, values...)
			return
		}
		return http.HandlerFunc(fn)
	}
}
