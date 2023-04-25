package bugsnag_structs

type DevBugsnagHandler struct {
}

func (handler *DevBugsnagHandler) Notify(
	err error,
) {
	// a noop so that bugsnag.Notify will not send notification to
	// bugsnag on development mode
}
