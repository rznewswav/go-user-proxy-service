package bugsnag

type silentHandler struct {
}

func (handler *silentHandler) Notify(
	_ error,
) {
	// a noop so that bugsnag.Notify will not send a notification to
	// bugsnag on development mode
}
