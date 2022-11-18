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

func TestTemplate_Render_join(t *testing.T) {
	assert := assert.New(t)

	tpl := MustCompile("{{join . \",\"}}")
	val, err := tpl.Render([]string{"A", "B", "C"})
	assert.NoError(err)
	assert.Equal("A,B,C", val)
}

func TestTemplate_Render_split(t *testing.T) {
	assert := assert.New(t)

	tpl := MustCompile("{{split . \",\"}}")
	val, err := tpl.Render("1,2,3,4,5")
	assert.NoError(err)
	assert.Equal("[1 2 3 4 5]", val)

	tpl = MustCompile(`
{{- $ids := split . "," -}}
[
{{- range $idx, $id := $ids}}
  {
    "id": {{$id}}
  }{{if lt $idx (sub (len $ids) 1)}},{{end -}}
{{end}}
]`)
	val, err = tpl.Render("1,2,3,4,5")
	assert.NoError(err)
	assert.Equal(`[
  {
    "id": 1
  },
  {
    "id": 2
  },
  {
    "id": 3
  },
  {
    "id": 4
  },
  {
    "id": 5
  }
]`, val)
}

func TestTemplate_Render_add(t *testing.T) {
	assert := assert.New(t)

	tpl := MustCompile("{{add . 3}}")
	val, err := tpl.Render(2)
	assert.NoError(err)
	assert.Equal("5", val)

	val, err = tpl.Render(2.5)
	assert.Error(err)

	tpl = MustCompile("{{add . 2.5}}")
	val, err = tpl.Render(2)
	assert.Error(err)
}

func TestTemplate_Render_sub(t *testing.T) {
	assert := assert.New(t)

	tpl := MustCompile("{{sub . 3}}")
	val, err := tpl.Render(2)
	assert.NoError(err)
	assert.Equal("-1", val)

	val, err = tpl.Render(2.5)
	assert.Error(err)

	tpl = MustCompile("{{sub . 2.5}}")
	val, err = tpl.Render(2)
	assert.Error(err)
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
