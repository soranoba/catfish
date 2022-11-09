package config

import (
	"github.com/soranoba/catfish/pkg/evaler"
	"github.com/soranoba/catfish/pkg/validator"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

type (
	Config struct {
		Routes []Route `yaml:"routes" json:"routes"`
	}
	Route struct {
		Method     string     `yaml:"method" json:"method" enums:"GET,POST,PUT,DELETE,*"`
		Path       string     `yaml:"path" json:"path" validate:"min=1"`
		ParserName string     `yaml:"parser" json:"parser" enums:"json,"`
		Response   []Response `yaml:"response" json:"response" required:"true" validate:"min=1"`
	}
	Response struct {
		Name      string            `yaml:"name" json:"name"`
		Condition *evaler.Expr      `yaml:"cond" json:"cond"`
		Delay     float64           `yaml:"delay" json:"delay"`
		Redirect  *URL              `yaml:"redirect" json:"redirect"`
		Status    int               `yaml:"status" json:"status" validate:"min=100,max=599"`
		Header    map[string]string `yaml:"header" json:"header"`
		Body      string            `yaml:"body" json:"body"`
	}
	URL struct {
		url.URL
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

func (u *URL) UnmarshalText(s []byte) error {
	if len(s) == 0 {
		return nil
	}

	url, err := url.Parse(string(s))
	if err != nil {
		return err
	}

	*u = URL{URL: *url}
	return nil
}

func (u URL) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}
