package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish/pkg/config"
	"github.com/soranoba/catfish/pkg/evaler"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type (
	HTTPHandler struct {
		config            config.Config
		mx                sync.Mutex
		routes            []*RouteData
		totalRequestCount uint64
	}
	RouteData struct {
		*Route
		parser            Parser
		presets           []*ResponsePreset
		routeRequestCount uint64
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
	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.mx.Lock()

	var (
		param     map[string]string
		preset    *ResponsePreset
		routePath string
		parser    Parser
		err       error
	)
	h.totalRequestCount += 1

	for _, route := range h.routes {
		if route.IsMatch(req, &param) {
			routePath = route.path
			parser = route.parser
			route.routeRequestCount += 1
			preset, err = ElectResponsePreset(route.presets, evaler.Args{
				"routeRequestCount": route.routeRequestCount,
				"totalRequestCount": h.totalRequestCount,
			})
			break
		}
	}

	if preset == nil {
		preset = defaultPreset
	}

	h.mx.Unlock()

	w.Header().Set("X-CATFISH-PATH", routePath)
	w.Header().Set("X-CATFISH-RESPONSE-PRESET-NAME", preset.Name)

	if err != nil {
		logrus.Errorf("Conditional expression failed: %v", err)
		h.failedWithError(w, err)
		return
	}

	time.Sleep(preset.Delay)

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
		logrus.Errorf("Template rendering failed: %v", err)
		h.failedWithError(w, err)
		return
	}

	for k, v := range preset.Header {
		w.Header().Set(k, v)
	}

	w.WriteHeader(preset.Status)
	w.Write(buf.Bytes())
}

func (h *HTTPHandler) failedWithError(w http.ResponseWriter, err error) {
	w.Header().Set("X-CATFISH-ERROR", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}
