package controllers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockController_ReplaceMiddleware(t *testing.T) {
	mockHeaderKey := "middleware type"
	var ActualMiddleware Handler[any] = func(
		Request Request[any],
		SetStatus SetStatus,
		SetHeader SetHeader) (Response any) {
		SetHeader(mockHeaderKey, "actual")
		return nil
	}

	var ReplacementMiddleware Handler[any] = func(
		Request Request[any],
		SetStatus SetStatus,
		SetHeader SetHeader) (Response any) {
		SetHeader(mockHeaderKey, "replacement")
		return nil
	}

	var Controller = C[any]().UseMiddleware(&ActualMiddleware)

	var MockController = Mock[any](Controller).ReplaceMiddleware(&ActualMiddleware, &ReplacementMiddleware)

	if !assert.Len(t, MockController.Middlewares, 1) {
		return
	}

	assert.Equal(t, MockController.Middlewares[0], &ReplacementMiddleware)
}
