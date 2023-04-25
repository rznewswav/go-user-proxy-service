package utils

func ArraySlice[T any](arr []T, from, until int) []T {
	if from < 0 {
		from = len(arr) + from
	}

	if until < 0 {
		until = len(arr) + until
	}

	if until > len(arr) {
		until = len(arr) - 1
	}

	if from < 0 || until < 0 {
		return []T{}
	}

	if from > until {
		return []T{}
	}

	return arr[from:until]
}
