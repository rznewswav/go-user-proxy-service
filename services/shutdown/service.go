package shutdown

import (
	"os"
	"service/services/logger"
	"sync"
	"syscall"

	eventemitter "github.com/vansante/go-event-emitter"
)

type Handler struct {
	count            int
	channel          chan int
	shutdownFinished map[int]bool
	name             map[int]*string
	eventEmitter     *eventemitter.Emitter
	writeLock        sync.RWMutex
}

func NewHandler() *Handler {
	channel := make(chan int)
	shutdownFinishedMap := make(map[int]bool)
	nameMap := make(map[int]*string)
	eventEmitter := eventemitter.NewEmitter(true)

	handler := Handler{
		count:            0,
		channel:          channel,
		shutdownFinished: shutdownFinishedMap,
		name:             nameMap,
		eventEmitter:     eventEmitter,
	}

	return &handler
}

var unnamedHook = "<unnamed>"

// NewShutdownHook Returns a function after the shutdown has finished
func (handler *Handler) NewShutdownHook(
	name string,
	onShutdown func(),
) func() {
	l := logger.For("shutdown")

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
			l.Debug(
				"onShutdown: %s",
				name,
			)
			onShutdown()
		},
	)

	l.Debug(
		"Added shutdown hook for event name: %s",
		name,
	)

	return func() {
		handler.channel <- id
	}
}

func (handler *Handler) IsAllShutdownHookCalled() bool {
	someHasNotFinished := false
	for k := range handler.shutdownFinished {
		someHasNotFinished = someHasNotFinished ||
			!handler.shutdownFinished[k]
	}
	return !someHasNotFinished
}

func (handler *Handler) Start(
	onFinishedShutdown func(),
) {
	l := logger.For("shutdown")
	for {
		if handler.IsAllShutdownHookCalled() {
			close(handler.channel)
			l.Info("shutting down...")
			onFinishedShutdown()
			return
		}
		shutdownId := <-handler.channel
		if handler.shutdownFinished[shutdownId] {
			name := handler.name[shutdownId]
			if name == nil {
				name = &unnamedHook
			}
			l.Warn(
				"Shutdown hook for %s was called more than once!",
				*name,
			)
		}

		handler.shutdownFinished[shutdownId] = true
	}
}

func (handler *Handler) RequestShutdown(
	reason string,
	isExpected ...bool,
) {
	l := logger.For("shutdown")
	if len(isExpected) > 0 && isExpected[0] {
		l.Info(
			"Silent shutdown was requested with reason: %s",
			reason,
		)
	} else {
		l.Error(
			"Shutdown was requested with reason: %s",
			reason,
		)
	}
	if p, err := os.FindProcess(syscall.Getpid()); err != nil {
		l.Error(
			"Shutdown was requested but current process id cannot be queried with reason: %s",
			err.Error(),
		)
	} else {
		p.Signal(os.Interrupt)
	}
}

func (handler *Handler) StartShutdown() {
	handler.eventEmitter.EmitEvent("shutdown")
}
