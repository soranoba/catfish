package main

import "net/url"

type (
	Context struct {
		Method     string
		URL        *url.URL
		Param      map[string]string
		Body       map[string]interface{}
		ParseError error
	}
)
