package utils

import "container/list"

func ListToArray[T any](l *list.List) []T {
	var arr []T = make([]T, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		arr[i] = *(e.Value.(*T))
		i += 1
	}
	return arr
}
