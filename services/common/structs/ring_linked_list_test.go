package structs_test

import (
	"service/services/common/structs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// should point back to head after complete iteration
func TestRingLinkedList(t *testing.T) {
	rll := new(structs.RingLinkedList)
	rll.Init(2)

	rll.Put("a")
	rll.Put("b")
	rll.Put("c")

	resultingArray := rll.Collect()
	assert.Len(t, resultingArray, 2)
	assert.Equal(t, resultingArray[0], "b")
	assert.Equal(t, resultingArray[1], "c")
}

func TestRingLinkedList2(t *testing.T) {
	rll := new(structs.RingLinkedList)
	rll.Init(2)

	rll.Put("a")
	rll.Put("b")
	rll.Put("c")
	rll.Put("d")

	resultingArray := rll.Collect()
	assert.Len(t, resultingArray, 2)
	assert.Equal(t, resultingArray[0], "c")
	assert.Equal(t, resultingArray[1], "d")
}

func TestRingLinkedListNonNil(t *testing.T) {
	rll := new(structs.RingLinkedList)
	rll.Init(2)

	rll.Put("a")
	rll.Put("b")
	rll.Put(nil)
	rll.Put("d")

	resultingArray := rll.CollectNonNil()
	assert.Len(t, resultingArray, 1)
	assert.Equal(t, resultingArray[0], "d")
}
