package main

import (
	"context"
	"os"

	iot "github.com/hierynomus/iot-monitor"
	"github.com/hierynomus/iot-monitor/pkg/monitor"
	"github.com/hierynomus/p1-monitor/internal/config"
	"github.com/hierynomus/p1-monitor/internal/dsmr"
	"github.com/hierynomus/p1-monitor/internal/p1metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := log.Logger.WithContext(context.Background())
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cfg := &config.Config{}
	if err := iot.StartMonitor(ctx, cfg, func() (*monitor.Monitor, error) {
		s, err := dsmr.NewDsmrReader(cfg.Serial)
		if err != nil {
			return nil, err
		}

		mp := p1metrics.Provider{}
		converter := dsmr.Converter{}

		return monitor.NewMonitor(ctx, cfg, s, mp, converter)
	}); err != nil {
		log.Fatal().Err(err).Msg("Failed to start monitor")
	}
}
