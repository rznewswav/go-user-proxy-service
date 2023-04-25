package utils

func FirstOfNonNil[T any](nilable ...*T) *T {
	for _, v := range nilable {
		if v != nil {
			return v
		}
	}
	return nil
}
