package config

import (
	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/hierynomus/p1-monitor/internal/dsmr"
)

var _ config.MonitorConfig = (*Config)(nil)

type Config struct {
	Serial dsmr.Config       `yaml:"serial" viper:"serial"`
	Http   config.HTTPConfig `yaml:"http" viper:"http"`
}

func (c *Config) HTTP() config.HTTPConfig {
	return c.Http
}

func (c *Config) RawMessageContentType() string {
	return "text/plain"
}
