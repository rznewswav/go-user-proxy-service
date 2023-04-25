package tests

import (
	"service/services/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayNonNil(t *testing.T) {
	str1 := "hello"
	str2 := "world"

	nilableStringArray := []*string{
		&str1,
		nil,
		&str2,
		nil,
	}

	filteredArray := utils.ArrayFilterNil(nilableStringArray)

	assert.Len(t, filteredArray, 2)
	assert.Equal(t, filteredArray[0], str1)
	assert.Equal(t, filteredArray[1], str2)
}
