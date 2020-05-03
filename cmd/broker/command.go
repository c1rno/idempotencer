package broker

import (
	"context"

	"github.com/c1rno/idempotencer/pkg/config"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/c1rno/idempotencer/pkg/metrics"
	"github.com/c1rno/idempotencer/pkg/signal"
	"github.com/spf13/cobra"
	_ "github.com/pebbe/zmq4"
)

var Command = &cobra.Command{
	Use:   `broker`,
	Short: `0MQ broker, needs to load-balancing`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.NewLogger(0)
		conf, _ := config.NewConfig(logger)
		logger = logging.NewLogger(conf.LogLevel)

		shutdown := metrics.RunMetricsSrv(conf.Metrics, func() errors.Error { return nil }, logger)
		defer shutdown(context.Background())


		signal.WaitShutdown(logger)
	},
}
