package main

import (
	"encoding/json"
	"io"
)

type (
	Parser interface {
		Parse(reader io.Reader, out interface{}) error
	}

	JsonParser struct{}
)

func NewParserWithName(name string) Parser {
	switch name {
	case "json":
		return &JsonParser{}
	default:
		return nil
	}
}

func (p *JsonParser) Parse(reader io.Reader, out interface{}) error {
	j := json.NewDecoder(reader)
	return j.Decode(out)
}
