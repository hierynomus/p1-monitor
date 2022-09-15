package config

import (
	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/hierynomus/p1-monitor/internal/dsmr"
)

var _ config.MonitorConfig = (*Config)(nil)

type Config struct {
	Serial    dsmr.Config       `yaml:"serial" viper:"serial"`
	Http      config.HTTPConfig `yaml:"http" viper:"http"` //nolint:revive
	Namespace string            `yaml:"namespace" viper:"namespace"`
	Metrics   map[string]Metric `yaml:"metrics" viper:"metrics"`
}

func (c *Config) HTTP() config.HTTPConfig {
	return c.Http
}

func (c *Config) RawMessageContentType() string {
	return "text/plain"
}

type Metric struct {
	Name string `yaml:"name"`
	Unit string `yaml:"unit"`
	// Type is one of: gauge, counter, histogram, summary
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	// Labels are the labels of the metric
	Labels map[string]string `yaml:"labels"`
}
