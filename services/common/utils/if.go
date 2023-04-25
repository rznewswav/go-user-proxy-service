package utils

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func Iff[T any](cond bool, vtrue, vfalse func() T) T {
	if cond {
		return vtrue()
	}
	return vfalse()
}
