package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestJsonParser_Parse(t *testing.T) {
	assert := assert.New(t)

	reader := strings.NewReader("{\"id\":1,\"data\":{\"name\":\"Alice\"}}")
	var data map[string]interface{}

	assert.NoError(NewParserWithName("json").Parse(reader, &data))
	assert.Equal(map[string]interface{}{
		"id": float64(1),
		"data": map[string]interface{}{
			"name": "Alice",
		},
	}, data)
}
