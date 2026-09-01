// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [Apache 2.0](LICENSE) license.
//
// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v1"
)

// Config represents the exporter configuration file
type Config struct {
	ReportFilePath string                      `yaml:"reportfile_path"`
	Metrics        map[string]MetricDefinition `yaml:"metrics"`
}

// MetricDefinition defines the source and conversions for a single metric
type MetricDefinition struct {
	MetricName  string `yaml:"name"`
	Description string `yaml:"description"`
	Converter   string `yaml:"converter"`
}

func loadConfigFromFile(path string) (*Config, error) {
	b, err := os.ReadFile(path) // #nosec G304 -- path is an operator-controlled CLI flag, not untrusted input
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
