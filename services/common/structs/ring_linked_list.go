package structs

import (
	"service/services/common/interfaces"

	"github.com/mariomac/gostream/stream"
)

type RingLinkedListItem struct {
	next  *RingLinkedListItem
	value any
}

type RingLinkedList struct {
	Size int
	Head *RingLinkedListItem
}

func (rll *RingLinkedList) Init(size int) {
	if size < 1 {
		rll.Size = 0
		rll.Head = nil
		return
	}

	rll.Size = size
	current := new(RingLinkedListItem)
	head := current
	rll.Head = current
	iterate := 1
	for iterate < size {
		new := new(RingLinkedListItem)
		current.next = new
		current = new
		iterate += 1
	}
	current.next = head
}

func (rll *RingLinkedList) Put(item any) {
	if rll.Size < 1 {
		return
	}
	next := rll.Head.next
	next.value = item
	rll.Head = next
}

func (rll *RingLinkedList) Collect() []any {
	arr := make([]any, rll.Size)
	// always start from next, not from current head
	// we want the current head to be placed at the last of the array
	next := rll.Head.next
	iterate := 0
	for iterate < rll.Size {
		arr[iterate] = next.value
		next = next.next
		iterate += 1
	}
	return arr
}

func (rll *RingLinkedList) CollectNonNil() []any {
	return stream.
		OfSlice(rll.Collect()).
		Filter(func(a any) bool {
			return a != nil
		}).ToSlice()
}

func (rll *RingLinkedList) CollectNonNilToString() []string {
	nonNilStream := stream.
		OfSlice(rll.Collect()).
		Filter(func(a any) bool {
			return a != nil
		})
	return stream.Map(nonNilStream, func(a any) string {
		if stringable, isStringable := a.(interfaces.Stringable); isStringable {
			return stringable.String()
		} else {
			return ""
		}
	}).ToSlice()
}
