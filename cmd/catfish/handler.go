package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish/pkg/config"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type (
	HTTPHandler struct {
		config config.Config
		mx     sync.Mutex
		routes []*RouteData
	}
	RouteData struct {
		*Route
		parser  Parser
		presets []*ResponsePreset
	}
	Context struct {
		Method     string
		URL        *url.URL
		Param      map[string]string
		Body       map[string]interface{}
		ParseError error
	}
)

var (
	defaultPreset *ResponsePreset
)

func init() {
	defaultPreset, _ = NewResponsePreset("", &config.Response{
		Status: 404,
		Body:   "Not Found\n",
	})
}

func NewHTTPHandler(conf *config.Config) (*HTTPHandler, error) {
	errors := make([]error, 0)
	routes := make([]*RouteData, 0)

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
		routes = append(routes, &RouteData{
			Route:   NewRoute(route.Method, route.Path),
			parser:  NewParserWithName(route.ParserName),
			presets: presets,
		})
	}

	if len(errors) > 0 {
		logrus.Error(errors)
	}

	return &HTTPHandler{
		config: *conf,
		routes: routes,
	}, nil
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handleRequest(w, req)
}

func (h *HTTPHandler) handleRequest(w http.ResponseWriter, req *http.Request) {
	h.mx.Lock()

	var param map[string]string
	var preset *ResponsePreset
	var routePath string
	var parser Parser
	for _, route := range h.routes {
		if route.IsMatch(req, &param) {
			routePath = route.path
			parser = route.parser
			preset = ElectResponsePreset(route.presets, defaultPreset)
			break
		}
	}

	if preset == nil {
		preset = defaultPreset
	}

	h.mx.Unlock()

	time.Sleep(preset.Delay)

	w.Header().Set("X-CATFISH-PATH", routePath)
	w.Header().Set("X-CATFISH-RESPONSE-PRESET-NAME", preset.Name)

	var body map[string]interface{}
	var parseError error
	if parser != nil {
		parseError = parser.Parse(req.Body, &body)
	}

	buf := new(bytes.Buffer)
	ctx := Context{
		Method:     req.Method,
		URL:        req.URL,
		Param:      param,
		Body:       body,
		ParseError: parseError,
	}
	if err := preset.BodyTemplate.Execute(buf, &ctx); err != nil {
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
