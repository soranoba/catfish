package main

import (
	"fmt"
	"html"
	"net/http"
)

type (
	HTTPHandler struct {
	}
)

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
