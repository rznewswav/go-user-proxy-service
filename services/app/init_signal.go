package app

import (
	"os"
	"os/signal"
	"service/services/logger"
	"service/services/shutdown"
	"syscall"
)

func InitSignal(
	shutdownHandler *shutdown.Handler,
) func() {
	l := logger.For("signal")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sigchan
		println("")
		l.Debug(
			"Caught s %v: terminating",
			s,
		)
		shutdownHandler.StartShutdown()
	}()
	l.Info("initialized signal handler")
	return func() {
		close(sigchan)
	}
}
