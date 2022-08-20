package main

import (
	"fmt"
	"github.com/soranoba/catfish-server/pkg/config"
	"html"
	"net/http"
	"sync"
)

type (
	HTTPHandler struct {
		config config.Config
		mx     sync.Mutex
		routes []*Route
	}
)

func NewHTTPHandler(conf *config.Config) *HTTPHandler {
	routes := make([]*Route, len(conf.Routes))
	for i, _ := range conf.Routes {
		routes[i] = NewRoute(conf.Routes[i].Method, conf.Routes[i].Path)
	}

	return &HTTPHandler{
		config: *conf,
		routes: routes,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

	h.mx.Lock()
	route, ctx, ok := h.getRoute(r)
	h.mx.Unlock()

	if !ok {
		return
	}

	fmt.Fprintf(w, "%s %#v", route.path, *ctx)
	_, _ = route, ctx
}

func (h *HTTPHandler) getRoute(r *http.Request) (*Route, *Context, bool) {
	var ctx Context
	for _, route := range h.routes {
		if route.IsMatch(r, &ctx) {
			return route, &ctx, true
		}
	}
	return nil, nil, false
}
