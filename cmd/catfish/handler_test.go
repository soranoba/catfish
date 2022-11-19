package main

import (
	"github.com/soranoba/catfish/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	suite.Run(t, &HandlerSuite{
		assert: assert.New(t),
	})
}

type HandlerSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *HandlerSuite) TestSuccess() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - name: Response name
        status: 202
        header:
          Content-Type: application/json
        body: "{}"
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
	resp := suite.do(req)
	body, _ := io.ReadAll(resp.Body)

	suite.assert.Equal(202, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
	suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
	suite.assert.Equal("Response name", resp.Header.Get(HeaderCatfishResponsePresetName))
	suite.assert.Equal("application/json", resp.Header.Get("Content-Type"))
	suite.assert.Equal("{}", string(body))
}

func (suite *HandlerSuite) TestFailure() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - name: Response name
        status: 200
        header:
          Content-Type: application/json
        body: "{{.NoParam}}"
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
	resp := suite.do(req)
	body, _ := io.ReadAll(resp.Body)

	suite.assert.Equal(500, resp.StatusCode)
	suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
	suite.assert.Equal("Response name", resp.Header.Get(HeaderCatfishResponsePresetName))
	suite.assert.Equal("template: :1:2: executing \"\" at <.NoParam>: can't evaluate field NoParam in type *main.Context", resp.Header.Get(HeaderCatfishError))
	suite.assert.Equal("", string(body))
}

func (suite *HandlerSuite) TestNoRoute() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 202
        body: "{}"
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL+"/about", nil)
	resp := suite.do(req)
	body, _ := io.ReadAll(resp.Body)

	suite.assert.Equal(404, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
	suite.assert.Equal("", resp.Header.Get(HeaderCatfishPath))
	suite.assert.Equal("", resp.Header.Get(HeaderCatfishResponsePresetName))
	suite.assert.Equal("Not Found\n", string(body))
}

func (suite *HandlerSuite) TestDelay() {
	func() {
		handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 200
        delay: 0.05
        body: "{}"
`)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		before := time.Now()
		resp := suite.do(req)
		duration := time.Now().Sub(before)
		suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
		suite.assert.GreaterOrEqual(duration, 50*time.Millisecond)
	}()

	func() {
		handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 200
        body: "{}"
`)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		before := time.Now()
		resp := suite.do(req)
		duration := time.Now().Sub(before)
		suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
		suite.assert.LessOrEqual(duration, 50*time.Millisecond)
	}()
}

func (suite *HandlerSuite) TestRedirect() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - name: redirect
        status: 302
        redirect: http://example.com
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL+"/users?key=value", nil)
	resp := suite.doWithoutRedirect(req)
	body, _ := io.ReadAll(resp.Body)

	suite.assert.Equal(302, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
	suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
	suite.assert.Equal("redirect", resp.Header.Get(HeaderCatfishResponsePresetName))
	suite.assert.Equal("<a href=\"http://example.com?key=value\">Found</a>.\n\n", string(body))

	resp = suite.do(req)
	suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
	suite.assert.Equal("", resp.Header.Get(HeaderCatfishPath))
	suite.assert.Equal("", resp.Header.Get(HeaderCatfishResponsePresetName))

	// NOTE: keep queries.
	suite.assert.Equal("http://example.com?key=value", resp.Request.URL.String())
}

func (suite *HandlerSuite) TestDelayWithRedirect() {
	func() {
		handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 302
        delay: 0.05
        redirect: http://example.com
`)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		before := time.Now()
		_ = suite.doWithoutRedirect(req)
		duration := time.Now().Sub(before)
		suite.assert.GreaterOrEqual(duration, 50*time.Millisecond)
	}()

	func() {
		handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 302
        redirect: http://example.com
`)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		before := time.Now()
		_ = suite.doWithoutRedirect(req)
		duration := time.Now().Sub(before)
		suite.assert.LessOrEqual(duration, 50*time.Millisecond)
	}()
}

func (suite *HandlerSuite) TestCond_totalRequestCount() {
	handler := suite.createHandler(`
routes:
  - method: PUT
    path: /users/:id
    response:
      - status: 200
        body: >
          {"name": "Bob"}
  - method: GET
    path: /users
    response:
      - status: 200
        cond: totalRequestCount == 1
        body: >
          [{"name": "Alice"}]
      - status: 200
        cond: totalRequestCount == 3
        body: >
          [{"name": "Alice"}, {"name": "Bob"}]
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("[{\"name\": \"Alice\"}]\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("PUT", srv.URL+"/users/2", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users/:id", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("{\"name\": \"Bob\"}\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("[{\"name\": \"Alice\"}, {\"name\": \"Bob\"}]\n", string(body))
	}()
}

func (suite *HandlerSuite) TestCond_routeRequestCount() {
	handler := suite.createHandler(`
routes:
  - method: PUT
    path: /users/:id
    response:
      - status: 200
        cond: routeRequestCount == 1
        body: >
          {"name": "Bob"}
  - method: GET
    path: /users
    response:
      - status: 200
        cond: routeRequestCount == 1
        body: >
          [{"name": "Alice"}]
      - status: 200
        cond: routeRequestCount == 2
        body: >
          [{"name": "Alice"}, {"name": "Bob"}]
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("[{\"name\": \"Alice\"}]\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("PUT", srv.URL+"/users/2", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users/:id", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("{\"name\": \"Bob\"}\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("[{\"name\": \"Alice\"}, {\"name\": \"Bob\"}]\n", string(body))
	}()
}

func (suite *HandlerSuite) TestCond_param() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users/:id
    response:
      - status: 200
        cond: param["id"] == 1
        body: >
          {"name": "Alice"}
      - status: 200
        cond: param["id"] == 2
        body: >
          {"name": "Bob"}
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users/1", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("/users/:id", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("{\"name\": \"Alice\"}\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users/2", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
		suite.assert.Equal("/users/:id", resp.Header.Get(HeaderCatfishPath))
		suite.assert.Equal("{\"name\": \"Bob\"}\n", string(body))
	}()
}

func (suite *HandlerSuite) TestCond_query() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 200
        cond: query["id"][0] == 1
        body: >
          {"name": "Alice"}
      - status: 200
        cond: query["id"][0] == 2
        body: >
          {"name": "Bob"}
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users?id=1", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode)
		suite.assert.Equal("{\"name\": \"Alice\"}\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users?id=2", nil)
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
		suite.assert.Equal("{\"name\": \"Bob\"}\n", string(body))
	}()
}

func (suite *HandlerSuite) TestCond_header() {
	// NOTE: Header names are always converted to CamelCase
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users
    response:
      - status: 200
        cond: header["Id"][0] == 1
        body: >
          {"name": "Alice"}
      - status: 200
        cond: header["Id"][0] == 2
        body: >
          {"name": "Bob"}
`)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		req.Header.Set("ID", "1")
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
		suite.assert.Equal("{\"name\": \"Alice\"}\n", string(body))
	}()
	func() {
		req, _ := http.NewRequest("GET", srv.URL+"/users", nil)
		req.Header.Set("ID", "2")
		resp := suite.do(req)
		body, _ := io.ReadAll(resp.Body)

		suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
		suite.assert.Equal("{\"name\": \"Bob\"}\n", string(body))
	}()
}

func (suite *HandlerSuite) TestBody() {
	handler := suite.createHandler(`
routes:
  - method: GET
    path: /users/:id/*path
    response:
      - status: 200
        body: >

          Method: {{.Method}}{{print "\n" -}}
          URL: {{.URL.String}}{{print "\n" -}}
          Content-Type: {{.Header.Get "Content-Type"}}{{print "\n" -}}
          Param(id): {{.Param.id}}{{print "\n" -}}
          Param(path): {{.Param.path}}

`)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL+"/users/1/2/3", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := suite.do(req)
	body, _ := io.ReadAll(resp.Body)

	suite.assert.Equal(200, resp.StatusCode, resp.Header.Get(HeaderCatfishError))
	suite.assert.Equal(`
Method: GET
URL: /users/1/2/3
Content-Type: application/json
Param(id): 1
Param(path): 2/3
`, string(body))
}

func (suite *HandlerSuite) createHandler(yaml string) *HTTPHandler {
	conf, err := config.LoadYaml(strings.NewReader(yaml))
	if !suite.assert.NoError(err) {
		suite.T().FailNow()
	}

	handler, err := NewHTTPHandler(conf)
	if !suite.assert.NoError(err) {
		suite.T().FailNow()
	}
	return handler
}

func (suite *HandlerSuite) do(req *http.Request) *http.Response {
	resp, err := (&http.Client{}).Do(req)
	if !suite.assert.NoError(err) {
		suite.T().FailNow()
	}
	return resp
}

func (suite *HandlerSuite) doWithoutRedirect(req *http.Request) *http.Response {
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if !suite.assert.NoError(err) {
		suite.T().FailNow()
	}
	return resp
}
