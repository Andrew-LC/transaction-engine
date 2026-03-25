package middleware

import (
	"log"
	"time"
	"net/http"
)


type LoggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size   int
}

func (lrw *LoggerResponseWriter) Write(val []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(val) 
	lrw.size += size
	return size, err
}
 
func (lrw *LoggerResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		loggerWrapper := &LoggerResponseWriter{
			ResponseWriter: w,
			statusCode:        200,
		}

		next.ServeHTTP(loggerWrapper, r)

		log.Printf(
			"[%s] %s %s | %d | %d bytes | %v\n",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			loggerWrapper.statusCode,
			loggerWrapper.size,
			time.Since(start),
		)
	})
}
