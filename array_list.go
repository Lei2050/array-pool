package arraypool

import (
	"math/bits"
)

const (
	// DefaultSegmentSize is the default size of each segment in the ArrayList.
	DefaultSegmentSize = 128
)

type Segment[T any] struct {
	arr   []T
	count int
}

// ArrayList is a generic data structure that stores elements in segments.
// It uses a segmented array approach to optimize memory usage and performance.
// Each segment has a fixed size, and the ArrayList automatically allocates new segments as needed.
//
// Type parameters:
// T - The type of elements stored in the ArrayList.
type ArrayList[T any] struct {
	// defaultValue is the default value of type T.
	// It is used to reset elements when they are removed.
	defaultValue T
	// count is the total number of elements currently stored in the ArrayList.
	count int
	// segments is a slice of Segment[T] that stores the actual elements.
	// Each segment has a fixed size defined by segmentSize.
	segments []Segment[T]
	// segmentIdx is the index of the current segment in the segments slice.
	// It is used to keep track of the current segment when adding or removing elements.
	segmentIdx int
	// segmentSize is the size of each segment in the ArrayList.
	// It is rounded up to the nearest power of 2 to optimize bitwise operations.
	segmentSize int
	// segmentSizeMask is a mask used to calculate the index within a segment.
	// It is equal to segmentSize - 1.
	segmentSizeMask int
	// segmentSizeShift is the number of bits to shift to calculate the segment index.
	// It is equal to the number of trailing zeros in segmentSize.
	segmentSizeShift int
}

// NewArrayList creates a new instance of ArrayList with the specified segment size.
// The segment size is rounded up to the nearest power of 2 to optimize bitwise operations.
//
// Parameters:
// segmentSize - The initial size of each segment in the ArrayList. Must be greater than 0.
//
// Returns:
// A pointer to the newly created ArrayList instance.
//
// Note:
// This function will panic if the segmentSize is less than or equal to 0.
func NewArrayList[T any](segmentSize int) *ArrayList[T] {
	if segmentSize <= 0 {
		segmentSize = DefaultSegmentSize
	}

	al := &ArrayList[T]{
		// Round the segment size up to the nearest power of 2
		segmentSize: int(NearestPowerOf2(uint(segmentSize))),
	}
	al.segmentSizeMask = al.segmentSize - 1
	al.segmentSizeShift = bits.TrailingZeros(uint(al.segmentSize))
	return al
}

// Add adds a new element to the ArrayList and returns a pointer to the newly added element.
// This method first allocates a new slot in the ArrayList using the Alloc method.
// Then it assigns the provided value to the allocated slot.
//
// Parameters:
// v - The value of type T to be added to the ArrayList.
//
// Returns:
// A pointer to the newly added element of type T.
func (al *ArrayList[T]) Add(v T) *T {
	ptr := al.Alloc()
	*ptr = v
	return ptr
}

// Alloc allocates a new slot in the ArrayList and returns a pointer to the allocated element.
// This method first checks if there are any segments in the list. If not, it creates a new segment.
// If the current segment is full, it either creates a new segment or moves to the next one.
// After ensuring there is space, it increments the count of the current segment and the total count of the list.
// Finally, it returns a pointer to the newly allocated element in the current segment.
//
// Returns:
// A pointer to the newly allocated element of type T.
func (al *ArrayList[T]) Alloc() *T {
	if len(al.segments) == 0 {
		// If there are no segments, create a new segment with the specified segment size
		al.segments = append(al.segments, Segment[T]{arr: make([]T, al.segmentSize)})
	} else {
		// Check if the current segment is full
		if al.segments[al.segmentIdx].count == al.segmentSize {
			if al.segmentIdx == len(al.segments)-1 {
				// If it is the last segment, create a new segment with the specified segment size
				al.segments = append(al.segments, Segment[T]{arr: make([]T, al.segmentSize)})
			}
			al.segmentIdx++
		}
	}

	segment := &al.segments[al.segmentIdx]
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
	lastPtr := al.GetRef(al.count - 1)
	// Replace the element at the specified index with the last element
	*al.GetRef(idx) = *lastPtr
	// Set the last element to the sentinel value
	*lastPtr = al.defaultValue
	// Remove the last element from the list
	al.RemoveLast()
}

// RemoveLast removes the last element from the ArrayList.
// This method checks if the list is empty. If not, it decrements the count of the current segment
// and the total count of the list. If the current segment becomes empty, it tries to move to the
// previous segment and releases the last empty segment if necessary.
//
// Note:
// This method does nothing if the list is already empty.
func (al *ArrayList[T]) RemoveLast() {
	if al.count <= 0 {
		return
	}

	segment := &al.segments[al.segmentIdx]
	segment.count--
	al.count--
	if segment.count > 0 {
		return
	}

	lastEmpty := al.segmentIdx
	// al.segmentIdx is always non-negative
	if al.segmentIdx > 0 {
		al.segmentIdx--
	}
	if len(al.segments)-1 == lastEmpty {
		// We try to keep an empty segment to prevent frequent segment allocation.
		return
	}

	// Move to the next segment (which is the last empty segment)
	lastEmpty++
	// Release the last empty segment by setting its array to nil and count to 0
	al.segments[lastEmpty].arr = nil
	al.segments[lastEmpty].count = 0
	// Remove the last empty segment from the list
	al.segments = al.segments[0:lastEmpty]
}

func (al *ArrayList[T]) Clear() {
	for i := range al.segments {
		al.segments[i].arr = nil
		al.segments[i].count = 0
	}
	al.segments = nil
	al.count = 0
	al.segmentIdx = 0
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
