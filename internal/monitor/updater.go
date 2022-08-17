package monitor

import (
	"context"
	"sync"

	p1http "github.com/hierynomus/p1-monitor/internal/http"
	"github.com/hierynomus/p1-monitor/internal/logging"
	p1dsmr "github.com/hierynomus/p1-monitor/pkg/dsmr"
	"github.com/hierynomus/p1-monitor/pkg/p1metrics"
	"github.com/roaldnefs/go-dsmr"
)

type Updater struct {
	WaitGroup *sync.WaitGroup
	ch        <-chan p1dsmr.RawTelegram
	handler   *p1http.TelegramHandler
	collector *p1metrics.Collector
}

func NewUpdater(ch <-chan p1dsmr.RawTelegram, handler *p1http.TelegramHandler, collector *p1metrics.Collector) *Updater {
	return &Updater{
		WaitGroup: &sync.WaitGroup{},
		ch:        ch,
		handler:   handler,
		collector: collector,
	}
}

func (u *Updater) Start(ctx context.Context) error {
	u.WaitGroup.Add(1)

	go u.run(ctx)

	return nil
}

func (u *Updater) run(ctx context.Context) {
	logger := logging.LoggerFor(ctx, "updater")
	defer u.WaitGroup.Done()

	for { //nolint:gosimple
		select {
		case t, ok := <-u.ch:
			if !ok {
				logger.Info().Msg("Updater channel closed, terminating!")
				return
			}

			parsedTelegram, err := dsmr.ParseTelegram(string(t))
			if err != nil {
				logger.Error().Err(err).Msg("Failed to parse telegram")
				continue
			}

			logger.Debug().Msg("Parsed telegram")

			u.handler.UpdateTelegram(string(t))
			u.collector.Update(parsedTelegram)
		}
	}
}
