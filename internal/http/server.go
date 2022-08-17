package http

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/hierynomus/p1-monitor/internal/logging"
)

type DsmrServer struct {
	ctx       context.Context
	config    Config
	WaitGroup *sync.WaitGroup
	srv       *http.Server
	handlers  map[string]http.Handler
}

func NewDsmrServer(ctx context.Context, c Config) *DsmrServer {
	return &DsmrServer{
		ctx:       ctx,
		config:    c,
		WaitGroup: &sync.WaitGroup{},
		srv: &http.Server{
			Addr:              c.ListenAddress,
			Handler:           nil,
			ReadTimeout:       c.Timeout,
			ReadHeaderTimeout: c.Timeout,
			WriteTimeout:      c.Timeout,
		},
		handlers: make(map[string]http.Handler),
	}
}

func (s *DsmrServer) AddHandler(name string, handler http.Handler) {
	s.handlers[name] = handler
}

func (s *DsmrServer) Start(ctx context.Context) error {
	s.WaitGroup.Add(1)

	go s.run(ctx)

	return nil
}

func (s *DsmrServer) Stop() error {
	return s.srv.Shutdown(s.ctx)
}

func (s *DsmrServer) run(_ context.Context) {
	defer s.WaitGroup.Done()
	logger := logging.LoggerFor(s.ctx, "dsmr-server")
	mux := http.NewServeMux()

	logger.Info().Msg("Starting dsmr server")

	for name, handler := range s.handlers {
		mux.Handle(name, handler)
	}

	s.srv.Handler = mux

	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.Info().Msg("dsmr server stopped")
	} else {
		logger.Error().Err(err).Msg("dsmr server failed")
	}
}
