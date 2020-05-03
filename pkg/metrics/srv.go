package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	OK = "OK"
)

type (
	Checker  func() errors.Error
	Shutdown func(context.Context) errors.Error
)

func RunMetricsSrv(conf Config, checker Checker, logger logging.Logger) Shutdown {
	http.Handle(conf.MetricsRoute, promhttp.Handler())
	http.Handle(
		conf.HealthCheckRoute,
		HealthCheckHandler{
			fn:  checker,
			log: logger,
		},
	)
	srv := &http.Server{Addr: conf.Socket}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Run metrics srv", map[string]interface{}{
			"health-check-route": conf.Socket + conf.HealthCheckRoute,
			"metrics-route":      conf.Socket + conf.MetricsRoute,
		})
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			terr := errors.NewError(errors.MetricsSrvFail, err)
			logger.Error(
				terr.Error(),
				map[string]interface{}{
					"err": terr,
				},
			)
		}
		logger.Info("Shutdown metrics srv", map[string]interface{}{
			"health-check-route": conf.Socket + conf.HealthCheckRoute,
			"metrics-route":      conf.Socket + conf.MetricsRoute,
		})
	}()
	return func(ctx context.Context) errors.Error {
		err := srv.Shutdown(ctx)
		wg.Wait()
		if err != nil {
			terr := errors.NewError(errors.MetricsSrvShutdownFail, err)
			logger.Error(
				terr.Error(),
				map[string]interface{}{
					"err": terr,
				},
			)
			return terr
		}
		return nil
	}
}

type HealthCheckHandler struct {
	http.Handler
	fn  Checker
	log logging.Logger
}

func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.fn()
	if err != nil {
		h.log.Error(
			err.Error(),
			map[string]interface{}{
				"err": err,
			},
		)
		http.Error(w, err.String(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, OK)
	}
}
