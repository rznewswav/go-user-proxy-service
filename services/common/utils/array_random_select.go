package utils

func ArrayRandomSelect[T any](arr []T) T {
	return arr[Random.Int63n(int64(len(arr)))]
}
