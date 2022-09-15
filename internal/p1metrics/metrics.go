package p1metrics

import (
	"github.com/hierynomus/iot-monitor/pkg/iot"
	"github.com/hierynomus/p1-monitor/internal/config"
)

var _ iot.MetricProvider = (*Provider)(nil)

type Provider struct {
	metrics map[string]iot.MetricCollector
}

func NewProvider(cfg *config.Config) *Provider {
	return &Provider{
		metrics: buildMetrics(cfg),
	}
}

func buildMetrics(cfg *config.Config) map[string]iot.MetricCollector {
	metrics := make(map[string]iot.MetricCollector)
	for name, metric := range cfg.Metrics {
		labels := map[string]string{}
		for k, v := range metric.Labels {
			labels[k] = v
		}
		labels["dsmr_key"] = name

		if metric.Type == "gauge" {
			metrics[name] = iot.NewGaugeL(cfg.Namespace, metric.Name, metric.Description, labels)
		} else if metric.Type == "counter" {
			metrics[name] = iot.NewCounterL(cfg.Namespace, metric.Name, metric.Description, labels)
		}
	}

	return metrics
}

func (p *Provider) Metrics() map[string]iot.MetricCollector {
	return p.metrics
}
