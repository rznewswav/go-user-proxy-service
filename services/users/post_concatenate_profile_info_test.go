package users

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"service/services/server/controllers"
	"testing"
)

func TestConcatenateProfileInfo(t *testing.T) {
	mockController := controllers.Mock[map[string]interface{}](ConcatenateProfileInfo)
	userProfileValue := "hello"

	var MockUserInfo controllers.Handler[any] = func(
		Request controllers.Request[any],
		SetStatus controllers.SetStatus,
		SetHeader controllers.SetHeader,
	) (Response any) {
		Request.Set(UserProfileToken, userProfileValue)
		return nil
	}

	mockController.ReplaceMiddleware(&AuthMiddleware, &MockUserInfo)
	response, status, _ := mockController.SendMockRequest()
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, response, "profile")
}
