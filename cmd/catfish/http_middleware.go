package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	httpResponseWriter struct {
		http.ResponseWriter
		statusCode int
	}
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapWriter := &httpResponseWriter{ResponseWriter: w, statusCode: 200}
		defer func() {
			logrus.WithFields(logrus.Fields{
				"@type":  "access",
				"status": wrapWriter.statusCode,
				"host":   r.Host,
				"method": r.Method,
				"uri":    r.RequestURI,
			}).Info()
		}()
		next.ServeHTTP(wrapWriter, r)
	})
}

func (w *httpResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
