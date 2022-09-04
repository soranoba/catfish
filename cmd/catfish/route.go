package main

import (
	"net/http"
	"regexp"
	"strings"
)

type (
	Route struct {
		method     string
		path       string
		paramNames []string
		re         *regexp.Regexp
	}
)

func NewRoute(method string, path string) *Route {
	regTexts := make([]string, 0)
	paramNames := make([]string, 0)
	for _, segment := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/"), "/") {
		if strings.HasPrefix(segment, ":") {
			paramNames = append(paramNames, strings.Trim(segment, ":"))
			regTexts = append(regTexts, "([^/]*)")
		} else if strings.HasPrefix(segment, "*") {
			paramNames = append(paramNames, strings.Trim(segment, "*"))
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
		method:     strings.ToUpper(method),
		path:       path,
		paramNames: paramNames,
		re:         regexp.MustCompile("^" + strings.Join(regTexts, "/") + "$"),
	}
}

func (r *Route) IsMatch(req *http.Request, paramOut *map[string]string) bool {
	if r == nil {
		return false
	}

	if r.method != "*" && r.method != strings.ToUpper(req.Method) {
		return false
	}

	matched := r.re.FindStringSubmatch(strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/"))
	if len(matched) == 0 {
		return false
	}

	param := make(map[string]string)
	for i := 0; i < len(r.paramNames); i++ {
		param[r.paramNames[i]] = matched[i+1]
	}

	if paramOut != nil {
		*paramOut = param
	}
	return true
}
