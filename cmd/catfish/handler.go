package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish/pkg/config"
	"github.com/soranoba/catfish/pkg/evaler"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type (
	AdminHTTPHandler struct {
		*HTTPHandler
		fileServer http.Handler
	}
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

	Variables struct {
		GlobalVariables GlobalVariables  `json:"global_variables"`
		RouteVariables  []RouteVariables `json:"route_variables"`
	}
	GlobalVariables struct {
		TotalRequestCount uint64 `json:"totalRequestCount"`
	}
	RouteVariables struct {
		Route     RouteVariablesKey   `json:"route"`
		Variables RouteVariablesValue `json:"variables"`
	}
	RouteVariablesKey struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	}
	RouteVariablesValue struct {
		RouteRequestCount uint64 `json:"routeRequestCount"`
	}
)

var (
	defaultPreset *ResponsePreset
	//go:embed static/public
	staticFs embed.FS
)

func init() {
	defaultPreset, _ = NewResponsePreset(&config.Response{
		Status: 404,
		Body:   "Not Found\n",
	})
}

func NewHTTPHandler(conf *config.Config) (*HTTPHandler, error) {
	errors := make([]error, 0)
	routes := make([]*RouteData, 0)

	for i, route := range conf.Routes {
		var presets []*ResponsePreset
		for _, res := range conf.Routes[i].Response {
			preset, err := NewResponsePreset(&res)
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

func (h *HTTPHandler) handleAdminRequest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		h.setCORSResponse(w)
		return
	}

	h.mx.Lock()
	defer h.mx.Unlock()

	switch req.Method {
	case http.MethodGet:
		switch req.URL.Path {
		case "/api/config":
			b, err := json.Marshal(h.config)
			if err != nil {
				h.failedWithError(w, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		case "/api/variables":
			routeVariables := make([]RouteVariables, len(h.routes))
			for i, route := range h.routes {
				routeVariables[i] = RouteVariables{
					Route: RouteVariablesKey{
						Method: route.method,
						Path:   route.path,
					},
					Variables: RouteVariablesValue{
						RouteRequestCount: route.routeRequestCount,
					},
				}
			}
			variables := &Variables{
				GlobalVariables: GlobalVariables{
					TotalRequestCount: h.totalRequestCount,
				},
				RouteVariables: routeVariables,
			}
			b, err := json.Marshal(variables)
			if err != nil {
				h.failedWithError(w, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	case http.MethodPut:
		switch req.URL.Path {
		case "/api/variables/reset":
			h.totalRequestCount = 0
			for _, route := range h.routes {
				route.routeRequestCount = 0
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *HTTPHandler) handleRequest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		h.setCORSResponse(w)
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
			preset, err = ElectResponsePreset(route.presets, evaler.Params{
				"routeRequestCount": route.routeRequestCount,
				"totalRequestCount": h.totalRequestCount,
				"param":             param,
				"query":             req.URL.Query(),
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

	if preset.Redirect != nil {
		http.Redirect(w, req, *preset.Redirect, preset.Status)
		return
	}

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

func (h *HTTPHandler) setCORSResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Origin,Accept,Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) failedWithError(w http.ResponseWriter, err error) {
	w.Header().Set("X-CATFISH-ERROR", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func NewAdminHTTPHandler(h *HTTPHandler) *AdminHTTPHandler {
	dir, err := fs.Sub(staticFs, "static/public")
	if err != nil {
		log.Fatal(err)
	}

	return &AdminHTTPHandler{
		HTTPHandler: h,
		fileServer:  http.FileServer(http.FS(dir)),
	}
}

func (h *AdminHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api/") {
		h.HTTPHandler.handleAdminRequest(w, req)
		return
	}
	h.fileServer.ServeHTTP(w, req)
}
