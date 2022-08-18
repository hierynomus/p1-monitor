package cmd

import (
	"github.com/hierynomus/p1-monitor/internal/config"
	"github.com/hierynomus/p1-monitor/internal/monitor"
	"github.com/hierynomus/p1-monitor/version"
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

		m, err := monitor.NewMonitor(cmd.Context(), cfg)
		if err != nil {
			return err
		}

		return m.Start(cmd.Context())
	}
}
