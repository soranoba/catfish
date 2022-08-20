package main

import (
	"net/http"
	"strings"
)

type (
	Route struct {
		method         string
		path           string
		pathComponents []string
	}

	Context struct {
		Param map[string]string
	}
)

func NewRoute(method string, path string) *Route {
	path = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	return &Route{
		method:         strings.ToUpper(method),
		path:           path,
		pathComponents: strings.Split(path, "/"),
	}
}

func (r *Route) IsMatch(req *http.Request, ctxOut *Context) bool {
	if r == nil {
		return false
	}

	if r.method != req.Method {
		return false
	}

	path := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/")
	pc := strings.Split(path, "/")
	if len(pc) != len(r.pathComponents) {
		return false
	}

	param := make(map[string]string)
	for i := 0; i < len(pc); i++ {
		if strings.HasPrefix(r.pathComponents[i], ":") {
			param[strings.TrimPrefix(r.pathComponents[i], ":")] = pc[i]
		} else if r.pathComponents[i] != pc[i] {
			return false
		}
	}

	if ctxOut != nil {
		*ctxOut = Context{Param: param}
	}
	return true
}
