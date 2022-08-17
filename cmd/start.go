package cmd

import (
	"github.com/hierynomus/p1-monitor/internal/config"
	p1http "github.com/hierynomus/p1-monitor/internal/http"
	"github.com/hierynomus/p1-monitor/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

func StartCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the P1 Monitor",
		Long:  "Start the P1 Monitor",
		RunE:  RunStart(cfg),
	}
}

func RunStart(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		logger := log.Ctx(cmd.Context())
		logger.Info().Str("version", version.Version).Str("commit", version.Commit).Str("date", version.Date).Msg("Starting P1 Monitor")

		telegramHandler := p1http.NewTelegramHandler(cmd.Context())
		server := p1http.NewDsmrServer(cmd.Context(), cfg.Http)
		server.AddHandler("/metrics", promhttp.Handler())
		server.AddHandler("/", telegramHandler)

		return nil
	}
}
