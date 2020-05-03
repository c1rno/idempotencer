package migrate

import (
	"context"

	"github.com/c1rno/idempotencer/pkg/config"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/c1rno/idempotencer/pkg/metrics"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   `migrate`,
	Short: `Setup persistence layer schema`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.NewLogger(0)
		conf, _ := config.NewConfig(logger)
		logger = logging.NewLogger(conf.LogLevel)

		shutdown := metrics.RunMetricsSrv(conf.Metrics, func() errors.Error { return nil }, logger)
		defer shutdown(context.Background())
	},
}
