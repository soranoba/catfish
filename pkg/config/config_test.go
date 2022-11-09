package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYamlFile(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	conf, err := LoadYamlFile(filepath.Join(dir, "testdata/config1.yml"))
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
          body: ""
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

func TestLoadYamlFile_invalidCond(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.Getwd()
	assert.NoError(err)

	_, err = LoadYamlFile(filepath.Join(dir, "testdata/config2.yml"))
	assert.EqualError(err, "parsing error: x x\t:1:3 - 1:4 unexpected Ident while scanning operator")
}
