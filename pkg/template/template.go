package template

import (
	"bytes"
	"io"
	"text/template"
)

type (
	Template struct {
		t *template.Template
	}
)

var (
	funcmap = template.FuncMap{}
)

func New(text string) (*Template, error) {
	t, err := template.New("").Funcs(funcmap).Parse(text)
	if err != nil {
		return nil, err
	}

	return &Template{
		t: t,
	}, nil
}

func (t *Template) Execute(w io.Writer, data interface{}) error {
	return t.t.Execute(w, data)
}

func (t *Template) Render(data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
