package template

import (
	"bytes"
	"io"
	"strings"
	"text/template"
)

type (
	Template struct {
		t   *template.Template
		raw string
	}
)

var (
	funcmap = template.FuncMap{
		"join":  strings.Join,
		"split": strings.Split,
		"add":   add,
		"sub":   sub,
	}
)

func MustCompile(text string) *Template {
	t, err := Compile(text)
	if err != nil {
		panic(err)
	}
	return t
}

func Compile(text string) (*Template, error) {
	t, err := compile(text)
	if err != nil {
		return nil, err
	}

	return &Template{
		t:   t,
		raw: text,
	}, nil
}

func (t *Template) Execute(w io.Writer, data interface{}) error {
	if t == nil {
		return nil
	}
	return t.t.Execute(w, data)
}

func (t *Template) Render(data interface{}) (string, error) {
	if t == nil {
		return "", nil
	}

	buf := new(bytes.Buffer)
	if err := t.t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t *Template) UnmarshalText(s []byte) error {
	raw := string(s)
	template, _ := compile(raw)
	*t = Template{
		t:   template,
		raw: raw,
	}
	return nil
}

func (t Template) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Template) Validate() error {
	if t == nil {
		return nil
	}

	if _, err := compile(t.raw); err != nil {
		return err
	}
	return nil
}

func (t *Template) String() string {
	return t.raw
}

func compile(text string) (*template.Template, error) {
	return template.New("").Funcs(funcmap).Parse(text)
}

func add(n1 int, n2 int) int {
	return n1 + n2
}

func sub(n1 int, n2 int) int {
	return n1 - n2
}
