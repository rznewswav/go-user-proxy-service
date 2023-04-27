package utils

import "github.com/mariomac/gostream/stream"

// ArrayFilterNil filters out all nil values in an array
// and returns an array of non-nil items
func ArrayFilterNil[T any](nilable []*T) []T {
	return StreamFilterNilAsArray(stream.OfSlice(nilable))
}
