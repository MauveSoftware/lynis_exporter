package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v1"
)

type Config struct {
	ReportFilePath string                      `yaml:"reportfile_path"`
	Metrics        map[string]MetricDefinition `yaml:"metrics"`
}

type MetricDefinition struct {
	MetricName  string `yaml:"name"`
	Description string `yaml:"description"`
	Converter   string `yaml:"converter"`
}

func loadConfigFromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config file")
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse config file")
	}

	return cfg, nil
}
