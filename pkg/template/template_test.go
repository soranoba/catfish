package template

import (
	"github.com/soranoba/valis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplate_Render(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		ID    uint
		Name  string
		Param map[string]string
	}

	tpl := MustCompile(`{
  "id": {{ .ID }},
  "name": "{{ .Name }}",
  "param": "{{ .Param.a }}"
}`)

	val, err := tpl.Render(User{ID: 1, Name: "Alice", Param: map[string]string{"a": "A"}})
	assert.NoError(err)
	assert.Equal("{\n  \"id\": 1,\n  \"name\": \"Alice\",\n  \"param\": \"A\"\n}", val)
}

func TestTemplate_Render_nil(t *testing.T) {
	assert := assert.New(t)

	var tpl *Template
	val, err := tpl.Render(nil)
	assert.Equal("", val)
	assert.NoError(err)
}

func TestTemplate_UnmarshalText(t *testing.T) {
	assert := assert.New(t)
	tpl := Template{}
	assert.NoError(tpl.UnmarshalText([]byte("{{ifif}}{{end}}")))
	assert.Error(tpl.Validate())
	assert.NoError(tpl.UnmarshalText([]byte("{{if .Name}}{{.Name}}{{end}}")))
	assert.NoError(tpl.Validate())

	type User struct {
		Name string
	}

	val, err := tpl.Render(User{Name: "Alice"})
	assert.Equal("Alice", val)
	assert.NoError(err)
}

func TestExpr_MarshalText(t *testing.T) {
	assert := assert.New(t)
	tpl := MustCompile("{{if .Name}}{{.Name}}{{end}}")
	val, err := tpl.MarshalText()
	assert.NoError(err)
	assert.Equal("{{if .Name}}{{.Name}}{{end}}", string(val))
}

func TestExpr_Validate(t *testing.T) {
	var _ valis.Validatable = &Template{}

	assert := assert.New(t)

	var tpl Template
	assert.NoError(tpl.UnmarshalText([]byte("{{ifif}}{{end}}")))
	assert.EqualError(
		tpl.Validate(),
		"template: :1: function \"ifif\" not defined",
	)
}
