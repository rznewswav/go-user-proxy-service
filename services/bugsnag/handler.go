package bugsnag

import (
	"fmt"
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
)

type Handler interface {
	Notify(err error)
}

type handler struct {
	Handler
}

func (handler *handler) Notify(source error) {
	err := bugsnag.Notify(source)
	if err != nil {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintf(
			os.Stderr,
			"FATAL: Cannot notify bugsnag for error\n  %s\n\nReceived error from bugsnag:\n  %s\n",
			source.Error(),
			err.Error(),
		)
	}
}
