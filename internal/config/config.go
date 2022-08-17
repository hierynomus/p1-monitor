package config

import (
	p1http "github.com/hierynomus/p1-monitor/internal/http"
	"github.com/hierynomus/p1-monitor/pkg/dsmr"
)

type Config struct {
	Serial dsmr.Config   `yaml:"serial" viper:"serial"`
	Http   p1http.Config `yaml:"http" viper:"http"`
}
