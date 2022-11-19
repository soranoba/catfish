package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadYamlFile(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	conf, err := LoadYamlFile(filepath.Join(dir, "testdata/config.yml"))
	assert.NoError(err)
	data, err := yaml.Marshal(conf)
	assert.NoError(err)
	assert.Equal(`routes:
    - method: GET
      path: /users/:id
      parser: json
      response:
        - name: "200"
          cond: "0.8"
          delay: 0.1
          redirect: null
          status: 200
          header:
            Content-Type: application/json
          body: |
            {
              "id": 1,
              "name": "Alice"
            }
        - name: "500"
          cond: "1.0"
          delay: 0
          redirect: null
          status: 500
          header:
            Content-Type: application/json
          body: |
            {
              "message": "Internal Server Error"
            }
    - method: POST
      path: /users
      parser: ""
      response:
        - name: "401"
          cond: "true"
          delay: 0
          redirect: null
          status: 401
          header:
            Content-Type: application/json
          body: |
            {
              "message": "Unauthorized"
            }
    - method: '*'
      path: /company
      parser: ""
      response:
        - name: company
          cond: null
          delay: 0
          redirect: https://soranoba.net
          status: 302
          header: {}
          body: null
    - method: '*'
      path: '*'
      parser: ""
      response:
        - name: default
          cond: null
          delay: 0
          redirect: null
          status: 404
          header:
            Content-Type: application/json
          body: |
            {
              "message": "Not Found"
            }
`, string(data))
}

func TestLoadYamlFile_invalidYaml(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	_, err = LoadYamlFile(filepath.Join(dir, "testdata/invalid_yaml.yml"))
	assert.EqualError(err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!seq into config.Config")
}

func TestLoadYamlFile_invalidType(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	_, err = LoadYamlFile(filepath.Join(dir, "testdata/invalid_type.yml"))
	assert.EqualError(err, "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `Interna...` into int")
}

func TestLoadYamlFile_invalidValidation(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	_, err = LoadYamlFile(filepath.Join(dir, "testdata/invalid_validation.yml"))
	assert.EqualError(err, `(inclusion) .Routes[0].Method is not included in [GET POST PUT DELETE *]
(lte) .Routes[0].Response[0].Status must be less than or equal to 599`)
}

func TestLoadYamlFile_invalidCond(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	_, err = LoadYamlFile(filepath.Join(dir, "testdata/invalid_cond.yml"))
	assert.EqualError(err, "(custom) .Routes[0].Response[0].Condition bad expression: 'x x'")
}

func TestLoadYamlFile_invalidBody(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	_, err = LoadYamlFile(filepath.Join(dir, "testdata/invalid_body.yml"))
	assert.EqualError(err, "(custom) .Routes[0].Response[0].Body template: :1: function \"ifif\" not defined")
}

func TestLoadYaml(t *testing.T) {
	assert := assert.New(t)

	conf, err := LoadYaml(strings.NewReader(`
routes:
  - method: GET
    path: /users/:id
    response:
      - name: "200"
        status: 200
        header:
          Content-Type: application/json
        body: |
          {
            "id": 1,
            "name": "Alice"
          }
`))
	assert.NoError(err)
	assert.Len(conf.Routes, 1)
}
