package utils

import "github.com/mariomac/gostream/stream"

func StreamFilterNilAsArray[T any](
	nilableStream stream.Stream[*T],
) []T {
	nonNilSlices := stream.Map(
		nilableStream.
			Filter(
				func(t *T) bool {
					return t != nil
				},
			),
		func(filteredStream *T) T {
			return *filteredStream
		},
	)

	return nonNilSlices.ToSlice()
}
