package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish-server/pkg/config"
	"net/http"
	"sync"
	"time"
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
			preset, err := NewResponsePreset(name, &res)
			if err != nil {
				errors = append(errors, err)
			} else {
				presets = append(presets, preset)
			}
		}
		routes[NewRoute(route.Method, route.Path)] = presets
	}

	defaultPreset, err := NewResponsePreset("", &conf.Default.Response)
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
	var routePath string
	for route, presets := range h.routes {
		if route.IsMatch(req, &ctx) {
			routePath = route.path
			preset = ElectResponsePreset(presets, h.defaultPreset)
			break
		}
	}

	if preset == nil {
		preset = h.defaultPreset
	}

	h.mx.Unlock()

	time.Sleep(preset.Delay)

	w.Header().Set("X-CATFISH-PATH", routePath)
	w.Header().Set("X-CATFISH-RESPONSE-PRESET-NAME", preset.Name)

	buf := new(bytes.Buffer)
	if err := preset.BodyTemplate.Execute(buf, ctx); err != nil {
		logrus.Warnf("Template rendering failed: %v", err)
		w.Header().Set("X-CATFISH-ERROR", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, v := range preset.Header {
		w.Header().Set(k, v)
	}

	w.WriteHeader(preset.Status)
	w.Write(buf.Bytes())
}
