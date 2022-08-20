package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRoute_IsMatch(t *testing.T) {
	assert := assert.New(t)

	routes := []*Route{
		NewRoute("get", "/users/:id"),
		NewRoute("GET", "/users/:id/"),
		NewRoute("GET", "users/:id/"),
		NewRoute("Get", "users/:id"),
	}
	var ctx Context
	for _, route := range routes {
		req, err := http.NewRequest("GET", "https://example.com/users/1#aaaa?a=b", nil)
		if assert.NoError(err) {
			assert.True(route.IsMatch(req, &ctx))
			assert.Equal(map[string]string{"id": "1"}, ctx.Param)
		}
		req, err = http.NewRequest("GET", "https://example.com/users/aaaa/", nil)
		if assert.NoError(err) {
			assert.True(route.IsMatch(req, &ctx))
			assert.Equal(map[string]string{"id": "aaaa"}, ctx.Param)
		}
	}

	routes = []*Route{
		NewRoute("POST", "/users/:id"),
		NewRoute("GET", "/users/"),
		NewRoute("GET", "/users/2"),
	}
	for _, route := range routes {
		req, err := http.NewRequest("GET", "https://example.com/users/1#aaaa?a=b", nil)
		if assert.NoError(err) {
			assert.False(route.IsMatch(req, &ctx))
		}
	}
}
