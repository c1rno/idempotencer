package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/c1rno/idempotencer/pkg/logging"
)

func WaitShutdown(logger logging.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	logger.Info("Shutdown", map[string]interface{}{
		"sig": <-ch,
	})
}
