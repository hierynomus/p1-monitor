package dsmr

import (
	"github.com/hierynomus/iot-monitor/pkg/iot"
	"github.com/roaldnefs/go-dsmr"
)

var _ iot.Converter = (*Converter)(nil)

type Converter struct {
}

func (c Converter) Convert(rawMessage iot.RawMessage) (iot.MetricMessage, error) {
	parsedTelegram, err := dsmr.ParseTelegram(string(rawMessage))
	if err != nil {
		return nil, err
	}

	msg := map[string]iot.Metric{}
	for k, v := range parsedTelegram.DataObjects {
		msg[k] = iot.Metric{
			Value: v.Value,
			Unit:  v.Unit,
		}
	}

	return msg, nil
}
