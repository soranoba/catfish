package template

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplate_Render(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		ID   uint
		Name string
	}

	tpl, err := New(`{
  "id": {{ .ID }},
  "name": "{{ .Name }}"
}`)
	assert.NoError(err)

	val, err := tpl.Render(User{ID: 1, Name: "Alice"})
	assert.NoError(err)
	assert.Equal("{\n  \"id\": 1,\n  \"name\": \"Alice\"\n}", val)
}
