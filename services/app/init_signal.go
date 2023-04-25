package app

import (
	"os"
	"os/signal"
	"service/services/logger"
	"service/services/shutdown"
	"syscall"
)

func InitSignal(
	shutdownHandler *shutdown.ShutdownHandler,
) func() {
	logger := logger.WithContext("signal")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		signal := <-sigchan
		println("")
		logger.Debug(
			"Caught signal %v: terminating",
			signal,
		)
		shutdownHandler.StartShutdown()
	}()
	logger.Info("initialized signal handler")
	return func() {
		close(sigchan)
	}
}
