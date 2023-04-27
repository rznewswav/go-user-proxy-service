package controllers

import (
	"github.com/google/uuid"
)

const RequestIdToken = "x-req-id"

var AssignRequestId Handler[any] = func(
	Request Request[any],
	SetStatus SetStatus,
	SetHeader SetHeader,
) (Response any) {
	existingRequestId, _ := Request.Get(RequestIdToken)
	if existingRequestId != nil && len(existingRequestId.(string)) > 0 {
		return
	}
	id := uuid.New()
	idString := id.String()

	Request.Set(RequestIdToken, idString)
	SetHeader(RequestIdToken, idString)
	return nil
}
