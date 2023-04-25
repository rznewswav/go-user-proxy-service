package bugsnag_interfaces

type IBugsnagHandler interface {
	Notify(err error)
}
