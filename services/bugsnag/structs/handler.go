package bugsnag_structs

import (
	"fmt"
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
)

type BugsnagHandler struct {
}

func (handler *BugsnagHandler) Notify(source error) {
	err := bugsnag.Notify(source)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"FATAL: Cannot notify bugsnag for error\n  %s\n\nReceived error from bugsnag:\n  %s\n",
			source.Error(),
			err.Error(),
		)
	}
}
