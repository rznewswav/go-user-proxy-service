package bugsnag

import driver "github.com/bugsnag/bugsnag-go/v2"

type BugsnagDecorator interface {
	Decorate(
		event *driver.Event,
		config *driver.Configuration,
	)
}
