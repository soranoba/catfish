package main

import (
	"net/http"
	"net/url"
)

type (
	Context struct {
		Method     string
		URL        *url.URL
		Header     http.Header
		Param      map[string]string
		Body       map[string]interface{}
		ParseError error
	}
)
