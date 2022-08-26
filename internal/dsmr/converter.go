package dsmr

import (
	"github.com/hierynomus/iot-monitor/pkg/converter"
	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
	"github.com/roaldnefs/go-dsmr"
)

var _ converter.Converter = (*Converter)(nil)

type Converter struct {
}

func (c Converter) Convert(rawMessage scraper.RawMessage) (exporter.MetricMessage, error) {
	parsedTelegram, err := dsmr.ParseTelegram(string(rawMessage))
	if err != nil {
		return nil, err
	}

	msg := map[string]exporter.Metric{}
	for k, v := range parsedTelegram.DataObjects {
		msg[k] = exporter.Metric{
			Value: v.Value,
			Unit:  v.Unit,
		}
	}

	return msg, nil
}
