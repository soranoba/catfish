package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRoute_IsMatch(t *testing.T) {
	assert := assert.New(t)

	NewRequest := func(method string, path string) *http.Request {
		req, err := http.NewRequest(method, path, nil)
		assert.NoError(err, "[%s] %s", method, path)
		return req
	}

	check := func(route *Route, req *http.Request, pathParameters map[string]string) {
		if req == nil {
			return
		}
		var param map[string]string
		assert.Equal(pathParameters != nil, route.IsMatch(req, &param), "[%s] %s", req.Method, req.URL)
		assert.Equal(pathParameters, param, "[%s] %s", req.Method, req.URL)
	}

	// HTTP Method is case insensitive
	check(
		NewRoute("Get", "/users"),
		NewRequest("GET", "https://example.com/users"),
		map[string]string{},
	)
	check(
		NewRoute("GET", "/users"),
		NewRequest("Get", "https://example.com/users"),
		map[string]string{},
	)

	// Special HTTP Method
	check(
		NewRoute("*", "/users"),
		NewRequest("GET", "https://example.com/users"),
		map[string]string{},
	)

	// Different HTTP Methods
	check(
		NewRoute("GET", "/users"),
		NewRequest("POST", "https://example.com/users"),
		nil,
	)

	// Path parameters (Single segment)
	check(
		NewRoute("GET", "/users/:id"),
		NewRequest("GET", "https://example.com/users/1"),
		map[string]string{"id": "1"},
	)
	check(
		NewRoute("GET", "/users/:id"),
		NewRequest("GET", "https://example.com/users/1/foolowers"),
		nil,
	)
	check(
		NewRoute("GET", "/users/:id"),
		NewRequest("GET", "https://example.com/users/"),
		nil,
	)
	check(
		NewRoute("GET", "/users/:id"),
		NewRequest("GET", "https://example.com/users//"),
		map[string]string{"id": ""},
	)
	check(
		NewRoute("GET", "/users/:id/followers/:id"),
		NewRequest("GET", "https://example.com/users/1/followers/2"),
		map[string]string{"id": "2"},
	)
	check(
		NewRoute("GET", "/users/:id/followers/:fid"),
		NewRequest("GET", "https://example.com/users/1/followers/2"),
		map[string]string{"id": "1", "fid": "2"},
	)

	// Path parameters (Multiple segment)
	check(
		NewRoute("GET", "/*"),
		NewRequest("GET", "https://example.com/users/1"),
		map[string]string{"": "users/1"},
	)
	check(
		NewRoute("GET", "/*rest"),
		NewRequest("GET", "https://example.com/users/1/"),
		map[string]string{"rest": "users/1"},
	)
	check(
		NewRoute("GET", "/*a/x"),
		NewRequest("GET", "https://example.com/a/b"),
		nil,
	)
	check(
		NewRoute("GET", "/*a/x/*b"),
		NewRequest("GET", "https://example.com/a/b/x"),
		map[string]string{"a": "a/b", "b": ""},
	)
	check(
		NewRoute("GET", "/*a/x/*b"),
		NewRequest("GET", "https://example.com/a/x/b/x/c/x"),
		map[string]string{"a": "a/x/b/x/c", "b": ""},
	)
	check(
		NewRoute("GET", "/*a/x/*b"),
		NewRequest("GET", "https://example.com/a/b/x/c/d"),
		map[string]string{"a": "a/b", "b": "c/d"},
	)

	// Trailing slash
	check(
		NewRoute("GET", "/users/:id/"),
		NewRequest("GET", "https://example.com/users/1"),
		map[string]string{"id": "1"},
	)
	check(
		NewRoute("GET", "/users/:id/"),
		NewRequest("GET", "https://example.com/users/1/"),
		map[string]string{"id": "1"},
	)
	check(
		NewRoute("GET", "/users/:id"),
		NewRequest("GET", "https://example.com/users/1"),
		map[string]string{"id": "1"},
	)
	check(
		NewRoute("GET", "/users/:id"),
		NewRequest("GET", "https://example.com/users/1/"),
		map[string]string{"id": "1"},
	)
}
