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
				Method: "GET",
				Path:   "/users/:id",
				Response: map[string]Response{
					"200": {
						Condition: "rand(0.8)",
						Status:    200,
						Body:      "{\n  \"id\": 1,\n  \"name\": \"Alice\"\n}\n",
					},
					"500": {
						Condition: "true",
						Status:    500,
						Body:      "{\n  \"message\": \"Internal Server Error\"\n}\n",
					},
				},
			},
			{
				Method: "POST",
				Path:   "/users",
				Response: map[string]Response{
					"401": {
						Condition: "true",
						Status:    401,
						Body:      "{\n  \"message\": \"Unauthorized\"\n}\n",
					},
				},
			},
		},
	}, *conf)
}
