package middlewares

import (
	"github.com/google/uuid"
	"service/services/server/constants"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
)

var AssignRequestId handlers.Handler[any] = func(
	Request req.Request[any],
) (Response resp.Response) {
	existingRequestId, _ := Request.Get(constants.RequestIdToken)
	if existingRequestId != nil && len(existingRequestId.(string)) > 0 {
		return
	}
	id := uuid.New()
	idString := id.String()

	Request.Set(constants.RequestIdToken, idString)
	return nil
}
