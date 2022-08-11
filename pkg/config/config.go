package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		Routes []Route `yaml:"routes"`
	}
	Route struct {
		Method   string              `yaml:"method" validate:"min=1"`
		Path     string              `yaml:"path" validate:"min=1"`
		Response map[string]Response `yaml:"response" required:"true" validate:"min=1"`
	}
	Response struct {
		Condition string            `yaml:"cond"`
		Status    int               `yaml:"status" validate:"min=100,max=599"`
		Header    map[string]string `yaml:"header"`
		Body      string            `yaml:"body"`
	}
)

func LoadYamlFile(filepath string) (*Config, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&conf); err != nil {
		return &conf, err
	}

	return &conf, nil
}
