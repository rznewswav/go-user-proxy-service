package tests

import (
	"service/services/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentFunc(t *testing.T) {
	val := utils.GetCurrentFuncName()

	if !assert.NotNil(t, val) {
		// hard fail
		return
	}

	assert.Equal(t, "TestGetCurrentFunc", val)
}
