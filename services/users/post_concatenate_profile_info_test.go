package users

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"service/services/server/controllers"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
	"testing"
)

func TestConcatenateProfileInfo(t *testing.T) {
	mockController := controllers.Mock[map[string]interface{}](ConcatenateProfileInfo)
	userProfileValue := "hello"

	var MockUserInfo handlers.Handler[any] = func(
		Request req.Request[any],
	) (Response resp.Response) {
		Request.Set(UserProfileToken, userProfileValue)
		return nil
	}

	mockController.ReplaceMiddleware(&AuthMiddleware, &MockUserInfo)
	response, status, _ := mockController.SendMockRequest(
		controllers.MockBody(map[string]interface{}{"hello": "world"}),
	)
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, response, "body")
	body := response["body"]
	assert.Contains(t, body, "hello")
	hello := body.(map[string]interface{})["hello"]
	assert.Equal(t, hello, "world")

	assert.Contains(t, response, "profile")
	profile := response["profile"]
	assert.Equal(t, userProfileValue, profile)
}
