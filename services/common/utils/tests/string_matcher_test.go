package tests

import (
	"service/services/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMatcherEn(t *testing.T) {
	matcher := utils.NewStringMatcher(
		[]string{"hello", "world"},
		true,
	)

	if !assert.NotNil(t, matcher) {
		// hard fail
		return
	}

	assert.True(t, matcher.HasInDocument("hello"))
	assert.True(t, matcher.HasInDocument("world"))
	assert.False(t, matcher.HasInDocument("llow"))
	assert.False(t, matcher.HasInDocument("malaysia"))
}

func TestStringMatcherZn(t *testing.T) {
	matcher := utils.NewStringMatcher(
		[]string{"武汉", "武汉病毒"},
		false,
	)

	if !assert.NotNil(t, matcher) {
		// hard fail
		return
	}

	assert.True(t, matcher.HasInDocument("行动武汉"))
	assert.True(t, matcher.HasInDocument("行动武汉病"))
	assert.False(t, matcher.HasInDocument("病行"))
	assert.False(t, matcher.HasInDocument("非同質化代幣"))
}
