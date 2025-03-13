package arraypool

import (
	"math/bits"
)

type Segment[T any] struct {
	arr   []T
	count int
}

type ArrayList[T any] struct {
	sentinel         T //哨兵
	count            int
	segments         []Segment[T]
	segmentSize      int
	segmentSizeMask  int
	segmentSizeShift int
}

func NewArrayList[T any](segmentSize int) *ArrayList[T] {
	if segmentSize <= 0 {
		panic("segmentSize must be greater than 0")
	}

	al := &ArrayList[T]{
		segmentSize: int(NearestPowerOf2(uint(segmentSize))),
	}
	al.segmentSizeMask = al.segmentSize - 1
	al.segmentSizeShift = bits.TrailingZeros(uint(al.segmentSize))
	return al
}

func (al *ArrayList[T]) Add(v T) *T {
	ptr := al.Alloc()
	*ptr = v
	return ptr
}

func (al *ArrayList[T]) Alloc() *T {
	if len(al.segments) == 0 || al.segments[len(al.segments)-1].count == al.segmentSize {
		al.segments = append(al.segments, Segment[T]{arr: make([]T, al.segmentSize)})
	}

	segment := &al.segments[len(al.segments)-1]
	segment.count++
	al.count++
	return &segment.arr[segment.count-1]
}

func (al *ArrayList[T]) Get(idx int) T {
	if idx >= al.count {
		panic("out of bound")
	}
	return al.segments[idx>>al.segmentSizeShift].arr[idx&al.segmentSizeMask]
}

func (al *ArrayList[T]) GetRef(idx int) *T {
	if idx >= al.count {
		panic("out of bound")
	}
	return &al.segments[idx>>al.segmentSizeShift].arr[idx&al.segmentSizeMask]
}

// FastRemoveAt removes the element at the specified index in the ArrayList.
// This method replaces the element at the given index with the last element in the list,
// and then removes the last element from the list.
// This operation is fast because it does not require shifting all the elements after the removed element.
//
// Parameters:
// idx - The index of the element to be removed.
//
// Note:
// This method will panic if the index is out of bounds.
func (al *ArrayList[T]) FastRemoveAt(idx int) {
	// Get a reference to the last element in the list
	// Bugfix: Here should be `al.count - 1` to get the correct last element
	lastPtr := al.GetRef(al.count - 1)
	// Replace the element at the specified index with the last element
	*al.GetRef(idx) = *lastPtr
	// Set the last element to the sentinel value
	*lastPtr = al.sentinel
	// Remove the last element from the list
	al.RemoveLast()
}

func (al *ArrayList[T]) RemoveLast() {
	if al.count <= 0 {
		return
	}

	segment := &al.segments[len(al.segments)-1]
	segment.count--
	al.count--
	if segment.count > 0 {
		return
	}

	//release last segment
	segment.arr = nil
	al.segments = al.segments[0 : len(al.segments)-1]
}

func (al *ArrayList[T]) Clear() {
	for i := range al.segments {
		al.segments[i].arr = nil
		al.segments[i].count = 0
	}
	al.segments = nil
	al.count = 0
}

func (al *ArrayList[T]) Count() int {
	return al.count
}

func NearestPowerOf2(n uint) uint {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
