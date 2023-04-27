package bugsnag

import driver "github.com/bugsnag/bugsnag-go/v2"

type Decorator interface {
	Decorate(
		event *driver.Event,
		config *driver.Configuration,
	)
}
