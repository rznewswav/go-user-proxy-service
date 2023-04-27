package logger

import (
	bugsnag_structs "service/services/bugsnag"
)

func SplitPayloadIntoFormatterAndNotifyableError(
	payload []any,
) ([]any, *bugsnag_structs.NotifiableError) {
	if len(payload) <= 0 {
		return payload, nil
	}

	lastOfPayload := payload[len(payload)-1]

	if ne, ok := lastOfPayload.(bugsnag_structs.NotifiableError); ok {
		return payload[:len(payload)-1], &ne
	}

	if ne, ok := lastOfPayload.(*bugsnag_structs.NotifiableError); ok {
		return payload[:len(payload)-1], ne
	}

	return payload, nil
}
