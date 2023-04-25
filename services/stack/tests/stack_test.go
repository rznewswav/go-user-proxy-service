package main

import (
	"service/services/stack"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	stackTrace := stack.GetStackTrace()

	for _, st := range stackTrace {
		assert.NotNil(t, st)
	}

	if !assert.Greater(
		t,
		len(stackTrace),
		0,
		"Generated stack trace must be at least one!",
	) {
		return
	}

	assert.Contains(t, stackTrace[0].Method, "TestStack")
}
