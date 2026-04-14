package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CustomDirective struct {
	Name        string   `yaml:"name"`
	Types       []string `yaml:"types"`
	ParamCount  int      `yaml:"param_count"`
	Description string   `yaml:"description"`
}

type Config struct {
	CustomDirectives []CustomDirective `yaml:"custom_directives"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
