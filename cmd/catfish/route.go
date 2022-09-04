package main

import (
	"net/http"
	"regexp"
	"strings"
)

type (
	Route struct {
		method         string
		path           string
		pathParameters []string
		re             *regexp.Regexp
	}

	Context struct {
		Param map[string]string
	}
)

func NewRoute(method string, path string) *Route {
	regTexts := make([]string, 0)
	pathParameters := make([]string, 0)
	for _, segment := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/"), "/") {
		if strings.HasPrefix(segment, ":") {
			pathParameters = append(pathParameters, strings.Trim(segment, ":"))
			regTexts = append(regTexts, "([^/]*)")
		} else if strings.HasPrefix(segment, "*") {
			pathParameters = append(pathParameters, strings.Trim(segment, "*"))
			if len(regTexts) == 0 {
				regTexts = append(regTexts, "(.*)")
			} else {
				regTexts[len(regTexts)-1] += "/?(.*)"
			}
		} else {
			regTexts = append(regTexts, regexp.QuoteMeta(segment))
		}
	}

	return &Route{
		method:         strings.ToUpper(method),
		path:           path,
		pathParameters: pathParameters,
		re:             regexp.MustCompile("^" + strings.Join(regTexts, "/") + "$"),
	}
}

func (r *Route) IsMatch(req *http.Request, ctxOut *Context) bool {
	if r == nil {
		return false
	}

	if r.method != strings.ToUpper(req.Method) {
		return false
	}

	matched := r.re.FindStringSubmatch(strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/"))
	if len(matched) == 0 {
		return false
	}

	param := make(map[string]string)
	for i := 0; i < len(r.pathParameters); i++ {
		param[r.pathParameters[i]] = matched[i+1]
	}

	if ctxOut != nil {
		*ctxOut = Context{Param: param}
	}
	return true
}
