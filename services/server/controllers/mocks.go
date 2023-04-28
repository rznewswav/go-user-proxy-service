package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"service/services/common/structs"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
)

type MockController[T any] struct {
	Controller[T]
}

func Mock[T any](controller Controller[T]) MockController[T] {
	cloned := controller.Clone()
	return MockController[T]{
		Controller: cloned,
	}
}

func (mc MockController[T]) ReplaceMiddleware(
	target *handlers.Handler[any],
	replaceWith *handlers.Handler[any],
) MockController[T] {
	for index, middleware := range mc.Middlewares {
		if target == middleware {
			mc.Middlewares[index] = replaceWith
		}
	}
	return mc
}

type MockConfig struct {
	Key   string
	Value string
}

type MockAddHeader MockConfig
type MockBody any

func (mc MockController[T]) SendMockRequest(opt ...any) (
	response gin.H,
	status int,
	header structs.StringDefaultedMap,
) {
	status = http.StatusOK

	var body T
	var requestHeaders structs.StringDefaultedMap

	for _, castable := range opt {
		switch casted := castable.(type) {
		case MockBody:
			body = casted.(T)
		case MockAddHeader:
			requestHeaders.Set(casted.Key, casted.Value)
		}
	}

	ctx := make(map[string]interface{})
	requestForMiddie := req.WrapRequestMockBody[any](
		body,
		ctx,
		&requestHeaders,
	)
	for index, middleware := range mc.Middlewares {
		middleResponse := (*middleware)(
			requestForMiddie,
		)

		if middleResponse == nil {
			continue
		}

		if castedResponse, castable := middleResponse.(resp.Response); castable {
			payload := castedResponse.GetResponsePayload()
			payload.Header.ForEach(func(key, value string) {
				header.Set(key, value)
			})
			status = payload.Status
			if castedResponse.Next() {
				continue
			} else {
				response = payload.Data
			}
		} else {
			panic(fmt.Sprintf("middleware of index %d is not returning data of type resp.Response", index))
		}
	}

	request := req.WrapRequestMockBody[T](
		body,
		ctx,
		&requestHeaders,
	)
	handlerResponse := mc.Handler(
		request,
	)

	if castedResponse, castable := handlerResponse.(resp.Response); castable {
		payload := castedResponse.GetResponsePayload()
		payload.Header.ForEach(func(key, value string) {
			header.Set(key, value)
		})
		status = payload.Status
		response = payload.Data
	} else {
		panic("handler is not returning data of type resp.Response")
	}
	return
}
