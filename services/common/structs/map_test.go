package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	var h StringDefaultedMap

	h.Set("hello", "world")

	if !assert.NotNil(t, h.H) {
		return
	}
	assert.Contains(t, *h.H, "hello")
	assert.Equal(t, (*h.H)["hello"], "world")
}

func TestGetHeaders(t *testing.T) {
	var h StringDefaultedMap

	h.Set("hello", "world")

	assert.Equal(t, h.Get("hello"), "world")
}

func TestGetNilHeaders(t *testing.T) {
	var h StringDefaultedMap

	h.Set("hello", "world")

	assert.Equal(t, h.Get("nil"), "")
}
