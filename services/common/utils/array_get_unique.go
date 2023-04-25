package utils

import (
	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/stream"
	"golang.org/x/exp/constraints"
)

func ArrayGetUnique[T constraints.Ordered](input ...T) []T {
	keysMap := make(map[T]bool)
	for _, v := range input {
		keysMap[v] = true
	}
	keysMapStream := stream.OfMap(keysMap)
	return stream.
		Map(keysMapStream, func(p item.Pair[T, bool]) T {
			return p.Key
		}).ToSlice()
}
