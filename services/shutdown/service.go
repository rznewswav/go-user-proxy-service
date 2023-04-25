package shutdown

import (
	"os"
	"service/services/logger"
	"sync"
	"syscall"

	eventemitter "github.com/vansante/go-event-emitter"
)

type ShutdownHandler struct {
	count            int
	channel          chan int
	shutdownFinished map[int]bool
	name             map[int]*string
	eventEmitter     *eventemitter.Emitter
	writeLock        sync.RWMutex
}

func NewHandler() *ShutdownHandler {
	channel := make(chan int)
	shutdownFinishedMap := make(map[int]bool)
	nameMap := make(map[int]*string)
	eventEmitter := eventemitter.NewEmitter(true)

	handler := ShutdownHandler{
		count:            0,
		channel:          channel,
		shutdownFinished: shutdownFinishedMap,
		name:             nameMap,
		eventEmitter:     eventEmitter,
	}

	return &handler
}

var unnamedHook = "<unnamed>"

/*
Returns a function after the shutdown has finished
*/
func (handler *ShutdownHandler) NewShutdownHook(
	name string,
	onShutdown func(),
) func() {
	logger := logger.WithContext("shutdown")

	// eliminate this error: concurrent map writes https://app.bugsnag.com/newswav/service-worker/errors/63e5798763b1730008822c00?filters[event.since]=30d&filters[error.status]=open&filters[app.release_stage]=production
	// using this solution: https://stackoverflow.com/questions/45585589/golang-fatal-error-concurrent-map-read-and-map-write
	handler.writeLock.RLock()
	id := handler.count
	handler.writeLock.RUnlock()

	handler.writeLock.Lock()
	handler.count += 1
	handler.shutdownFinished[id] = false
	handler.name[id] = &name
	handler.writeLock.Unlock()

	handler.eventEmitter.AddListener(
		"shutdown",
		func(arguments ...interface{}) {
			logger.Debug(
				"onShutdown: %s",
				name,
			)
			onShutdown()
		},
	)

	logger.Debug(
		"Added shutdown hook for event name: %s",
		name,
	)

	return func() {
		handler.channel <- id
	}
}

func (handler *ShutdownHandler) IsAllShutdownHookCalled() bool {
	someHasNotFinished := false
	for k := range handler.shutdownFinished {
		someHasNotFinished = someHasNotFinished ||
			!handler.shutdownFinished[k]
	}
	return !someHasNotFinished
}

func (handler *ShutdownHandler) Start(
	onFinishedShutdown func(),
) {
	logger := logger.WithContext("shutdown")
	for {
		if handler.IsAllShutdownHookCalled() {
			close(handler.channel)
			logger.Info("shutting down...")
			onFinishedShutdown()
			return
		}
		shutdownId := <-handler.channel
		if handler.shutdownFinished[shutdownId] {
			name := handler.name[shutdownId]
			if name == nil {
				name = &unnamedHook
			}
			logger.Warn(
				"Shutdown hook for %s was called more than once!",
				*name,
			)
		}

		handler.shutdownFinished[shutdownId] = true
	}
}

func (handler *ShutdownHandler) RequestShutdown(
	reason string,
	isExpected ...bool,
) {
	logger := logger.WithContext("shutdown")
	if len(isExpected) > 0 && isExpected[0] {
		logger.Info(
			"Silent shutdown was requested with reason: %s",
			reason,
		)
	} else {
		logger.Error(
			"Shutdown was requested with reason: %s",
			reason,
		)
	}
	if p, err := os.FindProcess(syscall.Getpid()); err != nil {
		logger.Error(
			"Shutdown was requested but current process id cannot be queried with reason: %s",
			err.Error(),
		)
	} else {
		p.Signal(os.Interrupt)
	}
}

func (handler *ShutdownHandler) StartShutdown() {
	handler.eventEmitter.EmitEvent("shutdown")
}
