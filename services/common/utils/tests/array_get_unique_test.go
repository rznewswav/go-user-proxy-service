package tests

import (
	"service/services/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayGetUnique(t *testing.T) {
	input := []string{
		"hello",
		"world",
		"hello",
	}

	output := utils.ArrayGetUnique(input...)
	assert.Len(t, output, 2)
	assert.Contains(t, output, "hello")
	assert.Contains(t, output, "world")
}
