package config

import (
	"github.com/soranoba/catfish/pkg/validator"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		Routes []Route `yaml:"routes"`
	}
	Route struct {
		Method   string              `yaml:"method" validate:"enums=GET POST PUT DELETE *"`
		Path     string              `yaml:"path" validate:"min=1"`
		Response map[string]Response `yaml:"response" required:"true" validate:"min=1"`
	}
	Response struct {
		Condition string            `yaml:"cond"`
		Delay     float64           `yaml:"delay"`
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

	v := validator.NewValidator()
	if err := v.Validate(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
