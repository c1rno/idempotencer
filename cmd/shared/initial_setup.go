package shared

import (
	"context"
	"sync"

	"github.com/c1rno/idempotencer/pkg/config"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/c1rno/idempotencer/pkg/metrics"
	"github.com/c1rno/idempotencer/pkg/signal"
)

type Setup struct {
	Waiter func()
	Log    logging.Logger
	Conf   config.Config
	Ctx    context.Context
	Cancel context.CancelFunc
	Wg     *sync.WaitGroup
}

func InitialSetup() (Setup, errors.Error) {
	var (
		s   Setup
		err error
	)
	s.Wg = &sync.WaitGroup{}
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	s.Log = logging.NewLogger(0)
	s.Conf, err = config.NewConfig(s.Log)
	if err != nil {
		return s, helpers.NewErrWithLog(s.Log, errors.UnknownError, err)
	}
	s.Log = logging.NewLogger(s.Conf.LogLevel)
	shutdown := metrics.RunMetricsSrv(
		s.Conf.Metrics,
		func() errors.Error { return nil },
		s.Log,
	)
	s.Waiter = func() {
		signal.WaitShutdown(s.Log)
		s.Cancel()
		shutdown(s.Ctx)
		s.Wg.Wait()
	}
	return s, nil
}
