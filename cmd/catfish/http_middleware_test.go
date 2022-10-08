package main

import (
	"net/http"
	"testing"
)

func TestHttpResponseWriter(t *testing.T) {
	var _ http.ResponseWriter = &httpResponseWriter{}
}
