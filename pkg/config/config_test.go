package config

import (
	"github.com/stretchr/testify/assert"
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
	assert.Equal(Config{
		Routes: []Route{
			{
				Method:     "GET",
				Path:       "/users/:id",
				ParserName: "json",
				Response: []Response{
					{
						Name:      "200",
						Condition: "0.8",
						Delay:     0.1,
						Status:    200,
						Header: map[string]string{
							"Content-Type": "application/json",
						},
						Body: "{\n  \"id\": 1,\n  \"name\": \"Alice\"\n}\n",
					},
					{
						Name:      "500",
						Condition: "1.0",
						Status:    500,
						Header: map[string]string{
							"Content-Type": "application/json",
						},
						Body: "{\n  \"message\": \"Internal Server Error\"\n}\n",
					},
				},
			},
			{
				Method: "POST",
				Path:   "/users",
				Response: []Response{
					{
						Name:      "401",
						Condition: "true",
						Status:    401,
						Header: map[string]string{
							"Content-Type": "application/json",
						},
						Body: "{\n  \"message\": \"Unauthorized\"\n}\n",
					},
				},
			},
			{
				Method: "*",
				Path:   "*",
				Response: []Response{
					{
						Name:   "default",
						Status: 404,
						Header: map[string]string{
							"Content-Type": "application/json",
						},
						Body: "{\n  \"message\": \"Not Found\"\n}\n",
					},
				},
			},
		},
	}, *conf)
}
