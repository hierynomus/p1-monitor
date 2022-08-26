package p1metrics

import (
	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "p1"
)

var _ exporter.MetricProvider = (*Provider)(nil)

type Provider struct{}

func NewGaugeL(name, help string, labels map[string]string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        name,
		Help:        help,
		ConstLabels: prometheus.Labels(labels),
	})
}

func labels(kvs ...string) map[string]string {
	m := make(map[string]string)

	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return m
}

// MetricBuilders contains builders for all object types in a DSMR frame.
func (p Provider) Metrics() map[string]exporter.MetricCollector {
	return map[string]exporter.MetricCollector{
		"1-0:1.8.1":   NewGaugeL("electricity_delivered_to_client_kwh", "Meter Reading electricity delivered to client in 0,001 kWh", labels("tariff", "1", "dsmr_key", "1-0:1.8.1")),
		"1-0:1.8.2":   NewGaugeL("electricity_delivered_to_client_kwh", "Meter Reading electricity delivered to client in 0,001 kWh", labels("tariff", "2", "dsmr_key", "1-0:1.8.2")),
		"1-0:2.8.1":   NewGaugeL("electricity_delivered_by_client_kwh", "Meter Reading electricity delivered by client in 0,001 kWh", labels("tariff", "1", "dsmr_key", "1-0:2.8.1")),
		"1-0:2.8.2":   NewGaugeL("electricity_delivered_by_client_kwh", "Meter Reading electricity delivered by client in 0,001 kWh", labels("tariff", "2", "dsmr_key", "1-0:2.8.2")),
		"0-0:96.14.0": NewGaugeL("tariff_indicator_electricity", "Tariff indicator electricity", labels("dsmr_key", "0-0:96.14.0")),
		"1-0:1.7.0":   NewGaugeL("electricity_power_delivered_kw", "Actual electricity power delivered (+P) in 1 Watt resolution", labels("dsmr_key", "1-0:1.7.0")),
		"1-0:2.7.0":   NewGaugeL("electricity_power_received_kw", "Actual electricity power received (-P) in 1 Watt resolution", labels("dsmr_key", "1-0:2.7.0")),
		"0-0:17.0.0":  NewGaugeL("threshold_electricity_kw", "The actual threshold Electricity in kW", labels("dsmr_key", "0-0:17.0.0")),
		"0-0:96.7.21": NewGaugeL("power_failures_total", "Number of power failures in any phase", labels("dsmr_key", "0-0:96.7.21")),
		"0-0:96.7.9":  NewGaugeL("long_power_failures_total", "Number of long power failures in any phase", labels("dsmr_key", "0-0:96.7.9")),
		"1-0:32.32.0": NewGaugeL("voltage_sags", "Number of voltage sags", labels("phase", "L1", "dsmr_key", "1-0:32.32.0")),
		"1-0:52.32.0": NewGaugeL("voltage_sags", "Number of voltage sags", labels("phase", "L2", "dsmr_key", "1-0:52.32.0")),
		"1-0:72:32.0": NewGaugeL("voltage_sags", "Number of voltage sags", labels("phase", "L3", "dsmr_key", "1-0:72:32.0")),
		"1-0:32.36.0": NewGaugeL("voltage_swells", "Number of voltage swells", labels("phase", "L1", "dsmr_key", "1-0:32.36.0")),
		"1-0:52.36.0": NewGaugeL("voltage_swells", "Number of voltage swells", labels("phase", "L2", "dsmr_key", "1-0:52.36.0")),
		"1-0:72.36.0": NewGaugeL("voltage_swells", "Number of voltage swells", labels("phase", "L3", "dsmr_key", "1-0:72.36.0")),
		"1-0:31.7.0":  NewGaugeL("current_a", "Instantaneous current in A resolution", labels("phase", "L1", "dsmr_key", "1-0:31.7.0")),
		"1-0:51.7.0":  NewGaugeL("current_a", "Instantaneous current in A resolution", labels("phase", "L2", "dsmr_key", "1-0:51.7.0")),
		"1-0:71.7.0":  NewGaugeL("current_a", "Instantaneous current in A resolution", labels("phase", "L3", "dsmr_key", "1-0:71.7.0")),
		"1-0:32.7.0":  NewGaugeL("voltage_v", "Instantaneous voltage in V resolution", labels("phase", "L1", "dsmr_key", "1-0:32.7.0")),
		"1-0:52.7.0":  NewGaugeL("voltage_v", "Instantaneous voltage in V resolution", labels("phase", "L2", "dsmr_key", "1-0:52.7.0")),
		"1-0:72.7.0":  NewGaugeL("voltage_v", "Instantaneous voltage in V resolution", labels("phase", "L3", "dsmr_key", "1-0:72.7.0")),
		"1-0:21.7.0":  NewGaugeL("active_power_delivered_kw", "Instantaneous active power (+P) in W resolution", labels("phase", "L1", "dsmr_key", "1-0:21.7.0")),
		"1-0:41.7.0":  NewGaugeL("active_power_delivered_kw", "Instantaneous active power (+P) in W resolution", labels("phase", "L2", "dsmr_key", "1-0:41.7.0")),
		"1-0:61.7.0":  NewGaugeL("active_power_delivered_kw", "Instantaneous active power (+P) in W resolution", labels("phase", "L3", "dsmr_key", "1-0:61.7.0")),
		"1-0:22.7.0":  NewGaugeL("active_power_received_kw", "Instantaneous active power (-P) in W resolution", labels("phase", "L1", "dsmr_key", "1-0:22.7.0")),
		"1-0:42.7.0":  NewGaugeL("active_power_received_kw", "Instantaneous active power (-P) in W resolution", labels("phase", "L2", "dsmr_key", "1-0:42.7.0")),
		"1-0:62.7.0":  NewGaugeL("active_power_received_kw", "Instantaneous active power (-P) in W resolution", labels("phase", "L3", "dsmr_key", "1-0:62.7.0")),
		"0-1:24.2.3":  NewGaugeL("gas_m3", "Actual gas volume delivered", labels("dsmr_key", "0-1:24.2.3")),
	}
}
