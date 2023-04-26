package utils

func ArrayGetOrNil[T any](arr []T, at int) *T {
	if at < 0 {
		at = len(arr) + at
	}
	if at >= len(arr) {
		return nil
	}

	return &arr[at]
}
