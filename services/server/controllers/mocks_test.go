package controllers

import (
	"github.com/stretchr/testify/assert"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
	"testing"
)

func TestMockController_ReplaceMiddleware(t *testing.T) {
	mockHeaderKey := "middleware type"
	var ActualMiddleware handlers.Handler[any] = func(
		Request req.Request[any]) (Response resp.Response) {
		return resp.N().Header(mockHeaderKey, "actual")
	}

	var ReplacementMiddleware handlers.Handler[any] = func(
		Request req.Request[any]) (Response resp.Response) {
		return resp.N().Header(mockHeaderKey, "replacement")
	}

	var Controller = C[any]().UseMiddleware(&ActualMiddleware)

	var MockController = Mock[any](Controller).ReplaceMiddleware(&ActualMiddleware, &ReplacementMiddleware)

	if !assert.Len(t, MockController.Middlewares, 1) {
		return
	}

	assert.Equal(t, MockController.Middlewares[0], &ReplacementMiddleware)
}
