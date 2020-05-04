package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns = "lynis_"
)

var (
	regex   *regexp.Regexp
	ageDesc *prometheus.Desc
)

func init() {
	regex = regexp.MustCompile(`^([^=]+)=(.*)$`)
	ageDesc = prometheus.NewDesc(ns+"report_age_seconds", "Report age in seconds", nil, nil)
}

type collector struct {
	cfg     *Config
	descs   []*prometheus.Desc
	metrics []prometheus.Metric
}

func newCollector(cfg *Config) *collector {
	return &collector{
		cfg:   cfg,
		descs: []*prometheus.Desc{ageDesc},
	}
}

func (c *collector) collect() error {
	info, err := os.Stat(c.cfg.ReportFilePath)
	if err != nil {
		return errors.Wrap(err, "cloud not get file info for report file")
	}

	age := float64(time.Since(info.ModTime()).Seconds())
	c.metrics = append(c.metrics, prometheus.MustNewConstMetric(ageDesc, prometheus.GaugeValue, age))

	b, err := ioutil.ReadFile(c.cfg.ReportFilePath)
	if err != nil {
		return errors.Wrap(err, "cloud not read report file")
	}

	r := bytes.NewReader(b)
	s := bufio.NewScanner(r)

	for s.Scan() {
		m := regex.FindStringSubmatch(s.Text())

		if len(m) == 0 {
			continue
		}

		err = c.parseForMetric(m[1], m[2])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *collector) parseForMetric(field, value string) error {
	def, found := c.cfg.Metrics[field]
	if !found {
		return nil
	}

	name := field
	if def.MetricName != "" {
		name = def.MetricName
	}

	desc := prometheus.NewDesc(ns+name, def.Description, nil, nil)
	c.descs = append(c.descs, desc)

	conv := converterByName(def.Converter)
	v, err := conv(value)
	if err != nil {
		errors.Wrapf(err, "could not parse value (%s) from field %s", value, field)
	}

	m := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v)
	c.metrics = append(c.metrics, m)

	return nil
}

// Describe implements prometheus.Collector interface
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descs {
		ch <- d
	}
}

// Collect implements prometheus.Collector interface
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.metrics {
		ch <- m
	}
}
