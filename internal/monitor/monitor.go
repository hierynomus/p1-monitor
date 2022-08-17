package monitor

import (
	"context"
	"os"
	"sync"
	"syscall"

	"github.com/hierynomus/p1-monitor/internal/config"
	p1http "github.com/hierynomus/p1-monitor/internal/http"
	"github.com/hierynomus/p1-monitor/internal/logging"
	"github.com/hierynomus/p1-monitor/pkg/dsmr"
	"github.com/hierynomus/p1-monitor/pkg/p1metrics"
	"github.com/ztrue/shutdown"
)

type Monitor struct {
	config    *config.Config
	reader    *dsmr.Reader
	server    *p1http.DsmrServer
	updater   *Updater
	WaitGroup *sync.WaitGroup
}

func NewMonitor(ctx context.Context, config *config.Config) (*Monitor, error) {
	ch := make(chan dsmr.RawTelegram)

	reader, err := dsmr.NewDsmrReader(config.Serial, ch)
	if err != nil {
		return nil, err
	}

	collector := p1metrics.NewCollector()

	server := p1http.NewDsmrServer(ctx, config.Http)
	if err != nil {
		return nil, err
	}

	handler := p1http.NewTelegramHandler(ctx)
	server.AddHandler("/", handler)
	promhandler, err := p1metrics.NewPrometheusHandler(ctx, collector)
	if err != nil {
		return nil, err
	}

	server.AddHandler("/metrics", promhandler)

	updater := NewUpdater(ch, handler, collector)

	return &Monitor{
		config:    config,
		reader:    reader,
		updater:   updater,
		server:    server,
		WaitGroup: &sync.WaitGroup{},
	}, nil
}

func (m *Monitor) Start(ctx context.Context) error {
	shutdown.AddWithParam(func(s os.Signal) {
		logger := logging.LoggerFor(ctx, "shutdown-hook")
		logger.Warn().Str("signal", s.String()).Msg("Received signal, shutting down")
		m.reader.Stop()

		m.reader.WaitGroup.Wait()
		m.updater.WaitGroup.Wait()
		if err := m.server.Stop(); err != nil {
			logger.Error().Err(err).Msg("Failed to gracefully stop server")
		}

		m.server.WaitGroup.Wait()

		logger.Info().Msg("All processes stopped, terminating!")
	})

	if err := m.reader.Start(ctx); err != nil {
		return err
	}

	if err := m.updater.Start(ctx); err != nil {
		return err
	}

	if err := m.server.Start(ctx); err != nil {
		return err
	}

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return nil
}
