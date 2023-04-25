package utils

import (
	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/stream"
	"golang.org/x/exp/constraints"
)

func MapAllValuesTrue[T constraints.Ordered](
	m map[T]bool,
) bool {
	return stream.
		OfMap(m).
		AllMatch(
			func(p item.Pair[T, bool]) bool {
				return p.Val
			},
		)
}
