package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish-server/pkg/config"
	"net/http"
	"sync"
)

type (
	HTTPHandler struct {
		config        config.Config
		mx            sync.Mutex
		routes        map[*Route][]*ResponsePreset
		defaultPreset *ResponsePreset
	}
)

func NewHTTPHandler(conf *config.Config) (*HTTPHandler, error) {
	errors := make([]error, 0)
	routes := make(map[*Route][]*ResponsePreset)

	for i, route := range conf.Routes {
		var presets []*ResponsePreset
		for name, res := range conf.Routes[i].Response {
			preset, err := NewResponsePreset(name, res)
			if err != nil {
				errors = append(errors, err)
			} else {
				presets = append(presets, preset)
			}
		}
		routes[NewRoute(route.Method, route.Path)] = presets
	}

	defaultPreset, err := NewResponsePreset("", conf.Default.Response)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		logrus.Error(errors)
	}

	return &HTTPHandler{
		config:        *conf,
		routes:        routes,
		defaultPreset: defaultPreset,
	}, nil
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handleRequest(w, req)
}

func (h *HTTPHandler) handleRequest(w http.ResponseWriter, req *http.Request) {
	h.mx.Lock()

	var ctx Context
	var preset *ResponsePreset
	var route *Route
	for r, presets := range h.routes {
		if r.IsMatch(req, &ctx) {
			route = r
			preset = electPreset(presets, h.defaultPreset)
			break
		}
	}

	h.mx.Unlock()

	if preset != nil {
		w.Header().Set("X-CATFISH-PATH", route.path)
	} else {
		w.Header().Set("X-CATFISH-PATH", "")
		preset = h.defaultPreset
	}

	for k, v := range preset.Header {
		w.Header().Set(k, v)
	}

	buf := new(bytes.Buffer)
	if err := preset.BodyTemplate.Execute(buf, ctx); err != nil {
		logrus.Warnf("Template rendering failed: %v", err)
		w.Header().Set("X-CATFISH-ERROR", err.Error())
	}

	w.WriteHeader(preset.Status)
	w.Write(buf.Bytes())
}
